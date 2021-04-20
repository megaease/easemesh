package resourcesyncer

import (
	"github.com/go-logr/logr"
	"github.com/go-test/deep"
	"github.com/imdario/mergo"
	"github.com/megaease/easemesh/mesh-operator/pkg/api/v1beta1"
	"github.com/megaease/easemesh/mesh-operator/pkg/syncer"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
	v1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"net/url"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"strconv"
	"strings"
)

const (
	agentVolumeName      = "easeagent-volume"
	agentVolumeMountPath = "/easeagent-volume"

	sidecarParamsVolumeName      = "sidecar-params-volume"
	sidecarParamsVolumeMountPath = "/sidecar-params-volume"
	sidecarInitContainerName     = "easegateway-sidecar-initializer"

	agentInitContainerName      = "easeagent-initializer"
	agentInitContainerImage     = "192.168.50.105:5001/megaease/easeagent-initializer:latest"
	agentInitContainerMountPath = "/easeagent-share-volume"

	easeAgentJar       = " -javaagent:" + agentVolumeMountPath + "/easeagent.jar -Deaseagent.log.conf=" + agentVolumeMountPath + "/log4j2.xml "
	jolokiaAgentJar    = " -javaagent:" + agentVolumeMountPath + "/jolokia.jar "
	javaAgentJarOption = easeAgentJar + jolokiaAgentJar

	javaToolOptionsEnvName = "JAVA_TOOL_OPTIONS"
	podIPEnvName           = "APPLICATION_IP"
	podNameEnvName         = "POD_NAME"

	k8sPodIPFieldPath   = "status.podIP"
	k8sPodNameFieldPath = "metadata.name"

	sidecarImageName                = "192.168.50.105:5001/megaease/easegateway:server-sidecar"
	sidecarContainerName            = "easegateway-sidecar"
	sidecarMountPath                = "/easegateway-sidecar"
	sidecarIngressPortName          = "sidecar-ingress"
	sidecarIngressPortContainerPort = 13001

	sidecarEgressPortName         = "sidecar-egress"
	sidecarEressPortContainerPort = 13002

	defaultJMXAliveProbe        = "http://localhost:8778/jolokia/exec/com.megaease.easeagent:type=ConfigManager/healthz"
	defaultAgentHttpServerProbe = "http://localhost:9900/health"

	clusterRoleReader           = "reader"
	defaultClusterRole          = clusterRoleReader
	defaultRequestTimeoutSecond = "10s"
	defaultName                 = "eg-name"

	sideCarMeshServicenameLabel = "mesh-servicename"
	sideCarAliveProbeLabel      = "alive-probe"
	sideCarApplicationPortLabel = "application-port"
	meshServiceLabelsLabel      = "mesh-service-labels"
)

type sideCarParams struct {
	ClusterJoinUrls       string            `yaml:"cluster-join-urls"`
	ClusterRequestTimeout string            `yaml:"cluster-request-timeout"`
	ClusterRole           string            `yaml:"cluster-role"`
	ClusterName           string            `yaml:"cluster-name"`
	Labels                map[string]string `yaml: "Labels"`
}

func (params *sideCarParams) String() string {

	str := " "
	for k, v := range params.Labels {
		str += " --Labels=" + k + "=" + v
	}

	str += " --cluster-request-timeout=" + params.ClusterRequestTimeout
	str += " --cluster-role=" + params.ClusterRole
	str += " --cluster-join-urls=" + params.ClusterJoinUrls
	str += " --cluster-name=" + params.ClusterName
	return str
}

func (params *sideCarParams) Yaml() (string, error) {
	bytes, err := yaml.Marshal(params)
	if err != nil {
		return "", errors.Errorf("obj should be a deployment but is a %T", err)
	}
	return string(bytes), nil
}

type deploySyncer struct {
	meshDeployment *v1beta1.MeshDeployment
	sideCarImage   string
	clusterJoinURL string
	clusterName    string
	scheme         *runtime.Scheme
	client         client.Client
}

