package resource

import (
	"github.com/megaease/easemesh-api/v1alpha1"
	"github.com/megaease/easemeshctl/cmd/client/resource/meta"
)

type (
	// ObservabilityOutputServer describes observability output server resource of the EaseMesh
	ObservabilityOutputServer struct {
		meta.MeshResource `yaml:",inline"`
		Spec              *v1alpha1.ObservabilityOutputServer `yaml:"spec" jsonschema:"required"`
	}
)

// ToV1Alpha1 converts a ObservabilityOutputServer resource to v1alpha1.ObservabilityOutputServer
func (r *ObservabilityOutputServer) ToV1Alpha1() (result *v1alpha1.ObservabilityOutputServer) {
	return r.Spec
}

// ToObservabilityOutputServer converts a v1alpha1.ObservabilityOutputServer resource to a ObservabilityOutputServer resource
func ToObservabilityOutputServer(serviceID string, output *v1alpha1.ObservabilityOutputServer) *ObservabilityOutputServer {
	result := &ObservabilityOutputServer{
		Spec: &v1alpha1.ObservabilityOutputServer{},
	}
	result.MeshResource = NewObservabilityOutputServerResource(DefaultAPIVersion, serviceID)
	result.Spec = output
	return result
}
