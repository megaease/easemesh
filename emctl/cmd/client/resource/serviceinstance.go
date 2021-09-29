/*
 * Copyright (c) 2017, MegaEase
 * All rights reserved.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package resource

import (
	"fmt"
	"strings"

	"github.com/megaease/easemesh-api/v1alpha1"
	"github.com/megaease/easemeshctl/cmd/client/resource/meta"
)

type (
	// ServiceInstance describes service instance resource of the EaseMesh
	ServiceInstance struct {
		meta.MeshResource `yaml:",inline"`
		Spec              *v1alpha1.ServiceInstance `yaml:"spec" jsonschema:"required"`
	}
)

var _ meta.TableObject = &ServiceInstance{}

// ParseName parses the name of service instance to service name and instance id.
func (si *ServiceInstance) ParseName() (serviceName, instanceID string, err error) {
	ss := strings.Split(si.Name(), "/")

	if len(ss) != 2 {
		return "", "", fmt.Errorf("invalid service instance name (format: serviceName/instanceID)")
	}

	return ss[0], ss[1], nil
}

// Columns returns the columns of ServiceInstance.
func (si *ServiceInstance) Columns() []*meta.TableColumn {
	if si.Spec == nil {
		return nil
	}

	return []*meta.TableColumn{
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

// ToServiceInstance converts a v1alpha1.ServiceInstance resource to a ServiceInstance resource.
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
