package handler

import (
	"reflect"
	"sync"
	"testing"

	"github.com/megaease/easemesh/mesh-shadow/pkg/object"
	appsV1 "k8s.io/api/apps/v1"
)

func Test_isShadowDeployment(t *testing.T) {
	deployment1 := fakeSourceDeployment()
	deployment2 := fakeShadowDeployment()

	type args struct {
		spec appsV1.DeploymentSpec
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "test1",
			args: args{
				spec: deployment1.Spec,
			},
			want: false,
		},
		{
			name: "test2",
			args: args{
				spec: deployment2.Spec,
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isShadowDeployment(tt.args.spec); got != tt.want {
				t.Errorf("isShadowDeployment() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestShadowServiceDeploySearcher_Search(t *testing.T) {
	searchChan := make(chan interface{})
	defer close(searchChan)

	searcher := &ShadowServiceDeploySearcher{
		KubeClient: prepareClientForTest(),
		ResultChan: searchChan,
	}

	sourceDeployment := fakeSourceDeployment()
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			select {
			case obj := <-searcher.ResultChan:
				if !reflect.DeepEqual(obj.(ShadowServiceBlock).deployObj, *sourceDeployment) {
					t.Errorf("Search Deployment Error, Searcher.Search() = %v, \n want %v", obj, sourceDeployment)
				}
				return
			}
		}
	}()

	shadowService := fakeShadowService()
	objs := []object.ShadowService{shadowService}
	searcher.Search(objs)
	wg.Wait()
}

const deploymentYaml01 = `
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

const deploymentYaml02 = `
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
