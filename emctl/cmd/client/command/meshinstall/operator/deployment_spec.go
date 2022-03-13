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

package operator

import (
	installbase "github.com/megaease/easemeshctl/cmd/client/command/meshinstall/base"

	"github.com/pkg/errors"
	appsV1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

type deploymentSpecFunc func(ctx *installbase.StageContext) *appsV1.Deployment

func operatorDeploymentSpec(ctx *installbase.StageContext) installbase.InstallFunc {
	deployment := deploymentConfigVolumeSpec(
		deploymentManagerContainerSpec(
			deploymentRBACContainerSpec(
				deploymentBaseSpec(deploymentInitialize(nil)))))(ctx)

	return func(ctx *installbase.StageContext) error {
		err := installbase.DeployDeployment(deployment, ctx.Client, ctx.Flags.MeshNamespace)
		if err != nil {
			return errors.Wrapf(err, "deployment operation %s failed", deployment.Name)
		}
		return err
	}
}

func meshOperatorLabels() map[string]string {
	selector := map[string]string{}
	selector["app"] = installbase.OperatorDeploymentName
	return selector
}

func deploymentInitialize(fn deploymentSpecFunc) deploymentSpecFunc {
	return func(ctx *installbase.StageContext) *appsV1.Deployment {
		return &appsV1.Deployment{}
	}
}

func deploymentBaseSpec(fn deploymentSpecFunc) deploymentSpecFunc {
	return func(ctx *installbase.StageContext) *appsV1.Deployment {
		spec := fn(ctx)

		labels := meshOperatorLabels()
		spec.Name = installbase.OperatorDeploymentName
		spec.Spec.Selector = &metav1.LabelSelector{
			MatchLabels: labels,
		}

		replicas := int32(ctx.Flags.EaseMeshOperatorReplicas)
		spec.Spec.Replicas = &replicas
		spec.Spec.Template.Labels = labels
		spec.Spec.Template.Spec.Containers = []v1.Container{}

		var v int64 = 65532 //?
		spec.Spec.Template.Spec.SecurityContext = &v1.PodSecurityContext{
			RunAsUser: &v,
		}
		return spec
	}
}

func deploymentRBACContainerSpec(fn deploymentSpecFunc) deploymentSpecFunc {
	return func(ctx *installbase.StageContext) *appsV1.Deployment {
		spec := fn(ctx)
		rbacContainer := v1.Container{}
		rbacContainer.Name = "kube-rbac-proxy"
		rbacContainer.Image = "gcr.io/kubebuilder/kube-rbac-proxy:v0.5.0"
		rbacContainer.Ports = []v1.ContainerPort{
			{
				Name:          "https",
				ContainerPort: 8443,
			},
		}
		rbacContainer.Args = []string{
			"--secure-listen-address=0.0.0.0:8443",
			"--upstream=http://127.0.0.1:8080/",
			"--logtostderr=true",
			"--v=10",
		}
		spec.Spec.Template.Spec.Containers = append(spec.Spec.Template.Spec.Containers, rbacContainer)
		return spec
	}
}

func deploymentConfigVolumeSpec(fn deploymentSpecFunc) deploymentSpecFunc {
	return func(ctx *installbase.StageContext) *appsV1.Deployment {
		spec := fn(ctx)
		spec.Spec.Template.Spec.Volumes = []v1.Volume{
			{
				Name: installbase.OperatorConfigMapName,
				VolumeSource: v1.VolumeSource{
					ConfigMap: &v1.ConfigMapVolumeSource{
						LocalObjectReference: v1.LocalObjectReference{
							Name: installbase.OperatorConfigMapName,
						},
					},
				},
			},
			{
				Name: installbase.OperatorSecretName,
				VolumeSource: v1.VolumeSource{
					Secret: &v1.SecretVolumeSource{
						SecretName: installbase.OperatorSecretName,
					},
				},
			},
		}
		return spec
	}
}

func deploymentManagerContainerSpec(fn deploymentSpecFunc) deploymentSpecFunc {
	return func(ctx *installbase.StageContext) *appsV1.Deployment {
		spec := fn(ctx)
		container, _ := installbase.AcceptContainerVisitor("operator-manager",
			ctx.Flags.ImageRegistryURL+"/"+ctx.Flags.EaseMeshOperatorImage,
			v1.PullIfNotPresent,
			newVisitor(ctx))

		spec.Spec.Template.Spec.Containers = append(spec.Spec.Template.Spec.Containers, *container)
		return spec
	}
}

func newVisitor(ctx *installbase.StageContext) installbase.ContainerVisitor {
	return &containerVisitor{ctx: ctx}
}

type containerVisitor struct {
	ctx *installbase.StageContext
}

func (v *containerVisitor) VisitorCommandAndArgs(c *v1.Container) (command []string, args []string) {
	return []string{installbase.OperatorCmd}, []string{installbase.OperatorArgs}
}

func (v *containerVisitor) VisitorContainerPorts(c *v1.Container) ([]v1.ContainerPort, error) {
	return []v1.ContainerPort{
		{
			Name:          installbase.OperatorMutatingWebhookPortName,
			ContainerPort: installbase.OperatorMutatingWebhookPort,
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

func (v *containerVisitor) VisitorVolumeMounts(c *v1.Container) ([]v1.VolumeMount, error) {
	return []v1.VolumeMount{
		{
			Name:      installbase.OperatorConfigMapName,
			MountPath: installbase.OperatorConfigMapVolumeMountPath,
			SubPath:   installbase.OperatorConfigMapVolumeMountSubPath,
		},
		{
			Name:      installbase.OperatorSecretName,
			MountPath: installbase.OperatorSecretVolumeMountPath,
		},
	}, nil
}

func (v *containerVisitor) VisitorVolumeDevices(c *v1.Container) ([]v1.VolumeDevice, error) {
	return nil, nil
}

func (v *containerVisitor) VisitorLivenessProbe(c *v1.Container) (*v1.Probe, error) {
	return &v1.Probe{
		Handler: v1.Handler{
			HTTPGet: &v1.HTTPGetAction{
				Path:   "/healthz",
				Port:   intstr.FromInt(8081),
				Scheme: "HTTP",
			},
		},
		InitialDelaySeconds: 15,
		PeriodSeconds:       20,
	}, nil
}

func (v *containerVisitor) VisitorReadinessProbe(c *v1.Container) (*v1.Probe, error) {
	return &v1.Probe{
		Handler: v1.Handler{
			HTTPGet: &v1.HTTPGetAction{
				Path:   "/readyz",
				Port:   intstr.FromInt(8081),
				Scheme: "HTTP",
			},
		},
		InitialDelaySeconds: 5,
		PeriodSeconds:       10,
	}, nil
}

func (v *containerVisitor) VisitorLifeCycle(c *v1.Container) (*v1.Lifecycle, error) {
	return nil, nil
}

func (v *containerVisitor) VisitorSecurityContext(c *v1.Container) (*v1.SecurityContext, error) {
	return nil, nil
}
