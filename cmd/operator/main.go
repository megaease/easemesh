/*
Copyright 2021 MegaEase.cn.
*/

package main

import (
	"flag"
	"io/ioutil"
	"os"

	meshv1beta1 "github.com/megaease/easemesh/mesh-operator/pkg/api/v1beta1"
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
	DefaultImageRegistryURL = "docker.io"
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

type ConfigSpec struct {
	ImageRegistryURL     string `yaml:"image-registry-url" jsonschema:"required"`
	ClusterName          string `yaml:"cluster-name" jsonschema:"required"`
	ClusterJoinURL       string `yaml:"cluster-join-urls" jsonschema:"required"`
	MetricsAddr          string `yaml:"metrics-bind-address" jsonschema:"required"`
	EnableLeaderElection bool   `yaml:"leader-elect" jsonschema:"required"`
	ProbeAddr            string `yaml:"health-probe-bind-address" jsonschema:"required"`
}

func main() {
	var imageRegistryURL string
	var clusterName string
	var clusterJoinURL string
	var metricsAddr string
	var enableLeaderElection bool
	var probeAddr string
	var configFile string

	flag.StringVar(&imageRegistryURL, "image-registry-url", DefaultImageRegistryURL, "The Registry URL of the Image.")
	flag.StringVar(&clusterName, "cluster-name", "", "The cluster-name of the eg master.")
	flag.StringVar(&clusterJoinURL, "cluster-join-urls", "", "The address the eg master binds to.")
	flag.StringVar(&metricsAddr, "metrics-bind-address", ":8080", "The address the metric endpoint binds to.")
	flag.StringVar(&probeAddr, "health-probe-bind-address", ":8081", "The address the probe endpoint binds to.")
	flag.BoolVar(&enableLeaderElection, "leader-elect", false, "Enable leader election for controller manager. "+
		"Enabling this will ensure there is only one active controller manager.")
	flag.StringVar(&configFile, "config", " ", "A yaml file config the operator. ")
	opts := zap.Options{
		Development: true,
	}
	opts.BindFlags(flag.CommandLine)
	flag.Parse()

	if configFile != " " {
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
		clusterJoinURL = spec.ClusterJoinURL
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

	if err = (&controllers.MeshDeploymentReconciler{
		Client:           mgr.GetClient(),
		Log:              ctrl.Log.WithName("controllers").WithName("MeshDeployment"),
		Scheme:           mgr.GetScheme(),
		ClusterJoinURL:   clusterJoinURL,
		ClusterName:      clusterName,
		ImageRegistryURL: imageRegistryURL,
		Recorder:         mgr.GetEventRecorderFor("controller.MeshDeployment"),
	}).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "MeshDeployment")
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
