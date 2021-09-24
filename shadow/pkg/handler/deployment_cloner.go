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
	"encoding/json"
	"reflect"

	"github.com/megaease/easemesh/mesh-shadow/pkg/object"
	"github.com/megaease/easemesh/mesh-shadow/pkg/utils"
	"github.com/pkg/errors"
	appsV1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ShadowDeploymentFunc type ShadowFunc func(ctx *object.CloneContext) error
type ShadowDeploymentFunc func() error

type cloneDeploymentSpecFunc func(sourceDeployment *appsV1.Deployment, shadowService *object.ShadowService) *appsV1.Deployment

func (cloner *ShadowServiceCloner) cloneDeployment(sourceDeployment *appsV1.Deployment, shadowService *object.ShadowService) ShadowDeploymentFunc {
	shadowDeployment := cloner.cloneDeploymentSpec(sourceDeployment, shadowService)
	return func() error {
		err := utils.DeployDeployment(shadowDeployment, cloner.KubeClient, shadowDeployment.Namespace)
		if err != nil {
			return errors.Wrapf(err, "Clone deployment %s for service %s failed", sourceDeployment.Name, shadowService.ServiceName)
		}
		return err
	}
}

func (cloner *ShadowServiceCloner) cloneDeploymentSpec(sourceDeployment *appsV1.Deployment, shadowService *object.ShadowService) *appsV1.Deployment {
	shadowDeployment := cloner.injectShadowConfiguration(
		cloner.shadowDeploymentBaseSpec(
			cloner.shadowDeploymentInitialize(nil)))(sourceDeployment, shadowService)
	return shadowDeployment
}

func (cloner *ShadowServiceCloner) injectShadowConfiguration(fn cloneDeploymentSpecFunc) cloneDeploymentSpecFunc {
	return func(sourceDeployment *appsV1.Deployment, shadowService *object.ShadowService) *appsV1.Deployment {
		deployment := fn(sourceDeployment, shadowService)

		shadowConfigs := make(map[string]interface{})
		shadowConfigs[databaseShadowConfigEnv] = shadowService.MySQL
		shadowConfigs[elasticsearchShadowConfigEnv] = shadowService.ElasticSearch
		shadowConfigs[redisShadowConfigEnv] = shadowService.Redis
		shadowConfigs[kafkaShadowConfigEnv] = shadowService.Kafka
		shadowConfigs[rabbitmqShadowConfigEnv] = shadowService.RabbitMQ

		newEnvs := make([]corev1.EnvVar, 0)
		for k, v := range shadowConfigs {
			env := generateShadowConfigEnv(k, v)
			if env != nil {
				newEnvs = append(newEnvs, *env)
			}
		}

		appContainerName, _ := sourceDeployment.Annotations[shadowAppContainerNameKey]
		appContainer, _ := findAppContainer(sourceDeployment.Spec.Template.Spec.Containers, appContainerName)
		appContainer.Env = injectEnvVars(appContainer.Env, newEnvs...)
		deployment.Spec.Template.Spec.Containers = injectContainers(deployment.Spec.Template.Spec.Containers, *appContainer)
		return deployment
	}
}

// findContainer returns the copy of the container,
// which means it won't change the original container when changing the result.
func findAppContainer(containers []corev1.Container, containerName string) (*corev1.Container, bool) {

	if containerName == "" {
		for i, c := range containers {
			if c.Name == sidecarContainerName {
				continue
			}
			return &containers[i], false
		}

	} else {
		for _, container := range containers {
			if container.Name == containerName {
				return &container, true
			}
		}
	}
	return nil, false
}

func injectContainers(containers []corev1.Container, elems ...corev1.Container) []corev1.Container {
	for _, elem := range elems {
		replaced := false
		for i, existedContainer := range containers {
			if existedContainer.Name == elem.Name {
				containers[i] = elem
				replaced = true
			}
		}
		if !replaced {
			containers = append(containers, elem)
		}
	}

	return containers
}

