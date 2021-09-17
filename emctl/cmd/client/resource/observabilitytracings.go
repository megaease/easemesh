package resource

import (
	"github.com/megaease/easemesh-api/v1alpha1"
	"github.com/megaease/easemeshctl/cmd/client/resource/meta"
)

type (
	// ObservabilityTracings describes observability tracings resource of the EaseMesh
	ObservabilityTracings struct {
		meta.MeshResource `yaml:",inline"`
		Spec              *v1alpha1.ObservabilityTracings `yaml:"spec" jsonschema:"required"`
	}
)

// ToV1Alpha1 converts a ObservabilityTracings resource to v1alpha1.ObservabilityTracings
func (r *ObservabilityTracings) ToV1Alpha1() (result *v1alpha1.ObservabilityTracings) {
	return r.Spec
}

// ToObservabilityTracings converts a v1alpha1.ObservabilityTracings resource to a ObservabilityTracings resource
func ToObservabilityTracings(serviceID string, tracing *v1alpha1.ObservabilityTracings) *ObservabilityTracings {
	result := &ObservabilityTracings{
		Spec: &v1alpha1.ObservabilityTracings{},
	}
	result.MeshResource = NewObservabilityTracingsResource(DefaultAPIVersion, serviceID)
	result.Spec = tracing
	return result
}
