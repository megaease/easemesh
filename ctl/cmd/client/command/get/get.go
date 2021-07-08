package get

import (
	"fmt"

	"github.com/megaease/easemeshctl/cmd/client/command/flags"
	"github.com/megaease/easemeshctl/cmd/client/command/meshclient"
	"github.com/megaease/easemeshctl/cmd/client/command/printer"
	"github.com/megaease/easemeshctl/cmd/client/resource"
	"github.com/megaease/easemeshctl/cmd/client/util"
	"github.com/megaease/easemeshctl/cmd/common"

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
				return fmt.Errorf("visit failed: %v", e)
			}

			resourceID := mo.Kind()
			if mo.Name() != "" {
				resourceID += "/" + mo.Name()
			}

			objects, err := WrapGetterByMeshObject(mo, meshclient.New(flags.Server), flags.Timeout).Get()
			if err != nil {
				return fmt.Errorf("%s get failed: %s", resourceID, err)
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
