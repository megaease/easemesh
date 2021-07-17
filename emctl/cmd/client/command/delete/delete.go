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

package delete

import (
	"fmt"

	"github.com/megaease/easemeshctl/cmd/client/command/flags"
	"github.com/megaease/easemeshctl/cmd/client/command/meshclient"
	"github.com/megaease/easemeshctl/cmd/client/resource"
	"github.com/megaease/easemeshctl/cmd/client/util"
	"github.com/megaease/easemeshctl/cmd/common"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

// Run is the entrypoint of the emctl delete sub command
func Run(cmd *cobra.Command, flags *flags.Delete) {
	visitorBulder := util.NewVisitorBuilder()

	cmdArgs := cmd.Flags().Args()

	if len(cmdArgs) == 0 && flags.YamlFile == "" {
		common.ExitWithErrorf("no resource specified")
	}

	if len(cmdArgs) != 0 {
		if flags.YamlFile != "" {
			common.ExitWithErrorf("file and command args are both specified")
		}
		if len(cmdArgs) != 2 {
			common.ExitWithErrorf("invalid command args: support <resource kind> <resource name>")
		}
		visitorBulder.CommandParam(&util.CommandOptions{
			Kind: cmdArgs[0],
			Name: cmdArgs[1],
		})
	}

	if flags.YamlFile != "" {
		visitorBulder.FilenameParam(&util.FilenameOptions{
			Recursive: flags.Recursive,
			Filenames: []string{flags.YamlFile},
		})
	}

	vss, err := visitorBulder.Do()
	if err != nil {
		common.ExitWithErrorf("build visitor failed: %s", err)
	}

	var errs []error
	for _, vs := range vss {
		err := vs.Visit(func(mo resource.MeshObject, e error) error {
			if e != nil {
				return errors.Wrap(e, "visit failed")
			}

			err := WrapDeleterByMeshObject(mo, meshclient.New(flags.Server), flags.Timeout).Delete()
			if err != nil {
				return errors.Wrapf(err, "%s/%s deleted failed", mo.Kind(), mo.Name())
			}

			fmt.Printf("%s/%s deleted successfully\n", mo.Kind(), mo.Name())
			return nil
		})

		if err != nil {
			errs = append(errs, err)
		}
	}

	if len(errs) > 0 {
		common.ExitWithErrorf("deleting resources has errors occurred")
	}
}
