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

		ClusterJoinURLs []string
		ClusterName     string
	}
)
