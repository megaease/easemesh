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
	"math/rand"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"strconv"
)

const (
	agentVolumeName      = "easeagent-volume"
	agentVolumeMountPath = "/easeagent-volume"

	agentInitContainerName      = "easeagent-initializer"
	agentInitContainerImage     = "192.168.50.105:5001/megaease/easeagent-initializer:latest"
	agentInitContainerMountPath = "/easeagent-share-volume"

	easeAgentJar       = "-javaagent:" + agentVolumeMountPath + "/easeagent.jar -Deaseagent.log.conf=" + agentVolumeMountPath + "/log4j2.xml"
	jolokiaAgentJar    = "-javaagent:" + agentVolumeMountPath + "/jolokia.jar "
	javaAgentJarOption = easeAgentJar + jolokiaAgentJar

	javaToolOptionsEnvName = "JAVA_TOOL_OPTIONS"
	podIPEnvName           = "APPLICATION_IP"

	k8sPodIPFieldPath = "status.podIP"

	sideCarImageName                = "192.168.50.105:5001/megaease/easegateway:server-sidecar"
	sideCarContainerName            = "easegateway-sidecar"
	sideCarMountPath                = "/easegateway-sidecar"
	sideCarIngressPortName          = "sidecar-ingress"
	sideCarIngressPortContainerPort = 13001

	sideCarEgressPortName         = "sidecar-egress"
	sideCarEressPortContainerPort = 13002

	defaultJMXAliveProbe = "http://localhost:8778/jolokia/exec/com.megaease.easeagent:type=ConfigManager/healthz"

	clusterRoleReader           = "reader"
	defaultClusterRole          = clusterRoleReader
	defaultRequestTimeoutSecond = "10s"
	defaultName                 = "eg-name"

	sideCarMeshServicenameLabel = "mesh-servicename"
	sideCarAliveProbeLabel      = "alive-probe"
	sideCarApplicationPortLabel = "application-port"
)

type sideCarParams struct {
	Name                  string            `yaml:"name"`
	ClusterJoinUrls       string            `yaml:"cluster-join-urls"`
	ClusterRequestTimeout string            `yaml:"cluster-request-timeout"`
	ClusterRole           string            `yaml:"cluster-role"`
	ClusterName           string            `yaml:"cluster-name"`
	Labels                map[string]string `yaml: "Labels,omitempty"`
}

func (params *sideCarParams) String() string {

	str := " "
	for k, v := range params.Labels {
		str += " --Labels=" + k + "=" + v
	}

	str += " --name=" + params.Name
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
		sideCarImage:   sideCarImageName,
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

	err = d.injectAgentVolumes(deploy)
	if err != nil {
		return errors.Wrap(err, "inject Agent Volume error")
	}

	err = d.completeAppContainerSpec(deploy)
	if err != nil {
		return errors.Wrap(err, "inject Agent Jar into Application Container error")
	}

	err = d.injectEaseAgentInitContainer(deploy)
	if err != nil {
		return errors.Wrap(err, "inject EaseAgent InitContainer error")
	}

	err = d.injectSideCarSpec(deploy)
	if err != nil {
		return errors.Wrap(err, "inject side car error")
	}

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
		if container.Name == sideCarContainerName {
			deploy.Spec.Template.Spec.Containers[index] = sideCarContainer
			return nil
		}
	}

	deploy.Spec.Template.Spec.Containers = append(deploy.Spec.Template.Spec.Containers, sideCarContainer)
	return nil
}

func (d *deploySyncer) completeSideCarSpec(deploy *v1.Deployment, sideCarContainer *corev1.Container) error {

	sideCarContainer.Name = sideCarContainerName

	command := "/opt/easegateway/bin/easegateway-server -f /easegateway-sidecar/eg-sidecar.yaml"
	sideCarContainer.Command = []string{"/bin/sh", "-c", command}
	sideCarContainer.Image = d.sideCarImage
	sideCarContainer.ImagePullPolicy = corev1.PullAlways
	d.injectPortIntoContainer(sideCarContainer, sideCarIngressPortName, sideCarIngressPort)
	d.injectPortIntoContainer(sideCarContainer, sideCarEgressPortName, sideCarEgressPort)
	d.injectEnvIntoContainer(sideCarContainer, podIPEnvName, podIPEnv)
	err := d.injectAgentVolumeMounts(sideCarContainer, sideCarMountPath)
	return err
}

