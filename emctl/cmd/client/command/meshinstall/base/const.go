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
	// --- Easegress itself related.

	// EasegressPrimaryClusterRole is the primary role name of Easegress.
	EasegressPrimaryClusterRole = "primary"
	// EasegressSecondaryClusterRole is the secondary role name of Easegress.
	EasegressSecondaryClusterRole = "secondary"

	// --- Control plane config related to path (putting them together makes it clearer).

	// ControlPlaneConfigMapName is the name of config map of control plane.
	ControlPlaneConfigMapName = "easemesh-control-plane-config"
	// ControlPlaneConfigMapKey is the key of data of config map of control plane.
	ControlPlaneConfigMapKey = "control-plane.yaml"
	// ControlPlaneConfigMapVolumeMountPath is the path of volume mouth of config map of control plane.
	ControlPlaneConfigMapVolumeMountPath = "/opt/easegress/config/control-plane.yaml"
	// ControlPlaneConfigMapVolumeMountSubPath is the subpath of volume mouth of config map of control plane.
	ControlPlaneConfigMapVolumeMountSubPath = "control-plane.yaml"
	// ControlPlaneHomeDir is home directory of control plane.
	ControlPlaneHomeDir = "/opt/easegress"
	// ControlPlaneDataDir is data directory of control plane.
	ControlPlaneDataDir = "/opt/easegress/control-plane-data"
	// ControlPlaneCmd is the essential command of control plane.
	ControlPlaneCmd = "/opt/easegress/bin/easegress-server -f /opt/easegress/config/control-plane.yaml"

	// --- Control plane StatefuleSet related.

	// ControlPlaneStatefulSetName is the name of control plane statefulset.
	ControlPlaneStatefulSetName = "easemesh-control-plane"
	// ControlPlaneStatefulSetClientPortName is the name of client port.
	ControlPlaneStatefulSetClientPortName = "client-port"
	// ControlPlaneStatefulSetPeerPortName is the name of peer port.
	ControlPlaneStatefulSetPeerPortName = "peer-port"
	// ControlPlaneStatefulSetAdminPortName is the name of admin port.
	ControlPlaneStatefulSetAdminPortName = "admin-port"
	// ControlPlanePVCName is the name of persisten volume claim control plane.
	ControlPlanePVCName = "control-plane-pvc"

	// --- Control Plane Service related.

	// ControlPlanePlubicServiceName is name of public service of control plane.
	ControlPlanePlubicServiceName = "easemesh-control-plane-public"
	// ControlPlaneHeadlessServiceName is name of headless service of control plane.
	ControlPlaneHeadlessServiceName = "easemesh-control-plane-hs"

	// --- Sidecar related.

	// SidecarHomeDir is the directory of sidecar.
	SidecarHomeDir = "/opt/easegress"

	// --- MeshController related.

	// MeshControllerName is the name of MeshController in EaseMesh.
	MeshControllerName = "easemesh-controller"
	// MeshControllerAPIPort is the API port of sidecar for handling local Eureka/Conslu/Nacos APIs.
	MeshControllerAPIPort = 13009

	// --- Operator Deployment related.

	// OperatorConfigMapName is the name of config map of operator.
	OperatorConfigMapName = "easemesh-operator-config"
	// OperatorConfigMapKey is the key of config map of operator.
	OperatorConfigMapKey = "operator.yaml"
	// OperatorConfigMapVolumeMountPath is the path of volume mouth of config map of operator.
	OperatorConfigMapVolumeMountPath = "/opt/operator/operator.yaml"
	// OperatorConfigMapVolumeMountSubPath is the subpath of volume mouth of config map in control plane config map.
	OperatorConfigMapVolumeMountSubPath = "operator.yaml"
	// OperatorSecretName is the name of secret of operator deployment.
	OperatorSecretName = "easemesh-operator-secret"
	// OperatorSecretVolumeMountPath is the secret directory of adminssion control of operator deployment.
	OperatorSecretVolumeMountPath = "/opt/operator/cert-volume"
	// OperatorSecretCertFileName is the cert filename of admission control of operator deployment.
	OperatorSecretCertFileName = "cert.pem"
	// OperatorSecretKeyFileName is the key filename of admission control of operator deployment.
	OperatorSecretKeyFileName = "key.pem"
	// OperatorCmd is the command of operator.
	OperatorCmd = "/manager"
	// OperatorArgs is the args of operator.
	OperatorArgs = "--config=/opt/operator/operator.yaml"

	// OperatorDeploymentName is the name of operator deployment.
	OperatorDeploymentName = "easemesh-operator"
	// OperatorServiceName is the name of service of operator deployment.
	OperatorServiceName = "easemesh-operator-service"
	// OperatorCSRName is the name of CertificateSigningRequest of operator deployment.
	OperatorCSRName = "easemesh-operator-csr"
	// OperatorMutatingWebhookName is the name of mutating-webhook of admission control of operator deployment.
	OperatorMutatingWebhookName = "easemesh-operator-mutating-webhook"
	// OperatorMutatingWebhookPath is the path of admission control of operator deployment.
	OperatorMutatingWebhookPath = "/mutate"
	// OperatorMutatingWebhookPortName is the name of mutating webhook port of admission control of operator deployment.
	OperatorMutatingWebhookPortName = "mutate-port"
	// OperatorMutatingWebhookPort is the port of adminssion control of operator deployment.
	OperatorMutatingWebhookPort = 9090

	// --- Operator injection related.

	// SidecarImageName is the imaget name of sidecar.
	SidecarImageName = "megaease/easegress:easemesh"
	// AgentInitializerImageName is the image name of agent initializer.
	AgentInitializerImageName = "megaease/easeagent-initializer:latest"
	// AgentLog4jConfigName is the file name of log4j config of agent.
	AgentLog4jConfigName = "log4j2.xml"

	// --- Ingress Controller related.

	// IngressControllerDeploymentName is the name of deployment of ingress controller.
	IngressControllerDeploymentName = "easemesh-ingress-controller"
	// IngressControllerDeploymentCmd is the essetial command of deployment of ingress controller.
	IngressControllerDeploymentCmd = "/opt/easegress/bin/easegress-server -f /opt/easegress/config/ingress-controller.yaml"
	// IngressControllerConfigMapName is the name of config map of ingress controller.
	IngressControllerConfigMapName = "easemesh-ingress-controller-config"
	// IngressControllerServiceName is the name of service of ingress controller
	IngressControllerServiceName = "easemesh-ingress-controller-service"
	// IngressControllerConfigMapKey is the key of data of config map of ingress controller.
	IngressControllerConfigMapKey = "ingress-controller.yaml"
	// IngressControllerConfigMapVolumeMountPath is the path of volume mouth of config map of ingress controller.
	IngressControllerConfigMapVolumeMountPath = "/opt/easegress/config/ingress-controller.yaml"
	// IngressControllerConfigMapVolumeMountSubPath is the subpath of volume mouth of config map of ingress controller.
	IngressControllerConfigMapVolumeMountSubPath = "control-plane.yaml"
	// IngressControllerHomeDir is home directory of control plane.
	IngressControllerHomeDir = "/opt/easegress"

	// --- Shadow Service related.

	// IngressControllerShadowServiceName is the name of shadow service of ingress controller.
	IngressControllerShadowServiceName = "easemesh-ingress-controller-shadowservice"

	// --- Kubernetes related.

	// DefaultKubeDir is the directory of Kubernetes config.
	DefaultKubeDir = ".kube"
	// DefaultKubernetesFilename is the file name of Kubernetes config.
	DefaultKubernetesFilename = "config"
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
