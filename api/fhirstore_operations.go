package api

import (
	"fmt"

	"github.com/imaware/fhir-operator/api/v1alpha1"
	fhirv1alpha1 "github.com/imaware/fhir-operator/api/v1alpha1"
	"google.golang.org/api/googleapi"
	"google.golang.org/api/healthcare/v1"
	ctrl "sigs.k8s.io/controller-runtime"
)

var (
	logger = ctrl.Log.WithName("fhirstore_operations.go")
)

const DATSET_ERROR_NOT_FOUND = 403
const FHIR_STORE_ERROR_NOT_FOUND = 404

// Perform a fhirstore healthcare api in-order to make a decision to delete a fhirstore if the
// fhirstore does exist
func ReadAndOrDeleteFHIRStore(fhirStoreGetCall FHIRStoreClientGetCall, fhirStoreDeleteCall FHIRStoreClientDeleteCall, fhirStore *fhirv1alpha1.FhirStore) error {
	// check if fhir store exists
	// if it exists we break early
	exists, err := fhirStoreExists(fhirStoreGetCall, fhirStore)
	if err != nil {
		return err
	}
	if !exists {
		logger.V(1).Info(fmt.Sprintf("Fhirstore %v does not exist skipping delete for resource %v in namespace %v", fhirStore.Spec.FhirStoreID, fhirStore.Name, fhirStore.Namespace))
		return nil
	}
	// means fhir store does not exist create it
	err = deleteFhirStore(fhirStoreDeleteCall, fhirStore)
	if err != nil {
		return err
	}
	return nil
}

// Perform a get on the fhir store for it's IAM policy settings and return the policy.
func ReadFHIRStoreIAMPolicy(fhirStoreIAMPolicyGetCall FHIRStoreClientIAMPolicyGetCall, fhirStore *fhirv1alpha1.FhirStore) (*healthcare.Policy, error) {
	policy, err := GetFHIRStoreIAMPolicy(fhirStoreIAMPolicyGetCall)
	if err != nil {
		logger.V(1).Error(err, fmt.Sprintf("Failed to get fhirstore iam policy for resource %v in namespace %v", fhirStore.Name, fhirStore.Namespace))
		fhirStore.Status.Status = FAILED
		fhirStore.Status.Message = FHIRStoreCreateFailedStatus(err.Error())
		return nil, fmt.Errorf(fmt.Sprintf("Failed to get fhirstore iam policy for resource %v in namespace %v", fhirStore.Name, fhirStore.Namespace))
	}
	return policy, nil
}

// Perform an update on the fhir store's IAM policy setting. update the fhirstore object's status based on a failed
// update action. An error will be returned if the API call fails
func CreateOrUpdateFHIRStoreIAMPolicy(fhirStoreIAMPolicyCreateOrUpdateCall FHIRStoreClientIAMPolicyCreateOrUpdateCall, fhirStore *fhirv1alpha1.FhirStore) error {
	_, err := UpdateFHIRStoreIAMPolicy(fhirStoreIAMPolicyCreateOrUpdateCall)
	if err != nil {
		logger.V(1).Error(err, fmt.Sprintf("Failed to create fhirstore %v IAM policy for resource %v in namespace %v", fhirStore.Spec.FhirStoreID, fhirStore.Name, fhirStore.Namespace))
		fhirStore.Status.Status = FAILED
		fhirStore.Status.Message = FHIRStoreCreateFailedStatus(err.Error())
		return fmt.Errorf(fmt.Sprintf("Failed to create fhirstore %v IAM policy for resource %v in namespace %v", fhirStore.Spec.FhirStoreID, fhirStore.Name, fhirStore.Namespace))
	}
	return nil
}

// Perform a patch on a fhir store to update its settings. An error is returned if the API call fails
func PatchFhirStore(fhirStorePatchCall FHIRStoreClientPatchCall, fhirStore *fhirv1alpha1.FhirStore) error {
	_, err := PatchFHIRStore(fhirStorePatchCall)
	if err != nil {
		logger.V(1).Error(err, fmt.Sprintf("Failed to patch fhirstore %v for resource %v in namespace %v", fhirStore.Spec.FhirStoreID, fhirStore.Name, fhirStore.Namespace))
		fhirStore.Status.Status = FAILED
		fhirStore.Status.Message = FHIRStoreCreateFailedStatus(err.Error())
		return fmt.Errorf(fmt.Sprintf("Failed to patch fhirstore %v  for resource %v in namespace %v", fhirStore.Spec.FhirStoreID, fhirStore.Name, fhirStore.Namespace))
	}
	return nil
}

