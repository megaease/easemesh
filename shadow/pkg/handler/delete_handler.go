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

	"github.com/megaease/easemesh/mesh-operator/pkg/api/v1beta1"
	"github.com/megaease/easemesh/mesh-shadow/pkg/object"
	"github.com/megaease/easemesh/mesh-shadow/pkg/utils"
	appv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	runTimeClient "sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	Separator = '/'
)

type Deleter interface {
	Delete(obj interface{})
	FindDeletableObjs(obj interface{})
}

type ShadowServiceDeleter struct {
	KubeClient    kubernetes.Interface
	RunTimeClient *runTimeClient.Client
	CRDClient     rest.Interface
	DeleteChan    chan interface{}
}

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

	case v1beta1.MeshDeployment:
		meshDeployment := obj.(v1beta1.MeshDeployment)
		err = utils.DeleteMeshDeployment(deleter.CRDClient, meshDeployment.Namespace, meshDeployment)
		if err != nil {
			log.Printf("Delete ShadowService's MeshDeployment failed. NameSpace: %s, Name: %s. error: %s", meshDeployment.Namespace, meshDeployment.Name, err)
		} else {
			log.Printf("Delete ShadowService's MeshDeployment Success. NameSpace: %s, Name: %s.", meshDeployment.Namespace, meshDeployment.Name)
		}
	}
}

func (deleter *ShadowServiceDeleter) FindDeletableObjs(obj interface{}) {
	shadowServiceList := obj.([]object.ShadowService)
	shadowServiceNameMap := make(map[string]int)
	for _, ss := range shadowServiceList {
		shadowServiceNameMap[namespacedName(ss.Namespace, ss.Name)] = 1
	}

	namespaces, err := utils.ListNameSpaces(deleter.KubeClient)
	if err != nil {
		log.Printf("List namespaces failed. error: %s", err)
	}

	listOptions := metav1.ListOptions{
		LabelSelector: shadowLabelKey + "=true",
	}
	for _, namespace := range namespaces {
		meshDeployments, err := utils.ListMeshDeployment(deleter.CRDClient, namespace.Name, listOptions)
		if err != nil {
			log.Printf("List MeshDeployment failed. error: %s", err)
		}
		if meshDeployments != nil {
			for _, meshDeployment := range meshDeployments.Items {
				if shadowServiceName, ok := meshDeployment.Annotations[shadowServiceNameAnnotationKey]; ok {
					if _, ok = shadowServiceNameMap[namespacedName(namespace.Name, shadowServiceName)]; ok {
						continue
					} else {
						deleter.DeleteChan <- meshDeployment
					}
				}
			}
		}

		deployments, err := utils.ListDeployments(namespace.Name, deleter.KubeClient, listOptions)
		if err != nil {
			log.Printf("List Deployment failed. error: %s", err)
		}
		for _, deployment := range deployments {
			if shadowServiceName, ok := deployment.Annotations[shadowServiceNameAnnotationKey]; ok {
				if _, ok = shadowServiceNameMap[namespacedName(namespace.Name, shadowServiceName)]; ok {
					continue
				} else {
					deleter.DeleteChan <- deployment
				}
			}
		}
	}

}

func namespacedName(namespace string, name string) string {
	return namespace + string(Separator) + name
}
