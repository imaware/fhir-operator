/*
Copyright 2021.

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

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// FhirResourceSpecFhirStoreSelector defines the FhirStore that this resource should apply to
type FhirResourceSpecFhirStoreSelector struct {
	// The FhirStore resource name to select for the resource
	Name string `json:"name"`
}

// FhirResourceSpec defines the desired state of FhirResource
type FhirResourceSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Type this points to the type of the FHIR resource such as ObservationDefinition
	ResourceType string `json:"resourceType"`
	// The representation of the FHIR resource in JSON format
	Representation string `json:"representation"`
	// The FhirStore that the resource will be applied to by the selector
	Selector FhirResourceSpecFhirStoreSelector `json:"selector"`
}

// FhirResourceStatus defines the observed state of FhirResource
type FhirResourceStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	Status  string `json:"Status"`
	Message string `json:"Message"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// FhirResource is the Schema for the fhirresources API
type FhirResource struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   FhirResourceSpec   `json:"spec,omitempty"`
	Status FhirResourceStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// FhirResourceList contains a list of FhirResource
type FhirResourceList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []FhirResource `json:"items"`
}

func init() {
	SchemeBuilder.Register(&FhirResource{}, &FhirResourceList{})
}
