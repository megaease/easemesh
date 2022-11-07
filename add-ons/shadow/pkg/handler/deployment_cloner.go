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
	"encoding/json"
	"fmt"
	"log"
	"reflect"
	"sort"
	"strings"

	"github.com/megaease/easemesh/mesh-shadow/pkg/object"
	installbase "github.com/megaease/easemeshctl/cmd/client/command/meshinstall/base"
	"github.com/pkg/errors"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ShadowDeploymentFunc type ShadowFunc func(ctx *object.CloneContext) error
type ShadowDeploymentFunc func() error

func (cloner *ShadowServiceCloner) cloneDeployment(sourceDeployment *appsv1.Deployment, shadowService *object.ShadowService) ShadowDeploymentFunc {
	shadowDeployment := cloner.cloneDeploymentSpec(sourceDeployment, shadowService)
	shadowConfigmaps := cloner.cloneConfigMapSpecs(shadowDeployment, shadowService)
	shadowSecrets := cloner.cloneSecretSpecs(shadowDeployment, shadowService)
	return func() error {
		for _, cm := range shadowConfigmaps {
			// Clean retrieval data.
			cm.ResourceVersion, cm.UID = "", ""

			err := installbase.DeployConfigMap(&cm, cloner.KubeClient, cm.Namespace)
			if err != nil {
				return errors.Wrapf(err, "deploy shadow configmap %s for service %s failed", sourceDeployment.Name, shadowService.ServiceName)
			}

			log.Printf("deploy configmap %s for service %s succeed", sourceDeployment.Name, shadowService.ServiceName)
		}

		for _, secret := range shadowSecrets {
			// Clean retrieval data.
			secret.ResourceVersion, secret.UID = "", ""

			err := installbase.DeploySecret(&secret, cloner.KubeClient, secret.Namespace)
			if err != nil {
				return errors.Wrapf(err, "deploy shadow secret %s for service %s failed", sourceDeployment.Name, shadowService.ServiceName)
			}

			log.Printf("deploy secret %s for service %s succeed", sourceDeployment.Name, shadowService.ServiceName)
		}

		err := installbase.DeployDeployment(shadowDeployment, cloner.KubeClient, shadowDeployment.Namespace)
		if err != nil {
			return errors.Wrapf(err, "deploy shadow deployment %s for service %s failed", sourceDeployment.Name, shadowService.ServiceName)
		}

		log.Printf("deploy deployment %s for service %s succeed", sourceDeployment.Name, shadowService.ServiceName)

		return nil
	}
}

func (cloner *ShadowServiceCloner) cloneDeploymentSpec(sourceDeployment *appsv1.Deployment, shadowService *object.ShadowService) *appsv1.Deployment {
	shadowDeployment := cloner.generateShadowDeployment(sourceDeployment, shadowService)
	cloner.decorateShadowDeploymentBaseSpec(shadowDeployment, sourceDeployment)
	cloner.decorateEnvs(shadowDeployment, sourceDeployment, shadowService)
	cloner.decorateVolumes(shadowDeployment, shadowService)

	return shadowDeployment
}

func (cloner *ShadowServiceCloner) cloneConfigMapSpecs(deployment *appsv1.Deployment, shadowService *object.ShadowService) []corev1.ConfigMap {
	configMaps := []corev1.ConfigMap{}
	for _, configMap := range shadowService.ConfigMaps {
		configMap.Name = shadowVolumeName(configMap.Name, deployment.Name)
		configMaps = append(configMaps, configMap)
	}

	return configMaps
}

func (cloner *ShadowServiceCloner) cloneSecretSpecs(deployment *appsv1.Deployment, shadowService *object.ShadowService) []corev1.Secret {
	secrets := []corev1.Secret{}
	for _, secret := range shadowService.Secrets {
		secret.Name = shadowVolumeName(secret.Name, deployment.Name)
		secrets = append(secrets, secret)
	}

	return secrets
}