func injectEnvVars(envVars []corev1.EnvVar, elems ...corev1.EnvVar) []corev1.EnvVar {
	for _, elem := range elems {
		replaced := false
		for i, existedEnvVar := range envVars {
			if existedEnvVar.Name == elem.Name {
				envVars[i] = elem
				replaced = true
			}
		}
		if !replaced {
			envVars = append(envVars, elem)
		}
	}
	return envVars
}

func generateShadowConfigEnv(envName string, config interface{}) *corev1.EnvVar {
	if config == nil || reflect.ValueOf(config).IsNil() {
		return nil
	}

	configJson, err := json.Marshal(config)
	if err != nil {
		return nil
	}

	env := &corev1.EnvVar{}
	env.Name = envName
	env.Value = string(configJson)
	return env

}

func (cloner *ShadowServiceCloner) shadowDeploymentInitialize(fn cloneDeploymentSpecFunc) cloneDeploymentSpecFunc {
	return func(sourceDeployment *appsV1.Deployment, shadowService *object.ShadowService) *appsV1.Deployment {
		return &appsV1.Deployment{
			TypeMeta: sourceDeployment.TypeMeta,
			ObjectMeta: metav1.ObjectMeta{
				Name:        shadowName(sourceDeployment.Name),
				Namespace:   sourceDeployment.Namespace,
				Labels:      sourceDeployment.Labels,
				Annotations: sourceDeployment.Annotations,
			},
		}

	}
}

func (cloner *ShadowServiceCloner) shadowDeploymentBaseSpec(fn cloneDeploymentSpecFunc) cloneDeploymentSpecFunc {
	return func(sourceDeployment *appsV1.Deployment, shadowService *object.ShadowService) *appsV1.Deployment {
		deployment := fn(sourceDeployment, shadowService)
		deployment.Spec = sourceDeployment.Spec

		labels := deployment.Spec.Selector.MatchLabels
		if labels == nil {
			labels = map[string]string{}
		}

		shadowServiceLabels := shadowServiceLabels()
		for k, v := range shadowServiceLabels {
			labels[k] = v
		}
		deployment.Spec.Selector = &metav1.LabelSelector{
			MatchLabels: labels,
		}
		deployment.Spec.Template.Labels = labels

		containers := deployment.Spec.Template.Spec.Containers
		deployment.Spec.Template.Spec.Containers = shadowContainers(containers)

		initContainers := deployment.Spec.Template.Spec.InitContainers
		deployment.Spec.Template.Spec.InitContainers = shadowInitContainers(initContainers)

		volumes := deployment.Spec.Template.Spec.Volumes
		deployment.Spec.Template.Spec.Volumes = shadowVolumes(volumes)
		return deployment
	}
}

func shadowName(name string) string {
	return name + shadowDeploymentNameSuffix
}
func shadowServiceLabels() map[string]string {
	selector := map[string]string{}
	selector[shadowLabelKey] = "true"
	return selector
}

func shadowContainers(containers []corev1.Container) []corev1.Container {
	newContainers := make([]corev1.Container, 0)
	for _, container := range containers {
		if container.Name != sidecarContainerName {
			newContainers = append(newContainers, shadowContainer(container))
		}
	}
	return newContainers
}

func shadowContainer(container corev1.Container) corev1.Container {
	mounts := container.VolumeMounts
	newMounts := make([]corev1.VolumeMount, 0)
	for _, mount := range mounts {
		if mount.Name != agentVolumeName {
			newMounts = append(newMounts, mount)
		}
	}
	container.VolumeMounts = newMounts
	return container
}

func shadowInitContainers(initContainers []corev1.Container) []corev1.Container {
	newInitContainers := make([]corev1.Container, 0)
	for _, initContainer := range initContainers {
		if initContainer.Name != initContainerName {
			newInitContainers = append(newInitContainers, initContainer)
		}
	}
	return newInitContainers
}

func shadowVolumes(volumes []corev1.Volume) []corev1.Volume {
	newVolumes := make([]corev1.Volume, 0)

	for _, volume := range volumes {
		if volume.Name == agentVolumeName || volume.Name == sidecarVolumeName {
			continue
		}
		newVolumes = append(newVolumes, volume)
	}
	return newVolumes
}
