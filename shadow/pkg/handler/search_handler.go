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

func (searcher *ShadowServiceDeploySearcher) Search(objs interface{}) {
	shadowServices := objs.([]object.ShadowService)
	if len(shadowServices) == 0 {
		return
	}

	shadowServicesNamespacesMap := make(map[string][]object.ShadowService)
	for _, shadowService := range shadowServices {
		if _, ok := shadowServicesNamespacesMap[shadowService.Namespace]; ok {
			shadowServicesNamespacesMap[shadowService.Namespace] = append(shadowServicesNamespacesMap[shadowService.Namespace], shadowService)
		} else {
			shadowServicesNamespacesMap[shadowService.Namespace] = []object.ShadowService{shadowService}
		}
	}

	for namespace, shadowServiceList := range shadowServicesNamespacesMap {

		shadowServiceNameMap := make(map[string]object.ShadowService)
		for _, ss := range shadowServiceList {
			shadowServiceNameMap[ss.ServiceName] = ss
		}

		meshDeploymentList, err := utils.ListMeshDeployment(searcher.CRDClient, namespace, metav1.ListOptions{})
		if err != nil {
			log.Printf("Query MeshDeployment for shadow service error. %s", err)
		}
		for _, meshDeployment := range meshDeploymentList.Items {
			if isShadowDeployment(meshDeployment.Spec.Deploy.DeploymentSpec) {
				continue
			}
			if ss, ok := shadowServiceNameMap[meshDeployment.Spec.Service.Name]; ok {
				searcher.ResultChan <- ServiceCloneBlock{
					ss,
					meshDeployment,
				}
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
			if serviceName, ok := annotations[MeshServiceAnnotation]; ok {
				if ss, ok := shadowServiceNameMap[serviceName]; ok {
					searcher.ResultChan <- ServiceCloneBlock{
						ss,
						deployment,
					}
				}

			}
		}
	}
}

func isShadowDeployment(spec appsV1.DeploymentSpec) bool {
	if shadowLabel, ok := spec.Selector.MatchLabels[shadowLabelKey]; ok {
		return shadowLabel == "true"
	}
	return false
}
