# Mesh Service  
![Architecture](../imgs/architecture.png)

1. EaseMesh divides the main components into two parts, one is **Control plane**, the other is **Data plane**. In the control plane, EaseMesh uses the Easegress cluster to form a reliable decision delivery and persistence unit. The data plane is composed of each mesh service with the user's business logic and EaseMesh's enhancement units, including EaseAgent and Easegress-sidecar. And there is also a Mesh Ingress unit for routing and handling South-North way traffic.   
2. EaseMesh is a solution for better service governance in the Java domain. `Mesh services` are first-class citizens. All governance is aimed at the mesh service concept. 
3. The `tenant` is used to group several mesh services of the same business domain. Mesh services can communicate with each other in the same tenant. In EaseMesh, there is a special global tenant that is visible to the entire mesh. Users can put some global, shared services in this special tenant.

## Create a mesh service
1. You can choose to deploy a new mesh service in an existing tenant, or creating a new one for it. Modifying the YAML content below and name it with `new_tenant.yaml`:
```yaml
name: ${your-tenant-name}
description: "this is a test tenant for EaseMesh demoing"
```
Applying it with cmd
```bash
$ ./eashmesh/bin/meshctl tenant create -f ./new_tenant.yaml
```

2. Creating your mesh service in EaseMesh. Note, we only need to add this new service's logic entity now. The actual business logic and the way to deploy will be introduced later. Modifying the YAML content below and name it with `new_service.yaml`:
* **Note: Please remember to change the YAML's placholders such as ${your-service-name} to your real service name before applying.**
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
Apply it with cmd
```bash
$ ./eashmesh/bin/meshctl service create -f ./new_service.yaml
```

