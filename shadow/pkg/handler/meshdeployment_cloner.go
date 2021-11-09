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
	"github.com/megaease/easemesh/mesh-operator/pkg/api/v1beta1"
	"github.com/megaease/easemesh/mesh-shadow/pkg/object"
	"github.com/megaease/easemesh/mesh-shadow/pkg/utils"
	"github.com/pkg/errors"
	appsV1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type CloneMeshDeploymentFunc func() error

func (cloner *ShadowServiceCloner) cloneMeshDeployment(sourceMeshDeployment *v1beta1.MeshDeployment, shadowService *object.ShadowService) CloneMeshDeploymentFunc {
	shadowMeshDeployment := cloner.decorateShadowMeshDeployment(sourceMeshDeployment, shadowService)
	return func() error {
		err := utils.DeployMesheployment(cloner.CRDClient, shadowMeshDeployment.Namespace, shadowMeshDeployment)
		if err != nil {
			return errors.Wrapf(err, "Clone mesh deployment %s for service %s failed", sourceMeshDeployment.Name, shadowService.ServiceName)
		}
		return err
	}
}

func (cloner *ShadowServiceCloner) generateShadowMeshDeployment(sourceMeshDeployment *v1beta1.MeshDeployment, shadowService *object.ShadowService) *v1beta1.MeshDeployment {
	if sourceMeshDeployment.Labels == nil {
		sourceMeshDeployment.Labels = map[string]string{}
	}
	injectShadowLabels(sourceMeshDeployment.Labels)

	if sourceMeshDeployment.Annotations == nil {
		sourceMeshDeployment.Annotations = map[string]string{}
	}
	injectShadowAnnotation(sourceMeshDeployment.Annotations, shadowService)
	return &v1beta1.MeshDeployment{
		TypeMeta: sourceMeshDeployment.TypeMeta,
		ObjectMeta: metav1.ObjectMeta{
			Name:        shadowName(sourceMeshDeployment.Name),
			Namespace:   sourceMeshDeployment.Namespace,
			Labels:      sourceMeshDeployment.Labels,
			Annotations: sourceMeshDeployment.Annotations,
		},
	}
}

func (cloner *ShadowServiceCloner) decorateShadowMeshDeployment(sourceMeshDeployment *v1beta1.MeshDeployment, shadowService *object.ShadowService) *v1beta1.MeshDeployment {
	shadowMeshDeployment := cloner.generateShadowMeshDeployment(sourceMeshDeployment, shadowService)
	shadowMeshDeployment.Spec.Service = sourceMeshDeployment.Spec.Service

	labels := shadowMeshDeployment.Labels
	if labels == nil {
		labels = make(map[string]string)
	}

	shadowServiceLabels := shadowServiceLabels()
	for k, v := range shadowServiceLabels {
		labels[k] = v
	}
	shadowMeshDeployment.Labels = labels

	deployment := &appsV1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      sourceMeshDeployment.Name,
			Namespace: sourceMeshDeployment.Namespace,
		},
		Spec: sourceMeshDeployment.Spec.Deploy.DeploymentSpec,
	}
	if sourceMeshDeployment.Spec.Service.AppContainerName != "" {
		deployment.Annotations[shadowAppContainerNameKey] = sourceMeshDeployment.Spec.Service.AppContainerName
	}
	shadowDeployment := cloner.cloneDeploymentSpec(deployment, shadowService)
	shadowMeshDeployment.Spec.Deploy.DeploymentSpec = shadowDeployment.Spec

	return shadowMeshDeployment
}