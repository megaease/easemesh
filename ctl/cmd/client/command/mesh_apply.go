package command

import (
	"github.com/megaease/easemeshctl/cmd/client/command/apply"
	"github.com/megaease/easemeshctl/cmd/client/command/flags"

	"github.com/spf13/cobra"
)

func ApplyCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "apply",
		Short:   "Apply a configuration to easemesh",
		Long:    "",
		Example: "emctl apply -f config.yaml",
	}

	flags := &flags.Apply{}
	flags.AttachCmd(cmd)

	cmd.Run = func(cmd *cobra.Command, args []string) {
		apply.Run(cmd, flags)
	}

	return cmd
}
