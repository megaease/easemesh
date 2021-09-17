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
)
