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

package utils

import (
	"context"

	appsV1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	ctrl "sigs.k8s.io/controller-runtime"
)

// NewKubernetesClient creates Kubernetes client set.
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

// ListDeployments list Deployment.
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
