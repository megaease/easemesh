package meshingress

import (
	installbase "github.com/megaease/easemeshctl/cmd/client/command/meshinstall/base"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	appsV1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/client-go/kubernetes"
)

type deploymentSpecFunc func(*installbase.InstallArgs) *appsV1.Deployment

func meshIngressLabel() map[string]string {
	selector := map[string]string{}
	selector["app"] = "Easegress-ingress"
	return selector
}

func deploymentSpec(args *installbase.InstallArgs) installbase.InstallFunc {
	deployment := deploymentConfigVolumeSpec(
		deploymentContainerSpec(
			deploymentBaseSpec(
				deploymentInitialize(nil))))(args)

	return func(cmd *cobra.Command, kubeClient *kubernetes.Clientset, args *installbase.InstallArgs) error {
		err := installbase.DeployDeployment(deployment, kubeClient, args.MeshNameSpace)
		if err != nil {
			return errors.Wrapf(err, "deployment operation %s failed", deployment.Name)
		}
		return err
	}
}

func deploymentInitialize(fn deploymentSpecFunc) deploymentSpecFunc {
	return func(args *installbase.InstallArgs) *appsV1.Deployment {
		return &appsV1.Deployment{}
	}
}

func deploymentBaseSpec(fn deploymentSpecFunc) deploymentSpecFunc {
	return func(args *installbase.InstallArgs) *appsV1.Deployment {
		spec := fn(args)
		spec.Name = installbase.DefaultMeshIngressControllerName
		spec.Spec.Selector = &metav1.LabelSelector{
			MatchLabels: meshIngressLabel(),
		}

		var replicas = int32(args.EasegressIngressReplicas)
		spec.Spec.Replicas = &replicas
		spec.Spec.Template.Labels = meshIngressLabel()
		spec.Spec.Template.Spec.Containers = []v1.Container{}
		return spec
	}
}

func deploymentContainerSpec(fn deploymentSpecFunc) deploymentSpecFunc {

	return func(args *installbase.InstallArgs) *appsV1.Deployment {

		spec := fn(args)
		container, _ := installbase.AcceptContainerVisistor("easegress-ingress",
			args.ImageRegistryURL+"/"+args.EasegressImage,
			v1.PullAlways,
			newVisitor(args))

		spec.Spec.Template.Spec.Containers = append(spec.Spec.Template.Spec.Containers, *container)
		return spec
	}
}

func deploymentConfigVolumeSpec(fn deploymentSpecFunc) deploymentSpecFunc {

	return func(args *installbase.InstallArgs) *appsV1.Deployment {
		spec := fn(args)
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
	args *installbase.InstallArgs
}

func newVisitor(args *installbase.InstallArgs) installbase.ContainerVisitor {
	return &containerVisitor{args}
}

func (v *containerVisitor) VisitorCommandAndArgs(c *v1.Container) (command []string, args []string) {

	return []string{"/bin/sh"},
		[]string{"-c", "/opt/easegress/bin/easegress-server -f /opt/eg-config/eg-ingress.yaml"}
}

func (v *containerVisitor) VisitorContainerPorts(c *v1.Container) ([]v1.ContainerPort, error) {

	return []v1.ContainerPort{
		{
			Name:          installbase.DefaultMeshAdminPortName,
			ContainerPort: installbase.DefaultMeshAdminPort,
		},
		{
			Name:          installbase.DefaultMeshClientPortName,
			ContainerPort: installbase.DefaultMeshClientPort,
		},
		{
			Name:          installbase.DefaultMeshPeerPortName,
			ContainerPort: installbase.DefaultMeshPeerPort,
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

	return &v1.Probe{
		Handler: v1.Handler{
			HTTPGet: &v1.HTTPGetAction{
				Host: "127.0.0.1",
				Port: intstr.FromInt(installbase.DefaultMeshAdminPort),
				Path: "/apis/v1/healthz",
			},
		},
		InitialDelaySeconds: 50,
	}, nil
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
