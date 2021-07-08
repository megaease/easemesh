package controlpanel

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/megaease/easemeshctl/cmd/client/command/flags"
	installbase "github.com/megaease/easemeshctl/cmd/client/command/meshinstall/base"

	yamljsontool "github.com/ghodss/yaml"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func configMapSpec(installFlags *flags.Install) installbase.InstallFunc {
	var host = "0.0.0.0"

	var config = installbase.EasegressConfig{
		Name:                    installbase.DefaultMeshControlPlaneName,
		ClusterName:             installbase.DefaultMeshControlPlaneName,
		ClusterRole:             installbase.WriterClusterRole,
		ClusterListenClientUrls: []string{"http://" + "0.0.0.0:" + strconv.Itoa(installFlags.EgClientPort)},
		ClusterListenPeerUrls:   []string{"http://" + "0.0.0.0:" + strconv.Itoa(installFlags.EgPeerPort)},
		ClusterJoinUrls:         []string{},
		ApiAddr:                 host + ":" + strconv.Itoa(installFlags.EgAdminPort),
		DataDir:                 "/opt/eg-data/data",
		WalDir:                  "",
		CpuProfileFile:          "",
		MemoryProfileFile:       "",
		LogDir:                  "/opt/eg-data/log",
		MemberDir:               "/opt/eg-data/member",
		StdLogLevel:             "INFO",
	}

	for i := 0; i < installFlags.EasegressControlPlaneReplicas; i++ {
		config.ClusterJoinUrls = append(config.ClusterJoinUrls,
			fmt.Sprintf("http://%s-%d.%s.%s:%d",
				installbase.DefaultMeshControlPlaneName,
				i,
				installbase.DefaultMeshControlPlaneHeadlessServiceName,
				installFlags.MeshNameSpace,
				installFlags.EgPeerPort))
	}

	configData := map[string]string{}
	configBytes, _ := yaml.Marshal(config)
	configData["eg-master.yaml"] = string(configBytes)

	buff, _ := yaml.Marshal(configData)
	configJson, _ := yamljsontool.YAMLToJSON(buff)

	var params map[string]string
	_ = json.Unmarshal(configJson, &params)

	configMap := &v1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      installbase.DefaultMeshControlPlaneConfig,
			Namespace: installFlags.MeshNameSpace,
		},
		Data: params,
	}

	return func(cmd *cobra.Command, kubeClient *kubernetes.Clientset, installFlags *flags.Install) error {
		err := installbase.DeployConfigMap(configMap, kubeClient, installFlags.MeshNameSpace)
		if err != nil {
			return err
		}
		return nil
	}
}
