package resource

import (
	"github.com/megaease/easemesh-api/v1alpha1"
)

const (
	DefaultAPIVersion = "mesh.megaease.com/v1alpha1"

	LoadBalanceRoundRobinPolicy = "roundRobin"
	DefaultSideIngressProtocol  = "http"
	DefaultSideEgressProtocol   = "http"
	DefaultSideIngressPort      = 13001
	DefaultSideEgressPort       = 13002

	KindService = "Service"
	KindCanary  = "Canary"

	KindObservabilityMetrics      = "ObservabilityMetrics"
	KindObservabilityTracings     = "ObservabilityTracings"
	KindObservabilityOutputServer = "ObservabilityOutputServer"

	KindTenant      = "Tenant"
	KindLoadBalance = "Loadbalance"
	KindResilience  = "Resilience"

	KindIngress = "Ingress"
)

type (
	VersionKind struct {
		APIVersion string `yaml:"apiVersion" jsonschema:"required"`
		Kind       string `yaml:"kind" jsonschema:"required"`
	}

	MetaData struct {
		Name   string            `yaml:"name" jsonschema:"required"`
		Labels map[string]string `yaml:"labels" jsonschema:"omitempty"`
	}

	MeshResource struct {
		VersionKind `yaml:",inline"`
		MetaData    MetaData `yaml:"metadata" jsonschema:"required"`
	}

	MeshObject interface {
		Name() string
		Kind() string
		APIVersion() string
		Labels() map[string]string
	}
)

