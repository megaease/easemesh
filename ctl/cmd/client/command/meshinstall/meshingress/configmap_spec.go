package meshingress

import (
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
	params := &installbase.EasegressReaderParams{}
	params.ClusterRole = installbase.ReaderClusterRole
	params.ClusterRequestTimeout = "10s"
	params.ClusterJoinUrls = "http://" + flags.DefaultMeshControlPlaneHeadfulServiceName + ":" + strconv.Itoa(installFlags.EgPeerPort)
	params.ClusterName = installbase.DefaultMeshControlPlaneName
	params.Name = "mesh-ingress"

	labels := make(map[string]string)
	labels["mesh-role"] = "ingress-controller"
	params.Labels = labels

	data := map[string]string{}
	ingressControllerConfig, err := yaml.Marshal(params)
	configMap := &v1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      installbase.DefaultMeshIngressConfig,
			Namespace: installFlags.MeshNameSpace,
		},
	}
	if err == nil {
		data["eg-ingress.yaml"] = string(ingressControllerConfig)
		configMap.Data = data
	}

	return func(cmd *cobra.Command, kubeClient *kubernetes.Clientset, installFlags *flags.Install) error {
		if err != nil {
			return errors.Wrapf(err, "Create MeshIngress %s configmap spec error", configMap.Name)
		}
		err = installbase.DeployConfigMap(configMap, kubeClient, installFlags.MeshNameSpace)
		if err != nil {
			return errors.Wrapf(err, "Deploy configmap %s error", configMap.Name)
		}
		return nil
	}
}
