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

const (
	apiURL = "/apis/v1"

	// MeshControllersURL is the mesh controller prefix.
	MeshControllersURL = apiURL + "/objects"

	// MeshControllerURL is the mesh controller path.
	MeshControllerURL = apiURL + "/objects/%s"

	// ConsulServiceRegistrysURL is the consul service registry path.
	ConsulServiceRegistrysURL = apiURL + "/objects"

	// ConsulServiceRegistryURL is the consul service registry path.
	ConsulServiceRegistryURL = apiURL + "/objects/%s"

	// MeshTenantsURL is the mesh tenant prefix.
	MeshTenantsURL = apiURL + "/mesh/tenants"

	// MeshTenantURL is the mesh tenant path.
	MeshTenantURL = apiURL + "/mesh/tenants/%s"

	// MeshServicesURL is mesh service prefix.
	MeshServicesURL = apiURL + "/mesh/services"

	// MeshServiceURL is the mesh service path.
	MeshServiceURL = apiURL + "/mesh/services/%s"

	// MeshServiceCanaryURL is the mesh service canary path.
	MeshServiceCanaryURL = apiURL + "/mesh/services/%s/canary"

	// MeshServiceResilienceURL is the mesh service resilience path.
	MeshServiceResilienceURL = apiURL + "/mesh/services/%s/resilience"

	// MeshServiceLoadBalanceURL is the mesh service load balance path.
	MeshServiceLoadBalanceURL = apiURL + "/mesh/services/%s/loadbalance"

	// MeshServiceOutputServerURL is the mesh service output server path.
	MeshServiceOutputServerURL = apiURL + "/mesh/services/%s/outputserver"

	// MeshServiceTracingsURL is the mesh service tracings path.
	MeshServiceTracingsURL = apiURL + "/mesh/services/%s/tracings"

	// MeshServiceMetricsURL is the mesh service metrics path.
	MeshServiceMetricsURL = apiURL + "/mesh/services/%s/metrics"

	// MeshServiceInstancesURL is the mesh service prefix.
	MeshServiceInstancesURL = apiURL + "/mesh/serviceinstances"

	// MeshServiceInstanceURL is the mesh service path.
	MeshServiceInstanceURL = apiURL + "/mesh/serviceinstances/%s/%s"

	// MeshIngressesURL is the mesh ingress prefix.
	MeshIngressesURL = apiURL + "/mesh/ingresses"

	// MeshIngressURL is the mesh ingress path.
	MeshIngressURL = apiURL + "/mesh/ingresses/%s"

	// MeshCustomResourceKindsURL is the mesh custom resource kind prefix.
	MeshCustomResourceKindsURL = apiURL + "/mesh/customresourcekinds"

	// MeshCustomResourceKindURL is the mesh custom resource kind path.
	MeshCustomResourceKindURL = apiURL + "/mesh/customresourcekinds/%s"

	// MeshAllCustomResourcesURL is the mesh custom resource.
	MeshAllCustomResourcesURL = apiURL + "/mesh/customresources"

	// MeshCustomResourcesURL is the mesh custom resource prefix.
	MeshCustomResourcesURL = apiURL + "/mesh/customresources/%s"

	// MeshCustomResourceURL is the mesh custom resource path.
	MeshCustomResourceURL = apiURL + "/mesh/customresources/%s/%s"
)
