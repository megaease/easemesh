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
	"path"

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
	"k8s.io/client-go/util/homedir"
	"sigs.k8s.io/controller-runtime/pkg/client"

	_ "k8s.io/client-go/plugin/pkg/client/auth"
	ctrl "sigs.k8s.io/controller-runtime"

	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

const (
	DefaultKubeDir          = ".kube"
	DefaultKubernetesConfig = "config"
)

var (
	DefaultKubernetesConfigDir  = path.Join(homedir.HomeDir(), DefaultKubeDir)
	DefaultKubernetesConfigPath = path.Join(DefaultKubernetesConfigDir, DefaultKubernetesConfig)
)

func NewKubernetesClient() (*kubernetes.Clientset, error) {
	kubeConfig, err := clientcmd.BuildConfigFromFlags("", DefaultKubernetesConfigPath)
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
	kubeConfig, err := clientcmd.BuildConfigFromFlags("", DefaultKubernetesConfigPath)
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

	k8sConfig, err := clientcmd.BuildConfigFromFlags("", DefaultKubernetesConfigPath)
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
