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

package deploymentmodifier

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/megaease/easemesh/mesh-operator/pkg/base"

	"github.com/pkg/errors"
	v1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
)

var (
	// Volumes stuff.
	volumes = []corev1.Volume{
		{
			Name: initContainerAgentVolumeName,
			VolumeSource: corev1.VolumeSource{
				EmptyDir: &corev1.EmptyDirVolumeSource{},
			},
		},
		{
			Name: initContainerSidecarVolumeName,
			VolumeSource: corev1.VolumeSource{
				EmptyDir: &corev1.EmptyDirVolumeSource{},
			},
		},
	}

	// Init container stuff.
	initContainerName      = "initializer"
	initContainerImageName = "megaease/easeagent-initializer:latest"

	initContainerAgentVolumeName        = "agent-volume"
	initContainerAgentVolumeMountPath   = "/agent-volume"
	initContainerSidecarVolumeName      = "sidecar-volume"
	initContainerSidecarVolumeMountPath = "/sidecar-volume"
	initContainerSidecarConfigPath      = "/sidecar-volume/sidecar-config.yaml"
	initContainerVolumeMounts           = []corev1.VolumeMount{
		{
			Name:      initContainerAgentVolumeName,
			MountPath: initContainerAgentVolumeMountPath,
		},
		{
			Name:      initContainerSidecarVolumeName,
			MountPath: initContainerSidecarVolumeMountPath,
		},
	}

	// Application container stuff.
	appContainerAgentVolumeName      = initContainerAgentVolumeName
	appContainerAgentVolumeMountPath = initContainerAgentVolumeMountPath
	appContainerVolumeMounts         = []corev1.VolumeMount{
		{
			Name:      appContainerAgentVolumeName,
			MountPath: appContainerAgentVolumeMountPath,
		},
	}

	appContainerJavaEnvName  = "JAVA_TOOL_OPTIONS"
	appContainerJavaEnvValue = fmt.Sprintf(" -javaagent:%s/easeagent.jar -Deaseagent.log.conf=%s/log4j2.xml ",
		appContainerAgentVolumeMountPath, appContainerAgentVolumeMountPath)
	appContainerEnvs = []corev1.EnvVar{
		{
			Name:  appContainerJavaEnvName,
			Value: appContainerJavaEnvValue,
		},
	}

	// Sidecar container stuff.
	sidecarContainerName      = "easemesh-sidecar"
	sidecarContainerImageName = func(baseRuntime *base.Runtime) string {
		if baseRuntime.SidecarImageName != "" {
			return baseRuntime.SidecarImageName
		}
		return "megaease/easegress:server-sidecar"
	}

	sidecarContainerVolumeName      = initContainerSidecarVolumeName
	sidecarContainerVolumeMountPath = initContainerSidecarVolumeMountPath
	sidecarContainerConfigPath      = initContainerSidecarConfigPath
	sidecarContainerVolumeMounts    = []corev1.VolumeMount{
		{
			Name:      sidecarContainerVolumeName,
			MountPath: sidecarContainerVolumeMountPath,
		},
	}

	sidecarContainerAppIPEnvName  = "APPLICATION_IP"
	sidecarContainerAppIPEnvValue = &corev1.EnvVarSource{
		FieldRef: &corev1.ObjectFieldSelector{
			FieldPath: "status.podIP",
		},
	}
	siecarContainerEnvs = []corev1.EnvVar{
		{
			Name:      sidecarContainerAppIPEnvName,
			ValueFrom: sidecarContainerAppIPEnvValue,
		},
	}

	sidecarContainerIngressPortName          = "sidecar-ingress"
	sidecarContainerIngressPortContainerPort = int32(13001)
	sidecarContainerEgressPortName           = "sidecar-egress"
	sidecarContainerEgressPortContainerPort  = int32(13002)
	sidecarContainerEurekaPortName           = "sidecar-eureka"
	sidecarContainerEurekaPortContainerPort  = int32(13009)
	sidecarContainerPorts                    = []corev1.ContainerPort{
		{
			Name:          sidecarContainerIngressPortName,
			ContainerPort: sidecarContainerIngressPortContainerPort,
		},
		{
			Name:          sidecarContainerEgressPortName,
			ContainerPort: sidecarContainerEgressPortContainerPort,
		},
		{
			Name:          sidecarContainerEurekaPortName,
			ContainerPort: sidecarContainerEurekaPortContainerPort,
		},
	}

	sidecarContainerCmd = []string{
		"/bin/sh",
		"-c",
		fmt.Sprintf("/opt/easegress/bin/easegress-server -f %s",
			initContainerSidecarConfigPath),
	}
)

