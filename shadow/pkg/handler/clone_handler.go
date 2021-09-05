package handler

import (
	"log"

	"github.com/megaease/easemesh/mesh-shadow/pkg/object/v1beta1"
	appv1 "k8s.io/api/apps/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const (

	// Init container stuff.
	initContainerName = "initializer"

	agentVolumeName   = "agent-volume"
	sidecarVolumeName = "sidecar-volume"

	// Sidecar container stuff.
	sidecarContainerName = "easemesh-sidecar"

	shadowLabelKey            = "mesh.megaease.com/shadow-service"
	shadowAppContainerNameKey = "mesh.megaease.com/app-container-name"

	shadowDeploymentNameSuffix = "-shadow"

	mysqlShadowConfigEnv         = "EASE_MYSQL_CONFIG"
	kafkaShadowConfigEnv         = "EASE_KAFKA_CONFIG"
	rabbitmqShadowConfigEnv      = "EASE_RABBITMQ_CONFIG"
	redisShadowConfigEnv         = "EASE_REDIS_CONFIG"
	elasticsearchShadowConfigEnv = "EASE_ELASTICSEARCH_CONFIG"
)

type CloneHandler struct {
	KubeClient    *kubernetes.Clientset
	RunTimeClient *client.Client
	CRDClient     *rest.RESTClient
}

func (handler *CloneHandler) Clone(obj interface{}) {

	var err error
	block := obj.(ServiceCloneBlock)
	switch block.deployObj.(type) {
	case appv1.Deployment:
		deployment := block.deployObj.(appv1.Deployment)
		err = handler.CloneDeployment(&deployment, &block.service)()
	case v1beta1.MeshDeployment:
		meshDeployment := block.deployObj.(v1beta1.MeshDeployment)
		err = handler.CloneMeshDeployment(&meshDeployment, &block.service)()
	}
	if err != nil {
		log.Printf("Create shadow service failed. service: %s error: %s", block.service.ServiceName, err)
	}
}
