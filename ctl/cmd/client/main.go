/*
 * Copyright (c) 2017, MegaEase
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
# Install EaseMesh Components
emctl install --clean-when-failed

# Apply Tenant (kind is case-insensitive in command line)
emctl apply -f tenant-001.yaml

# Apply Service
emctl apply -f service-001.yaml

# Apply LoadBalance
emctl apply -f loadbalance.yaml

# Apply Ingress
emctl apply -f ingress.yaml

# Get service.
emctl get service
emctl get service -o yaml
emctl get service service-001 -o json

# Get LoadBalance
emctl get loadbalance
emctl get loadbalance service-001 -o yaml


# Delete service
emctl delete service service-001
emctl delete service -f service-001.yaml

# Delete LoadBalance
emctl delete loadbalance service-001

# NOTE: The manipulation of the kinds attached to Service below is the same with LoadBalance:
# - Sidecar
# - Resilience
# - Canary
# - ObservabilityMetrics, ObservabilityTracings, ObservabilityOutputServer`

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
		command.DeleteCmd(),
		command.GetCmd(),
		completionCmd,
	)

	err := rootCmd.Execute()
	if err != nil {
		common.ExitWithError(err)
	}
}