// NewDeploymentSyncer return a syncer of the deployment, our operator will
// inject sidecar into the sub deployment spec of the MeshDeployment
func NewDeploymentSyncer(c client.Client, meshDeploy *v1beta1.MeshDeployment,
	scheme *runtime.Scheme, clusterJoinURL string, clusterName string, log logr.Logger) syncer.Interface {
	newSyncer := &deploySyncer{
		meshDeployment: meshDeploy,
		sideCarImage:   sidecarImageName,
		client:         c,
		clusterJoinURL: clusterJoinURL,
		clusterName:    clusterName,
	}

	obj := &v1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      meshDeploy.Name,
			Namespace: meshDeploy.Namespace,
		},
	}
	return syncer.New("Deployment", c, meshDeploy, obj, scheme, log, func() error {
		previous := obj.DeepCopy()
		err := newSyncer.realSyncFn(obj)
		diff := deep.Equal(previous, obj)
		log.V(1).Info("Diff", "diff", diff)
		return err
	})
}

func (d *deploySyncer) realSyncFn(obj client.Object) error {
	deploy, ok := obj.(*v1.Deployment)
	if !ok {
		return errors.Errorf("obj should be a deployment but is a %T", obj)
	}

	sourceDeploySpec := d.meshDeployment.Spec.Deploy.DeploymentSpec

	deploy.Name = d.meshDeployment.Name
	deploy.Namespace = d.meshDeployment.Namespace
	err := mergo.Merge(&deploy.Spec, &sourceDeploySpec, mergo.WithOverride)
	if err != nil {
		return errors.Wrap(err, "merge meshDeployment failed")
	}

	// FIXME: Labels in metadata of PodTemplate will be discarding by unknown reason, we temporarily
	// complement it with matchLabel of v1.DeploymentSpec

	if deploy.Spec.Template.ObjectMeta.Labels == nil {
		deploy.Spec.Template.ObjectMeta.Labels = d.meshDeployment.Spec.Deploy.DeploymentSpec.Selector.MatchLabels
	}

	d.injectVolumes(deploy)

	err = d.completeAppContainerSpec(deploy)
	if err != nil {
		return errors.Wrap(err, "Complete Application Container error")
	}

	err = d.injectInitContainers(deploy)
	if err != nil {
		return errors.Wrap(err, "inject InitContainer error")
	}

	err = d.injectSideCarSpec(deploy)
	if err != nil {
		return errors.Wrap(err, "inject side car error")
	}

	return nil
}

func (d *deploySyncer) injectVolumes(deploy *v1.Deployment) {
	d.injectVolumeIntoDeployment(deploy, easeAgentVolume)
	d.injectVolumeIntoDeployment(deploy, sideCarParamsVolume)
}

func (d *deploySyncer) injectVolumeIntoDeployment(deploy *v1.Deployment, fn func() corev1.Volume) {
	volume := fn()
	if len(deploy.Spec.Template.Spec.Volumes) == 0 {
		deploy.Spec.Template.Spec.Volumes = []corev1.Volume{volume}
		return
	}
	for index, v := range deploy.Spec.Template.Spec.Volumes {
		if v.Name == volume.Name {
			deploy.Spec.Template.Spec.Volumes[index] = volume
			return
		}
	}
	deploy.Spec.Template.Spec.Volumes = append(deploy.Spec.Template.Spec.Volumes, volume)
}

// completeAppContainerSpec add volumeMounts for mount AgentVolume and declare env for Java Application
func (d *deploySyncer) completeAppContainerSpec(deploy *v1.Deployment) error {

	appContainer, err := d.getAppContainer(deploy)
	if err != nil {
		return err
	}

	d.injectVolumeMountIntoContainer(appContainer, agentVolumeName, easeAgentVolumeMount)
	d.injectEnvIntoContainer(appContainer, javaToolOptionsEnvName, javaToolsOptionEnv)
	return nil
}

