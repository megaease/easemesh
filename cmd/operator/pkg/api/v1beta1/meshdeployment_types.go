/*
Copyright 2021 MegaEase.cn.
*/

package v1beta1

import (
	v1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// ServiceSpec describes mesh service properties
type ServiceSpec struct {
	//Name is mesh service name of the deployment
	Name             string `json:"name"`
	// +kubebuilder:validation:Optional
	// +kubebuilder:default=""
	AppContainerName string `json:"appContainerName"`
	//Labels is dedicated to labeling instance of deployment for traffic control
	// +kubebuilder:validation:Optional
	Labels map[string]string `json:"labels"`
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
	// Deploy describe a service desired state of the K8s deployment
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
