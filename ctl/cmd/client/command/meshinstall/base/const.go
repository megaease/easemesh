package installbase

const (
	ObjectsURL = "/apis/v1/objects"
	ObjectURL  = "/apis/v1/objects/%s"
	MemberList = "/apis/v1/status/members"
)

const (
	DefaultOperatorPath = "./manifests/easemesh-operator.yaml"

	DefaultMeshControlPlaneName                = "easemesh-control-plane"
	DefaultMeshClientPortName                  = "client-port"
	DefaultMeshPeerPortName                    = "peer-port"
	DefaultMeshAdminPortName                   = "admin-port"
	DefaultMeshControlPlanePlubicServiceName   = "easemesh-controlplane-public"
	DefaultMeshControlPlaneHeadlessServiceName = "easemesh-controlplane-hs"
	DefaultMeshControlPlaneServicePeerPort     = 2380
	DefaultMeshControlPlanelServiceAdminPort   = 2381

	DefaultMeshControlPlanePVName     = "easegress-control-plane-pv"
	DefaultMeshControlPlanePVHostPath = "/opt/easemesh"
	DefaultMeshControlPlaneConfig     = "easemesh-cluster-cm"

	DefaultMeshControllerName = "easemesh-controller"

	DefaultMeshOperatorName                         = "easemesh-operator"
	DefaultMeshOperatorControllerManagerServiceName = "mesh-operator-controller-manager-metrics-service"

	DefaultMeshIngressConfig         = "easemesh-ingress-config"
	DefaultMeshIngressService        = "easemesh-ingress-service"
	DefaultMeshIngressControllerName = "easemesh-ingress-easegress"

	// DefaultKubeDir represents default kubernetes client configuration directory
	DefaultKubeDir = ".kube"

	//
	DefaultKubernetesConfig = "config"
	WriterClusterRole       = "writer"
	ReaderClusterRole       = "reader"
)

type InstallPhase int

const (
	BeginPhase InstallPhase = iota
	EndPhase
	ErrorPhase
)
