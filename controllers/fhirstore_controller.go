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

	"github.com/imaware/fhir-operator/api"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

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
	// Get all valuable metadata for object
	fhirStore := &fhirv1alpha1.FhirStore{}
	err := r.Get(context.TODO(), req.NamespacedName, fhirStore)
	if err != nil {
		if errors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile request.
			// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
			// Return and don't requeue
			logger.V(1).Info("Fhirstore resource not found.")
			return ctrl.Result{}, nil
		}
		// Error reading the object - requeue the request.
		logger.V(1).Error(err, "Failed to get Fhirstore")
		return ctrl.Result{}, err
	}
	datasetGetCall, err := api.BuildDatasetGetCall(configFhirStore.GCPProject, configFhirStore.GCPLocation, fhirStore.Spec.DatasetID)
	if err != nil {
		logger.Error(err, "Failed to build Dataset get call", "fhirStoreID", fhirStore.Spec.FhirStoreID, "datasetID", fhirStore.Spec.DatasetID, "fhireStoreObjectName", fhirStore.Name)
		return ctrl.Result{}, nil
	}
	fhirStoreGetCall, err := api.BuildFHIRStoreGetCall(configFhirStore.GCPProject, configFhirStore.GCPLocation, fhirStore.Spec.DatasetID, fhirStore.Spec.FhirStoreID)
	if err != nil {
		logger.Error(err, "Failed to build Fhir store get call", "fhirStoreID", fhirStore.Spec.FhirStoreID, "datasetID", fhirStore.Spec.DatasetID, "fhireStoreObjectName", fhirStore.Name)
		return ctrl.Result{}, nil
	}
	// Check if the FhirStore instance is marked to be deleted, which is
	// indicated by the deletion timestamp being set.
	isFhirStoreMarkedToBeDeleted := fhirStore.GetDeletionTimestamp() != nil
	if isFhirStoreMarkedToBeDeleted {
		fhirStoreDeleteCall, err := api.BuildFHIRStoreDeleteCall(configFhirStore.GCPProject, configFhirStore.GCPLocation, fhirStore.Spec.DatasetID, fhirStore.Spec.FhirStoreID)
		if err != nil {
			logger.Error(err, "Failed to build Fhir store delete call", "fhirStoreID", fhirStore.Spec.FhirStoreID, "datasetID", fhirStore.Spec.DatasetID, "fhireStoreObjectName", fhirStore.Name)
			return ctrl.Result{}, err
		}
		if controllerutil.ContainsFinalizer(fhirStore, FHISTORE_FINALIZER) {
			// Run finalization logic for fhirstore finalizaer. If the
			// finalization logic fails, don't remove the finalizer so
			// that we can retry during the next reconciliation.
			err := api.ReadAndOrDeleteFHIRStore(datasetGetCall, fhirStoreGetCall, fhirStoreDeleteCall, fhirStore)
			r.Status().Update(ctx, fhirStore)
			if err != nil {
				logger.Error(err, "Something went wrong during FHIR store creation", "fhirStoreID", fhirStore.Spec.FhirStoreID, "datasetID", fhirStore.Spec.DatasetID, "fhireStoreObjectName", fhirStore.Name)
				return ctrl.Result{}, nil
			}
			// Remove fhirstore finalizer. Once all finalizers have been
			// removed, the object will be deleted.
			err = removeFhirStoreFinalizer(fhirStore, r, ctx)
			if err != nil {
				return ctrl.Result{}, err
			}
		}
		return ctrl.Result{}, nil
	}
	// Add finalizer for this CR
	err = addFhirStoreFinalizer(fhirStore, r, ctx)
	if err != nil {
		return ctrl.Result{}, nil
	}
	logger.Info("Checking if Fhirstore should be created", "fhirStoreID", fhirStore.Spec.FhirStoreID, "datasetID", fhirStore.Spec.DatasetID, "fhireStoreObjectName", fhirStore.Name)
	fhirStoreCreateCall, err := api.BuildFHIRStoreCreateCall(configFhirStore.GCPProject, configFhirStore.GCPLocation, fhirStore.Spec.DatasetID, "R4", fhirStore.Spec.FhirStoreID)
	if err != nil {
		logger.Error(err, "Failed to build Fhir store create call", "fhirStoreID", fhirStore.Spec.FhirStoreID, "datasetID", fhirStore.Spec.DatasetID, "fhireStoreObjectName", fhirStore.Name)
	}
	enqueue, err := api.ReadAndOrCreateFHIRStore(datasetGetCall, fhirStoreGetCall, fhirStoreCreateCall, fhirStore)
	r.Status().Update(ctx, fhirStore)
	if err != nil {
		logger.Error(err, "Something went wrong during FHIR store creation", "fhirStoreID", fhirStore.Spec.FhirStoreID, "datasetID", fhirStore.Spec.DatasetID, "fhireStoreObjectName", fhirStore.Name)
		return ctrl.Result{Requeue: enqueue}, err
	} else {
		return ctrl.Result{Requeue: enqueue}, nil
	}
}

// SetupWithManager sets up the controller with the Manager.
func (r *FhirStoreReconciler) SetupWithManager(mgr ctrl.Manager, conf *api.ConfigVars) error {
	configFhirStore = conf
	logger.V(1).Info("Starting reconcile loop for fhirstore_controller.go")
	return ctrl.NewControllerManagedBy(mgr).
		For(&fhirv1alpha1.FhirStore{}).
		Complete(r)
}

// Add a finalizer to the kubernetes fhirstor's metadata
//   finalizers:
//  - fhir.imaware.com/finalizer
// This allows for full cleanup of the fhirstore and any other dependencies
func addFhirStoreFinalizer(fhirStore *v1alpha1.FhirStore, r *FhirStoreReconciler, ctx context.Context) error {
	if !controllerutil.ContainsFinalizer(fhirStore, FHISTORE_FINALIZER) {
		logger.V(1).Info(fmt.Sprintf("Adding finalizer to resource %v", fhirStore.Name))
		controllerutil.AddFinalizer(fhirStore, FHISTORE_FINALIZER)
		err := r.Update(ctx, fhirStore)
		if err != nil {
			logger.V(1).Error(err, fmt.Sprintf("Failed to add finalizer for fhirstore resource %v", fhirStore.Name))
			return err
		}
		logger.V(1).Info(fmt.Sprintf("Added finalizer to resource %v", fhirStore.Name))
	} else {
		logger.V(1).Info(fmt.Sprintf("Finalizer already present on resource %v", fhirStore.Name))
	}
	return nil
}

// Remove a finalizer to the kubernetes fhirstor's metadata
//   finalizers:
//  - fhir.imaware.com/finalizer
// This allows for full removal of the resource from the kubernetes ETCD
func removeFhirStoreFinalizer(fhirStore *v1alpha1.FhirStore, r *FhirStoreReconciler, ctx context.Context) error {
	logger.V(1).Info(fmt.Sprintf("Removing finalizer from resource %v", fhirStore.Name))
	controllerutil.RemoveFinalizer(fhirStore, FHISTORE_FINALIZER)
	err := r.Update(ctx, fhirStore)
	if err != nil {
		logger.V(1).Error(err, fmt.Sprintf("Failed to update finalizer for fhirstore resource %v", fhirStore.Name))
		return err
	}
	logger.V(1).Info(fmt.Sprintf("Removed finalizer from resource %v", fhirStore.Name))
	return nil
}
