package main

import (
	"os"

	"github.com/megaease/easemeshctl/cmd/client/command"
	"github.com/megaease/easemeshctl/cmd/common"

	"github.com/spf13/cobra"
)

func init() {
	cobra.EnablePrefixMatching = true
}

var exampleUsage = ` # EaseMesh command line tool for management and operation
  emctl <subcommand> 
`

func main() {
	rootCmd := &cobra.Command{
		Use:        "emctl",
		Short:      "A command line tool for EaseMesh management and operation",
		Example:    exampleUsage,
		SuggestFor: []string{"emctl"},
	}

	completionCmd := &cobra.Command{
		Use:   "completion bash|zsh",
		Short: "Output shell completion code for the specified shell (bash or zsh)",
		Run: func(cmd *cobra.Command, args []string) {
			switch args[0] {
			case "bash":
				rootCmd.GenBashCompletion(os.Stdout)
			case "zsh":
				rootCmd.GenZshCompletion(os.Stdout)
			default:
				common.ExitWithErrorf("unsupported shell %s, expecting bash or zsh", args[0])
			}
		},
		Args: cobra.ExactArgs(1),
	}

	rootCmd.AddCommand(
		command.InstallCmd(),
		command.ResetCmd(),
		command.ApplyCmd(),
		completionCmd,
	)

	err := rootCmd.Execute()
	if err != nil {
		common.ExitWithError(err)
	}
}
