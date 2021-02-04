package resourcesyncer

import (
	"github.com/pkg/errors"
	v1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/megaease/easemesh/mesh-operator/api/v1beta1"
	"github.com/megaease/easemesh/mesh-operator/syncer"
)

type deploySyncer struct {
	deploy       v1beta1.MeshDeployment
	sideCarImage string
}

// NewDeploymentSyncer return a syncer of the deployment, our operator will
// inject sidecar into the sub deployment spec of the MeshDeployment
func NewDeploymentSyncer(c client.Client, meshDeploy *v1beta1.MeshDeployment, scheme *runtime.Scheme) syncer.Interface {
	deploy := &deploySyncer{}
	obj := &v1.Deployment{}
	return syncer.New("Deployment", c, meshDeploy, obj, scheme, func() error {
		return deploy.realSyncFn(obj)
	})
}

func (d *deploySyncer) realSyncFn(obj client.Object) error {
	deploy, ok := obj.(*v1.Deployment)
	if !ok {
		return errors.Errorf("obj should be a deployment but is a %T", obj)
	}
	deploy.Name = d.deploy.Name
	return nil
}
