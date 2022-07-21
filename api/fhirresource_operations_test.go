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

func Test_no_update_same_no_major_diff(t *testing.T) {
	var representation = `{"bar":"bar"}`
	var fhirResource = generateTestFhirResource("codesystem-0a510f81-c132-4b22-958c-6351fa14068d", "imaware-dev-store", representation, "CodeSystem")
	fhirResource.ResourceVersion = "foo"
	fhirResource.Status.LastObservedResourceVersion = "foo"
	fhirResource.Status.Status = CREATED
	var expected = false
	toBeUpdated := IsFhirResourceToBeUpdatedOrCreated(fhirResource)
	if expected != toBeUpdated {
		t.Errorf("expected boolean %t, got %t", expected, toBeUpdated)
	}

}

func Test_update_major_diff(t *testing.T) {
	var representation = `{"bar":"bar"}`
	var fhirResource = generateTestFhirResource("codesystem-0a510f81-c132-4b22-958c-6351fa14068d", "imaware-dev-store", representation, "CodeSystem")
	fhirResource.ResourceVersion = "foo"
	fhirResource.Status.LastObservedResourceVersion = "bar"
	fhirResource.Status.Status = CREATED
	var expected = true
	toBeUpdated := IsFhirResourceToBeUpdatedOrCreated(fhirResource)
	if expected != toBeUpdated {
		t.Errorf("expected boolean %t, got %t", expected, toBeUpdated)
	}

}
func Test_fhir_resrouce_is_to_not_be_updated_created_status(t *testing.T) {
	var fhirResource = &v1alpha1.FhirResource{}
	fhirResource.Status.Status = CREATED
	var expected = false
	toBeUpdated := IsFhirResourceToBeUpdatedOrCreated(fhirResource)
	if expected != toBeUpdated {
		t.Errorf("expected boolean %t, got %t", expected, toBeUpdated)
	}

}

func generateTestFhirResource(name string, selector string, jsonRepresentation string, resourceType string) *v1alpha1.FhirResource {
	fhirResource := &v1alpha1.FhirResource{}
	fhirResource.Spec.Selector.Name = selector
	fhirResource.Name = name
	fhirResource.Spec.ResourceType = resourceType
	fhirResource.Spec.Representation = jsonRepresentation
	return fhirResource
}
