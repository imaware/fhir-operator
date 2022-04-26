//go:build integration

package controllers

import (
	"fmt"
	"testing"

	"github.com/imaware/fhir-operator/api"
	"github.com/imaware/fhir-operator/api/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const DATASET_ID_IT = "demo-dataset"
const FHIRSTORE_ID = "it-fhir-store"
const IT_PROJECT = "imaware-test"
const LOCATION = "us-central1"

func Test_create_fhirstore_basic(t *testing.T) {
	configFhirStore = &api.ConfigVars{
		GCPProject:   "imaware-test",
		GCPLocation:  "us-central1",
		DebugEnabled: false,
	}

	fhirStore := &v1alpha1.FhirStore{
		ObjectMeta: metav1.ObjectMeta{
			Name:      FHIRSTORE_ID,
			Namespace: "default",
		},
		Spec: v1alpha1.FhirStoreSpec{
			FhirStoreID: FHIRSTORE_ID,
			DatasetID:   DATASET_ID_IT,
		},
	}
	tearDownFhirStore(fhirStore)
	datasetGetCall, err := api.BuildDatasetGetCall(IT_PROJECT, LOCATION, DATASET_ID_IT)
	if err != nil {
		t.Errorf("Failed to build dataset get call %v", err.Error())
	}
	fhirstoreGetCall, err := api.BuildFHIRStoreGetCall(IT_PROJECT, LOCATION, DATASET_ID_IT, FHIRSTORE_ID)
	if err != nil {
		t.Errorf("Failed to build fhir get call %v", err.Error())
	}
	fhirstoreCreateCall, err := api.BuildFHIRStoreCreateCall(IT_PROJECT, LOCATION, DATASET_ID_IT, "R4", FHIRSTORE_ID)
	if err != nil {
		t.Errorf("Failed to build create  call %v", err.Error())
	}

	_, err = createFhirStoreLoop(fhirStore, datasetGetCall, fhirstoreGetCall, fhirstoreCreateCall)
	if err != nil {
		t.Errorf("Failed to create  fhirstore %v", err.Error())
	}
}

func Test_create_fhirstore_bigquerry(t *testing.T) {
	configFhirStore = &api.ConfigVars{
		GCPProject:   "imaware-test",
		GCPLocation:  "us-central1",
		DebugEnabled: false,
	}

	fhirStore := &v1alpha1.FhirStore{
		ObjectMeta: metav1.ObjectMeta{
			Name:      FHIRSTORE_ID,
			Namespace: "default",
		},
		Spec: v1alpha1.FhirStoreSpec{
			FhirStoreID: FHIRSTORE_ID,
			DatasetID:   DATASET_ID_IT,
			Options: v1alpha1.FhirStoreSpecOptions{
				Bigquery: []v1alpha1.FhirStoreSpecOptionsBigquery{
					{
						Id: "bq://imaware-test.test",
					},
				},
			},
		},
	}
	tearDownFhirStore(fhirStore)
	datasetGetCall, err := api.BuildDatasetGetCall(IT_PROJECT, LOCATION, DATASET_ID_IT)
	if err != nil {
		t.Errorf("Failed to build dataset get call %v", err.Error())
	}
	fhirstoreGetCall, err := api.BuildFHIRStoreGetCall(IT_PROJECT, LOCATION, DATASET_ID_IT, FHIRSTORE_ID)
	if err != nil {
		t.Errorf("Failed to build fhir get call %v", err.Error())
	}
	fhirstoreCreateCall, err := api.BuildFHIRStoreCreateCall(IT_PROJECT, LOCATION, DATASET_ID_IT, "R4", FHIRSTORE_ID)
	if err != nil {
		t.Errorf("Failed to build create  call %v", err.Error())
	}

	_, err = createFhirStoreLoop(fhirStore, datasetGetCall, fhirstoreGetCall, fhirstoreCreateCall)
	if err != nil {
		t.Errorf("Failed to create  fhirstore %v", err.Error())
	}
	if fhirStore.Status.Status == api.FAILED {
		t.Errorf("Failed to create  fhirstore with bigquerry config")
	}
}

func tearDownFhirStore(fhirStore *v1alpha1.FhirStore) {

	datasetGetCall, err := api.BuildDatasetGetCall(IT_PROJECT, LOCATION, DATASET_ID_IT)
	if err != nil {
		fmt.Printf("Failed to build dataset get call %v", err.Error())
	}
	fhirstoreGetCall, err := api.BuildFHIRStoreGetCall(IT_PROJECT, LOCATION, DATASET_ID_IT, FHIRSTORE_ID)
	if err != nil {
		fmt.Printf("Failed to build fhir get call %v", err.Error())
	}

	_, err = deleteFhirStoreLoop(fhirStore, datasetGetCall, fhirstoreGetCall)
	if err != nil {
		fmt.Printf("Failed to create  fhirstore %v", err.Error())
	}
}
