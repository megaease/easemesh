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
	"github.com/megaease/easemesh-api/v1alpha1"
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

		Sidecar       *v1alpha1.Sidecar       `yaml:"sidecar" jsonschema:"required"`
		Resilience    *v1alpha1.Resilience    `yaml:"resilience" jsonschema:"omitempty"`
		Canary        *v1alpha1.Canary        `yaml:"canary" jsonschema:"omitempty"`
		LoadBalance   *v1alpha1.LoadBalance   `yaml:"loadBalance" jsonschema:"omitempty"`
		Observability *v1alpha1.Observability `yaml:"observability" jsonschema:"omitempty"`
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

// ToV1Alpha1 converts an Ingress resource to v1alpha1.Ingress
func (s *Service) ToV1Alpha1() *v1alpha1.Service {
	result := &v1alpha1.Service{}
	result.Name = s.Name()
	if s.Spec != nil {
		result.RegisterTenant = s.Spec.RegisterTenant
		result.Resilience = s.Spec.Resilience
		result.Canary = s.Spec.Canary
		result.LoadBalance = s.Spec.LoadBalance
		result.Sidecar = s.Spec.Sidecar
		result.Observability = s.Spec.Observability
	}
	return result
}

// ToService converts a v1alpha1.Service resource to a Service resource
func ToService(service *v1alpha1.Service) *Service {
	result := &Service{
		Spec: &ServiceSpec{},
	}
	result.MeshResource = NewServiceResource(DefaultAPIVersion, service.Name)
	result.Spec.RegisterTenant = service.RegisterTenant
	result.Spec.Sidecar = service.Sidecar
	result.Spec.Resilience = service.Resilience
	result.Spec.Canary = service.Canary
	result.Spec.LoadBalance = service.LoadBalance
	result.Spec.Observability = service.Observability
	return result
}
