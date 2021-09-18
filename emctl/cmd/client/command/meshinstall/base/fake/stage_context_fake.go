package fake

import (
	"github.com/megaease/easemeshctl/cmd/client/command/flags"
	installbase "github.com/megaease/easemeshctl/cmd/client/command/meshinstall/base"
	apiextensions "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	"k8s.io/client-go/kubernetes"
)

// NewStageContextForApply return a fake stageContext for apply subcommand
func NewStageContextForApply(client kubernetes.Interface,
	apiextensionsClient apiextensions.Interface) *installbase.StageContext {
	return &installbase.StageContext{
		Client:              client,
		APIExtensionsClient: apiextensionsClient,
		Flags: &flags.Install{
			// Easegress Control Plane params
			EasegressImage:                "megaease/easegress",
			EasegressControlPlaneReplicas: 3,

			EgClientPort:       2379,
			EgAdminPort:        2380,
			EgPeerPort:         2381,
			EgServicePeerPort:  2380,
			EgServiceAdminPort: 2381,

			MeshControlPlaneStorageClassName:      "easemesh-storage-class",
			MeshControlPlanePersistVolumeName:     installbase.DefaultMeshControlPlanePVName,
			MeshControlPlanePersistVolumeHostPath: installbase.DefaultMeshControlPlanePVHostPath,
			MeshControlPlanePersistVolumeCapacity: "3Gi",
			MeshControlPlaneCheckHealthzMaxTime:   1000,

			MeshIngressReplicas:    1,
			MeshIngressServicePort: 19527,
			EaseMeshRegistryType:   "eureka",
			HeartbeatInterval:      30000,

			// EaseMesh Operator params
			EaseMeshOperatorImage:    "megaease/easemesh-operator",
			EaseMeshOperatorReplicas: 1,
			SpecFile:                 "",

			OperationGlobal: &flags.OperationGlobal{
				MeshNamespace: "easemesh",
				EgServiceName: "easemesh-controlplane-svc",
			},
		},
	}
}