// Perform a read on the dataset and fhirstore healthcare api in-order to make a decision to create a fhirstore if the
// fhirstore does not exist
func ReadAndOrCreateFHIRStore(datasetGetCall DatastoreClientGetCall, fhirStoreGetCall FHIRStoreClientGetCall, fhirStoreCreateCall FHRIStoreClientCreateCall, fhirStore *fhirv1alpha1.FhirStore) error {
	// make sure dataset exists
	err := datasetExists(datasetGetCall, fhirStore)
	if err != nil {
		return err
	}
	// check if fhir store exists
	// if it exists we break early
	fhirStorExists, err := fhirStoreExists(fhirStoreGetCall, fhirStore)
	if err != nil {
		return err
	} else if !fhirStorExists {
		// means fhir store does not exist create it
		err = createFhirStore(fhirStoreCreateCall, fhirStore)
		if err != nil {
			return err
		}
	}
	return nil
}

// Create the fhirstore and update the fhirstore object's status based on a failed
// create. An error will be returned if the api call fails.
func createFhirStore(fhirStoreCreateCall FHRIStoreClientCreateCall, fhirStore *fhirv1alpha1.FhirStore) error {
	_, err := CreateFHIRStore(fhirStoreCreateCall)
	if err != nil {
		gcpErr, ok := err.(*googleapi.Error)
		if ok {
			logger.V(1).Error(gcpErr, fmt.Sprintf("Failed to create fhirstore %v for resource %v in namespace %v", fhirStore.Spec.FhirStoreID, fhirStore.Name, fhirStore.Namespace))
			fhirStore.Status.Status = FAILED
			fhirStore.Status.Message = FHIRStoreCreateFailedStatus(err.Error())
			return fmt.Errorf(fmt.Sprintf("Failed to create fhirstore %v for resource %v in namespace %v", fhirStore.Spec.FhirStoreID, fhirStore.Name, fhirStore.Namespace))
		} else {
			fhirStore.Status.Status = FAILED
			fhirStore.Status.Message = FHIRStoreCreateFailedStatus(err.Error())
			return fmt.Errorf(fmt.Sprintf("Create fhirstore %v internal error for resource %v in namespace %v: %v", fhirStore.Spec.FhirStoreID, fhirStore.Name, fhirStore.Namespace, err))
		}
	} else {
		logger.Info(fmt.Sprintf("FHIR store created %v in dataset %v for resource %v in namespace %v", fhirStore.Spec.FhirStoreID, fhirStore.Spec.DatasetID, fhirStore.Name, fhirStore.Namespace))
		fhirStore.Status.Status = CREATING
		fhirStore.Status.Message = FHIRStoreCreatingStatus(fhirStore.Spec.DatasetID, fhirStore.Spec.FhirStoreID)
		return nil
	}
}

// Delete the fhirstore and update the fhirstore object's status based on the
// delete if it fails. An error will be returned if the api call fails.
func deleteFhirStore(fhirStoreDeleteCall FHIRStoreClientDeleteCall, fhirStore *fhirv1alpha1.FhirStore) error {
	_, err := DeleteFHIRStore(fhirStoreDeleteCall)
	if err != nil {
		gcpErr, ok := err.(*googleapi.Error)
		if ok {
			logger.V(1).Error(gcpErr, fmt.Sprintf("Failed to delete fhirstore %v for resource %v in namespace %v", fhirStore.Spec.FhirStoreID, fhirStore.Name, fhirStore.Namespace))
			fhirStore.Status.Status = FAILED
			fhirStore.Status.Message = FHIRStoreDeleteFailedStatus(err.Error())
			return fmt.Errorf("Failed to delete FHIR store")
		} else {
			fhirStore.Status.Status = FAILED
			fhirStore.Status.Message = FHIRStoreDeleteFailedStatus(err.Error())
			return fmt.Errorf("Delete fhirstore internal error: %v", err)
		}
	} else {
		logger.Info(fmt.Sprintf("FHIR store deleted %v in dataset %v for resource %v in namesapce %v", fhirStore.Spec.FhirStoreID, fhirStore.Spec.DatasetID, fhirStore.Name, fhirStore.Namespace))
		return nil
	}
}

