/*
 * Copyright (c) 2021, MegaEase
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

package ingresscontroller

import (
	"github.com/megaease/easemeshctl/cmd/client/command/flags"
	installbase "github.com/megaease/easemeshctl/cmd/client/command/meshinstall/base"

	"github.com/pkg/errors"
	appsV1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type deploymentSpecFunc func(*installbase.StageContext) *appsV1.Deployment

func meshIngressLabel() map[string]string {
	selector := map[string]string{}
	selector["app"] = installbase.IngressControllerDeploymentName
	return selector
}

func deploymentSpec(ctx *installbase.StageContext) installbase.InstallFunc {
	deployment := deploymentConfigVolumeSpec(
		deploymentContainerSpec(
			deploymentBaseSpec(
				deploymentInitialize(nil))))(ctx)

	return func(ctx *installbase.StageContext) error {
		err := installbase.DeployDeployment(deployment, ctx.Client, ctx.Flags.MeshNamespace)
		if err != nil {
			return errors.Wrapf(err, "deploy %s failed", deployment.Name)
		}
		return err
	}
}

func deploymentInitialize(fn deploymentSpecFunc) deploymentSpecFunc {
	return func(ctx *installbase.StageContext) *appsV1.Deployment {
		return &appsV1.Deployment{}
	}
}

func deploymentBaseSpec(fn deploymentSpecFunc) deploymentSpecFunc {
	return func(ctx *installbase.StageContext) *appsV1.Deployment {
		spec := fn(ctx)
		spec.Name = installbase.IngressControllerDeploymentName
		spec.Spec.Selector = &metav1.LabelSelector{
			MatchLabels: meshIngressLabel(),
		}

		replicas := int32(ctx.Flags.MeshIngressReplicas)
		spec.Spec.Replicas = &replicas
		spec.Spec.Template.Labels = meshIngressLabel()
		spec.Spec.Template.Spec.Containers = []v1.Container{}
		return spec
	}
}

func deploymentContainerSpec(fn deploymentSpecFunc) deploymentSpecFunc {
	return func(ctx *installbase.StageContext) *appsV1.Deployment {
		spec := fn(ctx)
		container, _ := installbase.AcceptContainerVisitor(installbase.IngressControllerDeploymentName,
			ctx.Flags.ImageRegistryURL+"/"+ctx.Flags.EasegressImage,
			v1.PullIfNotPresent,
			newVisitor(ctx))

		spec.Spec.Template.Spec.Containers = append(spec.Spec.Template.Spec.Containers, *container)
		return spec
	}
}

func deploymentConfigVolumeSpec(fn deploymentSpecFunc) deploymentSpecFunc {
	return func(ctx *installbase.StageContext) *appsV1.Deployment {
		spec := fn(ctx)
		spec.Spec.Template.Spec.Volumes = []v1.Volume{
			{
				Name: installbase.IngressControllerConfigMapName,
				VolumeSource: v1.VolumeSource{
					ConfigMap: &v1.ConfigMapVolumeSource{
						LocalObjectReference: v1.LocalObjectReference{
							Name: installbase.IngressControllerConfigMapName,
						},
					},
				},
			},
		}
		return spec
	}
}

type containerVisitor struct {
	ctx *installbase.StageContext
}

func newVisitor(ctx *installbase.StageContext) installbase.ContainerVisitor {
	return &containerVisitor{ctx}
}

func (v *containerVisitor) VisitorCommandAndArgs(c *v1.Container) (command []string, args []string) {
	return []string{"/bin/sh"},
		[]string{"-c", installbase.IngressControllerDeploymentCmd}
}

func (v *containerVisitor) VisitorContainerPorts(c *v1.Container) ([]v1.ContainerPort, error) {
	return []v1.ContainerPort{
		{
			Name:          installbase.ControlPlaneStatefulSetAdminPortName,
			ContainerPort: flags.DefaultMeshAdminPort,
		},
		{
			Name:          installbase.ControlPlaneStatefulSetClientPortName,
			ContainerPort: flags.DefaultMeshClientPort,
		},
		{
			Name:          installbase.ControlPlaneStatefulSetPeerPortName,
			ContainerPort: flags.DefaultMeshPeerPort,
		},
	}, nil
}

func (v *containerVisitor) VisitorEnvs(c *v1.Container) ([]v1.EnvVar, error) {
	return []v1.EnvVar{
		{
			Name: "EG_NAME",
			ValueFrom: &v1.EnvVarSource{
				FieldRef: &v1.ObjectFieldSelector{
					FieldPath: "metadata.name",
				},
			},
		},
		{
			Name: "HOSTNAME",
			ValueFrom: &v1.EnvVarSource{
				FieldRef: &v1.ObjectFieldSelector{
					FieldPath: "metadata.name",
				},
			},
		},
		{
			Name: "APPLICATION_IP",
			ValueFrom: &v1.EnvVarSource{
				FieldRef: &v1.ObjectFieldSelector{
					FieldPath: "status.podIP",
				},
			},
		},
	}, nil
}

func (v *containerVisitor) VisitorEnvFrom(c *v1.Container) ([]v1.EnvFromSource, error) {
	return nil, nil
}

func (v *containerVisitor) VisitorResourceRequirements(c *v1.Container) (*v1.ResourceRequirements, error) {
	return nil, nil
}

func (v *containerVisitor) VisitorVolumeMounts(c *v1.Container) ([]v1.VolumeMount, error) {
	return []v1.VolumeMount{
		{
			Name:      installbase.IngressControllerConfigMapName,
			MountPath: installbase.IngressControllerConfigMapVolumeMountPath,
			SubPath:   installbase.IngressControllerConfigMapVolumeMountSubPath,
		},
	}, nil
}

func (v *containerVisitor) VisitorVolumeDevices(c *v1.Container) ([]v1.VolumeDevice, error) {
	return nil, nil
}

func (v *containerVisitor) VisitorLivenessProbe(c *v1.Container) (*v1.Probe, error) {
	/* FIXME: K8s probe report connection reset, but the port can be accessed via localhost/127.0.0.1
	maybe the default admin API port should listen on all interface instead of loopback address.

	return &v1.Probe{
		Handler: v1.Handler{
			HTTPGet: &v1.HTTPGetAction{
				Host: "localhost",
				Port: intstr.FromInt(installbase.DefaultMeshAdminPort),
				Path: "/apis/v2/healthz",
			},
		},
		InitialDelaySeconds: 50,
	}, nil
	*/
	return nil, nil
}

func (v *containerVisitor) VisitorReadinessProbe(c *v1.Container) (*v1.Probe, error) {
	return nil, nil
}

func (v *containerVisitor) VisitorLifeCycle(c *v1.Container) (*v1.Lifecycle, error) {
	return nil, nil
}

func (v *containerVisitor) VisitorSecurityContext(c *v1.Container) (*v1.SecurityContext, error) {
	return nil, nil
}
