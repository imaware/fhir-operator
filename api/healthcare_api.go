package api

import (
	"bytes"
	"context"
	"fmt"
	"net/http"

	"google.golang.org/api/googleapi"
	healthcare "google.golang.org/api/healthcare/v1"
)

type FHRIStoreClientCreateCall interface {
	Do(opts ...googleapi.CallOption) (*healthcare.FhirStore, error)
}

type DatastoreClientCreateCall interface {
	Do(opts ...googleapi.CallOption) (*healthcare.Operation, error)
}

type DatastoreClientGetCall interface {
	Do(opts ...googleapi.CallOption) (*healthcare.Dataset, error)
}

type FHIRStoreClientGetCall interface {
	Do(opts ...googleapi.CallOption) (*healthcare.FhirStore, error)
}

type FHIRStoreClientDeleteCall interface {
	Do(opts ...googleapi.CallOption) (*healthcare.Empty, error)
}

type FHIRStoreClientIAMPolicyGetCall interface {
	Do(opts ...googleapi.CallOption) (*healthcare.Policy, error)
}

type FHIRStoreClientPatchCall interface {
	Do(opts ...googleapi.CallOption) (*healthcare.FhirStore, error)
}
type FHIRStoreClientIAMPolicyCreateOrUpdateCall interface {
	Do(opts ...googleapi.CallOption) (*healthcare.Policy, error)
}
type FHIRStoreResourceClientUpdateCall interface {
	Do(opts ...googleapi.CallOption) (*http.Response, error)
}

type FHIRStoreResourceClientDeleteCall interface {
	Do(opts ...googleapi.CallOption) (*http.Response, error)
}

type FHIRStoreResourceClientGetCall interface {
	Do(opts ...googleapi.CallOption) (*http.Response, error)
}

func BuildFHIRStoreCreateCall(projectID string, location string, datasetID string, version string, fhirStoreID string) (*healthcare.ProjectsLocationsDatasetsFhirStoresCreateCall, error) {
	healthcareService, err := healthcare.NewService(context.Background())
	if err != nil {
		return nil, fmt.Errorf("get client error: %v", err)
	}
	fhirStore := &healthcare.FhirStore{Version: version, EnableUpdateCreate: true}
	storesService := healthcareService.Projects.Locations.Datasets.FhirStores
	parent := fmt.Sprintf("projects/%s/locations/%s/datasets/%s", projectID, location, datasetID)
	return storesService.Create(parent, fhirStore).FhirStoreId(fhirStoreID), nil
}

func BuildDatasetCreateCall(projectID string, location string, datasetID string) (*healthcare.ProjectsLocationsDatasetsCreateCall, error) {
	healthcareService, err := healthcare.NewService(context.Background())
	if err != nil {
		return nil, fmt.Errorf("get client error: %v", err)
	}
	datasetService := healthcareService.Projects.Locations.Datasets
	name := fmt.Sprintf("projects/%s/locations/%s", projectID, location)
	return datasetService.Create(name, &healthcare.Dataset{}).DatasetId(datasetID), nil
}

func BuildDatasetGetCall(projectID string, location string, datasetID string) (*healthcare.ProjectsLocationsDatasetsGetCall, error) {
	healthcareService, err := healthcare.NewService(context.Background())
	if err != nil {
		return nil, fmt.Errorf("get client error: %v", err)
	}
	datasetService := healthcareService.Projects.Locations.Datasets
	name := fmt.Sprintf("projects/%s/locations/%s/datasets/%s", projectID, location, datasetID)
	return datasetService.Get(name), nil
}

func BuildFHIRStoreGetCall(projectID string, location string, datasetID string, fhirStoreID string) (*healthcare.ProjectsLocationsDatasetsFhirStoresGetCall, error) {
	healthcareService, err := healthcare.NewService(context.Background())
	if err != nil {
		return nil, fmt.Errorf("get client error: %v", err)
	}
	datasetService := healthcareService.Projects.Locations.Datasets.FhirStores
	name := fmt.Sprintf("projects/%s/locations/%s/datasets/%s/fhirStores/%s", projectID, location, datasetID, fhirStoreID)
	return datasetService.Get(name), nil
}

