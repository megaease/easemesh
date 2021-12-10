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

package installbase

const (
	// ObjectsURL is url of objects.
	ObjectsURL = "/apis/v1/objects"
	// ObjectURL is url of object.
	ObjectURL = "/apis/v1/objects/%s"
	// MemberList is url of member list.
	MemberList = "/apis/v1/status/members"
)

const (
	// DefaultOperatorPath is path of default operation configuration.
	DefaultOperatorPath = "./manifests/easemesh-operator.yaml"

	// DefaultMeshControlPlaneName is the default control plane name of easemesh.
	DefaultMeshControlPlaneName = "easemesh-control-plane"
	// DefaultMeshClientPortName is the default port name of etcd client.
	DefaultMeshClientPortName = "client-port"
	// DefaultMeshPeerPortName is the default port name of etcd peer.
	DefaultMeshPeerPortName = "peer-port"
	// DefaultMeshAdminPortName is the default port name of etcd administration.
	DefaultMeshAdminPortName = "admin-port"
	// DefaultMeshControlPlanePlubicServiceName is the default exposed public service name of the EaseMesh.
	DefaultMeshControlPlanePlubicServiceName = "easemesh-controlplane-public"

	//DefaultMeshControlPlaneHeadlessServiceName is the default headless service name.
	DefaultMeshControlPlaneHeadlessServiceName = "easemesh-controlplane-hs"
	//DefaultMeshControlPlaneServicePeerPort  is the default value of etcd peer port.
	DefaultMeshControlPlaneServicePeerPort = 2380
	//DefaultMeshControlPlanelServiceAdminPort is the default value of etcd admin port.
	DefaultMeshControlPlanelServiceAdminPort = 2381

	//DefaultMeshControlPlanePVName is the default name of the persisten volume used by the control plane.
	DefaultMeshControlPlanePVName = "easegress-control-plane-pv"
	//DefaultMeshControlPlanePVHostPath is the default of path of pv.
	DefaultMeshControlPlanePVHostPath = "/opt/easemesh"
	//DefaultMeshControlPlaneConfig is the default configmap name of the easemesh control plane.
	DefaultMeshControlPlaneConfig = "easemesh-cluster-cm"

	//MeshControllerName is the mesh controller name in the easegress.
	MeshControllerName = "easemesh-controller"

	//DefaultMeshOperatorName is the default meshdeployment operator name of the EaseMesh.
	DefaultMeshOperatorName = "easemesh-operator"
	//DefaultMeshOperatorServiceName is the default service name of the meshdeployment operator.
	DefaultMeshOperatorServiceName = "easemesh-operator-service"
	//DefaultMeshOperatorSecretName is the default secret resource name of the meshdeployment operator.
	DefaultMeshOperatorSecretName = "easemesh-operator-secret"
	//DefaultMeshOperatorCSRName is the default CSR resource name of the meshdeployment operator.
	DefaultMeshOperatorCSRName = "easemesh-operator-csr"
	//DefaultMeshOperatorMutatingWebhookName  is the default operator mutating-webhook name of the adminission control.
	DefaultMeshOperatorMutatingWebhookName = "easemesh-operator-mutating-webhook"
	//DefaultMeshOperatorMutatingWebhookPath is the default path of the admission control for the EaseMesh.
	DefaultMeshOperatorMutatingWebhookPath = "/mutate"
	//DefaultMeshOperatorMutatingWebhookPort is the default port listened by the adminssion control for the EaseMesh.
	DefaultMeshOperatorMutatingWebhookPort = 9090
	//DefaultMeshOperatorCertDir is the default certs file localtion the adminssion control server.
	DefaultMeshOperatorCertDir = "/cert-volume"
	//DefaultMeshOperatorCertFileName is the default certs file name of the admission control server.
	DefaultMeshOperatorCertFileName = "cert.pem"
	//DefaultMeshOperatorKeyFileName is the default key name of the admission control server.
	DefaultMeshOperatorKeyFileName = "key.pem"

	// DefaultSidecarImageName is the default sidecar image name.
	DefaultSidecarImageName = "megaease/easegress:server-sidecar"
	//DefaultEaseagentInitializerImageName is the default easeagent initializer image name.
	DefaultEaseagentInitializerImageName = "megaease/easeagent-initializer:latest"
	//DefaultLog4jConfigName is the default log4j config file name of the easeagent.
	DefaultLog4jConfigName = "log4j2.xml"

	//DefaultMeshIngressConfig is the default configmap name of the meshingress.
	DefaultMeshIngressConfig = "easemesh-ingress-config"
	//DefaultMeshIngressService is the default service name of the meshingress.
	DefaultMeshIngressService = "easemesh-ingress-service"
	//DefaultMeshIngressControllerName is the default deployment name of the meshingress.
	DefaultMeshIngressControllerName = "easemesh-ingress-easegress"
	//DefaultShadowServiceControllerName is the default deployment name of the shadow service controller.
	DefaultShadowServiceControllerName = "easemesh-shadowservice-controller"

	// DefaultKubeDir represents default kubernetes client configuration directory.
	DefaultKubeDir = ".kube"

	//DefaultKubernetesConfig is the default config name of the K8s.
	DefaultKubernetesConfig = "config"
	//WriterClusterRole is the write role name of the Easegress.
	WriterClusterRole = "writer"
	//ReaderClusterRole is the read role name of the Easegress.
	ReaderClusterRole = "reader"
)

// InstallPhase is the phrase of installation.
type InstallPhase int

const (
	// BeginPhase is the phrase of beginning.
	BeginPhase InstallPhase = iota
	// EndPhase is the phrase of ending.
	EndPhase
	// ErrorPhase if the phrase of erroring handling.
	ErrorPhase
)
