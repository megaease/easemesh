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

package installbase

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/megaease/easemeshctl/cmd/common"

	admissionregv1 "k8s.io/api/admissionregistration/v1"
	appsV1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	apiextensions "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

var (
	scheme         = runtime.NewScheme()
	codecs         = serializer.NewCodecFactory(scheme)
	parameterCodec = runtime.NewParameterCodec(scheme)

	encoder = unstructured.NewJSONFallbackEncoder(codecs.LegacyCodec(scheme.PrioritizedVersionsAllGroups()...))

	metadataAccessor = meta.NewAccessor()
)

type (
	createResourceFunc func() error
	getResourceRunc    func() (runtime.Object, error)
	updateResourceFunc func() error

	// PredictFunc is the type of function to predict if the resource is ready.
	PredictFunc func(interface{}) bool

	// PodStatus is the status of Pod.
	PodStatus struct {
		Name            string
		ReadyContainer  int
		ExpectContainer int
		Status          string
		Restarts        int
	}

	deleteResourceFunc func(*kubernetes.Clientset, string, string, string) error

	// ListPodFunc is the type of function to list pod.
	ListPodFunc func(*kubernetes.Clientset, string) []PodStatus
)

// NewKubernetesClient creates Kubernetes client set.
func NewKubernetesClient() (*kubernetes.Clientset, error) {
	config, err := clientcmd.BuildConfigFromFlags("", DefaultKubernetesConfigPath)
	if err != nil {
		return nil, err
	}

	kubeClient, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}
	return kubeClient, nil
}

// NewKubernetesAPIExtensionsClient creates Kubernetes API extensions client.
func NewKubernetesAPIExtensionsClient() (*apiextensions.Clientset, error) {
	config, err := clientcmd.BuildConfigFromFlags("", DefaultKubernetesConfigPath)
	if err != nil {
		return nil, err
	}

	clientset, err := apiextensions.NewForConfig(config)
	if err != nil {
		return nil, err
	}
	return clientset, nil
}

func requestContext() context.Context     { return context.TODO() }
func createOptions() metav1.CreateOptions { return metav1.CreateOptions{} }
func getOptions() metav1.GetOptions       { return metav1.GetOptions{} }
func updateOptions() metav1.UpdateOptions { return metav1.UpdateOptions{} }

func adaptReplaceObject(old, new runtime.Object) error {
	oldAnnots, err := metadataAccessor.Annotations(old)
	if err != nil {
		return err
	}

	delete(oldAnnots, v1.LastAppliedConfigAnnotation)
	metadataAccessor.SetAnnotations(old, oldAnnots)

	lastConfig, err := json.Marshal(old)
	if err != nil {
		return err
	}

	newAnnots, err := metadataAccessor.Annotations(new)
	if err != nil {
		return err
	}
	if newAnnots == nil {
		newAnnots = make(map[string]string)
	}
	newAnnots[v1.LastAppliedConfigAnnotation] = string(lastConfig)

	err = metadataAccessor.SetAnnotations(new, newAnnots)
	if err != nil {
		return err
	}

	oldResourceVersion, err := metadataAccessor.ResourceVersion(old)
	if err == nil {
		metadataAccessor.SetResourceVersion(new, oldResourceVersion)
	}

	return nil
}

func deployResource(createFn createResourceFunc, updateFn updateResourceFunc) error {
	err := createFn()
	if err == nil {
		return nil
	}

	if !errors.IsAlreadyExists(err) {
		return err
	}

	return updateFn()
}

// DeployNamespace creates or updates Namespace.
func DeployNamespace(namespace *v1.Namespace, clientSet *kubernetes.Clientset) error {
	createFn := func() error {
		_, err := clientSet.CoreV1().Namespaces().
			Create(requestContext(), namespace, createOptions())
		return err
	}

	updateFn := func() error {
		oldObject, err := clientSet.CoreV1().Namespaces().
			Get(requestContext(), namespace.Name, getOptions())
		if err != nil {
			return err
		}

		err = adaptReplaceObject(oldObject, namespace)
		if err != nil {
			return err
		}

		_, err = clientSet.CoreV1().Namespaces().
			Update(requestContext(), namespace, updateOptions())
		return err
	}

	return deployResource(createFn, updateFn)
}

