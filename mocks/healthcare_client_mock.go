package mocks

import (
	"net/http"

	"google.golang.org/api/googleapi"
	"google.golang.org/api/healthcare/v1"
)

type MockDatasetCreateCall struct{}
type MockDatastoreGetCall struct{}
type MockDatasetGetCallBadRequest struct{}

type MockFhirCreateCall struct{}
type MockFhirDeleteCall struct{}
type MockFhirExportCall struct{}
type MockFhirGetCall struct{}
type MockFhirGetCallBadRequest struct{}
type MockFhirCreateCallBadRequest struct{}
type MockFhirDeleteCallBadRequest struct{}
type MockFhirCreateOrUpdateIAMPolicyCall struct{}
type MockFhirCreateOrUpdateIAMPolicyCallBadRequest struct{}
type MockFhirGetIAMPolicyCall struct{}
type MockFhirGetIAMPolicyCallBadRequest struct{}
type MockFhirPatchCall struct{}

type MockFhirResourceDeleteCallBadRequest struct{}
type MockFhirResourceDeleteCall struct{}
type MockFhirResourceUpdateCallBadRequest struct{}
type MockFhirResourceUpdateCall struct{}
type MockFhirResourceGetCallBadRequest struct{}
type MockFhirResourceGetCall struct{}
type MockFhirResourceGetCallReturnedResource struct{}

func (m *MockFhirCreateCall) Do(opts ...googleapi.CallOption) (*healthcare.FhirStore, error) {
	return nil, nil
}

func (m *MockFhirExportCall) Do(opts ...googleapi.CallOption) (*healthcare.Operation, error) {
	return nil, nil
}

func (m *MockDatasetCreateCall) Do(opts ...googleapi.CallOption) (*healthcare.Operation, error) {
	return nil, nil
}

func (m *MockDatastoreGetCall) Do(opts ...googleapi.CallOption) (*healthcare.Dataset, error) {
	return nil, nil
}

func (m *MockFhirPatchCall) Do(opts ...googleapi.CallOption) (*healthcare.FhirStore, error) {
	return nil, nil
}

func (m *MockFhirGetIAMPolicyCall) Do(opts ...googleapi.CallOption) (*healthcare.Policy, error) {
	bindings := []*healthcare.Binding{}
	var policy = &healthcare.Policy{}
	policy.Bindings = append(bindings, &healthcare.Binding{})
	return policy, nil
}

func (m *MockFhirGetIAMPolicyCallBadRequest) Do(opts ...googleapi.CallOption) (*healthcare.Policy, error) {
	googleError := &googleapi.Error{Code: 500}
	return nil, googleError
}

func (m *MockFhirCreateOrUpdateIAMPolicyCall) Do(opts ...googleapi.CallOption) (*healthcare.Policy, error) {
	bindings := []*healthcare.Binding{}
	var policy = &healthcare.Policy{}
	policy.Bindings = append(bindings, &healthcare.Binding{})
	return policy, nil
}

func (m *MockFhirCreateOrUpdateIAMPolicyCallBadRequest) Do(opts ...googleapi.CallOption) (*healthcare.Policy, error) {
	googleError := &googleapi.Error{Code: 500}
	return nil, googleError
}

func (m *MockFhirGetCall) Do(opts ...googleapi.CallOption) (*healthcare.FhirStore, error) {
	return nil, nil
}

func (m *MockFhirDeleteCall) Do(opts ...googleapi.CallOption) (*healthcare.Empty, error) {
	return nil, nil
}

func (m *MockFhirResourceDeleteCall) Do(opts ...googleapi.CallOption) (*http.Response, error) {
	return nil, nil
}

func (m *MockFhirResourceUpdateCall) Do(opts ...googleapi.CallOption) (*http.Response, error) {
	respone := &http.Response{
		StatusCode: 200,
		Body:       http.NoBody,
	}
	return respone, nil
}
func (m *MockFhirResourceGetCall) Do(opts ...googleapi.CallOption) (*http.Response, error) {
	respone := &http.Response{
		StatusCode: 400,
		Body:       http.NoBody,
	}
	return respone, nil
}

func (m *MockFhirResourceGetCallReturnedResource) Do(opts ...googleapi.CallOption) (*http.Response, error) {
	respone := &http.Response{
		StatusCode: 200,
		Body:       http.NoBody,
	}
	return respone, nil
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

func (m *MockFhirResourceDeleteCallBadRequest) Do(opts ...googleapi.CallOption) (*http.Response, error) {
	googleError := &googleapi.Error{}
	return nil, googleError
}

func (m *MockFhirResourceGetCallBadRequest) Do(opts ...googleapi.CallOption) (*http.Response, error) {
	googleError := &googleapi.Error{}
	return nil, googleError
}

func (m *MockFhirResourceUpdateCallBadRequest) Do(opts ...googleapi.CallOption) (*http.Response, error) {
	googleError := &googleapi.Error{}
	return nil, googleError
}
