package resourcesyncer

import (
	"github.com/go-logr/logr"
	"github.com/go-test/deep"
	"github.com/imdario/mergo"
	"github.com/megaease/easemesh/mesh-operator/pkg/api/v1beta1"
	"github.com/megaease/easemesh/mesh-operator/pkg/syncer"
	"github.com/megaease/easemesh/mesh-operator/pkg/util"
	"github.com/pkg/errors"
	v1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"strconv"
)

const (
	agentVolumeName      = "easeagent-volume"
	agentVolumeMountPath = "/easeagent-volume"

	agentInitContainerName      = "easeagent-initializer"
	agentInitContainerImage     = "192.168.50.105:5001/megaease/easeagent-initializer:latest"
	agentInitContainerMountPath = "/easeagent-share-volume"

	easeAgentJar       = "-javaagent:" + agentVolumeMountPath + "/easeagent.jar "
	jolokiaAgentJar    = "-javaagent:" + agentVolumeMountPath + "/jolokia.jar "
	javaAgentJarOption = easeAgentJar + jolokiaAgentJar

	javaToolOptionsEnvName = "JAVA_TOOL_OPTIONS"

	sideCarImageName     = "192.168.50.105:5001/megaease/easegateway:server-sidecar"
	sideCarContainerName = "easegateway-sidecar"

	defaultJMXAliveProbe = "http://localhost:8080/jolokia/exec/com.megaease.easeagent:type=ConfigManager/healthz"

	clusterRoleReader           = "reader"
	defaultClusterRole          = clusterRoleReader
	defaultRequestTimeoutSecond = 10
	defaultName                 = "eg-name"

	sideCarMeshServicenameLabel = "mesh-servicename"
	sideCarAliveProbeLabel      = "alive-probe"
)

type sideCarParams struct {
	labels                map[string]string
	name                  string
	clusterJoinUrl        string
	clusterRequestTimeout int
	clusterRole           string
}

func (params *sideCarParams) String() string {

	str := " "
	for k, v := range params.labels {
		str += " --labels=" + k + "=" + v
	}

	str += " --name=" + params.name
	str += " --cluster-request-timeout=" + strconv.Itoa(params.clusterRequestTimeout)
	str += " --cluster-role=" + params.clusterRole
	str += " --cluster-join-url=" + params.clusterJoinUrl
	return str
}

type deploySyncer struct {
	meshDeployment *v1beta1.MeshDeployment
	sideCarImage   string
	scheme         *runtime.Scheme
	client         client.Client
}

// NewDeploymentSyncer return a syncer of the deployment, our operator will
// inject sidecar into the sub deployment spec of the MeshDeployment
func NewDeploymentSyncer(c client.Client, meshDeploy *v1beta1.MeshDeployment,
	scheme *runtime.Scheme, log logr.Logger) syncer.Interface {
	newSyncer := &deploySyncer{
		meshDeployment: meshDeploy,
		sideCarImage:   sideCarImageName,
		client:         c,
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

	// FIXME: labels in metadata of PodTemplate will be discarding by unknown reason, we temporarily
	// complement it with matchLabel of v1.DeploymentSpec

	if deploy.Spec.Template.ObjectMeta.Labels == nil {
		deploy.Spec.Template.ObjectMeta.Labels = d.meshDeployment.Spec.Deploy.DeploymentSpec.Selector.MatchLabels
	}

	err = d.injectAgentVolumes(deploy)
	if err != nil {
		return errors.Wrap(err, "inject Agent Volume error")
	}

	err = d.injectEaseAgentInitContainer(deploy)
	if err != nil {
		return errors.Wrap(err, "inject EaseAgent InitContainer error")
	}

	err = d.injectSideCarSpec(deploy)
	if err != nil {
		return errors.Wrap(err, "inject side car error")
	}

	err = d.injectAgentJarIntoApp(&deploy.Spec.Template.Spec.Containers[0])
	if err != nil {
		return errors.Wrap(err, "inject Agent Jar into Application Container error")
	}

	return nil
}

func (d *deploySyncer) injectSideCarSpec(deploy *v1.Deployment) error {

	containers := deploy.Spec.Template.Spec.Containers
	for _, container := range containers {
		if container.Name == sideCarContainerName {
			return nil
		}
	}

	// Eg SideCar Params
	params, err := d.initSideCarParams()
	if err != nil {
		return err
	}

	sideCarContainer := corev1.Container{}
	sideCarContainer.Name = sideCarContainerName
	sideCarContainer.Image = d.sideCarImage

	if len(sideCarContainer.Args) == 0 {
		sideCarContainer.Args = []string{params.String()}
	} else {
		sideCarContainer.Args = append(sideCarContainer.Args, params.String())
	}
	deploy.Spec.Template.Spec.Containers = append(containers, sideCarContainer)
	return nil
}

func (d *deploySyncer) initSideCarParams() (*sideCarParams, error) {
	params := &sideCarParams{}
	params.clusterRole = defaultClusterRole
	params.name = defaultName
	params.clusterRequestTimeout = defaultRequestTimeoutSecond

	var aliveProbeURL string
	livenessProbe := d.meshDeployment.Spec.Deploy.DeploymentSpec.Template.Spec.Containers[0].LivenessProbe
	if livenessProbe != nil && livenessProbe.HTTPGet != nil {
		host := livenessProbe.HTTPGet.Host
		port := livenessProbe.HTTPGet.Port
		path := livenessProbe.HTTPGet.Path
		url := "http://" + host + port.StrVal + path
		aliveProbeURL = url
	} else {
		aliveProbeURL = defaultJMXAliveProbe
	}

	labels := make(map[string]string)
	labels[sideCarMeshServicenameLabel] = d.meshDeployment.Spec.Service.Name
	labels[sideCarAliveProbeLabel] = aliveProbeURL

	meshOperator, _ := util.GetEaseMeshOperator(d.client)
	masterJoinURL := meshOperator.GetEGMasterJoinURL(d.client)
	params.clusterJoinUrl = masterJoinURL

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
	command := "cp -r " + agentVolumeMountPath + "/. " + agentInitContainerMountPath
	initContainer.Command = []string{"/bin/sh", "-c", command}

	err := d.injectAgentVolumeMounts(&initContainer, agentInitContainerMountPath)
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
	} else {
		container.VolumeMounts = append(container.VolumeMounts, volumeMount)
	}

	return nil
}

// injectAgentJarIntoApp add volumeMounts for mount AgentVolume and declare JAVA_TOOL_OPTIONS env for Java Application
func (d *deploySyncer) injectAgentJarIntoApp(container *corev1.Container) error {

	err := d.injectAgentVolumeMounts(container, agentVolumeMountPath)
	if err != nil {
		return errors.Wrap(err, "inject agent volumeMounts error")
	}

	javaToolOptionsEnv := corev1.EnvVar{
		Name:  javaToolOptionsEnvName,
		Value: javaAgentJarOption,
	}
	if len(container.Env) == 0 {
		container.Env = []corev1.EnvVar{javaToolOptionsEnv}
	} else {
		container.Env = append(container.Env, javaToolOptionsEnv)
	}
	return nil
}
