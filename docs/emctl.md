# EaseMesh Command-Line

`emctl` is the dedicated command to handle resources of EaseMesh, which running in [Easegress](https://github.com/megaease/easegress) as `MeshController` who has different roles in different instances. `MeshController` will register its own admin API in `Easegress`, so the server flag in `emctl` keeps the same as Easegress's.

Running `emctl --help`  can get details about every subcommand. Here's the cheat sheet to show the common usage of it.

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

# Get service.
emctl get service
emctl get service -o yaml
emctl get service service-001 -o json

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