// DeployDeployment creates or updates Deployment.
func DeployDeployment(deployment *appsV1.Deployment, clientSet *kubernetes.Clientset, namespace string) error {
	createFn := func() error {
		_, err := clientSet.AppsV1().Deployments(namespace).
			Create(requestContext(), deployment, createOptions())
		return err
	}

	updateFn := func() error {
		oldObject, err := clientSet.AppsV1().Deployments(namespace).
			Get(requestContext(), deployment.Name, getOptions())
		if err != nil {
			return err
		}

		err = adaptReplaceObject(oldObject, deployment)
		if err != nil {
			return err
		}

		_, err = clientSet.AppsV1().Deployments(namespace).
			Update(requestContext(), deployment, updateOptions())
		return err
	}

	return deployResource(createFn, updateFn)
}

// DeployStatefulset creates or updates StatefulSet.
func DeployStatefulset(statefulset *appsV1.StatefulSet, clientSet *kubernetes.Clientset, namespace string) error {
	createFn := func() error {
		_, err := clientSet.AppsV1().StatefulSets(namespace).
			Create(requestContext(), statefulset, createOptions())
		return err
	}

	updateFn := func() error {
		oldObject, err := clientSet.AppsV1().StatefulSets(namespace).
			Get(requestContext(), statefulset.Name, getOptions())
		if err != nil {
			return err
		}

		err = adaptReplaceObject(oldObject, statefulset)
		if err != nil {
			return err
		}

		_, err = clientSet.AppsV1().StatefulSets(namespace).
			Update(requestContext(), statefulset, updateOptions())
		return err
	}

	return deployResource(createFn, updateFn)
}

// DeployService creates or updates Service.
func DeployService(service *v1.Service, clientSet *kubernetes.Clientset, namespace string) error {
	createFn := func() error {
		_, err := clientSet.CoreV1().Services(namespace).
			Create(requestContext(), service, createOptions())
		return err
	}

	updateFn := func() error {
		oldObject, err := clientSet.CoreV1().Services(namespace).
			Get(requestContext(), service.Name, getOptions())
		if err != nil {
			return err
		}

		err = adaptReplaceObject(oldObject, service)
		if err != nil {
			return err
		}

		// NOTE: https://github.com/helm/helm/issues/6378#issuecomment-557746499
		service.Spec.ClusterIP = oldObject.Spec.ClusterIP
		service.Spec.ClusterIPs = oldObject.Spec.ClusterIPs

		_, err = clientSet.CoreV1().Services(namespace).
			Update(requestContext(), service, updateOptions())
		return err
	}

	return deployResource(createFn, updateFn)
}

// DeployConfigMap creates or updates ConfigMap.
func DeployConfigMap(configMap *v1.ConfigMap, clientSet *kubernetes.Clientset, namespace string) error {
	createFn := func() error {
		_, err := clientSet.CoreV1().ConfigMaps(namespace).
			Create(requestContext(), configMap, createOptions())
		return err
	}

	updateFn := func() error {
		oldObject, err := clientSet.CoreV1().ConfigMaps(namespace).
			Get(requestContext(), configMap.Name, getOptions())
		if err != nil {
			return err
		}

		err = adaptReplaceObject(oldObject, configMap)
		if err != nil {
			return err
		}

		_, err = clientSet.CoreV1().ConfigMaps(namespace).
			Update(requestContext(), configMap, updateOptions())
		return err
	}

	return deployResource(createFn, updateFn)
}

// DeploySecret creates or updates Secret.
func DeploySecret(secret *v1.Secret, clientSet *kubernetes.Clientset, namespace string) error {

	createFn := func() error {
		_, err := clientSet.CoreV1().Secrets(namespace).
			Create(requestContext(), secret, createOptions())
		return err
	}

	updateFn := func() error {
		oldObject, err := clientSet.CoreV1().Secrets(namespace).
			Get(requestContext(), secret.Name, getOptions())
		if err != nil {
			return err
		}

		err = adaptReplaceObject(oldObject, secret)
		if err != nil {
			return err
		}

		_, err = clientSet.CoreV1().Secrets(namespace).
			Update(requestContext(), secret, updateOptions())

		return err
	}

	return deployResource(createFn, updateFn)
}

