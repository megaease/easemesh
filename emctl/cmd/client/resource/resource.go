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
	"github.com/megaease/easemeshctl/cmd/client/resource/meta"
)

// TODO: Split every kind of resource to a deidcate package,
// which contains its information as much as possible.
// Currently, adding a new resource is very long-winded to do it
// in many packages across the repository.

const (
	// DefaultAPIVersion is current apis version for the EaseMesh
	DefaultAPIVersion = "mesh.megaease.com/v1alpha1"

	// LoadBalanceRoundRobinPolicy is round robin policy
	LoadBalanceRoundRobinPolicy = "roundRobin"

	// DefaultSideIngressProtocol is default communication protocol for inbound traffic of the sidecar
	DefaultSideIngressProtocol = "http"

	// DefaultSideEgressProtocol is default communication protocol for outbound traffic of the sidecar
	DefaultSideEgressProtocol = "http"

	// DefaultSideIngressPort is default port listend by the sidecar for inbound traffic
	DefaultSideIngressPort = 13001

	// DefaultSideEgressPort is default port listend by the sidecar for outbound traffic
	DefaultSideEgressPort = 13002

	// KindMeshController is mesh controller kind of the EaseMesh control plane.
	KindMeshController = "MeshController"

	// KindService is service kind of the EaseMesh resource
	KindService = "Service"

	// KindServiceInstance is the service instance kind of EaseMesh resource.
	KindServiceInstance = "ServiceInstance"

	// KindCanary is canary kind of the EaseMesh resource
	KindCanary = "Canary"

	// KindObservabilityMetrics is observability metrics kind of the EaseMesh resource
	KindObservabilityMetrics = "ObservabilityMetrics"

	// KindObservabilityTracings is observability tracings kind of the EaseMesh resource
	KindObservabilityTracings = "ObservabilityTracings"

	// KindObservabilityOutputServer is observability output server kind of the EaseMesh resource
	KindObservabilityOutputServer = "ObservabilityOutputServer"

	// KindTenant is tenant kind of the EaseMesh resource
	KindTenant = "Tenant"

	// KindLoadBalance is loadbalance kind of the EaseMesh resource
	KindLoadBalance = "LoadBalance"

	// KindResilience is resilience kind of the EaseMesh resource
	KindResilience = "Resilience"

	// KindIngress is ingress kind of the EaseMesh resource
	KindIngress = "Ingress"

	// KindCustomResourceKind is the kind of 'custom resource kind' of the EaseMesh resource
	KindCustomResourceKind = "CustomResourceKind"
)

type (
	// ObjectCreator create a MeshObject
	ObjectCreator interface {
		NewFromKind(meta.VersionKind) (meta.MeshObject, error)
		NewFromResource(meta.MeshResource) (meta.MeshObject, error)
	}

	objectCreator struct{}
)

// NewObjectCreator creates an ObjectCreator
func NewObjectCreator() ObjectCreator {
	return &objectCreator{}
}

func (oc *objectCreator) NewFromKind(kind meta.VersionKind) (meta.MeshObject, error) {
	return oc.new(kind, meta.MetaData{})
}

func (oc *objectCreator) NewFromResource(resource meta.MeshResource) (meta.MeshObject, error) {
	return oc.new(resource.VersionKind, resource.MetaData)
}

func (oc *objectCreator) new(kind meta.VersionKind, metaData meta.MetaData) (meta.MeshObject, error) {
	apiVersion := kind.APIVersion
	if apiVersion == "" {
		apiVersion = DefaultAPIVersion
	}

	switch kind.Kind {
	case KindMeshController:
		return &MeshController{
			MeshResource: NewMeshControllerResource(apiVersion, metaData.Name),
		}, nil
	case KindService:
		return &Service{
			MeshResource: NewServiceResource(apiVersion, metaData.Name),
		}, nil
	case KindServiceInstance:
		return &ServiceInstance{
			MeshResource: NewServiceInstanceResource(apiVersion, metaData.Name),
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
	case KindCustomResourceKind:
		return &CustomResourceKind{
			MeshResource: NewCustomResourceKindResource(apiVersion, metaData.Name),
		}, nil
	default:
		return &CustomResource{
			MeshResource: NewMeshResource(apiVersion, kind.Kind, metaData.Name),
		}, nil
	}
}

// NewMeshControllerResource returns a MeshResouce with the mesh controller kind.
func NewMeshControllerResource(apiVersion, name string) meta.MeshResource {
	return NewMeshResource(apiVersion, KindMeshController, name)
}

// NewIngressResource returns a MeshResource with the ingress kind
func NewIngressResource(apiVersion, name string) meta.MeshResource {
	return NewMeshResource(apiVersion, KindIngress, name)
}

// NewServiceResource returns a MeshResource with the service kind
func NewServiceResource(apiVersion, name string) meta.MeshResource {
	return NewMeshResource(apiVersion, KindService, name)
}

// NewServiceInstanceResource returns a MeshResource with the service kind
func NewServiceInstanceResource(apiVersion, name string) meta.MeshResource {
	return NewMeshResource(apiVersion, KindServiceInstance, name)
}

// NewCanaryResource returns a MeshResource with the canary kind
func NewCanaryResource(apiVersion, name string) meta.MeshResource {
	return NewMeshResource(apiVersion, KindCanary, name)
}

// NewLoadBalanceResource returns a MeshResource with the loadbalance kind
func NewLoadBalanceResource(apiVersion, name string) meta.MeshResource {
	return NewMeshResource(apiVersion, KindLoadBalance, name)
}

// NewResilienceResource returns a MeshResource with the resilience kind
func NewResilienceResource(apiVersion, name string) meta.MeshResource {
	return NewMeshResource(apiVersion, KindResilience, name)
}

// NewObservabilityTracingsResource returns a MeshResource with the observability tracings kind
func NewObservabilityTracingsResource(apiVersion, name string) meta.MeshResource {
	return NewMeshResource(apiVersion, KindObservabilityTracings, name)
}

// NewObservabilityMetricsResource returns a MeshResource with the observability metrics kind
func NewObservabilityMetricsResource(apiVersion, name string) meta.MeshResource {
	return NewMeshResource(apiVersion, KindObservabilityMetrics, name)
}

// NewObservabilityOutputServerResource returns a MeshResource with the observability output service kind
func NewObservabilityOutputServerResource(apiVersion, name string) meta.MeshResource {
	return NewMeshResource(apiVersion, KindObservabilityOutputServer, name)
}

// NewTenantResource returns a MeshResource with the tenant kind
func NewTenantResource(apiVersion, name string) meta.MeshResource {
	return NewMeshResource(apiVersion, KindTenant, name)
}

// NewCustomResourceKindResource returns a MeshResource with the custom resource kind
func NewCustomResourceKindResource(apiVersion, name string) meta.MeshResource {
	return NewMeshResource(apiVersion, KindCustomResourceKind, name)
}

// NewMeshResource returns a generic MeshResource
func NewMeshResource(api, kind, name string) meta.MeshResource {
	return meta.MeshResource{
		VersionKind: meta.VersionKind{
			APIVersion: api,
			Kind:       kind,
		},
		MetaData: meta.MetaData{
			Name: name,
		},
	}
}
