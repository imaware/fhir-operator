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

// FhirGCSConnectorSpec defines the desired state of FhirGCSConnector
type FhirGCSConnectorSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// The name of the topic to consume from
	Topic string `json:"topic"`
	// Name of subscription to create
	SubscriptionName string `json:"subscriptionName"`
	// A filter to apply to subscription events
	Filter            string                           `json:"filter,omitempty"`
	FhirStoreSelector FhirGCSConnectorSpecFHIRSelector `json:"fhirStoreSelector"`
}

// info for binding a fhir resource to a fhir store
type FhirGCSConnectorSpecFHIRSelector struct {
	// name of the fhir store in cluster
	Name string `json:"name"`
}

// FhirGCSConnectorStatus defines the observed state of FhirGCSConnector
type FhirGCSConnectorStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	Status  string `json:"Status"`
	Message string `json:"Message"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// FhirGCSConnector is the Schema for the fhirgcsconnectors API
type FhirGCSConnector struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   FhirGCSConnectorSpec   `json:"spec,omitempty"`
	Status FhirGCSConnectorStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// FhirGCSConnectorList contains a list of FhirGCSConnector
type FhirGCSConnectorList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []FhirGCSConnector `json:"items"`
}

func init() {
	SchemeBuilder.Register(&FhirGCSConnector{}, &FhirGCSConnectorList{})
}
