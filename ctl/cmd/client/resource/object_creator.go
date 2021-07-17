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
	"github.com/pkg/errors"
)

type (
	// ObjectCreator create a MeshObject
	ObjectCreator interface {
		NewFromKind(VersionKind) (MeshObject, error)
		NewFromResource(MeshResource) (MeshObject, error)
	}

	objectCreator struct{}
)

// NewObjectCreator create a ObjectCreator
func NewObjectCreator() ObjectCreator {
	return &objectCreator{}
}

func (oc *objectCreator) NewFromKind(kind VersionKind) (MeshObject, error) {
	return oc.new(kind, MetaData{})
}

func (oc *objectCreator) NewFromResource(resource MeshResource) (MeshObject, error) {
	return oc.new(resource.VersionKind, resource.MetaData)
}

func (oc *objectCreator) new(kind VersionKind, metaData MetaData) (MeshObject, error) {
	apiVersion := kind.APIVersion
	if apiVersion == "" {
		apiVersion = DefaultAPIVersion
	}

	switch kind.Kind {
	case KindService:
		return &Service{
			MeshResource: NewServiceResource(apiVersion, metaData.Name),
		}, nil
	case KindTenant:
		return &Tenant{
			MeshResource: NewTenantResource(apiVersion, metaData.Name),
		}, nil
	case KindLoadBalance:
		return &LoadBalance{
			MeshResource: NewLoadBalanceResource(apiVersion, metaData.Name),
		}, nil
	case KindCanary:
		return &Canary{
			MeshResource: NewCanaryResource(apiVersion, metaData.Name),
		}, nil
	case KindObservabilityTracings:
		return &ObservabilityTracings{
			MeshResource: NewObservabilityTracingsResource(apiVersion, metaData.Name),
		}, nil
	case KindObservabilityOutputServer:
		return &ObservabilityOutputServer{
			MeshResource: NewObservabilityOutputServerResource(apiVersion, metaData.Name),
		}, nil
	case KindObservabilityMetrics:
		return &ObservabilityMetrics{
			MeshResource: NewObservabilityMetricsResource(apiVersion, metaData.Name),
		}, nil
	case KindResilience:
		return &Resilience{
			MeshResource: NewResilienceResource(apiVersion, metaData.Name),
		}, nil
	case KindIngress:
		return &Ingress{
			MeshResource: NewIngressResource(apiVersion, metaData.Name),
		}, nil
	default:
		return nil, errors.Errorf("unsupported kind %s", kind.Kind)
	}
}

// NewIngressResource return a MeshResource with the ingress kind
func NewIngressResource(apiVersion, name string) MeshResource {
	return NewMeshResource(apiVersion, KindIngress, name)
}

// NewServiceResource return a MeshResource with the service kind
func NewServiceResource(apiVersion, name string) MeshResource {
	return NewMeshResource(apiVersion, KindService, name)
}

// NewCanaryResource return a MeshResource with the canary kind
func NewCanaryResource(apiVersion, name string) MeshResource {
	return NewMeshResource(apiVersion, KindCanary, name)
}

// NewLoadBalanceResource return a MeshResource with the loadbalance kind
func NewLoadBalanceResource(apiVersion, name string) MeshResource {
	return NewMeshResource(apiVersion, KindLoadBalance, name)
}

// NewResilienceResource return a MeshResource with the resilience kind
func NewResilienceResource(apiVersion, name string) MeshResource {
	return NewMeshResource(apiVersion, KindResilience, name)
}

// NewObservabilityTracingsResource return a MeshResource with the observability tracings kind
func NewObservabilityTracingsResource(apiVersion, name string) MeshResource {
	return NewMeshResource(apiVersion, KindObservabilityTracings, name)
}

// NewObservabilityMetricsResource return a MeshResource with the observability metrics kind
func NewObservabilityMetricsResource(apiVersion, name string) MeshResource {
	return NewMeshResource(apiVersion, KindObservabilityMetrics, name)
}

// NewObservabilityOutputServerResource return a MeshResource with the observability output service kind
func NewObservabilityOutputServerResource(apiVersion, name string) MeshResource {
	return NewMeshResource(apiVersion, KindObservabilityOutputServer, name)
}

// NewTenantResource return a MeshResource with the tenant kind
func NewTenantResource(apiVersion, name string) MeshResource {
	return NewMeshResource(apiVersion, KindTenant, name)
}

// NewMeshResource return a generic MeshResource
func NewMeshResource(api, kind, name string) MeshResource {
	return MeshResource{
		VersionKind: VersionKind{
			APIVersion: api,
			Kind:       kind,
		},
		MetaData: MetaData{
			Name: name,
		},
	}
}
