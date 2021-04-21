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

![The topology of the spring-petclinic diagram](example/mesh-app-petclinic/backgroud/microservices-architecture-diagram.jpg)


### Environment require
#### Infrastructure version
* Linux kernel version 4.15+
* Kubernetes version 1.18+
* Mysql version 14.14+
####  Dependence check
1. Run cmd `kubectl get nodes` to make sure your k8s cluster is healthy. 
2. Run cmd `mysql -u$your_db_user -p$your_db_pass` to make sure application can connect to db successfully. 

### Start PetClinic in Easemesh with K8s:

1. Enter `./example/mesh-app-petclinic` dir, execute `./install.sh `
2. Using the db table schemes and records provided in [PetClinic example](https://github.com/spring-projects/spring-petclinic/tree/main/src/main/resources/db/mysql) to set up yours.
3. Run `kubectl get svc mesh-ingress `
easemesh will create a k8s `NodePort` type service for easemesh IngressGateway. Configure it into your traffic gateway's routing address,e.g., configure NGINX with
```
location /pet/ {
    proxy_pass http://$NodePortIP:$NodePortNum/;
    ...
}

```
4. Open browser with `$your_domain/pet/#!/welcome`, should see the welcome page of PetClinic. 

### Canary deployment

![EaseMesh Canary topology](example/mesh-app-petclinic/backgroud/canary-demo-diagram.png)

1. Colored your traffic with HTTP header `X-Canary: lv1`. This can be done by using NGINX's `proxy_set_header` cmd or with the help of other tools.
NGINX Example:
```
location /pet/ {
    set $canary_header "";
    if ($http_user_agent ~ Firefox) {
        set $canary_header "lv1";
    }

    if ($http_user_agent ~ Chrome) {
        set $canary_header "";
    }
...
```
* if user's browser is Firefox, it will be routed into canary version. The chrome user will visit the original page. 

2. Developing a canary customer service version, we add extra process to the city field of the customer data. The change can be checked via [this commitment](https://github.com/akwei/spring-petclinic-cloud/commit/3be54a2c7e63c955990cbc1e78dab029b516a3ec)

2. Deploy canary version's customer service with cmd `kubectl apply -f  ./example/mesh-app-petclinic/canary/customers-service-deployment-canary.yaml`
3. Open chrome with `$your_domain/pet/#!/owners`, the owner info page remained the same.
4. Visit the same URL with Firefox, should see the table with brand new city field which will be added "-US" suffix into every record. 
