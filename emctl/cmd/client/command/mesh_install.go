/*
 * Copyright (c) 2017, MegaEase
 * All rights reserved.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package command

import (
	stdcontext "context"
	"fmt"
	"io/ioutil"

	"github.com/megaease/easemeshctl/cmd/client/command/flags"
	installbase "github.com/megaease/easemeshctl/cmd/client/command/meshinstall/base"
	"github.com/megaease/easemeshctl/cmd/client/command/meshinstall/controlpanel"
	"github.com/megaease/easemeshctl/cmd/client/command/meshinstall/crd"
	"github.com/megaease/easemeshctl/cmd/client/command/meshinstall/installation"
	"github.com/megaease/easemeshctl/cmd/client/command/meshinstall/meshingress"
	"github.com/megaease/easemeshctl/cmd/client/command/meshinstall/operator"
	"github.com/megaease/easemeshctl/cmd/client/command/rcfile"
	"github.com/megaease/easemeshctl/cmd/common"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// InstallCmd is the entrypoint of the emctl installation
func InstallCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "install",
		Short:   "Deploy infrastructure components of the EaseMesh",
		Long:    "",
		Example: "emctl install --clean-when-failed",
	}
	flags := &flags.Install{}
	flags.AttachCmd(cmd)

	cmd.Run = func(cmd *cobra.Command, args []string) {
		if flags.SpecFile != "" {
			var buff []byte
			var err error
			buff, err = ioutil.ReadFile(flags.SpecFile)
			if err != nil {
				common.ExitWithErrorf("%s failed: %v", cmd.Short, err)
			}

			err = yaml.Unmarshal(buff, flags)
			if err != nil {
				common.ExitWithErrorf("%s failed: %v", cmd.Short, err)
			}
		}
		install(cmd, flags)
	}

	return cmd
}

func install(cmd *cobra.Command, flags *flags.Install) {
	var err error
	kubeClient, err := installbase.NewKubernetesClient()
	if err != nil {
		common.ExitWithErrorf("%s failed: %v", cmd.Short, err)
	}

	apiExtensionClient, err := installbase.NewKubernetesAPIExtensionsClient()
	if err != nil {
		common.ExitWithErrorf("%s failed: %v", cmd.Short, err)
	}

	context := &installbase.StageContext{
		Flags:               flags,
		Client:              kubeClient,
		Cmd:                 cmd,
		APIExtensionsClient: apiExtensionClient,
	}

	install := installation.New(
		installation.Wrap(crd.PreCheck, crd.Deploy, crd.Clear, crd.Describe),
		installation.Wrap(controlpanel.PreCheck, controlpanel.Deploy, controlpanel.Clear, controlpanel.Describe),
		installation.Wrap(operator.PreCheck, operator.Deploy, operator.Clear, operator.Describe),
		installation.Wrap(meshingress.PreCheck, meshingress.Deploy, meshingress.Clear, meshingress.Describe),
	)

	err = install.DoInstallStage(context)
	if err != nil {
		if flags.CleanWhenFailed {
			install.ClearResource(context)
		}
		common.ExitWithErrorf("install mesh infrastructure error: %s", err)
	}

	postInstall(context)

	fmt.Println("Done.")
}

func postInstall(context *installbase.StageContext) {
	namespace := context.Flags.MeshNamespace
	name := installbase.DefaultMeshControlPlanePlubicServiceName
	service, err := context.Client.CoreV1().Services(namespace).Get(stdcontext.TODO(), name, metav1.GetOptions{})
	if err != nil {
		common.OutputErrorf("ignored: get service %s/%s failed: %v", namespace, name, err)
		return
	}

	rc, err := rcfile.New()
	if err != nil {
		common.OutputErrorf("ignored: new rcfile failed: %v", err)
		return
	}

	nodes, err := context.Client.CoreV1().Nodes().List(stdcontext.TODO(), metav1.ListOptions{})
	if err != nil {
		common.OutputErrorf("ignored: get nodes information failed: %v", err)
		return
	}
	firstNodeIP := ""
	for _, n := range nodes.Items {
		for _, address := range n.Status.Addresses {
			if address.Type == v1.NodeInternalIP {
				firstNodeIP = address.Address
			}
		}
	}

	if firstNodeIP == "" {
		common.OutputErrorf("ignored: no candidate node ip can be selected")
		return
	}

	for _, port := range service.Spec.Ports {
		if port.Name == installbase.DefaultMeshAdminPortName {
			rc.Server = fmt.Sprintf("%s:%d", firstNodeIP, port.NodePort)
			break
		}
	}

	if rc.Server == "" {
		common.OutputErrorf("ignored: %s of service %s/%s not found", installbase.DefaultMeshAdminPortName, namespace, name)
		return
	}

	err = rc.Marshal()
	if err != nil {
		common.OutputError(err)
	} else {
		fmt.Printf("run commands file: %s\n", rc.Path())
	}
}
