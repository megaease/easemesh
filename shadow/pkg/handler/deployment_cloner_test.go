package handler

import (
	"bytes"
	"reflect"
	"testing"

	"github.com/megaease/easemesh/mesh-shadow/pkg/object"
	appsV1 "k8s.io/api/apps/v1"
	coreV1 "k8s.io/api/core/v1"
	k8Yaml "k8s.io/apimachinery/pkg/util/yaml"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
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
func fakeDeployment() *appsV1.Deployment{
	decoder := k8Yaml.NewYAMLOrJSONDecoder(bytes.NewReader([]byte(sourceDeploymentYaml)), 1000)
	sourceDeployment := &appsV1.Deployment{}
	_ = decoder.Decode(sourceDeployment)
	return sourceDeployment
}

func fakeClonedDeployment() *appsV1.Deployment{
	decoder := k8Yaml.NewYAMLOrJSONDecoder(bytes.NewReader([]byte(clonedDeploymentYaml)), 1000)
	clonedDeployment := &appsV1.Deployment{}
	_ = decoder.Decode(clonedDeployment)
	return clonedDeployment
}

func fakeSourceVolumes() []coreV1.Volume{
	deployment := fakeDeployment()
	return deployment.Spec.Template.Spec.Volumes
}

func fakeInitContainers() []coreV1.Container{
	deployment := fakeDeployment()
	return deployment.Spec.Template.Spec.InitContainers
}

func fakeContainers() []coreV1.Container{
	deployment := fakeDeployment()
	return deployment.Spec.Template.Spec.Containers
}

//
func TestShadowServiceCloner_cloneDeploymentSpec(t *testing.T) {
	type fields struct {
		KubeClient    kubernetes.Interface
		RunTimeClient *client.Client
		CRDClient     *rest.RESTClient
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
			name: "test",
			fields: fields{},
			args: args{
				sourceDeployment: deployment,
				shadowService: &shadowService,
			},
			want: clonedDeployment,

		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cloner := &ShadowServiceCloner{
				KubeClient:    tt.fields.KubeClient,
				RunTimeClient: tt.fields.RunTimeClient,
				CRDClient:     tt.fields.CRDClient,
			}
			got := cloner.cloneDeploymentSpec(tt.args.sourceDeployment, tt.args.shadowService)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("cloneDeploymentSpec() = %v, \n want %v", got, tt.want)
			}
		})
	}
}

