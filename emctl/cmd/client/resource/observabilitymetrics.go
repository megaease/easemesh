package resource

import (
	"github.com/megaease/easemesh-api/v1alpha1"
	"github.com/megaease/easemeshctl/cmd/client/resource/meta"
)

type (

	// ObservabilityMetrics describes observability metrics resource of the EaseMesh
	ObservabilityMetrics struct {
		meta.MeshResource `yaml:",inline"`
		Spec              *v1alpha1.ObservabilityMetrics `yaml:"spec" jsonschema:"required"`
	}
)

// ToV1Alpha1 converts a ObservabilityMetrics resource to v1alpha1.ObservabilityMetrics
func (r *ObservabilityMetrics) ToV1Alpha1() (result *v1alpha1.ObservabilityMetrics) {
	return r.Spec
}

// ToObservabilityMetrics converts a v1alpha1.ObservabilityMetrics resource to a ObservabilityMetrics resource
func ToObservabilityMetrics(serviceID string, metrics *v1alpha1.ObservabilityMetrics) *ObservabilityMetrics {
	result := &ObservabilityMetrics{
		Spec: &v1alpha1.ObservabilityMetrics{},
	}
	result.MeshResource = NewObservabilityMetricsResource(DefaultAPIVersion, serviceID)
	result.Spec = metrics
	return result
}
