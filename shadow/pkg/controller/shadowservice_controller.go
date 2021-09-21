package controller

import (
	"time"

	"github.com/megaease/easemesh/mesh-shadow/pkg/common"
	"github.com/megaease/easemesh/mesh-shadow/pkg/handler"
	"github.com/megaease/easemesh/mesh-shadow/pkg/object"
	"github.com/megaease/easemesh/mesh-shadow/pkg/syncer"
	"github.com/megaease/easemesh/mesh-shadow/pkg/utils"
	"github.com/pkg/errors"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type (
	// ShadowExecutorService is a service which orchestrator cloner and deployer
	// to accomplish generate shadow service
	ShadowExecutorService interface {
		Do()
	}

	ShadowServiceController struct {
		KubeClient    *kubernetes.Clientset
		RunTimeClient *client.Client
		CRDClient     *rest.RESTClient
		syncer        *syncer.Syncer
		searchHanler  *handler.SearchHandler
		cloneHandler  *handler.CloneHandler

		cloneChan chan interface{}
	}

	// ServiceConfig holds configuration of shadow service controller
	ServiceConfig struct {
		MeshServer string
		Interval   time.Duration
	}
	// Opt is option to control service configuration
	Opt func(sc *ServiceConfig) error
)

// NewController a CollectorService to collect K8s metrics
func New(opts ...Opt) (*ShadowServiceController, error) {
	config := ServiceConfig{}
	for _, opt := range opts {
		if err := opt(&config); err != nil {
			return nil, err
		}
	}

	kubernetesClient, err := utils.NewKubernetesClient()
	if err != nil {
		return nil, errors.Wrapf(err, "new kubernetes clientSet error")
	}
	runtimeClient, err := utils.NewRuntimeClient()
	if err != nil {
		return nil, errors.Wrapf(err, "new Controller Runtime client error")
	}
	crdRestClient, err := utils.NewCRDRestClient()
	if err != nil {
		return nil, errors.Wrapf(err, "new Resst client error")
	}

	cloneChan := make(chan interface{})

	searchHandler := &handler.SearchHandler{
		KubeClient:    kubernetesClient,
		RunTimeClient: &runtimeClient,
		CRDClient:     crdRestClient,
		MeshServer:    config.MeshServer,
		CloneChan:     cloneChan,
		Interval:      30 * time.Second,
		Registry:      common.NewCallbackRegistry(),
	}

	cloneHandler := &handler.CloneHandler{
		KubeClient:    kubernetesClient,
		RunTimeClient: &runtimeClient,
		CRDClient:     crdRestClient,
	}

	server := syncer.Server{
		10 * time.Second,
		config.MeshServer,
	}
	syncer, err := server.NewSyncer(1 * time.Minute)

	return &ShadowServiceController{kubernetesClient, &runtimeClient, crdRestClient,
		syncer, searchHandler, cloneHandler, cloneChan}, nil
}

// Do start shadow service query and clone data
func (s *ShadowServiceController) Do() <-chan struct{} {
	result := make(chan struct{})
	customObjectsChan, _ := s.syncer.Sync("ShadowService")
	go func() {
		for {
			select {
			case obj := <-customObjectsChan:
				for _, v := range obj {
					shadowService := v.(object.ShadowService)
					s.searchHanler.SearchOriginDeployments(shadowService)
				}
			}
		}
	}()

	go func() {
		for {
			select {
			case obj := <-s.cloneChan:
				s.cloneHandler.Clone(obj)
			}
		}
	}()
	return result
}