func marshalLabels(labels map[string]string) string {
	labelsSlice := []string{}
	for k, v := range labels {
		labelsSlice = append(labelsSlice, k+"="+v)
	}
	return strings.Join(labelsSlice, "&")
}

func initContainerCommand(service *MeshService) []string {
	// TODO: Adjust for label names:
	// alive-probe -> mesh-alive-probe-url
	// application-port -> mesh-application-port
	// mesh-service-labels: Use `,` as separator instead of `&`.
	// mesh-servicename -> mesh-service-name

	const cmdTemplate = `set -e
cp -r /easeagent-volume/* %s

echo 'name: %s
cluster-join-urls: http://easemesh-controlplane-svc.easemesh:2380
cluster-request-timeout: 10s
cluster-role: reader
cluster-name: easemesh-control-plane
labels:
  alive-probe: %s
  application-port: %d
  mesh-service-labels: %s
  mesh-servicename: %s
' > %s`

	cmd := fmt.Sprintf(cmdTemplate,
		initContainerAgentVolumeMountPath,

		service.Name,

		service.AliveProbeURL,
		service.ApplicationPort,
		marshalLabels(service.Labels),
		service.Name,

		initContainerSidecarConfigPath)

	return []string{"sh", "-c", cmd}
}

type (
	DeploymentModifier struct {
		*base.Runtime
		meshService *MeshService
		deploy      *v1.Deployment
	}

	MeshService struct {
		// Name is required.
		Name string

		// Labels is optional.
		Labels map[string]string

		// AppContainerName is optional.
		// If empty, it will be the first container.
		AppContainerName string

		// ApplicationPort is optional.
		// If empty, we choose the first container port.
		ApplicationPort uint16

		// AliveProbeURL is optional.
		AliveProbeURL string
	}
)

// New creates a DeployModifier.
func New(baseRuntime *base.Runtime, meshService *MeshService,
	deploy *v1.Deployment) *DeploymentModifier {

	return &DeploymentModifier{
		Runtime:     baseRuntime,
		meshService: meshService,
		deploy:      deploy,
	}
}

// Modify modifies the Deployment.
func (m *DeploymentModifier) Modify() error {
	err := m.setupMeshService()
	if err != nil {
		return errors.Wrap(err, "set up mesh service")
	}

	m.injectVolumes(volumes...)
	m.injectInitContainer()
	m.injectSidecarContainer()

	err = m.adaptAppContainerSpec()
	if err != nil {
		return errors.Wrap(err, "complete app container spec")
	}

	return nil
}

func (m *DeploymentModifier) setupMeshService() error {
	if len(m.deploy.Spec.Template.Spec.Containers) == 0 {
		return fmt.Errorf("empty containers")
	}

	var container *corev1.Container
	if m.meshService.AppContainerName == "" {
		container = &m.deploy.Spec.Template.Spec.Containers[0]
		m.meshService.AppContainerName = container.Name
	} else {
		var exists bool
		container, exists = findContainer(m.deploy.Spec.Template.Spec.Containers, m.meshService.AppContainerName)
		if !exists {
			return errors.Errorf("container %s not found", m.meshService.AppContainerName)
		}
	}

	if m.meshService.ApplicationPort == 0 {
		if len(container.Ports) == 0 {
			return errors.Errorf("container %s got zero container port", container.Name)
		}
		m.meshService.ApplicationPort = uint16(container.Ports[0].ContainerPort)
	}

	return nil
}

