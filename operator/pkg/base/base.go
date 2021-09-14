package base

import (
	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/record"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type (
	// Runtime carries base rutime for one controller.
	Runtime struct {
		Name             string
		Client           client.Client
		Scheme           *runtime.Scheme
		Recorder         record.EventRecorder
		Log              logr.Logger
		ImageRegistryURL string
		ImagePullPolicy  string
		SidecarImageName string
		// AgentInitializerImageName is the image name of the Agent initializer.
		AgentInitializerImageName string
		// Log4jConfigName is  the name of log4f config name.
		Log4jConfigName string

		ClusterJoinURLs []string
		ClusterName     string
	}
)
