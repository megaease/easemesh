package command

import (
	"github.com/megaease/easemeshctl/cmd/client/command/delete"
	"github.com/megaease/easemeshctl/cmd/client/command/flags"

	"github.com/spf13/cobra"
)

func DeleteCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "delete",
		Short:   "Delete resources of easemesh",
		Example: "emctl delete -f config.yaml | emctl delete service service-001",
	}

	flags := &flags.Delete{}
	flags.AttachCmd(cmd)

	cmd.Run = func(cmd *cobra.Command, args []string) {
		delete.Run(cmd, flags)
	}

	return cmd
}
