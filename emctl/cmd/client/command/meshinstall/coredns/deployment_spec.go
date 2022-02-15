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

package coredns

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

func coreDNSDeploymentSpec(ctx *installbase.StageContext) installbase.InstallFunc {
	deployment := deploymentConfigVolumeSpec(
		deploymentBaseSpec(deploymentInitialize(nil)))(ctx)

	return func(ctx *installbase.StageContext) error {
		err := installbase.DeployDeployment(deployment, ctx.Client, coreDNSNamespace)
		if err != nil {
			return errors.Wrapf(err, "deployment operation %s failed", deployment.Name)
		}
		return err
	}
}

func coreDNSLabels() map[string]string {
	selector := map[string]string{}
	selector["k8s-app"] = "kube-dns"
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

		labels := coreDNSLabels()
		spec.Name = "coredns"
		spec.Labels = labels
		spec.Spec.Selector = &metav1.LabelSelector{
			MatchLabels: labels,
		}
		replicas := int32(ctx.CoreDNSFlags.Replicas)
		spec.Spec.Replicas = &replicas

		spec.Spec.Template.Labels = labels

		spec.Spec.Template.Spec.Affinity = &v1.Affinity{
			NodeAffinity: &v1.NodeAffinity{
				PreferredDuringSchedulingIgnoredDuringExecution: []v1.PreferredSchedulingTerm{
					{
						Preference: v1.NodeSelectorTerm{
							MatchExpressions: []v1.NodeSelectorRequirement{
								{
									Key:      "node-role.kubernetes.io/master",
									Operator: "In",
									Values:   []string{""},
								},
							},
						},
						Weight: 100,
					},
				},
			},
		}

		spec.Spec.Template.Spec.PriorityClassName = "system-cluster-critical"
		spec.Spec.Template.Spec.RestartPolicy = "Always"
		spec.Spec.Template.Spec.ServiceAccountName = "coredns"
		spec.Spec.Template.Spec.DNSPolicy = "Default"

		var v int64 = 65532 //?
		spec.Spec.Template.Spec.SecurityContext = &v1.PodSecurityContext{
			RunAsUser: &v,
		}

		container, _ := installbase.AcceptContainerVisitor("coredns",
			ctx.CoreDNSFlags.Image,
			v1.PullIfNotPresent,
			newVisitor(ctx))

		spec.Spec.Template.Spec.Containers = append(spec.Spec.Template.Spec.Containers, *container)
		return spec
	}
}

func deploymentConfigVolumeSpec(fn deploymentSpecFunc) deploymentSpecFunc {
	var defaultMode int32 = 420
	return func(ctx *installbase.StageContext) *appsV1.Deployment {
		spec := fn(ctx)
		spec.Spec.Template.Spec.Volumes = []v1.Volume{
			{
				Name: "config-volume",
				VolumeSource: v1.VolumeSource{
					ConfigMap: &v1.ConfigMapVolumeSource{
						LocalObjectReference: v1.LocalObjectReference{
							Name: "coredns",
						},
						DefaultMode: &defaultMode,
						Items: []v1.KeyToPath{
							{
								Key:  "Corefile",
								Path: "Corefile",
							},
						},
					},
				},
			},
		}
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
	return nil, []string{"-conf", "/etc/coredns/Corefile"}
}

func (v *containerVisitor) VisitorContainerPorts(c *v1.Container) ([]v1.ContainerPort, error) {
	return []v1.ContainerPort{
		{
			Name:          "dns",
			ContainerPort: 53,
			Protocol:      v1.ProtocolUDP,
		},
		{
			Name:          "dns-tcp",
			ContainerPort: 53,
			Protocol:      v1.ProtocolTCP,
		},
		{
			Name:          "metrics",
			ContainerPort: 9153,
			Protocol:      v1.ProtocolTCP,
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
			Name:      "config-volume",
			MountPath: "/etc/coredns",
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
				Path:   "/health",
				Port:   intstr.FromInt(8080),
				Scheme: "HTTP",
			},
		},
		InitialDelaySeconds: 10,
		PeriodSeconds:       20,
		SuccessThreshold:    1,
		FailureThreshold:    10,
		TimeoutSeconds:      5,
	}, nil
}

func (v *containerVisitor) VisitorReadinessProbe(c *v1.Container) (*v1.Probe, error) {
	return &v1.Probe{
		Handler: v1.Handler{
			HTTPGet: &v1.HTTPGetAction{
				Path:   "/ready",
				Port:   intstr.FromInt(8181),
				Scheme: "HTTP",
			},
		},
		InitialDelaySeconds: 10,
		PeriodSeconds:       20,
		SuccessThreshold:    1,
		FailureThreshold:    10,
		TimeoutSeconds:      5,
	}, nil
}

func (v *containerVisitor) VisitorLifeCycle(c *v1.Container) (*v1.Lifecycle, error) {
	return nil, nil
}

func (v *containerVisitor) VisitorSecurityContext(c *v1.Container) (*v1.SecurityContext, error) {
	allowPrivilegeEscalation := false
	readOnlyRootFilesystem := true
	return &v1.SecurityContext{
		AllowPrivilegeEscalation: &allowPrivilegeEscalation,
		Capabilities: &v1.Capabilities{
			Add: []v1.Capability{
				"NET_BIND_SERVICE",
			},
			Drop: []v1.Capability{
				"all",
			},
		},
		ReadOnlyRootFilesystem: &readOnlyRootFilesystem,
	}, nil
}
