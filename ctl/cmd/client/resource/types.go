package resource

import (
	"github.com/megaease/easemesh-api/v1alpha1"
)

const (
	apiVersion = "mesh.megaease.com/v1alpha1"
)
const (
	LoadBalanceRoundRobinPolicy = "roundRobin"
	DefaultSideIngressProtocol  = "http"
	DefaultSideEgressProtocol   = "http"
	DefaultSideIngressPort      = 13001
	DefaultSideEgressPort       = 13002
)

const (
	KindService = "service"
	KindCanary  = "canary"
	//
	KindObservabilityMetrics      = "observabilityMetrics"
	KindObservabilityTracings     = "observabilityTracings"
	KindObservabilityOutputServer = "observabilityOutputServer"
	//
	KindTenant      = "tenant"
	KindLoadBalance = "loadbalance"
	KindResilience  = "resilience"

	//
	KindIngress = "ingress"
)

type (
	MetaData struct {
		Name   string            `yaml:"name" jsonschema:"required"`
		Labels map[string]string `yaml:"labels" jsonschema:"omitempty"`
	}

	MeshResource struct {
		APIVersion string   `yaml:"apiVersion" jsonschema:"required"`
		Kind       string   `yaml:"kind" jsonschema:"required"`
		MetaData   MetaData `yaml:"metadata" jsonschema:"required"`
	}

	MeshObject interface {
		Name() string
		GetKind() string
		GetAPIVersion() string
	}

	VersionKind struct {
		APIVersion string `yaml:"apiVersion" jsonschema:"required"`
		Kind       string `yaml:"kind" jsonschema:"required"`
	}
)

type (
	Tenant struct {
		MeshResource
		Spec v1alpha1.Tenant `yaml:"spec" jsonschema:"omitempty"`
	}

	ServiceSpec struct {
		RegisterTenant string                `yaml:"registerTenant" jsonschema:"omitempty"`
		LoadBalance    *v1alpha1.LoadBalance `yaml:"loadbalance" jsonschema:"omitempty"`
		SideCar        *v1alpha1.Sidecar     `yaml:"sideCar" jsonschema:"omitempty"`
	}

	Service struct {
		MeshResource
		Spec ServiceSpec `yaml:"spec" jsonschema:"omitempty"`
	}

	Canary struct {
		MeshResource
		Spec v1alpha1.Canary `yaml:"spec" jsonschema:"omitempty"`
	}

	ObservabilityTracings struct {
		MeshResource
		Spec v1alpha1.ObservabilityTracings `yaml:"spec" jsonschema:"omitempty"`
	}

	ObservabilityOutputServer struct {
		MeshResource
		Spec v1alpha1.ObservabilityOutputServer `yaml:"spec" jsonschema:"omitempty"`
	}

	ObservabilityMetrics struct {
		MeshResource
		Spec v1alpha1.ObservabilityMetrics `yaml:"spec" jsonschema:"omitempty"`
	}
	LoadBalance struct {
		MeshResource
		Spec v1alpha1.LoadBalance `yaml:"spec" jsonschema:"omitempty"`
	}

	Resilience struct {
		MeshResource
		Spec v1alpha1.Resilience `yaml:"spec" jsonschema:"omitempty"`
	}

	IngressSpec struct {
		Rules []*v1alpha1.IngressRule `yaml:"rules" jsonschema:"omitempty"`
	}

	Ingress struct {
		MeshResource
		Spec IngressSpec `yaml:"spec" jsonschema:"omitempty"`
	}
)

var _ MeshObject = &Service{}
var _ MeshObject = &Tenant{}
var _ MeshObject = &Canary{}
var _ MeshObject = &ObservabilityTracings{}
var _ MeshObject = &LoadBalance{}
var _ MeshObject = &Resilience{}
var _ MeshObject = &Ingress{}

func (s *ServiceSpec) Default() {
	if s.SideCar.DiscoveryType == "" {
		s.SideCar.DiscoveryType = "eureka"
	}

	if s.SideCar.Address == "" {
		s.SideCar.Address = "127.0.0.1"
	}

	if s.SideCar.IngressPort == 0 {
		s.SideCar.IngressPort = DefaultSideIngressPort
	}

	if s.SideCar.IngressProtocol == "" {
		s.SideCar.IngressProtocol = DefaultSideIngressProtocol
	}

	if s.SideCar.EgressPort == 0 {
		s.SideCar.EgressPort = DefaultSideEgressPort
	}
	if s.SideCar.EgressProtocol == "" {
		s.SideCar.EgressProtocol = DefaultSideEgressProtocol
	}
	if s.LoadBalance.Policy == "" {
		s.LoadBalance.Policy = LoadBalanceRoundRobinPolicy
	}
}

