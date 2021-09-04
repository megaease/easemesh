package handler

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/megaease/easemesh/mesh-shadow/pkg/common"
	"github.com/megaease/easemesh/mesh-shadow/pkg/common/client"
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
	// NotFoundError indicate that the resource does not existed
	NotFoundError = errors.Errorf("resource not found")
)

const (
	MeshServiceAnnotation = "mesh.megaease.com/service-name"
	apiURL                = "/apis/v1"
	MeshShadowServicesURL = apiURL + "/mesh/shadowservices"
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

func (searcher *SearchHandler) Start() error {
	searcher.Registry.Add("", nil, searcher.Interval, searcher.queryDeploymentForShadowServiceFunc())
	return nil
}

func (searcher *SearchHandler) queryDeploymentForShadowServiceFunc() common.CallbackFunc {
	return func(context map[string]string, executeContext map[string]interface{}, interval time.Duration) bool {
		shadowServices, err := searcher.queryShadowServices()
		// an error occurs in retrieving metrics process this time, just ignore it.
		if err != nil {
			log.Printf("Query shadow service from easemesh control plane error %s", err)
		}
		searcher.queryOriginDeployments(shadowServices)
		return true
	}
}

func (searcher *SearchHandler) queryOriginDeployments(shadowServices []object.ShadowService) {

	namespaceMap := make(map[string][]*object.ShadowService)

	for _, shadowService := range shadowServices {
		services, ok := namespaceMap[shadowService.Namespace]
		if ok {
			services = append(services, &shadowService)
			namespaceMap[shadowService.Namespace] = services
		} else {
			services = []*object.ShadowService{}
			services = append(services, &shadowService)
			namespaceMap[shadowService.Namespace] = services
		}
	}

	for namespace, services := range namespaceMap {
		serviceNameMap := make(map[string]*object.ShadowService)
		for _, service := range services {
			serviceNameMap[service.ServiceName] = service
		}
		meshDeploymentList, err := utils.ListMeshDeployment(namespace, searcher.CRDClient, metav1.ListOptions{})
		if err != nil {
			log.Printf("Query MeshDeployment for shadow service error. %s", err)
		}

		for _, meshDeployment := range meshDeploymentList.Items {
			if _, ok := serviceNameMap[meshDeployment.Spec.Service.Name]; ok {
				searcher.CloneChan <- ServiceCloneBlock{
					*serviceNameMap[meshDeployment.Spec.Service.Name],
					meshDeployment,
				}
				delete(serviceNameMap, meshDeployment.Spec.Service.Name)
			}
		}

		deployments, err := utils.ListDeployments(namespace, searcher.KubeClient, metav1.ListOptions{})
		if err != nil {
			// logger.Errorf("Query Deployment for shadow service error.", err)
		}

		for _, deployment := range deployments {
			annotations := deployment.Annotations
			if _, ok := annotations[MeshServiceAnnotation]; ok {
				serviceName := annotations[MeshServiceAnnotation]
				if _, ok = serviceNameMap[serviceName]; ok {
					searcher.CloneChan <- ServiceCloneBlock{
						*serviceNameMap[serviceName],
						deployment,
					}
					delete(serviceNameMap, serviceName)
				}
			}
		}

		for serviceName, _ := range serviceNameMap {
			log.Printf("The service doesn't have MeshDeployment or Deployment for run it. Service: %s", serviceName)
		}
	}
}

func (searcher *SearchHandler) queryShadowServices() ([]object.ShadowService, error) {
	shadowServices, err := searcher.listShadowService(context.TODO())
	if err != nil {
		return nil, err
	}
	return shadowServices, nil

}

func (searcher *SearchHandler) listShadowService(ctx context.Context) ([]object.ShadowService, error) {
	jsonClient := client.NewHTTPJSON()
	url := "http://" + searcher.MeshServer + MeshShadowServicesURL
	result, err := jsonClient.
		GetByContext(ctx, url, nil, nil).
		HandleResponse(func(b []byte, statusCode int) (interface{}, error) {
			if statusCode == http.StatusNotFound {
				return nil, errors.Wrap(NotFoundError, "list service")
			}

			if statusCode >= 300 || statusCode < 200 {
				return nil, errors.Errorf("call GET %s failed, return statuscode %d text %s", url, statusCode, string(b))
			}

			services := []object.ShadowService{}
			err := json.Unmarshal(b, &services)
			if err != nil {
				return nil, errors.Wrapf(err, "unmarshal shadow services result")
			}
			return services, nil
		})
	if err != nil {
		return nil, err
	}
	return result.([]object.ShadowService), err
}
