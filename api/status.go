package api

import "fmt"

const FAILED = "Failed"
const CREATING = "Creating"
const CREATED = "Created"

func DatasetNotFoundOrPermissionsInvalidStatus(datasetID string) string {
	return fmt.Sprintf("Dataset either does not exist with ID %v or permissions are not valid for service account", datasetID)
}

func FHIRStoreCreateFailedStatus(fhirStoreID string) string {
	return fmt.Sprintf("Failed to create FHIR store %v check fhir-controller pod", fhirStoreID)
}

func FHIRStoreDeleteFailedStatus(fhirStoreID string) string {
	return fmt.Sprintf("Failed to delete FHIR store %v check fhir-controller pod", fhirStoreID)
}

func GetInternalError() string {
	return fmt.Sprint("Error check fhir-controller pod")
}

func FHIRStoreCreatingStatus(datasetId string, fhirStoreID string) string {
	return fmt.Sprintf("Creating FHIR store with ID %v in dataset %v", fhirStoreID, datasetId)
}

func FHIRStoreCreatedStatus(fhirStoreID string) string {
	return fmt.Sprintf("FHIR store %v up and running", fhirStoreID)
}