func (d *deploySyncer) injectSideCarSpec(deploy *v1.Deployment) error {

	sideCarContainer := corev1.Container{}
	err := d.completeSideCarSpec(deploy, &sideCarContainer)
	if err != nil {
		return err
	}

	if len(deploy.Spec.Template.Spec.Containers) == 0 {
		deploy.Spec.Template.Spec.Containers = []corev1.Container{sideCarContainer}
		return nil
	}

	for index, container := range deploy.Spec.Template.Spec.Containers {
		if container.Name == sidecarContainerName {
			deploy.Spec.Template.Spec.Containers[index] = sideCarContainer
			return nil
		}
	}

	deploy.Spec.Template.Spec.Containers = append(deploy.Spec.Template.Spec.Containers, sideCarContainer)
	return nil
}

func (d *deploySyncer) completeSideCarSpec(deploy *v1.Deployment, sideCarContainer *corev1.Container) error {

	sideCarContainer.Name = sidecarContainerName

	command := "/opt/easegateway/bin/easegateway-server -f /easegateway-sidecar/eg-sidecar.yaml"
	sideCarContainer.Command = []string{"/bin/sh", "-c", command}
	sideCarContainer.Image = d.sideCarImage
	sideCarContainer.ImagePullPolicy = corev1.PullAlways
	d.injectPortIntoContainer(sideCarContainer, sidecarIngressPortName, sideCarIngressPort)
	d.injectPortIntoContainer(sideCarContainer, sidecarEgressPortName, sideCarEgressPort)
	d.injectEnvIntoContainer(sideCarContainer, podIPEnvName, podIPEnv)
	err := d.injectSidecarVolumeMounts(sideCarContainer, sidecarMountPath)
	return err
}

func (d *deploySyncer) initSideCarParams() (*sideCarParams, error) {
	params := &sideCarParams{}
	params.ClusterRole = defaultClusterRole
	params.ClusterRequestTimeout = defaultRequestTimeoutSecond

	labelSlice := []string{}
	for key, value := range d.meshDeployment.Spec.Service.Labels {
		labelSlice = append(labelSlice, key+"="+value)
	}

	meshServiceLabels := url.QueryEscape(strings.Join(labelSlice, "&"))

	labels := make(map[string]string)
	labels[sideCarMeshServicenameLabel] = d.meshDeployment.Spec.Service.Name
	labels[sideCarAliveProbeLabel] = defaultAgentHttpServerProbe
	labels[sideCarApplicationPortLabel] = ""
	labels[meshServiceLabelsLabel] = meshServiceLabels

	params.Labels = labels
	params.ClusterJoinUrls = d.clusterJoinURL
	params.ClusterName = d.clusterName
	return params, nil
}

func (d *deploySyncer) injectInitContainers(deploy *v1.Deployment) error {
	err := d.injectInitContainersIntoDeployment(deploy, agentInitContainerImage, d.easeAgentInitContainer)
	if err != nil {
		return errors.Wrap(err, "inject EaseAgent InitContainer error")
	}

	err = d.injectInitContainersIntoDeployment(deploy, sidecarImageName, d.sidecarInitContainer)
	if err != nil {
		return errors.Wrap(err, "inject sidecar InitContainer error")
	}
	return nil
}

func (d *deploySyncer) injectInitContainersIntoDeployment(deploy *v1.Deployment, containerImageName string, fn func(deploy *v1.Deployment) (corev1.Container, error)) error {

	initContainer, err := fn(deploy)
	if err != nil {
		return err
	}
	initContainers := deploy.Spec.Template.Spec.InitContainers
	if len(initContainers) == 0 {
		deploy.Spec.Template.Spec.InitContainers = []corev1.Container{initContainer}
	} else {
		for index, container := range initContainers {
			if container.Image == containerImageName {
				deploy.Spec.Template.Spec.InitContainers[index] = initContainer
				return nil
			}
		}
		deploy.Spec.Template.Spec.InitContainers = append(deploy.Spec.Template.Spec.InitContainers, initContainer)
	}
	return nil
}

