### 1.Introduction  
* EaseMesh supports multiple micro service governance features such as traffic management, resilience features and observability. 

#### 1.1 Architecture
* The EaseMesh architecture is divided into two components. There are Control plane and Data plane. 
* The Control Plane's responsibility is to manage and monitor all services inside mesh, and accept user's declaring specs through CLT or console. 
* The Data Plane is composed by the Sidecar and JavaAgent in every deployed service Kubernetes Pod. They combined together to host traffic and gather the produced metrics/logs/tracings transparently from user's business service. The EaseMesh accepts out-side-mesh traffic by providing an Mesh-ingress which is an Easegress Node too. 
* The EaseMesh runs above Kubernetes which is the most popular and powerful orchestration platform currently, so that it can focus on handling Service Mesh about features. 

* ![The architecture diagram](../imgs/architecture.png)

#### 1.2 Reply components
* [The Easegress](https://github.com/megaease/easegress). It's the all-rounder gateway system to have built-in distributed storage, traffic scheduling, high performance and observability. 
* [The EaseAgent](https://github.com/megaease/easeagent). It's an APM tool under the Java system, used in a distributed system developed by Java. It provides cross-service call chain tracking and performance information collection for distributed systems.
* [The EaseMonitor](https://github.com/megaease/easeservice-mgmt-monitor).  It's the monitor component of the service governance. It supports dashboard and plane config data, queries and aggregates multiple types metrics data, queries service trace aggregation and topology analysis. 
  
#### 1.3 Control Plane 
1. In order to provide high-available ServiceMesh Control Plane in distributed environment, it requires odd numbers of Easegress nodes to form a cluster. The user can deliver their service mesh requirements and obtain running status and information inside the mesh by using RESTful-API or CLI through the Control Plane.
2. Each Easegress nodes synchronize configuration to achieve final-consistency with the help of RAFT algorithm. 
3. Every Easegress nodes run an MeshController for mesh-related logic. 

#### 1.4 Data plane
1. The Sidecar, is also an Easegress node in every service's Pod. It stands by the user's business service which is implemented with Spring-cloud framework. The Sidecar watches the service's specification modifications and applies them locally. The user's business service will accept Ingress traffic and deliver Egress traffic through its Sidecar.  
2. The JavaAgent, is a no invasion, service based view and high performance solution to enhance service governance's observability features in Java domain. It collects `JDBC`,`HTTP Servlet`, `HTTP filter`, Spring Boot 2.2.x: `WebClient 、 RestTemplate、FeignClient`, RabbitMQ and Jedis' metrics. It also supports collecting access-log and tracing recording.

#### 1.5 Kubernetes Deployment
1. EaseMesh uses Kubernetes CRD(customer resource define) to create a customized Mesh-needed Kubernetes Deployment. 
2. This EaseMesh CRD can automatically inject Sidecar and add JavaAgent using command into environment variable into user's original Kubernetes deployment.

### 1.6. Install 

##### 1.6.1 Prerequisites
Before you begin, check the following prerequisites:

1. Deploy kubernetes cluster with 1.18+.
2. Download the [EaseMesh release](https://github.com/megaease/easemesh/releases). 


##### 1.6.2 Install EaseMesh with egctl
You can install the EaseMesh using the following command:

```bash
$ cd easemesh/install
$ bin/egctl mesh install
Easegress control plane deploy success. Waiting startup...
Easegress control plane startup success.
Starting mesh controller success.
Easegress Ingress deploy success.
EaseMesh Operator deploy success.
Done.

```

This command deploys the components with default configuration. You can pass parameters through the command line to modify the configuration. Like following command:

```bash
$ bin/egctl mesh install --mesh-namespace=mynamespace
```

You can get all the configurations through ``-h``:

```textmate
$ bin/egctl mesh install -h

Deploy EaseMesh Components

Usage:
  egctl mesh install [flags]

Examples:
egctl mesh install <args>

Flags:
      --easegress-control-plane-replicas int    (default 3)
      --easegress-image string                  (default "megaease/easegress:latest")
      --easegress-ingress-replicas int          (default 3)
      --easemesh-operator-image string            (default "megaease/easemesh:latest")
      --easemesh-operator-replicas int            (default 3)
      --eg-admin-port int                         (default 2381)
      --eg-client-port int                        (default 2379)
      --eg-control-plane-pv-capacity int         The capacity of the PersistentVolume for easegress control plane storage, the unit is Gib. (default 3)
      --eg-control-plane-pv-hostpath string      The host path of the PersistentVolume for easegress control plane storage. (default "/opt/easegress")
      --eg-control-plane-pv-name string          The name of PersistentVolume for easegress control plane storage. (default "easegress-control-plane-pv")
      --eg-peer-port int                          (default 2380)
      --eg-service-admin-port int                 (default 2381)
      --eg-service-name string                    (default "easegress-public")
      --eg-service-peer-port int                  (default 2380)
  -f, --file string                              A yaml file specifying the install params.
      --heartbeat-interval int                    (default 5)
  -h, --help                                     help for install
      --image-registry-url string                 (default "docker.io")
      --mesh-namespace string                     (default "easemesh")
      --registry-type string                      (default "eureka")

Global Flags:
  -o, --output string   Output format(json, yaml) (default "yaml")
      --server string   The address of the Easegress endpoint (default "localhost:2381")

```
- Parameters Description 

| ParameterName                    | type   | description                                                                                                                                                         |
| -------------------------------- | ------ | ------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| image-registry-url               | string | Docker Image registry address, EaseMesh use it to pull Easegress and EaseMeshOperator Image. You can replace it with your private registry. Default: ``docker.io``. |
| easegress-image                  | string | The Easegress Image Name. Default: ``megaease/easegress:latest``.                                                                                                   |
| easemesh-operator-image          | string | The EaseMeshOperator Image. Default: ``megaease/easemesh:latest``.                                                                                                  |
| mesh-namespace                   | string | The Kubernetes Namespace for deployment of EaseMesh. Default: ``easemesh``.                                                                                         |
| easegress-control-plane-replicas | int    | The replicas of Easegress Control Plane's statefulset. Default: ``3``.                                                                                              |
| easegress-ingress-replicas       | int    | The replicas of Easegress Ingress's deployment. Default: ``easemesh``.                                                                                              |
| easemesh-operator-replicas       | int    | The replicas of EaseMesh Operator's deployment.  Default: ``easemesh``.                                                                                             |
| eg-client-port                   | int    | Port of Easegress Control Plane listen on for client traffic. Default: ``2379``.                                                                                    |
| eg-peer-port                     | int    | Port of Easegress Control Plane listen on for peer traffic. Default: ``2380``.                                                                                      |
| eg-admin-port                    | int    | Port of Easegress Control Plane listen on for admin traffic. Default: ``2381``.                                                                                     |
| eg-service-name                  | string | The Kubernetes service for Easegress control plane pods. Default: ``easegress-public``                                                                              |
| eg-service-client-port           | int    | Port of the service for pods's client port. Default: ``2379``.                                                                                                      |  |
| eg-service-peer-port             | int    | Port of the service for pods's peer port. Default: ``2380``.                                                                                                        |  |
| eg-service-admin-port            | int    | Port of the service for pods's admin port. Default: ``2381``.                                                                                                       |  |
| eg-control-plane-pv-name         | string | The PersistentVolume for Easegress Control Plane storage. Default: ``easegress-control-plane-pv``.                                                                  |
| eg-control-plane-pv-hostpath     | string | The path on host for Easegress Control Plane's PersistentVolume storage. Default: ``/opt/easegress``.                                                               |
| eg-control-plane-pv-capacity     | int    | The capacity of Easegress Control Plane's PersistentVolume, the unit is Gib.  Default: ``3``.                                                                       |
| registry-type                    | string | The registry type for application service registry. Default: ``eureka``                                                                                             |
| heartbeat-interval               | int    | The interval for checking service heartbeat, the unit is second. Default: ``5``                                                                                     |
| file                             | string | The config file for above parameters.                                                                                                                               |
    

##### 1.6.3 Check what’s installed      
The ``egctl mesh install`` command deploys the ``meshdeployments.mesh.megaease.com `` of CRD, the Easegress ControlPlane of StatefulSet and the PersistentVolume required for its storage, 
EasegressIngress and EaseMeshOperator of Deployment and required ConfigMap, Service, etc.

You can check all resources using following command:

```textmate
$ kubectl get crd | grep meshdeployment
meshdeployments.mesh.megaease.com                    

$ kubectl get ns | grep easemesh
easemesh          Active   33s

$ kubectl get statefulsets.apps -n easemesh
NAME                        READY   AGE
easegress-control-plane   3/3     33s

$ kubectl get deployments.apps -n easemesh
NAME                  READY   UP-TO-DATE   AVAILABLE   AGE
easegress-ingress   3/3     3            3           33s
easemesh-operator     3/3     3            3           33s
mesh-operator-hahha   3/3     3            3           33s

$ kubectl get pv
NAME                           CAPACITY   ACCESS MODES   RECLAIM POLICY   STATUS   CLAIM                                                               STORAGECLASS   REASON   AGE
easegress-control-plane-pv   3          RWO            Retain           Bound    easemesh/easegress-control-plane-pv-easegress-control-plane-0                           44h

$ kubectl get pods -n easemesh
NAME                                   READY   STATUS    RESTARTS   AGE
easegress-control-plane-0            1/1     Running   0          33s
easegress-control-plane-1            1/1     Running   0          33s
easegress-control-plane-2            1/1     Running   0          33s
easegress-ingress-847b7bddbb-9q7nf   1/1     Running   0          33s
easegress-ingress-847b7bddbb-px9n7   1/1     Running   0          33s
easegress-ingress-847b7bddbb-wj8ss   1/1     Running   0          33s
easemesh-operator-5fd5d55f8f-6d5bj     2/2     Running   0          33s
easemesh-operator-5fd5d55f8f-g89f9     2/2     Running   0          33s
easemesh-operator-5fd5d55f8f-j6ksk     2/2     Running   0          33s

$ kubectl get svc -n easemesh
NAME                                               TYPE        CLUSTER-IP       EXTERNAL-IP   PORT(S)                      AGE
easegress-hs                                     ClusterIP   None             <none>        2381/TCP,2380/TCP,2379/TCP   33s
easegress-ingress                                NodePort    10.106.94.98     <none>        13010:30010/TCP              33s
easegress-public                                 ClusterIP   10.104.166.129   <none>        2381/TCP,2380/TCP,2379/TCP   33s
mesh-operator-controller-manager-metrics-service   ClusterIP   10.97.62.250     <none>        8443/TCP                     33s
```

##### 1.6.4 Access Easegress Control Plane
You can use the ``ClusterIP:AdminPort`` of easegress-public service to access Easegress Control Plane by ``egctl``, like following command:

```bash
# Query Objects
$ bin/egctl object list --server 10.104.166.129:2381
- heartbeatInterval: 5s
  kind: MeshController
  name: easemesh-controller
  registryType: eureka

```

Now you can deploy application in EaseMesh according to the following document.

### 2.Deploy application in EaseMesh 
#### 2.1 Background
* EaseMesh can apply Java Spring Cloud application with only limited configuration modifications. No code modifications or recompiling needed. 
* EaseMesh treats **MeshService** as the first-class citizen. 
* EaseMesh supports multiple-tenant naturally. 

#### 2.2 Steps
1. Create a new Tenant with a configure file named "my_tenant.yaml" content below
```
name: ${your_tenant_name} 
services:
createdAt: 2021-04-19T18:00:00.00Z
description: "demo tenant"
```
2. Apply it with cmd `eashmesh/bin/meshctl tenant create -f ./my_tenant.yaml`

3. Check the tenant's creation by running cmd `eashmesh/bin/meshctl tenant get ${your_tenant_name}` 

3. Create a new application the configure file named "my_service.yaml" with content below

```
name: ${your_service_name} 
registerTenant: ${your_tenant_name} 
loadBalance:
  policy: random
  HeaderHashKey:
sidecar:
  discoveryType: eureka
  address: "127.0.0.1"
  ingressPort: 13001
  ingressProtocol: http
  egressPort: 13002
  egressProtocol: http
```
4. Apply it with cmd `eashmesh/bin/meshctl service create -f ./my_service.yaml` 

5. Check the service's creation by running cmd `eashmesh/bin/meshctl service get ${your_service_name}`

6. Prepare your application image, and put it into the  your application the configure file named "my_meshdeployment.yaml", here we prepare a Java Spring-cloud application using Eureka discovery center:  

```
apiVersion: v1
kind: ConfigMap
metadata:
  name: ${your_configmap} 
  namespace: ${your_ns} 
data:
  application-sit-yml: |
    server:
      port: 8080
    spring:
      application:
        name:  $(your_service_name) 
      datasource:
        url: jdbc:mysql://mysql.default:3306/meshappdemo?allowPublicKeyRetrieval=true&useUnicode=true&characterEncoding=utf-8&useSSL=false&serverTimezone=UTC&verifyServerCertificate=false
        username: ${your_username} 
        password: {$your_password} 
      jpa:
        database-platform: org.hibernate.dialect.MySQL5InnoDBDialect
      sleuth:
        enabled: false
        web:
          servlet:
          enabled: false
    eureka:
      client:
        serviceUrl:
          defaultZone: http://127.0.0.1:13009/mesh/eureka
      instance:
        preferIpAddress: true
        lease-expiration-duration-in-seconds: 60
---
apiVersion: mesh.megaease.com/v1beta1
kind: MeshDeployment
metadata:
  namespace: ${your_ns} 
  name: ${your_service} 
spec:
  service:
    name: ${your_service} 
  deploy:
    replicas: 2 
    selector:
      matchLabels:
        app: ${your_service} 
    template:
      metadata:
        labels:
          app: ${your_service} 
      spec:
        containers:
        - image: ${your_image_url} 
          name: ${your_service} 
          imagePullPolicy: IfNotPresent
          lifecycle:
            preStop:
              exec:
                command: ["sh", "-c", "sleep 10"]
          command: ["/bin/sh"]
          args: ["-c", "java -server -Xmx1024m -Xms1024m -Dspring.profiles.active=sit -Djava.security.egd=file:/dev/./urandom -jar /application/application.jar"]
          resources:
            limits:
              cpu: 2000m
              memory: 1Gi
            requests:
              cpu: 200m
              memory: 256Mi
          ports:
          - containerPort: 8080
          volumeMounts:
          - mountPath: /application/application-sit.yml
            name: configmap-volume-0
            subPath: application-sit.yml
        volumes:
          - configMap:
              defaultMode: 420
              items:
                - key: application-sit-yml
                  path: application-sit.yml
              name: ${your_service} 
            name: ${your_service}-volume-0
        restartPolicy: Always
```
Check the Kubernetes creation by running cmd `kubectl get pod -n ${your_ns} ${your_service}`
* **Note**
1. The configmap section is optional, depends on whether your application need it or not.
2. The Kubernetes namespace is also optional, you can choice to use the "default" namespace. Once you decide to use a particular namespace, make sure it is already exist.(you can run `kuberctl create ns ${your_ns}` to create yours)
3. The Eureka URL is always `http://127.0.0.1:13009/mesh/eureka`. If you are using Consul, the URL will be `http://127.0.0.1:13009`. In Nacos scenario, the URL will be `http://127.0.0.1:13009/nacos/v1`
### 3. Traffic Management 
#### 3.1 Resilience 
* EaseMesh implements four key types of resilience features, RateLimiter, CircuitBreaker, Retryer and Timeout by following Resilience4j library.
##### 3.1.1 RateLimiter
* Background: RateLimiter can establish your services' high availability and reliability, also it can be used for scaling APIs.  Protect your servers from overwhelm by peek traffic. 
* Steps: 
1. Deploy your application in EaseMesh, use cmd `eashmesh/bin/meshctl service get ${your_service_name}` to see current mesh service configuration and cmd `kubectl get pods ${your_service_pod_name}` to get whole Kubernetes pods and make sure there are pods running for it. 
2. We want to limit an API by specified HTTP method `POST` and `GET` and its URL which starts with prefix `/prefix` for accepting 50 request for 100 millisecond in service side. EaseMesh also supports URL matching with exact matching and regular expression matching. Once one request hit the current cycle's limit but there still have historical accumulated token left, it should wait for 100 millisecond for trying to get permitted. Available token will refresh every 10 millisecond for one cycle. 
3. Get current service's resilience spec by using cmd `easemesh/bin/meshctl service resilience get ${your_service_name"`, Add a RateLimiter into the `rateLimiter` section, save it into a new yaml file named `rateLimiter.yaml` 
```
rateLimiter:
  policy:
  - name: default
    timeoutDuration: 100 
    limitRefreshPeriod: 10
    limitForPeriod: 50 
  defaultPolicyRef: default 
  urls:
  - methods: ["POST", "PUT"]
    url:
      prefix: /users/
    policyRef: default 
```
4. Update the service with cmd `easemesh/bin/meshctl service resilience update ${your_service_name}  -f rateLimiter.yaml`

5. Once one upstream client hit the service's RateLimiter, it will receive HTTP response with header `X-EG-Rate-Limiter: too-many-requests`. 
* Field description 

| FieldName                   | type         | description                                                                                                                   |
| --------------------------- | ------------ | ----------------------------------------------------------------------------------------------------------------------------- |
| policy[].name               | string       | the name of this policy                                                                                                       |
| policy[].timeoutDuration    | string       | The duration for one request should wait for a permission,e.g.,`500ms`.                                                       |
| policy[].limitRefreshPeriod | string       | The period of a limit refresh. After each period the rate limiter sets its permissions count back to the limitForPeriod value |
| policy[].limitForPeriod     | int          | The number of permissions available during one limit refresh period                                                           |
| defaultPolicyRef            | string       | default applied policy name                                                                                                   |
| urls[].methods              | string array | HTTP methods, "POST","PUT","DELETE","GET"....                                                                                 |
| urls[].url.prefix           | string       | URL matching with prefix                                                                                                      |
| urls[].url.exact            | string       | URL matching with exactly                                                                                                     |
| urls[].url.regex            | string       | URL matching with regular expression                                                                                          |
| urls[].url.policyRef        | string       | the reference policy name, if its empty, will look up the `defaultPolicyRef` policy                                           |

##### 3.1.2 CircuitBreaker  
* Background: CircuitBreaker is used for blocking all in-coming requests when the the failure numbers reach the limit. You can declare an CircuitBreaker with **COUNT_BASED** or **TIME_BASED** type.  It has three types of states, open, closed and half-close. One service can declare its desired CircuitBreaker, and the upstream clients will active the same CircuitBreaker locally when calling this service. 
* Steps: 
1. We want to protect an API by specified HTTP method `GET` and its URL start with prefix `/users/` with **COUNT_BASED** sliding window type CircuitBreaker. It's sliding window count size is 20, the called service's failure analyzing conditions is when the HTTP response code is **500** and its failure rate threshold is 50%.  
2. Get current service's resilience spec by using cmd `easemesh/bin/meshctl service resilience get ${your_service_name"`, Add a CircuitBreaker into the `circuitBreaker` section, save it into a new yaml file named `circuitBreaker.yaml` 
```
circuitBreaker:
  policies:
  - name: default
    slidingWindowType: COUNT_BASED
    failureRateThreshold: 50
    slowCallRateThreshold: 100
    countingNetworkError: false
    slidingWindowSize: 20
    permittedNumberOfCallsInHalfOpenState: 10
    minimumNumberOfCalls: 10
    slowCallDurationThreshold: 100ms
    maxWaitDurationInHalfOpenState: 60s
    waitDurationInOpenState: 60s
    failureStatusCodes: [500]
  defaultPolicyRef: default
  urls:
  - methods:
    - GET
    url:
      exact: ""
      prefix: /users/
      regex: ""
    policyRef: "" 
```
4. Update the service with cmd `easemesh/bin/meshctl service resilience update ${your_service_name}  -f circuitBreaker.yaml`
5. Once the client active the CircuitBreaker, the client will receive HTTP response header with field `X-EG-Circuit-Breaker: circuit-is-broken`. 

* Field description 

| FieldName                                        | type         | description                                                                                                                                                                                                                                                                                                                                                                                                         |
| ------------------------------------------------ | ------------ | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| policies[].name                                  | string       | the name of this policy                                                                                                                                                                                                                                                                                                                                                                                             |
| policies[].slidingWindowType                     | string       | COUNT_BASED or TIME_BASED                                                                                                                                                                                                                                                                                                                                                                                           |
| policies[].failureRateThreshold                  | int          | Configures the failure rate threshold in percentage.  When the failure rate is equal or greater than the threshold the CircuitBreaker transitions to open and starts short-circuiting calls.                                                                                                                                                                                                                        |
| policies[].slowCallRateThreshold                 | int          | Configures a threshold in percentage. The CircuitBreaker considers a call as slow when the call duration is greater than slowCallDurationThreshold When the percentage of slow calls is equal or greater the threshold, the CircuitBreaker transitions to open and starts short-circuiting calls.                                                                                                                   |
| policies[].countingNetworkError                  | bool         | If circuit breaker active in network failure situation or not                                                                                                                                                                                                                                                                                                                                                       |
| policies[].permittedNumberOfCallsInHalfOpenState | int          | Configures the number of permitted calls when the CircuitBreaker is half open.                                                                                                                                                                                                                                                                                                                                      |
| policies[].minimumNumberOfCalls                  | int          | Configures the minimum number of calls which are required (per sliding window period) before the CircuitBreaker can calculate the error rate or slow call rate.  For example, if minimumNumberOfCalls is 10, then at least 10 calls must be recorded, before the failure rate can be calculated. If only 9 calls have been recorded the CircuitBreaker will not transition to open even if all 9 calls have failed. |
| policies[].maxWaitDurationInHalfOpenState        | int          | Configures a maximum wait duration which controls the longest amount of time a CircuitBreaker could stay in Half Open state, before it switches to open. Value 0 means Circuit Breaker would wait infinitely in HalfOpen State until all permitted calls have been completed.                                                                                                                                       |
| policies[].waitDurationInOpenState               | string       | The time that the CircuitBreaker should wait before transitioning from open to half-open,e.g.,`60000ms`.                                                                                                                                                                                                                                                                                                            |
| policies[].failureStatusCodes                    | int array    | The array of HTTP response code                                                                                                                                                                                                                                                                                                                                                                                     |
| defaultPolicyRef                                 | string       | default applied policy name, if its empty, will look up the `defaultPolicyRef` policy                                                                                                                                                                                                                                                                                                                               |
| urls[].methods                                   | string array | HTTP methods, "POST","PUT","DELETE","GET"....                                                                                                                                                                                                                                                                                                                                                                       |
| urls[].url.prefix                                | string       | URL matching with prefix                                                                                                                                                                                                                                                                                                                                                                                            |
| urls[].url.exact                                 | string       | URL matching with exactly                                                                                                                                                                                                                                                                                                                                                                                           |
| urls[].url.regex                                 | string       | URL matching with regular expression                                                                                                                                                                                                                                                                                                                                                                                |
| urls[].url.policyRef                             | string       | the reference policy name, if its empty, will look up the `defaultPolicyRef` policy                                                                                                                                                                                                                                                                                                                                 |
##### 3.1.3 Timeout(TimeLimiter)  
* Background: Timeout is the amount of time the client should wait for replies from a given service, it will be running in upstream clients and declared in downstream relied services. 
* Steps:
1. We want to cancel an API calling by specified HTTP method `GET` and its URL start with prefix `/users/` with 100 milliseconds.
2. Get current service's resilience spec by using cmd `easemesh/bin/meshctl service resilience get ${your_service_name}"`, Add a TimeLimiter into the `timeLimiter` section, save it into a new yaml file named `timeLimiter.yaml` .
```
timeLimiter:
  defaultTimeoutDuration: 600ms 
  urls:
  - methods: ["POST", "PUT"]
    url:
      prefix: /users/
    timeoutDuration: 100ms
```
4. Update the service with cmd `easemesh/bin/meshctl service resilience update ${your_service_name}  -f timeLimiter.yaml`
5. Once the client active the CircuitBreaker, the client will receive HTTP response header with field `X-EG-Time-Limiter: time-out`. 

* Field description 

| FieldName                  | type         | description                                      |
| -------------------------- | ------------ | ------------------------------------------------ |
| defaultTimeoutDuration     | string       | the default duration for timeout, e.g.,`500ms`.  |
| urls[].methods             | string array | HTTP methods, "POST","PUT","DELETE","GET"....    |
| urls[].url.prefix          | string       | URL matching with prefix                         |
| urls[].url.exact           | string       | URL matching with exactly                        |
| urls[].url.regex           | string       | URL matching with regular expression             |
| urls[].url.timeoutDuration | string       | the duration for this API's timeout,e.g.`600ms`. |

##### 3.1.4 Retryer  
* Background: Retryer can perform an API calling retry when the service HTTP response code indicated its in temporary unavailable states. The up-stream client should make sure this API is idempotent. The service can declare an Retryer for its desired APIs and active in client side. 
* Steps:
1. We want to use an Retryer for calling one API by specified HTTP method `GET` and its URL start with prefix `/users/`. It can retry at most 3 times, each try should wait 10 millisecond with exponential back off policy. 
2. Get current service's resilience spec by using cmd `easemesh/bin/meshctl service resilience get ${your_service_name}"`, Add a Retryer into the `retryer` section, save it into a new yaml file named `retryer.yaml` .
```
retryer:
  policies:
    - name: default
      maxAttempts: 3
      waitDuration: 10ms
      backoffPolicy: ExponentialBackOff
      countingNetworkError: false 
      failureStatusCodes:
      - 500
        503 
    - name: usersAPIPolicy
      maxAttempts: 3
      waitDuration: 10
      backoffPolicy: RandomBackOff
      randomizationFactor: 0.5
  defaultPolicyRef: default       
  urls:
  - methods: ["POST", "PUT"]
    url:
      prefix: /users/
    policyRef: usersAPIPolicy 
```
3. Once the client uses retryer successfully, the client will receive HTTP response header with field `X-EG-Time-Limiter: time-out`.
* Field description 

| FieldName                       | type         | description                                                                      |
| ------------------------------- | ------------ | -------------------------------------------------------------------------------- |
| policies[].name                 | string       | the name of this retry policies.                                                 |
| policies[].maxAttempts          | int          | The maximum number of attempts (including the initial call as the first attempt) |
| policies[].waitDuration         | string       | A based and fixed wait duration between retry attempts.                          |
| policies[].backoffPolicy        | string       | `ExponentialBackOff` or `RandomBackOff`                                          |
| policies[].randomizationFactor  | float        | float value between 0 and 1                                                      |
| policies[].countingNetworkError | bool         | If retry in network failure situation                                            |
| policies[].failureStatusCodes   | int array    | An HTTP statue codes array when retryer can perform                              |
| defaultPolicyRef                | string       | default applied policy name                                                      |
| urls[].methods                  | string array | HTTP methods, "POST","PUT","DELETE","GET"....                                    |
| urls[].url.prefix               | string       | URL matching with prefix                                                         |
| urls[].url.exact                | string       | URL matching with exactly                                                        |
| urls[].url.regex                | string       | URL matching with regular expression                                             |
| urls[].url.policyRef            | string       | the desired apply retry policy name                                              |

#### 3.2 Canary deployment
* Background: When new version of service called canary version want to be applied into formal environment, after unit testing, integration testing and regression testing, we still need to deploy these canary version's instances with small amount to accept some real and colored traffic. The colored traffic means when some targeted users with specified labels, the traffic gateway will color this user's traffic with desired labels. When this new instances deal with colored traffic for some while and become stable, we can scale the canary version's number to replace the former version's service instances. 
* Steps:
1. We want to add a canary version with mesh service label `version: canary`, and they will handle the colored traffic which has `X-Mesh-Canary: lv1` HTTP header. 
2. Deploy the canary version with instance label and new image URL
```
apiVersion: mesh.megaease.com/v1beta1
kind: MeshDeployment
metadata:
  namespace: ${your_ns} 
  name: ${your_service}-canary 
spec:
  service:
    name: ${your_service} 
    # labels for this canary instances
    labels:
    - version: canary
  deploy:
    replicas: 2 
    selector:
      matchLabels:
        app: ${your_service} 
    template:
      metadata:
        labels:
          app: ${your_service} 
      spec:
        containers:
        # the canary service's new image URL
        - image: ${your_image_new_url} 
          name: ${your_service} 
          imagePullPolicy: IfNotPresent
          lifecycle:
            preStop:
              exec:
                command: ["sh", "-c", "sleep 10"]
          command: ["/bin/sh"]
          args: ["-c", "java -server -Xmx1024m -Xms1024m -Dspring.profiles.active=sit -Djava.security.egd=file:/dev/./urandom -jar /application/application.jar"]
          resources:
            limits:
              cpu: 2000m
              memory: 1Gi
            requests:
              cpu: 200m
              memory: 256Mi
          ports:
          - containerPort: 8080
          volumeMounts:
          - mountPath: /application/application-sit.yml
            name: configmap-volume-0
            subPath: application-sit.yml
        volumes:
          - configMap:
              defaultMode: 420
              items:
                - key: application-sit-yml
                  path: application-sit.yml
              name: ${your_service} 
            name: ${your_service}-volume-0
        restartPolicy: Always
```
Use `kubectl get pod -l app: ${your_service}`, to make sure original version's service instances and canary version's running status:
```
NAME                                      READY   STATUS    RESTARTS   AGE
${your_service}-6c59797565-qv927      2/2     Running   0          8d
${your_service}-6c59797565-wmgw7      2/2     Running   0          8d
${your_service}-canary-84586f7675-lhrr5      2/2     Running   0          5min 
${your_service}-canary-7fbbfd777b-hbshm      2/2     Running   0          5min 
```
3. Apply the canary rule for your services with yaml named `canary.yaml` as below
```
canary:
  canaryRule:
  - serviceLabels:
      version: canary
    filter:
      headers:
        X-Mesh-Canary:
          values:
          - lv1
```
Use `easemesh/bin/meshctl service update ${your_service} canary -f ./canary.yaml` to apply this canary rule.
4. Make sure your traffic gateway color your target user's visit traffic with HTTP header `X-Mesh-Canary: lv1`.
5. You can use cmd `kubectl scale deployment ${your_service} --replicas=${increased_nums}` to scale the canary version's instances number. 
6. After fully real traffic testing, we can use `easemesh/bin/meshctl service instance list ${your_service}` to get whole instances list for your service. `eashmesh/bin/meshctl service instance clearLabel ${the_canary_service_instances_id}` to make the canary version instances become the new stable version 
7. Use cmd `eashmesh/bin/meshctl service instance offline ${the_original_service_instances_id}` to expel the old version's instances.

* Field description 

| FieldName                            | type   | description                                                                          |
| ------------------------------------ | ------ | ------------------------------------------------------------------------------------ |
| canaryRule[].serviceLabels           | map    | The canary instances' label                                                          |
| canaryRule[].filter[].headers.values | string | The exact matching string value for colored traffic's HTTP header value              |
| canaryRule[].filter[].headers.regexp | string | The regular expression matching string value for colored traffic's HTTP header value |

#### 3.3 Ingress Gateway   
* Background: MeshIngress is the rule to describe how traffic will be routed into mesh's internal after traffic gateway.  
* Step:
1. Deploy your service according #2 section.
2. We want to route HTTP traffic with HOST `${your_service}.com`, prepare the ingress rule named `ingress-rule.yaml` as below
```
name: ${the_ingress_rule_name} 
rules:
- host: ${your_service}.com 
  paths:
  - path: /
    backend: ${your_service}
```
Deploy it with `easemesh/bin/meshctl ingress create -f ./ingress-rule.yaml` 

* Filed description

| FieldName                    | type   | description                                                    |
| ---------------------------- | ------ | -------------------------------------------------------------- |
| name                         | string | The ingress rule's name                                        |
| rule[].host                  | string | The HOST value of your service visit URL                       |
| rule[].paths[].path          | string | The HTTP request path value                                    |
| rule[].paths[].rewriteTarget | string | The regular expression for rewriting the original request path |

### 4. Observability 
* Background:  In order to achieve better micro-services governance, EaseMesh need to provide observability of service behavior. It can empower operator/developer to troubleshoot, maintain, and optimize their applications.  
 
#### 4.1 Output Kafka 
* Background: EaseMesh will linkage EaseMonitor for aggregating and dealing with all services' metrics/logs/tracings, so we need a bridge which is **Kafka** to connect this two product. EaseMesh use JavaAgent technology to collecting all desired data, and output them into the bridge Kafka. 
* Steps:
1. Prepare the Kafka's visit URL and the configuration yaml named `output.yaml` as below:
```
outputServer:
  enabled: true
  bootstrapServer: ${your_kafka_host_one}:9093,${your_kafka_host_two}:9093,${your_kafka_host_three}:9093
  timeout: 30000
```
2. Update it into your service with cmd `easemesh/bin/meshctl service update ${your_service} observability -f ./output.yaml`
#### 4.2 Distributed Tracing
* Background: EaseMesh generates distributed trace spans for each services inside mesh. The operator/developer can fully understanding service dependencies and request flows. Also EaseMesh supports many kinds of tracing recording, including HTTP-Request, Remote-Invoking, Kafka, JDBC, Redis and RabbitMQ. 
* Steps:
1. We want to enable all tracing recording. Prepare the Tracing configuration named `tracing.yaml` as below:
```
tracings:
  enabled: true
  sampleByQPS: 30
  output:
    enabled: true
    reportThread: 1
    topic: log-tracing
    messageMaxBytes: 999900
    queuedMaxSpans: 1000
    queuedMaxSize: 1000000
    messageTimeout: 1000
  request:
    enabled: true
    servicePrefix: httpRequest
  remoteInvoke:
    enabled: true
    servicePrefix: remoteInvoke
  kafka:
    enabled: true
    servicePrefix: kafka
  jdbc:
    enabled: true
    servicePrefix: jdbc
  redis:
    enabled: true
    servicePrefix: redis
  rabbit:
    enabled: true
    servicePrefix: rabbit
``` 
2. Update it into your service with cmd `easemesh/bin/meshctl service update ${your_service} observability -f ./tracing.yaml`

* Field description 

| FieldName              | type   | description                                                                                 |
| ---------------------- | ------ | ------------------------------------------------------------------------------------------- |
| enabled                | bool   | Enabled this service's global tracing reporting switch                                      |
| sampleByQPS            | int    | Collects sample by QPS threshold, more than sampleByQPS value's requests won't be collected |
| output.enabled         | bool   | Enabled output to observability Kafka or not                                                |
| output.reportThread    | int    | The number of reporting Java threads                                                        |
| output.topic           | string | The output Kafka's topic name                                                               |
| output.messageMaxBytes | int    | The output Kafka's message max bytes                                                        |
| output.queuedMaxSpans  | int    | The output Kafka's queue max span number                                                    |
| output.queuedMaxSize   | int    | The output Kafka's queue max size                                                           |
| output.messageTimeout  | int    | The output Kafka's message timeout                                                          |

3. View the tracing recording in MegaEase portal: ![The tracing diagram](/imgs/tracing.png)

#### 4.3 Metrics & AccessLog 
* Background: EaseMesh collects service-level metrics for monitoring services communication inside mesh. The Metrics cover throughput ratio, executions error ratio, executions latency, response distribution and so on. Also EaseMesh supports many kinds of metrics recording, including Access-Log, HTTP-Request, Remote-Invoking, Kafka, JDBC, Redis and RabbitMQ. 
* Steps:
1. We want to enable all variable metrics reporting. Prepare the metrics configuration named `metrics.yaml` as below:
```
metrics:
  enabled: true
  access:
    enabled: true
    interval: 30000
    topic: application-log
  request:
    enabled: true
    interval: 30000
    topic: application-meter
  jdbcStatement:
    enabled: true
    interval: 30000
    topic: application-meter
  jdbcConnection:
    enabled: true
    interval: 30000
    topic: application-meter
  rabbit:
    enabled: true
    interval: 50000
    topic: platform-meter
  kafka:
    enabled: true
    interval: 40000
    topic: platform-meter
  redis:
    enabled: true
    interval: 70000
    topic: platform-meter
  jvmGc:
    enabled: true
    interval: 30000
    topic: platform-meter
  jvmMemory:
    enabled: true
    interval: 30000
    topic: platform-meter
  md5Dictionary:
    enabled: true
    interval: 30000000000
    topic: application-meter
``` 
2. Update it into your service with cmd `easemesh/bin/meshctl service update ${your_service} observability -f ./metrics.yaml`

* Field description 

| FieldName       | type   | description                                                                 |
| --------------- | ------ | --------------------------------------------------------------------------- |
| enabled         | bool   | Enabled this service's global metrics reporting switch                      |
| access.enabled  | bool   | Enabled access log metrics section or not                                   |
| access.interval | int    | The access log reporting interval, it's millisecond, default value is 30000 |
| access.topic    | string | The access log reporting to which Kafka topic                               |
                        

3. View the metrics in MegaEase portal: ![The Metrics diagram](/imgs/metrics.png))