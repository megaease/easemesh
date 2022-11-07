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

	"k8s.io/client-go/kubernetes"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const (

	// Init container stuff.
	initContainerName = "initializer"

	agentVolumeName   = "agent-volume"
	sidecarVolumeName = "sidecar-volume"

	// Sidecar container stuff.
	sidecarContainerName = "easemesh-sidecar"

	shadowServiceNameAnnotationKey = "mesh.megaease.com/shadow-service-name"

	// The value is comma-separated list of configmap names.
	// E.g: cm01-nginx-shadow,cm02-nginx-shadow
	shadowConfigMapsAnnotationKey = "mesh.megaease.com/shadow-configmaps"
	// The value is comma-separated list of secret names.
	// E.g: cm01-nginx-shadow,cm02-nginx-shadow
	shadowSecretsAnnotationKey = "mesh.megaease.com/shadow-secrets"

	shadowServiceVersionLabelAnnotationKey         = "mesh.megaease.com/service-labels"
	shadowServiceVersionLabelAnnotationValueFormat = "canary-name=%s"

	shadowServiceCanaryLabelKey        = "canary-name"
	shadowServiceCanaryHeader          = "X-Mesh-Shadow"
	shadowServiceCanaryDefaultPriority = 5

	shadowLabelKey            = "mesh.megaease.com/shadow-service"
	shadowAppContainerNameKey = "mesh.megaease.com/app-container-name"

	databaseShadowConfigEnv      = "EASE_RESOURCE_DATABASE"
	kafkaShadowConfigEnv         = "EASE_RESOURCE_KAFKA"
	rabbitmqShadowConfigEnv      = "EASE_RESOURCE_RABBITMQ"
	redisShadowConfigEnv         = "EASE_RESOURCE_REDIS"
	elasticsearchShadowConfigEnv = "EASE_RESOURCE_ELASTICSEARCH"

	meshServiceAnnotation = "mesh.megaease.com/service-name"

	separator = '/'
)

// Cloner is used to clone existed object.
type Cloner interface {
	Clone(obj interface{})
}

// ShadowServiceCloner clone Deployment according to ShadowService.
type ShadowServiceCloner struct {
	KubeClient    kubernetes.Interface
	RunTimeClient *client.Client
}

// Clone execute clone operation if there has ShadowService object.
func (cloner *ShadowServiceCloner) Clone(obj interface{}) {
	block := obj.(ShadowServiceBlock)
	if block.deployment == nil {
		log.Printf("shadow service %s: deployment is nil, skip clone", block.shadowService.ServiceName)
		return
	}

	err := cloner.cloneDeployment(block.deployment, &block.shadowService)()
	if err != nil {
		log.Printf("shodow service %s: clone failed: %s", block.shadowService.ServiceName, err)
	} else {
		log.Printf("shodow service %s: clone succeed", block.shadowService.ServiceName)
	}
}
