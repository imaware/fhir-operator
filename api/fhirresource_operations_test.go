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

func Test_fhir_resrouce_is_to_be_updated(t *testing.T) {
	var fhirResource = &v1alpha1.FhirResource{}
	var expected = true
	toBeUpdated, err := IsFhirResourceToBeUpdatedOrCreated(fhirResource)
	if err != nil {
		t.Errorf("expected no error received %v", err)
	}
	if expected != toBeUpdated {
		t.Errorf("expected boolean %t, got %t", expected, toBeUpdated)
	}

}

func Test_no_update_same_no_major_diff(t *testing.T) {
	var representation = `{"foo":"bar"}`
	var lastApplied = `{"apiVersion":"fhir.imaware.com/v1alpha1","kind":"FhirResource","metadata":{"annotations":{},"name":"codesystem-0a510f81-c132-4b22-958c-6351fa14068d","namespace":"imaware-dev"},"spec":{"representation":"{\"foo\":\"bar\"}","resourceType":"CodeSystem","selector":{"name":"imaware-dev-store"}}}`
	var fhirResource = generateTestFhirResource("codesystem-0a510f81-c132-4b22-958c-6351fa14068d", "imaware-dev-store", representation, "CodeSystem")
	fhirResource.Annotations = map[string]string{
		"kubectl.kubernetes.io/last-applied-configuration": lastApplied,
	}
	var expected = false
	toBeUpdated, err := IsFhirResourceToBeUpdatedOrCreated(fhirResource)
	if err != nil {
		t.Errorf("expected no error received %v", err)
	}
	if expected != toBeUpdated {
		t.Errorf("expected boolean %t, got %t", expected, toBeUpdated)
	}

}

func Test_update_major_diff(t *testing.T) {
	var representation = `{"bar":"bar"}`
	var lastApplied = `{"apiVersion":"fhir.imaware.com/v1alpha1","kind":"FhirResource","metadata":{"annotations":{},"name":"codesystem-0a510f81-c132-4b22-958c-6351fa14068d","namespace":"imaware-dev"},"spec":{"representation":"{\"foo\":\"bar\"}","resourceType":"CodeSystem","selector":{"name":"imaware-dev-store"}}}`
	var fhirResource = generateTestFhirResource("codesystem-0a510f81-c132-4b22-958c-6351fa14068d", "imaware-dev-store", representation, "CodeSystem")
	fhirResource.Annotations = map[string]string{
		"kubectl.kubernetes.io/last-applied-configuration": lastApplied,
	}
	var expected = true
	toBeUpdated, err := IsFhirResourceToBeUpdatedOrCreated(fhirResource)
	if err != nil {
		t.Errorf("expected no error received %v", err)
	}
	if expected != toBeUpdated {
		t.Errorf("expected boolean %t, got %t", expected, toBeUpdated)
	}

}
func Test_fhir_resrouce_is_to_not_be_updated_created_status(t *testing.T) {
	var fhirResource = &v1alpha1.FhirResource{}
	fhirResource.Status.Status = CREATED
	var expected = false
	toBeUpdated, err := IsFhirResourceToBeUpdatedOrCreated(fhirResource)
	if err != nil {
		t.Errorf("expected no error received %v", err)
	}
	if expected != toBeUpdated {
		t.Errorf("expected boolean %t, got %t", expected, toBeUpdated)
	}

}

func Test_fhir_resource_bad_json(t *testing.T) {
	var representation = `{"bar":"bar"}`
	var lastApplied = "bad_json"
	var fhirResource = generateTestFhirResource("codesystem-0a510f81-c132-4b22-958c-6351fa14068d", "imaware-dev-store", representation, "CodeSystem")
	fhirResource.Annotations = map[string]string{
		"kubectl.kubernetes.io/last-applied-configuration": lastApplied,
	}
	var expected = true
	toBeUpdated, err := IsFhirResourceToBeUpdatedOrCreated(fhirResource)
	if err == nil {
		t.Errorf("expected error received no error")
	}
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
