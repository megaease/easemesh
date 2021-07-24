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

package main

import (
	"flag"
	"io/ioutil"
	"os"

	"github.com/spf13/pflag"

	meshv1beta1 "github.com/megaease/easemesh/mesh-operator/pkg/api/v1beta1"
	"github.com/megaease/easemesh/mesh-operator/pkg/base"
	"github.com/megaease/easemesh/mesh-operator/pkg/controllers"

	// Import all Kubernetes client auth plugins (e.g. Azure, GCP, OIDC, etc.)
	// to ensure that exec-entrypoint and run can make use of them.

	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	_ "k8s.io/client-go/plugin/pkg/client/auth"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/healthz"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"

	// +kubebuilder:scaffold:imports
	"gopkg.in/yaml.v2"
)

const (
	// DefaultImageRegistryURL is the default image registry URL.
	DefaultImageRegistryURL = "docker.io"

	// DefaultSidecarImageName is the default sidecar image name.
	DefaultSidecarImageName = "megaease/easegress:server-sidecar"
)

var (
	scheme   = runtime.NewScheme()
	setupLog = ctrl.Log.WithName("setup")
)

func init() {
	utilruntime.Must(clientgoscheme.AddToScheme(scheme))

	utilruntime.Must(meshv1beta1.AddToScheme(scheme))
	// +kubebuilder:scaffold:scheme
}

// ConfigSpec is the config specification.
type ConfigSpec struct {
	ImageRegistryURL string `yaml:"image-registry-url" jsonschema:"required"`
	ClusterName      string `yaml:"cluster-name" jsonschema:"required"`
	// TODO: Make it to []string along with install configmap,
	// so it only supports one url for now.
	ClusterJoinURLs      string `yaml:"cluster-join-urls" jsonschema:"required"`
	MetricsAddr          string `yaml:"metrics-bind-address" jsonschema:"required"`
	EnableLeaderElection bool   `yaml:"leader-elect" jsonschema:"required"`
	ProbeAddr            string `yaml:"health-probe-bind-address" jsonschema:"required"`
}

func main() {
	var imageRegistryURL string
	var sidecarImageName string
	var clusterName string
	var clusterJoinURLs string
	var metricsAddr string
	var enableLeaderElection bool
	var probeAddr string
	var configFile string

	pflag.StringVar(&imageRegistryURL, "image-registry-url", DefaultImageRegistryURL, "The image registry URL")
	pflag.StringVar(&sidecarImageName, "sidecar-image-name", DefaultSidecarImageName, "The sidecar image name.")
	pflag.StringVar(&clusterName, "cluster-name", "", "The name of the Easegress cluster.")
	pflag.StringVar(&clusterJoinURLs, "cluster-join-urls", "", "The addresses to join the Easegress.")
	pflag.StringVar(&metricsAddr, "metrics-bind-address", ":8080", "The address the metric endpoint binds to.")
	pflag.StringVar(&probeAddr, "health-probe-bind-address", ":8081", "The address the probe endpoint binds to.")
	pflag.BoolVar(&enableLeaderElection, "leader-elect", false, "Enable leader election for controller manager. "+
		"Enabling this will ensure there is only one active controller manager.")
	pflag.StringVar(&configFile, "config", " ", "A yaml file config the operator. ")

	pflag.Parse()

	opts := zap.Options{
		Development: true,
	}

	opts.BindFlags(flag.CommandLine)
	pflag.Parse()

	if configFile != "" {
		config, err := ioutil.ReadFile(configFile)
		if err != nil {
			setupLog.Error(err, "Read configFile error, %v", err)
			os.Exit(1)
		}
		spec := &ConfigSpec{}
		err = yaml.Unmarshal(config, spec)
		if err != nil {
			setupLog.Error(err, "Read configFile error, %v", err)
			os.Exit(1)
		}

		imageRegistryURL = spec.ImageRegistryURL
		clusterName = spec.ClusterName
		clusterJoinURLs = spec.ClusterJoinURLs
		metricsAddr = spec.MetricsAddr
		probeAddr = spec.ProbeAddr
		enableLeaderElection = spec.EnableLeaderElection
	}

	ctrl.SetLogger(zap.New(zap.UseFlagOptions(&opts)))

	mgr, err := ctrl.NewManager(ctrl.GetConfigOrDie(), ctrl.Options{
		Scheme:                 scheme,
		MetricsBindAddress:     metricsAddr,
		Port:                   9443,
		HealthProbeBindAddress: probeAddr,
		LeaderElection:         enableLeaderElection,
		LeaderElectionID:       "870093a3.megaease.com",
	})
	if err != nil {
		setupLog.Error(err, "unable to start manager")
		os.Exit(1)
	}

	baseRuntime := base.Runtime{
		Client:           mgr.GetClient(),
		Scheme:           mgr.GetScheme(),
		Recorder:         mgr.GetEventRecorderFor("controller.MeshDeployment"),
		ImageRegistryURL: imageRegistryURL,
		SidecarImageName: sidecarImageName,

		ClusterJoinURLs: []string{clusterJoinURLs},
		ClusterName:     clusterName,
	}

	meshDeploymentRuntime := baseRuntime
	meshDeploymentRuntime.Name = "MeshDeployment"
	meshDeploymentRuntime.Log = ctrl.Log.WithName("controllers").WithName("MeshDeployment")
	meshDeploymentReconciler := &controllers.MeshDeploymentReconciler{Runtime: &meshDeploymentRuntime}
	meshDeploymentReconciler.SetupWithManager(mgr)
	if err != nil {
		setupLog.Error(err, "create controller of MeshDeployment failed")
		os.Exit(1)
	}

	deploymentRuntime := baseRuntime
	deploymentRuntime.Name = "Deployment"
	deploymentRuntime.Log = ctrl.Log.WithName("controllers").WithName("Deployment")
	deploymentReconciler := &controllers.DeploymentReconciler{Runtime: &meshDeploymentRuntime}
	deploymentReconciler.SetupWithManager(mgr)
	if err != nil {
		setupLog.Error(err, "create controller of Deployment failed")
		os.Exit(1)
	}

	// +kubebuilder:scaffold:builder

	if err := mgr.AddHealthzCheck("health", healthz.Ping); err != nil {
		setupLog.Error(err, "unable to set up health check")
		os.Exit(1)
	}
	if err := mgr.AddReadyzCheck("check", healthz.Ping); err != nil {
		setupLog.Error(err, "unable to set up ready check")
		os.Exit(1)
	}

	setupLog.Info("starting manager")
	if err := mgr.Start(ctrl.SetupSignalHandler()); err != nil {
		setupLog.Error(err, "problem running manager")
		os.Exit(1)
	}
}
