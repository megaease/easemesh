## Background

* EaseMesh uses Easegress-based sidecar inside Kubernetes Pod for **traffic hosting** and  EaseAgent for metrics reporting and RESTful-API-based RPC enhancement. 
* EaseMesh only supports Java SpringCloud ecosystem's application natively currently.

### EaseMesh traffic hosting

There are three types of traffic that are managed by EaseMesh. 

* First, the **RESTful-API HTTP traffic** for RPC inside the mesh. This traffic is invoked by Java applications with popular RPC frameworks, such as Feign, RestTemplate, and so on. EaseAgent will enhance this traffic by adding the target RPC server's name inside the HTTP header for telling the sidecar of the real handler.
* Second, the **Health-checking HTTP traffic**. This traffic is sent from the sidecar to the Java application's additional port opened by EaseAgent.  The complete URI is `http://localhost:9000/health` by default. This `9000` port is opened by EaseAgent, sidecar will query this URI period for checking the liveness of the Java application. After successfully deployed, sidecar will registry this instance into EaseMesh automatically after confirming the HTTP 200 success return by this URI.
* Third, the **Service-discovery traffic**. This traffic is invoked by the Java spring cloud application's RPC framework. During the lifetime of the Java application, sidecar will work as the Java application's service registry and discovery center. EaseMesh sidecar implements Eureka/Consul/Naocs APIs for hosting the Java application's registry and discovery requests. To make the sidecar server the registry and discovery center, value it with `http://localhost:13009` inside the Java application's  XML. The port `13009` is listened by sidecar for handling Eureka/Consul/Nacos APIs. 

The ports used by EaseMesh sidecar+agnet system

| Role    | Port  | Description                                                                                                                 |
| ------- | ----- | --------------------------------------------------------------------------------------------------------------------------- |
| Sidecar | 13001 | The default Ingress port listened by sidecar for handing over traffic to local Java application                             |
| Sidecar | 13002 | The default egress port listened by sidecar for routing local Java applications RPC request to another Java application     |
| Sidecar | 13009 | The default registry and discovery port listened by sidecar, for handling local Java application's Eureka/Conslu/Nacos APIs |
| Agent   | 9000  | The default health port listened by Agent queried by sidecar for checking the liveness of Java application                  |



## Problem

* Figuring out the standard for supporting multiple-language programs running inside EaseMesh.

### Analysis

* To support the none-Java-spring-cloud-based RESTful-API application, we had demoed a DNS-enhancement way for supporting Java spring boot application. Can we reuse this way to support Golang-based or RUST-based  RESTful-API applications? 

