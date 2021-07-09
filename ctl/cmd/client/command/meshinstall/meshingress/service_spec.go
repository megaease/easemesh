package meshingress

import (
	"github.com/megaease/easemeshctl/cmd/client/command/flags"
	installbase "github.com/megaease/easemeshctl/cmd/client/command/meshinstall/base"

	"github.com/spf13/cobra"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/client-go/kubernetes"
)

func serviceSpec(installFlags *flags.Install) installbase.InstallFunc {
	service := &v1.Service{}
	service.Name = installbase.DefaultMeshIngressService

	service.Spec.Ports = []v1.ServicePort{
		{
			Port:       installFlags.MeshIngressServicePort,
			Protocol:   v1.ProtocolTCP,
			TargetPort: intstr.IntOrString{IntVal: installFlags.MeshIngressServicePort},
		},
	}
	service.Spec.Selector = meshIngressLabel()
	service.Spec.Type = v1.ServiceTypeNodePort
	return func(cmd *cobra.Command, kubeClient *kubernetes.Clientset, installFlags *flags.Install) error {
		err := installbase.DeployService(service, kubeClient, installFlags.MeshNameSpace)
		return err
	}
}
