/*
 * Copyright (c) 2017, MegaEase
 * All rights reserved.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package controlpanel

import (
	"fmt"

	"github.com/megaease/easemeshctl/cmd/client/command/flags"
	installbase "github.com/megaease/easemeshctl/cmd/client/command/meshinstall/base"
	"github.com/megaease/easemeshctl/cmd/common"
	"github.com/pkg/errors"

	appsV1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type statefulsetSpecFunc func(ctx *installbase.StageContext) *appsV1.StatefulSet

func statefulsetSpec(ctx *installbase.StageContext) installbase.InstallFunc {
	statefulSet := statefulsetPVCSpec(
		statefulsetContainerSpec(
			baseStatefulSetSpec(
				initialStatefulSetSpec(nil))))(ctx)

	return func(ctx *installbase.StageContext) error {
		err := installbase.DeployStatefulset(statefulSet, ctx.Client, ctx.Flags.MeshNamespace)
		if err != nil {
			return errors.Wrapf(err, "deploy statefulset %s failed", statefulSet.ObjectMeta.Name)
		}
		return nil
	}
}

func initialStatefulSetSpec(fn statefulsetSpecFunc) statefulsetSpecFunc {
	return func(ctx *installbase.StageContext) *appsV1.StatefulSet {
		return &appsV1.StatefulSet{}
	}
}

func baseStatefulSetSpec(fn statefulsetSpecFunc) statefulsetSpecFunc {
	return func(ctx *installbase.StageContext) *appsV1.StatefulSet {
		spec := fn(ctx)
		labels := meshControlPlaneLabel()
		spec.Name = installbase.DefaultMeshControlPlaneName
		spec.Spec.ServiceName = installbase.DefaultMeshControlPlaneHeadlessServiceName

		spec.Spec.Selector = &metav1.LabelSelector{
			MatchLabels: labels,
		}

		var replicas = int32(ctx.Flags.EasegressControlPlaneReplicas)
		spec.Spec.Replicas = &replicas
		spec.Spec.Template.Labels = labels
		spec.Spec.Template.Spec.Volumes = []v1.Volume{
			{
				Name: installbase.DefaultMeshControlPlaneConfig,
				VolumeSource: v1.VolumeSource{
					ConfigMap: &v1.ConfigMapVolumeSource{
						LocalObjectReference: v1.LocalObjectReference{
							Name: installbase.DefaultMeshControlPlaneConfig,
						},
					},
				},
			},
		}
		return spec
	}
}
func statefulsetPVCSpec(fn statefulsetSpecFunc) statefulsetSpecFunc {
	return func(ctx *installbase.StageContext) *appsV1.StatefulSet {
		spec := fn(ctx)
		pvc := v1.PersistentVolumeClaim{}
		pvc.Name = installbase.DefaultMeshControlPlanePVName
		pvc.Spec.AccessModes = []v1.PersistentVolumeAccessMode{v1.ReadWriteOnce}
		pvc.Spec.StorageClassName = &ctx.Flags.MeshControlPlaneStorageClassName

		pvc.Spec.Resources.Requests = v1.ResourceList{
			v1.ResourceStorage: resource.MustParse(ctx.Flags.MeshControlPlanePersistVolumeCapacity),
		}
		spec.Spec.VolumeClaimTemplates = []v1.PersistentVolumeClaim{pvc}
		return spec
	}
}

func statefulsetContainerSpec(fn statefulsetSpecFunc) statefulsetSpecFunc {
	return func(ctx *installbase.StageContext) *appsV1.StatefulSet {
		spec := fn(ctx)
		container, err := installbase.AcceptContainerVisitor("easegress",
			ctx.Flags.ImageRegistryURL+"/"+ctx.Flags.EasegressImage,
			v1.PullIfNotPresent,
			newContainerVisistor(ctx))
		if err != nil {
			common.ExitWithErrorf("generate mesh controlpanel container spec failed: %s", err)
			return nil
		}

		spec.Spec.Template.Spec.Containers = []v1.Container{*container}
		return spec
	}
}

type containerVisitor struct {
	ctx *installbase.StageContext
}

var _ installbase.ContainerVisitor = &containerVisitor{}

func (m *containerVisitor) VisitorCommandAndArgs(c *v1.Container) (command []string, args []string) {
	return []string{"/bin/sh"},
		[]string{
			"-c",
			"/opt/easegress/bin/easegress-server -f /opt/eg-config/eg-master.yaml"}
}

func (m *containerVisitor) VisitorContainerPorts(c *v1.Container) ([]v1.ContainerPort, error) {
	return []v1.ContainerPort{
		{
			Name:          installbase.DefaultMeshAdminPortName,
			ContainerPort: flags.DefaultMeshAdminPort,
		},
		{
			Name:          installbase.DefaultMeshClientPortName,
			ContainerPort: flags.DefaultMeshClientPort,
		},
		{
			Name:          installbase.DefaultMeshPeerPortName,
			ContainerPort: flags.DefaultMeshPeerPort,
		},
	}, nil
}

func (m *containerVisitor) VisitorEnvs(c *v1.Container) ([]v1.EnvVar, error) {
	return []v1.EnvVar{
		{
			Name: "EG_NAME",
			ValueFrom: &v1.EnvVarSource{
				FieldRef: &v1.ObjectFieldSelector{
					FieldPath: "metadata.name",
				},
			},
		},
		// We set a unreachable host to --advertise-clients-urls and
		// initial-advertise-peer-urls as we need a consistency configuration
		// for all Easegress instance. The real cluster-advertise-client-url
		// and cluster-initial-peer-url will be passed through environment
		// `EG_CLUSTER_ADVERTISE_CLIENT_URL` and `EG_CLUSTER_INITIAL_ADVERTISE_PEER_URLS`
		{
			// Kubernetes leverage shell syntax to help refering another environment
			Name:  "EG_CLUSTER_ADVERTISE_CLIENT_URLS",
			Value: fmt.Sprintf("http://$(EG_NAME).%s.%s:%d", installbase.DefaultMeshControlPlaneHeadlessServiceName, m.ctx.Flags.MeshNamespace, m.ctx.Flags.EgClientPort),
		},
		{
			Name:  "EG_CLUSTER_INITIAL_ADVERTISE_PEER_URLS",
			Value: fmt.Sprintf("http://$(EG_NAME).%s.%s:%d", installbase.DefaultMeshControlPlaneHeadlessServiceName, m.ctx.Flags.MeshNamespace, m.ctx.Flags.EgPeerPort),
		},
	}, nil
}

func (m *containerVisitor) VisitorEnvFrom(c *v1.Container) ([]v1.EnvFromSource, error) {
	// do nothing
	return nil, nil
}

func (m *containerVisitor) VisitorResourceRequirements(c *v1.Container) (*v1.ResourceRequirements, error) {
	cpuRequest, err := resource.ParseQuantity("100m")
	if err != nil {
		return nil, err
	}
	memoryRequest, err := resource.ParseQuantity("1Gi")
	if err != nil {
		return nil, err
	}

	cpuLimit, err := resource.ParseQuantity("1000m")
	if err != nil {
		return nil, err
	}
	memoryLimit, err := resource.ParseQuantity("2Gi")
	if err != nil {
		return nil, err
	}

	return &v1.ResourceRequirements{
		Requests: map[v1.ResourceName]resource.Quantity{
			v1.ResourceCPU:    cpuRequest,
			v1.ResourceMemory: memoryRequest,
		},
		Limits: map[v1.ResourceName]resource.Quantity{
			v1.ResourceCPU:    cpuLimit,
			v1.ResourceMemory: memoryLimit,
		},
	}, nil
}

func (m *containerVisitor) VisitorVolumeMounts(c *v1.Container) ([]v1.VolumeMount, error) {
	return []v1.VolumeMount{
		{
			Name:      installbase.DefaultMeshControlPlanePVName,
			MountPath: "/opt/eg-data/",
		},
		{
			Name:      installbase.DefaultMeshControlPlaneConfig,
			MountPath: "/opt/eg-config/eg-master.yaml",
			SubPath:   "eg-master.yaml",
		},
	}, nil
}

func (m *containerVisitor) VisitorVolumeDevices(c *v1.Container) ([]v1.VolumeDevice, error) {
	// do nothing
	return nil, nil
}

func (m *containerVisitor) VisitorLivenessProbe(c *v1.Container) (*v1.Probe, error) {
	// do nothing
	return nil, nil
}

func (m *containerVisitor) VisitorReadinessProbe(c *v1.Container) (*v1.Probe, error) {

	// The initialization of the etcd's cluster depended on the domain name,
	// but domain name register rely on pod ready status, and pod ready
	// status rely on the successful initialization of etcd's cluster.
	// The situation produces a cycle dependency, so we disabled K8s
	// readiness probe

	// return &v1.Probe{
	// 	Handler: v1.Handler{
	// 		HTTPGet: &v1.HTTPGetAction{
	// 			Host: "127.0.0.1",
	// 			Port: intstr.FromInt(m.ctx.Flags.EgAdminPort),
	// 			Path: "/apis/v1/healthz",
	// 		},
	// 	},
	// 	InitialDelaySeconds: 10,
	// }, nil
	return nil, nil
}

func (m *containerVisitor) VisitorLifeCycle(c *v1.Container) (*v1.Lifecycle, error) {
	// do nothing
	return nil, nil
}

func (m *containerVisitor) VisitorSecurityContext(c *v1.Container) (*v1.SecurityContext, error) {
	// do nothing
	return nil, nil
}

func newContainerVisistor(ctx *installbase.StageContext) installbase.ContainerVisitor {
	return &containerVisitor{ctx: ctx}
}