// findContainer returns the copy of the container,
// which means it won't change the original container when changing the result.
func findContainer(containers []corev1.Container, containerName string) (*corev1.Container, bool) {
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
	if config == nil {
		return nil
	}

	env := &corev1.EnvVar{
		Name: envName,
	}
	switch c := config.(type) {
	case string:
		env.Value = c
	default:
		value := reflect.ValueOf(c)
		if value.Kind() == reflect.Ptr && value.IsNil() {
			return nil
		}

		configJSON, err := json.Marshal(config)
		if err != nil {
			return nil
		}
		env.Value = string(configJSON)
	}

	return env
}

func (cloner *ShadowServiceCloner) generateShadowDeployment(sourceDeployment *appsv1.Deployment, shadowService *object.ShadowService) *appsv1.Deployment {
	deployment := &appsv1.Deployment{
		TypeMeta: sourceDeployment.TypeMeta,
		ObjectMeta: metav1.ObjectMeta{
			Name:        shadowDeploymentName(sourceDeployment.Name, shadowService),
			Namespace:   sourceDeployment.Namespace,
			Labels:      sourceDeployment.Labels,
			Annotations: sourceDeployment.Annotations,
		},
	}

	if deployment.Labels == nil {
		deployment.Labels = map[string]string{}
	}
	injectShadowLabels(deployment.Labels)

	if deployment.Annotations == nil {
		deployment.Annotations = map[string]string{}
	}
	injectShadowAnnotation(deployment, shadowService)

	return deployment
}

func (cloner *ShadowServiceCloner) decorateShadowDeploymentBaseSpec(deployment *appsv1.Deployment, sourceDeployment *appsv1.Deployment) *appsv1.Deployment {
	deployment.Spec = sourceDeployment.Spec

	matchLabels := deployment.Spec.Selector.MatchLabels
	if matchLabels == nil {
		matchLabels = map[string]string{}
	}

	injectShadowLabels(matchLabels)
	deployment.Spec.Selector = &metav1.LabelSelector{
		MatchLabels: matchLabels,
	}

	sourceLabels := deployment.Spec.Template.Labels
	if sourceLabels == nil {
		sourceLabels = make(map[string]string)
	}
	for k, v := range matchLabels {
		sourceLabels[k] = v
	}
	deployment.Spec.Template.Labels = sourceLabels

	containers := deployment.Spec.Template.Spec.Containers
	deployment.Spec.Template.Spec.Containers = shadowContainers(containers)

	initContainers := deployment.Spec.Template.Spec.InitContainers
	deployment.Spec.Template.Spec.InitContainers = shadowInitContainers(initContainers)

	volumes := deployment.Spec.Template.Spec.Volumes
	deployment.Spec.Template.Spec.Volumes = shadowVolumes(volumes)
	return deployment
}

type envList []corev1.EnvVar

func (e envList) Len() int           { return len(e) }
func (e envList) Less(i, j int) bool { return e[i].Name < e[j].Name }
func (e envList) Swap(i, j int)      { e[i], e[j] = e[j], e[i] }

func (cloner *ShadowServiceCloner) decorateEnvs(deployment *appsv1.Deployment, sourceDeployment *appsv1.Deployment, shadowService *object.ShadowService) *appsv1.Deployment {
	envs := shadowEnvs(shadowService)

	appContainerName, _ := sourceDeployment.Annotations[shadowAppContainerNameKey]
	appContainer, _ := findContainer(deployment.Spec.Template.Spec.Containers, appContainerName)
	appContainer.Env = injectEnvVars(appContainer.Env, envs...)
	deployment.Spec.Template.Spec.Containers = injectContainers(deployment.Spec.Template.Spec.Containers, *appContainer)
	return deployment
}

func shadowVolumeName(volumeName, shadowDeploymentName string) string {
	// shadowDeploymentName contains suffix -shadow already.
	// There is no need to repeat it in the shadow volume name.
	return fmt.Sprintf("%s-%s", volumeName, shadowDeploymentName)
}