func (d *deploySyncer) easeAgentInitContainer(deploy *v1.Deployment) (corev1.Container, error) {

	initContainer := corev1.Container{}

	initContainer.Name = agentInitContainerName
	initContainer.Image = agentInitContainerImage
	initContainer.ImagePullPolicy = corev1.PullAlways

	command := "cp -r " + agentVolumeMountPath + "/. " + agentInitContainerMountPath
	initContainer.Command = []string{"/bin/sh", "-c", command}

	err := d.injectAgentVolumeMounts(&initContainer, agentInitContainerMountPath)
	if err != nil {
		return initContainer, errors.Wrap(err, "inject agent volumeMounts error")
	}
	return initContainer, nil

}

func (d *deploySyncer) sidecarInitContainer(deploy *v1.Deployment) (corev1.Container, error) {

	initContainer := corev1.Container{}

	initContainer.Name = sidecarInitContainerName
	initContainer.Image = sidecarImageName
	initContainer.ImagePullPolicy = corev1.PullAlways

	params, err := d.initSideCarParams()
	if err != nil {
		return initContainer, err
	}

	appContainer, err := d.getAppContainer(deploy)
	if err != nil {
		return initContainer, err
	}

	if len(appContainer.Ports) != 0 {
		port := appContainer.Ports[0].ContainerPort
		params.Labels[sideCarApplicationPortLabel] = strconv.Itoa(int(port))
	}

	livenessProbe := appContainer.LivenessProbe
	if livenessProbe != nil && livenessProbe.HTTPGet != nil {
		host := livenessProbe.HTTPGet.Host
		port := livenessProbe.HTTPGet.Port
		path := livenessProbe.HTTPGet.Path
		aliveProbeURL := "http://" + host + port.StrVal + path
		params.Labels[sideCarAliveProbeLabel] = aliveProbeURL
	}

	d.injectEnvIntoContainer(&initContainer, podNameEnvName, podNameEnv)
	s, err := params.Yaml()
	if err != nil {
		return initContainer, err
	}

	command := "echo name: $POD_NAME >> /opt/eg-sidecar.yaml; echo '" + s + "' >> /opt/eg-sidecar.yaml; cp -r /opt/. " + sidecarParamsVolumeMountPath
	initContainer.Command = []string{"/bin/sh", "-c", command}

	d.injectVolumeMountIntoContainer(&initContainer, sidecarParamsVolumeName, sidecarVolumeMount)

	return initContainer, nil

}

// injectAgentVolumeMounts add volumeMounts for mount AgentVolume which containing the jar into container
func (d *deploySyncer) injectAgentVolumeMounts(container *corev1.Container, mountPath string) error {

	volumeMount := corev1.VolumeMount{}
	volumeMount.Name = agentVolumeName
	volumeMount.MountPath = mountPath

	if len(container.VolumeMounts) == 0 {
		container.VolumeMounts = []corev1.VolumeMount{volumeMount}
		return nil
	}
	for index, vm := range container.VolumeMounts {
		if vm.Name == agentVolumeName {
			container.VolumeMounts[index] = volumeMount
			return nil
		}
	}
	container.VolumeMounts = append(container.VolumeMounts, volumeMount)
	return nil
}

func (d *deploySyncer) injectSidecarVolumeMounts(container *corev1.Container, mountPath string) error {

	volumeMount := corev1.VolumeMount{}
	volumeMount.Name = sidecarParamsVolumeName
	volumeMount.MountPath = mountPath

	if len(container.VolumeMounts) == 0 {
		container.VolumeMounts = []corev1.VolumeMount{volumeMount}
		return nil
	}
	for index, vm := range container.VolumeMounts {
		if vm.Name == sidecarParamsVolumeName {
			container.VolumeMounts[index] = volumeMount
			return nil
		}
	}
	container.VolumeMounts = append(container.VolumeMounts, volumeMount)
	return nil
}

