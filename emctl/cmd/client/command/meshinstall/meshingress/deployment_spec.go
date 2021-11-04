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

package meshingress

import (
	"github.com/megaease/easemeshctl/cmd/client/command/flags"
	installbase "github.com/megaease/easemeshctl/cmd/client/command/meshinstall/base"

	"github.com/pkg/errors"
	appsV1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type deploymentSpecFunc func(*flags.Install) *appsV1.Deployment

func meshIngressLabel() map[string]string {
	selector := map[string]string{}
	selector["app"] = "Easegress-ingress"
	return selector
}

func deploymentSpec(ctx *installbase.StageContext) installbase.InstallFunc {
	deployment := deploymentConfigVolumeSpec(
		deploymentContainerSpec(
			deploymentBaseSpec(
				deploymentInitialize(nil))))(ctx.Flags)

	return func(ctx *installbase.StageContext) error {
		err := installbase.DeployDeployment(deployment, ctx.Client, ctx.Flags.MeshNamespace)
		if err != nil {
			return errors.Wrapf(err, "deployment operation %s failed", deployment.Name)
		}
		return err
	}
}

func deploymentInitialize(fn deploymentSpecFunc) deploymentSpecFunc {
	return func(installFlags *flags.Install) *appsV1.Deployment {
		return &appsV1.Deployment{}
	}
}

func deploymentBaseSpec(fn deploymentSpecFunc) deploymentSpecFunc {
	return func(installFlags *flags.Install) *appsV1.Deployment {
		spec := fn(installFlags)
		spec.Name = installbase.DefaultMeshIngressControllerName
		spec.Spec.Selector = &metav1.LabelSelector{
			MatchLabels: meshIngressLabel(),
		}

		var replicas = int32(installFlags.MeshIngressReplicas)
		spec.Spec.Replicas = &replicas
		spec.Spec.Template.Labels = meshIngressLabel()
		spec.Spec.Template.Spec.Containers = []v1.Container{}
		return spec
	}
}

func deploymentContainerSpec(fn deploymentSpecFunc) deploymentSpecFunc {

	return func(installFlags *flags.Install) *appsV1.Deployment {

		spec := fn(installFlags)
		container, _ := installbase.AcceptContainerVisitor("easegress-ingress",
			installFlags.ImageRegistryURL+"/"+installFlags.EasegressImage,
			v1.PullIfNotPresent,
			newVisitor(installFlags))

		spec.Spec.Template.Spec.Containers = append(spec.Spec.Template.Spec.Containers, *container)
		return spec
	}
}

func deploymentConfigVolumeSpec(fn deploymentSpecFunc) deploymentSpecFunc {

	return func(installFlags *flags.Install) *appsV1.Deployment {
		spec := fn(installFlags)
		spec.Spec.Template.Spec.Volumes = []v1.Volume{
			{
				Name: "eg-ingress-config",
				VolumeSource: v1.VolumeSource{
					ConfigMap: &v1.ConfigMapVolumeSource{
						LocalObjectReference: v1.LocalObjectReference{
							Name: installbase.DefaultMeshIngressConfig,
						},
					},
				},
			},
		}
		return spec
	}
}

type containerVisitor struct {
	installFlags *flags.Install
}

func newVisitor(installFlags *flags.Install) installbase.ContainerVisitor {
	return &containerVisitor{installFlags}
}

func (v *containerVisitor) VisitorCommandAndArgs(c *v1.Container) (command []string, installFlags []string) {

	return []string{"/bin/sh"},
		[]string{"-c", "/opt/easegress/bin/easegress-server -f /opt/eg-config/eg-ingress.yaml"}
}

func (v *containerVisitor) VisitorContainerPorts(c *v1.Container) ([]v1.ContainerPort, error) {

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

func (v *containerVisitor) VisitorEnvs(c *v1.Container) ([]v1.EnvVar, error) {

	return nil, nil
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
			Name:      "eg-ingress-config",
			MountPath: "/opt/eg-config/eg-ingress.yaml",
			SubPath:   "eg-ingress.yaml",
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
				Path: "/apis/v1/healthz",
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
