package controlpanel

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/megaease/easemeshctl/cmd/client/command/flags"
	installbase "github.com/megaease/easemeshctl/cmd/client/command/meshinstall/base"
	"github.com/megaease/easemeshctl/cmd/common"
	"github.com/megaease/easemeshctl/cmd/common/client"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"k8s.io/client-go/kubernetes"
)

func provisionEaseMeshControlPanel(cmd *cobra.Command, kubeClient *kubernetes.Clientset, installFlags *flags.Install) error {

	entrypoints, err := installbase.GetMeshControlPanelEntryPoints(kubeClient, installFlags.MeshNameSpace,
		installbase.DefaultMeshControlPlanePlubicServiceName,
		installbase.DefaultMeshAdminPortName)
	if err != nil {
		return errors.Wrap(err, "get mesh control panel entrypoint failed")
	}

	meshControllerConfig := installbase.MeshControllerConfig{
		Name:              installbase.DefaultMeshControllerName,
		Kind:              flags.MeshControllerKind,
		RegistryType:      installFlags.EaseMeshRegistryType,
		HeartbeatInterval: strconv.Itoa(installFlags.HeartbeatInterval) + "s",
		IngressPort:       installFlags.MeshIngressServicePort,
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

func clearEaseMeshControlPanelProvision(cmd *cobra.Command, kubeClient *kubernetes.Clientset, installFlags *flags.Install) {

	entrypoints, err := installbase.GetMeshControlPanelEntryPoints(kubeClient, installFlags.MeshNameSpace,
		installbase.DefaultMeshControlPlanePlubicServiceName,
		installbase.DefaultMeshAdminPortName)
	if err != nil {
		common.OutputErrorf("clear: get mesh control panel entrypoint failed %s", err)
		return
	}

	for _, entrypoint := range entrypoints {
		url := fmt.Sprintf(entrypoint+installbase.ObjectURL, installbase.DefaultMeshControllerName)
		_, err = client.NewHTTPJSON().
			Delete(url, nil, time.Second*5, nil).
			HandleResponse(func(body []byte, statusCode int) (interface{}, error) {
				if statusCode == http.StatusNotFound {
					return nil, nil
				}
				if statusCode >= 400 {
					return nil, errors.Errorf("setup EaseMesh controller panel error, controller panel return statusCode %d, body: %s", statusCode, string(body))
				}
				return nil, nil
			})
		if err != nil {
			common.OutputErrorf("delete mesh controller configuration from %s failed %s", url, err)
		}
	}
}
