/*
 * Copyright (c) 2017, MegaEase
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

package controller

import (
	"sync"
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
		Do(stopChan <-chan struct{})
	}

	ShadowServiceController struct {
		kubeClient    kubernetes.Interface
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

// Do start shadow service sync and clone.
func (s *ShadowServiceController) Do(wg *sync.WaitGroup, stopChan <-chan struct{}) {
	shadowServiceChan, _ := s.syncer.Sync(ShadowServiceKind)

	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			select {
			case <-stopChan:
				s.syncer.Close()
				return
			case obj := <-shadowServiceChan:
				s.searcher.Search(obj)
			}
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			select {
			case <-stopChan:
				return
			case obj := <-s.cloneChan:
				s.cloner.Clone(obj)
			}
		}
	}()
}