func (m *MeshResource) Name() string {
	return m.MetaData.Name
}

func (m *MeshResource) GetKind() string {
	return m.Kind
}

func (m *MeshResource) GetAPIVersion() string {
	return m.APIVersion
}

func (s *Ingress) ToV1Alpha1() (result v1alpha1.Ingress) {
	result.Name = s.Name()
	result.Rules = s.Spec.Rules
	return
}

func (s *Service) ToV1Alpha1() (result v1alpha1.Service) {
	result.Name = s.Name()
	result.RegisterTenant = s.Spec.RegisterTenant
	result.LoadBalance = s.Spec.LoadBalance
	result.Sidecar = s.Spec.SideCar
	return
}

func (t *Tenant) ToV1Alpha1() (result v1alpha1.Tenant) {
	result.Name = t.Name()
	result.Services = t.Spec.Services
	result.Description = t.Spec.Description
	return result
}

func (l *LoadBalance) ToV1Alpha1() (result v1alpha1.LoadBalance) {
	result = l.Spec
	return
}

func (c *Canary) ToV1Alpha1() (result v1alpha1.Canary) {
	result = c.Spec
	return
}

func (r *Resilience) ToV1Alpha1() (result v1alpha1.Resilience) {
	result = r.Spec
	return
}

func (r *ObservabilityTracings) ToV1Alpha1() (result v1alpha1.ObservabilityTracings) {
	result = r.Spec
	return
}

func (r *ObservabilityOutputServer) ToV1Alpha1() (result v1alpha1.ObservabilityOutputServer) {
	result = r.Spec
	return
}

func (r *ObservabilityMetrics) ToV1Alpha1() (result v1alpha1.ObservabilityMetrics) {
	result = r.Spec
	return
}

func ToIngress(ingress *v1alpha1.Ingress) (result Ingress) {
	result.MeshResource = ForIngressResource(ingress.Name)
	result.Spec.Rules = ingress.Rules
	return
}

func ToService(service *v1alpha1.Service) (result Service) {
	result.MeshResource = ForServiceMeshResource(service.Name)
	result.Spec.SideCar = service.Sidecar
	result.Spec.LoadBalance = service.LoadBalance
	return
}

func ToCanary(name string, canary *v1alpha1.Canary) (result Canary) {
	result.MeshResource = ForCanaryMeshResource(name)
	result.Spec.CanaryRules = canary.CanaryRules
	return
}

func ToObservabilityTracings(serviceID string, tracing *v1alpha1.ObservabilityTracings) (result ObservabilityTracings) {
	result.MeshResource = ForObservabilityTracingsResource(serviceID)
	result.Spec = *tracing
	return
}

func ToObservabilityMetrics(serviceID string, metrics *v1alpha1.ObservabilityMetrics) (result ObservabilityMetrics) {
	result.MeshResource = ForObservabilityMetricsResource(serviceID)
	result.Spec = *metrics
	return
}

func ToObservabilityOutputServer(serviceID string, output *v1alpha1.ObservabilityOutputServer) (result ObservabilityOutputServer) {
	result.MeshResource = ForObservabilityOutputServerResource(serviceID)
	result.Spec = *output
	return
}

func ToLoadbalance(name string, loadBalance *v1alpha1.LoadBalance) (result LoadBalance) {
	result.MeshResource = ForLoadBalanceResource(name)
	result.Spec = *loadBalance
	return
}

func ToTenant(tenant *v1alpha1.Tenant) (result Tenant) {
	result.MeshResource = ForTenantResource(tenant.Name)
	result.Spec.Services = tenant.Services
	return
}

func ToResilience(name string, resilience *v1alpha1.Resilience) (result Resilience) {
	result.MeshResource = ForResilienceResource(name)
	result.Spec.RateLimiter = resilience.RateLimiter
	result.Spec.Retryer = resilience.Retryer
	result.Spec.CircuitBreaker = resilience.CircuitBreaker
	result.Spec.TimeLimiter = resilience.TimeLimiter
	return
}
