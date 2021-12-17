package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
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
		return nil, fmt.Errorf("Get client error: %v", err)
	}
	storesService := healthcareService.Projects.Locations.Datasets.FhirStores
	parent := fmt.Sprintf("projects/%s/locations/%s/datasets/%s", projectID, location, datasetID)
	return storesService.Create(parent, &healthcare.FhirStore{Version: version, EnableUpdateCreate: true}).FhirStoreId(fhirStoreID), nil
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

func BuildFHIRStoreResourceGetCall(projectID string, location string, datasetID string, fhirStoreID string, resourceType string, resourceID string) (*healthcare.ProjectsLocationsDatasetsFhirStoresFhirReadCall, error) {
	healthcareService, err := healthcare.NewService(context.Background())
	if err != nil {
		return nil, fmt.Errorf("Get client error: %v", err)
	}
	fhirService := healthcareService.Projects.Locations.Datasets.FhirStores.Fhir
	name := fmt.Sprintf("projects/%s/locations/%s/datasets/%s/fhirStores/%s/fhir/%s/%s", projectID, location, datasetID, fhirStoreID, resourceType, resourceID)
	call := fhirService.Read(name)
	call.Header().Set("Content-Type", "application/fhir+json;charset=utf-8")
	return call, nil
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

func BuildFHIRStoreResourceDeleteCall(projectID string, location string, datasetID string, fhirStoreID string, resourceType string, resourceID string) (*healthcare.ProjectsLocationsDatasetsFhirStoresFhirDeleteCall, error) {
	healthcareService, err := healthcare.NewService(context.Background())
	if err != nil {
		return nil, fmt.Errorf("Get client error: %v", err)
	}
	fhirService := healthcareService.Projects.Locations.Datasets.FhirStores.Fhir
	name := fmt.Sprintf("projects/%s/locations/%s/datasets/%s/fhirStores/%s/fhir/%s/%s", projectID, location, datasetID, fhirStoreID, resourceType, resourceID)
	return fhirService.Delete(name), nil
}

func BuildFHIRStoreResourceUpdateCall(projectID string, location string, datasetID string, fhirStoreID string, resourceRepresentation string, resourceType string, resourceID string) (*healthcare.ProjectsLocationsDatasetsFhirStoresFhirUpdateCall, error) {
	healthcareService, err := healthcare.NewService(context.Background())
	if err != nil {
		return nil, fmt.Errorf("Get client error: %v", err)
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

func GetFHIRResource(fhirResourceGetCall FHIRStoreResourceClientGetCall) (*http.Response, error) {
	return fhirResourceGetCall.Do()
}

func DeleteFHIRResource(fhirResourceDeleteCall FHIRStoreResourceClientDeleteCall) (*http.Response, error) {
	return fhirResourceDeleteCall.Do()
}

func UpdateFHIRResource(fhirStoreResourceUpdateCall FHIRStoreResourceClientUpdateCall) (*http.Response, error) {
	return fhirStoreResourceUpdateCall.Do()
}

// updateFHIRResource updates an FHIR resource to be active or not.
func updateFHIRResource(w io.Writer, projectID, location, datasetID, fhirStoreID, resourceType, fhirResourceID string, active bool) error {
	ctx := context.Background()

	healthcareService, err := healthcare.NewService(ctx)
	if err != nil {
		return fmt.Errorf("healthcare.NewService: %v", err)
	}

	fhirService := healthcareService.Projects.Locations.Datasets.FhirStores.Fhir

	// The following payload works with a Patient resource and is not
	// intended to work with other types of FHIR resources. If necessary,
	// supply a new payload with data that corresponds to the FHIR resource
	// you are updating.
	payload := map[string]interface{}{
		"resourceType": resourceType,
		"id":           fhirResourceID,
		"active":       active,
	}
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("json.Encode: %v", err)
	}

	name := fmt.Sprintf("projects/%s/locations/%s/datasets/%s/fhirStores/%s/fhir/%s/%s", projectID, location, datasetID, fhirStoreID, resourceType, fhirResourceID)

	call := fhirService.Update(name, bytes.NewReader(jsonPayload))
	call.Header().Set("Content-Type", "application/fhir+json;charset=utf-8")
	resp, err := call.Do()
	if err != nil {
		return fmt.Errorf("Update: %v", err)
	}
	defer resp.Body.Close()

	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("could not read response: %v", err)
	}

	if resp.StatusCode > 299 {
		return fmt.Errorf("Update: status %d %s: %s", resp.StatusCode, resp.Status, respBytes)
	}
	fmt.Fprintf(w, "%s", respBytes)

	return nil
}

func fhirGetPatientEverything(w io.Writer, projectID, location, datasetID, fhirStoreID, fhirResourceID string) error {
	ctx := context.Background()

	healthcareService, err := healthcare.NewService(ctx)
	if err != nil {
		return fmt.Errorf("healthcare.NewService: %v", err)
	}
	fhirService := healthcareService.Projects.Locations.Datasets.FhirStores.Fhir
	name := fmt.Sprintf("projects/%s/locations/%s/datasets/%s/fhirStores/%s/fhir/Patient/%s", projectID, location, datasetID, fhirStoreID, fhirResourceID)

	resp, err := fhirService.PatientEverything(name).Do()
	if err != nil {
		return fmt.Errorf("PatientEverything: %v", err)
	}

	defer resp.Body.Close()

	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("could not read response: %v", err)
	}

	if resp.StatusCode > 299 {
		return fmt.Errorf("PatientEverything: status %d %s: %s", resp.StatusCode, resp.Status, respBytes)
	}
	fmt.Fprintf(w, "%s", respBytes)

	return nil
}

func ddeleteFHIRResource(w io.Writer, projectID, location, datasetID, fhirStoreID, resourceType, fhirResourceID string) error {
	ctx := context.Background()

	healthcareService, err := healthcare.NewService(ctx)
	if err != nil {
		return fmt.Errorf("healthcare.NewService: %v", err)
	}

	fhirService := healthcareService.Projects.Locations.Datasets.FhirStores.Fhir

	name := fmt.Sprintf("projects/%s/locations/%s/datasets/%s/fhirStores/%s/fhir/%s/%s", projectID, location, datasetID, fhirStoreID, resourceType, fhirResourceID)

	if _, err := fhirService.Delete(name).Do(); err != nil {
		return fmt.Errorf("Delete: %v", err)
	}

	fmt.Fprintf(w, "Deleted %q", name)

	return nil
}