func BuildFHIRStoreResourceGetCall(projectID string, location string, datasetID string, fhirStoreID string, resourceType string, resourceID string) (*healthcare.ProjectsLocationsDatasetsFhirStoresFhirReadCall, error) {
	healthcareService, err := healthcare.NewService(context.Background())
	if err != nil {
		return nil, fmt.Errorf("get client error: %v", err)
	}
	fhirService := healthcareService.Projects.Locations.Datasets.FhirStores.Fhir
	name := fmt.Sprintf("projects/%s/locations/%s/datasets/%s/fhirStores/%s/fhir/%s/%s", projectID, location, datasetID, fhirStoreID, resourceType, resourceID)
	call := fhirService.Read(name)
	call.Header().Set("Content-Type", "application/fhir+json;charset=utf-8")
	return call, nil
}

func BuildFhirStoreGetIAMPolicyRequest(projectID string, location string, datasetID string, fhirStoreID string) (*healthcare.ProjectsLocationsDatasetsFhirStoresGetIamPolicyCall, error) {
	healthcareService, err := healthcare.NewService(context.Background())
	if err != nil {
		return nil, fmt.Errorf("get client error: %v", err)
	}
	fhirService := healthcareService.Projects.Locations.Datasets.FhirStores
	name := fmt.Sprintf("projects/%s/locations/%s/datasets/%s/fhirStores/%s", projectID, location, datasetID, fhirStoreID)
	return fhirService.GetIamPolicy(name), nil
}

func BuildFhirStoreUpdateOrCreateIAMPolicyRequest(projectID string, location string, datasetID string, fhirStoreID string, fhirStorePolicy *healthcare.Policy) (*healthcare.ProjectsLocationsDatasetsFhirStoresSetIamPolicyCall, error) {
	healthcareService, err := healthcare.NewService(context.Background())
	if err != nil {
		return nil, fmt.Errorf("get client error: %v", err)
	}
	fhirService := healthcareService.Projects.Locations.Datasets.FhirStores
	name := fmt.Sprintf("projects/%s/locations/%s/datasets/%s/fhirStores/%s", projectID, location, datasetID, fhirStoreID)
	req := &healthcare.SetIamPolicyRequest{
		Policy: fhirStorePolicy,
	}
	return fhirService.SetIamPolicy(name, req), nil
}

func BuildFHIRStoreDeleteCall(projectID string, location string, datasetID string, fhirStoreID string) (*healthcare.ProjectsLocationsDatasetsFhirStoresDeleteCall, error) {
	healthcareService, err := healthcare.NewService(context.Background())
	if err != nil {
		return nil, fmt.Errorf("get client error: %v", err)
	}
	datasetService := healthcareService.Projects.Locations.Datasets.FhirStores
	name := fmt.Sprintf("projects/%s/locations/%s/datasets/%s/fhirStores/%s", projectID, location, datasetID, fhirStoreID)
	return datasetService.Delete(name), nil
}

func BuildFhirStorePatchCall(projectID string, location string, datasetID string, fhirStoreID string, fhirStore *healthcare.FhirStore) (*healthcare.ProjectsLocationsDatasetsFhirStoresPatchCall, error) {
	healthcareService, err := healthcare.NewService(context.Background())
	if err != nil {
		return nil, fmt.Errorf("get client error: %v", err)
	}
	fhirService := healthcareService.Projects.Locations.Datasets.FhirStores
	name := fmt.Sprintf("projects/%s/locations/%s/datasets/%s/fhirStores/%s", projectID, location, datasetID, fhirStoreID)
	return fhirService.Patch(name, fhirStore), nil
}