[![](https://mermaid.ink/img/eyJjb2RlIjoic2VxdWVuY2VEaWFncmFtXG4gICAgSmF2YUFQUCAtPj4gK2NvcmVETlMgOiBhc2tpbmcgdGhlIGRvbWFpbiBhbmFseXNpc1xuICAgIGNvcmVETlMgLT4-ICtFdGNkIDogc2VhcmNoIHNlcnZpY2UgaW4gRWFzZU1lc2ggRXRjZFxuICAgIEV0Y2QgLT4-IC1jb3JlRE5TIDogcmV0dXJuIGxvY2FsIHNpZGVjYXIgYWRkcmVzcyBpZiBpdCdzIGEgbWVzaCBzZXJ2aWNlc1xuICAgIGNvcmVETlMgLT4-ICAtSmF2YUFQUCA6IHJldHVybiB0aGUgbG9jYWwgc2lkZWNhciBhZGRyXG4gICAgSmF2YUFQUCAtPj4gK2xvY2FsU2lkZWNhciA6IFJFU1RmdWwgcmVxdWVzdFxuICAgIGxvY2FsU2lkZWNhciAtPj4gK3RhcmdldFNpZGVjYXIgOiByb3V0aW5nIHRvIHRhcmdldCBzZXJ2ZXIncyBzaWRlY2FyXG4gICAgdGFyZ2V0U2lkZWNhciAtPj4gK3RhcmdldEphdmFBUFA6IHJvdXRpbmcgdG8gdGhlIHJlYWwgaGFuZGxlclxuICAgIHRhcmdldEphdmFBUFAgLT4-IC10YXJnZXRTaWRlY2FyOiByZXR1cm4gdGhlIHJlc3VsdFxuICAgIHRhcmdldFNpZGVjYXIgLT4-IC1sb2NhbFNpZGVjYXI6IHJldHVybiB0aGUgcmVzdWx0XG4gICAgbG9jYWxTaWRlY2FyIC0-PiAtSmF2YUFQUCA6IHJldHVybiB0aGUgcmVzc3VsdCIsIm1lcm1haWQiOnsidGhlbWUiOiJkZWZhdWx0In0sInVwZGF0ZUVkaXRvciI6ZmFsc2UsImF1dG9TeW5jIjp0cnVlLCJ1cGRhdGVEaWFncmFtIjpmYWxzZX0)](https://mermaid-js.github.io/mermaid-live-editor/edit/#eyJjb2RlIjoic2VxdWVuY2VEaWFncmFtXG4gICAgSmF2YUFQUCAtPj4gK2NvcmVETlMgOiBhc2tpbmcgdGhlIGRvbWFpbiBhbmFseXNpc1xuICAgIGNvcmVETlMgLT4-ICtFdGNkIDogc2VhcmNoIHNlcnZpY2UgaW4gRWFzZU1lc2ggRXRjZFxuICAgIEV0Y2QgLT4-IC1jb3JlRE5TIDogcmV0dXJuIGxvY2FsIHNpZGVjYXIgYWRkcmVzcyBpZiBpdCdzIGEgbWVzaCBzZXJ2aWNlc1xuICAgIGNvcmVETlMgLT4-ICAtSmF2YUFQUCA6IHJldHVybiB0aGUgbG9jYWwgc2lkZWNhciBhZGRyXG4gICAgSmF2YUFQUCAtPj4gK2xvY2FsU2lkZWNhciA6IFJFU1RmdWwgcmVxdWVzdFxuICAgIGxvY2FsU2lkZWNhciAtPj4gK3RhcmdldFNpZGVjYXIgOiByb3V0aW5nIHRvIHRhcmdldCBzZXJ2ZXIncyBzaWRlY2FyXG4gICAgdGFyZ2V0U2lkZWNhciAtPj4gK3RhcmdldEphdmFBUFA6IHJvdXRpbmcgdG8gdGhlIHJlYWwgaGFuZGxlclxuICAgIHRhcmdldEphdmFBUFAgLT4-IC10YXJnZXRTaWRlY2FyOiByZXR1cm4gdGhlIHJlc3VsdFxuICAgIHRhcmdldFNpZGVjYXIgLT4-IC1sb2NhbFNpZGVjYXI6IHJldHVybiB0aGUgcmVzdWx0XG4gICAgbG9jYWxTaWRlY2FyIC0-PiAtSmF2YUFQUCA6IHJldHVybiB0aGUgcmVzc3VsdCIsIm1lcm1haWQiOiJ7XG4gIFwidGhlbWVcIjogXCJkZWZhdWx0XCJcbn0iLCJ1cGRhdGVFZGl0b3IiOmZhbHNlLCJhdXRvU3luYyI6dHJ1ZSwidXBkYXRlRGlhZ3JhbSI6ZmFsc2V9)

* To support non-registry-discovery dependent on Java spring boot application, EaseMesh enhances Kubernetes' coreDNS with add a plugin for finding services inside EaseMesh's Etcd. We can reuse this method for none-Java-based programs. 
* EaseAgent uses Java Byte Buddy-based technology for collecting several application metrics. This requires a JVM-liked software architecture. This observability will be sacrificed for the none-Java-spring-cloud-based RESTful-API application.

### Protocol

To support the none-Java-spring-cloud-based RESTful-API application, regardless of which programming is used. The application must follow the protocol below


1. It must serve as standard RESTful-API for handling requesting or invoking RPC. 

2. It must use a domain for discovering in RESTful-API RPC.
   > Requirement: use coreDNS with easemesh specific plugin
   >              allowed domain format
   >                1: only service name
   >                2: regex rule:  ^(|(\w+\.)+)vet-services\.(\w+)\.svc\..+$
   > ​                 e.g.  _tcp.vet-services.easemesh.svc.cluster.local
   > ​                       vet-services.easemesh.svc.cluster.local
   > ​                       _zip._tcp.vet-services.easemesh.svc.com

3. It must serve the `http://localhost:9000/health` URI for EaseMesh health checking. (Only HTTP 200 return is required, regardless of the body content)

4. It must reserve ports `13001` , `13002` and `13009` for local sidecar usage.

If an application obeys the protocol above, then EaseMesh can run it inside with sacrificed observability regardless of the implements programming language.
