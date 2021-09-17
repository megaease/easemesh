package resource

import (
	"strings"

	"github.com/megaease/easemesh-api/v1alpha1"
	"github.com/megaease/easemeshctl/cmd/client/command/printer"
	"github.com/megaease/easemeshctl/cmd/client/resource/meta"
)

type (
	// Tenant describes tenant resource of the EaseMesh
	Tenant struct {
		meta.MeshResource `yaml:",inline"`
		Spec              *TenantSpec `yaml:"spec" jsonschema:"required"`
	}

	// TenantSpec describes whats service resided in
	TenantSpec struct {
		Services    []string `yaml:"services" jsonschema:"omitempty"`
		Description string   `yaml:"description" jsonschema:"omitempty"`
	}
)

var _ printer.TableObject = &Service{}

func (t *Tenant) Columns() []*printer.TableColumn {
	if t.Spec == nil {
		return nil
	}

	return []*printer.TableColumn{
		{
			Name:  "Services",
			Value: strings.Join(t.Spec.Services, ","),
		},
		{
			Name:  "Description",
			Value: t.Spec.Description,
		},
	}
}

// ToV1Alpha1 converts an Ingress resource to v1alpha1.Ingress
func (t *Tenant) ToV1Alpha1() *v1alpha1.Tenant {
	result := &v1alpha1.Tenant{}
	result.Name = t.Name()
	if t.Spec != nil {
		result.Services = t.Spec.Services
		result.Description = t.Spec.Description
	}
	return result
}

// ToTenant converts a v1alpha1.Tenant resource to a Tenant resource
func ToTenant(tenant *v1alpha1.Tenant) *Tenant {
	result := &Tenant{
		Spec: &TenantSpec{},
	}
	result.MeshResource = NewTenantResource(DefaultAPIVersion, tenant.Name)
	result.Spec.Services = tenant.Services
	result.Spec.Description = tenant.Description
	return result
}
