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

// FhirStoreSpec defines the desired state of FhirStore
type FhirStoreSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// fhirStoreID is the name that the fhir store will be called
	FhirStoreID string `json:"fhirStoreID"`
	// datasetID is the name of the dataset the fhirstore will be put in
	DatasetID string `json:"datasetID"`
}

// FhirStoreStatus defines the observed state of FhirStore
type FhirStoreStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	Status  string `json:"Status"`
	Message string `json:"Message"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// FhirStore is the Schema for the fhirstores API
type FhirStore struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   FhirStoreSpec   `json:"spec,omitempty"`
	Status FhirStoreStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// FhirStoreList contains a list of FhirStore
type FhirStoreList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []FhirStore `json:"items"`
}

func init() {
	SchemeBuilder.Register(&FhirStore{}, &FhirStoreList{})
}