func (m *DeploymentModifier) injectVolumes(volumes ...corev1.Volume) {
	for _, volume := range volumes {
		replaced := false
		for i, existedVolume := range m.deploy.Spec.Template.Spec.Volumes {
			if existedVolume.Name == volume.Name {
				m.deploy.Spec.Template.Spec.Volumes[i] = volume
				replaced = true
				break
			}
		}

		if !replaced {
			m.deploy.Spec.Template.Spec.Volumes = append(m.deploy.Spec.Template.Spec.Volumes, volume)
		}
	}

}

func (m *DeploymentModifier) injectInitContainer() {
	initContainer := corev1.Container{
		Name:            initContainerName,
		Image:           m.completeImageURL(initContainerImageName),
		ImagePullPolicy: corev1.PullPolicy(m.ImagePullPolicy),
		Command:         initContainerCommand(m.meshService),
		VolumeMounts:    initContainerVolumeMounts,
	}

	m.deploy.Spec.Template.Spec.InitContainers = injectContainers(m.deploy.Spec.Template.Spec.InitContainers, initContainer)
}

func (m *DeploymentModifier) adaptAppContainerSpec() error {
	containers := m.deploy.Spec.Template.Spec.Containers
	if len(containers) == 0 {
		return errors.Errorf("zero containers")
	}

	// NOTE: m.meshService.AppContainerName must not be empty after setupMeshService.
	appContainer, existed := findContainer(m.deploy.Spec.Template.Spec.Containers, m.meshService.AppContainerName)
	if !existed {
		return errors.Errorf("container %s not found", m.meshService.AppContainerName)
	}

	appContainer.VolumeMounts = injectVolumeMounts(appContainer.VolumeMounts, appContainerVolumeMounts...)
	appContainer.Env = injectEnvVars(appContainer.Env, appContainerEnvs...)

	m.deploy.Spec.Template.Spec.Containers = injectContainers(m.deploy.Spec.Template.Spec.Containers, *appContainer)

	return nil
}

func (m *DeploymentModifier) injectSidecarContainer() {
	sidecarContainer := corev1.Container{
		Name:            sidecarContainerName,
		Image:           m.completeImageURL(sidecarContainerImageName(m.Runtime)),
		ImagePullPolicy: corev1.PullPolicy(m.ImagePullPolicy),
		Command:         sidecarContainerCmd,
		VolumeMounts:    sidecarContainerVolumeMounts,
		Env:             siecarContainerEnvs,
		Ports:           sidecarContainerPorts,
	}

	m.deploy.Spec.Template.Spec.Containers = injectContainers(m.deploy.Spec.Template.Spec.Containers, sidecarContainer)
}

func (m *DeploymentModifier) completeImageURL(imageName string) string {
	return filepath.Join(m.ImageRegistryURL, imageName)
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

// findContainer returns the copy of the container,
// which means it won't change the original container when changing the result.
func findContainer(containers []corev1.Container, name string) (*corev1.Container, bool) {
	for _, container := range containers {
		if container.Name == name {
			return &container, true
		}
	}
	return nil, false
}

func injectVolumeMounts(volumeMounts []corev1.VolumeMount, elems ...corev1.VolumeMount) []corev1.VolumeMount {
	for _, elem := range elems {
		replaced := false
		for i, existedVolumeMount := range volumeMounts {
			if existedVolumeMount.Name == elem.Name {
				volumeMounts[i] = elem
				replaced = true
			}
		}
		if !replaced {
			volumeMounts = append(volumeMounts, elem)
		}
	}

	return volumeMounts
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

func injectContainerPorts(containerPorts []corev1.ContainerPort, elems ...corev1.ContainerPort) []corev1.ContainerPort {
	for _, elem := range elems {
		replaced := false
		for i, existedContainerPort := range containerPorts {
			if existedContainerPort.Name == elem.Name {
				containerPorts[i] = elem
				replaced = true
			}
		}
		if !replaced {
			containerPorts = append(containerPorts, elem)
		}
	}

	return containerPorts
}