func (cloner *ShadowServiceCloner) decorateVolumes(deployment *appsv1.Deployment, shadowService *object.ShadowService) {
	for i, volume := range deployment.Spec.Template.Spec.Volumes {
		shadowName := shadowVolumeName(volume.Name, deployment.Name)

		if volume.ConfigMap != nil {
			deployment.Spec.Template.Spec.Volumes[i].ConfigMap.Name = shadowName
		}

		if volume.Secret != nil {
			deployment.Spec.Template.Spec.Volumes[i].Secret.SecretName = shadowName
		}
	}
}

func shadowEnvs(shadowService *object.ShadowService) []corev1.EnvVar {
	envs := make(map[string]interface{})
	envs[databaseShadowConfigEnv] = shadowService.MySQL
	envs[elasticsearchShadowConfigEnv] = shadowService.ElasticSearch
	envs[redisShadowConfigEnv] = shadowService.Redis
	envs[kafkaShadowConfigEnv] = shadowService.Kafka
	envs[rabbitmqShadowConfigEnv] = shadowService.RabbitMQ

	// User defined envs.
	for k, v := range shadowService.Envs {
		envs[k] = v
	}

	newEnvs := make([]corev1.EnvVar, 0)
	for k, v := range envs {
		env := generateShadowConfigEnv(k, v)
		if env != nil {
			newEnvs = append(newEnvs, *env)
		}
	}

	sort.Sort(envList(newEnvs))

	return newEnvs
}

func shadowDeploymentName(name string, shadowService *object.ShadowService) string {
	return fmt.Sprintf("%s-%s", name, shadowService.CanaryName())
}

func sourceName(name string, ss *object.ShadowService) string {
	return strings.TrimSuffix(name, fmt.Sprintf("-%s", ss.CanaryName()))
}

func injectShadowLabels(labels map[string]string) {
	labels[shadowLabelKey] = "true"
}

func injectShadowAnnotation(deployment *appsv1.Deployment,
	service *object.ShadowService,
) {
	deployment.Annotations[shadowServiceNameAnnotationKey] = service.Name
	deployment.Annotations[shadowServiceVersionLabelAnnotationKey] = fmt.Sprintf(
		shadowServiceVersionLabelAnnotationValueFormat,
		service.CanaryName())

	shadowConfigMapIDs := []string{}
	for _, configMap := range service.ConfigMaps {
		shadowConfigMapID := fmt.Sprintf("%s/%s", configMap.Namespace, shadowVolumeName(configMap.Name, deployment.Name))
		shadowConfigMapIDs = append(shadowConfigMapIDs, shadowConfigMapID)
	}
	deployment.Annotations[shadowConfigMapsAnnotationKey] = strings.Join(shadowConfigMapIDs, ",")

	shadowSecretIDs := []string{}
	for _, secret := range service.Secrets {
		shadowSecretID := fmt.Sprintf("%s/%s", secret.Namespace, shadowVolumeName(secret.Name, deployment.Name))
		shadowSecretIDs = append(shadowSecretIDs, shadowSecretID)
	}
	deployment.Annotations[shadowSecretsAnnotationKey] = strings.Join(shadowSecretIDs, ",")
}

func shadowContainers(containers []corev1.Container) []corev1.Container {
	newContainers := make([]corev1.Container, 0)
	for _, container := range containers {
		// Prune sidecar container in case of repeated insertion from operator.
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
		// Prune volumes generated from mesh operator.
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
		// Prune volumes generated from mesh operator.
		if initContainer.Name != initContainerName {
			newInitContainers = append(newInitContainers, initContainer)
		}
	}
	return newInitContainers
}

func shadowVolumes(volumes []corev1.Volume) []corev1.Volume {
	newVolumes := make([]corev1.Volume, 0)

	for _, volume := range volumes {
		// Prune volumes generated from mesh operator.
		if volume.Name == agentVolumeName || volume.Name == sidecarVolumeName {
			continue
		}
		newVolumes = append(newVolumes, volume)
	}
	return newVolumes
}
