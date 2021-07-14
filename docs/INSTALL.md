# EaseMesh Installation

- [EaseMesh Installation](#easemesh-installation)
  - [Prerequisites](#prerequisites)
    - [Infrastructure components of the EaseMesh](#infrastructure-components-of-the-easemesh)
      - [Easegress](#easegress)
      - [EaseAgent](#easeagent)
      - [EaseMesh Operator](#easemesh-operator)
      - [EaseMesh command line tool - emctl](#easemesh-command-line-tool---emctl)
      - [Build the EaseMesh client tools from scratch](#build-the-easemesh-client-tools-from-scratch)
    - [Environments](#environments)
      - [K8s and Connectivity](#k8s-and-connectivity)
      - [Persistent Volume](#persistent-volume)
  - [Installation](#installation)
    - [Install EaseMesh](#install-easemesh)
    - [Reset environment](#reset-environment)
  - [Trouble Shooting](#trouble-shooting)

This document gives the instructions to install all infrastructure components (except K8s) required by the EaseMesh.

## Prerequisites


### Infrastructure components of the EaseMesh

The dependencies of the EaseMesh are the Easegress, EaseAgent, and EaseMesh Operators.

- Easegress
- EaseAgent
- EaseMesh Operator

The EaseMesh severely depends on the K8s, all dependencies must be packaged as images, so you need to prepare all three docker images.

#### Easegress

The Easegress plays multiple roles in the EaseMesh, including:

- Control plane manages services' registry, configurations for the Mesh.
- Sidecar takes over all traffic in and out of the application container. 
- Ingress takes over all traffic entered the cluster of the K8s.

You can found Easegress from [here](https://github.com/megaease/easegress/releases). The latest image has been uploaded by us. You can download it by `docker pull`.

```
docker pull megaease/easegress
```

> If you want to build Easegress from scratch, you cant refer to [here](https://github.com/megaease/easegress/blob/main/README.md#setting-up-easegress)

#### EaseAgent

The EaseAgent is a JavaAgent whose responsibility is collecting metrics and tracing information. It provides observability for services in the Java ecosystems. In a typical scenario, the EaseAgent is a jar package, you can download it from [release](https://github.com/megaease/easeagent/releases). As the EaseMesh relies on the K8s, we provide a dedicated image for automatically injecting the EaseAgent Jar file from the docker image.

The latest image has been uploaded by us. You can download it by `docker pull`.

```
docker pull megaease/easeagent-initializer
```

#### EaseMesh Operator

EaseMesh Operator is the CRD Operator whose responsibility is when a MeshDeployment (custom resource) is deployed, it renders a deployment spec from the custom resource and injects JavaAgent to service application and a sidecar container to take over all traffics. EaseMesh Operator follow the [Operator pattern](https://kubernetes.io/docs/concepts/extend-kubernetes/operator/).
EaseMesh Operator implement a control loop which repeatedly compare the desired state (deployment resources which are injected with a sidecar container) of the cluster to its actual state.

For convenience, we provide the EaseMesh Operator docker image in docker hub. you can download image via `docker pull`.

```
docker pull megaease/easemesh-operator
```

> You can build easemesh operator image from scratch, please refer [here](https://github.com/megaease/easemesh/tree/main/operator#how-to-build-it)

Ensuring you have all three images (accessing docker hub without any problems), we will begin our journey.

#### EaseMesh command line tool - emctl

#### Build the EaseMesh client tools from scratch

The `emctl` can be built from source code, you may follow this guide to build it. 


> The emctl is implemented in Golang language, you need prepare golang (1.16+) dev environment in advance


1. clone source code
```bash
git clone https://github.com/megaease/easemesh
```

2. build emctl executables

```bash
cd easemesh/ctl && make
```

3. If no errors occurred, the target was built in `bin/` directory, named with `emctl`


### Environments

#### K8s and Connectivity

EaseMesh severely depends on the K8s, in order to install the EaseMesh, you must confirm you have a healthy K8s cluster and sufficient resources (at least three work nodes by default).

The installation of the EaseMesh needs admin privilege of K8s cluster, the  `emctl` looks for a file named config in the $HOME/.kube directory. The emctl use it to communicate with K8s API server of a cluster.

By default, The control plane of the EaseMesh exposes its service via K8s' [NodePort](https://kubernetes.io/docs/concepts/services-networking/service/#nodeport), so you should ensure that the node running `emctl` can access nodes of the K8s cluster 

#### Persistent Volume

We deploy the control plane of the EaseMesh in a K8s cluster, as the control plane needs to persistent configuration, the persistent volume resource needs to be introduced.

The default replicas of the control plane is three, so three PVs is required. The capacity of each PV must greater than 3Gi by default. We provide a template spec of PV here for your referring:

```yaml
apiVersion: v1
kind: PersistentVolume
metadata:
  name: easemesh-storage-pv-1
spec:
  capacity:
    storage: 4Gi
  volumeMode: Filesystem
  accessModes:
  - ReadWriteOnce
  persistentVolumeReclaimPolicy: Delete
  storageClassName: easemesh-storage
  local:
    path: {{specific_path_you_need_to_substituted}}
  nodeAffinity:
    required:
      nodeSelectorTerms:
      - matchExpressions:
        - key: kubernetes.io/hostname
          operator: In
          values:
          - {{node_name_with_specific_path}}
---
apiVersion: v1
kind: PersistentVolume
metadata:
  name: easemesh-storage-pv-2
spec:
  capacity:
    storage: 4Gi
  volumeMode: Filesystem
  accessModes:
  - ReadWriteOnce
  persistentVolumeReclaimPolicy: Delete
  storageClassName: easemesh-storage
  local:
    path: {{specific_path_you_need_to_substituted}}
  nodeAffinity:
    required:
      nodeSelectorTerms:
      - matchExpressions:
        - key: kubernetes.io/hostname
          operator: In
          values:
          - {{node_name_with_specific_path}}
---
apiVersion: v1
kind: PersistentVolume
metadata:
  name: easemesh-storage-pv-3
spec:
  capacity:
    storage: 4Gi
  volumeMode: Filesystem
  accessModes:
  - ReadWriteOnce
  persistentVolumeReclaimPolicy: Delete
  storageClassName: easemesh-storage
  local:
    path: {{specific_path_you_need_to_substituted}}
  nodeAffinity:
    required:
      nodeSelectorTerms:
      - matchExpressions:
        - key: kubernetes.io/hostname
          operator: In
          values:
          - {{node_name_with_specific_path}}
---
apiVersion: storage.k8s.io/v1
kind: StorageClass
metadata:
  name: easemesh-storage
provisioner: kubernetes.io/no-provisioner
reclaimPolicy: Delete
volumeBindingMode: WaitForFirstConsumer

```

Changing contents according to your environment, provisioning it to your K8s cluster. You must ensuring all PVs created normally.


> We leverage [local volume](https://kubernetes.io/docs/concepts/storage/volumes/#local) to persistent control plane data.

## Installation

### Install EaseMesh


If all prerequisites are fulfilled properly, we can begin to install the EaseMesh. Although there are servals steps that setups required components, it can be done by one command. Steps are:

- **Setup the control plane of the EaseMesh**: EaseMesh leverage the Easegress implementing control plane, the installation will deploy Easegress as `statefulset` resource of the K8s
- **Provision mesh controller**: Mesh controller is implemented in the Easegress, but it is disabled by default. Installation needs to apply a configuration to enable the Mesh controller.
- **Create the `Custom Resource Definition`**: EaseMesh leverages custom resource definition to manage an application. Installation needs to create it in advance.
- **Provision the Operator**: The Operator is used to reconciling custom resources to deployment resources.
- **Provision Mesh ingress**: Mesh ingress is used to take over traffics entering in K8s cluster 

The installation will run all steps one by one. 

Install command is:

```bash
emctl install
```

If you want to speed up your installation, you can tag all three images and uploaded them into your local private docker registry. Specific private docker registry to install, just simply add an extra argument.

```bash
emctl install --image-registry-url {your_private_docker_registry_address}
```

more arguments can be discovered via :

```bash
emctl install --help
```


### Reset environment

if you want to remove the EaseMesh, just run the command:

```
emctl reset
```

> PVC and PV resources will not be reclaimed by default. you need to delete them manually.


## Trouble Shooting

There are a known issue, during the EaseMesh installation, the embed `etcd` of the Easegress can't be setup correctly. If you encounter this situation, just delete all persistent contents in nodes
