package api

import (
	"testing"

	"github.com/imaware/fhir-operator/api/v1alpha1"
	"github.com/imaware/fhir-operator/mocks"
)

var (
	mockDatasetBadGetCall  = &mocks.MockDatasetGetCallBadRequest{}
	mockDatasetGoodGetCall = &mocks.MockDatastoreGetCall{}
	mockFhirGoodCreateCall = &mocks.MockFhirCreateCall{}
	mockFhirBadCreateCall  = &mocks.MockFhirCreateCallBadRequest{}
	mockFhirBadGetCall     = &mocks.MockFhirGetCallBadRequest{}
	mockFhirBadDeleteCall  = &mocks.MockFhirDeleteCallBadRequest{}
	mockFhirGoodDeleteCall = &mocks.MockFhirDeleteCall{}
)

func Test_readandor_create_failed_get_dataset_request(t *testing.T) {
	var fhirStore = &v1alpha1.FhirStore{}
	var expected = false
	actual, err := ReadAndOrCreateFHIRStore(mockDatasetBadGetCall, nil, nil, fhirStore)
	if err == nil {
		t.Error("returned no error, wanted an error")
	}
	if actual != expected {
		t.Errorf("expected boolean %v, got boolean %v", expected, actual)
	}
}

func Test_readandor_create_create_fhirstore_fhir_request(t *testing.T) {
	var fhirStore = &v1alpha1.FhirStore{}
	var expected = true
	actual, err := ReadAndOrCreateFHIRStore(mockDatasetGoodGetCall, mockFhirBadGetCall, mockFhirGoodCreateCall, fhirStore)
	if err != nil {
		t.Errorf("returned error %v, but wanted no error", err.Error())
	}
	if actual != expected {
		t.Errorf("expected boolean %v, got boolean %v", expected, actual)
	}
}

func Test_readandor_create_failed_create_fhirstore_fhir_request(t *testing.T) {
	var fhirStore = &v1alpha1.FhirStore{}
	var expected = false
	actual, err := ReadAndOrCreateFHIRStore(mockDatasetGoodGetCall, mockFhirBadGetCall, mockFhirBadCreateCall, fhirStore)
	if err == nil {
		t.Error("returned no error, wanted an error")
	}
	if actual != expected {
		t.Errorf("expected boolean %v, got boolean %v", expected, actual)
	}
}

func Test_readandor_delete_failed_delete_fhirstor_fhir_request(t *testing.T) {
	var fhirStore = &v1alpha1.FhirStore{}
	err := ReadAndOrDeleteFHIRStore(mockDatasetGoodGetCall, mockFhirGoodCreateCall, mockFhirBadDeleteCall, fhirStore)
	if err == nil {
		t.Error("no error returned and wanted an error")
	}
}

func Test_readandor_delete_delete_fhirstor_fhir_request(t *testing.T) {
	var fhirStore = &v1alpha1.FhirStore{}
	err := ReadAndOrDeleteFHIRStore(mockDatasetGoodGetCall, mockFhirGoodCreateCall, mockFhirGoodDeleteCall, fhirStore)
	if err != nil {
		t.Errorf("error returned %v and wanted no error", err.Error())
	}
}
