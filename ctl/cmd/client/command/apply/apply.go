package apply

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
	vss, err := util.NewVisitorBuilder().
		FilenameParam(&util.FilenameOptions{Recursive: args.Recursive, Filenames: []string{args.YamlFile}}).
		Do()

	if err != nil {
		common.ExitWithErrorf("parse spec files error: %s", err)
	}
	var errs []error
	for _, vs := range vss {
		vs.Visit(func(mo resource.MeshObject, e error) error {
			if e != nil {
				common.OutputErrorInfo("visit an error %s", e)
				errs = append(errs, e)
				return nil
			}
			err := WrapApplierByMeshObject(mo, meshclient.New(args.Server), args.Timeout).Apply()
			if err != nil {
				errs = append(errs, err)
				common.OutputErrorInfo("apply %s of resource %s failed: %s\n", mo.Name(), mo.GetKind(), err)
			} else {
				fmt.Printf("%s of resource %s has been successfully applied\n", mo.Name(), mo.GetKind())
			}
			return nil
		})
	}

	if len(errs) > 0 {
		common.ExitWithErrorf("appling resources has errors occurred")
	}
}
