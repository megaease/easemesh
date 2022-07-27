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

package sidecarinjector

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/megaease/easemesh/mesh-operator/pkg/base"
	"github.com/megaease/easemesh/mesh-operator/pkg/util/labelstool"

	"github.com/pkg/errors"
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
	initContainerImageName = func(customImage string, spec *meshControllerSpec) string {
		if customImage != "" {
			return customImage
		}

		if spec.AgentInitializerImageName != "" {
			return spec.AgentInitializerImageName
		}

		return "megaease/easeagent-initializer:latest"
	}

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
	appContainerJavaEnvValue = func(spec *meshControllerSpec) string {
		log4jConfigName := "log4j2.xml"
		if spec.Log4jConfigName != "" {
			log4jConfigName = spec.Log4jConfigName
		}
		return fmt.Sprintf(" -javaagent:%s/easeagent.jar -Deaseagent.log.conf=%s/%s ",
			appContainerAgentVolumeMountPath, appContainerAgentVolumeMountPath, log4jConfigName)
	}

	// Sidecar container stuff.
	sidecarContainerName      = "easemesh-sidecar"
	sidecarContainerImageName = func(customImage string, spec *meshControllerSpec) string {
		if customImage != "" {
			return customImage
		}
		if spec.SidecarImageName != "" {
			return spec.SidecarImageName
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
	sidecarContainerEnvs = []corev1.EnvVar{
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

func initContainerCommand(service *MeshService) []string {
	const cmdTemplate = `set -e
cp -r /easeagent-volume/* %s

echo 'name: %s
cluster-name: easemesh-control-plane
cluster-role: secondary
cluster-request-timeout: 10s
cluster:
  primary-listen-peer-urls: http://easemesh-control-plane-service.easemesh:2380
labels:
  mesh-alive-probe: %s
  mesh-application-port: %d
  mesh-service-labels: %s
  mesh-service-name: %s
' > %s`

	cmd := fmt.Sprintf(cmdTemplate,
		initContainerAgentVolumeMountPath,

		service.Name,

		service.AliveProbeURL,
		service.ApplicationPort,
		labelstool.Marshal(service.Labels),
		service.Name,

		initContainerSidecarConfigPath)

	return []string{"sh", "-c", cmd}
}

type (
	// SidecarInjector is sidecar injector for pod.
	SidecarInjector struct {
		runtime     *base.Runtime
		dynamicSpec *dynamicSpec
		meshService *MeshService
		pod         *corev1.PodSpec
	}

	// MeshService descirbes the service for SidecarInjector.
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

		// InitContainerImage could overlap the default image of the init container
		InitContainerImage string

		// SidecarImage could overlap the default image of the sidecar
		SidecarImage string
	}
)

// New creates a SidecarInjector.
func New(baseRuntime *base.Runtime, meshService *MeshService, pod *corev1.PodSpec) *SidecarInjector {
	return &SidecarInjector{
		runtime:     baseRuntime,
		dynamicSpec: newDynamicSpec(baseRuntime),
		meshService: meshService,
		pod:         pod,
	}
}

// Inject injects sidecar to the pod.
// It is idempotent.
func (m *SidecarInjector) Inject() error {
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

func (m *SidecarInjector) setupMeshService() error {
	if len(m.pod.Containers) == 0 {
		return fmt.Errorf("empty containers")
	}

	var container *corev1.Container
	if m.meshService.AppContainerName == "" {
		for i, c := range m.pod.Containers {
			// NOTE: Kubernetes will append renamed app container
			// behind existed sidecar container, so we need to ignore.
			if c.Name == sidecarContainerName {
				continue
			}

			container = &m.pod.Containers[i]
			m.meshService.AppContainerName = container.Name
			break
		}
		if container == nil {
			return errors.Errorf("no app container")
		}
	} else {
		var exists bool
		container, exists = findContainer(m.pod.Containers, m.meshService.AppContainerName)
		if !exists {
			return errors.Errorf("container %s not found", m.meshService.AppContainerName)
		}
	}

	if m.meshService.AppContainerName == sidecarContainerName {
		return errors.Errorf("app container name is conflict with sidecar: %s", sidecarContainerName)
	}

	if m.meshService.ApplicationPort == 0 {
		if len(container.Ports) == 0 {
			return errors.Errorf("container %s got zero container port", container.Name)
		}
		m.meshService.ApplicationPort = uint16(container.Ports[0].ContainerPort)
	}

	return nil
}

func (m *SidecarInjector) injectVolumes(volumes ...corev1.Volume) {
	for _, volume := range volumes {
		replaced := false
		for i, existedVolume := range m.pod.Volumes {
			if existedVolume.Name == volume.Name {
				m.pod.Volumes[i] = volume
				replaced = true
				break
			}
		}

		if !replaced {
			m.pod.Volumes = append(m.pod.Volumes, volume)
		}
	}
}

func (m *SidecarInjector) injectInitContainer() {
	initContainer := corev1.Container{
		Name:            initContainerName,
		Image:           m.completeImageURL(initContainerImageName(m.meshService.InitContainerImage, m.dynamicSpec.spec())),
		ImagePullPolicy: corev1.PullPolicy(m.dynamicSpec.spec().ImagePullPolicy),
		Command:         initContainerCommand(m.meshService),
		VolumeMounts:    initContainerVolumeMounts,
	}

	m.pod.InitContainers = injectContainers(m.pod.InitContainers, initContainer)
}

func (m *SidecarInjector) adaptAppContainerSpec() error {
	containers := m.pod.Containers
	if len(containers) == 0 {
		return errors.Errorf("zero containers")
	}

	// NOTE: m.meshService.AppContainerName must not be empty after setupMeshService.
	appContainer, existed := findContainer(m.pod.Containers, m.meshService.AppContainerName)
	if !existed {
		return errors.Errorf("container %s not found", m.meshService.AppContainerName)
	}

	appContainer.VolumeMounts = injectVolumeMounts(appContainer.VolumeMounts, appContainerVolumeMounts...)

	appContainerEnvs := []corev1.EnvVar{
		{
			Name:  appContainerJavaEnvName,
			Value: appContainerJavaEnvValue(m.dynamicSpec.spec()),
		},
	}

	appContainer.Env = injectEnvVars(appContainer.Env, appContainerEnvs...)

	m.pod.Containers = injectContainers(m.pod.Containers, *appContainer)

	return nil
}

func (m *SidecarInjector) injectSidecarContainer() {
	sidecarContainer := corev1.Container{
		Name:            sidecarContainerName,
		Image:           m.completeImageURL(sidecarContainerImageName(m.meshService.SidecarImage, m.dynamicSpec.spec())),
		ImagePullPolicy: corev1.PullPolicy(m.dynamicSpec.spec().ImagePullPolicy),
		Command:         sidecarContainerCmd,
		VolumeMounts:    sidecarContainerVolumeMounts,
		Env:             sidecarContainerEnvs,
		Ports:           sidecarContainerPorts,
	}

	m.pod.Containers = injectContainers(m.pod.Containers, sidecarContainer)
}

func (m *SidecarInjector) completeImageURL(imageName string) string {
	if strings.Contains(imageName, ".") { // imageName with repository name, we don't add prefix
		return imageName
	}
	return filepath.Join(m.dynamicSpec.spec().ImageRegistryURL, imageName)
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
