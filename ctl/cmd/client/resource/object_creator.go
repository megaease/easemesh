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
	ObjectCreator interface {
		NewFromKind(VersionKind) (MeshObject, error)
		NewFromResource(MeshResource) (MeshObject, error)
	}

	objectCreator struct{}
)

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

func NewIngressResource(apiVersion, name string) MeshResource {
	return NewMeshResource(apiVersion, KindIngress, name)
}

func NewServiceResource(apiVersion, name string) MeshResource {
	return NewMeshResource(apiVersion, KindService, name)
}

func NewCanaryResource(apiVersion, name string) MeshResource {
	return NewMeshResource(apiVersion, KindCanary, name)
}

func NewLoadBalanceResource(apiVersion, name string) MeshResource {
	return NewMeshResource(apiVersion, KindLoadBalance, name)
}

func NewResilienceResource(apiVersion, name string) MeshResource {
	return NewMeshResource(apiVersion, KindResilience, name)
}

func NewObservabilityTracingsResource(apiVersion, name string) MeshResource {
	return NewMeshResource(apiVersion, KindObservabilityTracings, name)
}

func NewObservabilityMetricsResource(apiVersion, name string) MeshResource {
	return NewMeshResource(apiVersion, KindObservabilityMetrics, name)
}

func NewObservabilityOutputServerResource(apiVersion, name string) MeshResource {
	return NewMeshResource(apiVersion, KindObservabilityOutputServer, name)
}
func NewTenantResource(apiVersion, name string) MeshResource {
	return NewMeshResource(apiVersion, KindTenant, name)
}

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
