package controller

import (
	"time"

	"github.com/megaease/easemesh/mesh-shadow/pkg/handler"
	"github.com/megaease/easemesh/mesh-shadow/pkg/syncer"
	"github.com/megaease/easemesh/mesh-shadow/pkg/utils"
	"github.com/pkg/errors"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	ShadowServiceKind = "ShadowService"
)

type (
	// ShadowServiceExecutor is executor which orchestrator cloner and deployer for run shadow service.
	ShadowServiceExecutor interface {
		Do()
	}

	ShadowServiceController struct {
		kubeClient    *kubernetes.Clientset
		runTimeClient *client.Client
		crdClient     *rest.RESTClient

		syncer   *syncer.ShadowServiceSyncer
		cloner   handler.Cloner
		searcher handler.Searcher

		cloneChan chan interface{}
	}

	// Config holds configuration of ShadowServiceController.
	Config struct {
		MeshServer     string
		PullInterval   time.Duration
		RequestTimeout time.Duration
	}
	// Opt is option to control EaseMesh control plane.
	Opt func(sc *Config) error
)

// NewShadowServiceController create ShadowServiceController for execute ShadowService processing.
func NewShadowServiceController(opts ...Opt) (*ShadowServiceController, error) {
	config := Config{}
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

	shadowServiceCloner := handler.ShadowServiceCloner{
		KubeClient:    kubernetesClient,
		RunTimeClient: &runtimeClient,
		CRDClient:     crdRestClient,
	}

	shadowServiceSearcher := handler.ShadowServiceDeploySearcher{
		KubeClient:    kubernetesClient,
		RunTimeClient: &runtimeClient,
		CRDClient:     crdRestClient,
		ResultChan:    cloneChan,
	}

	shadowServiceSyncer, err := syncer.NewSyncer(config.MeshServer, config.RequestTimeout, config.PullInterval)

	return &ShadowServiceController{kubernetesClient, &runtimeClient, crdRestClient, shadowServiceSyncer, &shadowServiceCloner, &shadowServiceSearcher, cloneChan}, nil
}

func Init() {
	// shadowServiceKind := object.CustomObjectKind{
	// 	Name: ShadowServiceKind,
	// 	JsonSchema: "{" +
	// 		"\"name\": \"string\",  " +
	// 		"\"namespace\": \"string\", " +
	// 		"\"serviceName\": \"string\", " +
	// 		"\"mysql\": {\"uris\": \"[]string\",\"userName\": \"string\", \"password\": \"string\"}, " +
	// 		"\"kafka\": {\"uris\": \"[]string\"}, " +
	// 		"\"redis\": {\"uris\": \"[]string\",\"userName\": \"string\", \"password\": \"string\"}, " +
	// 		"\"rabbitMq\": {\"uris\": \"[]string\",\"userName\": \"string\", \"password\": \"string\"}, " +
	// 		"\"elasticSearch\": {\"uris\": \"[]string\",\"userName\": \"string\", \"password\": \"string\"}" +
	// 		"}",
	// }
}

// Do start shadow service sync and clone.
func (s *ShadowServiceController) Do() <-chan struct{} {
	result := make(chan struct{})
	customObjectsChan, _ := s.syncer.Sync(ShadowServiceKind)
	go func() {
		for {
			select {
			case obj := <-customObjectsChan:
				s.searcher.Search(obj)
			}
		}
	}()

	go func() {
		for {
			select {
			case obj := <-s.cloneChan:
				s.cloner.Clone(obj)
			}
		}
	}()
	return result
}
