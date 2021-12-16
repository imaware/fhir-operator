package api

import (
	"context"
	"fmt"

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

func BuildFHIRStoreCreateCall(projectID string, location string, datasetID string, version string, fhirStoreID string) (*healthcare.ProjectsLocationsDatasetsFhirStoresCreateCall, error) {
	healthcareService, err := healthcare.NewService(context.Background())
	if err != nil {
		return nil, fmt.Errorf("Get client error: %v", err)
	}
	storesService := healthcareService.Projects.Locations.Datasets.FhirStores
	parent := fmt.Sprintf("projects/%s/locations/%s/datasets/%s", projectID, location, datasetID)
	return storesService.Create(parent, &healthcare.FhirStore{Version: version}).FhirStoreId(fhirStoreID), nil
}

func BuildDatasetCreateCall(projectID string, location string, datasetID string) (*healthcare.ProjectsLocationsDatasetsCreateCall, error) {
	healthcareService, err := healthcare.NewService(context.Background())
	if err != nil {
		return nil, fmt.Errorf("Get client error: %v", err)
	}
	datasetService := healthcareService.Projects.Locations.Datasets
	name := fmt.Sprintf("projects/%s/locations/%s", projectID, location)
	return datasetService.Create(name, &healthcare.Dataset{}).DatasetId(datasetID), nil
}

func BuildDatasetGetCall(projectID string, location string, datasetID string) (*healthcare.ProjectsLocationsDatasetsGetCall, error) {
	healthcareService, err := healthcare.NewService(context.Background())
	if err != nil {
		return nil, fmt.Errorf("Get client error: %v", err)
	}
	datasetService := healthcareService.Projects.Locations.Datasets
	name := fmt.Sprintf("projects/%s/locations/%s/datasets/%s", projectID, location, datasetID)
	return datasetService.Get(name), nil
}

func BuildFHIRStoreGetCall(projectID string, location string, datasetID string, fhirStoreID string) (*healthcare.ProjectsLocationsDatasetsFhirStoresGetCall, error) {
	healthcareService, err := healthcare.NewService(context.Background())
	if err != nil {
		return nil, fmt.Errorf("Get client error: %v", err)
	}
	datasetService := healthcareService.Projects.Locations.Datasets.FhirStores
	name := fmt.Sprintf("projects/%s/locations/%s/datasets/%s/fhirStores/%s", projectID, location, datasetID, fhirStoreID)
	return datasetService.Get(name), nil
}

func BuildFHIRStoreDeleteCall(projectID string, location string, datasetID string, fhirStoreID string) (*healthcare.ProjectsLocationsDatasetsFhirStoresDeleteCall, error) {
	healthcareService, err := healthcare.NewService(context.Background())
	if err != nil {
		return nil, fmt.Errorf("Get client error: %v", err)
	}
	datasetService := healthcareService.Projects.Locations.Datasets.FhirStores
	name := fmt.Sprintf("projects/%s/locations/%s/datasets/%s/fhirStores/%s", projectID, location, datasetID, fhirStoreID)
	return datasetService.Delete(name), nil
}

func CreateFHIRStore(fhirseStoreCreateCall FHRIStoreClientCreateCall) (*healthcare.FhirStore, error) {
	return fhirseStoreCreateCall.Do()
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

func DeleteFHIRStore(fhirStoreDeleteCall FHIRStoreClientDeleteCall) (*healthcare.Empty, error) {
	return fhirStoreDeleteCall.Do()
}
