/*
Copyright 2019 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1beta1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// MyKindSpec defines the desired state of MyKind
type MyKindSpec struct {
	// Important: Run "make" to regenerate code after modifying this file

	// DeploymentName is the name of the Deployment resource that the
	// controller should create.
	// This field must be specified.
	// +kubebuilder:validation:MaxLength=64
	DeploymentName string `json:"deploymentName"`

	// Replicas is the number of replicas that should be specified on the
	// Deployment resource that the controller creates.
	// If not specified, one replica will be created.
	// +optional
	// +kubebuilder:validation:Minimum=0
	Replicas *int32 `json:"replicas,omitempty"`
}

// MyKindStatus defines the observed state of MyKind
type MyKindStatus struct {
	// Important: Run "make" to regenerate code after modifying this file

	// ReadyReplicas is the number of 'ready' replicas observed on the
	// Deployment resource created for this MyKind resource.
	// +optional
	// +kubebuilder:validation:Minimum=0
	ReadyReplicas int32 `json:"readyReplicas,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// MyKind is the Schema for the mykinds API
type MyKind struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   MyKindSpec   `json:"spec,omitempty"`
	Status MyKindStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// MyKindList contains a list of MyKind
type MyKindList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []MyKind `json:"items"`
}

func init() {
	SchemeBuilder.Register(&MyKind{}, &MyKindList{})
}
