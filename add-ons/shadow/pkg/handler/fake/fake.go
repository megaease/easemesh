package fake

import (
	"bytes"

	"github.com/megaease/easemesh/mesh-shadow/pkg/object"
	"github.com/megaease/easemeshctl/cmd/client/resource"
	appsV1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	k8Yaml "k8s.io/apimachinery/pkg/util/yaml"
)

const fakeCanaryYaml = `
apiVersion: mesh.megaease.com/v1alpha1
kind: ServiceCanary
metadata:
  name: shadow-service-canary
spec:
  priority: 5 
  selector:
    matchServices: [service1, service2, service3]
    matchInstanceLabels: {version: shadow}
  trafficRules:
    headers:
      X-Mesh-Shadow:
        exact: shadow
`

const fakeDeleteCanaryYaml = `
apiVersion: mesh.megaease.com/v1alpha1
kind: ServiceCanary
metadata:
  name: shadow-service-canary
spec:
  priority: 5 
  selector:
    matchServices: [service1, service2]
    matchInstanceLabels: {version: shadow}
  trafficRules:
    headers:
      X-Mesh-Shadow:
        exact: shadow
`
const sourceDeploymentYaml = `
apiVersion: apps/v1
kind: Deployment
metadata:
  annotations:
    mesh.megaease.com/service-name: visits-service
  name: visits-service
  namespace: spring-petclinic
spec:
  replicas: 1
  selector:
    matchLabels:
      app: visits-service
  strategy:
    rollingUpdate:
      maxSurge: 25%
      maxUnavailable: 25%
    type: RollingUpdate
  template:
    metadata:
      labels:
        app: visits-service
    spec:
      containers:
      - args:
        - -c
        - java -server -Xmx1024m -Xms1024m -Dspring.profiles.active=sit -Djava.security.egd=file:/dev/./urandom  org.springframework.boot.loader.JarLauncher
        command:
        - /bin/sh
        env:
        - name: JAVA_TOOL_OPTIONS
          value: ' -javaagent:/agent-volume/easeagent.jar -Deaseagent.log.conf=/agent-volume/log4j2.xml '
        image: megaease/spring-petclinic-visits-service:latest
        imagePullPolicy: Always
        lifecycle:
          preStop:
            exec:
              command:
              - sh
              - -c
              - sleep 10
        name: visits-service
        ports:
        - containerPort: 8080
          protocol: TCP
        resources:
          limits:
            cpu: "2"
            memory: 1Gi
          requests:
            cpu: 200m
            memory: 256Mi
        volumeMounts:
        - mountPath: /application/application-sit.yml
          name: configmap-volume-0
          subPath: application-sit.yml
        - mountPath: /agent-volume
          name: agent-volume
      - command:
        - /bin/sh
        - -c
        - /opt/easegress/bin/easegress-server -f /sidecar-volume/sidecar-config.yaml
        env:
        - name: APPLICATION_IP
          valueFrom:
            fieldRef:
              apiVersion: v1
              fieldPath: status.podIP
        image: 172.20.2.189:5001/megaease/easegress:server-sidecar
        imagePullPolicy: IfNotPresent
        name: easemesh-sidecar
        ports:
        - containerPort: 13001
          name: sidecar-ingress
          protocol: TCP
        - containerPort: 13002
          name: sidecar-egress
          protocol: TCP
        - containerPort: 13009
          name: sidecar-eureka
          protocol: TCP
        volumeMounts:
        - mountPath: /sidecar-volume
          name: sidecar-volume
      initContainers:
      - command:
        - sh
        - -c
        - "set -e\ncp -r /easeagent-volume/* /agent-volume\n\necho 'name: visits-service\ncluster-join-urls:
          http://easemesh-controlplane-svc.easemesh:2380\ncluster-request-timeout:
          10s\ncluster-role: reader\ncluster-name: easemesh-control-plane\nlabels:\n
          \ alive-probe: http://localhost:9900/health\n  application-port: 8080\n
          \ mesh-service-labels: \n  mesh-servicename: visits-service\n' > /sidecar-volume/sidecar-config.yaml"
        image: 172.20.2.189:5001/megaease/easeagent-initializer:latest
        imagePullPolicy: IfNotPresent
        name: initializer
        volumeMounts:
        - mountPath: /agent-volume
          name: agent-volume
        - mountPath: /sidecar-volume
          name: sidecar-volume
      volumes:
      - configMap:
          items:
          - key: application-sit-yml
            path: application-sit.yml
          name: visits-service
        name: configmap-volume-0
      - emptyDir: {}
        name: agent-volume
      - emptyDir: {}
        name: sidecar-volume
`

