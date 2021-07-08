package operator

import (
	"fmt"
	"strconv"

	"github.com/megaease/easemeshctl/cmd/client/command/flags"
	installbase "github.com/megaease/easemeshctl/cmd/client/command/meshinstall/base"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func configMapSpec(installFlags *flags.Install) installbase.InstallFunc {

	cfg := installbase.MeshOperatorConfig{
		ImageRegistryURL:     installFlags.ImageRegistryURL,
		ClusterName:          installbase.DefaultMeshControlPlaneName,
		ClusterJoinURLs:      "http://" + flags.DefaultMeshControlPlaneHeadfulServiceName + "." + installFlags.MeshNameSpace + ":" + strconv.Itoa(installFlags.EgPeerPort),
		MetricsAddr:          "127.0.0.1:8080",
		EnableLeaderElection: false,
		ProbeAddr:            ":8081",
	}

	configMap := &v1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      meshOperatorConfigMap,
			Namespace: installFlags.MeshNameSpace,
		},
	}
	operatorConfig, err := yaml.Marshal(cfg)
	if err == nil {
		// error will left for high order function to jude
		data := map[string]string{}
		data["operator-config.yaml"] = string(operatorConfig)
		configMap.Data = data
	}

	return func(cmd *cobra.Command, client *kubernetes.Clientset, installFlags *flags.Install) error {
		if err != nil {
			return errors.Wrap(err, "ConfigMap build error")
		}
		err = installbase.DeployConfigMap(configMap, client, installFlags.MeshNameSpace)
		if err != nil {
			return fmt.Errorf("create configMap failed: %v ", err)
		}
		return err
	}
}
