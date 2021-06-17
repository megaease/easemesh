package operator

import (
	"fmt"
	"strconv"

	installbase "github.com/megaease/easemeshctl/cmd/client/command/meshinstall/base"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func configMapSpec(args *installbase.InstallArgs) installbase.InstallFunc {

	cfg := installbase.MeshOperatorConfig{
		ImageRegistryURL:     args.ImageRegistryURL,
		ClusterName:          installbase.DefaultMeshControlPlaneName,
		ClusterJoinURLs:      "http://" + installbase.DefaultMeshControlPlaneHeadfulServiceName + "." + args.MeshNameSpace + ":" + strconv.Itoa(args.EgPeerPort),
		MetricsAddr:          "127.0.0.1:8080",
		EnableLeaderElection: false,
		ProbeAddr:            ":8081",
	}

	configMap := &v1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      meshOperatorConfigMap,
			Namespace: args.MeshNameSpace,
		},
	}
	operatorConfig, err := yaml.Marshal(cfg)
	if err == nil {
		// error will left for high order function to jude
		data := map[string]string{}
		data["operator-config.yaml"] = string(operatorConfig)
		configMap.Data = data
	}

	return func(cmd *cobra.Command, client *kubernetes.Clientset, args *installbase.InstallArgs) error {
		if err != nil {
			return errors.Wrap(err, "ConfigMap build error")
		}
		err = installbase.DeployConfigMap(configMap, client, args.MeshNameSpace)
		if err != nil {
			return fmt.Errorf("create configMap failed: %v ", err)
		}
		return err
	}
}