const shadowDeploymentYaml = `
apiVersion: apps/v1
kind: Deployment
metadata:
  annotations:
    mesh.megaease.com/service-labels: version=shadow
    mesh.megaease.com/service-name: visits-service
    mesh.megaease.com/shadow-service-name: shadow-visits-service
  labels:
    mesh.megaease.com/shadow-service: "true"
  name: visits-service-shadow
  namespace: spring-petclinic
spec:
  replicas: 1
  selector:
    matchLabels:
      app: visits-service
      mesh.megaease.com/shadow-service: "true"
  strategy:
    rollingUpdate:
      maxSurge: 25%
      maxUnavailable: 25%
    type: RollingUpdate
  template:
    metadata:
      creationTimestamp: null
      labels:
        app: visits-service
        mesh.megaease.com/shadow-service: "true"
    spec:
      initContainers: []
      containers:
      - args:
        - -c
        - java -server -Xmx1024m -Xms1024m -Dspring.profiles.active=sit -Djava.security.egd=file:/dev/./urandom  org.springframework.boot.loader.JarLauncher
        command:
        - /bin/sh
        env:
        - name: JAVA_TOOL_OPTIONS
          value: ' -javaagent:/agent-volume/easeagent.jar -Deaseagent.log.conf=/agent-volume/log4j2.xml '
        - name: EASEMESH_TAGS
          value: '{"label.local":"shadow"}'
        image: megaease/spring-petclinic-visits-service:latest
        imagePullPolicy: Always
        lifecycle:
          preStop:
            exec:
              command:
              - sh
              - -c
              - sleep 10
        name: visits-service
        ports:
        - containerPort: 8080
          protocol: TCP
        resources:
          limits:
            cpu: "2"
            memory: 1Gi
          requests:
            cpu: 200m
            memory: 256Mi
        volumeMounts:
        - mountPath: /application/application-sit.yml
          name: configmap-volume-0
          subPath: application-sit.yml
      volumes:
      - configMap:
          items:
          - key: application-sit-yml
            path: application-sit.yml
          name: visits-service
        name: configmap-volume-0
`

// NewServiceCanary create fake ServiceCanary for test.
func NewServiceCanary() *resource.ServiceCanary {
	decoder := k8Yaml.NewYAMLOrJSONDecoder(bytes.NewReader([]byte(fakeCanaryYaml)), 1000)
	serviceCanary := &resource.ServiceCanary{}
	_ = decoder.Decode(serviceCanary)
	return serviceCanary
}

// NewDeletedServiceCanary create fake ServiceCanary for test.
func NewDeletedServiceCanary() *resource.ServiceCanary {
	decoder := k8Yaml.NewYAMLOrJSONDecoder(bytes.NewReader([]byte(fakeDeleteCanaryYaml)), 1000)
	serviceCanary := &resource.ServiceCanary{}
	_ = decoder.Decode(serviceCanary)
	return serviceCanary
}

// NewShadowService create fake ShadowService for test.
func NewShadowService() object.ShadowService {
	shadowService := object.ShadowService{
		Name:        "shadow-visits-service",
		ServiceName: "visits-service",
		Namespace:   "spring-petclinic",
	}
	return shadowService
}

// NewNamespace create fake NameSpace for test.
func NewNamespace() *corev1.Namespace {
	ns1 := &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: "spring-petclinic",
		},
	}
	return ns1
}

// NewSourceDeployment create fake SourceDeployment for test.
func NewSourceDeployment() *appsV1.Deployment {
	decoder := k8Yaml.NewYAMLOrJSONDecoder(bytes.NewReader([]byte(sourceDeploymentYaml)), 1000)
	sourceDeployment := &appsV1.Deployment{}
	_ = decoder.Decode(sourceDeployment)
	return sourceDeployment
}

// NewShadowDeployment create fake Deployment for test.
func NewShadowDeployment() *appsV1.Deployment {
	decoder := k8Yaml.NewYAMLOrJSONDecoder(bytes.NewReader([]byte(shadowDeploymentYaml)), 1000)
	clonedDeployment := &appsV1.Deployment{}
	_ = decoder.Decode(clonedDeployment)
	return clonedDeployment
}