func (d *deploySyncer) getAppContainer(deploy *v1.Deployment) (*corev1.Container, error) {
	if d.meshDeployment.Spec.Service.AppContainerName == "" {
		return &d.meshDeployment.Spec.Deploy.Template.Spec.Containers[0], nil
	}
	for index, container := range deploy.Spec.Template.Spec.Containers {
		if container.Name == d.meshDeployment.Spec.Service.AppContainerName {
			return &deploy.Spec.Template.Spec.Containers[index], nil
		}
	}
	return nil, errors.Errorf("Application container do not exists. Please confirm application container name is %s.", d.meshDeployment.Spec.Service.AppContainerName)
}

func (d *deploySyncer) injectEnvIntoContainer(container *corev1.Container, envName string, fn func() corev1.EnvVar) {
	env := fn()
	if len(container.Env) == 0 {
		container.Env = []corev1.EnvVar{env}
		return
	}
	for index, env := range container.Env {
		if env.Name == envName {
			container.Env[index] = env
			return
		}
	}
	container.Env = append(container.Env, env)

}

func (d *deploySyncer) injectPortIntoContainer(container *corev1.Container, portName string, fn func() corev1.ContainerPort) {
	port := fn()
	if len(container.Ports) == 0 {
		container.Ports = []corev1.ContainerPort{port}
		return
	}
	for index, p := range container.Ports {
		if p.Name == portName {
			container.Ports[index] = port
			return
		}
	}
	container.Ports = append(container.Ports, port)

}

func (d *deploySyncer) injectVolumeMountIntoContainer(container *corev1.Container, volumeName string, fn func() corev1.VolumeMount) {
	volumeMount := fn()
	if len(container.VolumeMounts) == 0 {
		container.VolumeMounts = []corev1.VolumeMount{volumeMount}
		return
	}

	for index, vm := range container.VolumeMounts {
		if vm.Name == volumeName {
			container.VolumeMounts[index] = volumeMount
			return
		}
	}
	container.VolumeMounts = append(container.VolumeMounts, volumeMount)
}

func sideCarParamsVolume() corev1.Volume {
	volume := corev1.Volume{}
	volume.Name = sidecarParamsVolumeName
	volume.EmptyDir = &corev1.EmptyDirVolumeSource{}
	return volume
}

func easeAgentVolume() corev1.Volume {
	volume := corev1.Volume{}
	volume.Name = agentVolumeName
	volume.EmptyDir = &corev1.EmptyDirVolumeSource{}
	return volume
}

func javaToolsOptionEnv() corev1.EnvVar {
	env := corev1.EnvVar{
		Name:  javaToolOptionsEnvName,
		Value: javaAgentJarOption,
	}
	return env
}

func podIPEnv() corev1.EnvVar {
	varSource := &corev1.EnvVarSource{
		FieldRef: &corev1.ObjectFieldSelector{
			FieldPath: k8sPodIPFieldPath,
		},
	}

	env := corev1.EnvVar{
		Name:      podIPEnvName,
		ValueFrom: varSource,
	}
	return env
}

func podNameEnv() corev1.EnvVar {
	varSource := &corev1.EnvVarSource{
		FieldRef: &corev1.ObjectFieldSelector{
			FieldPath: k8sPodNameFieldPath,
		},
	}

	env := corev1.EnvVar{
		Name:      podNameEnvName,
		ValueFrom: varSource,
	}
	return env
}

func sideCarIngressPort() corev1.ContainerPort {
	port := corev1.ContainerPort{
		Name:          sidecarIngressPortName,
		ContainerPort: sidecarIngressPortContainerPort,
	}
	return port
}

func sideCarEgressPort() corev1.ContainerPort {
	port := corev1.ContainerPort{
		Name:          sidecarEgressPortName,
		ContainerPort: sidecarEressPortContainerPort,
	}
	return port
}

func easeAgentVolumeMount() corev1.VolumeMount {
	volumeMount := corev1.VolumeMount{}
	volumeMount.Name = agentVolumeName
	volumeMount.MountPath = agentVolumeMountPath
	return volumeMount
}

func sidecarVolumeMount() corev1.VolumeMount {
	volumeMount := corev1.VolumeMount{}
	volumeMount.Name = sidecarParamsVolumeName
	volumeMount.MountPath = sidecarParamsVolumeMountPath
	return volumeMount
}
