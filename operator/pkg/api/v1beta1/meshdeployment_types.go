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

package v1beta1

import (
	v1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ServiceSpec describes mesh service properties
type ServiceSpec struct {
	// Name is mesh service name of the deployment.
	Name string `json:"name"`
	// AppContainerName is the container name of application.
	//
	// +kubebuilder:validation:Optional
	// +kubebuilder:default=""
	AppContainerName string `json:"appContainerName"`

	// Labels is dedicated to label instance of deployment for traffic control.
	// +kubebuilder:validation:Optional
	Labels map[string]string `json:"labels"`

	// AliveProbeURL is alive probe url.
	// +kubebuilder:validation:Optional
	AliveProbeURL string `json:"aliveProbeURL"`

	// ApplicationPort is the listening port of applicaiton.
	// +kubebuilder:validation:Optional
	ApplicationPort uint16 `json:"applicationPort"`
}

// DeploySpec is the specification of the desired behavior of the Deployment.
type DeploySpec struct {

	// Number of desired pods. This is a pointer to distinguish between explicit
	// zero and not specified. Defaults to 1.
	// +optional
	v1.DeploymentSpec `json:",inline"`
}

// MeshDeploymentSpec defines the desired state of MeshDeployment
type MeshDeploymentSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	Service ServiceSpec `json:"service"`
	// Deploy describes a service desired state of the K8s deployment.
	Deploy DeploySpec `json:"deploy,omitempty"`
}

// MeshDeploymentStatus defines the observed state of MeshDeployment
type MeshDeploymentStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:path=meshdeployments,scope=Namespaced

// MeshDeployment is the Schema for the meshdeployments API
type MeshDeployment struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   MeshDeploymentSpec   `json:"spec,omitempty"`
	Status MeshDeploymentStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// MeshDeploymentList contains a list of MeshDeployment
type MeshDeploymentList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []MeshDeployment `json:"items"`
}

func init() {
	SchemeBuilder.Register(&MeshDeployment{}, &MeshDeploymentList{})
}
