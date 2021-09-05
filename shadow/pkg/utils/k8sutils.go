/*
 * Copyright (c) 2017, MegaEase
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

package utils

import (
	"context"

	"github.com/megaease/easemesh/mesh-shadow/pkg/config"
	"github.com/megaease/easemesh/mesh-shadow/pkg/object/v1beta1"
	appsV1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	apiextensions "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/client-go/kubernetes"
	"sigs.k8s.io/controller-runtime/pkg/client"

	_ "k8s.io/client-go/plugin/pkg/client/auth"
	ctrl "sigs.k8s.io/controller-runtime"

	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

func NewKubernetesClient() (*kubernetes.Clientset, error) {
	kubeConfig, err := clientcmd.BuildConfigFromFlags("", config.DefaultKubernetesConfigPath)
	if err != nil {
		return nil, err
	}

	kubeClient, err := kubernetes.NewForConfig(kubeConfig)
	if err != nil {
		return nil, err
	}
	return kubeClient, nil
}

func NewKubernetesAPIExtensionsClient() (*apiextensions.Clientset, error) {
	kubeConfig, err := clientcmd.BuildConfigFromFlags("", config.DefaultKubernetesConfigPath)
	if err != nil {
		return nil, err
	}

	clientset, err := apiextensions.NewForConfig(kubeConfig)
	if err != nil {
		return nil, err
	}
	return clientset, nil
}

func NewRuntimeClient() (client.Client, error) {
	mgr, err := ctrl.NewManager(ctrl.GetConfigOrDie(), ctrl.Options{
		Scheme: runtime.NewScheme(),
	})
	if err != nil {
		return nil, err
	}

	runTimeClient := mgr.GetClient()
	return runTimeClient, err
}

func NewCRDRestClient() (*rest.RESTClient, error) {

	k8sConfig, err := clientcmd.BuildConfigFromFlags("", config.DefaultKubernetesConfigPath)
	if err != nil {
		return nil, err
	}
	err = v1beta1.AddToScheme(scheme.Scheme)
	if err != nil {
		return nil, err
	}

	crdConfig := *k8sConfig
	crdConfig.ContentConfig.GroupVersion = &schema.GroupVersion{Group: v1beta1.GroupVersion.Group, Version: v1beta1.GroupVersion.Version}
	crdConfig.APIPath = "/apis"
	crdConfig.NegotiatedSerializer = serializer.NewCodecFactory(scheme.Scheme)
	crdConfig.UserAgent = rest.DefaultKubernetesUserAgent()
	crdRestClient, err := rest.UnversionedRESTClientFor(&crdConfig)
	return crdRestClient, err

}

func ListMeshDeployment(namespace string, client *rest.RESTClient, options metav1.ListOptions) (*v1beta1.MeshDeploymentList, error) {

	result := v1beta1.MeshDeploymentList{}
	err := client.
		Get().
		Namespace(namespace).
		Resource("meshdeployments").
		VersionedParams(&options, scheme.ParameterCodec).
		Do(context.TODO()).
		Into(&result)

	return &result, err
}

func GetMeshDeployment(namespace string, name string, client *rest.RESTClient, options metav1.GetOptions) (*v1beta1.MeshDeployment, error) {
	result := v1beta1.MeshDeployment{}
	err := client.
		Get().
		Namespace(namespace).
		Resource("meshdeployments").
		Name(name).
		VersionedParams(&options, scheme.ParameterCodec).Do(context.TODO()).Into(&result)

	return &result, err
}

func CreateMeshDeployment(namespace string, meshDeployment v1beta1.MeshDeployment, client *rest.RESTClient) (*v1beta1.MeshDeployment, error) {
	result := v1beta1.MeshDeployment{}
	err := client.
		Post().
		Namespace(namespace).
		Name(meshDeployment.Name).
		Resource("meshdeployments").
		Body(&meshDeployment).
		Do(context.TODO()).
		Into(&result)

	return &result, err
}

func UpdateMeshDeployment(namespace string, meshDeployment v1beta1.MeshDeployment, client *rest.RESTClient) (*v1beta1.MeshDeployment, error) {
	result := v1beta1.MeshDeployment{}
	err := client.
		Put().
		Namespace(namespace).
		Name(meshDeployment.Name).
		Resource("meshdeployments").
		Body(&meshDeployment).
		Do(context.TODO()).
		Into(&result)

	return &result, err
}

func GetNamespace(name string, clientSet *kubernetes.Clientset) (*v1.Namespace, error) {
	namespace, err := clientSet.CoreV1().Namespaces().Get(context.TODO(), name, metav1.GetOptions{})
	return namespace, err
}

func CreateNamespace(namespace *v1.Namespace, clientSet *kubernetes.Clientset) error {
	_, err := clientSet.CoreV1().Namespaces().Get(context.TODO(), namespace.Name, metav1.GetOptions{})
	if err != nil && errors.IsNotFound(err) {
		_, err := clientSet.CoreV1().Namespaces().Create(context.TODO(), namespace, metav1.CreateOptions{})
		if err != nil && errors.IsAlreadyExists(err) {
			return nil
		}
		return err
	}
	return nil
}

func ListDeployments(namespace string, clientSet *kubernetes.Clientset, options metav1.ListOptions) ([]appsV1.Deployment, error) {
	deploymentList, err := clientSet.AppsV1().Deployments(namespace).List(context.TODO(), options)
	if err != nil {
		return nil, err
	}
	return deploymentList.Items, nil
}

func GetDeployments(namespace string, name string, clientSet *kubernetes.Clientset, options metav1.GetOptions) (*appsV1.Deployment, error) {
	deployment, err := clientSet.AppsV1().Deployments(namespace).Get(context.TODO(), name, options)
	if err != nil {
		return nil, err
	}
	return deployment, nil
}

func applyResource(createFunc func() error, updateFunc func() error) error {
	err := createFunc()
	if err != nil && errors.IsAlreadyExists(err) {
		err = updateFunc()
	}
	return err
}

func DeployDeployment(deployment *appsV1.Deployment, clientSet *kubernetes.Clientset, namespace string) error {
	return applyResource(
		func() error {
			_, err := clientSet.AppsV1().Deployments(namespace).Create(context.TODO(), deployment, metav1.CreateOptions{})
			return err
		},
		func() error {
			_, err := clientSet.AppsV1().Deployments(namespace).Update(context.TODO(), deployment, metav1.UpdateOptions{})
			return err
		})
}

func DeployMesheployment(namespace string, deployment *v1beta1.MeshDeployment, client *rest.RESTClient) error {
	return applyResource(
		func() error {
			_, err := CreateMeshDeployment(namespace, *deployment, client)
			return err
		},
		func() error {
			meshDeployment, err := GetMeshDeployment(namespace, deployment.Name, client, metav1.GetOptions{})
			if err != nil {
				return err
			}
			deployment.ResourceVersion = meshDeployment.ResourceVersion
			_, err = UpdateMeshDeployment(namespace, *deployment, client)
			return err
		})
}

// func DeployStatefulset(statefulSet *appsV1.StatefulSet, clientSet *kubernetes.Clientset, namespace string) error {
// 	return applyResource(
// 		func() error {
// 			_, err := clientSet.AppsV1().StatefulSets(namespace).Create(context.TODO(), statefulSet, metav1.CreateOptions{})
// 			return err
// 		},
// 		func() error {
// 			_, err := clientSet.AppsV1().StatefulSets(namespace).Update(context.TODO(), statefulSet, metav1.UpdateOptions{})
// 			return err
// 		},
// 	)
// }
//
// func DeployService(service *v1.Service, clientSet *kubernetes.Clientset, namespace string) error {
// 	return applyResource(
// 		func() error {
// 			_, err := clientSet.CoreV1().Services(namespace).Create(context.TODO(), service, metav1.CreateOptions{})
// 			return err
// 		},
// 		func() error {
// 			_, err := clientSet.CoreV1().Services(namespace).Update(context.TODO(), service, metav1.UpdateOptions{})
// 			return err
// 		},
// 	)
// }
//
// func DeployConfigMap(configMap *v1.ConfigMap, clientSet *kubernetes.Clientset, namespace string) error {
// 	return applyResource(
// 		func() error {
// 			_, err := clientSet.CoreV1().ConfigMaps(namespace).Create(context.TODO(), configMap, metav1.CreateOptions{})
// 			return err
// 		},
// 		func() error {
// 			_, err := clientSet.CoreV1().ConfigMaps(namespace).Update(context.TODO(), configMap, metav1.UpdateOptions{})
// 			return err
// 		})
// }
//
// func DeploySecret(secret *v1.Secret, clientSet *kubernetes.Clientset, namespace string) error {
// 	return applyResource(
// 		func() error {
// 			_, err := clientSet.CoreV1().Secrets(namespace).Create(context.TODO(), secret, metav1.CreateOptions{})
// 			return err
// 		},
// 		func() error {
// 			_, err := clientSet.CoreV1().Secrets(namespace).Update(context.TODO(), secret, metav1.UpdateOptions{})
// 			return err
// 		})
// }
//
// func DeployMutatingWebhookConfig(mutatingWebhookConfig *admissionregv1.MutatingWebhookConfiguration, clientSet *kubernetes.Clientset, namespace string) error {
// 	return applyResource(
// 		func() error {
// 			_, err := clientSet.AdmissionregistrationV1().MutatingWebhookConfigurations().Create(context.TODO(), mutatingWebhookConfig, metav1.CreateOptions{})
// 			return err
// 		},
// 		func() error {
// 			_, err := clientSet.AdmissionregistrationV1().MutatingWebhookConfigurations().Update(context.TODO(), mutatingWebhookConfig, metav1.UpdateOptions{})
// 			return err
// 		})
// }
//
// func ListPersistentVolume(clientSet *kubernetes.Clientset) (*v1.PersistentVolumeList, error) {
// 	return clientSet.CoreV1().PersistentVolumes().List(context.TODO(), metav1.ListOptions{})
// }
//
// func DeployRole(role *rbacv1.Role, clientSet *kubernetes.Clientset, namespace string) error {
// 	return applyResource(
// 		func() error {
// 			_, err := clientSet.RbacV1().Roles(namespace).Create(context.TODO(), role, metav1.CreateOptions{})
// 			return err
// 		},
// 		func() error {
// 			_, err := clientSet.RbacV1().Roles(namespace).Update(context.TODO(), role, metav1.UpdateOptions{})
// 			return err
// 		})
// }
//
// func DeployRoleBinding(roleBinding *rbacv1.RoleBinding, clientSet *kubernetes.Clientset, namespace string) error {
// 	_, err := clientSet.RbacV1().RoleBindings(namespace).Get(context.TODO(), roleBinding.Name, metav1.GetOptions{})
// 	if err != nil && errors.IsNotFound(err) {
// 		_, err = clientSet.RbacV1().RoleBindings(namespace).Create(context.TODO(), roleBinding, metav1.CreateOptions{})
// 	} else {
// 		_, err = clientSet.RbacV1().RoleBindings(namespace).Update(context.TODO(), roleBinding, metav1.UpdateOptions{})
// 	}
// 	return err
// }
//
// func DeployClusterRole(clusterRole *rbacv1.ClusterRole, clientSet *kubernetes.Clientset) error {
// 	return applyResource(
// 		func() error {
// 			_, err := clientSet.RbacV1().ClusterRoles().Create(context.TODO(), clusterRole, metav1.CreateOptions{})
// 			return err
// 		},
// 		func() error {
// 			_, err := clientSet.RbacV1().ClusterRoles().Update(context.TODO(), clusterRole, metav1.UpdateOptions{})
// 			return err
// 		})
// }
//
// func DeployClusterRoleBinding(clusterRoleBinding *rbacv1.ClusterRoleBinding, clientSet *kubernetes.Clientset) error {
// 	return applyResource(
// 		func() error {
// 			_, err := clientSet.RbacV1().ClusterRoleBindings().Create(context.TODO(), clusterRoleBinding, metav1.CreateOptions{})
// 			return err
// 		},
// 		func() error {
// 			_, err := clientSet.RbacV1().ClusterRoleBindings().Update(context.TODO(), clusterRoleBinding, metav1.UpdateOptions{})
// 			return err
// 		})
// }
//
// func DeployCustomResourceDefinition(crd *apiextensionsv1.CustomResourceDefinition, clientSet *apiextensions.Clientset) error {
// 	return applyResource(
// 		func() error {
// 			_, err := clientSet.ApiextensionsV1().CustomResourceDefinitions().Create(context.TODO(), crd, metav1.CreateOptions{})
// 			return err
// 		},
// 		func() error {
// 			_, err := clientSet.ApiextensionsV1().CustomResourceDefinitions().Update(context.TODO(), crd, metav1.UpdateOptions{})
// 			return err
// 		})
// }
//
// type PredictFunc func(interface{}) bool
//
// func StatefulsetReadyPredict(object interface{}) (ready bool) {
// 	statefulset, ok := object.(*appsV1.StatefulSet)
// 	if !ok {
// 		return
// 	}
// 	return statefulset.Status.ReadyReplicas == *statefulset.Spec.Replicas
// }
// func DeploymentReadyPredict(object interface{}) (ready bool) {
// 	deploy, ok := object.(*appsV1.Deployment)
// 	if !ok {
// 		return
// 	}
// 	return deploy.Status.ReadyReplicas == *deploy.Spec.Replicas
// }
// func CheckStatefulsetResourceStatus(client *kubernetes.Clientset, namespace, resourceName string, predict PredictFunc) (bool, error) {
// 	statefulset, err := client.AppsV1().StatefulSets(namespace).Get(context.TODO(), resourceName, metav1.GetOptions{})
// 	if err != nil {
// 		if errors.IsNotFound(err) {
// 			return false, nil
// 		}
// 		return false, err
// 	}
// 	return predict(statefulset), nil
// }
//
// func CheckDeploymentResourceStatus(client *kubernetes.Clientset, namespace, name string, predict PredictFunc) (bool, error) {
//
// 	deploy, err := client.AppsV1().Deployments(namespace).Get(context.TODO(), name, metav1.GetOptions{})
// 	if err != nil {
// 		if errors.IsNotFound(err) {
// 			return false, nil
// 		}
// 		return false, err
// 	}
// 	return predict(deploy), nil
// }
//
// func GetMeshControlPanelEntryPoints(client *kubernetes.Clientset, namespace, resourceName, portName string) ([]string, error) {
// 	service, err := client.CoreV1().Services(namespace).Get(context.TODO(), resourceName, metav1.GetOptions{})
// 	if err != nil {
// 		return nil, err
// 	}
//
// 	nodes, err := client.CoreV1().Nodes().List(context.TODO(), metav1.ListOptions{})
// 	if err != nil {
// 		return nil, err
// 	}
//
// 	var nodePort int32
// 	for _, p := range service.Spec.Ports {
// 		if p.Name == portName {
// 			nodePort = p.NodePort
// 			break
// 		}
// 	}
// 	entrypoints := []string{}
// 	for _, n := range nodes.Items {
// 		for _, i := range n.Status.Addresses {
// 			if i.Type == v1.NodeInternalIP {
// 				entrypoints = append(entrypoints, "http://"+i.Address+":"+strconv.Itoa(int(nodePort)))
// 			}
// 		}
// 	}
// 	return entrypoints, nil
// }
//
// func BatchDeployResources(ctx *StageContext, installFuncs []InstallFunc) error {
// 	for _, fn := range installFuncs {
// 		err := fn.Deploy(ctx)
// 		if err != nil {
// 			return err
// 		}
// 	}
// 	return nil
// }
//
// func DeleteStatefulsetResource(client *kubernetes.Clientset, resource, namespace, name string) error {
// 	err := client.AppsV1().StatefulSets(namespace).Delete(context.Background(), name, metav1.DeleteOptions{})
// 	if err != nil && !errors.IsNotFound(err) {
// 		return err
// 	}
// 	return nil
// }
//
// func DeleteAppsV1Resource(client *kubernetes.Clientset, resource, namespace, name string) error {
// 	err := client.AppsV1().RESTClient().Delete().Resource(resource).Namespace(namespace).Name(name).Do(context.Background()).Error()
// 	if err != nil && !errors.IsNotFound(err) {
// 		return err
// 	}
// 	return nil
// }
//
// func DeleteCoreV1Resource(client *kubernetes.Clientset, resource, namespace, name string) error {
// 	err := client.CoreV1().RESTClient().Delete().Resource(resource).Namespace(namespace).Name(name).Do(context.Background()).Error()
// 	if err != nil && !errors.IsNotFound(err) {
// 		return err
// 	}
// 	return nil
// }
//
// func DeleteRbacV1Resources(client *kubernetes.Clientset, resources, namespace, name string) error {
// 	err := client.RbacV1().RESTClient().Delete().Resource(resources).Namespace(namespace).Name(name).Do(context.Background()).Error()
// 	if err != nil && !errors.IsNotFound(err) {
// 		return err
// 	}
// 	return nil
// }
//
// func DeleteAdmissionregV1Resources(client *kubernetes.Clientset, resources, namespace, name string) error {
// 	// NOTE: RESTClinet can't find mutatingwebhookconfigurations resource.
// 	err := client.AdmissionregistrationV1().MutatingWebhookConfigurations().Delete(context.TODO(), name, metav1.DeleteOptions{})
// 	if err != nil && !errors.IsNotFound(err) {
// 		return err
// 	}
// 	return nil
// }
//
// func DeleteCertificateV1Beta1Resources(client *kubernetes.Clientset, resources, namespace, name string) error {
// 	// NOTE: RESTClinet can't find csr resource.
// 	err := client.CertificatesV1beta1().CertificateSigningRequests().Delete(context.Background(), name, metav1.DeleteOptions{})
// 	if err != nil && !errors.IsNotFound(err) {
// 		return err
// 	}
// 	return nil
// }
//
// func DeleteCRDResource(client *apiextensions.Clientset, name string) error {
// 	err := client.ApiextensionsV1().CustomResourceDefinitions().Delete(context.Background(), name, metav1.DeleteOptions{})
// 	if err != nil && !errors.IsNotFound(err) {
// 		return err
// 	}
// 	return nil
// }
//
// type deleteResourceFunc func(*kubernetes.Clientset, string, string, string) error
//
// func DeleteResources(client *kubernetes.Clientset, resourceAndName [][]string, namespace string, deletefunc deleteResourceFunc) {
// 	for _, s := range resourceAndName {
// 		err := deletefunc(client, s[0], namespace, s[1])
// 		if err != nil {
// 			common.OutputErrorf("clear resource %s of %s in %s error: %s\n", s[1], s[0], namespace, err)
// 		}
// 	}
// }
//
// type PodStatus struct {
// 	Name            string
// 	ReadyContainer  int
// 	ExpectContainer int
// 	Status          string
// 	Restarts        int
// }
//
// func listPodByLabels(client *kubernetes.Clientset, labels map[string]string, namespace string) []PodStatus {
// 	i := 0
// 	labelSelector := ""
// 	for k, v := range labels {
//
// 		labelSelector += fmt.Sprintf("%s=%s", k, v)
// 		i++
// 		if i < len(labels) {
// 			labelSelector += ","
// 		}
// 	}
//
// 	podList, err := client.CoreV1().Pods(namespace).List(context.Background(), metav1.ListOptions{LabelSelector: labelSelector})
// 	if err != nil {
// 		common.OutputErrorf("Ignore error %s", err)
// 		return nil
// 	}
//
// 	podStatus := []PodStatus{}
// 	for _, p := range podList.Items {
// 		ready, restart := readyAndRestartCount(p.Status.ContainerStatuses)
// 		podStatus = append(podStatus, PodStatus{
// 			Name:            p.Name,
// 			ReadyContainer:  ready,
// 			ExpectContainer: len(p.Spec.Containers),
// 			Status:          string(p.Status.Phase),
// 			Restarts:        restart,
// 		})
// 	}
//
// 	return podStatus
// }
//
// func readyAndRestartCount(statuses []v1.ContainerStatus) (readyCount, restartCount int) {
// 	for _, s := range statuses {
// 		if s.Ready {
// 			readyCount++
// 		}
// 		readyCount += int(s.RestartCount)
// 	}
// 	return
// }
//
// type ListPodFunc func(*kubernetes.Clientset, string) []PodStatus
//
// func AdaptListPodFunc(labels map[string]string) ListPodFunc {
// 	return func(client *kubernetes.Clientset, namespace string) []PodStatus {
// 		return listPodByLabels(client, labels, namespace)
// 	}
// }
//
// func FormatPodStatus(client *kubernetes.Clientset, namespace string, fn ListPodFunc) (format string) {
// 	pods := fn(client, namespace)
// 	format += fmt.Sprintf("%-45s%-8s%-15s%-10s\n", "Name", "Ready", "Status", "Restarts")
// 	for _, p := range pods {
// 		format += fmt.Sprintf("%-45s%-8s%-15s%-10d\n", p.Name, fmt.Sprintf("%d/%d", p.ReadyContainer, p.ExpectContainer), p.Status, p.Restarts)
// 	}
// 	return
// }
