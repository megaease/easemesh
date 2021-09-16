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
)

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

type (
	// VersionKind holds version and kind information for APIs
	VersionKind struct {
		APIVersion string `yaml:"apiVersion" yaml:"apiVersion" jsonschema:"omitempty"`
		Kind       string `yaml:"kind" yaml:"kind" jsonschema:"required"`
	}

	// MetaData is meta data for resources of the EaseMesh
	MetaData struct {
		Name   string            `yaml:"name" yaml:"name" jsonschema:"required"`
		Labels map[string]string `yaml:"labels,omitempty" yaml:"labels,omitempty" jsonschema:"omitempty"`
	}

	// MeshResource holds common information for a resource of the EaseMesh
	MeshResource struct {
		VersionKind `yaml:",inline" yaml:",inline"`
		MetaData    MetaData `yaml:"metadata" yaml:"metadata" jsonschema:"required"`
	}

	// MeshObject describes what's feature of a comman EaseMesh object
	MeshObject interface {
		Name() string
		Kind() string
		APIVersion() string
		Labels() map[string]string
	}

	// MeshController is the spec of MeshController on Easegress.
	MeshController struct {
		MeshResource        `yaml:",inline"`
		MeshControllerAdmin `yaml:",inline"`
	}

	// MeshControllerV1Alpha1 is the v1alphv1 version of mesh controller.
	MeshControllerV1Alpha1 struct {
		Kind                string `yaml:"kind"`
		Name                string `yaml:"name"`
		MeshControllerAdmin `yaml:",inline"`
	}

	// MeshControllerAdmin is the admin config of mesh controller.
	MeshControllerAdmin struct {
		// HeartbeatInterval is the interval for one service instance reporting its heartbeat.
		HeartbeatInterval string `yaml:"heartbeatInterval" jsonschema:"required,format=duration"`

		// RegistryTime indicates which protocol the registry center accepts.
		RegistryType string `yaml:"registryType" jsonschema:"required"`

		// APIPort is the port for worker's API server
		APIPort int `yaml:"apiPort" jsonschema:"required"`

		// IngressPort is the port for http server in mesh ingress
		IngressPort int `yaml:"ingressPort" jsonschema:"required"`

		// ExternalServiceRegistry is the external service registry name.
		ExternalServiceRegistry string `yaml:"externalServiceRegistry" jsonschema:"omitempty"`
	}

	// Tenant describes tenant resource of the EaseMesh
	Tenant struct {
		MeshResource `yaml:",inline"`
		Spec         *TenantSpec `yaml:"spec" jsonschema:"required"`
	}

	// TenantSpec describes whats service resided in
	TenantSpec struct {
		Services    []string `yaml:"services" jsonschema:"omitempty"`
		Description string   `yaml:"description" jsonschema:"omitempty"`
	}

	// Service describes service resource of the EaseMesh
	Service struct {
		MeshResource `yaml:",inline"`
		Spec         *ServiceSpec `yaml:"spec" jsonschema:"required"`
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

	ServiceInstance struct {
		MeshResource `yaml:",inline"`
		Spec         *v1alpha1.ServiceInstance `yaml:"spec" jsonschema:"required"`
	}

	ServiceInstanceSpec struct {
		RegistryName string            `yaml:"registryName" jsonschema:"required"`
		ServiceName  string            `yaml:"serviceName" jsonschema:"required"`
		InstanceID   string            `yaml:"instanceID" jsonschema:"required"`
		IP           string            `yaml:"ip" jsonschema:"required"`
		Port         uint32            `yaml:"port" jsonschema:"required"`
		RegistryTime string            `yaml:"registryTime" jsonschema:"omitempty"`
		Labels       map[string]string `yaml:"labels" jsonschema:"omitempty"`
		Status       string            `yaml:"status" jsonschema:"omitempty"`
	}

	// Canary describes canary resource of the EaseMesh
	Canary struct {
		MeshResource `yaml:",inline"`
		Spec         *v1alpha1.Canary `yaml:"spec" jsonschema:"required"`
	}

	// ObservabilityTracings describes observability tracings resource of the EaseMesh
	ObservabilityTracings struct {
		MeshResource `yaml:",inline"`
		Spec         *v1alpha1.ObservabilityTracings `yaml:"spec" jsonschema:"required"`
	}

	// ObservabilityOutputServer describes observability output server resource of the EaseMesh
	ObservabilityOutputServer struct {
		MeshResource `yaml:",inline"`
		Spec         *v1alpha1.ObservabilityOutputServer `yaml:"spec" jsonschema:"required"`
	}

	// ObservabilityMetrics describes observability metrics resource of the EaseMesh
	ObservabilityMetrics struct {
		MeshResource `yaml:",inline"`
		Spec         *v1alpha1.ObservabilityMetrics `yaml:"spec" jsonschema:"required"`
	}

	// LoadBalance describes loadbalance resource of the EaseMesh
	LoadBalance struct {
		MeshResource `yaml:",inline"`
		Spec         *v1alpha1.LoadBalance `yaml:"spec" jsonschema:"required"`
	}

	// Resilience describes resilience resource of the EaseMesh
	Resilience struct {
		MeshResource `yaml:",inline"`
		Spec         *v1alpha1.Resilience `yaml:"spec" jsonschema:"required"`
	}

	// Ingress describes ingress resource of the EaseMesh
	Ingress struct {
		MeshResource `yaml:",inline"`
		Spec         *IngressSpec `yaml:"spec" jsonschema:"required"`
	}

	// IngressSpec wraps all route rules
	IngressSpec struct {
		Rules []*v1alpha1.IngressRule `yaml:"rules" jsonschema:"required"`
	}
)

var _ MeshObject = &MeshController{}
var _ MeshObject = &Service{}
var _ MeshObject = &Tenant{}
var _ MeshObject = &Canary{}
var _ MeshObject = &ObservabilityTracings{}
var _ MeshObject = &LoadBalance{}
var _ MeshObject = &Resilience{}
var _ MeshObject = &Ingress{}

// Default set default value for ServiceSpec
func (s *ServiceSpec) Default() {
	if s.Sidecar.DiscoveryType == "" {
		s.Sidecar.DiscoveryType = "eureka"
	}

	if s.Sidecar.Address == "" {
		s.Sidecar.Address = "127.0.0.1"
	}

	if s.Sidecar.IngressPort == 0 {
		s.Sidecar.IngressPort = DefaultSideIngressPort
	}

	if s.Sidecar.IngressProtocol == "" {
		s.Sidecar.IngressProtocol = DefaultSideIngressProtocol
	}

	if s.Sidecar.EgressPort == 0 {
		s.Sidecar.EgressPort = DefaultSideEgressPort
	}
	if s.Sidecar.EgressProtocol == "" {
		s.Sidecar.EgressProtocol = DefaultSideEgressProtocol
	}
	if s.LoadBalance.Policy == "" {
		s.LoadBalance.Policy = LoadBalanceRoundRobinPolicy
	}
}

// Name returns name of the EaseMesh resource
func (m *MeshResource) Name() string {
	return m.MetaData.Name
}

// Kind returns kind of the EaseMesh resource
func (m *MeshResource) Kind() string {
	return m.VersionKind.Kind
}

// APIVersion returns api version of the EaseMesh resource
func (m *MeshResource) APIVersion() string {
	return m.VersionKind.APIVersion
}

// Labels returns labels of the EaseMesh resource
func (m *MeshResource) Labels() map[string]string {
	return m.MetaData.Labels
}

// ToV1Alpha1 converts MeshController resouce to v1alpha1.
func (mc *MeshController) ToV1Alpha1() *MeshControllerV1Alpha1 {
	return &MeshControllerV1Alpha1{
		Kind:                mc.Kind(),
		Name:                mc.Name(),
		MeshControllerAdmin: mc.MeshControllerAdmin,
	}
}

// ToV1Alpha1 converts an Ingress resource to v1alpha1.Ingress
func (ing *Ingress) ToV1Alpha1() *v1alpha1.Ingress {
	result := &v1alpha1.Ingress{}
	result.Name = ing.Name()
	if ing.Spec != nil {
		result.Rules = ing.Spec.Rules
	}
	return result
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

func (si *ServiceInstance) ParseName() (serviceName, instanceID string, err error) {
	ss := strings.Split(si.Name(), "/")

	if len(ss) != 2 {
		return "", "", fmt.Errorf("invalid service instance name (format: serviceName/instanceID)")
	}

	return ss[0], ss[1], nil
}

// ToV1Alpha1 converts a Service resource to v1alpha1.ServiceInstance.
func (si *ServiceInstance) ToV1Alpha1() *v1alpha1.ServiceInstance {
	return si.Spec
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

// ToV1Alpha1 converts a loadbalance resource to v1alpha1.LoadBalance
func (l *LoadBalance) ToV1Alpha1() *v1alpha1.LoadBalance {
	return l.Spec
}

// ToV1Alpha1 converts a Canary resource to v1alpha1.Canary
func (c *Canary) ToV1Alpha1() *v1alpha1.Canary {
	return c.Spec
}

// ToV1Alpha1 converts a Resilience resource to v1alpha1.Resilience
func (r *Resilience) ToV1Alpha1() *v1alpha1.Resilience {
	return r.Spec
}

// ToV1Alpha1 converts a ObservabilityTracings resource to v1alpha1.ObservabilityTracings
func (r *ObservabilityTracings) ToV1Alpha1() (result *v1alpha1.ObservabilityTracings) {
	return r.Spec
}

// ToV1Alpha1 converts a ObservabilityOutputServer resource to v1alpha1.ObservabilityOutputServer
func (r *ObservabilityOutputServer) ToV1Alpha1() (result *v1alpha1.ObservabilityOutputServer) {
	return r.Spec
}

// ToV1Alpha1 converts a ObservabilityMetrics resource to v1alpha1.ObservabilityMetrics
func (r *ObservabilityMetrics) ToV1Alpha1() (result *v1alpha1.ObservabilityMetrics) {
	return r.Spec
}

// ToIngress converts a v1alpha1.Ingress resource to an Ingress resource
func ToIngress(ingress *v1alpha1.Ingress) *Ingress {
	result := &Ingress{
		Spec: &IngressSpec{},
	}
	result.MeshResource = NewIngressResource(DefaultAPIVersion, ingress.Name)
	result.Spec.Rules = ingress.Rules
	return result
}

// ToMeshController converts a MeshControllerV1Alpha1 resouce to a MeshController resource.
func ToMeshController(meshController *MeshControllerV1Alpha1) *MeshController {
	return &MeshController{
		MeshResource:        NewMeshResource(DefaultAPIVersion, meshController.Kind, meshController.Name),
		MeshControllerAdmin: meshController.MeshControllerAdmin,
	}
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

func ToServiceInstance(instance *v1alpha1.ServiceInstance) *ServiceInstance {
	result := &ServiceInstance{
		Spec: &v1alpha1.ServiceInstance{},
	}

	name := fmt.Sprintf("%s/%s", instance.ServiceName, instance.InstanceID)

	result.MeshResource = NewServiceInstanceResource(DefaultAPIVersion, name)
	result.Spec = instance

	return result
}

// ToCanary converts a v1alpha1.Canary resource to a Canary resource
func ToCanary(name string, canary *v1alpha1.Canary) *Canary {
	result := &Canary{
		Spec: &v1alpha1.Canary{},
	}
	result.MeshResource = NewCanaryResource(DefaultAPIVersion, name)
	result.Spec.CanaryRules = canary.CanaryRules
	return result
}

// ToObservabilityTracings converts a v1alpha1.ObservabilityTracings resource to a ObservabilityTracings resource
func ToObservabilityTracings(serviceID string, tracing *v1alpha1.ObservabilityTracings) *ObservabilityTracings {
	result := &ObservabilityTracings{
		Spec: &v1alpha1.ObservabilityTracings{},
	}
	result.MeshResource = NewObservabilityTracingsResource(DefaultAPIVersion, serviceID)
	result.Spec = tracing
	return result
}

// ToObservabilityMetrics converts a v1alpha1.ObservabilityMetrics resource to a ObservabilityMetrics resource
func ToObservabilityMetrics(serviceID string, metrics *v1alpha1.ObservabilityMetrics) *ObservabilityMetrics {
	result := &ObservabilityMetrics{
		Spec: &v1alpha1.ObservabilityMetrics{},
	}
	result.MeshResource = NewObservabilityMetricsResource(DefaultAPIVersion, serviceID)
	result.Spec = metrics
	return result
}

// ToObservabilityOutputServer converts a v1alpha1.ObservabilityOutputServer resource to a ObservabilityOutputServer resource
func ToObservabilityOutputServer(serviceID string, output *v1alpha1.ObservabilityOutputServer) *ObservabilityOutputServer {
	result := &ObservabilityOutputServer{
		Spec: &v1alpha1.ObservabilityOutputServer{},
	}
	result.MeshResource = NewObservabilityOutputServerResource(DefaultAPIVersion, serviceID)
	result.Spec = output
	return result
}

// ToLoadBalance converts a v1alpha1.LoadBalance resource to a LoadBalance resource
func ToLoadBalance(name string, loadBalance *v1alpha1.LoadBalance) *LoadBalance {
	result := &LoadBalance{
		Spec: &v1alpha1.LoadBalance{},
	}
	result.MeshResource = NewLoadBalanceResource(DefaultAPIVersion, name)
	result.Spec = loadBalance
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

// ToResilience converts a v1alpha1.Resilience resource to a Resilience resource
func ToResilience(name string, resilience *v1alpha1.Resilience) *Resilience {
	result := &Resilience{
		Spec: &v1alpha1.Resilience{},
	}
	result.MeshResource = NewResilienceResource(DefaultAPIVersion, name)
	result.Spec.RateLimiter = resilience.RateLimiter
	result.Spec.Retryer = resilience.Retryer
	result.Spec.CircuitBreaker = resilience.CircuitBreaker
	result.Spec.TimeLimiter = resilience.TimeLimiter
	return result
}
