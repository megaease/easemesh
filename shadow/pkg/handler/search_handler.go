package handler

import (
	"log"
	"time"

	"github.com/megaease/easemesh/mesh-shadow/pkg/common"
	"github.com/megaease/easemesh/mesh-shadow/pkg/object"
	"github.com/megaease/easemesh/mesh-shadow/pkg/utils"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	runTimeClient "sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/pkg/errors"
)

var (
	// ConflictError indicate that the resource already exists
	ConflictError = errors.Errorf("resource already exists")
	// NotFoundError indicate that the resource does not exist
	NotFoundError = errors.Errorf("resource not found")
)

const (
	MeshServiceAnnotation   = "mesh.megaease.com/service-name"
	apiURL                  = "/apis/v1"
	MeshShadowServicesURL   = apiURL + "/mesh/shadowservices"
	MeshCustomObjetWatchURL = apiURL + "/mesh/watchCustomObjects/{kind}"
	MeshCustomObjectsURL    = apiURL + "/mesh/customObjects/{kind}"
)

type SearchHandler struct {
	KubeClient    *kubernetes.Clientset
	RunTimeClient *runTimeClient.Client
	CRDClient     *rest.RESTClient
	MeshServer    string
	CloneChan     chan interface{}

	Interval time.Duration
	Registry *common.CallbackRegistry
}

type ServiceCloneBlock struct {
	service   object.ShadowService
	deployObj interface{}
}

func (searcher *SearchHandler) SearchOriginDeployments(shadowService object.ShadowService) {

	namespace := shadowService.Namespace
	serviceName := shadowService.ServiceName
	meshDeploymentList, err := utils.ListMeshDeployment(namespace, searcher.CRDClient, metav1.ListOptions{})
	if err != nil {
		log.Printf("Query MeshDeployment for shadow service error. %s", err)
	}

	for _, meshDeployment := range meshDeploymentList.Items {
		if meshDeployment.Spec.Service.Name == serviceName {
			searcher.CloneChan <- ServiceCloneBlock{
				shadowService,
				meshDeployment,
			}
		}
	}

	deployments, err := utils.ListDeployments(namespace, searcher.KubeClient, metav1.ListOptions{})
	if err != nil {
		log.Printf("Query Deployment for shadow service error. %s", err)
	}

	for _, deployment := range deployments {
		annotations := deployment.Annotations
		if _, ok := annotations[MeshServiceAnnotation]; ok {
			if serviceName == annotations[MeshServiceAnnotation] {
				searcher.CloneChan <- ServiceCloneBlock{
					shadowService,
					deployment,
				}
			}
		}
	}
	log.Printf("The service doesn't have MeshDeployment or Deployment for run it. Service: %s", serviceName)
}
