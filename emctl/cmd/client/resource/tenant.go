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

// Columns returns the columns of Tenant.
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
