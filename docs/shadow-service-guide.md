# Full Stack Stress With Shadow Service

With the increasing of hardware performance, network bandwidth and user data, traditional stand-alone applications are being replaced by software systems based on the network.

But while bringing more powerful computing capability, such a software system also introduces complexities which were never found in a stand-alone application. Today, a typical software system can include tens to tens of thousands modules. Moreover, to make the development and deployment process more faster, these modules will be developed by different teams in different languages, which also makes the communication between them more complicated.

At the same time, the business logic has also undergone great changes compared with the past, for example, the promotion of Black Friday will make the system bear several times or even dozens of times the daily pressure. As a result, stress testing is becoming increasingly important, but traditional testing methods are also becoming less and less adaptable.

## Issues with a dedicated test environment

Using a 1:1 test environment is a very good stress test method in the stand-alone era, but becomes very impractical in the age of the Internet.

The first issue is cost. To test a stand-alone application, the cost of a dedicated testing computer is affordable for an independent developer in most cases; but today's system requires too many resources, servers, bandwitdh, electricity and server room, the cost of building an 1:1 test environment according to the production environment is beyond the affordability of most companies.

Even if the company is rich enough for an 1:1 test environment, it is challenging to keep this test environment completely consistent with the production environment.

Because it is a test environment, people will keep deploy the test versions of their modules on it, so there will inevitably be programmers who have not restored it to the production version after testing their own modules, and in the long run, the difference between the test environment and the production environment will become larger and larger, resulting in distorted test results.

Moreover, because the cost is so high, I don't think any company can trench up such a test environment for each development group, that is, everyone must share the same environment. At this time, if there is no excellent internal coordination mechanism, the tests conducted by multiple development groups at the same time will also affect the test results.

There is also the problem of test data, if there is no guarantee that the data of the test environment is close to or even the same as the production environment, the test results will not be trusted. For example, a Weibo-like system, ordinary users like me generally only have a few dozen or hundreds of followers, so I send a message and do whatever I want to easily notify all followers. But for a million big V, the situation will be very different. Therefore, we cannot simply use simulated data for testing.

In order for the test results to be realistic and reliable, it is best to import complete production data into the test environment. At first glance, this only requires a simple backup to restore the database, but the production environment contains a lot of sensitive data, and copying it to the test environment at will will undoubtedly greatly increase the risk of data leakage.

## Testing issues in a production environment

Because there are many difficulties in using the test environment for stress testing, people have turned their attention to the production environment and tried to directly use the pressure trough period on the production environment for testing. But it's an intrusive solution that involves modifying or even redefining business logic, so it's also a huge challenge.

As shown in the following figure, the blue box is the original business logic, and the orange box is the new logic to achieve this scenario. On the surface, these logics only need to add a few if/else to implement, but in reality they are much more complex.

![diagram-01](./imgs/shadow-service-guide-01.png)

Assuming that we want to modify an online shopping system, the process of user shopping and placing orders should involve a series of modules such as users, orders, and payments.

If we were to modify the user module, in that diamond box, how would we be able to tell whether we should go through the test logic or the production logic? The more common method is to specify a range of user IDs in advance, if it is a user of this range, go to the test logic, otherwise take the production logic.

After the user module, the logic goes to the order module, at this time, we may still want to judge whether the test logic should be taken through the user ID, but the actual situation may be: after a series of complex processing processes, the order module can not see the user information at all, so this road is not passable.

In order to distinguish between normal orders and test orders, the user module is required to perform special processing in the orange box, such as adding special marks to the order number. However, in a complex system, it is not easy for the user module to know all the modules that the subsequent process will go through, so in order not to affect the normal production logic, it takes a lot of effort to ensure the normal transmission of the test state, not to mention whether to access different data sets, whether to simulate third-party services, and so on.

Obviously, the amount of work required for this modification of business logic is proportional to the number of function points. But in addition to the huge workload, the more serious problem is that after the hard work of revision, who can guarantee that all the changes that need to be made have been changed and correct? And in the event of an omission or error, the risk of destroying the production system is too great.

## The Solution

As can be seen from the previous analysis, traditional test methods are either costly and do not have accurate data, or are heavy and risk of disrupting production systems. Therefore, MegaEase believes that to solve the problem of full-site stress testing in the network era, a completely new approach must be used, and the key to this approach lies in "three consistency" and "four isolation".

