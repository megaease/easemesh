package command

import (
	installbase "github.com/megaease/easemeshctl/cmd/client/command/meshinstall/base"

	"github.com/spf13/cobra"
)

func baseCmdArgs(cmd *cobra.Command, args *installbase.InstallArgs) {
	cmd.Flags().StringVar(&args.MeshNameSpace, "mesh-namespace", installbase.DefaultMeshNameSpace, "")
	cmd.Flags().StringVar(&args.EgServiceName, "mesh-control-plane-service-name", installbase.DefaultMeshControlPlaneHeadfulServiceName, "")
}
