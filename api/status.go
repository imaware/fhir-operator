package api

import "fmt"

const FAILED = "Failed"
const CREATING = "Creating"
const CREATED = "Created"
const DELETED = "Deleted"
const PENDING = "Pending"

func DatasetNotFoundOrPermissionsInvalidStatus(datasetID string, errorString error) string {
	return fmt.Sprintf("Dataset either does not exist with ID %v or permissions are not valid for service account: %v", datasetID, errorString)
}

func FHIRStoreCreateFailedStatus(errorString string) string {
	return fmt.Sprintf("Failed to create FHIR store  due to [ %v ]", errorString)
}

func FHIRStorePatchFailedStatus(errorString string) string {
	return fmt.Sprintf("Failed to patch FHIR store due to [ %v ]", errorString)
}

func FHIRStoreDeleteFailedStatus(errorString string) string {
	return fmt.Sprintf("Failed to delete FHIR store due to [ %v ]", errorString)
}

func FHIRStoreExportFailedStatus(errorString string) string {
	return fmt.Sprintf("Failed to export FHIR store due to [ %v ]", errorString)
}

func GetInternalError(errorString string) string {
	return fmt.Sprintf("Internal error: [ %v ]", errorString)
}

func FHIRStoreFailedIAMPolicy(errorString string) string {
	return fmt.Sprintf("Failed to create and or update FHIR store iam policy due to [ %v ]", errorString)
}

func FHIRStoreFailedPatch(errorString string) string {
	return fmt.Sprintf("Failed to update FHIR store due to [ %v ]", errorString)
}

func FHIRStoreCreatingStatus(datasetId string, fhirStoreID string) string {
	return fmt.Sprintf("Creating FHIR store with ID %v in dataset %v", fhirStoreID, datasetId)
}

func FHIRStoreCreatedStatus(fhirStoreID string) string {
	return fmt.Sprintf("FHIR store %v up and running", fhirStoreID)
}

func FHIRStoreResourceCreateOrUpdateFailedStatus(resourceName string, errorString string) string {
	return fmt.Sprintf("Failed to create and or update FHIR resource %v due to [ %v ]", resourceName, errorString)
}

func FHIRStoreResourceCreatedorUpdatedStatus() string {
	return "Resource created and or updated"
}

func FHIRStoreResourceDeleteFailedStatus(resourceName string, errorString string) string {
	return fmt.Sprintf("Failed to delete FhirStore resource %v due to [ %v ]", resourceName, errorString)
}

func FHIRStoreResourceDeletetatus() string {
	return "Resource deleted"
}

func FHIRStoreResourcePendingOnFhirStoreStatus(fhirStoreName string) string {
	return fmt.Sprintf("Waiting on FhirStore %v to be in %v status", fhirStoreName, CREATED)
}

func FHIRStoreResourcePendingOnParentResourceStatus() string {
	return "Waiting on parent resource to be created"
}

func FHIRGCSConnectorFailedSubscriptionToTopicStatus(topicName string, errorString string) string {
	return fmt.Sprintf("Failed to subscribe to topic %s due to [ %s ]", topicName, errorString)
}

func FHIRGCSConnectorFailedSubscriptionCreationStatus(subscriptionName string, errorString string) string {
	return fmt.Sprintf("Failed to create subscription %s due to [ %s ]", subscriptionName, errorString)
}

func FHIRGCSConnectorFailedToGetTopicStatus(topicName string, errorString string) string {
	return fmt.Sprintf("Failed to get topic %s due to [ %s ]", topicName, errorString)
}

func FHIRGCSConnectorConsumingFromTopicStatus(subscriptionName string, topicName string) string {
	return fmt.Sprintf("Subscription %s consuming from topic %s", subscriptionName, topicName)
}
