package command

import (
	"time"

	"github.com/megaease/easemeshctl/cmd/client/command/delete"

	"github.com/spf13/cobra"
)

func getDeleteArgs(cmd *cobra.Command, args *delete.Arguments) error {
	var timeOutInMills int64
	cmd.Flags().StringVarP(&args.YamlFile, "file", "f", "", "A location contained the EaseMesh resource files (YAML format) to delete, could be a file, directory, or URL")
	cmd.Flags().StringVarP(&args.Server, "server", "s", "127.0.0.1:2381", "An address to access the EaseMesh control plane")
	cmd.Flags().Int64VarP(&timeOutInMills, "timeout", "t", 30000, "A duration that limit max time out for requesting the EaseMesh control plane, in millseconds unit (default: 30000)")
	cmd.Flags().BoolVarP(&args.Recursive, "recursive", "r", true, "Whether to recursively iterate all sub-directories and files of the location (default: true)")
	args.Timeout = time.Millisecond * time.Duration(timeOutInMills)

	return nil
}

func DeleteCmd() *cobra.Command {
	var deleteArgs delete.Arguments

	cmd := &cobra.Command{
		Use:     "delete",
		Short:   "Delete resources of easemesh",
		Example: "emctl delete -f config.yaml | emctl delete service service-001",
		Run: func(cmd *cobra.Command, args []string) {
			delete.Run(cmd, &deleteArgs)
		},
	}

	getDeleteArgs(cmd, &deleteArgs)

	return cmd
}
