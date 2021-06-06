package installbase

const (
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
    type: "DirectoryOrCreate"
    `
)

const (
	ObjectsURL = "/apis/v1/objects"
	ObjectURL  = "/apis/v1/objects/%s"
	MemberList = "/apis/v1/status/members/"
)

const (
	DefaultOperatorPath = "./manifests/easemesh-operator.yaml"
	// DefaultMeshNamespace all mesh infrastructure component should be deployed in this namespace
	DefaultMeshNameSpace = "easemesh"

	DefaultMeshControlPlaneName = "easemesh-control-plane"

	DefaultMeshControlPlaneReplicas = 3
	DefaultMeshIngressReplicas      = 1
	DefaultMeshOperatorReplicas     = 1
	DefaultMeshClientPortName       = "client-port"
	DefaultMeshPeerPortName         = "peer-port"
	DefaultMeshAdminPortName        = "admin-port"
	DefaultMeshClientPort           = 2379
	DefaultMeshPeerPort             = 2380
	DefaultMeshAdminPort            = 2381

	DefaultMeshControlPlanePlubicServiceName   = "easemesh-controlplane-public"
	DefaultMeshControlPlaneHeadlessServiceName = "easemesh-controlplane-hs"
	DefaultMeshControlPlaneHeadfulServiceName  = "easemesh-controlplane-svc"
	DefaultMeshControlPlaneServicePeerPort     = 2380
	DefaultMeshControlPlanelServiceAdminPort   = 2381
	DefaultMeshControlPlaneCheckHealthzMaxTime = 60

	DefaultMeshControlPlaneStorageClassName      = "easemesh-storage"
	DefaultMeshControlPlanePVName                = "easegress-control-plane-pv"
	DefaultMeshControlPlanePVHostPath            = "/opt/easemesh"
	DefaultMeshControlPlanePersistVolumeCapacity = "3Gi" // 3 Gib
	DefaultMeshControlPlaneConfig                = "easemesh-cluster-cm"

	// EaseMesh Controller default Params
	DefaultMeshRegistryType   = EurekaRegistryType
	DefaultHeartbeatInterval  = 5
	DefaultMeshControllerName = "easemesh-controller"
	MeshControllerKind        = "MeshController"

	DefaultMeshOperatorName                         = "easemesh-operator"
	DefaultMeshOperatorControllerManagerServiceName = "mesh-operator-controller-manager-metrics-service"

	// Ingress resource name
	DefaultMeshIngressConfig         = "easemesh-ingress-config"
	DefaultMeshIngressService        = "easemesh-ingress-service"
	DefaultMeshIngressControllerName = "easemesh-ingress-easegress"
)

const (
	// EurekaRegistryType is a constant that represents the default service registry type is eureka
	EurekaRegistryType string = "eureka"

	// DefaultKubeDir represents default kubernetes client configuration directory
	DefaultKubeDir = ".kube"

	//
	DefaultKubernetesConfig = "config"
	WriterClusterRole       = "writer"
	ReaderClusterRole       = "reader"
)

const (
	DefaultEasegressImage        = "megaease/easegress:latest"
	DefaultEaseMeshOperatorImage = "megaease/easemesh-operator:latest"
	DefaultImageRegistryURL      = "docker.io"
)

type InstallPhase int

const (
	BeginPhase InstallPhase = iota
	EndPhase
	ErrorPhase
)
