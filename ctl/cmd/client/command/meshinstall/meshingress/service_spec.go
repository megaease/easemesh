package meshingress

import (
	installbase "github.com/megaease/easemeshctl/cmd/client/command/meshinstall/base"

	"github.com/spf13/cobra"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/client-go/kubernetes"
)

func serviceSpec(args *installbase.InstallArgs) installbase.InstallFunc {
	service := &v1.Service{}
	service.Name = installbase.DefaultMeshIngressService

	service.Spec.Ports = []v1.ServicePort{
		{
			// FIXME: Don't specific a fix nodeport port number,
			// in the future, we need to specific a nodeport port via client argument
			Port:       13010,
			Protocol:   v1.ProtocolTCP,
			TargetPort: intstr.IntOrString{IntVal: 13010},
		},
	}
	service.Spec.Selector = meshIngressLabel()
	service.Spec.Type = v1.ServiceTypeNodePort
	return func(cmd *cobra.Command, kubeClient *kubernetes.Clientset, args *installbase.InstallArgs) error {
		err := installbase.DeployService(service, kubeClient, args.MeshNameSpace)
		return err
	}
}
