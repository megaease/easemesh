/*
Copyright 2021 MegaEase.cn.
*/

package controllers

import (
	"context"

	"github.com/go-logr/logr"
	"github.com/juju/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	meshv1beta1 "github.com/megaease/easemesh/mesh-operator/api/v1beta1"
	"github.com/megaease/easemesh/mesh-operator/controllers/resourcesyncer"
	"github.com/megaease/easemesh/mesh-operator/syncer"
)

const (
	controllerName = "controller.meshdeployment"
)

var (
	log = logf.Log.WithName(controllerName)
)

// MeshDeploymentReconciler reconciles a MeshDeployment object
type MeshDeploymentReconciler struct {
	client.Client
	Log      logr.Logger
	Scheme   *runtime.Scheme
	Recorder record.EventRecorder
}

// +kubebuilder:rbac:groups=mesh.megaease.com,resources=meshdeployments,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=mesh.megaease.com,resources=meshdeployments/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=mesh.megaease.com,resources=meshdeployments/finalizers,verbs=update

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
	deploy := &meshv1beta1.MeshDeployment{}
	err := r.Client.Get(context.TODO(), req.NamespacedName, deploy)
	if err != nil {
		if errors.IsNotFound(err) {
			// Object not found, return. Created objects are automatically garbage collected
			return reconcile.Result{}, nil
		}
		// Error reading the object - requeue the request
		return reconcile.Result{}, err
	}

	log := log.WithValues("key", req.NamespacedName)
	log.V(1).Info("reconcile meshdeployment")

	deploySyncer := resourcesyncer.NewDeploymentSyncer(r.Client, deploy, r.Scheme)
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
		Complete(r)
}
