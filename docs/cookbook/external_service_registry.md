
# External Service Registry

When architecture style moves to service mesh, there could be a middle status for the technical stack. Which is that there would be both legacy services and mesh services co-existing. The two different kinds of services also want to call each other.

EaseMesh uses internal component Etcd in the control plane to play the role of service registry, and the external services may use Consul, Nacos, Eureka, Zookeeper as a service registry. Based on service registry controllers from Easegress, we develop a registry syncer to synchronize services between internal Etcd and the external service registry.

As an example, we use `emctl` to apply consul service registry:

```bash
$ emctl apply -f consul-service-registry.yaml
```

consul-service-registry.yaml

```yaml
apiVersion: mesh.megaease.com/v1alpha1
kind: ConsulServiceRegistry
metadata:
  name: consul-service-registry
address: consul-server-0.consul-server.default:8500
scheme: http
datacenter: ""
token: ""
namespace: ""
syncInterval: 5s
serviceTags: []
```

Then we need to specify external service registry name in MeshController:

```bash
$ emctl apply -f mesh-controller.yaml
```

mesh-controller.yaml

```yaml
apiVersion: mesh.megaease.com/v1alpha1
kind: MeshController
metadata:
  name: easemesh-controller
apiPort: 13009
ingressPort: 19527
heartbeatInterval: 5s
# +
externalServiceRegistry: consul-service-registry
registryType: consul
```

And we could use `emctl get service && emctl get serviceinstance` to query the service and instance information from different sources.

> NOTICE: The registry syncer transforms the service entry format between different registries. The connectability needs to be guaranteed by themselves. For example, services run in the same Kubernetes cluster would be connectable without other operations.
