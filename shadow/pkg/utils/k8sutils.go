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

	appsV1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
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

	"github.com/megaease/easemesh/mesh-operator/pkg/api/v1beta1"
)

func NewKubernetesClient() (*kubernetes.Clientset, error) {
	kubeConfig, err := ctrl.GetConfig()
	if err != nil {
		return nil, err
	}

	kubeClient, err := kubernetes.NewForConfig(kubeConfig)
	if err != nil {
		return nil, err
	}
	return kubeClient, nil
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
	kubeConfig, err := ctrl.GetConfig()
	if err != nil {
		return nil, err
	}
	err = v1beta1.AddToScheme(scheme.Scheme)
	if err != nil {
		return nil, err
	}

	crdConfig := *kubeConfig
	crdConfig.ContentConfig.GroupVersion = &schema.GroupVersion{Group: v1beta1.GroupVersion.Group, Version: v1beta1.GroupVersion.Version}
	crdConfig.APIPath = "/apis"
	crdConfig.NegotiatedSerializer = serializer.NewCodecFactory(scheme.Scheme)
	crdConfig.UserAgent = rest.DefaultKubernetesUserAgent()
	crdRestClient, err := rest.UnversionedRESTClientFor(&crdConfig)
	return crdRestClient, err

}

func ListMeshDeployment(client rest.Interface, namespace string, options metav1.ListOptions) (*v1beta1.MeshDeploymentList, error) {
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

func GetMeshDeployment(client rest.Interface, namespace string, name string, options metav1.GetOptions) (*v1beta1.MeshDeployment, error) {
	result := v1beta1.MeshDeployment{}
	err := client.
		Get().
		Namespace(namespace).
		Resource("meshdeployments").
		Name(name).
		VersionedParams(&options, scheme.ParameterCodec).Do(context.TODO()).Into(&result)

	return &result, err
}

func CreateMeshDeployment(client rest.Interface, namespace string, meshDeployment v1beta1.MeshDeployment) (*v1beta1.MeshDeployment, error) {
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

func UpdateMeshDeployment(client rest.Interface, namespace string, meshDeployment v1beta1.MeshDeployment) (*v1beta1.MeshDeployment, error) {
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

func ListDeployments(namespace string, clientSet kubernetes.Interface, options metav1.ListOptions) ([]appsV1.Deployment, error) {
	deploymentList, err := clientSet.AppsV1().Deployments(namespace).List(context.TODO(), options)
	if err != nil {
		return nil, err
	}
	return deploymentList.Items, nil
}

// DeleteDeployment delete Deployment.
func DeleteDeployment(namespace string, name string, clientSet kubernetes.Interface, options metav1.DeleteOptions) error {
	return clientSet.AppsV1().Deployments(namespace).Delete(context.TODO(), name, options)
}

// ListNameSpaces lists namespaces.
func ListNameSpaces(clientSet kubernetes.Interface) ([]corev1.Namespace, error) {
	namespaceList, err := clientSet.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	return namespaceList.Items, nil
}

func DeleteMeshDeployment(client rest.Interface, namespace string, meshDeployment v1beta1.MeshDeployment) error {
	return client.
		Delete().
		Namespace(namespace).
		Name(meshDeployment.Name).
		Resource("meshdeployments").
		Do(context.TODO()).Error()
}

func applyResource(createFunc func() error, updateFunc func() error) error {
	err := createFunc()
	if err != nil && errors.IsAlreadyExists(err) {
		err = updateFunc()
	}
	return err
}

func DeployMesheployment(client rest.Interface, namespace string, deployment *v1beta1.MeshDeployment) error {
	return applyResource(
		func() error {
			_, err := CreateMeshDeployment(client, namespace, *deployment)
			return err
		},
		func() error {
			meshDeployment, err := GetMeshDeployment(client, namespace, deployment.Name, metav1.GetOptions{})
			if err != nil {
				return err
			}
			deployment.ResourceVersion = meshDeployment.ResourceVersion
			_, err = UpdateMeshDeployment(client, namespace, *deployment)
			return err
		})
}
