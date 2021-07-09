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

package controllers

import (
	"context"

	meshv1beta1 "github.com/megaease/easemesh/mesh-operator/pkg/api/v1beta1"
	"github.com/megaease/easemesh/mesh-operator/pkg/controllers/resourcesyncer"
	"github.com/megaease/easemesh/mesh-operator/pkg/syncer"

	"github.com/go-logr/logr"
	"github.com/juju/errors"
	v1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

// MeshDeploymentReconciler reconciles a MeshDeployment object
type MeshDeploymentReconciler struct {
	client.Client
	Log              logr.Logger
	Scheme           *runtime.Scheme
	Recorder         record.EventRecorder
	ClusterJoinURL   string
	ImageRegistryURL string
	ClusterName      string
}

// +kubebuilder:rbac:groups=mesh.megaease.com,resources=meshdeployments,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=mesh.megaease.com,resources=meshdeployments/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=mesh.megaease.com,resources=meshdeployments/finalizers,verbs=update
// +kubebuilder:rbac:groups=apps,resources=deployments,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=core,resources=pods,verbs=get;list;

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the MeshDeployment object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.7.0/pkg/reconcile
func (r *MeshDeploymentReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = r.Log.WithValues("meshdeployment", req.NamespacedName)
	// your logic here
	meshDeploy := &meshv1beta1.MeshDeployment{}
	err := r.Client.Get(context.TODO(), req.NamespacedName, meshDeploy)
	if err != nil {
		if errors.IsNotFound(err) {
			// Object not found, return. Created objects are automatically garbage collected
			return reconcile.Result{}, nil
		}
		// Error reading the object - requeue the request
		return reconcile.Result{}, err
	}

	log := r.Log.WithValues("key", req.NamespacedName)
	log.V(1).Info("deploy is", "meshdeployment", meshDeploy)

	deploySyncer := resourcesyncer.NewDeploymentSyncer(r.Client, meshDeploy, r.Scheme, r.ClusterJoinURL, r.ClusterName, r.Log, r.ImageRegistryURL)
	err = syncer.Sync(context.TODO(), deploySyncer, r.Recorder)
	if err != nil {
		log.V(1).Info("sync deployment resource error")
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
