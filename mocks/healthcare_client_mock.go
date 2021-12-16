package mocks

import (
	"google.golang.org/api/googleapi"
	"google.golang.org/api/healthcare/v1"
)

type MockFhirCreateCall struct{}
type MockDatasetCreateCall struct{}
type MockFhirGetCall struct{}
type MockDatastoreGetCall struct{}
type MockFhirDeleteCall struct{}
type MockDatasetGetCallBadRequest struct{}
type MockFhirGetCallBadRequest struct{}
type MockFhirCreateCallBadRequest struct{}
type MockFhirDeleteCallBadRequest struct{}

func (m *MockFhirCreateCall) Do(opts ...googleapi.CallOption) (*healthcare.FhirStore, error) {
	return nil, nil
}

func (m *MockDatasetCreateCall) Do(opts ...googleapi.CallOption) (*healthcare.Operation, error) {
	return nil, nil
}

func (m *MockDatastoreGetCall) Do(opts ...googleapi.CallOption) (*healthcare.Dataset, error) {
	return nil, nil
}

func (m *MockFhirGetCall) Do(opts ...googleapi.CallOption) (*healthcare.FhirStore, error) {
	return nil, nil
}

func (m *MockFhirDeleteCall) Do(opts ...googleapi.CallOption) (*healthcare.Empty, error) {
	return nil, nil
}

func (m *MockDatasetGetCallBadRequest) Do(opts ...googleapi.CallOption) (*healthcare.Dataset, error) {
	googleError := &googleapi.Error{Code: 500}
	return nil, googleError
}

func (m *MockFhirGetCallBadRequest) Do(opts ...googleapi.CallOption) (*healthcare.FhirStore, error) {
	googleError := &googleapi.Error{Code: 404}
	return nil, googleError
}

func (m *MockFhirCreateCallBadRequest) Do(opts ...googleapi.CallOption) (*healthcare.FhirStore, error) {
	googleError := &googleapi.Error{}
	return nil, googleError
}

func (m *MockFhirDeleteCallBadRequest) Do(opts ...googleapi.CallOption) (*healthcare.Empty, error) {
	googleError := &googleapi.Error{}
	return nil, googleError
}