// DeployMutatingWebhookConfig creates or updates WebHookConfig.
func DeployMutatingWebhookConfig(mutatingWebhookConfig *admissionregv1.MutatingWebhookConfiguration, clientSet *kubernetes.Clientset, namespace string) error {
	createFn := func() error {
		_, err := clientSet.AdmissionregistrationV1().MutatingWebhookConfigurations().
			Create(requestContext(), mutatingWebhookConfig, createOptions())
		return err
	}

	updateFn := func() error {
		oldObject, err := clientSet.AdmissionregistrationV1().MutatingWebhookConfigurations().
			Get(requestContext(), mutatingWebhookConfig.Name, getOptions())
		if err != nil {
			return err
		}

		err = adaptReplaceObject(oldObject, mutatingWebhookConfig)
		if err != nil {
			return err
		}

		_, err = clientSet.AdmissionregistrationV1().MutatingWebhookConfigurations().
			Update(requestContext(), mutatingWebhookConfig, updateOptions())
		return err
	}

	return deployResource(createFn, updateFn)
}

// ListPersistentVolume lists persistent volumes.
func ListPersistentVolume(clientSet *kubernetes.Clientset) (*v1.PersistentVolumeList, error) {
	return clientSet.CoreV1().PersistentVolumes().List(requestContext(), metav1.ListOptions{})
}

// DeployRole creates or updates Role.
func DeployRole(role *rbacv1.Role, clientSet *kubernetes.Clientset, namespace string) error {
	createFn := func() error {
		_, err := clientSet.RbacV1().Roles(namespace).
			Create(requestContext(), role, createOptions())
		return err
	}

	updateFn := func() error {
		oldObject, err := clientSet.RbacV1().Roles(namespace).
			Get(requestContext(), role.Name, getOptions())
		if err != nil {
			return err
		}

		err = adaptReplaceObject(oldObject, role)
		if err != nil {
			return err
		}

		_, err = clientSet.RbacV1().Roles(namespace).
			Update(requestContext(), role, updateOptions())
		return err
	}

	return deployResource(createFn, updateFn)
}

// DeployRoleBinding creates or updates RoleBinding.
func DeployRoleBinding(roleBinding *rbacv1.RoleBinding, clientSet *kubernetes.Clientset, namespace string) error {
	createFn := func() error {
		_, err := clientSet.RbacV1().RoleBindings(namespace).
			Create(requestContext(), roleBinding, createOptions())
		return err
	}

	updateFn := func() error {
		oldObject, err := clientSet.RbacV1().RoleBindings(namespace).
			Get(requestContext(), roleBinding.Name, getOptions())
		if err != nil {
			return err
		}

		err = adaptReplaceObject(oldObject, roleBinding)
		if err != nil {
			return err
		}

		_, err = clientSet.RbacV1().RoleBindings(namespace).
			Update(requestContext(), roleBinding, updateOptions())
		return err
	}

	return deployResource(createFn, updateFn)
}

// DeployClusterRole creates or updates ClusterRole.
func DeployClusterRole(clusterRole *rbacv1.ClusterRole, clientSet *kubernetes.Clientset) error {
	createFn := func() error {
		_, err := clientSet.RbacV1().ClusterRoles().
			Create(requestContext(), clusterRole, createOptions())
		return err
	}

	updateFn := func() error {
		oldObject, err := clientSet.RbacV1().ClusterRoles().
			Get(requestContext(), clusterRole.Name, getOptions())
		if err != nil {
			return err
		}

		err = adaptReplaceObject(oldObject, clusterRole)
		if err != nil {
			return err
		}

		_, err = clientSet.RbacV1().ClusterRoles().
			Update(requestContext(), clusterRole, updateOptions())
		return err
	}

	return deployResource(createFn, updateFn)
}

// DeployClusterRoleBinding creates or updates ClusterRoleBinding.
func DeployClusterRoleBinding(clusterRoleBinding *rbacv1.ClusterRoleBinding, clientSet *kubernetes.Clientset) error {
	createFn := func() error {
		_, err := clientSet.RbacV1().ClusterRoleBindings().
			Create(requestContext(), clusterRoleBinding, createOptions())
		return err
	}

	updateFn := func() error {
		oldObject, err := clientSet.RbacV1().ClusterRoleBindings().
			Get(requestContext(), clusterRoleBinding.Name, getOptions())
		if err != nil {
			return err
		}

		err = adaptReplaceObject(oldObject, clusterRoleBinding)
		if err != nil {
			return err
		}

		_, err = clientSet.RbacV1().ClusterRoleBindings().
			Update(requestContext(), clusterRoleBinding, updateOptions())
		return err
	}

	return deployResource(createFn, updateFn)
}

