package resource

import "github.com/pkg/errors"

type ObjectCreator interface {
	New(*VersionKind) (MeshObject, error)
}

type createObjectFunc func(*VersionKind) (MeshObject, error)

func NewObjectCreator() ObjectCreator {
	fn := func(kind *VersionKind) (MeshObject, error) {
		switch kind.Kind {
		case KindService:
			return &Service{}, nil
		case KindTenant:
			return &Tenant{}, nil
		case KindLoadBalance:
			return &LoadBalance{}, nil
		case KindCanary:
			return &Canary{}, nil
		case KindObservabilityTracings:
			return &ObservabilityTracings{}, nil
		case KindObservabilityOutputServer:
			return &ObservabilityOutputServer{}, nil
		case KindObservabilityMetrics:
			return &ObservabilityMetrics{}, nil
		case KindResilience:
			return &Resilience{}, nil
		case KindIngress:
			return &Ingress{}, nil
		}
		return nil, errors.Errorf("unknown VersionKind object version: %s, kind: %s", kind.APIVersion, kind.Kind)
	}
	var f createObjectFunc = fn
	return f
}

func (c createObjectFunc) New(k *VersionKind) (MeshObject, error) {
	return c(k)
}
func ForIngressResource(service string) MeshResource {
	return meshResource(apiVersion, KindIngress, service)
}

func ForServiceMeshResource(service string) MeshResource {
	return meshResource(apiVersion, KindService, service)
}

func ForCanaryMeshResource(id string) MeshResource {
	return meshResource(apiVersion, KindCanary, id)
}

func ForLoadBalanceResource(id string) MeshResource {
	return meshResource(apiVersion, KindLoadBalance, id)
}

func ForResilienceResource(id string) MeshResource {
	return meshResource(apiVersion, KindResilience, id)
}

func ForObservabilityTracingsResource(id string) MeshResource {
	return meshResource(apiVersion, KindObservabilityTracings, id)
}

func ForObservabilityMetricsResource(id string) MeshResource {
	return meshResource(apiVersion, KindObservabilityMetrics, id)
}

func ForObservabilityOutputServerResource(id string) MeshResource {
	return meshResource(apiVersion, KindObservabilityOutputServer, id)
}
func ForTenantResource(id string) MeshResource {
	return meshResource(apiVersion, KindTenant, id)
}

func meshResource(api, kind, id string) MeshResource {
	return MeshResource{
		APIVersion: api,
		Kind:       kind,
		MetaData: MetaData{
			Name: id,
		},
	}
}
