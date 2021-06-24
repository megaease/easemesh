# EaseMesh Manual  

- [EaseMesh Manual](#easemesh-manual)
  - [Mesh service](#mesh-service)
  - [MeshDeployment](#meshdeployment)
  - [Sidecar Traffic](#sidecar-traffic)
    - [Inbound](#inbound)
    - [Outbound](#outbound)
    - [Sidecar Configuration](#sidecar-configuration)
  - [Resilience](#resilience)
    - [CircuitBreaker](#circuitbreaker)
    - [RateLimiter](#ratelimiter)
    - [Retryer](#retryer)
    - [TimeLimiter](#timelimiter)
  - [Observability](#observability)
    - [Tracing](#tracing)
      - [Turn-on tracing](#turn-on-tracing)
      - [Turn-off tracing](#turn-off-tracing)
    - [Metrics](#metrics)
      - [Turn-on metrics reporting](#turn-on-metrics-reporting)
      - [Turn-off metrics reporting](#turn-off-metrics-reporting)
    - [Log](#log)
      - [Turn-on Log](#turn-on-log)
      - [Turn-off Log](#turn-off-log)


 EaseMesh divides the main components into two parts, one is **Control plane**, the other is **Data plane**. In the control plane, EaseMesh uses the Easegress cluster to form a reliable decision delivery and persistence unit. The data plane is composed of each mesh service with the user's business logic and EaseMesh's enhancement units, including EaseAgent and Easegress-sidecar. And there is also a Mesh Ingress unit for routing and handling South-North way traffic.

![Architecture](../imgs/architecture.png)

## Mesh service

Services are the first-class citizens of the EaseMesh. Developers need to breakdown their business logic into small units and implement it as services.

A service could have co-exist multiple versions, a version of the service is a [MeshDeployment](#meshdeployment) which binds to Kubernetes Deployment resource

The `tenant` is used to group several services of the same business domain. Services can communicate with each other in the same tenant. In EaseMesh, there is a special global tenant that is visible to the entire mesh. Users can put some global, shared services in this special tenant.

> ** Note: **
> All specs in the EaseMesh are written in Yaml formation
> Please remember to change the YAML's placeholders such as ${your-service-name} to your real service name before applying.

1. **Create a tenant for services** You can choose to deploy a new mesh service in an existing tenant, or creating a new one for it. Modify example YAML content below and apply it 


```yaml
name: ${your-tenant-name}
description: "This is a test tenant for EaseMesh demoing"
```

>Tenant Spec reference: https://github.com/megaease/easemesh-api/blob/master/v1alpha1/meshmodel.md#easemesh.v1alpha1.Tenant 

2. **Create a service and specify which tenant the service belonged to**. Creating your mesh service in EaseMesh. Note, we only need to add this new service's logic entity now. The actual business logic and the way to deploy will be introduced later. Modify example YAML content below and apply it 


```yaml
name: ${your-service-name}
registerTenant: ${your-tenant-name}
loadBalance:
  policy: roundRobin
  HeaderHashKey:
sidecar:
  discoveryType: eureka
  address: "127.0.0.1"
  ingressPort: 13001
  ingressProtocol: http
  egressPort: 13002
  egressProtocol: http
```

>Service Spec Reference: https://github.com/megaease/easemesh-api/blob/master/v1alpha1/meshmodel.md#easemesh.v1alpha1.Service

3. With steps 1 and 2, now we have a new tenant and a new mesh service. They are both logic units without actual processing entities. EaseMesh relies on Kubernetes to transparent the resource management and deployment details. In K8s, we need to build the business logic (your **Java Spring Cloud application**) into an image and tell K8s the number of your instances and resources, with so-called declarative API, mostly in a YAML form. We will use a K8s [ConfigMap](https://kubernetes.io/docs/concepts/configuration/configmap/) to store your application's configurations and an [Custom Resource Define(CRD)](https://kubernetes.io/docs/concepts/extend-kubernetes/api-extension/custom-resources/) called `MeshDeployment` to describe the way you want your application instances run in K8s. Here is a Java Spring Cloud application example that visiting MySQL, Eureka for service discovery/register. Preparing the deployment YAML by modifying content below, and applying it 

```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: ${your-configmap-name} 
  namespace: ${your-ns-name} 
data:
  application-sit-yml: |
    server:
      port: 8080
    spring:
      application:
        name:  $(your-service-name} 
      datasource:
        url: jdbc:mysql://mysql.default:3306/${your_db_name}?allowPublicKeyRetrieval=true&useUnicode=true&characterEncoding=utf-8&useSSL=false&serverTimezone=UTC&verifyServerCertificate=false
        username: ${your-db-username} 
        password: {$your-db-password} 
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
  namespace: ${your-ns-name} 
  name: ${your-service-name} 
spec:
  service:
    name: ${your-service-name}
  deploy:
    replicas: 2 
    selector:
      matchLabels:
        app: ${your-service-name}
    template:
      metadata:
        labels:
          app: ${your-service-name} 
      spec:
        containers:
        - image: ${your-image-url} 
          name: ${your-service-name} 
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
              name: ${your-service-name} 
            name: ${your_service}-volume-0
        restartPolicy: Always
```

4. For service register/discovery, EaseMesh supports three mainstream solutions, Eureka/Consul/Nacos. Check out the corresponding configuration URL below:

| Name   | URL In Mesh deployment configuration |
| ------ | ------------------------------------ |
| Eureka | http://127.0.0.1:13009/mesh/eureka   |
| Consul | http://127.0.0.1:13009               |
| Nacos  | http://127.0.0.1:13009/nacos/v1      |

5. Communications between internal mesh services can be done through Spring Cloud's recommended clients, such as `WebClient`, `RestTemplate`, and `FeignClient`. The original HTTP domain-based RPC remains unchanged. Please notice, EaseMesh will host the Ease-West way traffic by its mesh service name, so it is necessary to keep the mesh service name the same as the original Spring Cloud application name for HTTP domain-based RPC.

## MeshDeployment

EaseMesh relies on Kubernetes for managing service instances and the resources they require. For example, we can scale the number of instances with the help of Kubernetes. In fact, EaseMesh uses a mechanism called [Kubernetes  Custom Resource Define(CRD)](https://kubernetes.io/docs/concepts/extend-kubernetes/api-extension/custom-resources/) to combine the service metadata used by EaseMesh and Kubernetes original deployment. MeshDeployment can be used not only to deploy and manage service instances, it can also help us implement the canary deployment. 

MeshDeployment wraps native K8s [Deployment](https://kubernetes.io/docs/concepts/workloads/controllers/deployment) resources. The contents of `spec.deploy` section in the MeshDeployment spec is fully K8s deployments spec definition.


The canary deployment is a pattern for rolling out releases to a subset of servers. The idea is to first deploy the change to a small subset of servers, test it with real users' traffic, and then roll the change out to the rest of the servers. The canary deployment serves as an early warning indicator with less impact on downtime: if the canary deployment fails, the rest of the servers aren't impacted. In order to be safer, we can divide traffic into two kinds, normal traffic, and colored traffic. Only the colored traffic will be routed to the canary instance. The traffic can be colored with the users' model, then setting into standard HTTP header fields. 

![canary-deployment](./../imgs/canary-deployment.png)

1. Preparing new business logic with a new application image. Adding a new `MeshDeployment`, we would like to separate the original mesh server's instances from the new canary instances by labeling `version: canary` to canary instances. Modify example YAML content below and apply it 

```yaml
apiVersion: mesh.megaease.com/v1beta1
kind: MeshDeployment
metadata:
  namespace: ${your-ns-name} 
  name: ${your_service-name}-canary 
spec:
  service:
    name: ${your-service-name} 
    labels:
    - version: canary       # These map is used to label these canary instances
  deploy:                   # K8s native deployment spec contents
    replicas: 2 
    selector:
      matchLabels:
        app: ${your-service-name}   #Note! service name should remain the same with the origin mesh service  
    template:
      metadata:
        labels:
          app: ${your-service-name} 
      spec:
        containers:
        - image: ${your-image-new-url}    # the canary instance's new image URL
          name: ${your-service-name} 
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
              name: ${your-service-name} 
            name: ${your_service}-volume-0
        restartPolicy: Always

```
2. Checking the original normal instances and canary instances with cmd 

```bash
$ kubectl get pod -l app: ${your-service-name}

NAME                                      READY   STATUS    RESTARTS   AGE
${your-service-name}-6c59797565-qv927      2/2     Running   0          8d
${your-service-name}-6c59797565-wmgw7      2/2     Running   0          8d
${your-service-name}-canary-84586f7675-lhrr5      2/2     Running   0          5min 
${your-service-name}-canary-7fbbfd777b-hbshm      2/2     Running   0          5min 
```

3. When canary instances are ready for work, it's time to set the policy for traffic-matching. In this example, we would like to color traffic for the canary instance with HTTP header filed `X-Mesh-Canary: lv1` (Note, we want exact matching here, can be set to a Regular Expression) and all mesh service's APIs are the canary targets. Modify example canary rule YAML content below and apply it 

```yaml
canary:
  canaryRules:
  - serviceInstanceLabels:
      version: canary   # The canary instances must have this `version: canary` label.
    headers:
        X-Mesh-Canary:
          exact: lv1    # The colored traffic with this exact matching HTTP header value.
    urls:
      - methods: ["GET","POST","PUT","DELETE"]
        url:
          prefix: "/"  # Routing colored traffic to canary instances all HTTP APIs.

```

> CanaryRule spec reference: https://github.com/megaease/easemesh-api/blob/master/v1alpha1/meshmodel.md#easemesh.v1alpha1.CanaryRule


4. Visiting your mesh service with and without HTTP header `X-Mesh-Canary: lv1`, the colored traffic will be handled by canary instances.

## Sidecar Traffic

In `EaseMesh`, we use `EaseMeshController` based on `Easegress` to play the `Sidecar` role. As a sidecar, the mesh controller will handle inbound and outbound traffic. The inbound traffic means business traffic from outside to sidecar, and the outbound traffic means business traffic from sidecar to outside. We make them clean by the simple diagram:

**InBound Traffic**

![inbound traffic](../imgs/inbound-traffic.png)

### Inbound
MeshController will create a dedicated pipeline to handle inbound traffic:
1. Accept business traffic from outside in one port.
2. Use RateLimiter (See below) to do rate limiting.
3. Transport traffic to the service.


**OutBound Traffic**

![outbound traffic](../imgs/outbound-traffic.png)

### Outbound
MeshController will create dedicated pipelines to handle outbound traffic:
1. Accept business traffic from service in one port.
2. Use CircuitBreaker, Retryer, TimeLimiter to do protection for receiving services according to their own config.
3. Use load balance to choose the service instance.
4. Transport traffic to the chosen service instance.



> The diagram above is a logical direction of the **request** of traffic, the responses flow in the opposite direction which is the same category with corresponding requests.

Please notice the sidecar only handle business request traffic, which means it doesn't hijack traffic to:
1. Middleware, such as Redis, Kafka, Mysql, etc.
2. Any other control plane, such as the `Istio` pilot.

But we are well compatible with the Java ecosystem, so we adapt the mainstream service discovery registry like Eureka, Nacos, and Consul. We do hijack traffic to the service discovery, so it's required that the service **changes service registry address to sidecar address** in the startup-config.

### Sidecar Configuration
* **Note: Please remember to change the YAML's placeholders to your real service name tenant name.**

```yaml
name: ${your-service-name}
registerTenant: ${your-tenant-name} 
loadBalance:
  policy: roundRobin
sidecar:
  discoveryType: eureka
  address: "127.0.0.1"
  # Inbound traffic: OtherServices/Gateway -> Sidecar(http://127.0.0.1:13001) -> Service
  ingressPort: 13001
  ingressProtocol: http
  # Outbound traffic: Service -> Sidecar(http://127.0.0.1:13002) -> OtherServices
  # The OtherServices means multiple service instances under roundRobin policy.
  egressPort: 13002
  egressProtocol: http
```
> Sidecar Spec reference :https://github.com/megaease/easemesh-api/blob/master/v1alpha1/meshmodel.md#easemesh.v1alpha1.Sidecar


## Resilience

We borrow the core concept of the mature JAVA fault tolerate library [resilience4j](https://resilience4j.readme.io/) to implement the resilience. With the pipeline-filter(plugin) model of Easegress, We can assemble any of them together. Besides the function of each protection, we must know which side the protection takes effect in the Mesh scenario. We use the 2 clean terms: **sender** and **receiver** (of the request).

- sender: In general sender is a client which shots requests to the server
- receiver: In general receiver is a server that receives requests


### CircuitBreaker

In Mesh, `CircuitBreaker` takes effect in **sender** side, in another word, it applies on outbound traffic. For example, 

```yaml
name: ${your-service-name}
registerTenant: ${your-tenant-name}
resilience:
  circuitBreaker:
    policies:
      - name: count-based-example
        slidingWindowType: COUNT_BASED
        failureRateThreshold: 50
        slidingWindowSize: 100
        failureStatusCodes: [500, 503, 504]
    urls:
      - methods: [GET]
        url:
        prefix: /service-b/
        policyRef: count-based-example
      - methods: [GET, POST]
        url:
        prefix: /service-c/
        policyRef: count-based-example
```

> CircuitBreaker Spec reference: https://github.com/megaease/easemesh-api/blob/master/v1alpha1/meshmodel.md#easemesh.v1alpha1.CircuitBreaker 

The `sender` is `${your-service-name}`, `receiver`  side contains `service-b`  and `service-c`. So the circuit-breaker takes effect in `${your-service-name}`, and all responses from both `service-b` and `service-c` count to one circuit breaker here. Of course if the items of `urls` reference to different policies, the counting process will be in the respective circuit-breaker.

### RateLimiter

In Mesh, `RateLimiter` takes effect in `receiver` side, in another word, it applies on inbound traffic. For example:

```yaml
name: ${your-service-name}
registerTenant: ${your-tenant-name}
resilience:
  policies:
    - name: policy-example
      timeoutDuration: 100ms
      limitRefreshPeriod: 10ms
      limitForPeriod: 50
      defaultPolicyRef: policy-example
  urls:
    - methods: [GET, POST, PUT, DELETE]
      url:
        regex: ^/pets/\d+$
      policyRef: policy-example
```

> RateLimiter Spec reference :https://github.com/megaease/easemesh-api/blob/master/v1alpha1/meshmodel.md#easemesh.v1alpha1.RateLimiter

So all inbound traffic of `${your-service-name}` will be rate-limited by it, when the traffic character matches the policy. Please notice outbound traffic **from** `${your-service-name}` has no relationship with the rate limiter.

### Retryer

In Mesh, `Retry` takes effect in `sender` side, in another word, it applies on outbound traffic. For example :

```yaml
name: ${your-service-name}
registerTenant: ${your-tenant-name}
policies:
  - name: policy-example
    maxAttempts: 3
    waitDuration: 500ms
    failureStatusCodes: [500, 503, 504]
    defaultPolicyRef: policy-example
  urls:
    - methods: [GET, POST, PUT, DELETE]
      url:
        prefix: /books/
      policyRef: policy-example
```

> Retryer Spec reference: https://github.com/megaease/easemesh-api/blob/master/v1alpha1/meshmodel.md#easemesh.v1alpha1.Retryer 

All matching outbound traffic **from** `${your-service-name}` will be retried if the response code is one of `500`, `503`, and `504`.

### TimeLimiter
In Mesh, `TimeLimiter` takes effect in `sender` side. For example:

```yaml
name: ${your-service-name}
registerTenant: ${your-tenant-name}
urls:
- methods: [POST]
  url:
    exact: /users/1
  timeoutDuration: 500ms
```
> TimeLimiter Spec reference: https://github.com/megaease/easemesh-api/blob/master/v1alpha1/meshmodel.md#easemesh.v1alpha1.TimeLimiter

All matching outbound traffic **from** `${your-service-name}` have a timeout in `500ms`.


## Observability

Observability for micro-services in EaseMesh can be cataloged into three areas, distributed tracing, metrics, and logging. Users can see the details of a request, such as the complete request path,  invocation dependencies, and latency of each sub-requests so that issues can be diagnosed. Metrics can reflect the health level of the system and summarize its state. Logging is used to provide more details based on the requested access for helping resolving issues.  

### Tracing
* Tracing is disabled by default in EaseMesh. It can be enabled dynamically during the lifetime of the mesh services. Currently, the EaseMesh follow [OpenZipkin B3 specification](https://github.com/openzipkin/b3-propagation) to supports tracing these kinds of invocation :

| Name           | Description                                                                                                                                                                                                                                                              |
| -------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ |
| HTTP based RPC | Information about communication between mesh service via HTTP protocol, such as latency, status code, request path and so on. Currently, EaseMesh supports tracing for `WebClient`, `RestTemplate` and `FeignClient`, the more HTTP RPC libraries will be supported soon |
| JDBC           | Information about MySQL SQL execution latency, statement, results and so on.                                                                                                                                                                                             |
| Redis          | Information about Redis command latency, key, and so on.                                                                                                                                                                                                                 |
| RabbitMQ       | Information about RabbitMQ command latency, topic, routine key and so on.                                                                                                                                                                                                |
| Kafka          | Information about Kafka topics' latency and so on.                                                                                                                                                                                                                       |

* EaseMesh relies on `EaseAgent` for non-intrusive collecting span data, and Kafka to store all collected tracing data. 

#### Turn-on tracing

1. Configuring mesh service's `ObservabilityOutputServer` to enable EaseMesh output tracing related data into Kafka. Modify example YAML below, and apply it 

```yaml
 outputServer:
  enabled: true
  bootstrapServer: ${your_kafka_host_one}:9093,${your_kafka_host_two}:9093,${your_kafka_host_three}:9093
  timeout: 30000   
```
> OutputServer spec reference: https://github.com/megaease/easemesh-api/blob/master/v1alpha1/meshmodel.md#easemesh.v1alpha1.ObservabilityOutputServer 

2. Finding the desired enable tracing service protocol in [ObservabilityTracings](https://github.com/megaease/easemesh-api/blob/master/v1alpha1/meshmodel.md#easemesh.v1alpha1.ObservabilityTracings) structure. For example, turning on the switch in `ObservabilityTracings.remoteInvoke`  can record mesh service's HTTP RPC tracing data. Also, EaseMesh allows users to configure how Java Agent should report tracing data, such as the reporting sample rate, reporting thread numbers in JavaAgent, and so on. **Note: the reporting configuration is global inside one mesh service's tracing** . Modify example YAML below and applying it 

```yaml
tracings:
  enabled: true              # The global enable switch
  sampleByQPS: 30            # The data above QPS 30 will be ignore 
  output:
    enabled: true            # Enabling Kafka reporting  
    reportThread: 1          # Using one thread to report in JavaAgent
    topic: log-tracing       # The reporting Kafka topic name 
    messageMaxBytes: 999900  # 
    queuedMaxSpans: 1000
    queuedMaxSize: 1000000
    messageTimeout: 1000
  request:
    enabled: false 
    servicePrefix: httpRequest
  remoteInvoke:
    enabled: true                # Turing on this switch for RPC tracing only
    servicePrefix: remoteInvoke
  kafka:
    enabled: false 
    servicePrefix: kafka
  jdbc:
    enabled: false 
    servicePrefix: jdbc
  redis:
    enabled: false  
    servicePrefix: redis
  rabbit:
    enabled: false  
    servicePrefix: rabbit
```

>ObservabilityTracings spec reference: https://github.com/megaease/easemesh-api/blob/master/v1alpha1/meshmodel.md#easemesh.v1alpha1.ObservabilityTracings

4. Tracing data are organized as spans, each span is stored in the backend storage service, which provides online analysis and computing functions. MegaEase provides a sophisticated view to help users rapidly diagnosing problems. Checking the web console for your mesh service's RPC tracing information:

![tracing](../imgs/tracing.png)

#### Turn-off tracing

1. If you want to disable tracing for one mesh service, then set this mesh service's global [tracing switch](https://github.com/megaease/easemesh-api/blob/master/v1alpha1/meshmodel.md#easemesh.v1alpha1.ObservabilityTracings) to `off`. For example, you can prepare YAML as below and apply it 

```yaml
tracings:
  enabled: false             # The global enable switch
  sampleByQPS: 30            # The data above QPS 30 will be ignore 
  output:
    enabled: tru

    ......

```

2. For only disabling one tracing feature for one mesh service, find the corresponding section, then turn off its switch is enough. For example, to shut down one mesh service's Redis tracing feature, you can prepare YAML as bellow and apply it 

```yaml
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
    enabled: false               # Turing off this switch for not tracing Redis invocation 
    servicePrefix: redis
```  


### Metrics

* The EaseMesh leverage [EaseAgent( JavaAgent based on Java Byte buddy technology)](https://github.com/megaease/easeagent) to collect mesh services' application metrics in a non-intrusive way. It will collect the data from a service perspective with very low CPU, memory, I/O resource consumption. The supported metric types including:

* For the metric details for every type, checkout the EaseAgent's [develop-guide.md](https://github.com/megaease/easeagent/blob/master/doc/development-guide.md).

* Here are the metics that EaseMesh already supported:


| Name              | Description                                                                                                                                                                                                                                                                                                 |
| ----------------- | ----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| HTTP request      | The mesh service's HTTP APIs metrics, such as m1/m5/m15 rate(`m1` indicates the The http request executions per second `exponentially-weighted moving average` in last 1 minute ), URL, total counts, error counts, p99(The http-request execution duration in milliseconds for 99% user) and so on.        |
| JDBC Statement    | The mesh service's JDBC statement metrics, such as the signature(used for complete SQL sentence matching), JDBC total count, JDBC m1 rate(The JDBC method executions per second `exponentially-weighted moving average` in last 1 minute.), TopN JDBC M1 error rate, JDBC P99 execution duration and so on. |
| JDBC Connection   | The mesh service's JDBC connection  metrics, such as URL, JDBC Connect total count, JDBC Connect M1 rate, JDBC Connect min execution duration, JDBC Connect min execution duration JDBC Connect P99 execution duration and so on                                                                            |
| JVM Memory        | The mesh service's JVM memory related metrics such as JVM initial memory, JVM used memory, JVM committed memory and JVM max memory.                                                                                                                                                                         |
| JVM GC            | The mesh service's JVM GC related metrics such as JVM Gc time, JVM Gc collection times and JVM Gc times per second.                                                                                                                                                                                         |
| Kafka Client      | The mesh service's Kafka client metics such as (**Note:** this is not the reported target Kafka, the user's application usage's Kafka) topic name, Kafka producer throughput(M1), kafka consumer throughput(M1), producer min execution duration and so on.                                                 |
| RabbitMq Producer | The mesh service's RabbitMQ producer's metics such as rabbit exchange, producer M1 rate, producer P99 execution duration and so on.                                                                                                                                                                         |
| RabbitMq Consumer | The mesh service's RabbitMQ consumer's metics such as rabbit exchange, consumer M1 rate, consumer P99 execution duration and so on.                                                                                                                                                                         |
| Redis             | The mesh service's Redis client's metics such as redis P25 execution duration, redis M1 count, redis P99 execution duration and so on.                                                                                                                                                                      |
| MD5 Dictionary    | The mesh service's JDBC statement's complete SQL sentences and MD5 values.                                                                                                                                                                                                                                  |
#### Turn-on metrics reporting  

* EaseMesh also reports the mesh service's Metrics into the Kafka used by Tracing. So you can check out how to enable the output Kafka in the Tracing section. 

1. Finding the desired enable metrics type in `ObservabilityMetrics` structure. For example, turning on switch in `ObservabilityMetrics.request`  can report mesh service's HTTP request-related metrics.Modify example YAML below and apply it   

```yaml
metrics:
  enabled: true                  # the global metrics reporting switch 
  access:
    enabled: false  
    interval: 30000
    topic: application-log
  request:
    enabled: true                 # the enable target metrics, HTTP request related   
    interval: 30000               # the interval between reporting, million seconds 
    topic: application-meter      # the reporting target kafka's topic name
  jdbcStatement:
    enabled: false 
    interval: 30000
    topic: application-meter
  jdbcConnection:
    enabled: false 
    interval: 30000
    topic: application-meter
  rabbit:
    enabled: false  
    interval: 50000
    topic: platform-meter
  kafka:
    enabled: false 
    interval: 40000
    topic: platform-meter
  redis:
    enabled: false  
    interval: 70000
    topic: platform-meter
  jvmGc:
    enabled: false  
    interval: 30000
    topic: platform-meter
  jvmMemory:
    enabled: false  
    interval: 30000
    topic: platform-meter
  md5Dictionary:
    enabled: false 
    interval: 30000000000
    topic: application-meter
```

> Metrics Spec reference: https://github.com/megaease/easemesh-api/blob/master/v1alpha1/meshmodel.md#easemesh.v1alpha1.ObservabilityMetrics 

2. Checking the web console for your mesh service's HTTP request metrics 

![metrics](../imgs/metrics.png)

#### Turn-off metrics reporting

1. If you want to disable metrics reporting for one mesh service, then set this mesh service's global `metrics reporting switch` to `off`. For example prepare YAML as below and apply it 

```yaml
metrics:
  enabled: false             # The global enable switch
  access: 
    ......

```

2. For only disabling one type of metrics reporting for one mesh service, find the corresponding section, then turn off its switch is enough. For example, to shut down one mesh service's HTTP request metrics reporting, you can prepare YAML as bellow and apply it 

```yaml
metrics:
  enabled: true                  # the global metrics reporting switch 
  access:
    enabled: false  
    interval: 30000
    topic: application-log
  request:
    enabled: false                # the disable target metrics, HTTP request related   
    interval: 30000               # the interval between reporting, million seconds 
    topic: application-meter      # the reporting target kafka's topic name
    ....
```


### Log
* Access log is also disabled by default in EaseMesh. The access log is used to recording details of HTTP APIs of mesh service been requested.  

#### Turn-on Log
* EaseMesh also reports the mesh service's Logs into the Kafka used by Tracing. So you can check out how to enable the output Kafka in the Tracing section.  

1. Finding the `access` section in [ObservabilityMetrics](https://github.com/megaease/easemesh-api/blob/master/v1alpha1/meshmodel.md#easemesh.v1alpha1.ObservabilityMetrics) structure. Turning on switch in `ObservabilityMetrics.access`  can enable access logging for mesh service's HTTP APIs. Modify example YAML below and apply it 

```yaml
metrics:
  enabled: true                  # the global metrics reporting switch 
  access:
    enabled: true                # the enable target metrics, HTTP request related    
    interval: 30000              # the interval between reporting, million seconds 
    topic: application-log       # the reporting target kafka's topic name
  request:
    enabled: false 
    interval: 30000               
    topic: application-meter      
  jdbcStatement:
    enabled: false 
    interval: 30000
    topic: application-meter
  jdbcConnection:
    enabled: false 
    interval: 30000
    topic: application-meter
  rabbit:
    enabled: false  
    interval: 50000
    topic: platform-meter
  kafka:
    enabled: false 
    interval: 40000
    topic: platform-meter
  redis:
    enabled: false  
    interval: 70000
    topic: platform-meter
  jvmGc:
    enabled: false  
    interval: 30000
    topic: platform-meter
  jvmMemory:
    enabled: false  
    interval: 30000
    topic: platform-meter
  md5Dictionary:
    enabled: false 
    interval: 30000000000
    topic: application-meter
```

> AccessLog reference: https://github.com/megaease/easemesh-api/blob/master/v1alpha1/meshmodel.md#easemesh.v1alpha1.ObservabilityMetrics

2. Checking the web console for your mesh service's HTTP log 

![access log](../imgs/accesslog.png)


#### Turn-off Log 
1. For disabling access logging for one mesh service, find the `access` section, then turn off its switch. For example, to shut down one mesh service's HTTP APIs' logging, you can modify YAML as bellow and apply it 

```yaml
metrics:
  enabled: true                  # the global metrics reporting switch 
  access:
    enabled: false               # disable this service's logging
    interval: 30000
    topic: application-log
    ....
```
