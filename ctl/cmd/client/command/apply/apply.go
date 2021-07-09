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
			common.OutputError(err)
			errs = append(errs, err)
		}
	}

	if len(errs) > 0 {
		common.ExitWithErrorf("applying resources has errors occurred")
	}
}
