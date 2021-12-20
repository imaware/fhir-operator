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

// Core logic to delete the fhir resource.
// First we delete the resource from the database and check for if the resource still exists in the database.
// The reson for this is the delete will always return a 200. An error will be returned if the delete or get call fails, or if the resource
// is still found in the database after delete
func DeleteAndReadFHIRStoreResource(fhirStoreResourceDeleteCall FHIRStoreResourceClientDeleteCall, fhirResourceGetCall FHIRStoreResourceClientGetCall, fhirStoreResource *fhirv1alpha1.FhirResource) error {
	// Delete response body will always return a 200 if a fail or not so need to do a get
	// to make sure the resource is actually deleted
	_, err := DeleteFHIRResource(fhirStoreResourceDeleteCall)
	if err != nil {
		fhirResourceLogger.V(1).Error(err, fmt.Sprintf("Delete call failed for resource %v in namespace %v", fhirStoreResource.Name, fhirStoreResource.Namespace))
		fhirStoreResource.Status.Status = FAILED
		fhirStoreResource.Status.Message = FHIRStoreResourceDeleteFailedStatus(fhirStoreResource.Name, err.Error())
	} else {
		// need to do a get to make sure the resource was actually deleted
		var exists bool
		exists, err = fhirResourceExists(fhirResourceGetCall, fhirStoreResource)
		// resource not found succesful delete
		if err == nil && !exists {
			fhirResourceLogger.Info(fmt.Sprintf("Deleted fhir resource %v in namespace %v", fhirStoreResource.Name, fhirStoreResource.Namespace))
			fhirStoreResource.Status.Status = DELETED
			fhirStoreResource.Status.Message = FHIRStoreResourceDeletetatus()
			// resource is still in the database, something went wrong
		} else {
			fhirStoreResource.Status.Status = FAILED
			fhirStoreResource.Status.Message = FHIRStoreResourceDeleteFailedStatus(fhirStoreResource.Name, fmt.Sprintf("Failed to delete fhir resource %v in namespace %v as it still exists", fhirStoreResource.Name, fhirStoreResource.Namespace))
			err = fmt.Errorf("Failed to delete fhir resource %v in namespace %v as it still exists", fhirStoreResource.Name, fhirStoreResource.Namespace)
		}
	}
	return err
}

// Core logic to create the fhir resource.
// if the parent fhir resource is not yet in the database we requeue the event else we try to create it
// and fail on an api call error. An error returned does not mean we requeue the request, just an indicator of a
// problem in the process. The boolean is the only indicator of an event requeue.
func CreateOrUpdateFHIRResource(fhirStoreResourceUpdateCall FHIRStoreResourceClientUpdateCall, fhirStoreResource *fhirv1alpha1.FhirResource) (bool, error) {
	resp, err := UpdateFHIRResource(fhirStoreResourceUpdateCall)
	if err != nil {
		fhirResourceLogger.V(1).Error(err, fmt.Sprintf("Create or update resource call failed for resource %v in namespace %v", fhirStoreResource.Name, fhirStoreResource.Namespace))
		fhirStoreResource.Status.Status = FAILED
		fhirStoreResource.Status.Message = FHIRStoreResourceCreateOrUpdateFailedStatus(fhirStoreResource.Name, err.Error())
		return false, err
	}
	// get the response body for logging purposes, if we can't just set to an empty string
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
		// Something is either wrong with the representation or an internal error on the api build
	} else if resp.StatusCode > BAD_RESPONSE_CODE_THRESH {
		fhirStoreResource.Status.Status = FAILED
		fhirStoreResource.Status.Message = FHIRStoreResourceCreateOrUpdateFailedStatus(fhirStoreResource.Name, string(respBytes))
		return false, fmt.Errorf("Failed to create or update resource %v in namespace %v", fhirStoreResource.Name, fhirStoreResource.Namespace)
	}
	// Resource created
	fhirResourceLogger.V(1).Info(fmt.Sprintf("Created or updated fhir resource %v in namespace %v", fhirStoreResource.Name, fhirStoreResource.Namespace))
	fhirStoreResource.Status.Status = CREATED
	fhirStoreResource.Status.Message = FHIRStoreResourceCreatedorUpdatedStatus()
	return false, nil
}

// Parse a json string representation of a fhir resource and return the id associated with
// the resoure. An error is returned if the json is invalid.
func GetFHIRIResourceID(jsonString string) (string, error) {
	var ID = "id"
	resourceJson, err := utils.JsonStringToMap(jsonString)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%v", resourceJson[ID]), nil
}

// Check if the fhir resource exists in the fhir store. An error is returned on a bad api call.
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