type (
	Tenant struct {
		MeshResource `yaml:",inline"`
		Services     []string `yaml:"services" jsonschema:"omitempty"`
		Description  string   `yaml:"description" jsonschema:"omitempty"`
	}

	ServiceSpec struct {
		RegisterTenant string `yaml:"registerTenant" jsonschema:"required"`

		Sidecar       *v1alpha1.Sidecar       `yaml:"sidecar" jsonschema:"required"`
		Resilience    *v1alpha1.Resilience    `yaml:"resilience" jsonschema:"omitempty"`
		Canary        *v1alpha1.Canary        `yaml:"canary" jsonschema:"omitempty"`
		LoadBalance   *v1alpha1.LoadBalance   `yaml:"loadBalance" jsonschema:"omitempty"`
		Observability *v1alpha1.Observability `yaml:"observability" jsonschema:"omitempty"`
	}

	Service struct {
		MeshResource `yaml:",inline"`
		Spec         *ServiceSpec `yaml:"spec" jsonschema:"omitempty"`
	}

	Canary struct {
		MeshResource `yaml:",inline"`
		Spec         *v1alpha1.Canary `yaml:"spec" jsonschema:"omitempty"`
	}

	ObservabilityTracings struct {
		MeshResource `yaml:",inline"`
		Spec         *v1alpha1.ObservabilityTracings `yaml:"spec" jsonschema:"omitempty"`
	}

	ObservabilityOutputServer struct {
		MeshResource `yaml:",inline"`
		Spec         *v1alpha1.ObservabilityOutputServer `yaml:"spec" jsonschema:"omitempty"`
	}

	ObservabilityMetrics struct {
		MeshResource `yaml:",inline"`
		Spec         *v1alpha1.ObservabilityMetrics `yaml:"spec" jsonschema:"omitempty"`
	}
	LoadBalance struct {
		MeshResource `yaml:",inline"`
		Spec         *v1alpha1.LoadBalance `yaml:"spec" jsonschema:"omitempty"`
	}

	Resilience struct {
		MeshResource `yaml:",inline"`
		Spec         *v1alpha1.Resilience `yaml:"spec" jsonschema:"omitempty"`
	}

	IngressSpec struct {
		Rules []*v1alpha1.IngressRule `yaml:"rules" jsonschema:"omitempty"`
	}

	Ingress struct {
		MeshResource `yaml:",inline"`
		Spec         *IngressSpec `yaml:"spec" jsonschema:"omitempty"`
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

func (m *MeshResource) Name() string {
	return m.MetaData.Name
}

func (m *MeshResource) Kind() string {
	return m.VersionKind.Kind
}

func (m *MeshResource) APIVersion() string {
	return m.VersionKind.APIVersion
}

func (m *MeshResource) Labels() map[string]string {
	return m.MetaData.Labels
}

func (ing *Ingress) ToV1Alpha1() *v1alpha1.Ingress {
	result := &v1alpha1.Ingress{}
	result.Name = ing.Name()
	if ing.Spec != nil {
		result.Rules = ing.Spec.Rules
	}
	return result
}

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

func (t *Tenant) ToV1Alpha1() *v1alpha1.Tenant {
	result := &v1alpha1.Tenant{}
	result.Name = t.Name()
	result.Services = t.Services
	result.Description = t.Description
	return result
}

func (l *LoadBalance) ToV1Alpha1() *v1alpha1.LoadBalance {
	return l.Spec
}

func (c *Canary) ToV1Alpha1() *v1alpha1.Canary {
	return c.Spec
}

func (r *Resilience) ToV1Alpha1() *v1alpha1.Resilience {
	return r.Spec
}

func (r *ObservabilityTracings) ToV1Alpha1() (result *v1alpha1.ObservabilityTracings) {
	return r.Spec
}

func (r *ObservabilityOutputServer) ToV1Alpha1() (result *v1alpha1.ObservabilityOutputServer) {
	return r.Spec
}

func (r *ObservabilityMetrics) ToV1Alpha1() (result *v1alpha1.ObservabilityMetrics) {
	return r.Spec
}

func ToIngress(ingress *v1alpha1.Ingress) *Ingress {
	result := &Ingress{
		Spec: &IngressSpec{},
	}
	result.MeshResource = NewIngressResource(DefaultAPIVersion, ingress.Name)
	result.Spec.Rules = ingress.Rules
	return result
}

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

func ToCanary(name string, canary *v1alpha1.Canary) *Canary {
	result := &Canary{
		Spec: &v1alpha1.Canary{},
	}
	result.MeshResource = NewCanaryResource(DefaultAPIVersion, name)
	result.Spec.CanaryRules = canary.CanaryRules
	return result
}

func ToObservabilityTracings(serviceID string, tracing *v1alpha1.ObservabilityTracings) *ObservabilityTracings {
	result := &ObservabilityTracings{
		Spec: &v1alpha1.ObservabilityTracings{},
	}
	result.MeshResource = NewObservabilityTracingsResource(DefaultAPIVersion, serviceID)
	result.Spec = tracing
	return result
}

func ToObservabilityMetrics(serviceID string, metrics *v1alpha1.ObservabilityMetrics) *ObservabilityMetrics {
	result := &ObservabilityMetrics{
		Spec: &v1alpha1.ObservabilityMetrics{},
	}
	result.MeshResource = NewObservabilityMetricsResource(DefaultAPIVersion, serviceID)
	result.Spec = metrics
	return result
}

func ToObservabilityOutputServer(serviceID string, output *v1alpha1.ObservabilityOutputServer) *ObservabilityOutputServer {
	result := &ObservabilityOutputServer{
		Spec: &v1alpha1.ObservabilityOutputServer{},
	}
	result.MeshResource = NewObservabilityOutputServerResource(DefaultAPIVersion, serviceID)
	result.Spec = output
	return result
}

func ToLoadbalance(name string, loadBalance *v1alpha1.LoadBalance) *LoadBalance {
	result := &LoadBalance{
		Spec: &v1alpha1.LoadBalance{},
	}
	result.MeshResource = NewLoadBalanceResource(DefaultAPIVersion, name)
	result.Spec = loadBalance
	return result
}

func ToTenant(tenant *v1alpha1.Tenant) *Tenant {
	result := &Tenant{}
	result.MeshResource = NewTenantResource(DefaultAPIVersion, tenant.Name)
	result.Services = tenant.Services
	result.Description = tenant.Description
	return result
}

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
