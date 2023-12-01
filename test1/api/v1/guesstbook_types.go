/*
Copyright 2023.

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

package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// GuesstbookSpec defines the desired state of Guesstbook
type GuesstbookSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Foo is an example field of Guesstbook. Edit guesstbook_types.go to remove/update
	Foo       string `json:"foo,omitempty"`
	FirstName string `json:"firstName,omitempty"` // helin 20231201
	LastName  string `json:"lastName,omitempty"`  // helin 20231201
}

// GuesstbookStatus defines the observed state of Guesstbook
type GuesstbookStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	Status1 string `json:"status1"` // helin 20231201
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// Guesstbook is the Schema for the guesstbooks API
type Guesstbook struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   GuesstbookSpec   `json:"spec,omitempty"`
	Status GuesstbookStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// GuesstbookList contains a list of Guesstbook
type GuesstbookList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Guesstbook `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Guesstbook{}, &GuesstbookList{})
}
