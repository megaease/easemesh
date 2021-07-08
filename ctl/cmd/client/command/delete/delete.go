package delete

import (
	"fmt"

	"github.com/megaease/easemeshctl/cmd/client/command/flags"
	"github.com/megaease/easemeshctl/cmd/client/command/meshclient"
	"github.com/megaease/easemeshctl/cmd/client/resource"
	"github.com/megaease/easemeshctl/cmd/client/util"
	"github.com/megaease/easemeshctl/cmd/common"

	"github.com/spf13/cobra"
)

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
				return fmt.Errorf("visit failed: %v", e)
			}

			err := WrapDeleterByMeshObject(mo, meshclient.New(flags.Server), flags.Timeout).Delete()
			if err != nil {
				return fmt.Errorf("%s/%s deleted failed: %s", mo.Kind(), mo.Name(), err)
			}

			fmt.Printf("%s/%s deleted successfully\n", mo.Kind(), mo.Name())
			return nil
		})

		if err != nil {
			common.OutputError(err)
			errs = append(errs, err)
		}
	}

	if len(errs) > 0 {
		common.ExitWithErrorf("deleting resources has errors occurred")
	}
}
