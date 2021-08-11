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
	"io/ioutil"
	"os"

	"github.com/spf13/pflag"

	meshv1beta1 "github.com/megaease/easemesh/mesh-operator/pkg/api/v1beta1"
	"github.com/megaease/easemesh/mesh-operator/pkg/base"
	"github.com/megaease/easemesh/mesh-operator/pkg/controllers"
	"github.com/megaease/easemesh/mesh-operator/pkg/hook"

	// Import all Kubernetes client auth plugins (e.g. Azure, GCP, OIDC, etc.)
	// to ensure that exec-entrypoint and run can make use of them.

	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	_ "k8s.io/client-go/plugin/pkg/client/auth"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/healthz"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
	"sigs.k8s.io/controller-runtime/pkg/webhook"

	// +kubebuilder:scaffold:imports
	"gopkg.in/yaml.v2"
)

const (
	// DefaultImageRegistryURL is the default image registry URL.
	DefaultImageRegistryURL = "docker.io"

	// DefaultSidecarImageName is the default sidecar image name.
	DefaultSidecarImageName = "megaease/easegress:server-sidecar"

	// DefaultImagePullPolicy is the default image pull policy.
	DefaultImagePullPolicy = "IfNotPresent"
)

var (
	scheme = runtime.NewScheme()
)

func init() {
	utilruntime.Must(clientgoscheme.AddToScheme(scheme))

	utilruntime.Must(meshv1beta1.AddToScheme(scheme))
	// +kubebuilder:scaffold:scheme
}

// ConfigSpec is the config specification.
type ConfigSpec struct {
	ImageRegistryURL     string   `yaml:"image-registry-url" jsonschema:"required"`
	ClusterName          string   `yaml:"cluster-name" jsonschema:"required"`
	ClusterJoinURLs      []string `yaml:"cluster-join-urls" jsonschema:"required"`
	MetricsAddr          string   `yaml:"metrics-bind-address" jsonschema:"required"`
	EnableLeaderElection bool     `yaml:"leader-elect" jsonschema:"required"`
	ProbeAddr            string   `yaml:"health-probe-bind-address" jsonschema:"required"`
	WebhookPort          uint16   `yaml:"webhook-port" jsonschema:"required"`
	CertDir              string   `yaml:"cert-dir" jsonschema:"required"`
	CertName             string   `yaml:"cert-name" jsonschema:"required"`
	KeyName              string   `yaml:"key-name" jsonschema:"required"`
}

func main() {
	// TODO: Make flags/specfile parsing more maintainable.

	var (
		imageRegistryURL     string
		sidecarImageName     string
		imagePullPolicy      string
		clusterName          string
		clusterJoinURLs      []string
		metricsAddr          string
		enableLeaderElection bool
		configFile           string
		probeAddr            string
		webhookPort          uint16
		certDir              string
		certName             string
		keyName              string
	)

	pflag.StringVar(&imageRegistryURL, "image-registry-url", DefaultImageRegistryURL, "The image registry URL")
	pflag.StringVar(&sidecarImageName, "sidecar-image-name", DefaultSidecarImageName, "The sidecar image name.")
	pflag.StringVar(&imagePullPolicy, "image-pull-policy", DefaultImagePullPolicy, "The image pull policy. (support Always, IfNotPresent, Never)")
	pflag.StringVar(&clusterName, "cluster-name", "", "The name of the Easegress cluster.")
	pflag.StringSliceVar(&clusterJoinURLs, "cluster-join-urls", []string{"http://easemesh-controlplane-svc.easemesh:2380"}, "The addresses to join the Easegress.")
	pflag.StringVar(&metricsAddr, "metrics-bind-address", ":8080", "The address the metric endpoint binds to.")
	pflag.StringVar(&probeAddr, "health-probe-bind-address", ":8081", "The address the probe endpoint binds to.")
	pflag.BoolVar(&enableLeaderElection, "leader-elect", false, "Enable leader election for controller manager. "+
		"Enabling this will ensure there is only one active controller manager.")
	pflag.StringVar(&configFile, "config", " ", "A yaml file config the operator. ")
	pflag.StringVar(&certDir, "cert-dir", "/cert-volume", "The TLS cert directory.")
	pflag.StringVar(&certName, "cert-file", "cert.pem", "The TLS cert file name.")
	pflag.StringVar(&keyName, "key-file", "key.pem", "The TLS key file name.")
	pflag.Uint16Var(&webhookPort, "webhook-port", 9090, "Webhook port listening on.")

	pflag.Parse()

	ctrl.SetLogger(zap.New(zap.UseDevMode(false)))
	setupLog := ctrl.Log.WithName("setup")

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
		enableLeaderElection = spec.EnableLeaderElection
		probeAddr = spec.ProbeAddr
		webhookPort = spec.WebhookPort
		certDir = spec.CertDir
		certName = spec.CertName
		keyName = spec.KeyName
	}

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
		ImagePullPolicy:  imagePullPolicy,
		SidecarImageName: sidecarImageName,

		ClusterJoinURLs: clusterJoinURLs,
		ClusterName:     clusterName,
	}

	// Create MeshDeploymentReconciler.
	meshDeploymentRuntime := baseRuntime
	meshDeploymentRuntime.Name = "MeshDeployment"
	meshDeploymentRuntime.Log = ctrl.Log.WithName("controllers").WithName("MeshDeployment")
	meshDeploymentReconciler := &controllers.MeshDeploymentReconciler{Runtime: &meshDeploymentRuntime}
	meshDeploymentReconciler.SetupWithManager(mgr)
	if err != nil {
		setupLog.Error(err, "create controller of MeshDeployment failed")
		os.Exit(1)
	}

	// Create a webhook server.
	webhookRuntime := baseRuntime
	webhookRuntime.Name = "Webhook"
	webhookRuntime.Log = ctrl.Log.WithName("webhook").WithName("mutate")
	webhookMutate := hook.NewMutateHook(&webhookRuntime)
	webhookServer := &webhook.Server{
		Port:     int(webhookPort),
		CertDir:  certDir,
		CertName: certName,
		KeyName:  keyName,
	}

	webhookServer.Register("/mutate", webhookMutate.Admission)

	if err := mgr.Add(webhookServer); err != nil {
		setupLog.Error(err, "unable to set up webhook server")
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