// DeployCustomResourceDefinition creates or updates CustomResourceDefinition.
func DeployCustomResourceDefinition(crd *apiextensionsv1.CustomResourceDefinition, clientSet *apiextensions.Clientset) error {
	createFn := func() error {
		_, err := clientSet.ApiextensionsV1().CustomResourceDefinitions().
			Create(requestContext(), crd, createOptions())
		return err
	}

	updateFn := func() error {
		oldObject, err := clientSet.ApiextensionsV1().CustomResourceDefinitions().
			Get(requestContext(), crd.Name, getOptions())
		if err != nil {
			return err
		}

		err = adaptReplaceObject(oldObject, crd)
		if err != nil {
			return err
		}

		_, err = clientSet.ApiextensionsV1().CustomResourceDefinitions().
			Update(requestContext(), crd, updateOptions())
		return err
	}

	return deployResource(createFn, updateFn)
}

// StatefulsetReadyPredict returns if the StatefultSet is ready.
func StatefulsetReadyPredict(object interface{}) (ready bool) {
	statefulset, ok := object.(*appsV1.StatefulSet)
	if !ok {
		return
	}
	return statefulset.Status.ReadyReplicas == *statefulset.Spec.Replicas
}

// DeploymentReadyPredict returns if the Deployment is ready.
func DeploymentReadyPredict(object interface{}) (ready bool) {
	deploy, ok := object.(*appsV1.Deployment)
	if !ok {
		return
	}
	return deploy.Status.ReadyReplicas == *deploy.Spec.Replicas
}

// CheckStatefulsetResourceStatus checks if the StatefulSet is ready.
func CheckStatefulsetResourceStatus(client *kubernetes.Clientset, namespace, resourceName string, predict PredictFunc) (bool, error) {
	statefulset, err := client.AppsV1().StatefulSets(namespace).Get(requestContext(), resourceName, metav1.GetOptions{})
	if err != nil {
		if errors.IsNotFound(err) {
			return false, nil
		}
		return false, err
	}
	return predict(statefulset), nil
}

// CheckDeploymentResourceStatus checks if the Deployment is ready.
func CheckDeploymentResourceStatus(client *kubernetes.Clientset, namespace, name string, predict PredictFunc) (bool, error) {
	deploy, err := client.AppsV1().Deployments(namespace).Get(requestContext(), name, metav1.GetOptions{})
	if err != nil {
		if errors.IsNotFound(err) {
			return false, nil
		}
		return false, err
	}
	return predict(deploy), nil
}

