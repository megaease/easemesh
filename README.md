# EaseMesh
A service mesh compatible with the Spring Cloud ecosystem. Using [EaseGateway](https://github.com/megaease/easegateway) as a sidecar for service management & [EaseAgent](https://github.com/megaease/easeagent) as a monitor for service observability.
- [EaseMesh](#easemesh)
  - [Overview](#overview)
    - [Purposes and Principles](#purposes-and-principles)
      - [Purposes](#purposes)
      - [Principles](#principles)
    - [Architecture Diagram](#architecture-diagram)
    - [Features](#features)
    - [Dependent Projects](#dependent-projects)
  - [Quick Start](#quick-start)
    - [Environment Requirement](#environment-requirement)
      - [Infrastructure Version](#infrastructure-version)
      - [Dependence Checking](#dependence-checking)
    - [Installation](#installation)
    - [Examples](#examples)
      - [Overview](#overview-1)
        - [Start PetClinic in EaseMesh](#start-petclinic-in-easemesh)
      - [Canary Deployment](#canary-deployment)
      - [Undeploy](#undeploy)
  - [License](#license)

## Overview 
### Purposes and Principles
#### Purposes
* **Service mesh compatible with Spring Cloud ecosystem:** The applications developed in Spring Cloud ecosystem migrates to Service Mesh is quite painful. There are conflicting between Kubernetes-based Service Mesh and Spring Cloud in the service discovery, resilience design, etc. EaseMesh aims to solve this problem.
* **Microservices governance enhancement:** Although the microservices architecture makes services loosely coupled, helps CI/CD, it also complicates the system maintenance. EaseMesh can provide rich profiling data and problem-solving features for better microservices governance.

#### Principles
* **Spring Cloud Compatibility: ** Spring Cloud domain service management and resilient desgin.
* **No Code Changes:** Using sidecar & Java-agent for completed service governance and integrated observability.
* **Service Insight:** Service running metrics/tracing/logs monitoring. 
 

### Architecture Diagram
* ![The architecture diagram](/imgs/architecture.png)
### Features
* **Non-intrusive Design**: Zero code modification for Java Spring Cloud application migration, only small configuration update needed.
* **Java Register/Discovery**: Compatible with popular Java Spring Cloud ecosystem's Service register/discovery（Eureka/Consul/Nacos). 
* **Traffic Orchestration**: Coloring & Scheduling east-west and north-south traffic to configured services. 
* **Resource Management**: Rely on Kubernetes platform for CPU/Memory resources management. 
* **Canary Deployment**: Routing requests based on colored traffic and different versions of the service.
* **Resilience**: Including Timeout/CircuitBreaker/Retryer/Limiter, completely follow sophisticated resilience design.
* **Observability**: Including Metrics/Tracing/Log,e.g. HTTP Response code distribution, JVM GC counts, JDBC fully SQL sentences, Kafka/RabbitMQ/Redis metrics, open tracing records, access logs, and so on. With such abundant and services-oriented data, developers/operators can diagnosis where the true problems happened, and immediately take corresponding actions.
### Dependent Projects
1. [MegaEase EaseAgent](https://github.com/megaease/easeagen) 
2. [MegaEase EaseGateway](https://github.com/megaease/easegateway) 

## Quick Start
### Environment Requirement 
#### Infrastructure Version
* Linux kernel version 4.15+
* Kubernetes version 1.18+
* Mysql version 14.14+
####  Dependence Checking
1. Running `kubectl get nodes` to check your Kubernetes cluster's healthy. 
2. Running  `mysql -u$your_db_user -p$your_db_pass` to check the connection to your DB. 

### Installation
1. Registering K8s mesh-deployment CRD, and starting EaseMesh control-plane, IngressGateway with commands below:
```shell
$ cd ./install
$ ./egctl mesh install
```
**Note:** EaseMesh installation needs EaseGateway and EaseAgent's image. They are provided in Docker Hub. If you want to get them from your private image repository, run `./egctl mesh install --image-registry-url ${your_image-registry-url}` instead. 

2. Checking control plane and ingress gateway's status 
```shell
$ kubectl get pod mesh-ingress-${random-suffix}   
NAME              READY  STATUS  RESTARTS  AGE
mesh-ingress-${random-suffix}  1/1   Running  0     18h

$ kubectl get pod easegateway-cluster-0-${random-suffix}
NAME              READY  STATUS  RESTARTS  AGE
easegateway-cluster-0-${random-suffix}  1/1   Running  0     18h

$ kubectl get pod easegateway-cluster-1-${random-suffix}
NAME              READY  STATUS  RESTARTS  AGE
easegateway-cluster-1-${random-suffix}  1/1   Running  0     18h

$ kubectl get pod easegateway-cluster-2-${random-suffix}
NAME              READY  STATUS  RESTARTS  AGE
easegateway-cluster-2-${random-suffix}  1/1   Running  0     18h
```
3. Verifying the EaseMesh operator
```shell
$ kubectl get crd | grep meshdeployment              
meshdeployments.mesh.megaease.com       2021-03-18T02:54:15Z
```
### Examples 
#### Overview
*  [Spring Cloud PetClinic](https://github.com/spring-petclinic/spring-petclinic-cloud) microservice example.
* It uses Spring Cloud Gateway, Spring Cloud Circuit Breaker, Spring Cloud Config, Spring Cloud Sleuth, Resilience4j, Micrometer and Eureka Service Discovery from Spring Cloud Netflix technology stack.

![The topology migration diagram](imgs/topology-migration.png)


##### Start PetClinic in EaseMesh
1. Running  `./example/mesh-app-petclinic/deploy.sh`. 
2. Using the DB table schemes and records from [PetClinic example](https://github.com/spring-projects/spring-petclinic/tree/main/src/main/resources/db/mysql) to set up yours.
3. Running `kubectl get svc mesh-ingress` , then configure the NodPort IP address and port number into your traffic gateway's routing address,e.g., add config to NGINX with
```
location /pet/ {
        proxy_pass http://$NodePortIP:$NodePortNum/;
            ...
}

```
4. Visiting PetClinic website with `$your_domain/pet/#!/welcome` 

#### Canary Deployment
![EaseMesh Canary topology](./imgs/canary-deployment.png)
1. Coloring traffic with HTTP header `X-Canary: lv1` by using Chrome browser's **ModHeader** plugin. Then EaseMesh will route this colored traffic into the Customer service's canary version instance. 
2. Developing a canary version of Customer service to add an extra suffix to the city field for each record. 
```
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
3. Building the canary Customer service's image, and update it into `./example/mesh-app-petclinic/canary/customers-service-deployment-canary.yaml` file's line [#L22](https://github.com/megaease/easemesh/blob/main/example/mesh-app-petclinic/canary/customers-service-deployment-canary.yaml#L22). 
4. Running `kubectl apply -f  ./example/mesh-app-petclinic/canary/customers-service-deployment-canary.yaml`
5. Turning on the chrome **ModHeader** plugin to color the traffic, then visit PetClinic website. You can see the change to the table which adds an "-US" suffix to every city record. 
![plugin](./imgs/chrome_plugin.png)
#### Undeploy
* Running `./example/mesh-app-petclinic/undeploy.sh`.

## License
EaseMesh is under the Apache 2.0 license. See the [LICENSE](./LICENSE) file for details.
