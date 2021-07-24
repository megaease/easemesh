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
	"strconv"

	"github.com/megaease/easemesh/mesh-operator/pkg/base"
	"github.com/megaease/easemesh/mesh-operator/pkg/syncer"
	"github.com/pkg/errors"

	v1 "k8s.io/api/apps/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

const (
	annotationPrefix              = "mesh.megaease.com/"
	annotationEnableKey           = annotationPrefix + "enable"
	annotationServiceNameKey      = annotationPrefix + "service-name"
	annotationAppContainerNameKey = annotationPrefix + "app-container-name"
	annotationApplicationPortKey  = annotationPrefix + "application-port"
	annotationAliveProbeURLKey    = annotationPrefix + "alive-probe-url"
)

// DeploymentReconciler reconciles Deployment object with mesh condition.
type DeploymentReconciler struct {
	*base.Runtime
}

// +kubebuilder:rbac:groups=apps,resources=deployments,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=core,resources=pods,verbs=get;list;

// Reconcile reconciles Deloyment with mesh condition.
func (r *DeploymentReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	r.Log.WithValues("DeploymentID", req.NamespacedName)

	deploy := &v1.Deployment{}
	err := r.Client.Get(context.TODO(), req.NamespacedName, deploy)
	if err != nil {
		if apierrors.IsNotFound(err) {
			r.Log.Info("Deployment not found", "id", req.NamespacedName)
			return reconcile.Result{}, nil
		}
		r.Log.Error(err, "get Deployment", "id", req.NamespacedName)
		return reconcile.Result{}, err
	}

	if deploy.Annotations[annotationEnableKey] != "true" {
		return reconcile.Result{}, nil
	}

	r.Log.Info("syncing Deployment", "id", req.NamespacedName)

	mutateFn := func() error {
		applicationPortValue := deploy.Annotations[annotationApplicationPortKey]
		var applicationPort uint16
		if applicationPortValue != "" {
			port, err := strconv.ParseUint(applicationPortValue, 10, 16)
			if err != nil {
				return errors.Wrapf(err, "parse application port %s", applicationPortValue)
			}
			applicationPort = uint16(port)
		}

		service := &meshService{
			Name:             deploy.Name,
			Labels:           deploy.Labels,
			AppContainerName: deploy.Annotations[annotationAppContainerNameKey],
			AliveProbeURL:    deploy.Annotations[annotationAliveProbeURLKey],
			ApplicationPort:  applicationPort,
		}
		modifier := newDeploymentModifier(r.Runtime, service, deploy)

		return modifier.modify()
	}

	meshDeploymentSyncer := syncer.New(r.Runtime, deploy, deploy, mutateFn)
	err = syncer.Sync(context.TODO(), meshDeploymentSyncer, r.Recorder)
	if err != nil {
		r.Log.V(1).Error(err, "sync MeshDeployment")
	}

	return ctrl.Result{}, err
}

// SetupWithManager sets up the controller with the Manager.
func (r *DeploymentReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&v1.Deployment{}).
		Owns(&v1.Deployment{}).
		Complete(r)
}
