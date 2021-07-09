package operator

import (
	"github.com/megaease/easemeshctl/cmd/client/command/flags"
	installbase "github.com/megaease/easemeshctl/cmd/client/command/meshinstall/base"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	appsV1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/client-go/kubernetes"
)

type deploymentSpecFunc func(installFlags *flags.Install) *appsV1.Deployment

func operatorDeploymentSpec(installFlags *flags.Install) installbase.InstallFunc {

	deployment := deploymentConfigVolumeSpec(
		deploymentManagerContainerSpec(
			deploymentRBACContainerSpec(
				deploymentBaseSpec(deploymentInitialize(nil)))))(installFlags)

	return func(cmd *cobra.Command, kubeClient *kubernetes.Clientset, installFlags *flags.Install) error {
		err := installbase.DeployDeployment(deployment, kubeClient, installFlags.MeshNameSpace)
		if err != nil {
			return errors.Wrapf(err, "deployment operation %s failed", deployment.Name)
		}
		return err
	}
}

func meshOperatorLabels() map[string]string {
	selector := map[string]string{}
	selector["easemesh-operator"] = "operator-manager"
	return selector
}

func deploymentInitialize(fn deploymentSpecFunc) deploymentSpecFunc {
	return func(installFlags *flags.Install) *appsV1.Deployment {
		return &appsV1.Deployment{}
	}
}

func deploymentBaseSpec(fn deploymentSpecFunc) deploymentSpecFunc {
	return func(installFlags *flags.Install) *appsV1.Deployment {
		spec := fn(installFlags)

		labels := meshOperatorLabels()
		spec.Name = installbase.DefaultMeshOperatorName
		spec.Spec.Selector = &metav1.LabelSelector{
			MatchLabels: labels,
		}

		var replicas = int32(installFlags.EaseMeshOperatorReplicas)
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
	return func(installFlags *flags.Install) *appsV1.Deployment {
		spec := fn(installFlags)
		rbacContainer := v1.Container{}
		rbacContainer.Name = "kube-rbac-proxy"
		rbacContainer.Image = "gcr.io/kubebuilder/kube-rbac-proxy:v0.5.0"
		rbacContainer.Ports = []v1.ContainerPort{
			{
				Name:          "https",
				ContainerPort: int32(8443),
			},
		}
		rbacContainer.Args = []string{
			"--secure-listen-address=0.0.0.0:8443",
			"--upstream=http://127.0.0.1:8080/",
			"--logtostderr=true",
			"--v=10",
		}
		spec.Spec.Template.Spec.Containers =
			append(spec.Spec.Template.Spec.Containers, rbacContainer)
		return spec
	}
}

func deploymentConfigVolumeSpec(fn deploymentSpecFunc) deploymentSpecFunc {
	return func(installFlags *flags.Install) *appsV1.Deployment {
		spec := fn(installFlags)
		spec.Spec.Template.Spec.Volumes = []v1.Volume{
			{
				Name: "config-volume",
				VolumeSource: v1.VolumeSource{
					ConfigMap: &v1.ConfigMapVolumeSource{
						LocalObjectReference: v1.LocalObjectReference{
							Name: meshOperatorConfigMap,
						},
					},
				},
			},
		}
		return spec
	}

}

func deploymentManagerContainerSpec(fn deploymentSpecFunc) deploymentSpecFunc {

	return func(installFlags *flags.Install) *appsV1.Deployment {
		spec := fn(installFlags)
		container, _ := installbase.AcceptContainerVisistor("operator-manager",
			installFlags.ImageRegistryURL+"/"+installFlags.EaseMeshOperatorImage,
			v1.PullAlways,
			newVisitor(installFlags))

		spec.Spec.Template.Spec.Containers =
			append(spec.Spec.Template.Spec.Containers, *container)
		return spec
	}
}

func newVisitor(installFlags *flags.Install) installbase.ContainerVisitor {
	return &containerVisitor{installFlags: installFlags}
}

type containerVisitor struct {
	installFlags *flags.Install
}

func (v *containerVisitor) VisitorCommandAndArgs(c *v1.Container) (command []string, installFlags []string) {
	return []string{"/manager"},
		[]string{"--config=/opt/mesh/operator-config.yaml"}
}

func (v *containerVisitor) VisitorContainerPorts(c *v1.Container) ([]v1.ContainerPort, error) {
	return nil, nil
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
			MountPath: "/opt/mesh/operator-config.yaml",
			SubPath:   "operator-config.yaml",
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