func (d *deploySyncer) initSideCarParams() (*sideCarParams, error) {
	params := &sideCarParams{}
	params.ClusterRole = defaultClusterRole
	params.Name = d.meshDeployment.Spec.Service.Name + "-" + strconv.Itoa(rand.Int())
	params.ClusterRequestTimeout = defaultRequestTimeoutSecond

	labels := make(map[string]string)
	labels[sideCarMeshServicenameLabel] = d.meshDeployment.Spec.Service.Name
	labels[sideCarAliveProbeLabel] = defaultJMXAliveProbe
	labels[sideCarApplicationPortLabel] = ""

	params.Labels = labels
	params.ClusterJoinUrls = d.clusterJoinURL
	params.ClusterName = d.clusterName
	return params, nil
}

// injectAgentVolumes add a empty volume for storage agent jar
func (d *deploySyncer) injectAgentVolumes(deploy *v1.Deployment) error {

	agentVolume := corev1.Volume{}
	agentVolume.Name = agentVolumeName
	agentVolume.EmptyDir = &corev1.EmptyDirVolumeSource{}

	volumes := deploy.Spec.Template.Spec.Volumes

	if len(volumes) == 0 {
		deploy.Spec.Template.Spec.Volumes = []corev1.Volume{agentVolume}
	} else {
		for _, volume := range volumes {
			if volume.Name == agentVolumeName && volume.EmptyDir != nil {
				return nil
			}
		}
		deploy.Spec.Template.Spec.Volumes = append(deploy.Spec.Template.Spec.Volumes, agentVolume)
	}

	return nil
}

// injectEaseAgentInitContainer add a InitContainer of K8S for download agent jars
func (d *deploySyncer) injectEaseAgentInitContainer(deploy *v1.Deployment) error {

	initContainer := corev1.Container{}

	initContainer.Name = agentInitContainerName
	initContainer.Image = agentInitContainerImage
	initContainer.ImagePullPolicy = corev1.PullAlways

	params, err := d.initSideCarParams()
	if err != nil {
		return err
	}

	appContainer, err := d.getAppContainer(deploy)
	if err != nil {
		return err
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

	s, err := params.Yaml()
	if err != nil {
		return err
	}

	command := "echo '" + s + "' > /easeagent-share-volume/eg-sidecar.yaml; cp -r " + agentVolumeMountPath + "/. " + agentInitContainerMountPath
	initContainer.Command = []string{"/bin/sh", "-c", command}

	err = d.injectAgentVolumeMounts(&initContainer, agentInitContainerMountPath)
	if err != nil {
		return errors.Wrap(err, "inject agent volumeMounts error")
	}

	initContainers := deploy.Spec.Template.Spec.InitContainers
	if len(initContainers) == 0 {
		deploy.Spec.Template.Spec.InitContainers = []corev1.Container{initContainer}
	} else {
		for _, container := range initContainers {
			if container.Image == agentInitContainerImage {
				return nil
			}
		}
		deploy.Spec.Template.Spec.InitContainers = append(deploy.Spec.Template.Spec.InitContainers, initContainer)
	}

	return nil
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

// completeAppContainerSpec add volumeMounts for mount AgentVolume and declare env for Java Application
func (d *deploySyncer) completeAppContainerSpec(deploy *v1.Deployment) error {

	appContainer, err := d.getAppContainer(deploy)
	if err != nil {
		return err
	}

	err = d.injectAgentVolumeMounts(appContainer, agentVolumeMountPath)
	if err != nil {
		return errors.Wrap(err, "inject agent volumeMounts error")
	}
	d.injectEnvIntoContainer(appContainer, javaToolOptionsEnvName, javaToolsOptionEnv)
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

func sideCarIngressPort() corev1.ContainerPort {
	port := corev1.ContainerPort{
		Name:          sideCarIngressPortName,
		ContainerPort: sideCarIngressPortContainerPort,
	}
	return port
}

func sideCarEgressPort() corev1.ContainerPort {
	port := corev1.ContainerPort{
		Name:          sideCarEgressPortName,
		ContainerPort: sideCarEressPortContainerPort,
	}
	return port
}