3. With steps 1 and 2, now we have a new tenant and a new mesh service. They are both logic units without actual processing entities. EaseMesh relies on Kubernetes to transparent the resource management and deployment details. In K8s, we need to build the business logic (your **Java Spring Cloud application**) into an image and tell K8s the number of your instances and resources, with so-called declarative API, mostly in a YAML form. We will use a K8s [ConfigMap](https://kubernetes.io/docs/concepts/configuration/configmap/) to store your application's configurations and an [Custom Rousce Define(CRD)](https://kubernetes.io/docs/concepts/extend-kubernetes/api-extension/custom-resources/) called `MeshDeployment` to describe the way you want your application instances run in K8s. Here is a Java Spring Cloud application example that visiting MySQL, Eureka for service discovery/register. Preparing the deployment YAML by modifying content below, and naming it with `mesh_deployment.yaml`
```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: ${your_configmap_name} 
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
Using cmd to deploy it in K8s 
```bash
$ kubectl apply -f ./mesh_deployment.yaml
```
Checking the deployment result in K8s by running cmd 
```bash
kubectl get pod -n ${your-ns-name} ${your-service-name}
```

4. For service register/discovery, EaseMesh supports three mainstream solutions, Eureka/Consul/Nacos. Check out the corresponding URL below:

| Name   | URL In Mesh deployement configuration |
| ------ | ------------------------------------- |
| Erueka | http://127.0.0.1:13009/mesh/eureka    |
| Consul | http://127.0.0.1:13009                |
| Nacos  | http://127.0.0.1:13009/nacos/v1       |

5. Communications between internal mesh services can be done through Spring Cloud's recommended clients, such as `WebClient`, `RestTemplate`, and `FeignClient`. The original HTTP domain-based RPC remains unchanged. Please notice, EaseMesh will host the Ease-West way traffic by its mesh service name, so it is necessary to keep the mesh service name the same as the original Spring Cloud application name for HTTP domain-based RPC.

## Canary deployment 
![canary-deployment](./../imgs/canary-deployment.png)

Canary deployments are a pattern for rolling out releases to a subset of servers. The idea is to first deploy the change to a small subset of servers, test it with real users' traffic, and then roll the change out to the rest of the servers. The canary deployment serves as an early warning indicator with less impact on downtime: if the canary deployment fails, the rest of the servers aren't impacted. In order to be safer, we can divide traffic into two kinds, normal traffic, and colored traffic. Only the colored traffic will be routed to the canary instance. The traffic can be colored with the users' model, then setting into standard HTTP header fields. 

1. Preparing new business logic with a new application image. Adding a new `MeshDepoloyment`, we would like to separate the original mesh server's instances from the new canary instances by labeling `version: canary` to canary instances. Modify YAML content below and named it as `canary-deployment.yaml`
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
  deploy:
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

3. When canary instances are ready for work, it's time to set the policy for traffic-matching. In this example, we would like to color traffic for the canary instance with HTTP header filed `X-Mesh-Canary: lv1`(Note, we want exact matching here, can be set to a Regular Expression) and all mesh service's APIs are the canary targets. Modifying the canary rule YAML content below and named it with `canary-rule.yaml`
```yaml
canary:
  canaryRules:
  - serviceLabels:
      version: canary   # The canary instances must have this `version: canary` label.
    headers:
        X-Mesh-Canary:
          exact: lv1    # The colored traffic with this exatc matching HTTP header value.
    urls:
      - methods: ["*"]
        url:
          prefix: "/"  # Routing colored traffic to cananry instances all HTTP APIs.

```
Applying it with cmd 
```bash
$ ./easemesh/bin/meshctl service update ${your-service-name} canary -f ./canary-rule.yaml
```

4. Visiting your mesh service with and without HTTP header `X-Mesh-Canary: lv1`, the colored traffic will be handled by canary instances.

## Sidecar Traffic
In `EaseMesh`, we use `EaseMeshController` based on `Easegress` to play the `Sidecar` role. As a sidecar, the mesh controller will handle inbound and outbound traffic. The inbound traffic means business traffic from outside to sidecar, and the outbound traffic means business traffic from sidecar to outside. We make them clean by the simple diagram:
`Inbound traffic: OtherServices/Gateway -> Sidecar -> Service`
`Outbound traffic: Service -> Sidecar -> OtherServices`
The diagram above is a logical direction of the **request** of traffic, the responses flow in the opposite direction which is the same category with corresponding requests.
### Inbound
MeshController will create a dedicated pipeline to handle inbound traffic:
1. Accept business traffic from outside in one port.
2. Use RateLimiter (See below) to do rate limiting.
3. Transport traffic to the service.
### Outbound
MeshController will create dedicated pipelines to handle outbound traffic:
1. Accept business traffic from service in one port.
2. Use CircuitBreaker, Retryer, TimeLimiter to do protection for receiving services according to their own config.
3. Use load balance to choose the service instance.
4. Transport traffic to the chosen service instance.
Please notice the sidecar only handle business traffic, which means it doesn't hijack traffic to:
1. Middlewares, such as Redis, Kafka, Mysql, etc.
2. Any other control plane, such as Istio pilot.
But we are compatible with the Java ecosystem, so we adapt the mainstream service discovery registry like Eureka, Nacos, and Consul. We do hijack traffic to the service discovery, so it's required that the service **changes service registry address to sidecar address** in the startup-config.
### Config
* **Note: Please remember to change the YAML's placholders to your real service name tenant name.**
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
Checking out the [mesh service](https://github.com/megaease/easemesh-api/blob/master/v1alpha1/meshmodel.proto#38) and [sidecar](https://github.com/megaease/easemesh-api/blob/master/v1alpha1/meshmodel.proto#87) structure definitions for more field descriptions. 

## Resilience
We leverage the core concept of the mature fault tolerate library [resilience4j](https://resilience4j.readme.io/) to do the resilience. We choose the pipeline-filter model of Easegress to let us assemble any of them together. Besides the function of every protection, we must know which side the protection takes effect in the Mesh scenario. We use the 2 clean terms: **sender** and **receiver** (of the request).
### CircuitBreaker
The original [Martin Flower article](https://martinfowler.com/bliki/CircuitBreaker.html) describes it well enough. And [Filter CircuitBreaker](https://github.com/megaease/easegress/blob/main/doc/filters.md#circuitbreaker)  gives config and examples of it clearly.
In Mesh, `CircuitBreaker` takes effect in **sender** side. For example:
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
The `sender` is `${your-service-name}`, `receiver`  side contains `service-b`  and `service-c`. So the circuit breaker takes effect in `${your-service-name}`, and all responses from both `service-b` and `service-c` count to one circuit breaker here. Of course if the items of `urls` reference to different policies, the counting process will be in the respective circuit breaker.
### RateLimiter
[Filter RateLimiter](https://github.com/megaease/easegress/blob/main/doc/filters.md#ratelimiter) gives config and examples of it.
In Mesh, `RateLimter` takes effect in `receiver` side. For example:
```yaml
name: ${your-service-name}
registerTenant: tenant-exmaple
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
So all matching requests **to** `${your-service-name}` matching will go into the rate limiter `policy-example`. Please notice the requests **from** `${your-service-name}` has no relationship with the rate limiter.
### Retryer
[Filter Retryer](https://github.com/megaease/easegress/blob/main/doc/filters.md#retry) gives config and examples of it.
In Mesh, `Retry` takes effect in `sender` side. For example:
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
All matching requests **from** `${your-service-name}` will be retried if the response code is one of `500`, `503`, and `504`.
### TimeLimiter
[Filter TimeLimiter](https://github.com/megaease/easegress/blob/main/doc/filters.md#timelimiter) gives config and examples of it.
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
All matching requests **from** `${your-service-name}` have a timeout in `500ms`.


# Observability
Observability for microservices in EaseMesh can be cataloged into three areas, distributed tracing, metrics, and logging. Users can see the details of a request, such as the complete request path,  invocation dependencies, and latency of each sub-requests so that problems can be diagnosed. Metrics can reflect the health level of the system and summarize its state. Logging is used to provide more details based on the requested access.  

## Tracing
* Tracing is disabled by default in EaseMesh. It can be enabled dynamically during the lifetime of the mesh services. Currently, EaseMesh supports tracing these kinds of service protocols:

| Name           | Description                                                                                                                                                                                                           |
| -------------- | --------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| HTTP based RPC | Informations about communication between mesh service via HTTP protocol, such as latency, status code, request path and so on. Currently, EaseMesh supports tracing for `WebClient`, `RestTemplate` and `FeignClient` |
| JDBC           | Informations about MySQL SQL execution times, latency, successful rate and so on.                                                                                                                                     |
| Redis          | Informations about Redis command execution times, latency, successful rate and so on.                                                                                                                                 |
| RabbitMQ       | Informations about RabbitMQ command execution times, latency, successful rate and so on.                                                                                                                              |
| Kafka          | Informations about Kafka topics' read/write times, latency, successful rate and so on.                                                                                                                                |
| HTTP APIs      | Informations about the mesh service's HTTP APIs' latency, successful rate, status code and so on.                                                                                                                     |

* EaseMesh relies on `EaseAgent` for non-intrusive collecting metrics, and Kafka to store all collected tracing data. 
### Turn-on tracing
1. Configuring mesh service's [ObservabilityOutputServer](https://github.com/megaease/easemesh-api/blob/master/v1alpha1/meshmodel.proto#L312) to enable EaseMesh output tracing related data into Kafka. Modifying example YAML below, and name it with `outputservice.yaml` 
```yaml
 outputServer:
  enabled: true
  bootstrapServer: ${your_kafka_host_one}:9093,${your_kafka_host_two}:9093,${your_kafka_host_three}:9093
  timeout: 30000   
```
2. Updating mesh service's observability configurations with cmd 
```bash
$ ./easemesh/bin/meshctl service update ${your-service-name} observability -f ./outputservice.yaml
```
3. Finding the desired enable tracing service protocol in [ObservabilityTracings](https://github.com/megaease/easemesh-api/blob/master/v1alpha1/meshmodel.proto#L357) structure. For example, turning on the switch in `ObservabilityTracings.remoteInvoke`  can record mesh service's HTTP RPC tracing data. Also, EaseMesh allows users to configure how Java Agent should report tracing data, such as the reporting sample rate, reporting thread numbers in JavaAgent, and so on. **Note: the reporting configuration is global inside one mesh service's tracing** . Modifying example YAML below and name it with `enableRPCInvoke.yaml`   
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
4. Applying YAML with cmd 
```bash
$ ./easemesh/bin/meshctl service update ${your-service-name} observability -f ./enableRPCInvoke.yaml
```
5. Checking the web console for your mesh service's RPC tracing informations:
![tracing](../imgs/tracing.png)

### Turn-off tracing
1. If you want to disable tracing for one mesh service, then set this mesh service's global [tracing switch](https://github.com/megaease/easemesh-api/blob/master/v1alpha1/meshmodel.proto#L359) to `off`. For example prepare YAML as below and name it with `disableTracing.yaml`
```yaml
tracings:
  enabled: false             # The global enable switch
  sampleByQPS: 30            # The data above QPS 30 will be ignore 
  output:
    enabled: tru

    ......

```
Then applying it for your mesh service 
```bash
$ ./easemesh/bin/meshctl service update ${your-service-name} observability -f ./disableTracing.yaml
```
2. For only disabling one tracing feature for one mesh service, find the corresponding section, then turn off its switch is enough. For example, to shut down one mesh service's Redis tracing feature, you can prepare YAML as bellow and name it with `disableRedisTracing.yaml`
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
Applying YAML with cmd 
```bash
$ ./easemesh/bin/meshctl service update ${your-service-name} observability -f ./disableRedisTracing.yaml
```

## Metrics
* The EaseMesh uses EaseAgent(JavaAgent based on Java Byte buddy technology) to collect mesh services' basic metrics in a non-intrusive way. It will collect the data from a service perspective with very low CPU, memory, I/O resource usage. The supported metric types including:
* For the metric details for every type, checkout the EaseAgent's [develop-guide.md](https://github.com/megaease/easeagent/blob/master/doc/development-guide.md).


| Name              | Description                                                                                                                                                                                                                                                                                                 |
| ----------------- | ----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| HTTP request      | The mesh service's HTTP APIs metrics, such as m1/m5/m15 rate(`m1` indicates the The http request executions per second `exponentially-weighted moving average` in last 1 minute ), URL, total counts, error counts, p99(The http-request execution duration in milliseconds for 99% user) and so on.        |
| JDBC Statement    | The mesh service's JDBC statement metrics, such as the signature(used for complete SQL sentence matching), JDBC total count, JDBC m1 rate(The JDBC method executions per second `exponentially-weighted moving average` in last 1 minute.), Topn JDBC M1 error rate, JDBC P99 execution duration and so on. |
| JDBC Connection   | The mesh service's JDBC connection  metrics, such as URL, JDBC Connect total count, JDBC Connect M1 rate, JDBC Connect min execution duration, JDBC Connect min execution duration JDBC Connect P99 execution duration and so on                                                                            |
| JVM Memory        | The mesh service's JVM memory related metrics such as JVM initial memory, JVM used memory, JVM committed memory and JVM max memory.                                                                                                                                                                         |
| JVM GC            | The mesh service's JVM GC related metrics such as JVM Gc time, JVM Gc collection times and JVM Gc times per second.                                                                                                                                                                                         |
| Kafka Client      | The mesh service's Kafka client metics such as (no the reported target Kafka, the user's application usage's Kafka) topic name, Kafka producer throughput(M1), kafka consumer throughput(M1), producer min execution duration and so on.                                                                    |
| RabbitMq Producer | The mesh service's RabbitMQ producer's metics such as rabbit exchange, producer M1 rate, producer P99 execution duration and so on.                                                                                                                                                                         |
| RabbitMq Consumer | The mesh service's RabbitMQ consumer's metics such as rabbit exchange, consumer M1 rate, consumer P99 execution duration and so on.                                                                                                                                                                         |
| Redis             | The mesh service's Redis client's metics such as redis P25 execution duration, redis M1 count, redis P99 execution duration and so on.                                                                                                                                                                      |
| MD5 Dictionary    | The mesh service's JDBC statement's complete SQL sentences and MD5 values.                                                                                                                                                                                                                                  |
### Turn-on metrics reporting  
* EaseMesh also reports the mesh service's Metrics into the Kafka used by Tracing. So you can check out how to enable the output Kafka in the Tracing section. 

1. Finding the desired enable metrics type in [ObservabilityMetrics](https://github.com/megaease/easemesh-api/blob/master/v1alpha1/meshmodel.proto#L399) structure. For example, turning on switch in `ObservabilityMetrics.request`  can report mesh service's HTTP request-related metrics.Modifying example YAML below and name it with `enableHTTPRequest.yaml`  
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

2. Applying YAML with cmd 
```bash
$ ./easemesh/bin/meshctl service update ${your-service-name} observability -f ./enableHTTPRequest.yaml
```
3. Checking the web console for your mesh service's HTTP request metrics 
![metrics](../imgs/metrics.png)

### Turn-off metrics reporting
1. If you want to disable metrics reporting for one mesh service, then set this mesh service's global [metrics reporting switch](https://github.com/megaease/easemesh-api/blob/master/v1alpha1/meshmodel.proto#L390) to `off`. For example prepare YAML as below and name it with `disableMetrics.yaml`
```yaml
metrics:
  enabled: false             # The global enable switch
  access: 
    ......

```
Then applying it for your mesh service 
```bash
$ ./easemesh/bin/meshctl service update ${your-service-name} observability -f ./disableMetrics.yaml
```
2. For only disabling one type of metrics reporting for one mesh service, find the corresponding section, then turn off its switch is enough. For example, to shut down one mesh service's HTTP request metrics reporting, you can prepare YAML as bellow and name it with `disableHTTPReqMetrics.yaml`

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
then applying YAML with cmd 
```bash
$ ./easemesh/bin/meshctl service update ${your-service-name} observability -f ./disableHTTPReqMetrics.yaml
```

## Log
* Access log is also disabled by default in EaseMesh. The access log is used to recording details of HTTP APIs of mesh service been requested.  

### Turn-on Log
* EaseMesh also reports the mesh service's Logs into the Kafka used by Tracing. So you can check out how to enable the output Kafka in the Tracing section.  

1. Finding the `access` section in [ObservabilityMetrics](https://github.com/megaease/easemesh-api/blob/master/v1alpha1/meshmodel.proto#L399) structure. Turning on switch in `ObservabilityMetrics.access`  can enable access logging for mesh service's HTTP APIs.Modifying example YAML below and name it with `enableLog.yaml`  
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

2. Applying YAML with cmd 
```bash
$ ./easemesh/bin/meshctl service update ${your-service-name} observability -f ./enableLog.yaml
```
3. Checking the web console for your mesh service's HTTP log 
![metrics](../imgs/accesslog.png)


### Turn-off Log 
1. For disabling access logging for one mesh service, find the `access` section, then turn off its switch. For example, to shut down one mesh service's HTTP APIs logging, you can prepare YAML as bellow and name it with `disableLog.yaml`

```yaml
metrics:
  enabled: true                  # the global metrics reporting switch 
  access:
    enabled: false               # disable this service's logging
    interval: 30000
    topic: application-log
    ....
```
then applying YAML with cmd 
```bash
$ ./easemesh/bin/meshctl service update ${your-service-name} observability -f ./disableLog.yaml
```
 