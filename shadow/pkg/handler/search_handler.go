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

package handler

import (
	"log"

	"github.com/megaease/easemesh/mesh-shadow/pkg/object"
	"github.com/megaease/easemesh/mesh-shadow/pkg/utils"
	appsV1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	runTimeClient "sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	MeshServiceAnnotation = "mesh.megaease.com/service-name"
)

type Searcher interface {
	Search(obj interface{})
}

type ShadowServiceDeploySearcher struct {
	KubeClient    kubernetes.Interface
	RunTimeClient *runTimeClient.Client
	CRDClient     *rest.RESTClient
	ResultChan    chan interface{}
}

type ServiceCloneBlock struct {
	service   object.ShadowService
	deployObj interface{}
}

func (searcher *ShadowServiceDeploySearcher) Search(obj interface{}) {
	shadowService := obj.(object.ShadowService)
	namespace := shadowService.Namespace
	serviceName := shadowService.ServiceName

	meshDeploymentList, err := utils.ListMeshDeployment(namespace, searcher.CRDClient, metav1.ListOptions{})
	if err != nil {
		log.Printf("Query MeshDeployment for shadow service error. %s", err)
	}
	for _, meshDeployment := range meshDeploymentList.Items {
		if isShadowDeployment(meshDeployment.Spec.Deploy.DeploymentSpec) {
			continue
		}
		if meshDeployment.Spec.Service.Name == serviceName {
			searcher.ResultChan <- ServiceCloneBlock{
				shadowService,
				meshDeployment,
			}
			return
		}
	}

	deployments, err := utils.ListDeployments(namespace, searcher.KubeClient, metav1.ListOptions{})
	if err != nil {
		log.Printf("Query Deployment for shadow service error. %s", err)
	}
	for _, deployment := range deployments {
		if isShadowDeployment(deployment.Spec) {
			continue
		}
		annotations := deployment.Annotations
		if _, ok := annotations[MeshServiceAnnotation]; ok {
			if serviceName == annotations[MeshServiceAnnotation] {
				searcher.ResultChan <- ServiceCloneBlock{
					shadowService,
					deployment,
				}
				return
			}
		}
	}
	log.Printf("The service doesn't have MeshDeployment or Deployment for run it. Service: %s, NameSpace: %s, "+
		"ShadowService: %s", serviceName, namespace, shadowService.Name)
}

func isShadowDeployment(spec appsV1.DeploymentSpec) bool {
	if shadowLabel, ok := spec.Selector.MatchLabels[shadowLabelKey]; ok {
		return shadowLabel == "true"
	}
	return false
}
