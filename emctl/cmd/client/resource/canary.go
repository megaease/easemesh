package resource

import (
	"github.com/megaease/easemesh-api/v1alpha1"
	"github.com/megaease/easemeshctl/cmd/client/resource/meta"
)

type (
	// Canary describes canary resource of the EaseMesh
	Canary struct {
		meta.MeshResource `yaml:",inline"`
		Spec              *v1alpha1.Canary `yaml:"spec" jsonschema:"required"`
	}
)

// ToV1Alpha1 converts a Canary resource to v1alpha1.Canary
func (c *Canary) ToV1Alpha1() *v1alpha1.Canary {
	return c.Spec
}

// ToCanary converts a v1alpha1.Canary resource to a Canary resource
func ToCanary(name string, canary *v1alpha1.Canary) *Canary {
	result := &Canary{
		Spec: &v1alpha1.Canary{},
	}
	result.MeshResource = NewCanaryResource(DefaultAPIVersion, name)
	result.Spec.CanaryRules = canary.CanaryRules
	return result
}
