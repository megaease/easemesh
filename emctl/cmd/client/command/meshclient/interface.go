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

package meshclient

import (
	"context"

	"github.com/megaease/easemeshctl/cmd/client/resource"
)

// MeshClient is a client for accessing the EaseMesh control plane service
type MeshClient interface {
	V1Alpha1() V1Alpha1Interface
}

// V1Alpha1Interface is an interface that aggregates all resources accessor for the EaseMesh
type V1Alpha1Interface interface {
	MeshControllerGetter
	TenantGetter
	ServiceGetter
	ServiceInstanceGetter
	LoadbalanceGetter
	CanaryGetter
	ObservabilityGetter
	ResilienceGetter
	IngressGetter
	CustomResourceKindGetter
	CustomResourceGetter
}

// MeshControllerGetter represents a mesh controller resource accessor
type MeshControllerGetter interface {
	MeshController() MeshControllerInterface
}

// ServiceGetter represents a Service resource accessor
type ServiceGetter interface {
	Service() ServiceInterface
}

// ServiceInstanceGetter represents a Service resource accessor
type ServiceInstanceGetter interface {
	ServiceInstance() ServiceInstanceInterface
}

// ResilienceGetter represents a Resilience resource accessor
type ResilienceGetter interface {
	Resilience() ResilienceInterface
}

// IngressGetter represents an Ingress resource accessor
type IngressGetter interface {
	Ingress() IngressInterface
}

// CustomResourceKindGetter represents an CustomResourceKind accessor
type CustomResourceKindGetter interface {
	CustomResourceKind() CustomResourceKindInterface
}

// CustomResourceGetter represents an CustomResource accessor
type CustomResourceGetter interface {
	CustomResource() CustomResourceInterface
}

// MeshControllerInterface captures the set of operations for interacting with the EaseMesh REST apis of the mesh controller resource.
type MeshControllerInterface interface {
	Get(context.Context, string) (*resource.MeshController, error)
	Patch(context.Context, *resource.MeshController) error
	Create(context.Context, *resource.MeshController) error
	Delete(context.Context, string) error
	List(context.Context) ([]*resource.MeshController, error)
}

// ServiceInterface captures the set of operations for interacting with the EaseMesh REST apis of the service resource.
type ServiceInterface interface {
	Get(context.Context, string) (*resource.Service, error)
	Patch(context.Context, *resource.Service) error
	Create(context.Context, *resource.Service) error
	Delete(context.Context, string) error
	List(context.Context) ([]*resource.Service, error)
}

// ServiceInstanceInterface captures the set of operations for interacting with the EaseMesh REST apis of the service instance resource.
type ServiceInstanceInterface interface {
	Get(ctx context.Context, serviceName, instanceID string) (*resource.ServiceInstance, error)
	Delete(ctx context.Context, serviceName, instanceID string) error
	List(context.Context) ([]*resource.ServiceInstance, error)
}

// ResilienceInterface captures the set of operations for interacting with the EaseMesh REST apis of the resilience resource.
type ResilienceInterface interface {
	Get(context.Context, string) (*resource.Resilience, error)
	Patch(context.Context, *resource.Resilience) error
	Create(context.Context, *resource.Resilience) error
	Delete(context.Context, string) error
	List(context.Context) ([]*resource.Resilience, error)
}

// IngressInterface captures the set of operations for interacting with the EaseMesh REST apis of the ingress resource.
type IngressInterface interface {
	Get(context.Context, string) (*resource.Ingress, error)
	Patch(context.Context, *resource.Ingress) error
	Create(context.Context, *resource.Ingress) error
	Delete(context.Context, string) error
	List(context.Context) ([]*resource.Ingress, error)
}

// CustomResourceKindInterface captures the set of operations for interacting with the EaseMesh REST apis of the custom resource kind.
type CustomResourceKindInterface interface {
	Get(context.Context, string) (*resource.CustomResourceKind, error)
	Patch(context.Context, *resource.CustomResourceKind) error
	Create(context.Context, *resource.CustomResourceKind) error
	Delete(context.Context, string) error
	List(context.Context) ([]*resource.CustomResourceKind, error)
}

// CustomResourceInterface captures the set of operations for interacting with the EaseMesh REST apis of the custom resource.
type CustomResourceInterface interface {
	Get(context.Context, string, string) (*resource.CustomResource, error)
	Patch(context.Context, *resource.CustomResource) error
	Create(context.Context, *resource.CustomResource) error
	Delete(context.Context, string, string) error
	List(context.Context, string) ([]*resource.CustomResource, error)
}
