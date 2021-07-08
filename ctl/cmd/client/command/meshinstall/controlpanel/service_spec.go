package controlpanel

import (
	"github.com/megaease/easemeshctl/cmd/client/command/flags"
	installbase "github.com/megaease/easemeshctl/cmd/client/command/meshinstall/base"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/client-go/kubernetes"
)

func meshControlPanelLabel() map[string]string {
	selector := map[string]string{}
	selector["mesh-controlpanel-app"] = "easegress-mesh-controlpanel"
	return selector
}

func serviceSpec(installFlags *flags.Install) installbase.InstallFunc {

	labels := meshControlPanelLabel()

	headlessService := &v1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      installbase.DefaultMeshControlPlaneHeadlessServiceName,
			Namespace: installFlags.MeshNameSpace,
		},
	}

	headlessService.Spec.ClusterIP = "None"
	headlessService.Spec.Selector = labels
	headlessService.Spec.Ports = []v1.ServicePort{
		{
			Name:       installbase.DefaultMeshAdminPortName,
			Port:       int32(installFlags.EgAdminPort),
			TargetPort: intstr.IntOrString{IntVal: 2381},
		},
		{
			Name:       installbase.DefaultMeshPeerPortName,
			Port:       int32(installFlags.EgPeerPort),
			TargetPort: intstr.IntOrString{IntVal: 2380},
		},
		{
			Name:       installbase.DefaultMeshClientPortName,
			Port:       int32(installFlags.EgClientPort),
			TargetPort: intstr.IntOrString{IntVal: 2379},
		},
	}

	headfulService := &v1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      installFlags.EgServiceName,
			Namespace: installFlags.MeshNameSpace,
		},
	}

	headfulService.Spec.Selector = labels
	headfulService.Spec.Ports = []v1.ServicePort{
		{
			Name:       installbase.DefaultMeshAdminPortName,
			Port:       int32(installFlags.EgAdminPort),
			TargetPort: intstr.IntOrString{IntVal: 2381},
		},
		{
			Name:       installbase.DefaultMeshPeerPortName,
			Port:       int32(installFlags.EgPeerPort),
			TargetPort: intstr.IntOrString{IntVal: 2380},
		},
		{
			Name:       installbase.DefaultMeshClientPortName,
			Port:       int32(installFlags.EgClientPort),
			TargetPort: intstr.IntOrString{IntVal: 2379},
		},
	}

	service := &v1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      installbase.DefaultMeshControlPlanePlubicServiceName,
			Namespace: installFlags.MeshNameSpace,
		},
	}
	service.Spec.Ports = []v1.ServicePort{
		{
			Name:       installbase.DefaultMeshAdminPortName,
			Port:       int32(installFlags.EgAdminPort),
			TargetPort: intstr.IntOrString{IntVal: 2381},
		},
		{
			Name:       installbase.DefaultMeshPeerPortName,
			Port:       int32(installFlags.EgPeerPort),
			TargetPort: intstr.IntOrString{IntVal: 2380},
		},
		{
			Name:       installbase.DefaultMeshClientPortName,
			Port:       int32(installFlags.EgClientPort),
			TargetPort: intstr.IntOrString{IntVal: 2379},
		},
	}

	// FIXME: for test we leverage nodeport for expose controlpanel service
	// for production, we will give users options to switch to Loadbalance or ingress
	service.Spec.Type = v1.ServiceTypeNodePort
	service.Spec.Selector = labels

	return func(cmd *cobra.Command, client *kubernetes.Clientset, installFlags *flags.Install) error {
		err := installbase.DeployService(headlessService, client, installFlags.MeshNameSpace)
		if err != nil {
			return errors.Wrap(err, "deploy easemesh controlpanel inner service failed")
		}
		err = installbase.DeployService(service, client, installFlags.MeshNameSpace)
		if err != nil {
			return errors.Wrap(err, "deploy easemesh controlpanel public service failed")
		}

		err = installbase.DeployService(headfulService, client, installFlags.MeshNameSpace)
		if err != nil {
			return errors.Wrap(err, "deploy easemesh controlpanel headful service failed")
		}
		return nil
	}
}
