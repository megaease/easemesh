package delete

import (
	"fmt"
	"time"

	"github.com/megaease/easemeshctl/cmd/client/command/meshclient"
	"github.com/megaease/easemeshctl/cmd/client/resource"
	"github.com/megaease/easemeshctl/cmd/client/util"
	"github.com/megaease/easemeshctl/cmd/common"

	"github.com/spf13/cobra"
)

type Arguments struct {
	Server    string
	YamlFile  string
	Recursive bool
	Timeout   time.Duration
}

func Run(cmd *cobra.Command, args *Arguments) {
	visitorBulder := util.NewVisitorBuilder()

	cmdArgs := cmd.Flags().Args()

	if len(cmdArgs) == 0 && args.YamlFile == "" {
		common.ExitWithErrorf("no resource specified")
	}

	if len(cmdArgs) != 0 {
		if args.YamlFile != "" {
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

	if args.YamlFile != "" {
		visitorBulder.FilenameParam(&util.FilenameOptions{
			Recursive: args.Recursive,
			Filenames: []string{args.YamlFile},
		})
	}

	vss, err := visitorBulder.Do()
	if err != nil {
		common.ExitWithErrorf("build visitor failed: %s", err)
	}

	var errs []error
	for _, vs := range vss {
		vs.Visit(func(mo resource.MeshObject, e error) error {
			if e != nil {
				common.OutputErrorf("visit failed: %v", e)
				errs = append(errs, e)
				return nil
			}

			err := WrapDeleterByMeshObject(mo, meshclient.New(args.Server), args.Timeout).Delete()
			if err != nil {
				errs = append(errs, err)
				common.OutputErrorf("%s/%s deleted failed: %s\n", mo.Kind(), mo.Name(), err)
			} else {
				fmt.Printf("%s/%s deleted successfully\n", mo.Kind(), mo.Name())
			}
			return nil
		})
	}

	if len(errs) > 0 {
		common.ExitWithErrorf("deleting resources has errors occurred")
	}
}
