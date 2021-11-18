package handler

import (
	"bytes"
	"reflect"
	"testing"

	"github.com/megaease/easemesh/mesh-shadow/pkg/object"
	appsV1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	k8Yaml "k8s.io/apimachinery/pkg/util/yaml"
	"k8s.io/client-go/kubernetes"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func fakeShadowService() object.ShadowService {
	shadowService := object.ShadowService{
		Name:        "shadow-visits-service",
		ServiceName: "visits-service",
		Namespace:   "spring-petclinic",
	}
	return shadowService
}
func fakeNameSpace() *corev1.Namespace {
	ns1 := &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: "spring-petclinic",
		},
	}
	return ns1
}

func fakeDeployment() *appsV1.Deployment {
	decoder := k8Yaml.NewYAMLOrJSONDecoder(bytes.NewReader([]byte(sourceDeploymentYaml)), 1000)
	sourceDeployment := &appsV1.Deployment{}
	_ = decoder.Decode(sourceDeployment)
	return sourceDeployment
}

func fakeClonedDeployment() *appsV1.Deployment {
	decoder := k8Yaml.NewYAMLOrJSONDecoder(bytes.NewReader([]byte(clonedDeploymentYaml)), 1000)
	clonedDeployment := &appsV1.Deployment{}
	_ = decoder.Decode(clonedDeployment)
	return clonedDeployment
}

func TestShadowServiceCloner_cloneDeploymentSpec(t *testing.T) {
	type fields struct {
		KubeClient    kubernetes.Interface
		RunTimeClient *client.Client
	}

	deployment := fakeDeployment()
	shadowService := fakeShadowService()
	clonedDeployment := fakeClonedDeployment()
	type args struct {
		sourceDeployment *appsV1.Deployment
		shadowService    *object.ShadowService
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *appsV1.Deployment
	}{
		{
			name:   "test",
			fields: fields{},
			args: args{
				sourceDeployment: deployment,
				shadowService:    &shadowService,
			},
			want: clonedDeployment,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cloner := &ShadowServiceCloner{
				KubeClient:    tt.fields.KubeClient,
				RunTimeClient: tt.fields.RunTimeClient,
			}
			got := cloner.cloneDeploymentSpec(tt.args.sourceDeployment, tt.args.shadowService)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("cloneDeploymentSpec() = %v, \n want %v", got, tt.want)
			}
		})
	}
}

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

const clonedDeploymentYaml = `
apiVersion: apps/v1
kind: Deployment
metadata:
  annotations:
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
