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

package apply

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

// Run is the entrypoint of the emctl apply subcommand
func Run(cmd *cobra.Command, flags *flags.Apply) {
	if flags.YamlFile == "" {
		common.ExitWithErrorf("no resource specified")
	}

	vss, err := util.NewVisitorBuilder().
		FilenameParam(&util.FilenameOptions{
			Recursive: flags.Recursive,
			Filenames: []string{flags.YamlFile},
		}).
		Do()

	if err != nil {
		common.ExitWithErrorf("build visitor failed: %v", err)
	}

	var errs []error
	for _, vs := range vss {
		err := vs.Visit(func(mo resource.MeshObject, e error) error {
			if e != nil {
				return errors.Wrap(e, "visit failed")
			}

			err := WrapApplierByMeshObject(mo, meshclient.New(flags.Server), flags.Timeout).Apply()
			if err != nil {
				return fmt.Errorf("%s/%s applied failed: %s", mo.Kind(), mo.Name(), err)
			}

			fmt.Printf("%s/%s applied successfully\n", mo.Kind(), mo.Name())
			return nil
		})

		if err != nil {
			errs = append(errs, err)
		}
	}

	if len(errs) > 0 {
		common.ExitWithErrorf("applying resources has errors occurred")
	}
}
