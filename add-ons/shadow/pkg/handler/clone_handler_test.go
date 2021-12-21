package handler

import (
	"testing"

	shadowfake "github.com/megaease/easemesh/mesh-shadow/pkg/handler/fake"
)

func TestShadowServiceCloner_Clone(t *testing.T) {

	cloner := &ShadowServiceCloner{
		KubeClient: prepareClientForTest(),
	}

	shadowService := shadowfake.NewShadowService()
	sourceDeployment := shadowfake.NewSourceDeployment()

	serviceCloneBlock := ShadowServiceBlock{
		service:   shadowService,
		deployObj: sourceDeployment,
	}
	cloner.Clone(serviceCloneBlock)
}
