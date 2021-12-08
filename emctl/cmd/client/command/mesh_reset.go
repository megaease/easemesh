/*
 * Copyright (c) 2021, MegaEase
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
	"github.com/megaease/easemeshctl/cmd/client/command/flags"
	installbase "github.com/megaease/easemeshctl/cmd/client/command/meshinstall/base"
	"github.com/megaease/easemeshctl/cmd/client/command/meshinstall/controlpanel"
	"github.com/megaease/easemeshctl/cmd/client/command/meshinstall/crd"
	"github.com/megaease/easemeshctl/cmd/client/command/meshinstall/installation"
	"github.com/megaease/easemeshctl/cmd/client/command/meshinstall/meshingress"
	"github.com/megaease/easemeshctl/cmd/client/command/meshinstall/operator"
	"github.com/megaease/easemeshctl/cmd/client/command/meshinstall/shadowservice"
	"github.com/megaease/easemeshctl/cmd/common"

	"github.com/spf13/cobra"
)

func reset(cmd *cobra.Command, resetFlags *flags.Reset) {
	kubeClient, err := installbase.NewKubernetesClient()
	if err != nil {
		common.ExitWithErrorf("%s failed: %v", cmd.Short, err)
	}

	apiExtensionClient, err := installbase.NewKubernetesAPIExtensionsClient()
	if err != nil {
		common.ExitWithErrorf("%s failed: %v", cmd.Short, err)
	}

	clearFuncs := []installation.ClearFunc{
		shadowservice.Clear,
		meshingress.Clear,
		operator.Clear,
		controlpanel.Clear,
		crd.Clear,
	}

	stageContext := installbase.StageContext{
		Cmd:                 cmd,
		Client:              kubeClient,
		Flags:               &flags.Install{OperationGlobal: resetFlags.OperationGlobal},
		APIExtensionsClient: apiExtensionClient,
		ClearFuncs:          nil,
	}

	for _, f := range clearFuncs {
		err := f(&stageContext)
		if err != nil {
			common.OutputErrorf("ignored a reseting resource error %s", err)
		}
	}
}

// ResetCmd invoke reset sub command entrypoint
func ResetCmd() *cobra.Command {
	flags := &flags.Reset{}

	cmd := &cobra.Command{
		Use:     "reset",
		Short:   "Reset infrastructure components of the EaseMesh",
		Long:    "",
		Example: "emctl reset",
	}

	flags.AttachCmd(cmd)
	cmd.Run = func(cmd *cobra.Command, args []string) {
		reset(cmd, flags)
	}

	return cmd
}
