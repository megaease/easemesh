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

package get

import (
	"github.com/megaease/easemeshctl/cmd/client/command/flags"
	"github.com/megaease/easemeshctl/cmd/client/command/meshclient"
	"github.com/megaease/easemeshctl/cmd/client/command/printer"
	"github.com/megaease/easemeshctl/cmd/client/resource"
	"github.com/megaease/easemeshctl/cmd/client/util"
	"github.com/megaease/easemeshctl/cmd/common"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

func Run(cmd *cobra.Command, flags *flags.Get) {
	switch flags.OutputFormat {
	case "table", "yaml", "json":
	default:
		common.ExitWithErrorf("unsupported output format %s (support table, yaml, json)",
			flags.OutputFormat)
	}

	visitorBulder := util.NewVisitorBuilder()

	cmdArgs := cmd.Flags().Args()

	switch len(cmdArgs) {
	case 0:
		common.ExitWithErrorf("no resource specified")
	case 1:
		visitorBulder.CommandParam(&util.CommandOptions{
			Kind: cmdArgs[0],
		})
	case 2:
		visitorBulder.CommandParam(&util.CommandOptions{
			Kind: cmdArgs[0],
			Name: cmdArgs[1],
		})
	default:
		common.ExitWithErrorf("invalid command args: support <resource kind> [resource name]")
	}

	vss, err := visitorBulder.Do()
	if err != nil {
		common.ExitWithErrorf("build visitor failed: %s", err)
	}

	printer := printer.New(flags.OutputFormat)
	var errs []error
	for _, vs := range vss {
		err := vs.Visit(func(mo resource.MeshObject, e error) error {
			if e != nil {
				return errors.Wrap(e, "visit failed")
			}

			resourceID := mo.Kind()
			if mo.Name() != "" {
				resourceID += "/" + mo.Name()
			}

			objects, err := WrapGetterByMeshObject(mo, meshclient.New(flags.Server), flags.Timeout).Get()
			if err != nil {
				return errors.Wrapf(err, "%s get failed", resourceID)
			}

			printer.PrintObjects(objects)

			return nil
		})

		if err != nil {
			common.OutputError(err)
			errs = append(errs, err)
		}
	}

	if len(errs) > 0 {
		common.ExitWithErrorf("getting resources has errors occurred")
	}
}
