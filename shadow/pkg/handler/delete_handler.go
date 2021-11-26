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
	appv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// Deleter is used to find and delete useless Objs.
type Deleter interface {
	Delete(obj interface{})
	FindDeletableObjs(obj interface{})
}

// ShadowServiceDeleter find and delete useless ShadowService's Deployment.
type ShadowServiceDeleter struct {
	KubeClient kubernetes.Interface
	DeleteChan chan interface{}
}

// Delete execute delete operation.
func (deleter *ShadowServiceDeleter) Delete(obj interface{}) {
	var err error
	switch obj.(type) {
	case appv1.Deployment:
		deployment := obj.(appv1.Deployment)
		err = utils.DeleteDeployment(deployment.Namespace, deployment.Name, deleter.KubeClient, metav1.DeleteOptions{})
		if err != nil {
			log.Printf("Delete ShadowService's Deployment failed. NameSpace: %s, Name: %s. error: %s", deployment.Namespace, deployment.Name, err)
		} else {
			log.Printf("Delete ShadowService's Deployment Success. NameSpace: %s, Name: %s.", deployment.Namespace, deployment.Name)
		}
	}
}

// FindDeletableObjs finds objects that can be deleted and send it for deletion.
func (deleter *ShadowServiceDeleter) FindDeletableObjs(obj interface{}) {
	shadowServiceList := obj.([]object.ShadowService)
	shadowServiceNameMap := make(map[string]object.ShadowService)
	for _, ss := range shadowServiceList {
		shadowServiceNameMap[namespacedName(ss.Namespace, ss.Name)] = ss
	}

	namespaces, err := utils.ListNameSpaces(deleter.KubeClient)
	if err != nil {
		log.Printf("List namespaces failed. error: %s", err)
	}

	for _, namespace := range namespaces {
		deleter.findDeletableDeployments(namespace.Name, shadowServiceNameMap)
	}

}

func shadowListOptions() metav1.ListOptions {
	return metav1.ListOptions{
		LabelSelector: shadowLabelKey + "=true",
	}
}

func (deleter *ShadowServiceDeleter) findDeletableDeployments(namespace string, shadowServiceNameMap map[string]object.ShadowService) {
	allDeployments, err := utils.ListDeployments(namespace, deleter.KubeClient, metav1.ListOptions{})
	allDeploymentsMap := make(map[string]appv1.Deployment)
	if err != nil {
		log.Printf("List Deployment failed. error: %s", err)
		return
	}

	for _, deployment := range allDeployments {
		allDeploymentsMap[deployment.Name] = deployment
	}

	sourceDeploymentExistsFn := func(name string, serviceName string) bool {
		deploy, ok := allDeploymentsMap[name]
		if !ok {
			return false
		}

		sourceServiceName, ok := deploy.Annotations[meshServiceAnnotation]
		if ok && sourceServiceName == serviceName {
			return true
		}
		return false
	}

	shadowDeployments, err := utils.ListDeployments(namespace, deleter.KubeClient, shadowListOptions())
	if err != nil {
		log.Printf("List Deployment failed. error: %s", err)
		return
	}

	for _, deployment := range shadowDeployments {
		if shadowServiceName, ok := deployment.Annotations[shadowServiceNameAnnotationKey]; ok {
			if !shadowServiceExists(namespacedName(namespace, shadowServiceName), shadowServiceNameMap) {
				deleter.DeleteChan <- deployment
				continue
			}
			shadowService, _ := shadowServiceNameMap[namespacedName(namespace, shadowServiceName)]
			if !sourceDeploymentExistsFn(sourceName(deployment.Name), shadowService.ServiceName) {
				deleter.DeleteChan <- deployment
				continue
			}
		}
	}
}

// If ShadowService is deleted, the shadow deployment need to be deleted.
func shadowServiceExists(namespacedName string, shadowServiceNameMap map[string]object.ShadowService) bool {
	_, ok := shadowServiceNameMap[namespacedName]
	return ok
}

func namespacedName(namespace string, name string) string {
	return namespace + string(separator) + name
}
