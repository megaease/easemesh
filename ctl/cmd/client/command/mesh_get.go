package command

import (
	"time"

	"github.com/megaease/easemeshctl/cmd/client/command/get"

	"github.com/spf13/cobra"
)

func getGetArgs(cmd *cobra.Command, args *get.Arguments) error {
	var timeOutInMills int64
	cmd.Flags().StringVarP(&args.Server, "server", "s", "127.0.0.1:2381", "An address to access the EaseMesh control plane")
	cmd.Flags().Int64VarP(&timeOutInMills, "timeout", "t", 30000, "A duration that limit max time out for requesting the EaseMesh control plane, in millseconds unit (default: 30000)")
	cmd.Flags().StringVarP(&args.OutputFormat, "output", "o", "table", "Output format (support table, yaml, json)")
	args.Timeout = time.Millisecond * time.Duration(timeOutInMills)

	return nil
}

func GetCmd() *cobra.Command {
	var getArgs get.Arguments

	cmd := &cobra.Command{
		Use:     "get",
		Short:   "Get resources of easemesh",
		Example: "emctl get -f config.yaml | emctl get service service-001",
		Run: func(cmd *cobra.Command, args []string) {
			get.Run(cmd, &getArgs)
		},
	}

	getGetArgs(cmd, &getArgs)

	return cmd
}
