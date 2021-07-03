package command

import (
	"time"

	"github.com/megaease/easemeshctl/cmd/client/command/apply"

	"github.com/spf13/cobra"
)

func getApplyArgs(cmd *cobra.Command, args *apply.Arguments) error {
	var timeOutInMills int64
	cmd.Flags().StringVarP(&args.YamlFile, "file", "f", "", "A location contained the EaseMesh resource files (YAML format) to apply, could be a file, directory, or URL")
	cmd.Flags().StringVarP(&args.Server, "server", "s", "127.0.0.1:2381", "An address to access the EaseMesh control plane")
	cmd.Flags().Int64VarP(&timeOutInMills, "timeout", "t", 30000, "A duration that limit max time out for requesting the EaseMesh control plane, in millseconds unit (default: 30000)")
	cmd.Flags().BoolVarP(&args.Recursive, "recursive", "r", true, "Whether to recursively iterate all sub-directories and files of the location (default: true)")
	args.Timeout = time.Millisecond * time.Duration(timeOutInMills)
	return nil
}

func ApplyCmd() *cobra.Command {
	var applyArgs apply.Arguments

	cmd := &cobra.Command{
		Use:     "apply",
		Short:   "Apply a configuration to easemesh",
		Long:    "",
		Example: "emctl apply -f config.yaml",
		Run: func(cmd *cobra.Command, args []string) {
			apply.Run(cmd, &applyArgs)
		},
	}

	getApplyArgs(cmd, &applyArgs)
	return cmd
}
