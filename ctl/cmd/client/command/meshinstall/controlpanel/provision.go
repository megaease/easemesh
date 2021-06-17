package controlpanel

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	installbase "github.com/megaease/easemeshctl/cmd/client/command/meshinstall/base"
	"github.com/megaease/easemeshctl/cmd/common/client"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"k8s.io/client-go/kubernetes"
)

func provisionEaseMeshControlPanel(cmd *cobra.Command, kubeClient *kubernetes.Clientset, args *installbase.InstallArgs) error {

	entrypoints, err := installbase.GetMeshControlPanelEntryPoints(kubeClient, args.MeshNameSpace,
		installbase.DefaultMeshControlPlanePlubicServiceName,
		installbase.DefaultMeshAdminPortName)
	if err != nil {
		return errors.Wrap(err, "get mesh control panel entrypoint failed")
	}

	meshControllerConfig := installbase.MeshControllerConfig{
		Name:              installbase.DefaultMeshControllerName,
		Kind:              installbase.MeshControllerKind,
		RegistryType:      args.EaseMeshRegistryType,
		HeartbeatInterval: strconv.Itoa(args.HeartbeatInterval) + "s",
		IngressPort:       args.MeshIngressServicePort,
	}

	configBody, err := json.Marshal(meshControllerConfig)
	if err != nil {
		return fmt.Errorf("startUp MeshController failed: %v", err)
	}

	for _, entrypoint := range entrypoints {
		url := entrypoint + installbase.ObjectsURL
		_, err = client.NewHTTPJSON().
			Post(url, configBody, time.Second*5, nil).
			HandleResponse(func(body []byte, statusCode int) (interface{}, error) {
				if statusCode >= 400 {
					return nil, errors.Errorf("setup EaseMesh controller panel error, controller panel return statusCode %d, body: %s", statusCode, string(body))
				}
				return nil, nil
			})
		if err == nil {
			return nil
		}
	}

	return errors.Wrapf(err, "call EaseMesh control panel %v error", entrypoints)
}
