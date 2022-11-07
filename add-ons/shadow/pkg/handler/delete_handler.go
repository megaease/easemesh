/*
 * Copyright (c) 2021, MegaEase
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
	"strings"

	"github.com/megaease/easemesh/mesh-shadow/pkg/object"
	"github.com/megaease/easemesh/mesh-shadow/pkg/utils"
	installbase "github.com/megaease/easemeshctl/cmd/client/command/meshinstall/base"

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
	block := obj.(ShadowServiceBlock)
	if block.deployment == nil {
		log.Printf("shadow service %s: deployment is nil, skip delete",
			block.shadowService.ServiceName)
		return
	}

	deployment := block.deployment

	err := utils.DeleteDeployment(deployment.Namespace, deployment.Name, deleter.KubeClient, metav1.DeleteOptions{})
	if err != nil {
		log.Printf("delete deployment %s/%s failed: %s", deployment.Namespace, deployment.Name, err)
	} else {
		log.Printf("delete deployment %s/%s succeed", deployment.Namespace, deployment.Name)
	}

	shadowConfigMapIDs := deployment.Annotations[shadowConfigMapsAnnotationKey]
	for _, shadowConfigMapID := range strings.Split(shadowConfigMapIDs, ",") {
		id := strings.Split(shadowConfigMapID, "/")
		if len(id) != 2 {
			log.Printf("invalid shadow configmap id: %s", id)
			continue
		}

		ns, name := id[0], id[1]

		configMapResource := [][]string{{"configmaps", name}}
		installbase.DeleteResources(deleter.KubeClient, configMapResource,
			ns, installbase.DeleteCoreV1Resource)
	}

	shadowSecretIDs := deployment.Annotations[shadowSecretsAnnotationKey]
	for _, shadowSecretID := range strings.Split(shadowSecretIDs, ",") {
		id := strings.Split(shadowSecretID, "/")
		if len(id) != 2 {
			log.Printf("invalid shadow secret id: %s", id)
			continue
		}

		ns, name := id[0], id[1]

		secretResource := [][]string{{"secrets", name}}
		installbase.DeleteResources(deleter.KubeClient, secretResource,
			ns, installbase.DeleteCoreV1Resource)
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
		deployment := deployment.DeepCopy()

		if shadowServiceName, ok := deployment.Annotations[shadowServiceNameAnnotationKey]; ok {
			shadowService, _ := shadowServiceNameMap[namespacedName(namespace, shadowServiceName)]
			if !shadowServiceExists(namespacedName(namespace, shadowServiceName), shadowServiceNameMap) {
				deleter.DeleteChan <- ShadowServiceBlock{
					shadowService: object.ShadowService{
						Name: shadowServiceName,
					},
					deployment: deployment,
				}
				continue
			}

			deploymentSourceName := sourceName(deployment.Name, &shadowService)
			sourceExisted := sourceDeploymentExistsFn(deploymentSourceName, shadowService.ServiceName)

			if !sourceExisted {
				deleter.DeleteChan <- ShadowServiceBlock{
					shadowService: shadowService,
					deployment:    deployment,
				}
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
