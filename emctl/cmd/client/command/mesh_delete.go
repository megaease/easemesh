/*
 * Copyright (c) 2021, MegaEase
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

package command

import (
	"github.com/megaease/easemeshctl/cmd/client/command/delete"
	"github.com/megaease/easemeshctl/cmd/client/command/flags"

	"github.com/spf13/cobra"
)

// DeleteCmd invokes delete sub command entrypoint
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
