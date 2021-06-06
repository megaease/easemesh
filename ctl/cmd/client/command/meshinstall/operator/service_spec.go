package operator

import (
	installbase "github.com/megaease/easemeshctl/cmd/client/command/meshinstall/base"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/client-go/kubernetes"
)

func serviceSpec(args *installbase.InstallArgs) installbase.InstallFunc {
	labels := meshOperatorLabels()

	service := &v1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      installbase.DefaultMeshOperatorControllerManagerServiceName,
			Namespace: args.MeshNameSpace,
		},
	}
	service.Spec.Ports = []v1.ServicePort{
		{
			Name:       "https",
			Port:       int32(8443),
			TargetPort: intstr.IntOrString{StrVal: "https"},
		},
	}
	service.Spec.Selector = labels
	return func(cmd *cobra.Command, kubeClient *kubernetes.Clientset, args *installbase.InstallArgs) error {
		err := installbase.DeployService(service, kubeClient, args.MeshNameSpace)
		if err != nil {
			return errors.Wrapf(err, "Create operator service %s error", args.MeshNameSpace)
		}
		return err
	}
}