Three-consistency refers to business consistency, data consistency and resource consistency. That is to say, the test system and the production system should be exactly the same, only in this way can accurate test data be obtained. Realistically, 100% consistency is not easy to do, for example, we usually can't ask a third party to cooperate with us to test, so we can only use a simulated method to replace some third-party dependencies. But we still need to do the best possible to ensure the consistency of the two systems.

Four isolations refer to service isolation, data isolation, traffic isolation, and resource isolation. These isolations are all designed to completely separate the production and test systems and avoid their mutual influence. 

Obviously, the three consistency solves the problem of the accuracy of the test results, while the four isolation ensures that the test process does not affect the production system.

Based on the above concept, MegaEase implements the Shadow Service feature in EaseMesh. Using this feature, users can easily create a replica of all services in the system, except for the shadow tag, these replicas are exactly the same as the original service, thus ensuring business consistency and business isolation. At the same time, Shadow Service also automatically creates a Canary rule that forwards requests with `X-Mesh-Shadow: shadow` headers as test requests to the service replica and sends other requests to the original service for traffic isolation.

In terms of data, Shadow Service can replace the connection information of various middleware including MySQL, Kafka, Redis, etc. according to the configuration, and change the sending destination of data requests, which ensures data isolation. The user can directly copy the production data as test data to ensure data consistency.

Resource consistency and resource isolation mean that the test system should use the same resources as the production system specifications, but should not share the same set of resources. It's mostly a hardware issue, but Kubernetes has given a very good answer at the software level, and EaseMesh is built on top of Kubernetes, so by deploying a copy of the service into a new POD, resource consistency and resource isolation are guaranteed.

## How to use shadow service

Below we use an order payment scenario to introduce the specific use of Shadow Service. This scenario involves three services, User, Order, and Payment, and the User and Order Services use their own MySQL databases, as shown in the following figure (where the payment service eventually calls a third-party service to complete the payment, but is not shown in the figure).

![diagram-02](./imgs/shadow-service-guide-02.png)

To test it, we need to first create two copies of the database:

![diagram-03](./imgs/shadow-service-guide-03.png)

Then, use the following emctl command to create a shadow copy of the user and order service, a Canary rule (created automatically, so not included in the configuration below), and the Mock payment service. Note that when we created a copy of the User and Orders service, we used the database that pointed to the database copy that we just created.

```bash
echo '
kind: ShadowService
apiVersion: mesh.megaease.com/v1alpla1
metadata:
  name: shadow-user-service
spec:
  serviceName: user-service
  namespace: megaease-mall
  mysql:
    uris: "jdbc:mysql://172.20.2.216:3306/shadow_user_db..."
    userName: "megaease"
    password: "megaease"

---

kind: ShadowService
apiVersion: mesh.megaease.com/v1alpla1
metadata:
  name: shadow-order-service
spec:
  serviceName: order-service
  namespace: megaease-mall
  mysql:
    uris: "jdbc:mysql://172.20.2.216:3306/shadow_order_db..."
    userName: "megaease"
    password: "megaease"

---

kind: Mock
apiVersion: mesh.megaease.com/v1alpha1
metadata:
  name: payment-service
  registerTenant: megaease-mall
spec:
  enabled: false
  rules:
    - match:
        pathPrefix: /
        headers:
          X-Mesh-Shadow:
            exact: shadow
      code: 200
      headers:
        Content-Type: application/json
      body: '{"result":"succeeded"}'
' | emctl apply
```

This completes the creation of the entire Shadow Service, and the final system architecture looks like this:

![diagram-04](./imgs/shadow-service-guide-04.png)
## Pros of shadow service

Using Shadow Service for stress testing has the following advantages, in addition to the already mentioned results being accurate and not affecting the production system:

* 0 code modification: all through the configuration is complete, no need to modify any code, there is no risk of bugs.
* Low cost: In the case of using a cloud server, the hardware resources used for testing are applied for with the application, and they are released when they are used up, and only need to pay for the actual use period.
*Secure: Although production data is used during testing, the test system and the production system are in the same security domain, so there is no increased risk of data leakage.
