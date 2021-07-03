package command

import (
	installbase "github.com/megaease/easemeshctl/cmd/client/command/meshinstall/base"
	"github.com/megaease/easemeshctl/cmd/client/command/meshinstall/controlpanel"
	"github.com/megaease/easemeshctl/cmd/client/command/meshinstall/crd"
	"github.com/megaease/easemeshctl/cmd/client/command/meshinstall/installation"
	"github.com/megaease/easemeshctl/cmd/client/command/meshinstall/meshingress"
	"github.com/megaease/easemeshctl/cmd/client/command/meshinstall/operator"
	"github.com/megaease/easemeshctl/cmd/common"

	"github.com/spf13/cobra"
)

func reset(cmd *cobra.Command, args *installbase.InstallArgs) {
	kubeClient, err := installbase.NewKubernetesClient()
	if err != nil {
		common.ExitWithErrorf("%s failed: %v", cmd.Short, err)
	}

	apiExtensionClient, err := installbase.NewKubernetesApiExtensionsClient()
	if err != nil {
		common.ExitWithErrorf("%s failed: %v", cmd.Short, err)
	}

	clearFuncs := []installation.ClearFunc{
		meshingress.Clear,
		operator.Clear,
		controlpanel.Clear,
		crd.Clear,
	}

	stageContext := installbase.StageContext{
		Cmd:                 cmd,
		Client:              kubeClient,
		Arguments:           *args,
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

func ResetCmd() *cobra.Command {
	iargs := &installbase.InstallArgs{}
	cmd := &cobra.Command{
		Use:     "reset",
		Short:   "Reset infrastructure components of the EaseMesh",
		Long:    "",
		Example: "emctl reset",
		Run: func(cmd *cobra.Command, args []string) {
			reset(cmd, iargs)
		},
	}

	baseCmdArgs(cmd, iargs)
	return cmd
}
