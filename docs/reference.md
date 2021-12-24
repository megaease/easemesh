## EaseMesh Reference


### Tenant and Service

- Tenant
- Service

| <p align="left">Tenant</p> | <p align="left">Service</p> |
|-|-|
|<p align="left">Tenant is the logic group of mesh services. [Tenant specification](https://github.com/megaease/easemesh-api/blob/master/v1alpha1/meshmodel.md#easemesh.v1alpha1.Tenant) provide a scope for service's names. Names of service need to be unique within a tenant.</p>|<p align="left">[Service specification]() define the unique service name in a tenant, a service is may related to other resource such as observability, resilience and traffics</p>|


### Traffic

| <p align="left">Traffic Split</p> | <p align="left">Load Balance </p>|
|-|-|
| <p>[Canary specification](https://github.com/megaease/easemesh-api/blob/master/v1alpha1/meshmodel.md#easemesh.v1alpha1.Canary)  describes how to correctly split traffic to service instance</p>|<p> [LoadBalance specification](https://github.com/megaease/easemesh-api/blob/master/v1alpha1/meshmodel.md#easemesh.v1alpha1.LoadBalance) describes how to load balance between instance</p> |

| <p align="left">Traffic Control</p> | <p align="left">HTTP Route Group </p> |
|-|-|
| <p>[Traffic Control](https://github.com/megaease/easemesh-api/blob/master/v1alpha1/meshmodel.md#easemesh.v1alpha1.TrafficTarget) describes how to control the access to services.</p> |<p>[HTTP Route Group](https://github.com/megaease/easemesh-api/blob/main/v1alpha1/meshmodel.md#easemesh.v1alpha1.HTTPRouteGroup) describes the HTTP traffic rules.</p>|

### Observability

Observability consists of three components:
- ObservabilityOutputServer
- ObservabilityTracings
- ObservabilityMetrics

|<p align="left">ObservabilityMetrics</p>|<p align="left">ObservabilityTracings</p>|<p align="left">ObservabilityOutputServer</p>|
|-|-|-|
|<p align="left">[ObservabilityMetrics specification](https://github.com/megaease/easemesh-api/blob/master/v1alpha1/meshmodel.md#easemesh.v1alpha1.ObservabilityMetrics) describe how to control metrics collection</p>|<p align="left">[ObservabilityTracings specification](https://github.com/megaease/easemesh-api/blob/master/v1alpha1/meshmodel.md#easemesh.v1alpha1.ObservabilityTracings) describe how to control tracing collection</p>|<p align="left">[ObservabilityOutputServer specification](https://github.com/megaease/easemesh-api/blob/master/v1alpha1/meshmodel.md#easemesh.v1alpha1.ObservabilityOutputServer) describe how to report tracing and metrics to backend</p>| 



### Resilience

Resilience configures four key types of features
- RateLimiter.
- Retryer 

|<p align="left">RateLimiter</p>|<p align="left">Retryer</p>|
|-|-|
|<p align="left">RateLimiter specification describe the sidecar how to rate limits the inbound traffic</p>|<p align="left">Retryer specification describe the sidecar how to issue a repeat request </p>|

- CircuitBreaker
- TimeLimiter.

|<p align="left">CircuitBreaker</p>|<p align="left">TimeLimiter</p>|
|-|-|
|<p align="left">CircuitBreaker specification describes the sidecar how to circuit break a downstream service</p>|<p align="left">TimeLimiter specification describes the sidecar how to control request time out </p>|


### Ingress
Ingress is the spec of mesh ingress.

|<p align="left">Ingress </p>|
|-|
|<p align="left">[Ingress specification](https://github.com/megaease/easemesh-api/blob/master/v1alpha1/meshmodel.md#ingress) describes how to router the traffic (or request) that came from outside to appropriate destinations (service instances) <br /></p>|

### Sidecar

|<p align="left">Ingress </p>|
|-|
|[Sidecar specification]() describe the sidecar behavior, for example, define the sidecar should listen on what's port for inbound/outbound traffic|

### Extensibility

|<p align="left">CustomResourceKind</p>|
|-|
|<p align="left">[CustomResourceKind](https://github.com/megaease/easemesh-api/blob/main/v1alpha1/meshmodel.md#easemesh.v1alpha1.CustomResourceKind) describes the specification of a Custom Resource. Shadow Service is an implementation of Custom Resource.</p>|


### Shadow Service

[Shadow Service](./shadow-service-manual.md) is an implementation of Custom Resource, its `CustomResourceKind` is defined as:

```yaml
kind: CustomResourceKind
apiVersion: mesh.megaease.com/v1alpla1
metadata:
  name: ShadowService
spec:
  jsonSchema:
    type: object
    properties:
      name:
        type: string
      namespace:
        type: string
      serviceName:
        type: string
      mysql:
        type: object
        properties:
          uris:
            type: string
          userName:
            type: string
          password:
            type: string
      kafka:
        type: object
        properties:
          uris:
            type: string
      redis:
        type: object
        properties:
          uris:
            type: string
          userName:
            type: string
          password:
            type: string
      rabbitMq:
        type: object
        properties:
          uris:
            type: string
          userName:
            type: string
          password:
            type: string
      elasticSearch:
        type: object
        properties:
          uris:
            type: string
          userName:
            type: string
          password:
            type: string
```

---
See [EaseMesh Reference](https://github.com/megaease/easemesh-api/blob/master/v1alpha1/meshmodel.md#easemesh.v1alpha1.Service)