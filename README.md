
# EaseMesh

A service mesh compatible with the Spring Cloud ecosystem. Using [Easegress](https://github.com/megaease/easegress) as a sidecar for service management & [EaseAgent](https://github.com/megaease/easeagent) as a monitor for service observability.

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
    - [7.2 Canary Deployment](#72-canary-deployment)
    - [7.3 Clean](#73-clean)
  - [8. Roadmap](#8-roadmap)
  - [9. License](#9-license)

## 1. Purposes

Why do we reinvent another wheel?

- **Service mesh compatible with Spring Cloud ecosystem:** The micro-services developed in Spring Cloud ecosystem have their own service registry/discovery system, this is quite different with Kubernetes ecosystem which uses the DNS as the service discovery. Currently, the major Service Mesh solution (e.g. Istio) using the Kubernetes domain technology. So, this is painful and conflicted with Java Spring Cloud domain. EaseMesh aims to make Service Mesh compatible with Java Spring Cloud completely.

- **Integrated Observability:** Currently Kubernetes-based service mesh only can see the ingress/egress traffic, and it has no idea what's happened in service/application. So, combining with Java Agent technology, we can have the full capability to observe everything inside and outside of service/application.

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
cd ctl && make
export PATH=$(pwd)/bin:${PATH}
```

### 7.1 Start PetClinic in EaseMesh

1. Apply mesh configuration files

```bash
emctl apply -f https://raw.githubusercontent.com/megaease/easemesh-spring-petclinic/main/mesh-conf/a-pet-tenant.yaml
emctl apply -f https://raw.githubusercontent.com/megaease/easemesh-spring-petclinic/main/mesh-conf/api-gateway.yaml
emctl apply -f https://raw.githubusercontent.com/megaease/easemesh-spring-petclinic/main/mesh-conf/customers.yaml
emctl apply -f https://raw.githubusercontent.com/megaease/easemesh-spring-petclinic/main/mesh-conf/ingress.yaml
emctl apply -f https://raw.githubusercontent.com/megaease/easemesh-spring-petclinic/main/mesh-conf/vets.yaml
emctl apply -f https://raw.githubusercontent.com/megaease/easemesh-spring-petclinic/main/mesh-conf/visits.yaml
```

2. Create namespace: `kubectl create namespace spring-petclinic`
3. Apply petclinic stack

```bash
kubectl apply -f https://raw.githubusercontent.com/megaease/easemesh-spring-petclinic/main/mesh-deployments/api-gateway-deployment.yaml
kubectl apply -f https://raw.githubusercontent.com/megaease/easemesh-spring-petclinic/main/mesh-deployments/customers-service-deployment.yaml
kubectl apply -f https://raw.githubusercontent.com/megaease/easemesh-spring-petclinic/main/mesh-deployments/vets-service-deployment.yaml
kubectl apply -f https://raw.githubusercontent.com/megaease/easemesh-spring-petclinic/main/mesh-deployments/visits-service-deployment.yaml
```

4. Use the DB table schemes and records from [PetClinic example](https://github.com/spring-projects/spring-petclinic/tree/main/src/main/resources/db/mysql) to set up yours.
5. Run `kubectl get service -n easemesh easemesh-ingress-service` , then configure the NodPort IP address and port number into your traffic gateway's routing address, e.g, add config to NGINX:

```plain
location /pet/ {
    proxy_pass http://$NodePortIP:$NodePortNum/;
}
```

**Note:**  that the PetClinic website should be routed by the  `/`  subpath, or it should use  `NGINX`'s replacing response content feature for correcting resource URL:

```plain
location /pet/ {
    proxy_pass http://$NodePortIP:$NodePortNum/;
    sub_filter 'href="/' 'href="/pet/';
	sub_filter 'src="/' 'src="/pet/';
	sub_filter_once  off;
}
```

6. Visiting PetClinic website with `$your_domain/pet/#!/welcome`

### 7.2 Canary Deployment

![EaseMesh Canary topology](./imgs/canary-deployment.png)

1. Coloring traffic with HTTP header `X-Canary: lv1` by using Chrome browser's **ModHeader** plugin. Then EaseMesh will route this colored traffic into the Customer service's canary version instance.

2. Apply mesh configuration file: `emctl apply -f https://raw.githubusercontent.com/megaease/easemesh-spring-petclinic/main/canary/customer-canary.yaml`

3. Developing a canary version of Customer service to add an extra suffix to the city field for each record.

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

4. Building the canary Customer service's image, and update image address in `https://github.com/megaease/easemesh-spring-petclinic/blob/main/canary/customers-service-deployment-canary.yaml`. Or just use our default canary image which already was in it.

5. Apply canary deployment: `kubectl apply -f https://raw.githubusercontent.com/megaease/easemesh-spring-petclinic/main/canary/customers-service-deployment-canary.yaml`

6. Turning on the chrome **ModHeader** plugin to color the traffic, then visit PetClinic website. You can see the change to the table which adds an "-US" suffix to every city record.

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

## 9. License

EaseMesh is under the Apache 2.0 license. See the [LICENSE](./LICENSE) file for details.
