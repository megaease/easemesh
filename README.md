# EaseMesh

EaseMesh is a service mesh that is compatible with the Spring Cloud ecosystem. It is based on [Easegress](https://github.com/megaease/easegress) for the sidecar of service management and [EaseAgent](https://github.com/megaease/easeagent) for the monitor of service observing.

<a href="https://megaease.com/easemesh">
    <img src="./imgs/easemesh.svg"
        alt="EaseMesh logo" title="EaseMesh" height="175" width="175" align="right"/>
</a>

- [EaseMesh](#easemesh)
  - [1. Purposes](#1-purposes)
  - [2. Principles](#2-principles)
  - [3. Architecture](#3-architecture)
  - [4. Features](#4-features)
  - [5. Dependent Projects](#5-dependent-projects)
  - [6. Quick Start](#6-quick-start)
    - [6.1 Environment Requirement](#61-environment-requirement)
    - [6.2 Sanity Checking](#62-sanity-checking)
    - [6.3 Installation](#63-installation)
  - [7. Demonstration](#7-demonstration)
    - [7.1 Start PetClinic in EaseMesh](#71-start-petclinic-in-easemesh)
      - [7.1.1 Step 1: Apply mesh configuration](#711-step-1-apply-mesh-configuration)
      - [7.1.2 Step 2: Create namespace](#712-step-2-create-namespace)
      - [7.1.3 Step 4: Setup Database](#713-step-4-setup-database)
      - [7.1.4 Step 3: Apply petclinic stack](#714-step-3-apply-petclinic-stack)
      - [7.1.5 Get exposed port of `EaseMesh ingress` service](#715-get-exposed-port-of-easemesh-ingress-service)
      - [7.1.6 Step 5: Configure reverse proxy](#716-step-5-configure-reverse-proxy)
        - [7.1.6.1 Config reverse proxy via Easegress](#7161-config-reverse-proxy-via-easegress)
        - [7.1.6.2 Config reverse proxy via Easegress](#7162-config-reverse-proxy-via-easegress)
    - [7.2 Canary Deployment](#72-canary-deployment)
      - [7.2.1  Step 1: Coloring traffic](#721--step-1-coloring-traffic)
      - [7.2.2 Step 2: Apply canary configuration of the EaseMesh](#722-step-2-apply-canary-configuration-of-the-easemesh)
      - [7.2.3 Step 3:  Prepare a canary version of the application](#723-step-3--prepare-a-canary-version-of-the-application)
      - [7.2.4 Step 4: Build canary image](#724-step-4-build-canary-image)
      - [7.2.5 Step 5. Deploy canary version](#725-step-5-deploy-canary-version)
      - [7.2.6 Step 6: Sending coloring traffic](#726-step-6-sending-coloring-traffic)
    - [7.3 Clean](#73-clean)
  - [8. Roadmap](#8-roadmap)
  - [9. Contributing](#9-contributing)
  - [10. License](#10-license)
  - [11. User Manual](#11-user-manual)

## 1. Purposes

Why do we reinvent another wheel?

**Service mesh compatible with Spring Cloud ecosystem:** Micro-service in Spring Cloud ecosystem has its own service registry/discovery components. It is quite different from Kubernetes ecosystem using DNS for service discovery. The major Service Mesh solution (e.g. Istio) using the Kubernetes domain technology. It is painful and conflicted with Java Spring Cloud ecosystem. EaseMesh aims to make Service Mesh compatible with Java Spring Cloud completely.

- **Integrated Observability:** Currently Kubernetes-based service mesh only can see the ingress/egress traffic, and it has no idea what's happened in service/application. So, combining with Java Agent technology, we can have the full capability to observe everything inside and outside of service/application.

- **Sophisticated capability of traffic split:**  The EaseMesh has the sophisticated capability of traffic split, it can split traffic of a request chain into not only first service but also last. The capability could be applied in the canary deployment, online production testing scenarios.

> Shortly, **the EaseMesh leverages the Kubernetes sidecar and Java Agent techniques to make Java applications have service governance and integrated observability without change a line of source code**.

## 2. Principles

- **Spring Cloud Compatibility:** Spring Cloud domain service management and resilient design.
- **No Code Changes:** Using sidecar & Java-agent for completed service governance and integrated observability.
- **Service Insight:** Service running metrics/tracing/logs monitoring.

## 3. Architecture

![The architecture diagram](./imgs/architecture.png)

## 4. Features

- **Non-intrusive Design**: Zero code modification for Java Spring Cloud application migration, only small configuration update needed.
- **Java Register/Discovery**: Compatible with popular Java Spring Cloud ecosystem's Service register/discovery（Eureka/Consul/Nacos).
- **Traffic Orchestration**: Coloring & Scheduling east-west and north-south traffic to configured services.
- **Resource Management**: Rely on Kubernetes platform for CPU/Memory resources management.
- **Canary Deployment**: Routing requests based on colored traffic and different versions of the service.
- **Resilience**: Including Timeout/CircuitBreaker/Retryer/Limiter, completely follow sophisticated resilience design.
- **Observability**: Including Metrics/Tracing/Log,e.g. HTTP Response code distribution, JVM GC counts, JDBC fully SQL sentences, Kafka/RabbitMQ/Redis metrics, open tracing records, access logs, and so on. With such abundant and services-oriented data, developers/operators can diagnosis where the true problems happened, and immediately take corresponding actions.

## 5. Dependent Projects

1. [MegaEase EaseAgent](https://github.com/megaease/easeagent)
2. [MegaEase Easegress](https://github.com/megaease/easegress)

## 6. Quick Start

### 6.1 Environment Requirement

- Linux kernel version 4.15+
- Kubernetes version 1.18+
- Mysql version 14.14+

### 6.2 Sanity Checking

- Running `kubectl get nodes` to check your Kubernetes cluster's healthy.

### 6.3 Installation

Please check out [install.md](./docs/install.md) to install EaseMesh.

## 7. Demonstration

- [Spring Cloud PetClinic](https://github.com/spring-petclinic/spring-petclinic-cloud) microservice example.

- It uses Spring Cloud Gateway, Spring Cloud Circuit Breaker, Spring Cloud Config, Spring Cloud Sleuth, Resilience4j, Micrometer and Eureka Service Discovery from Spring Cloud Netflix technology stack.

![The topology migration diagram](imgs/topology-migration.png)

Prepare the `emctl`

```bash
git clone https://github.com/megaease/easemesh
cd emctl && make
export PATH=$(pwd)/bin:${PATH}
```

### 7.1 Start PetClinic in EaseMesh

#### 7.1.1 Step 1: Apply mesh configuration

Apply the EaseMesh configuration files

```bash
emctl apply -f https://raw.githubusercontent.com/megaease/easemesh-spring-petclinic/main/mesh-conf/a-pet-tenant.yaml
emctl apply -f https://raw.githubusercontent.com/megaease/easemesh-spring-petclinic/main/mesh-conf/api-gateway.yaml
emctl apply -f https://raw.githubusercontent.com/megaease/easemesh-spring-petclinic/main/mesh-conf/customers.yaml
emctl apply -f https://raw.githubusercontent.com/megaease/easemesh-spring-petclinic/main/mesh-conf/ingress.yaml
emctl apply -f https://raw.githubusercontent.com/megaease/easemesh-spring-petclinic/main/mesh-conf/vets.yaml
emctl apply -f https://raw.githubusercontent.com/megaease/easemesh-spring-petclinic/main/mesh-conf/visits.yaml
```

#### 7.1.2 Step 2: Create namespace

leverage kubectl to create `spring-petclinic` namespace

```bash
kubectl create namespace spring-petclinic
```

#### 7.1.3 Step 4: Setup Database

Petclinic needs to access database, the default is memory database. But in the EaseMesh quick start, you need to prepare a Mysql database for the demo.

Use the DB table schemes and records from [PetClinic example](https://github.com/spring-projects/spring-petclinic/tree/main/src/main/resources/db/mysql) to set up yours.

#### 7.1.4 Step 3: Apply petclinic stack

Deploy petclinic resources to k8s cluster, we have developed an [operator](./operator/README.md) to manage the custom resource (MeshDeployment) of the EaseMesh. `Meshdeployment` contains a K8s' complete deployment spec and a piece of extra information about the service.

> The Operator of the EaseMesh will automatically inject a sidecar to pod and a JavaAgent into the application's JVM

```bash
kubectl apply -f https://raw.githubusercontent.com/megaease/easemesh-spring-petclinic/main/mesh-deployments/api-gateway-deployment.yaml
kubectl apply -f https://raw.githubusercontent.com/megaease/easemesh-spring-petclinic/main/mesh-deployments/customers-service-deployment.yaml
kubectl apply -f https://raw.githubusercontent.com/megaease/easemesh-spring-petclinic/main/mesh-deployments/vets-service-deployment.yaml
kubectl apply -f https://raw.githubusercontent.com/megaease/easemesh-spring-petclinic/main/mesh-deployments/visits-service-deployment.yaml
```

> ATTENTION: There is a ConfigMap spec in yaml file, it describes how to connected the database for applications. You need to change as per your real environment

#### 7.1.5 Get exposed port of `EaseMesh ingress` service

```bash
kubectl get service -n easemesh easemesh-ingress-service
```
> **ATTENTION**: From the output, you may notice exposed port of the Ingress service. If you needn't use reverse proxy service, you can directly access pet-clinic application via http://{your_host}:{exposed_port}/

#### 7.1.6 Step 5: Configure reverse proxy

> **ATTENTION**: The step is optional. It can be omitted when you have no requirements about reverse proxy.

##### 7.1.6.1 Config reverse proxy via Easegress

> **ATTENTION**: Only for scenarios that the Easegress acts as the role of reverse proxy service

If you leverage the [Easegress](https://github.com/megaease/easegress) as a reverse proxy service, the following configuration can be applied.

HTTP Server spec (file name: http-server.yaml):

```yaml
kind: HTTPServer
name: spring-petclinic-example
port: 443
https: true
keepAlive: true
keepAliveTimeout: 75s
maxConnection: 10240
cacheSize: 0
certs:
  key: {add your certs information to here}
rules:
  - paths:
    - pathPrefix: /
      backend: http-petclinic-pipeline
```

HTTP Pipeline spec (file name: http-petclinic-pipeline.yaml):

```yaml
name: http-petclinic-pipeline
kind: HTTPPipeline
flow:
  - filter: requestAdaptor
  - filter: proxy
filters:
  - name: requestAdaptor
    kind: RequestAdaptor
    method: ""
    path: null
    header:
      del: []
      set:
        Host: "{you host name, can be omitted}"
        X-Forwarded-Proto: "https"
        Connection: "upgrade"
      add:
        Host: "{you host name, can be omitted}"
  - name: proxy
    kind: Proxy
    mainPool:
      servers:
      - url: http://{node1_of_k8s_cluster}:{port_exposed_by_ingress_service}
      - url: http://{node2_of_k8s_cluster}:{port_exposed_by_ingress_service}
      loadBalance:
        policy: roundRobin
```

Change contents in `{}` as per your environment, and apply it via Easegress client command tool `egctl`:

```bash
egctl apply -f http-server.yaml
egctl apply -f http-petclinic-pipeline.yaml
```

> **egctl** is the client command line of the Easegress

Visiting PetClinic website with `$your_domain/#!/welcome`

##### 7.1.6.2 Config reverse proxy via Easegress

> **ATTENTION**: Only for scenarios that the Nginx acts as the role of reverse proxy service

if you leverage the Nginx as a reverse proxy service, the following configuration should be added.

Then configure the NodPort IP address and port number into your traffic gateway's routing address, e.g, add config to NGINX:

```plain
location /pet/ {
    proxy_pass http://{node1_of_k8s_cluster}:{port_exposed_by_ingress_service}/;
}
```

> **ATTENTION:**  that the PetClinic website should be routed by the  `/`  subpath, or it should use  `NGINX`'s replacing response content feature for correcting resource URL:

```plain
location /pet/ {
    proxy_pass http://{node1_of_k8s_cluster}:{port_exposed_by_ingress_service/;
    sub_filter 'href="/' 'href="/pet/';
    sub_filter 'src="/' 'src="/pet/';
    sub_filter_once  off;
}
```

Visiting PetClinic website with `$your_domain/pet/#!/welcome`.

### 7.2 Canary Deployment

Canary deployment demonstrates how to route coloring traffic (request) to a canary version of the specific service.

![EaseMesh Canary topology](./imgs/canary-deployment.png)

- `Customer Service (v2)` is the canary version service.
- The line of red color in the diagram represents coloring traffic (request).
- The coloring traffic is correctly routed into canary version service after it has passed through the first service (API Gateway).

#### 7.2.1  Step 1: Coloring traffic

Coloring traffic with HTTP header `X-Canary: lv1` by using Chrome browser's **[ModHeader](https://chrome.google.com/webstore/detail/modheader/idgpnmonknjnojddfkpgkljpfnnfcklj?hl=en)** plugin. Then EaseMesh will route this colored traffic into the Customer service's canary version instance.

#### 7.2.2 Step 2: Apply canary configuration of the EaseMesh

Apply mesh configuration file:

```bash
emctl apply -f https://raw.githubusercontent.com/megaease/easemesh-spring-petclinic/main/canary/customer-canary.yaml`
```

#### 7.2.3 Step 3:  Prepare a canary version of the application

> **ATTENTION**  You can skip the step, we have provides the canary image to docker hub `megaease/spring-petclinic-customers-service:canary` you can found it in the docker hub.

Developing a canary version of Customer service to add an extra suffix to the city field for each record.

```diff
diff --git a/spring-petclinic-customers-service/src/main/java/org/springframework/samples/petclinic/customers/model/Owner.java b/spring-petclinic-customers-service/src/main/java/org/springframework/samples/petclinic/customers/model/Owner.java
index 360e765..cc2df3d 100644
--- a/spring-petclinic-customers-service/src/main/java/org/springframework/samples/petclinic/customers/model/Owner.java
+++ b/spring-petclinic-customers-service/src/main/java/org/springframework/samples/petclinic/customers/model/Owner.java
@@ -99,7 +99,7 @@ public class Owner {
    }

    public String getAddress() {
-        return this.address;
+        return this.address + " - US";
    }

    public void setAddress(String address) {k
```

#### 7.2.4 Step 4: Build canary image

> **ATTENTION**  You can skip the step, we have provides the canary image to docker hub `megaease/spring-petclinic-customers-service:canary` you can found it in the docker hub.

Building the canary Customer service's image, and update image version in `https://github.com/megaease/easemesh-spring-petclinic/blob/main/canary/customers-service-deployment-canary.yaml`. Or just use our default canary image which already was in it.

#### 7.2.5 Step 5. Deploy canary version

Being similar to [7.1.4](#714-step-3-apply-petclinic-stack),  we leverage kubectl to deploy the canary version of `MeshDeployment`

```bash
kubectl apply -f https://raw.githubusercontent.com/megaease/easemesh-spring-petclinic/main/canary/customers-service-deployment-canary.yaml`
```

> **ATTENTION**: There is a ConfigMap spec in echo yaml spec, it describes how to connect the database for applications. You need to change its contents for your environment.

#### 7.2.6 Step 6: Sending coloring traffic

Turning on the chrome **ModHeader** plugin to color the traffic, then visit PetClinic website. You can see the change to the table which adds an "-US" suffix to every city record.

![plugin](./imgs/chrome_plugin.png)

### 7.3 Clean

- Run `kubectl delete namespace spring-petclinic`.
- Run

```bash
emctl delete ingress pet-ingress
emctl delete service api-gateway
emctl delete service customers-service
emctl delete service vets-service
emctl delete service visits-service
emctl delete tenant pet
```

## 8. Roadmap

See [EaseMesh Roadmap](./docs/Roadmap.md) for details.

## 9. Contributing

See [MegaEase Community](https://github.com/megaease/community) to follow our contributing details.

## 10. License

EaseMesh is under the Apache 2.0 license. See the [LICENSE](./LICENSE) file for details.

## 11. User Manual

See [EaseMesh User Manual](./docs/user_manual.md) for details.
