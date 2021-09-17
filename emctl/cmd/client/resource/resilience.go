package resource

import (
	"github.com/megaease/easemesh-api/v1alpha1"
	"github.com/megaease/easemeshctl/cmd/client/resource/meta"
)

type (
	// Resilience describes resilience resource of the EaseMesh
	Resilience struct {
		meta.MeshResource `yaml:",inline"`
		Spec              *v1alpha1.Resilience `yaml:"spec" jsonschema:"required"`
	}
)

// ToV1Alpha1 converts a Resilience resource to v1alpha1.Resilience
func (r *Resilience) ToV1Alpha1() *v1alpha1.Resilience {
	return r.Spec
}

// ToResilience converts a v1alpha1.Resilience resource to a Resilience resource
func ToResilience(name string, resilience *v1alpha1.Resilience) *Resilience {
	result := &Resilience{
		Spec: &v1alpha1.Resilience{},
	}
	result.MeshResource = NewResilienceResource(DefaultAPIVersion, name)
	result.Spec.RateLimiter = resilience.RateLimiter
	result.Spec.Retryer = resilience.Retryer
	result.Spec.CircuitBreaker = resilience.CircuitBreaker
	result.Spec.TimeLimiter = resilience.TimeLimiter
	return result
}
