/*
 * Copyright (c) 2021, MegaEase
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
	"github.com/megaease/easemesh-api/v2alpha1"
	"github.com/megaease/easemeshctl/cmd/client/resource/meta"
)

type (
	// Service describes service resource of the EaseMesh
	Service struct {
		meta.MeshResource `yaml:",inline"`
		Spec              *ServiceSpec `yaml:"spec" jsonschema:"required"`
	}

	// ServiceSpec describes details of the service resource
	ServiceSpec struct {
		RegisterTenant string `yaml:"registerTenant" jsonschema:"required"`

		Sidecar       *v2alpha1.Sidecar       `yaml:"sidecar" jsonschema:"required"`
		Mock          *v2alpha1.Mock          `yaml:"mock" jsonschema:"omitempty"`
		Resilience    *v2alpha1.Resilience    `yaml:"resilience" jsonschema:"omitempty"`
		LoadBalance   *v2alpha1.LoadBalance   `yaml:"loadBalance" jsonschema:"omitempty"`
		Observability *v2alpha1.Observability `yaml:"observability" jsonschema:"omitempty"`
	}
)

var _ meta.TableObject = &Service{}

// Columns returns the columns of Service.
func (s *Service) Columns() []*meta.TableColumn {
	if s.Spec == nil {
		return nil
	}

	return []*meta.TableColumn{
		{
			Name:  "Tenant",
			Value: s.Spec.RegisterTenant,
		},
	}
}

// ToV2Alpha1 converts an Ingress resource to v2alpha1.Ingress
func (s *Service) ToV2Alpha1() *v2alpha1.Service {
	result := &v2alpha1.Service{}
	result.Name = s.Name()
	if s.Spec != nil {
		result.RegisterTenant = s.Spec.RegisterTenant
		result.Resilience = s.Spec.Resilience
		result.LoadBalance = s.Spec.LoadBalance
		result.Mock = s.Spec.Mock
		result.Sidecar = s.Spec.Sidecar
		result.Observability = s.Spec.Observability
	}
	return result
}

// ToService converts a v2alpha1.Service resource to a Service resource
func ToService(service *v2alpha1.Service) *Service {
	result := &Service{
		Spec: &ServiceSpec{},
	}
	result.MeshResource = NewServiceResource(DefaultAPIVersion, service.Name)
	result.Spec.RegisterTenant = service.RegisterTenant
	result.Spec.Sidecar = service.Sidecar
	result.Spec.Resilience = service.Resilience
	result.Spec.Mock = service.Mock
	result.Spec.LoadBalance = service.LoadBalance
	result.Spec.Observability = service.Observability
	return result
}
