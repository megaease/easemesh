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
)

// Searcher is used to search existing objects for cloning.
type Searcher interface {
	Search(obj interface{})
}

// ShadowServiceDeploySearcher search objects from Kubernetes according to ShadowService.
type ShadowServiceDeploySearcher struct {
	KubeClient kubernetes.Interface
	ResultChan chan interface{}
}

// ServiceCloneBlock is passed by searcher to cloner to perform clone operation.
type ServiceCloneBlock struct {
	service   object.ShadowService
	deployObj interface{}
}

// Search finds objects from kubernetes based on ShadowService and sends them to cloner.
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
		searcher.searchDeployment(namespace, shadowServiceNameMap)
	}
}

func (searcher *ShadowServiceDeploySearcher) searchDeployment(namespace string, shadowServiceNameMap map[string]object.ShadowService) {
	deployments, err := utils.ListDeployments(namespace, searcher.KubeClient, metav1.ListOptions{})
	if err != nil {
		log.Printf("Query Deployment for shadow service error. %s", err)
	}
	for _, deployment := range deployments {
		if isShadowDeployment(deployment.Spec) {
			continue
		}
		annotations := deployment.Annotations
		if serviceName, ok := annotations[meshServiceAnnotation]; ok {
			if ss, ok := shadowServiceNameMap[serviceName]; ok {
				searcher.ResultChan <- ServiceCloneBlock{
					ss,
					deployment,
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
