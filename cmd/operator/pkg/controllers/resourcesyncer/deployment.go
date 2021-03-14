package resourcesyncer

import (
	"github.com/megaease/easemesh/mesh-operator/pkg/api/v1beta1"
	"github.com/megaease/easemesh/mesh-operator/pkg/syncer"
	"strconv"

	"github.com/go-logr/logr"
	"github.com/go-test/deep"
	"github.com/imdario/mergo"
	"github.com/pkg/errors"
	v1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	agentVolumeName      = "easeagent-volume"
	agentVolumeMountPath = "/easeagent-volume"

	agentInitContainerName  = "easeagent-initializer"
	agentInitContainerImage = "easeagent-initializer:latest"

	easeAgentJar       = "-javaagent:" + agentVolumeMountPath + "/easeagent.jar "
	jolokiaAgentJar    = "-javaagent:" + agentVolumeMountPath + "/jolokia.jar "
	javaAgentJarOption = easeAgentJar + jolokiaAgentJar

	javaToolOptionsEnvName = "JAVA_TOOL_OPTIONS"

	sideCarImageName = "easegateway:latest"

	defaultJMXAliveProbe = "http://localhost:8080/jolokia/exec/com.megaease.easeagent:type=ConfigManager/healthz"

	clusterRoleReader           = "reader"
	defaultClusterRole          = clusterRoleReader
	defaultRequestTimeoutSecond = 10
	defaultName                 = "eg-name"
)

type sideCarParams struct {
	meshServiceName       string
	aliveProbeURL         string
	name                  string
	clusterJoinUrl        string
	clusterRequestTimeout int
	clusterRole           string
}

type deploySyncer struct {
	meshDeployment *v1beta1.MeshDeployment
	sideCarImage   string
	scheme         *runtime.Scheme
}

// NewDeploymentSyncer return a syncer of the deployment, our operator will
// inject sidecar into the sub deployment spec of the MeshDeployment
func NewDeploymentSyncer(c client.Client, meshDeploy *v1beta1.MeshDeployment,
	scheme *runtime.Scheme, log logr.Logger) syncer.Interface {
	newSyncer := &deploySyncer{
		meshDeployment: meshDeploy,
		sideCarImage:   sideCarImageName,
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
	if len(deploy.Spec.Template.Spec.Containers) != 2 {
		newContainers := make([]corev1.Container, 2, 2)
		err := mergo.Merge(&deploy.Spec.Template.Spec.Containers[0],
			&newContainers[0],
			mergo.WithOverride)
		if err != nil {
			return errors.Wrap(err, "copy default container error")
		}
	}

	// Eg SideCar Params
	params, err := d.initSideCarParams()
	if err != nil {
		return err
	}

	sideCarContainer := corev1.Container{}
	sideCarContainer.Name = params.name
	sideCarContainer.Image = d.sideCarImage

	// TODO: Add command and args

	deploy.Spec.Template.Spec.Containers[1] = sideCarContainer
	return nil
}

func (d *deploySyncer) initSideCarParams() (*sideCarParams, error) {
	params := &sideCarParams{}
	params.clusterRole = defaultClusterRole
	params.meshServiceName = d.meshDeployment.Spec.Service.Name

	serviceLabels := d.meshDeployment.Spec.Service.Labels
	if name, ok := serviceLabels["name"]; ok {
		params.name = name
	} else {
		params.name = defaultName
	}

	clusterRequestTimeout, ok := d.meshDeployment.Labels["clusterRequestTimeout"]
	if ok {
		timeout, err := strconv.Atoi(clusterRequestTimeout)
		if err != nil {
			return nil, errors.Wrap(err, "MeshDeployment labels error, clusterRequestTimeout expected is int, but get other type")
		}
		params.clusterRequestTimeout = timeout
	} else {
		params.clusterRequestTimeout = defaultRequestTimeoutSecond
	}

	livenessProbe := d.meshDeployment.Spec.Deploy.DeploymentSpec.Template.Spec.Containers[0].LivenessProbe
	if livenessProbe != nil && livenessProbe.HTTPGet != nil {
		host := livenessProbe.HTTPGet.Host
		port := livenessProbe.HTTPGet.Port
		path := livenessProbe.HTTPGet.Path
		url := "http://" + host + port.StrVal + path
		params.aliveProbeURL = url
	} else {
		params.aliveProbeURL = defaultJMXAliveProbe
	}

	// TODO:  query cluster-join-url from eg master
	return params, nil
}

// injectAgentVolumes add a empty volume for storage agent jar
func (d *deploySyncer) injectAgentVolumes(deploy *v1.Deployment) error {

	agentVolume := corev1.Volume{}
	agentVolume.Name = agentVolumeName
	agentVolume.EmptyDir = &corev1.EmptyDirVolumeSource{}

	if len(deploy.Spec.Template.Spec.Volumes) == 0 {
		deploy.Spec.Template.Spec.Volumes = []corev1.Volume{agentVolume}
	} else {
		deploy.Spec.Template.Spec.Volumes = append(deploy.Spec.Template.Spec.Volumes, agentVolume)
	}

	return nil
}

// injectEaseAgentInitContainer add a InitContainer of K8S for download agent jars
func (d *deploySyncer) injectEaseAgentInitContainer(deploy *v1.Deployment) error {

	initContainer := corev1.Container{}

	initContainer.Name = agentInitContainerName
	initContainer.Image = agentInitContainerImage

	err := d.injectAgentVolumeMounts(&initContainer)
	if err != nil {
		return errors.Wrap(err, "inject agent volumeMounts error")
	}

	if len(deploy.Spec.Template.Spec.InitContainers) != 1 {
		deploy.Spec.Template.Spec.InitContainers = []corev1.Container{initContainer}
	} else {
		deploy.Spec.Template.Spec.InitContainers = append(deploy.Spec.Template.Spec.InitContainers, initContainer)
	}

	return nil
}

// injectAgentVolumeMounts add volumeMounts for mount AgentVolume which containing the jar into container
func (d *deploySyncer) injectAgentVolumeMounts(container *corev1.Container) error {

	volumeMount := corev1.VolumeMount{}
	volumeMount.Name = agentVolumeName
	volumeMount.MountPath = agentVolumeMountPath

	if len(container.VolumeMounts) == 0 {
		container.VolumeMounts = []corev1.VolumeMount{volumeMount}
	} else {
		container.VolumeMounts = append(container.VolumeMounts, volumeMount)
	}

	return nil
}

// injectAgentJarIntoApp add volumeMounts for mount AgentVolume and declare JAVA_TOOL_OPTIONS env for Java Application
func (d *deploySyncer) injectAgentJarIntoApp(container *corev1.Container) error {

	err := d.injectAgentVolumeMounts(container)
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
