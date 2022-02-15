/*
 * Copyright (c) 2021, MegaEase
 * All rights reserved.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
package installbase

import (
	"bytes"
	"fmt"
	"testing"

	admissionregv1 "k8s.io/api/admissionregistration/v1"
	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	apiextensions "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	extensionfake "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset/fake"
	k8serr "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	k8yaml "k8s.io/apimachinery/pkg/util/yaml"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/fake"
	k8stesting "k8s.io/client-go/testing"
)

func TestNewKubernetesClient(t *testing.T) {
	NewKubernetesClient()
	NewKubernetesAPIExtensionsClient()
}

func prepareClientForTest() kubernetes.Interface {
	var result runtime.Object
	client := fake.NewSimpleClientset()
	client.PrependReactor("create", "*", func(action k8stesting.Action) (handled bool, ret runtime.Object, err error) {
		result = action.(k8stesting.CreateAction).GetObject()

		return true, action.(k8stesting.CreateAction).GetObject(), k8serr.NewAlreadyExists(schema.GroupResource{
			Resource: "Namespace",
			Group:    "v1",
		}, "na")
	})

	client.PrependReactor("update", "*", func(action k8stesting.Action) (handled bool, ret runtime.Object, err error) {
		return true, action.(k8stesting.UpdateAction).GetObject(), nil
	})

	client.PrependReactor("get", "*", func(action k8stesting.Action) (handled bool, ret runtime.Object, err error) {
		return true, result, nil
	})

	return client
}

func TestDeployNameSpace(t *testing.T) {
	ns := v1.Namespace{ObjectMeta: metav1.ObjectMeta{
		Name: "na",
	}}

	client := prepareClientForTest()
	err := DeployNamespace(&ns, client)
	if err != nil {
		t.Fatalf("deploy namespace error: %s", err)
	}
}

func TestDeployDeployment(t *testing.T) {
	deploy := appsv1.Deployment{}
	client := prepareClientForTest()
	err := DeployDeployment(&deploy, client, "poc")
	if err != nil {
		t.Fatalf("deploy deploy error: %s", err)
	}
}

func TestDeployStatefulset(t *testing.T) {
	statefulset := appsv1.StatefulSet{}
	client := prepareClientForTest()
	err := DeployStatefulset(&statefulset, client, "poc")
	if err != nil {
		t.Fatalf("deploy statefulset error: %s", err)
	}
}

func TestDeployRole(t *testing.T) {
	role := rbacv1.Role{}
	client := prepareClientForTest()
	err := DeployRole(&role, client, "poc")
	if err != nil {
		t.Fatalf("deploy role error: %s", err)
	}
}

func TestDeploySecret(t *testing.T) {
	secret := v1.Secret{}
	client := prepareClientForTest()
	err := DeploySecret(&secret, client, "poc")
	if err != nil {
		t.Fatalf("deploy secret error: %s", err)
	}
}

func TestDeployConfigMap(t *testing.T) {
	object := v1.ConfigMap{}
	client := prepareClientForTest()
	err := DeployConfigMap(&object, client, "poc")
	if err != nil {
		t.Fatalf("deploy configmap error: %s", err)
	}
}

func TestDeployService(t *testing.T) {
	object := v1.Service{}
	client := prepareClientForTest()
	err := DeployService(&object, client, "poc")
	if err != nil {
		t.Fatalf("deploy service error: %s", err)
	}
}

func TestClusterRole(t *testing.T) {
	object := rbacv1.ClusterRole{}
	client := prepareClientForTest()
	err := DeployClusterRole(&object, client)
	if err != nil {
		t.Fatalf("deploy clusterrole error: %s", err)
	}
}

func TestDeployClusterRolebinding(t *testing.T) {
	object := rbacv1.ClusterRoleBinding{}
	client := prepareClientForTest()
	err := DeployClusterRoleBinding(&object, client)
	if err != nil {
		t.Fatalf("deploy clusterrolebinding error: %s", err)
	}
}

func TestDeployRolebinding(t *testing.T) {
	object := rbacv1.RoleBinding{}
	client := prepareClientForTest()
	err := DeployRoleBinding(&object, client, "poc")
	if err != nil {
		t.Fatalf("deploy rolebinding error: %s", err)
	}
}

func TestDeployMutatingWebhookConfig(t *testing.T) {
	object := admissionregv1.MutatingWebhookConfiguration{}
	client := prepareClientForTest()
	err := DeployMutatingWebhookConfig(&object, client, "poc")
	if err != nil {
		t.Fatalf("deploy mutatingwebhookconfig error: %s", err)
	}
}

func TestListPersistentVolume(t *testing.T) {
	client := prepareClientForTest()
	_, err := ListPersistentVolume(client)
	if err != nil {
		t.Fatalf("list persistent error: %s", err)
	}
}

func prepareAPIExtensionClientForTest() apiextensions.Interface {
	client := extensionfake.NewSimpleClientset()
	var result runtime.Object
	client.PrependReactor("create", "*", func(action k8stesting.Action) (handled bool, ret runtime.Object, err error) {
		result = action.(k8stesting.CreateAction).GetObject()

		return true, action.(k8stesting.CreateAction).GetObject(), k8serr.NewAlreadyExists(schema.GroupResource{
			Resource: "Namespace",
			Group:    "v1",
		}, "na")
	})

	client.PrependReactor("update", "*", func(action k8stesting.Action) (handled bool, ret runtime.Object, err error) {
		return true, action.(k8stesting.UpdateAction).GetObject(), nil
	})

	client.PrependReactor("get", "*", func(action k8stesting.Action) (handled bool, ret runtime.Object, err error) {
		return true, result, nil
	})
	return client
}

func TestDeployCustomResource(t *testing.T) {
	object := apiextensionsv1.CustomResourceDefinition{}
	client := prepareAPIExtensionClientForTest()
	err := DeployCustomResourceDefinition(&object, client)
	if err != nil {
		t.Fatalf("deploy customeresource error: %s", err)
	}
}

func TestStatefulsetPredict(t *testing.T) {
	s := appsv1.StatefulSet{}

	s.Status.ReadyReplicas = 1
	s.Spec.Replicas = &s.Status.ReadyReplicas

	if !StatefulsetReadyPredict(&s) {
		t.Fatalf("expect statefulset is read")
	}

	a := struct{}{}
	if StatefulsetReadyPredict(a) {
		t.Fatalf("expect statefulset is not read")
	}
}

func TestDeployPredict(t *testing.T) {
	s := appsv1.Deployment{}

	s.Status.ReadyReplicas = 1
	s.Spec.Replicas = &s.Status.ReadyReplicas

	if !DeploymentReadyPredict(&s) {
		t.Fatalf("expect statefulset is read")
	}

	a := struct{}{}
	if DeploymentReadyPredict(a) {
		t.Fatalf("expect statefulset is not read")
	}
}

func TestCheckStatefulsetResourcestatus(t *testing.T) {
	client := prepareClientForTest()

	if ok, err := CheckStatefulsetResourceStatus(client,
		"ns", "easemesh-control-plane",
		func(o interface{}) bool {
			return true
		}); !ok || err != nil {
		t.Fatalf("check statefulset resource should ok")
	}
}

func TestCheckDeploymentResourcestatus(t *testing.T) {
	client := prepareClientForTest()

	if ok, err := CheckDeploymentResourceStatus(client,
		"ns", "easemesh-operator",
		func(o interface{}) bool {
			return true
		}); !ok || err != nil {
		t.Fatalf("check deployment resource should ok")
	}
}

func TestGetMeshControlPlaneEndpoints(t *testing.T) {
	client := fake.NewSimpleClientset()

	client.PrependReactor("get", "services",
		func(action k8stesting.Action) (handled bool, ret runtime.Object, err error) {
			service := v1.Service{}
			dec := k8yaml.NewYAMLOrJSONDecoder(bytes.NewReader([]byte(serviceSpec)), 1000)
			err = dec.Decode(&service)
			if err != nil {
				t.Fatal(fmt.Sprintf("Error while decoding YAML object. Err was: %s", err))
			}
			return true, &service, nil
		})
	client.PrependReactor("list", "nodes",
		func(action k8stesting.Action) (handled bool, ret runtime.Object, err error) {
			node := v1.Node{}
			dec := k8yaml.NewYAMLOrJSONDecoder(bytes.NewReader([]byte(nodeSpec)), 1000)
			err = dec.Decode(&node)
			obj := v1.NodeList{}
			obj.Items = append(obj.Items, node)
			return true, &obj, nil
		})

	e, err := GetMeshControlPlaneEndpoints(client, "easemesh", "", "https")
	if err != nil {
		t.Fatalf("getmeshcontrolplaneendpoint should be successful, but error: %s", err)
	}

	if len(e) == 0 {
		t.Fatalf("should get a controlplane endpoints, but 0")
	}
}

func TestBachDeployResource(t *testing.T) {
	installFunc := []InstallFunc{
		func(ctx *StageContext) error { return nil },
	}
	err := BatchDeployResources(nil, installFunc)
	if err != nil {
		t.Fatalf("batch deploy resource should succeed")
	}
}

func TestDeleteStatefulsetResource(t *testing.T) {
	client := prepareClientForTest()
	err := DeleteStatefulsetResource(client, "na", "easemesh", "easemesh-control-plane")
	if err != nil {
		t.Fatalf("delete statefulset resource error: %s", err)
	}
}

func TestDeleteAdmissionregV1Resource(t *testing.T) {
	client := prepareClientForTest()
	err := DeleteAdmissionregV1Resources(client, "na", "easemesh", "easemesh-control-plane")
	if err != nil {
		t.Fatalf("delete adminssionregv1resource resource error: %s", err)
	}
}

func TestDeleteCRDResource(t *testing.T) {
	client := prepareAPIExtensionClientForTest()
	err := DeleteCRDResource(client, "easemesh-operator")
	if err != nil {
		t.Fatalf("delete adminssionregv1resource resource error: %s", err)
	}
}

func TestDeleteCertificateV1Resource(t *testing.T) {
	client := prepareClientForTest()
	err := DeleteCertificateV1Beta1Resources(client, "na", "easemesh", "easemesh-control-plane")
	if err != nil {
		t.Fatalf("delete adminssionregv1resource resource error: %s", err)
	}
}

func TestDeleteResource(t *testing.T) {
	client := prepareClientForTest()
	resourceName := [][]string{{"a", "b"}}
	deletefunc := DeleteStatefulsetResource
	DeleteResources(client, resourceName, "easemesh", deletefunc)
}

func TestAdaptListPodFunc(t *testing.T) {
	client := fake.NewSimpleClientset()
	labels := map[string]string{
		"app": "easestack-ingress-controller",
	}

	client.PrependReactor("list", "pods",
		func(action k8stesting.Action) (handled bool, ret runtime.Object, err error) {
			obj := v1.PodList{}
			pod := v1.Pod{}
			pod.Labels = labels
			dec := k8yaml.NewYAMLOrJSONDecoder(bytes.NewReader([]byte(serviceSpec)), 1000)
			err = dec.Decode(&pod)
			if err != nil {
				t.Fatal(fmt.Sprintf("Error while decoding YAML object. Err was: %s", err))
			}
			obj.Items = append(obj.Items, pod)
			return true, &obj, nil
		})

	statuses := AdaptListPodFunc(labels)(client, "easemesh")
	if len(statuses) == 0 {
		t.Fatalf("AdaptListPodFunc return number pod status, but result is 0")
	}

	format := FormatPodStatus(client, "mesh", AdaptListPodFunc(labels))
	if format == "" {
		t.Fatalf("format should has contents")
	}
}

func TestDeleteAppV1Resource(t *testing.T) {
	// TODO only when upstream fix RESTClient() bug, I can test it
	// https://github.com/kubernetes/client-go/blob/release-1.22/kubernetes/typed/core/v1/fake/fake_core_client.go#L97
}

const (
	serviceSpec = `
spec:
  clusterIP: 10.233.33.216
  ports:
  - name: https
    port: 8443
    protocol: TCP
    targetPort: 8443
  - name: mutate-webhook
    port: 9090
    protocol: TCP
    targetPort: 9090
  selector:
    easemesh-operator: operator-manager
  sessionAffinity: None
  type: ClusterIP
`

	nodeSpec = `
  spec:
    podCIDR: 10.233.66.0/24
    podCIDRs:
    - 10.233.66.0/24
  status:
    addresses:
    - address: 10.0.20.103
      type: InternalIP
    - address: kube-3
      type: Hostname
`

	podSpec = `
spec:
  containers:
  - command:
    - /bin/sh
    - -c
    - |-
      echo "name: $POD_NAME" > /easegress-ingress/config.yaml && echo '
      cluster-request-timeout: 10s
      cluster-role: writer
      api-addr: "0.0.0.0:2381"
      debug: false
      cluster-name: easegress-ingress-controller
                ' >> /easegress-ingress/config.yaml && /opt/easegress/bin/easegress-server -f /easegress-ingress/config.yaml
    env:
    - name: POD_NAME
      valueFrom:
        fieldRef:
          apiVersion: v1
          fieldPath: metadata.name
    - name: POD_NAMESPACE
      valueFrom:
        fieldRef:
          apiVersion: v1
          fieldPath: metadata.namespace
    image: megaease/easegress:easemesh
    imagePullPolicy: Always
    name: easestack-ingress-controller
    resources: {}
    terminationMessagePath: /dev/termination-log
    terminationMessagePolicy: File
    volumeMounts:
    - mountPath: /easegress-ingress
      name: ingress-params-volume
    - mountPath: /var/run/secrets/kubernetes.io/serviceaccount
      name: easestack-ingress-controller-token-6nng4
      readOnly: true
  dnsPolicy: ClusterFirst
  enableServiceLinks: true
  nodeName: kube-1
  preemptionPolicy: PreemptLowerPriority
  priority: 0
  restartPolicy: Always
  schedulerName: default-scheduler
  securityContext: {}
  serviceAccount: easestack-ingress-controller
  serviceAccountName: easestack-ingress-controller
  terminationGracePeriodSeconds: 30
  tolerations:
  - effect: NoExecute
    key: node.kubernetes.io/not-ready
    operator: Exists
    tolerationSeconds: 300
  - effect: NoExecute
    key: node.kubernetes.io/unreachable
    operator: Exists
    tolerationSeconds: 300
  volumes:
  - emptyDir: {}
    name: ingress-params-volume
  - name: easestack-ingress-controller-token-6nng4
    secret:
      defaultMode: 420
      secretName: easestack-ingress-controller-token-6nng4
`
)