func BuildFHIRStoreResourceDeleteCall(projectID string, location string, datasetID string, fhirStoreID string, resourceType string, resourceID string) (*healthcare.ProjectsLocationsDatasetsFhirStoresFhirDeleteCall, error) {
	healthcareService, err := healthcare.NewService(context.Background())
	if err != nil {
		return nil, fmt.Errorf("get client error: %v", err)
	}
	fhirService := healthcareService.Projects.Locations.Datasets.FhirStores.Fhir
	name := fmt.Sprintf("projects/%s/locations/%s/datasets/%s/fhirStores/%s/fhir/%s/%s", projectID, location, datasetID, fhirStoreID, resourceType, resourceID)
	return fhirService.Delete(name), nil
}

func BuildFHIRStoreResourceUpdateCall(projectID string, location string, datasetID string, fhirStoreID string, resourceRepresentation string, resourceType string, resourceID string) (*healthcare.ProjectsLocationsDatasetsFhirStoresFhirUpdateCall, error) {
	healthcareService, err := healthcare.NewService(context.Background())
	if err != nil {
		return nil, fmt.Errorf("get client error: %v", err)
	}
	resourceRepresentationBytes := []byte(resourceRepresentation)
	fhirService := healthcareService.Projects.Locations.Datasets.FhirStores.Fhir
	name := fmt.Sprintf("projects/%s/locations/%s/datasets/%s/fhirStores/%s/fhir/%s/%s", projectID, location, datasetID, fhirStoreID, resourceType, resourceID)
	call := fhirService.Update(name, bytes.NewReader(resourceRepresentationBytes))
	call.Header().Set("Content-Type", "application/fhir+json;charset=utf-8")
	return call, nil
}

func CreateFHIRStore(fhirseStoreCreateCall FHRIStoreClientCreateCall) (*healthcare.FhirStore, error) {
	return fhirseStoreCreateCall.Do()
}

func PatchFHIRStore(fhirStorePatchCall FHIRStoreClientPatchCall) (*healthcare.FhirStore, error) {
	return fhirStorePatchCall.Do()
}
func UpdateFHIRStoreIAMPolicy(fhirStoreCreateOrUpdateIAMPolicyCall FHIRStoreClientIAMPolicyCreateOrUpdateCall) (*healthcare.Policy, error) {
	return fhirStoreCreateOrUpdateIAMPolicyCall.Do()
}

func CreateDataset(datastoreCreateCall DatastoreClientCreateCall) (*healthcare.Operation, error) {
	return datastoreCreateCall.Do()
}

func GetDataset(datastoreGetCall DatastoreClientGetCall) (*healthcare.Dataset, error) {
	return datastoreGetCall.Do()
}

func GetFHIRStore(fhirStoreGetCall FHIRStoreClientGetCall) (*healthcare.FhirStore, error) {
	return fhirStoreGetCall.Do()
}

func GetFHIRStoreIAMPolicy(fhirStoreGetIAMPolicyCall FHIRStoreClientIAMPolicyGetCall) (*healthcare.Policy, error) {
	return fhirStoreGetIAMPolicyCall.Do()
}

func DeleteFHIRStore(fhirStoreDeleteCall FHIRStoreClientDeleteCall) (*healthcare.Empty, error) {
	return fhirStoreDeleteCall.Do()
}

func GetFHIRResource(fhirResourceGetCall FHIRStoreResourceClientGetCall) (*http.Response, error) {
	return fhirResourceGetCall.Do()
}

func DeleteFHIRResource(fhirResourceDeleteCall FHIRStoreResourceClientDeleteCall) (*http.Response, error) {
	return fhirResourceDeleteCall.Do()
}

func UpdateFHIRResource(fhirStoreResourceUpdateCall FHIRStoreResourceClientUpdateCall) (*http.Response, error) {
	return fhirStoreResourceUpdateCall.Do()
}
