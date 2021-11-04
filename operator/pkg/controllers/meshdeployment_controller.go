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

package controllers

import (
	"context"

	meshv1beta1 "github.com/megaease/easemesh/mesh-operator/pkg/api/v1beta1"
	"github.com/megaease/easemesh/mesh-operator/pkg/base"
	"github.com/megaease/easemesh/mesh-operator/pkg/sidecarinjector"
	"github.com/megaease/easemesh/mesh-operator/pkg/syncer"

	"github.com/imdario/mergo"
	"github.com/pkg/errors"
	v1 "k8s.io/api/apps/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

// MeshDeploymentReconciler reconciles a MeshDeployment object
type MeshDeploymentReconciler struct {
	*base.Runtime
}

// +kubebuilder:rbac:groups=mesh.megaease.com,resources=meshdeployments,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=mesh.megaease.com,resources=meshdeployments/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=mesh.megaease.com,resources=meshdeployments/finalizers,verbs=update
// +kubebuilder:rbac:groups=apps,resources=deployments,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=core,resources=pods,verbs=get;list;

// Reconcile reconciles MeshDeployment.
func (r *MeshDeploymentReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	r.Log.WithValues("MeshDeploymentID", req.NamespacedName)

	meshDeploy := &meshv1beta1.MeshDeployment{}
	err := r.Client.Get(ctx, req.NamespacedName, meshDeploy)
	if err != nil {
		if apierrors.IsNotFound(err) {
			r.Log.Info("MeshDeployment not found", "id", req.NamespacedName)
			return reconcile.Result{}, nil
		}
		r.Log.Error(err, "get MeshDeployment", "id", req.NamespacedName)
		return reconcile.Result{}, err
	}

	deploy := &v1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      meshDeploy.Name,
			Namespace: meshDeploy.Namespace,
		},
	}

	r.Log.Info("syncing MeshDeployment", "id", req.NamespacedName)

	mutateFn := func() error {
		sourceDeploySpec := meshDeploy.Spec.Deploy.DeploymentSpec

		err := mergo.Merge(&deploy.Spec, &sourceDeploySpec, mergo.WithOverride)
		if err != nil {
			return errors.Wrap(err, "merge MeshDeployment into Deployment failed")
		}

		// FIXME: The decoder of client.Get() won't unmarhsal the Labels strangely.
		// Now updating vendors will cause kinds of broken dependencies.
		// Reference:
		//   https://github.com/kubernetes/klog/issues/253
		//   https://github.com/kubernetes/klog/pull/242
		//   https://github.com/kubernetes-sigs/controller-runtime/issues/1538
		deploy.Spec.Template.ObjectMeta.Labels = sourceDeploySpec.Selector.MatchLabels

		service := &sidecarinjector.MeshService{
			Name:             meshDeploy.Name,
			Labels:           meshDeploy.Spec.Service.Labels,
			AppContainerName: meshDeploy.Spec.Service.AppContainerName,
			AliveProbeURL:    meshDeploy.Spec.Service.AliveProbeURL,
			ApplicationPort:  meshDeploy.Spec.Service.ApplicationPort,
		}
		injector := sidecarinjector.New(r.Runtime, service, &deploy.Spec.Template.Spec)

		return injector.Inject()
	}

	meshDeploymentSyncer := syncer.New(r.Runtime, meshDeploy, deploy, mutateFn)
	err = syncer.Sync(context.TODO(), meshDeploymentSyncer, r.Recorder)
	if err != nil {
		r.Log.V(1).Error(err, "sync MeshDeployment")
	}

	return ctrl.Result{}, err
}

// SetupWithManager sets up the controller with the Manager.
func (r *MeshDeploymentReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&meshv1beta1.MeshDeployment{}).
		Owns(&v1.Deployment{}).
		Complete(r)
}
