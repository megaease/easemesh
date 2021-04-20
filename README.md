# easemesh
A service mesh implementation for connecting, secure, control, and observe services of spring-cloud.

## Install 
```
cd ./install
./egctl mesh install
```
* It will register K8s mesh-deployment CRD, and start easemesh control-plane, IngressGateway. 

## Quick Start
### Background
* SprintCloud PetClinic  [github link](https://github.com/spring-petclinic/spring-petclinic-cloud) micro service example.
* It uses Spring Cloud Gateway, Spring Cloud Circuit Breaker, Spring Cloud Config, Spring Cloud Sleuth, Resilience4j, Micrometer and Eureka Service Discovery from Sprint Cloud Netflix technology stack.

### Environment require
* Linux kernel version 4.15+
* Kubernetes version 1.18+
* Mysql version 14.14+

### Start PetClinic in Easemesh with K8s:
1. Enter `./example/script` dir, execute `./install.sh `
2. In `./example/script` dir, execute `./init_db.sh` to configure database which is mysql with contents and tables provided in [PetClinic example](https://github.com/spring-projects/spring-petclinic/tree/main/src/main/resources/db/mysql).
3. Open browser with `$your_domain/pet/#!/welcome`, should see the welcome page of PetClinic. 

### Canary deployment
1. Colored your traffic with HTTP header `X-Canary: lv1`. This can be done by setting NGINX's `proxy_set_header` or with the help of other tools.
2. Deploy canary version's customer service with cmd `kubectl apply -f  ./example/mesh-app-petclinic/k8s/customers-service-deployment-canary.yaml`
3. Open browser with `$your_domain/pet/#!/welcome`, go to customer page, should see the table with brand new city field. 
