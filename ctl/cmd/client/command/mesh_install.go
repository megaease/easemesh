package command

import (
	"fmt"
	"io/ioutil"

	"github.com/megaease/easemeshctl/cmd/client/command/flags"
	installbase "github.com/megaease/easemeshctl/cmd/client/command/meshinstall/base"
	"github.com/megaease/easemeshctl/cmd/client/command/meshinstall/controlpanel"
	"github.com/megaease/easemeshctl/cmd/client/command/meshinstall/crd"
	"github.com/megaease/easemeshctl/cmd/client/command/meshinstall/installation"
	"github.com/megaease/easemeshctl/cmd/client/command/meshinstall/meshingress"
	"github.com/megaease/easemeshctl/cmd/client/command/meshinstall/operator"
	"github.com/megaease/easemeshctl/cmd/common"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

func InstallCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "install",
		Short:   "Deploy infrastructure components of the EaseMesh",
		Long:    "",
		Example: "emctl install <args>",
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

	apiExtensionClient, err := installbase.NewKubernetesApiExtensionsClient()
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

	fmt.Println("Done.")
}
