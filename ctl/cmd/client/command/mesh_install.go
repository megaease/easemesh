package command

import (
	"fmt"
	"io/ioutil"

	installbase "github.com/megaease/easemeshctl/cmd/client/command/meshinstall/base"
	"github.com/megaease/easemeshctl/cmd/client/command/meshinstall/controlpanel"
	"github.com/megaease/easemeshctl/cmd/client/command/meshinstall/crd"
	"github.com/megaease/easemeshctl/cmd/client/command/meshinstall/installation"
	"github.com/megaease/easemeshctl/cmd/client/command/meshinstall/meshingress"
	"github.com/megaease/easemeshctl/cmd/client/command/meshinstall/operator"
	"github.com/megaease/easemeshctl/cmd/common"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

func InstallCmd() *cobra.Command {
	iArgs := &installbase.InstallArgs{}
	cmd := &cobra.Command{
		Use:     "install",
		Short:   "Deploy infrastructure components of the EaseMesh",
		Long:    "",
		Example: "emctl install <args>",
		Run: func(cmd *cobra.Command, args []string) {
			if iArgs.SpecFile != "" {
				var buff []byte
				var err error
				buff, err = ioutil.ReadFile(iArgs.SpecFile)
				if err != nil {
					common.ExitWithErrorf("%s failed: %v", cmd.Short, err)
				}

				err = yaml.Unmarshal(buff, &iArgs)
				if err != nil {
					common.ExitWithErrorf("%s failed: %v", cmd.Short, err)
				}
			}
			install(cmd, iArgs)
		},
	}

	addInstallArgs(cmd, iArgs)
	return cmd
}

func addInstallArgs(cmd *cobra.Command, args *installbase.InstallArgs) {
	baseCmdArgs(cmd, args)
	cmd.Flags().IntVar(&args.EgClientPort, "mesh-control-plane-client-port", installbase.DefaultMeshClientPort, "Mesh control plane client port for remote accessing")
	cmd.Flags().IntVar(&args.EgAdminPort, "mesh-control-plane-admin-port", installbase.DefaultMeshAdminPort, "Port of mesh control plane admin for management")
	cmd.Flags().IntVar(&args.EgPeerPort, "mesh-control-plane-peer-port", installbase.DefaultMeshPeerPort, "Port of mesh control plane for consensus each other")
	cmd.Flags().IntVar(&args.MeshControlPlaneCheckHealthzMaxTime,
		"mesh-control-plane-check-healthz-max-time",
		installbase.DefaultMeshControlPlaneCheckHealthzMaxTime,
		"Max timeout in second for checking control panel component whether ready or not (default 60 seconds)")

	cmd.Flags().IntVar(&args.EgServicePeerPort, "mesh-control-plane-service-peer-port", installbase.DefaultMeshPeerPort, "")
	cmd.Flags().IntVar(&args.EgServiceAdminPort, "mesh-control-plane-service-admin-port", installbase.DefaultMeshAdminPort, "")

	// cmd.Flags().StringVar(&args.EGControlPlanePersistVolumeName, "eg-control-plane-pv-name", installbase.DefaultEgControlPlanePVName, egControlPlanePVNameHelpStr)
	// cmd.Flags().StringVar(&args.EGControlPlanePersistVolumeHostPath, "eg-control-plane-pv-hostpath", installbase.DefaultEgControlPlanePVHostPath, egControlPlanePVHostPathHelpStr)
	cmd.Flags().StringVar(&args.MeshControlPlaneStorageClassName, "mesh-storage-class-name", installbase.DefaultMeshControlPlaneStorageClassName, "")
	cmd.Flags().StringVar(&args.MeshControlPlanePersistVolumeCapacity, "mesh-control-plane-pv-capacity", installbase.DefaultMeshControlPlanePersistVolumeCapacity,
		installbase.MeshControlPlanePVNotExistedHelpStr)

	cmd.Flags().StringVar(&args.EaseMeshRegistryType, "registry-type", installbase.DefaultMeshRegistryType, installbase.MeshRegistryTypeHelpStr)
	cmd.Flags().IntVar(&args.HeartbeatInterval, "heartbeat-interval", installbase.DefaultHeartbeatInterval, "")

	cmd.Flags().StringVar(&args.ImageRegistryURL, "image-registry-url", installbase.DefaultImageRegistryURL, "")
	cmd.Flags().StringVar(&args.EasegressImage, "easegress-image", installbase.DefaultEasegressImage, "")
	cmd.Flags().StringVar(&args.EaseMeshOperatorImage, "easemesh-operator-image", installbase.DefaultEaseMeshOperatorImage, "")

	cmd.Flags().IntVar(&args.EasegressControlPlaneReplicas, "easemesh-control-plane-replicas", installbase.DefaultMeshControlPlaneReplicas, "")
	cmd.Flags().IntVar(&args.EasegressIngressReplicas, "easeemesh-ingress-replicas", installbase.DefaultMeshIngressReplicas, "")
	cmd.Flags().IntVar(&args.EaseMeshOperatorReplicas, "easemesh-operator-replicas", installbase.DefaultMeshOperatorReplicas, "")
	cmd.Flags().StringVarP(&args.SpecFile, "file", "f", "", "A yaml file specifying the install params.")
	cmd.Flags().BoolVar(&args.CleanWhenFailed, "clean-when-failed", true, "Clean resources when installation failed, default true")
}

func install(cmd *cobra.Command, args *installbase.InstallArgs) {
	var err error
	kubeClient, err := installbase.NewKubernetesClient()
	if err != nil {
		common.ExitWithErrorf("%s failed: %v", cmd.Short, err)
	}

	apiExtensionClient, err := installbase.NewKubernetesApiExtensionsClient()
	if err != nil {
		common.ExitWithErrorf("%s failed: %v", cmd.Short, err)
	}

	context := &installbase.StageContext{
		Arguments:           *args,
		Client:              kubeClient,
		Cmd:                 cmd,
		APIExtensionsClient: apiExtensionClient,
	}

	install := installation.New(
		installation.Wrap(crd.PreCheck, crd.Deploy, crd.Clear, crd.Describe),
		installation.Wrap(controlpanel.PreCheck, controlpanel.Deploy, controlpanel.Clear, controlpanel.Describe),
		installation.Wrap(operator.PreCheck, operator.Deploy, operator.Clear, operator.Describe),
		installation.Wrap(meshingress.PreCheck, meshingress.Deploy, meshingress.Clear, meshingress.Describe),
	)

	err = install.DoInstallStage(context)
	if err != nil {
		if args.CleanWhenFailed {
			install.ClearResource(context)
		}
		common.ExitWithErrorf("install mesh infrastructure error: %s", err)
	}

	fmt.Println("Done.")
}
