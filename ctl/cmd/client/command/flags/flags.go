package flags

import (
	"time"

	"github.com/spf13/cobra"
)

const (
	// DefaultMeshNamespace all mesh infrastructure component should be deployed in this namespace
	DefaultMeshNameSpace = "easemesh"

	DefaultMeshControlPlaneReplicas = 3
	DefaultMeshIngressReplicas      = 1
	DefaultMeshOperatorReplicas     = 1
	DefaultMeshClientPort           = 2379
	DefaultMeshPeerPort             = 2380
	DefaultMeshAdminPort            = 2381

	DefaultMeshControlPlaneHeadfulServiceName  = "easemesh-controlplane-svc"
	DefaultMeshControlPlaneCheckHealthzMaxTime = 60

	DefaultMeshControlPlaneStorageClassName      = "easemesh-storage"
	DefaultMeshControlPlanePersistVolumeCapacity = "3Gi" // 3 Gib

	// EaseMesh Controller default Params
	DefaultMeshRegistryType  = "eureka"
	DefaultHeartbeatInterval = 5
	MeshControllerKind       = "MeshController"

	// Ingress resource name
	DefaultMeshIngressServicePort = 19527

	MeshControlPlanePVNameHelpStr     = "The name of PersistentVolume for EaseMesh control plane storage."
	MeshControlPlanePVHostPathHelpStr = "The host path of the PersistentVolume for EaseMesh control plane storage."
	MeshControlPlanePVCapacityHelpStr = "The capacity of the PersistentVolume for EaseMesh control plane storage."

	MeshRegistryTypeHelpStr = "The registry type for application service registry, one of: eureka|consul|nacos."

	MeshControlPlaneStartupFailedHelpStr = `
		EaseMesh Control Plane deploy failed. Please check the K8S resource under the %s namespace for finding errors as follows:

		$ kubectl get statefulsets.apps -n %s

		$ kubectl get pods -n %s
	`
	MeshControlPlanePVNotExistedHelpStr = `

PersistentVolume does not have enough resources, the required number is %d, but only %d.
EaseMesh control plane needs PersistentVolume to store data. You need to create PersistentVolume in advance and specify its storageClassName as %s.

You can create PersistentVolume by the following definition:

apiVersion: v1
kind: PersistentVolume
metadata:
  labels:
    app: easemesh
  name: easemesh-pv
spec:
  storageClassName: %s
  accessModes:
  - {ReadWriteOnce}
  capacity:
    storage: {%s}
  hostPath:
    path: {/opt/easemesh/}
    type: "DirectoryOrCreate"`

	DefaultEasegressImage        = "megaease/easegress:latest"
	DefaultEaseMeshOperatorImage = "megaease/easemesh-operator:latest"
	DefaultImageRegistryURL      = "docker.io"
)

type (
	OperationGlobal struct {
		MeshNameSpace string
		EgServiceName string
	}

	Install struct {
		*OperationGlobal

		ImageRegistryURL string

		CleanWhenFailed bool

		// Easegress Control Plane params
		EasegressImage                string
		EasegressControlPlaneReplicas int
		EgClientPort                  int
		EgAdminPort                   int
		EgPeerPort                    int

		EgServicePeerPort  int
		EgServiceAdminPort int

		MeshControlPlaneStorageClassName      string
		MeshControlPlanePersistVolumeName     string
		MeshControlPlanePersistVolumeHostPath string
		MeshControlPlanePersistVolumeCapacity string
		MeshControlPlaneCheckHealthzMaxTime   int

		MeshIngressReplicas    int
		MeshIngressServicePort int32

		// EaseMesh Controller  params
		EaseMeshRegistryType string
		HeartbeatInterval    int

		// EaseMesh Operator params
		EaseMeshOperatorImage    string
		EaseMeshOperatorReplicas int

		SpecFile string
	}

	Reset struct {
		*OperationGlobal
	}

	AdminGlobal struct {
		Server  string
		Timeout time.Duration
	}

	AdminFileInput struct {
		YamlFile  string
		Recursive bool
	}

	Apply struct {
		*AdminGlobal
		*AdminFileInput
	}

	Delete struct {
		*AdminGlobal
		*AdminFileInput
	}

	Get struct {
		*AdminGlobal
		OutputFormat string
	}
)

