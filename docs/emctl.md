# EaseMesh Command-Line

`emctl` is the dedicated command to handle resources of EaseMesh, which runs in [Easegress](https://github.com/megaease/easegress) MeshController who has different roles in different instances. `MeshController` will register its own admin API in `Easegress`, so the server flag in `emctl` keeps the same as Easegress's.

Running `emctl --help`  or `emctl help <subcommand>` can get details about every subcommand.

## emctl install

Deploy infrastructure components of the EaseMesh.

```bash
emctl install [flags]

# Examples
emctl install --mesh-namespace mesh-demo --clean-when-failed
```

| Flags                                           | Shorthand | Description                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                | Description |
| ----------------------------------------------- | --------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ | ----------- |
| --clean-when-failed                             |           | Clean resources when installation failed (default true)                                                                                                                                                                                                                                                                                                                                                                                                                                                                                    |             |
| --easegress-image string                        |           | Easegress image name (default "megaease/easegress:latest")                                                                                                                                                                                                                                                                                                                                                                                                                                                                                 |             |
| --easemesh-control-plane-replicas int           |           | Mesh control plane reaplicas (default 3)                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                   |             |
| --easemesh-ingress-replicas int                 |           | Mesh ingress controller replicas (default 1)                                                                                                                                                                                                                                                                                                                                                                                                                                                                                               |             |
| --easemesh-operator-image string                |           | Mesh operator image name (default "megaease/easemesh-operator:latest")                                                                                                                                                                                                                                                                                                                                                                                                                                                                     |             |
| --easemesh-operator-replicas int                |           | Mesh operator controller replicas (default 1)                                                                                                                                                                                                                                                                                                                                                                                                                                                                                              |             |
| --file string                                   | -f        | A yaml file specifying the install params                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                  |             |
| --heartbeat-interval int                        |           | Heartbeart interval for mesh service (default 5)                                                                                                                                                                                                                                                                                                                                                                                                                                                                                           |             |
| --help                                          | -h        | help for install                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                           |             |
| --image-registry-url string                     |           | Image registry URL (default "docker.io")                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                   |             |
| --mesh-control-plane-admin-port int             |           | Port of mesh control plane admin for management (default 2381)                                                                                                                                                                                                                                                                                                                                                                                                                                                                             |             |
| --mesh-control-plane-check-healthz-max-time int |           | Max timeout in second for checking control panel component whether ready or not (default 60)                                                                                                                                                                                                                                                                                                                                                                                                                                               |             |
| --mesh-control-plane-client-port int            |           | Mesh control plane client port for remote accessing (default 2379)                                                                                                                                                                                                                                                                                                                                                                                                                                                                         |             |
| --mesh-control-plane-peer-port int              |           | Port of mesh control plane for consensus each other (default 2380)                                                                                                                                                                                                                                                                                                                                                                                                                                                                         |             |
| --mesh-control-plane-pv-capacity string         |           | EaseMesh control plane needs PersistentVolume to store data. You need to create PersistentVolume in advance and specify its storageClassName as the value of --mesh-storage-class-name.  You can create PersistentVolume by the following definition:  apiVersion: v1 kind: PersistentVolume metadata:   labels:     app: easemesh   name: easemesh-pv spec:   storageClassName: {easemesh-storage}   accessModes:   - {ReadWriteOnce}   capacity:     storage: {3Gi}   hostPath:     path: {/opt/easemesh/}     type: "DirectoryOrCreate" |             |
| --mesh-control-plane-service-admin-port int     |           | Port of Easegress admin address (default 2381)                                                                                                                                                                                                                                                                                                                                                                                                                                                                                             |             |
| --mesh-control-plane-service-name string        |           | Mesh control plane service name (default "easemesh-controlplane-svc")                                                                                                                                                                                                                                                                                                                                                                                                                                                                      |             |
| --mesh-control-plane-service-peer-port int      |           | Port of Easegress cluster peer (default 2380)                                                                                                                                                                                                                                                                                                                                                                                                                                                                                              |             |
| --mesh-ingress-service-port int32               |           | Port of mesh ingress controller (default 19527)                                                                                                                                                                                                                                                                                                                                                                                                                                                                                            |             |
| --mesh-namespace string                         |           | EaseMesh namespace in kubernetes (default "easemesh")                                                                                                                                                                                                                                                                                                                                                                                                                                                                                      |             |
| --mesh-storage-class-name string                |           | Mesh storage class name (default "easemesh-storage")                                                                                                                                                                                                                                                                                                                                                                                                                                                                                       |             |
| --registry-type string                          |           | The registry type for application service registry, support eureka, consul, nacos (default "eureka")                                                                                                                                                                                                                                                                                                                                                                                                                                       |             |

## emctl reset

Reset infrastructure components of the EaseMesh

```bash
emctl reset [flags]

# Examples
emctl reset --mesh-namespace mesh-demo
```

| Flags                                    | Shorthand | Description                                                           |
| ---------------------------------------- | --------- | --------------------------------------------------------------------- |
| --help                                   | -h        | help for reset                                                        |
| --mesh-control-plane-service-name string |           | Mesh control plane service name (default "easemesh-controlplane-svc") |
| --mesh-namespace string                  |           | EaseMesh namespace in kubernetes (default "easemesh")                 |

## emctl apply

Apply a configuration to easemesh.

```bash
emctl apply [flags]

# Examples
emctl apply -f config.yaml
```

| Flags              | Shorthand | Description                                                                                                 |
| ------------------ | --------- | ----------------------------------------------------------------------------------------------------------- |
| --file string      | -f        | A location contained the EaseMesh resource files (YAML format) to apply, could be a file, directory, or URL |
| --help             | -h        | help for apply                                                                                              |
| --recursive        | -r        | Whether to recursively iterate all sub-directories and files of the location (default true)                 |
| --server string    | -s        | An address to access the EaseMesh control plane (default "127.0.0.1:2381")                                  |
| --timeout duration | -t        | A duration that limit max time out for requesting the EaseMesh control plane (default 30s)                  |

## emctl get

Get resources of easemesh.

```bash
emctl get [flags]

# Examples
emctl get -f config.yaml
emctl get service service-001
```

| Flags              | Shorthand | Description                                                                                |
| ------------------ | --------- | ------------------------------------------------------------------------------------------ |
| --help             | -h        | help for get                                                                               |
| --output string    | -o        | Output format (support table, yaml, json) (default "table")                                |
| --server string    | -r        | An address to access the EaseMesh control plane (default "127.0.0.1:2381")                 |
| --timeout duration | -t        | A duration that limit max time out for requesting the EaseMesh control plane (default 30s) |

## emctl delete

Delete resources of easemesh.

```bash
emctl delete [flags]

# Examples
emctl delete -f config.yaml
emctl delete service service-001
```

| Flags              | Shorthand | Description                                                                                                 |
| ------------------ | --------- | ----------------------------------------------------------------------------------------------------------- |
| --file string      | -f        | A location contained the EaseMesh resource files (YAML format) to apply, could be a file, directory, or URL |
| --help             | -h        | help for delete                                                                                             |
| --recursive        | -r        | Whether to recursively iterate all sub-directories and files of the location (default true)                 |
| --server string    | -s        | An address to access the EaseMesh control plane (default "127.0.0.1:2381")                                  |
| --timeout duration | -t        | A duration that limit max time out for requesting the EaseMesh control plane (default 30s)                  |

## Cheatsheet

```bash
# Install EaseMesh Components
emctl install --clean-when-failed

# Apply Tenant
echo 'apiVersion: mesh.megaease.com/v1alpha1                                                        [~/code/easemesh/ctl]
kind: Tenant
metadata:
  name: tenant-001
  labels: {}
services: []
description: tenant-001' | emctl apply -f -

# Apply Tenant
emctl apply -f tenant-001.yaml

# Apply Service
emctl apply -f service-001.yaml

# Apply Service
echo 'kind: Service
metadata:
  name: service-001
spec:
  registerTenant: "tenant-001"
  sidecar: {}' | emctl apply -f -

# Apply LoadBalance
echo 'apiVersion: mesh.megaease.com/v1alpha1                                                        [~/code/easemesh/ctl]
kind: LoadBalance
metadata:
  name: service-001
spec:
  policy: random' | emctl apply -f -

# Apply Ingress
echo 'apiVersion: mesh.megaease.com/v1alpha1
kind: Ingress
metadata:
  name: service-001
  labels: {}
spec:
  rules:
  - paths:
    - path: .*
      backend: service-001' | emctl apply -f -

# Get Tenant (kind is case-insensitive in command line)
emctl get tenant -o yaml

# Get service
emctl get service
emctl get service -o yaml
emctl get service service-001 -o json

# Get LoadBalance
emctl get loadbalance
emctl get loadbalance service-001 -o yaml

# Delete service
emctl delete service service-001
emctl delete service -f service-001.yaml

# Delete LoadBalance
emctl delete loadbalance service-001

# NOTE: The manipulation of the kinds attached to Service below is the same with LoadBalance:
# - Sidecar
# - Resilience
# - Canary
# - ObservabilityMetrics, ObservabilityTracings, ObservabilityOutputServer
```