// Check if the Dataset exists in the GCP healthcare API and update the fhirstore object's status based on the
// check. An error will be returned if the api call fails.
func datasetExists(datasetGetCall DatastoreClientGetCall, fhirStore *fhirv1alpha1.FhirStore) error {
	_, err := GetDataset(datasetGetCall)
	if err != nil {
		gcpErr, ok := err.(*googleapi.Error)
		if ok {
			code := gcpErr.Code
			if code == DATSET_ERROR_NOT_FOUND {
				logger.V(1).Error(gcpErr, fmt.Sprintf("Failed to get datastore %v for resource %v in namesapce %v. This can be due to either bad permissions or resource does not exist", fhirStore.Spec.DatasetID, fhirStore.Name, fhirStore.Namespace))
				fhirStore.Status.Status = FAILED
				fhirStore.Status.Message = DatasetNotFoundOrPermissionsInvalidStatus(fhirStore.Spec.DatasetID, gcpErr)
				return fmt.Errorf("Invalid credentials or datastore does not exist")
			} else {
				logger.V(1).Error(err, "Get dataset call failed")
				fhirStore.Status.Status = FAILED
				fhirStore.Status.Message = GetInternalError(err.Error())
				return fmt.Errorf("Get dataset call failed")
			}
		} else {
			fhirStore.Status.Status = FAILED
			fhirStore.Status.Message = GetInternalError(err.Error())
			return fmt.Errorf("Get dataset internal error: %v", err)
		}
	} else {
		logger.V(1).Info(fmt.Sprintf("Found dataset with ID %v", fhirStore.Spec.DatasetID))
		return nil
	}
}

// Check if the Fhirstore exists in the GCP healthcare API and update the fhirstore object's status based on the
// check. An error will be returned if the api call fails.
func fhirStoreExists(fhirstoreGetCall FHIRStoreClientGetCall, fhirStore *fhirv1alpha1.FhirStore) (bool, error) {
	_, err := GetFHIRStore(fhirstoreGetCall)
	if err != nil {
		gcpErr, ok := err.(*googleapi.Error)
		if ok {
			code := gcpErr.Code
			// 403 can mean either bad credentials or it does not exist, try to create
			if code == FHIR_STORE_ERROR_NOT_FOUND {
				logger.V(1).Info(fmt.Sprintf("Fhirstore %v not found", fhirStore.Spec.FhirStoreID))
				return false, nil
			} else {
				logger.V(1).Error(err, "GCP Healthcare FHIR GET api error")
				fhirStore.Status.Status = FAILED
				fhirStore.Status.Message = GetInternalError(err.Error())
				return false, fmt.Errorf("GCP Healthcare FHIR GET API error")
			}
		} else {
			fhirStore.Status.Status = FAILED
			fhirStore.Status.Message = GetInternalError(err.Error())
			return false, fmt.Errorf("Get fhirstore error: %v", err)
		}
	} else {
		logger.V(1).Info(fmt.Sprintf("Fhirstore %v exists", fhirStore.Spec.FhirStoreID))
		return true, nil
	}
}

// Given the auth spec of the fhirStore generate the google policy bindings to attach to the
// fhir store IAM api
func GenerateIAMPolicyBindings(newBindings map[string]v1alpha1.FhirStoreSpecAuth) []*healthcare.Binding {
	var policyBindings []*healthcare.Binding
	for role, members := range newBindings {
		policyBinding := &healthcare.Binding{
			Role:    role,
			Members: members.Members,
		}
		policyBindings = append(policyBindings, policyBinding)
	}
	return policyBindings
}

// Given the bigquerryConfig spec fo the fhirStore generate the corresponding streaming configs for the fhirstore object
func GenerateFhirStoreBigQueryConfigs(bigquerryConfigs []v1alpha1.FhirStoreSpecOptionsBigquery) []*healthcare.StreamConfig {

	streamingConfigs := []*healthcare.StreamConfig{}
	for _, config := range bigquerryConfigs {
		streamingConfig := &healthcare.StreamConfig{
			BigqueryDestination: &healthcare.GoogleCloudHealthcareV1FhirBigQueryDestination{
				DatasetUri:       config.Id,
				WriteDisposition: "WRITE_APPEND",
				SchemaConfig: &healthcare.SchemaConfig{
					SchemaType: "ANALYTICS",
				},
			},
		}
		streamingConfigs = append(streamingConfigs, streamingConfig)
	}
	return streamingConfigs
}
