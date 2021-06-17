package controlpanel

import (
	installbase "github.com/megaease/easemeshctl/cmd/client/command/meshinstall/base"

	"github.com/spf13/cobra"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func namespaceSpec(args *installbase.InstallArgs) installbase.InstallFunc {
	ns := &v1.Namespace{ObjectMeta: metav1.ObjectMeta{
		Name:   args.MeshNameSpace,
		Labels: map[string]string{},
	}}
	return func(cmd *cobra.Command, client *kubernetes.Clientset, args *installbase.InstallArgs) error {
		err := installbase.CreateNameSpace(ns, client)
		if err != nil && !errors.IsAlreadyExists(err) {
			return err
		}
		return nil
	}
}
