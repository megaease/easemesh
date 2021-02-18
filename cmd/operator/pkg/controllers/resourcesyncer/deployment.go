package resourcesyncer

import (
	"github.com/megaease/easemesh/mesh-operator/pkg/api/v1beta1"
	"github.com/megaease/easemesh/mesh-operator/pkg/syncer"

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

type deploySyncer struct {
	deploy       *v1beta1.MeshDeployment
	sideCarImage string
	scheme       *runtime.Scheme
}

// NewDeploymentSyncer return a syncer of the deployment, our operator will
// inject sidecar into the sub deployment spec of the MeshDeployment
func NewDeploymentSyncer(c client.Client, meshDeploy *v1beta1.MeshDeployment,
	scheme *runtime.Scheme, log logr.Logger) syncer.Interface {
	deploy := &deploySyncer{
		deploy:       meshDeploy,
		sideCarImage: "nginx:1.14.2",
	}

	obj := &v1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      meshDeploy.Name,
			Namespace: meshDeploy.Namespace,
		},
	}
	return syncer.New("Deployment", c, meshDeploy, obj, scheme, log, func() error {
		previous := obj.DeepCopy()
		err := deploy.realSyncFn(obj)
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

	sourceDeploySpec := d.deploy.Spec.Deploy.DeploymentSpec

	deploy.Name = d.deploy.Name
	deploy.Namespace = d.deploy.Namespace
	err := mergo.Merge(&deploy.Spec, &sourceDeploySpec, mergo.WithOverride)
	if err != nil {
		return errors.Wrap(err, "merge deploy failed")
	}

	// FIXME: labels in metadata of PodTemplate will be discarding by unknown reason, we temporarily
	// complement it with matchLabel of v1.DeploymentSpec

	if deploy.Spec.Template.ObjectMeta.Labels == nil {
		deploy.Spec.Template.ObjectMeta.Labels = d.deploy.Spec.Deploy.DeploymentSpec.Selector.MatchLabels
	}

	if len(deploy.Spec.Template.Spec.Containers) != 2 {
		deploy.Spec.Template.Spec.Containers = make([]corev1.Container, 2, 2)
	}

	err = mergo.Merge(&deploy.Spec.Template.Spec.Containers[0],
		&sourceDeploySpec.Template.Spec.Containers[0],
		mergo.WithOverride)
	if err != nil {
		return errors.Wrap(err, "copy default container error")
	}

	err = d.injectSideCarSpec(&deploy.Spec.Template.Spec.Containers[1])
	if err != nil {
		return errors.Wrap(err, "inject side car error")
	}

	return nil
}

func (d *deploySyncer) injectSideCarSpec(targetContainer *corev1.Container) error {
	// Replace with real logic, this is a stub code for test

	targetContainer.Name = "easemesh-sidecar"
	targetContainer.Image = "nginx:1.14.2"
	targetContainer.Ports = []corev1.ContainerPort{
		{
			ContainerPort: 80,
		},
	}

	return nil
}
