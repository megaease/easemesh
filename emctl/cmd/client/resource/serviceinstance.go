package resource

import (
	"fmt"
	"strings"

	"github.com/megaease/easemesh-api/v1alpha1"
	"github.com/megaease/easemeshctl/cmd/client/command/printer"
	"github.com/megaease/easemeshctl/cmd/client/resource/meta"
)

type (
	ServiceInstance struct {
		meta.MeshResource `yaml:",inline"`
		Spec              *v1alpha1.ServiceInstance `yaml:"spec" jsonschema:"required"`
	}
)

var _ printer.TableObject = &ServiceInstance{}

func (si *ServiceInstance) ParseName() (serviceName, instanceID string, err error) {
	ss := strings.Split(si.Name(), "/")

	if len(ss) != 2 {
		return "", "", fmt.Errorf("invalid service instance name (format: serviceName/instanceID)")
	}

	return ss[0], ss[1], nil
}

// Columns returns self-defined columns.
func (si *ServiceInstance) Columns() []*printer.TableColumn {
	if si.Spec == nil {
		return nil
	}

	return []*printer.TableColumn{
		{
			Name:  "RegistryName",
			Value: si.Spec.RegistryName,
		},
		{
			Name:  "IP",
			Value: si.Spec.Ip,
		},
		{
			Name:  "Port",
			Value: fmt.Sprintf("%d", si.Spec.Port),
		},
		{
			Name:  "Status",
			Value: si.Spec.Status,
		},
	}
}

func ToServiceInstance(instance *v1alpha1.ServiceInstance) *ServiceInstance {
	result := &ServiceInstance{
		Spec: &v1alpha1.ServiceInstance{},
	}

	name := fmt.Sprintf("%s/%s", instance.ServiceName, instance.InstanceID)

	result.MeshResource = NewServiceInstanceResource(DefaultAPIVersion, name)
	result.MeshResource.MetaData.Labels = instance.Labels
	result.Spec = instance

	return result
}

// ToV1Alpha1 converts a Service resource to v1alpha1.ServiceInstance.
func (si *ServiceInstance) ToV1Alpha1() *v1alpha1.ServiceInstance {
	return si.Spec
}