// GetMeshControlPlaneEndpoints gets the endpoints of EaseMesh control plane.
func GetMeshControlPlaneEndpoints(client *kubernetes.Clientset, namespace, resourceName, portName string) ([]string, error) {
	service, err := client.CoreV1().Services(namespace).Get(requestContext(), resourceName, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	nodes, err := client.CoreV1().Nodes().List(requestContext(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	var nodePort int32
	for _, p := range service.Spec.Ports {
		if p.Name == portName {
			nodePort = p.NodePort
			break
		}
	}
	entrypoints := []string{}
	for _, n := range nodes.Items {
		for _, i := range n.Status.Addresses {
			if i.Type == v1.NodeInternalIP {
				entrypoints = append(entrypoints, "http://"+i.Address+":"+strconv.Itoa(int(nodePort)))
			}
		}
	}
	return entrypoints, nil
}

// BatchDeployResources deploy resources in batches.
func BatchDeployResources(ctx *StageContext, installFuncs []InstallFunc) error {
	for _, fn := range installFuncs {
		err := fn.Deploy(ctx)
		if err != nil {
			return err
		}
	}
	return nil
}

// DeleteStatefulsetResource deletes Statefulset.
func DeleteStatefulsetResource(client *kubernetes.Clientset, resource, namespace, name string) error {
	err := client.AppsV1().StatefulSets(namespace).Delete(context.Background(), name, metav1.DeleteOptions{})
	if err != nil && !errors.IsNotFound(err) {
		return err
	}
	return nil
}

// DeleteAppsV1Resource deletes resources within group AppV1.
func DeleteAppsV1Resource(client *kubernetes.Clientset, resource, namespace, name string) error {
	err := client.AppsV1().RESTClient().Delete().Resource(resource).Namespace(namespace).Name(name).Do(context.Background()).Error()
	if err != nil && !errors.IsNotFound(err) {
		return err
	}
	return nil
}

// DeleteCoreV1Resource deletes resources within group CoreV1.
func DeleteCoreV1Resource(client *kubernetes.Clientset, resource, namespace, name string) error {
	err := client.CoreV1().RESTClient().Delete().Resource(resource).Namespace(namespace).Name(name).Do(context.Background()).Error()
	if err != nil && !errors.IsNotFound(err) {
		return err
	}
	return nil
}

// DeleteRbacV1Resources deletes resources within group RbacV1.
func DeleteRbacV1Resources(client *kubernetes.Clientset, resources, namespace, name string) error {
	err := client.RbacV1().RESTClient().Delete().Resource(resources).Namespace(namespace).Name(name).Do(context.Background()).Error()
	if err != nil && !errors.IsNotFound(err) {
		return err
	}
	return nil
}

// DeleteAdmissionregV1Resources deletes resources within group AdmissionregV1.
func DeleteAdmissionregV1Resources(client *kubernetes.Clientset, resources, namespace, name string) error {
	// NOTE: RESTClinet can't find mutatingwebhookconfigurations resource.
	err := client.AdmissionregistrationV1().MutatingWebhookConfigurations().Delete(requestContext(), name, metav1.DeleteOptions{})
	if err != nil && !errors.IsNotFound(err) {
		return err
	}
	return nil
}

// DeleteCertificateV1Beta1Resources deletes resources within group CertificateV1Beta1.
func DeleteCertificateV1Beta1Resources(client *kubernetes.Clientset, resources, namespace, name string) error {
	// NOTE: RESTClinet can't find csr resource.
	err := client.CertificatesV1beta1().CertificateSigningRequests().Delete(context.Background(), name, metav1.DeleteOptions{})
	if err != nil && !errors.IsNotFound(err) {
		return err
	}
	return nil
}

// DeleteCRDResource deletes resources within group CustomResourceDefinitions.
func DeleteCRDResource(client *apiextensions.Clientset, name string) error {
	err := client.ApiextensionsV1().CustomResourceDefinitions().Delete(context.Background(), name, metav1.DeleteOptions{})
	if err != nil && !errors.IsNotFound(err) {
		return err
	}
	return nil
}

// DeleteResources deletes resources.
func DeleteResources(client *kubernetes.Clientset, resourceAndName [][]string, namespace string, deletefunc deleteResourceFunc) {
	for _, s := range resourceAndName {
		err := deletefunc(client, s[0], namespace, s[1])
		if err != nil {
			common.OutputErrorf("clear resource %s of %s in %s error: %s\n", s[1], s[0], namespace, err)
		}
	}
}

func listPodByLabels(client *kubernetes.Clientset, labels map[string]string, namespace string) []PodStatus {
	i := 0
	labelSelector := ""
	for k, v := range labels {

		labelSelector += fmt.Sprintf("%s=%s", k, v)
		i++
		if i < len(labels) {
			labelSelector += ","
		}
	}

	podList, err := client.CoreV1().Pods(namespace).List(context.Background(), metav1.ListOptions{LabelSelector: labelSelector})
	if err != nil {
		common.OutputErrorf("Ignore error %s", err)
		return nil
	}

	podStatus := []PodStatus{}
	for _, p := range podList.Items {
		ready, restart := readyAndRestartCount(p.Status.ContainerStatuses)
		podStatus = append(podStatus, PodStatus{
			Name:            p.Name,
			ReadyContainer:  ready,
			ExpectContainer: len(p.Spec.Containers),
			Status:          string(p.Status.Phase),
			Restarts:        restart,
		})
	}

	return podStatus
}

func readyAndRestartCount(statuses []v1.ContainerStatus) (readyCount, restartCount int) {
	for _, s := range statuses {
		if s.Ready {
			readyCount++
		}
		readyCount += int(s.RestartCount)
	}
	return
}

// AdaptListPodFunc adapts the ListPodFunc with labels.
func AdaptListPodFunc(labels map[string]string) ListPodFunc {
	return func(client *kubernetes.Clientset, namespace string) []PodStatus {
		return listPodByLabels(client, labels, namespace)
	}
}

// FormatPodStatus formats PodStatus.
func FormatPodStatus(client *kubernetes.Clientset, namespace string, fn ListPodFunc) (format string) {
	pods := fn(client, namespace)
	format += fmt.Sprintf("%-45s%-8s%-15s%-10s\n", "Name", "Ready", "Status", "Restarts")
	for _, p := range pods {
		format += fmt.Sprintf("%-45s%-8s%-15s%-10d\n", p.Name, fmt.Sprintf("%d/%d", p.ReadyContainer, p.ExpectContainer), p.Status, p.Restarts)
	}
	return
}
