package api

import (
	"testing"

	"github.com/imaware/fhir-operator/api/v1alpha1"
	"github.com/imaware/fhir-operator/mocks"
)

var (
	mockFhirResourceGoodDeleteCall       = &mocks.MockFhirResourceDeleteCall{}
	mockFhirResourceBadDeleteCall        = &mocks.MockFhirResourceDeleteCallBadRequest{}
	mockFhirResourceGoodGetCall          = &mocks.MockFhirResourceGetCall{}
	mockFhirResourceBadGetCall           = &mocks.MockFhirResourceGetCallBadRequest{}
	mockFhirResourceGoodUpdateCall       = &mocks.MockFhirResourceUpdateCall{}
	mockFhirResourceBadUpdateCall        = &mocks.MockFhirResourceUpdateCallBadRequest{}
	mockFhirResourceReturnedResourceCall = &mocks.MockFhirResourceGetCallReturnedResource{}
)

func Test_delete_and_read_fhir_resource(t *testing.T) {
	var fhirResource = &v1alpha1.FhirResource{}
	err := DeleteAndReadFHIRStoreResource(mockFhirResourceGoodDeleteCall, mockFhirResourceGoodGetCall, fhirResource)
	if err != nil {
		t.Errorf("returned an error %v, wanted no error", err.Error())
	}
}

func Test_delete_and_read_failed_fhir_resource(t *testing.T) {
	var fhirResource = &v1alpha1.FhirResource{}
	err := DeleteAndReadFHIRStoreResource(mockFhirResourceGoodDeleteCall, mockFhirResourceReturnedResourceCall, fhirResource)
	if err == nil {
		t.Error("returned no error, wanted no error")
	}
}

func Test_create_or_update_fhir_resource(t *testing.T) {
	var fhirResource = &v1alpha1.FhirResource{}
	var expected = false
	enqueu, err := CreateOrUpdateFHIRResource(mockFhirResourceGoodUpdateCall, fhirResource)
	if err != nil {
		t.Errorf("returned an error %v wanted no error", err)
	}
	if enqueu != expected {
		t.Errorf("expected boolean %t, got %t", expected, enqueu)
	}
}

func Test_create_or_update_failed_fhir_resource(t *testing.T) {
	var fhirResource = &v1alpha1.FhirResource{}
	var expected = false
	enqueu, err := CreateOrUpdateFHIRResource(mockFhirResourceBadUpdateCall, fhirResource)
	if err == nil {
		t.Error("returned no error, wanted an error")
	}
	if enqueu != expected {
		t.Errorf("expected boolean %t, got %t", expected, enqueu)
	}
}
