/*
Copyright 2021.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controllers

import (
	"context"
	"fmt"
	"time"

	"github.com/imaware/fhir-operator/api"
	"google.golang.org/api/healthcare/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/predicate"

	"github.com/imaware/fhir-operator/api/utils"
	"github.com/imaware/fhir-operator/api/v1alpha1"
	fhirv1alpha1 "github.com/imaware/fhir-operator/api/v1alpha1"
)

var (
	logger          = ctrl.Log.WithName("fhirstore_controller.go")
	configFhirStore *api.ConfigVars
)

// FhirStoreReconciler reconciles a FhirStore object
type FhirStoreReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

type HealthCareInfoStruct struct {
	fhireStoreObjectName string
	fhirStoreID          string
	datasetID            string
}

const FHISTORE_FINALIZER = "fhir.imaware.com/finalizer"

//+kubebuilder:rbac:groups=fhir.imaware.com,resources=fhirstores,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=fhir.imaware.com,resources=fhirstores/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=fhir.imaware.com,resources=fhirstores/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the FhirStore object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.10.0/pkg/reconcile
func (r *FhirStoreReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	var result = ctrl.Result{}
	var err error
	// Get all valuable metadata for object
	fhirStore := &fhirv1alpha1.FhirStore{}
	err = r.Get(context.TODO(), req.NamespacedName, fhirStore)
	if err != nil {
		if errors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile request.
			// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
			// Return and don't requeue
			logger.V(1).Info("Fhirstore resource not found.")
			// set err to nil to not requeue the request
			err = nil
		} else {
			logger.V(1).Error(err, "Failed to get Fhirstore")
		}
	} else {
		var datasetGetCall *healthcare.ProjectsLocationsDatasetsGetCall
		datasetGetCall, err = api.BuildDatasetGetCall(configFhirStore.GCPProject, configFhirStore.GCPLocation, fhirStore.Spec.DatasetID)
		if err != nil {
			logger.V(1).Error(err, "Failed to build Dataset get call", "fhirStoreID", fhirStore.Spec.FhirStoreID, "datasetID", fhirStore.Spec.DatasetID, "fhireStoreObjectName", fhirStore.Name)
		}
		var fhirStoreGetCall *healthcare.ProjectsLocationsDatasetsFhirStoresGetCall
		fhirStoreGetCall, err = api.BuildFHIRStoreGetCall(configFhirStore.GCPProject, configFhirStore.GCPLocation, fhirStore.Spec.DatasetID, fhirStore.Spec.FhirStoreID)
		if err != nil {
			logger.V(1).Error(err, "Failed to build Fhir store get call", "fhirStoreID", fhirStore.Spec.FhirStoreID, "datasetID", fhirStore.Spec.DatasetID, "fhireStoreObjectName", fhirStore.Name)
		}
		if datasetGetCall != nil && fhirStoreGetCall != nil {
			// Check if the FhirStore instance is marked to be deleted, which is
			// indicated by the deletion timestamp being set.
			isFhirStoreMarkedToBeDeleted := fhirStore.GetDeletionTimestamp() != nil
			if isFhirStoreMarkedToBeDeleted {
				if fhirStore.Spec.Options.PreventDelete {
					logger.Info(fmt.Sprintf("Fhirstore %v can not be deleted %v in namesapce %v as preventDelete option set will remove store", fhirStore.Spec.FhirStoreID, fhirStore.Name, fhirStore.Namespace))
					utils.RemoveFinalizer(fhirStore, FHISTORE_FINALIZER)
					updateError := r.Update(ctx, fhirStore)
					if updateError != nil {
						logger.V(1).Error(err, fmt.Sprintf("Failed to add finalizer for fhirstore resource %v", fhirStore.Name))
						err = updateError
					} else {
						err = nil
					}
				} else {
					result, err = deleteFhirStoreLoop(fhirStore, datasetGetCall, fhirStoreGetCall)
					if err == nil {
						utils.RemoveFinalizer(fhirStore, FHISTORE_FINALIZER)
						updateError := r.Update(ctx, fhirStore)
						if updateError != nil {
							logger.V(1).Error(err, fmt.Sprintf("Failed to remove finalizer for fhirstore resource %v", fhirStore.Name))
							err = updateError
						}
					}
				}
			} else {
				// Add finalizer for this CR
				utils.AddFinalizer(fhirStore, FHISTORE_FINALIZER)
				updateError := r.Update(ctx, fhirStore)
				if err != nil {
					logger.V(1).Error(err, fmt.Sprintf("Failed to add finalizer for fhirstore resource %v", fhirStore.Name))
					err = updateError
				} else {
					var fhirStoreCreateCall *healthcare.ProjectsLocationsDatasetsFhirStoresCreateCall
					fhirStoreCreateCall, err = api.BuildFHIRStoreCreateCall(configFhirStore.GCPProject, configFhirStore.GCPLocation, fhirStore.Spec.DatasetID, "R4", fhirStore.Spec.FhirStoreID)
					if err != nil {
						logger.V(1).Error(err, "Failed to build Fhir store create call", "fhirStoreID", fhirStore.Spec.FhirStoreID, "datasetID", fhirStore.Spec.DatasetID, "fhireStoreObjectName", fhirStore.Name)
					} else {
						result, err = createFhirStoreLoop(fhirStore, datasetGetCall, fhirStoreGetCall, fhirStoreCreateCall)
					}
				}

				if checkExportCondition(fhirStore) {
					result, err = exportFhirStoreLoop(fhirStore)
				}

			}
		}

		r.Status().Update(ctx, fhirStore)
	}
	return result, err
}

// SetupWithManager sets up the controller with the Manager.
func (r *FhirStoreReconciler) SetupWithManager(mgr ctrl.Manager, conf *api.ConfigVars) error {
	configFhirStore = conf
	pred := predicate.GenerationChangedPredicate{}
	logger.V(1).Info("Starting reconcile loop for fhirstore_controller.go")
	return ctrl.NewControllerManagedBy(mgr).
		For(&fhirv1alpha1.FhirStore{}).WithEventFilter(pred).
		Complete(r)
}

func createFhirStoreLoop(fhirStore *v1alpha1.FhirStore, datasetGetCall api.DatastoreClientGetCall, fhirStoreGetCall api.FHIRStoreClientGetCall, fhirStoreCreateCall api.FHRIStoreClientCreateCall) (ctrl.Result, error) {
	logger.Info("Checking if Fhirstore should be created", "fhirStoreID", fhirStore.Spec.FhirStoreID, "datasetID", fhirStore.Spec.DatasetID, "fhireStoreObjectName", fhirStore.Name)
	err := api.ReadAndOrCreateFHIRStore(datasetGetCall, fhirStoreGetCall, fhirStoreCreateCall, fhirStore)
	if err != nil {
		logger.Error(err, "Something went wrong during FHIR store creation", "fhirStoreID", fhirStore.Spec.FhirStoreID, "datasetID", fhirStore.Spec.DatasetID, "fhireStoreObjectName", fhirStore.Name)
		return ctrl.Result{}, err
	}
	// create IAM policy if needed
	if len(fhirStore.Spec.Auth) > 0 {
		fhirStoreIAMPolicyGetCall, err := api.BuildFhirStoreGetIAMPolicyRequest(configFhirStore.GCPProject, configFhirStore.GCPLocation, fhirStore.Spec.DatasetID, fhirStore.Spec.FhirStoreID)
		if err != nil {
			logger.V(1).Error(err, "Failed to build Fhir store IAM get call", "fhirStoreID", fhirStore.Spec.FhirStoreID, "datasetID", fhirStore.Spec.DatasetID, "fhireStoreObjectName", fhirStore.Name)
			return ctrl.Result{}, err
		}
		// Get the policy to make sure we dont miss anything in updating or creating policy
		policy, err := api.ReadFHIRStoreIAMPolicy(fhirStoreIAMPolicyGetCall, fhirStore)
		if err != nil {
			logger.Error(err, "Something went wrong during FHIR store IAM policy read", "fhirStoreID", fhirStore.Spec.FhirStoreID, "datasetID", fhirStore.Spec.DatasetID, "fhireStoreObjectName", fhirStore.Name)
			return ctrl.Result{}, nil
		}
		policy.Bindings = api.GenerateIAMPolicyBindings(fhirStore.Spec.Auth)
		fhirStoreIAMPolicyCreateOrUpdateCall, err := api.BuildFhirStoreUpdateOrCreateIAMPolicyRequest(configFhirStore.GCPProject, configFhirStore.GCPLocation, fhirStore.Spec.DatasetID, fhirStore.Spec.FhirStoreID, policy)
		if err != nil {
			logger.V(1).Error(err, "Failed to build Fhir store IAM update call", "fhirStoreID", fhirStore.Spec.FhirStoreID, "datasetID", fhirStore.Spec.DatasetID, "fhireStoreObjectName", fhirStore.Name)
			return ctrl.Result{}, err
		}
		// create or update the policy for the fhir store
		err = api.CreateOrUpdateFHIRStoreIAMPolicy(fhirStoreIAMPolicyCreateOrUpdateCall, fhirStore)
		if err != nil {
			logger.Error(err, "Something went wrong during FHIR store IAM policy create or update", "fhirStoreID", fhirStore.Spec.FhirStoreID, "datasetID", fhirStore.Spec.DatasetID, "fhireStoreObjectName", fhirStore.Name)
			fhirStore.Status.Status = api.FAILED
			fhirStore.Status.Message = api.FHIRStoreFailedIAMPolicy(err.Error())
			return ctrl.Result{}, nil
		}

	}
	// configure the big querry streaming if needed
	if len(fhirStore.Spec.Options.Bigquery) > 0 {
		fhirStoreGCP, err := api.GetFHIRStore(fhirStoreGetCall)
		if err != nil {
			logger.Error(err, "Something went wrong during FHIR store read", "fhirStoreID", fhirStore.Spec.FhirStoreID, "datasetID", fhirStore.Spec.DatasetID, "fhireStoreObjectName", fhirStore.Name)
			return ctrl.Result{}, nil
		}
		fhirStoreStreamingConfigs := api.GenerateFhirStoreBigQueryConfigs(fhirStore.Spec.Options.Bigquery)
		fhirStoreGCP.StreamConfigs = fhirStoreStreamingConfigs
		fhirStorePatchCall, err := api.BuildFhirStorePatchCall(configFhirStore.GCPProject, configFhirStore.GCPLocation, fhirStore.Spec.DatasetID, fhirStore.Spec.FhirStoreID, fhirStoreGCP)
		if err != nil {
			logger.V(1).Error(err, "Failed to build Fhir store patch call", "fhirStoreID", fhirStore.Spec.FhirStoreID, "datasetID", fhirStore.Spec.DatasetID, "fhireStoreObjectName", fhirStore.Name)
			return ctrl.Result{}, err
		}
		// allow for updating the streamConfigs path in the resource
		fhirStorePatchCall.UpdateMask("streamConfigs")

		// create or update the policy for the fhir store
		err = api.PatchFhirStore(fhirStorePatchCall, fhirStore)
		if err != nil {
			logger.Error(err, "Something went wrong during FHIR store IAM policy create or update", "fhirStoreID", fhirStore.Spec.FhirStoreID, "datasetID", fhirStore.Spec.DatasetID, "fhireStoreObjectName", fhirStore.Name)
			fhirStore.Status.Status = api.FAILED
			fhirStore.Status.Message = api.FHIRStoreFailedPatch(err.Error())
			return ctrl.Result{}, nil
		}
	}
	// if we are at this point the fhir store is up and running
	logger.Info(fmt.Sprintf("Fhirstore %v created for resource %v in namesapce %v", fhirStore.Spec.FhirStoreID, fhirStore.Name, fhirStore.Namespace))
	fhirStore.Status.Status = api.CREATED
	fhirStore.Status.Message = api.FHIRStoreCreatedStatus(fhirStore.Spec.FhirStoreID)
	return ctrl.Result{}, nil
}

func deleteFhirStoreLoop(fhirStore *fhirv1alpha1.FhirStore, datasetGetCall api.DatastoreClientGetCall, fhirStoreGetCall api.FHIRStoreClientGetCall) (ctrl.Result, error) {
	fhirStoreDeleteCall, err := api.BuildFHIRStoreDeleteCall(configFhirStore.GCPProject, configFhirStore.GCPLocation, fhirStore.Spec.DatasetID, fhirStore.Spec.FhirStoreID)
	if err != nil {
		logger.Error(err, "Failed to build Fhir store delete call", "fhirStoreID", fhirStore.Spec.FhirStoreID, "datasetID", fhirStore.Spec.DatasetID, "fhireStoreObjectName", fhirStore.Name)
		return ctrl.Result{}, err
	}
	err = api.ReadAndOrDeleteFHIRStore(fhirStoreGetCall, fhirStoreDeleteCall, fhirStore)
	//r.Status().Update(ctx, fhirStore)
	if err != nil {
		logger.Error(err, "Something went wrong during FHIR store delete", "fhirStoreID", fhirStore.Spec.FhirStoreID, "datasetID", fhirStore.Spec.DatasetID, "fhireStoreObjectName", fhirStore.Name)
		return ctrl.Result{}, nil
	}
	return ctrl.Result{}, nil
}

func checkExportCondition(fhirStore *v1alpha1.FhirStore) bool {

	// check if exports are enabled
	if !fhirStore.Spec.ExportOptions.EnableExports {
		return false
	}

	lastExport := fhirStore.Status.LastExported

	// retry failed exports no matter the time
	if lastExport == "FAILED" {
		logger.V(1).Info("Failed to export %s last time, reconciling")
		return true
	}

	longFormat := "2006-01-02 15:04:05.999999999 -0700 MST"
	lastExportTime, err := time.Parse(longFormat, lastExport)
	if err != nil {
		logger.Error(err, "Status error: invalid time status")
		return false
	}

	// convert frequency string to time object
	frequency, err := time.ParseDuration(fhirStore.Spec.ExportOptions.Frequency)
	if err != nil {
		logger.Error(err, "Config error: invalid frequency parameter")
		return false
	}

	if time.Since(lastExportTime) > frequency {
		return true
	}

	return false
}

func exportFhirStoreLoop(fhirStore *v1alpha1.FhirStore) (ctrl.Result, error) {
	if fhirStore == nil {
		return ctrl.Result{}, nil
	}

	if configFhirStore == nil {
		return ctrl.Result{}, nil
	}

	fhirStoreExportCall, err := api.BuildFHIRStoreExportCall(configFhirStore.GCPProject, configFhirStore.GCPLocation, fhirStore.Spec.DatasetID, fhirStore.Spec.FhirStoreID, fhirStore.Spec.ExportOptions.Location)
	if err != nil {
		logger.Error(err, "Failed to build Fhir store export call", "fhirStoreID", fhirStore.Spec.FhirStoreID, "datasetID", fhirStore.Spec.DatasetID, "fhireStoreObjectName", fhirStore.Name)
		return ctrl.Result{}, err
	}

	err = api.ExportFhirStore(fhirStoreExportCall, fhirStore)
	//r.Status().Update(ctx, fhirStore)
	if err != nil {
		logger.Error(err, "Something went wrong during FHIR store export", "fhirStoreID", fhirStore.Spec.FhirStoreID, "datasetID", fhirStore.Spec.DatasetID, "fhireStoreObjectName", fhirStore.Name)
		return ctrl.Result{}, err
	}
	return ctrl.Result{}, nil
}