//
// func Test_findContainer(t *testing.T) {
//
// 	deployment := fakeDeployment()
// 	type args struct {
// 		containers    []coreV1.Container
// 		containerName string
// 	}
//
// 	tests := []struct {
// 		name  string
// 		args  args
// 		want  *coreV1.Container
// 		want1 bool
// 	}{
// 		// TODO: Add test cases.
// 		{
// 			name: "visits-service",
// 			args: args{
// 				containers: deployment.Spec.Template.Spec.Containers,
// 				containerName: "visits-service",
// 			},
// 			want: &deployment.Spec.Template.Spec.Containers[0],
// 			want1: true,
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			container, ok := findContainer(tt.args.containers, tt.args.containerName)
// 			if !reflect.DeepEqual(container, tt.want) {
// 				t.Errorf("findContainer() got = %v, want %v", container, tt.want)
// 			}
// 			if ok != tt.want1 {
// 				t.Errorf("findContainer() got1 = %v, want %v", ok, tt.want1)
// 			}
// 		})
// 	}
// }
//
//
// func Test_injectShadowLabels(t *testing.T) {
// 	labels := map[string]string{}
// 	tests := []struct {
// 		name string
// 		args map[string]string
// 	}{
// 		// TODO: Add test cases.
// 		{
// 			name: "test",
// 			args: labels,
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			injectShadowLabels(tt.args)
// 			value, ok := labels[shadowLabelKey]
// 			if !ok || value  != "true"  {
// 				t.Errorf("injectShadowLabels() failed")
// 			}
// 		})
// 	}
// }
//
//
// func Test_shadowContainers(t *testing.T) {
// 	container := &coreV1.Container{}
// 	k8Yaml.Unmarshal([]byte(sourceContainerYaml), container)
//
// 	containers := fakeContainers()
// 	tests := []struct {
// 		name string
// 		args []coreV1.Container
// 		want []coreV1.Container
// 	}{
// 		{
// 			name: "test",
// 			args: containers,
// 			want: []coreV1.Container{*container},
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			if got := shadowContainers(tt.args); !reflect.DeepEqual(got, tt.want) {
// 				t.Errorf("shadowContainers() = %v, want %v", got, tt.want)
// 			}
// 		})
// 	}
// }
//
// func Test_shadowInitContainers(t *testing.T) {
// 	initContainers := fakeInitContainers()
// 	tests := []struct {
// 		name string
// 		args []coreV1.Container
// 		want []coreV1.Container
// 	}{
// 		{
// 			name: "test",
// 			args: initContainers,
// 			want: []coreV1.Container{},
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			if got := shadowInitContainers(tt.args); !reflect.DeepEqual(got, tt.want) {
// 				t.Errorf("shadowInitContainers() = %v, want %v", got, tt.want)
// 			}
// 		})
// 	}
// }
//
// func Test_shadowName(t *testing.T) {
//
// 	testName1 := "visit-service"
// 	testName2 := "vet-service"
// 	testName3 := "custom-service"
// 	tests := []struct {
// 		name string
// 		want string
// 	}{
// 		{
// 			name: testName1,
// 			want: testName1 + shadowDeploymentNameSuffix,
// 		},
// 		{
// 			name: testName2,
// 			want: testName2 + shadowDeploymentNameSuffix,
// 		},
// 		{
// 			name: testName3,
// 			want: testName3 + shadowDeploymentNameSuffix,
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			if got := shadowName(tt.name); got != tt.want {
// 				t.Errorf("shadowName() = %v, want %v", got, tt.want)
// 			}
// 		})
// 	}
// }
//
// func Test_shadowVolumes(t *testing.T) {
// 	fakeVolumes := fakeSourceVolumes()
// 	expectVolume := &coreV1.Volume{}
// 	k8Yaml.Unmarshal([]byte(sourceVolumeYaml), expectVolume)
//
// 	tests := []struct {
// 		name string
// 		args []coreV1.Volume
// 		want []coreV1.Volume
// 	}{
// 		{
// 			name: "test",
// 			args: fakeVolumes,
// 			want: []coreV1.Volume{*expectVolume},
// 		},
// 	}
//
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			if got := shadowVolumes(tt.args); !reflect.DeepEqual(got, tt.want) {
// 				t.Errorf("shadowVolumes() = %v, want %v", got, tt.want)
// 			}
// 		})
// 	}
// }
//
// func Test_sourceName(t *testing.T) {
//
// 	testName1 := "visit-service"
// 	testName2 := "vet-service"
// 	testName3 := "custom-service"
// 	tests := []struct {
// 		name string
// 		want string
// 	}{
// 		{
// 			name: testName1 + shadowDeploymentNameSuffix,
// 			want: testName1,
// 		},
// 		{
// 			name: testName2 + shadowDeploymentNameSuffix,
// 			want: testName2,
// 		},
// 		{
// 			name: testName3 + shadowDeploymentNameSuffix,
// 			want: testName3,
// 		},
// 	}
//
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			if got := sourceName(tt.name); got != tt.want {
// 				t.Errorf("sourceName() = %v, want %v", got, tt.want)
// 			}
// 		})
// 	}
// }

const sourceVolumeYaml = `
name: configmap-volume-0
configMap:
  items:
  - key: application-sit-yml
    path: application-sit.yml
  name: visits-service
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

const sourceContainerYaml = `
        args:
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
`