func (i *Install) AttachCmd(cmd *cobra.Command) {
	i.OperationGlobal = &OperationGlobal{}
	i.OperationGlobal.AttachCmd(cmd)
	cmd.Flags().IntVar(&i.EgClientPort, "mesh-control-plane-client-port", DefaultMeshClientPort, "Mesh control plane client port for remote accessing")
	cmd.Flags().IntVar(&i.EgAdminPort, "mesh-control-plane-admin-port", DefaultMeshAdminPort, "Port of mesh control plane admin for management")
	cmd.Flags().IntVar(&i.EgPeerPort, "mesh-control-plane-peer-port", DefaultMeshPeerPort, "Port of mesh control plane for consensus each other")
	cmd.Flags().IntVar(&i.MeshControlPlaneCheckHealthzMaxTime,
		"mesh-control-plane-check-healthz-max-time",
		DefaultMeshControlPlaneCheckHealthzMaxTime,
		"Max timeout in second for checking control panel component whether ready or not (default 60 seconds)")

	cmd.Flags().IntVar(&i.EgServicePeerPort, "mesh-control-plane-service-peer-port", DefaultMeshPeerPort, "")
	cmd.Flags().IntVar(&i.EgServiceAdminPort, "mesh-control-plane-service-admin-port", DefaultMeshAdminPort, "")

	// cmd.Flags().StringVar(&i.EGControlPlanePersistVolumeName, "eg-control-plane-pv-name", DefaultEgControlPlanePVName, egControlPlanePVNameHelpStr)
	// cmd.Flags().StringVar(&i.EGControlPlanePersistVolumeHostPath, "eg-control-plane-pv-hostpath", DefaultEgControlPlanePVHostPath, egControlPlanePVHostPathHelpStr)
	cmd.Flags().StringVar(&i.MeshControlPlaneStorageClassName, "mesh-storage-class-name", DefaultMeshControlPlaneStorageClassName, "")
	cmd.Flags().StringVar(&i.MeshControlPlanePersistVolumeCapacity, "mesh-control-plane-pv-capacity", DefaultMeshControlPlanePersistVolumeCapacity,
		MeshControlPlanePVNotExistedHelpStr)

	cmd.Flags().Int32Var(&i.MeshIngressServicePort, "mesh-ingress-service-port", DefaultMeshIngressServicePort, "A port on which mesh ingress controller listening")

	cmd.Flags().StringVar(&i.EaseMeshRegistryType, "registry-type", DefaultMeshRegistryType, MeshRegistryTypeHelpStr)
	cmd.Flags().IntVar(&i.HeartbeatInterval, "heartbeat-interval", DefaultHeartbeatInterval, "")

	cmd.Flags().StringVar(&i.ImageRegistryURL, "image-registry-url", DefaultImageRegistryURL, "")
	cmd.Flags().StringVar(&i.EasegressImage, "easegress-image", DefaultEasegressImage, "")
	cmd.Flags().StringVar(&i.EaseMeshOperatorImage, "easemesh-operator-image", DefaultEaseMeshOperatorImage, "")

	cmd.Flags().IntVar(&i.EasegressControlPlaneReplicas, "easemesh-control-plane-replicas", DefaultMeshControlPlaneReplicas, "")
	cmd.Flags().IntVar(&i.MeshIngressReplicas, "easemesh-ingress-replicas", DefaultMeshIngressReplicas, "")
	cmd.Flags().IntVar(&i.EaseMeshOperatorReplicas, "easemesh-operator-replicas", DefaultMeshOperatorReplicas, "")
	cmd.Flags().StringVarP(&i.SpecFile, "file", "f", "", "A yaml file specifying the install params.")
	cmd.Flags().BoolVar(&i.CleanWhenFailed, "clean-when-failed", true, "Clean resources when installation failed, default true")
}

func (r *Reset) AttachCmd(cmd *cobra.Command) {
	r.OperationGlobal = &OperationGlobal{}
	r.OperationGlobal.AttachCmd(cmd)
}

func (o *OperationGlobal) AttachCmd(cmd *cobra.Command) {
	cmd.Flags().StringVar(&o.MeshNameSpace, "mesh-namespace", DefaultMeshNameSpace, "EaseMesh namespace in kubernetes")
	cmd.Flags().StringVar(&o.EgServiceName, "mesh-control-plane-service-name", DefaultMeshControlPlaneHeadfulServiceName, "")
}

func (a *AdminGlobal) AttachCmd(cmd *cobra.Command) {
	cmd.Flags().StringVarP(&a.Server, "server", "s", "127.0.0.1:2381", "An address to access the EaseMesh control plane")
	cmd.Flags().DurationVarP(&a.Timeout, "timeout", "t", 30*time.Second, "A duration that limit max time out for requesting the EaseMesh control plane")
}

func (a *AdminFileInput) AttachCmd(cmd *cobra.Command) {
	cmd.Flags().StringVarP(&a.YamlFile, "file", "f", "", "A location contained the EaseMesh resource files (YAML format) to apply, could be a file, directory, or URL")
	cmd.Flags().BoolVarP(&a.Recursive, "recursive", "r", true, "Whether to recursively iterate all sub-directories and files of the location")
}

func (a *Apply) AttachCmd(cmd *cobra.Command) {
	a.AdminGlobal = &AdminGlobal{}
	a.AdminGlobal.AttachCmd(cmd)

	a.AdminFileInput = &AdminFileInput{}
	a.AdminFileInput.AttachCmd(cmd)
}

func (d *Delete) AttachCmd(cmd *cobra.Command) {
	d.AdminGlobal = &AdminGlobal{}
	d.AdminGlobal.AttachCmd(cmd)

	d.AdminFileInput = &AdminFileInput{}
	d.AdminFileInput.AttachCmd(cmd)
}

func (g *Get) AttachCmd(cmd *cobra.Command) {
	g.AdminGlobal = &AdminGlobal{}
	g.AdminGlobal.AttachCmd(cmd)

	cmd.Flags().StringVarP(&g.OutputFormat, "output", "o", "table", "Output format (support table, yaml, json)")
}
