/*
 * Copyright (c) 2021, MegaEase
 * All rights reserved.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package fake

import (
	"github.com/megaease/easemeshctl/cmd/client/command/flags"
	installbase "github.com/megaease/easemeshctl/cmd/client/command/meshinstall/base"
	apiextensions "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	"k8s.io/client-go/kubernetes"
)

// NewStageContextForApply return a fake stageContext for apply subcommand
func NewStageContextForApply(client kubernetes.Interface,
	apiextensionsClient apiextensions.Interface,
) *installbase.StageContext {
	return &installbase.StageContext{
		Client:              client,
		APIExtensionsClient: apiextensionsClient,
		CoreDNSFlags:        &flags.CoreDNS{},
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
			MeshControlPlanePersistVolumeName:     installbase.ControlPlanePVCName,
			MeshControlPlanePersistVolumeHostPath: "/opt/easemesh",
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
				EgServiceName: "easemesh-control-plane-svc",
			},
		},
	}
}
