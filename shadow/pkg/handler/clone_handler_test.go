package handler

import (
	"testing"
)

func TestShadowServiceCloner_Clone(t *testing.T) {
	shadowService := fakeShadowService()
	deployment := fakeDeployment()

	serviceCloneBlock := ServiceCloneBlock{
		service:   shadowService,
		deployObj: deployment,
	}

	cloner := &ShadowServiceCloner{
		KubeClient:    nil,
		RunTimeClient: nil,
	}
	cloner.Clone(serviceCloneBlock)
}
