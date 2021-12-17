package api

import (
	"fmt"
	"io/ioutil"

	"github.com/imaware/fhir-operator/api/utils"
	fhirv1alpha1 "github.com/imaware/fhir-operator/api/v1alpha1"
	ctrl "sigs.k8s.io/controller-runtime"
)

var (
	fhirResourceLogger = ctrl.Log.WithName("fhirresources_operations.go")
)

// A Failed status code is above this for resource calls
const BAD_RESPONSE_CODE_THRESH = 299

// This response code is not a failure, just means that the parent resource has not been created yet
const UPDATE_RESPONSE_CODE_REFERENCE_NOT_FOUND = 400

func DeleteAndReadFHIRStoreResource(fhirStoreResourceDeleteCall FHIRStoreResourceClientDeleteCall, fhirResourceGetCall FHIRStoreResourceClientGetCall, fhirStoreResource *fhirv1alpha1.FhirResource) error {
	err := deleteFHIRResource(fhirStoreResourceDeleteCall, fhirStoreResource)
	if err != nil {
		return err
	}
	// need to do a get to make sure the resource was actually deleted
	exists, err := fhirResourceExists(fhirResourceGetCall, fhirStoreResource)
	if err != nil {
		return err
	}
	if !exists {
		fhirResourceLogger.Info(fmt.Sprintf("Deleted fhir resource %v in namespace %v", fhirStoreResource.Name, fhirStoreResource.Namespace))
		fhirStoreResource.Status.Status = DELETED
		fhirStoreResource.Status.Message = FHIRStoreResourceDeletetatus()
		return nil
	} else {
		fhirResourceLogger.Info(fmt.Sprintf("Failed to delete fhir resource %v in namespace %v", fhirStoreResource.Name, fhirStoreResource.Namespace))
		fhirStoreResource.Status.Status = FAILED
		fhirStoreResource.Status.Message = FHIRStoreResourceDeleteFailedStatus(fhirStoreResource.Name, fmt.Sprintf("Failed to delete fhir resource %v in namespace %v as it still exists", fhirStoreResource.Name, fhirStoreResource.Namespace))
		return fmt.Errorf("Failed to delete fhir resource %v in namespace %v as it still exists", fhirStoreResource.Name, fhirStoreResource.Namespace)
	}
}

func CreateOrUpdateFHIRResource(fhirStoreResourceUpdateCall FHIRStoreResourceClientUpdateCall, fhirStoreResource *fhirv1alpha1.FhirResource) (bool, error) {
	return createOrUpdateFHIRResource(fhirStoreResourceUpdateCall, fhirStoreResource)
}

func GetFHIRIResourceID(jsonString string) (string, error) {
	var ID = "id"
	resourceJson, err := utils.JsonStringToMap(jsonString)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%v", resourceJson[ID]), nil
}

func createOrUpdateFHIRResource(fhirStoreResourceUpdateCall FHIRStoreResourceClientUpdateCall, fhirStoreResource *fhirv1alpha1.FhirResource) (bool, error) {
	resp, err := UpdateFHIRResource(fhirStoreResourceUpdateCall)
	if err != nil {
		fhirResourceLogger.V(1).Error(err, fmt.Sprintf("Create or update resource call failed for resource %v in namespace %v", fhirStoreResource.Name, fhirStoreResource.Namespace))
		fhirStoreResource.Status.Status = FAILED
		fhirStoreResource.Status.Message = FHIRStoreResourceCreateOrUpdateFailedStatus(fhirStoreResource.Name, err.Error())
		return false, err
	}
	defer resp.Body.Close()
	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fhirResourceLogger.V(1).Error(err, "could not read response")
		respBytes = []byte("")
	}
	fhirResourceLogger.V(1).Info(fmt.Sprintf("Update: status %d %s: %s", resp.StatusCode, resp.Status, respBytes))
	// requeue the event as it means the parent resource has not been created yet
	if resp.StatusCode == UPDATE_RESPONSE_CODE_REFERENCE_NOT_FOUND {
		fhirStoreResource.Status.Status = PENDING
		fhirStoreResource.Status.Message = FHIRStoreResourcePendingOnParentResourceStatus()
		return true, nil
	} else if resp.StatusCode > BAD_RESPONSE_CODE_THRESH {
		fhirStoreResource.Status.Status = FAILED
		fhirStoreResource.Status.Message = FHIRStoreResourceCreateOrUpdateFailedStatus(fhirStoreResource.Name, string(respBytes))
		return false, fmt.Errorf("Failed to create or update resource %v in namespace %v", fhirStoreResource.Name, fhirStoreResource.Namespace)
	}
	fhirResourceLogger.Info(fmt.Sprintf("Created or updated fhir resource %v in namespace %v", fhirStoreResource.Name, fhirStoreResource.Namespace))
	fhirStoreResource.Status.Status = CREATED
	fhirStoreResource.Status.Message = FHIRStoreResourceCreatedorUpdatedStatus()
	return false, nil
}

func deleteFHIRResource(fhirStoreResourceDeleteCall FHIRStoreResourceClientDeleteCall, fhirStoreResource *fhirv1alpha1.FhirResource) error {
	// Delete response body will always return a 200 if a fail or not so need to do a get
	// to make sure the resource is actually deleted
	_, err := DeleteFHIRResource(fhirStoreResourceDeleteCall)
	if err != nil {
		fhirResourceLogger.V(1).Error(err, fmt.Sprintf("Delete call failed for resource %v in namespace %v", fhirStoreResource.Name, fhirStoreResource.Namespace))
		fhirStoreResource.Status.Status = FAILED
		fhirStoreResource.Status.Message = FHIRStoreResourceDeleteFailedStatus(fhirStoreResource.Name, err.Error())
		return err
	}
	return nil
}

func fhirResourceExists(fhirStoreResourceGetCall FHIRStoreResourceClientGetCall, fhirStoreResource *fhirv1alpha1.FhirResource) (bool, error) {
	resp, err := GetFHIRResource(fhirStoreResourceGetCall)
	if err != nil {
		fhirResourceLogger.V(1).Error(err, fmt.Sprintf("Get resource call failed for resource %v in namespace %v", fhirStoreResource.Name, fhirStoreResource.Namespace))
		return false, err
	}
	defer resp.Body.Close()
	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fhirResourceLogger.V(1).Error(err, "could not read response")
		respBytes = []byte("")
	}
	fhirResourceLogger.V(1).Info(fmt.Sprintf("Get: status %d %s: %s", resp.StatusCode, resp.Status, string(respBytes)))
	if resp.StatusCode > BAD_RESPONSE_CODE_THRESH {
		fhirResourceLogger.V(1).Info(fmt.Sprintf("No FHIR resource found for resource %v in namespace %v", fhirStoreResource.Name, fhirStoreResource.Namespace))
		return false, nil
	}
	fhirResourceLogger.V(1).Info(fmt.Sprintf("Found fhir resource %v in namespace %v", fhirStoreResource.Name, fhirStoreResource.Namespace))
	return true, nil

}
