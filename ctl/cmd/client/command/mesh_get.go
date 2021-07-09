package command

import (
	"github.com/megaease/easemeshctl/cmd/client/command/flags"
	"github.com/megaease/easemeshctl/cmd/client/command/get"

	"github.com/spf13/cobra"
)

func GetCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "get",
		Short:   "Get resources of easemesh",
		Example: "emctl get -f config.yaml | emctl get service service-001",
	}

	flags := &flags.Get{}
	flags.AttachCmd(cmd)

	cmd.Run = func(cmd *cobra.Command, args []string) {
		get.Run(cmd, flags)
	}

	return cmd
}
