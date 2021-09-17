package resource

import (
	"github.com/megaease/easemesh-api/v1alpha1"
	"github.com/megaease/easemeshctl/cmd/client/command/printer"
	"github.com/megaease/easemeshctl/cmd/client/resource/meta"
)

type (
	// LoadBalance describes loadbalance resource of the EaseMesh
	LoadBalance struct {
		meta.MeshResource `yaml:",inline"`
		Spec              *v1alpha1.LoadBalance `yaml:"spec" jsonschema:"required"`
	}
)

var _ printer.TableObject = &LoadBalance{}

func (l *LoadBalance) Columns() []*printer.TableColumn {
	if l.Spec == nil {
		return nil
	}

	return []*printer.TableColumn{
		{
			Name:  "Policy",
			Value: l.Spec.Policy,
		},
		{
			Name:  "HeaderHashKey",
			Value: l.Spec.HeaderHashKey,
		},
	}
}

// ToV1Alpha1 converts a loadbalance resource to v1alpha1.LoadBalance
func (l *LoadBalance) ToV1Alpha1() *v1alpha1.LoadBalance {
	return l.Spec
}

// ToLoadBalance converts a v1alpha1.LoadBalance resource to a LoadBalance resource
func ToLoadBalance(name string, loadBalance *v1alpha1.LoadBalance) *LoadBalance {
	result := &LoadBalance{
		Spec: &v1alpha1.LoadBalance{},
	}
	result.MeshResource = NewLoadBalanceResource(DefaultAPIVersion, name)
	result.Spec = loadBalance
	return result
}
