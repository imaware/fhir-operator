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

	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

	"github.com/imaware/fhir-operator/api"
	"github.com/imaware/fhir-operator/api/v1alpha1"
	fhirv1alpha1 "github.com/imaware/fhir-operator/api/v1alpha1"
)

var (
	fhirResourceLogger = ctrl.Log.WithName("fhirresource_controller.go")
	configFhirResource *api.ConfigVars
	// 3 seconds
	pendingResourceDuration time.Duration = 3000000000
)

// FhirResourceReconciler reconciles a FhirResource object
type FhirResourceReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=fhir.imaware.com,resources=fhirresources,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=fhir.imaware.com,resources=fhirresources/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=fhir.imaware.com,resources=fhirresources/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the FhirResource object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.10.0/pkg/reconcile
func (r *FhirResourceReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	fhirResource := &v1alpha1.FhirResource{}
	err := r.Get(context.TODO(), req.NamespacedName, fhirResource)
	if err != nil {
		if errors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile request.
			// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
			// Return and don't requeue
			logger.V(1).Info("FhirResource resource not found.")
			return ctrl.Result{}, nil
		}
		// Error reading the object - requeue the request.
		logger.V(1).Error(err, "Failed to get FhirResource")
		return ctrl.Result{}, err
	}
	// make sure the fhirstore exists and is in CREATED state
	fhirStoreRequest := ctrl.Request{
		NamespacedName: types.NamespacedName{
			Name:      fhirResource.Spec.Selector.Name,
			Namespace: fhirResource.Namespace,
		},
	}
	//get fhirresource id
	fhirResourceID, err := api.GetFHIRIResourceID(fhirResource.Spec.Representation)
	if err != nil {
		logger.V(1).Error(err, "Failed to get resourceID", "resource", fhirResource.Name, "namespace", fhirResource.Namespace)
		fhirResource.Status.Status = api.PENDING
		fhirResource.Status.Message = api.FHIRStoreResourceCreateOrUpdateFailedStatus(fhirResource.Spec.Selector.Name, err.Error())
		r.Status().Update(ctx, fhirResource)
		return ctrl.Result{}, nil
	}
	fhirStore := getNamespacedFhirStore(fhirStoreRequest, r, ctx)
	// Check if the FhirStore instance is marked to be deleted, which is
	// indicated by the deletion timestamp being set.
	isFhirResourceMarkedToBeDeleted := fhirResource.GetDeletionTimestamp() != nil
	if isFhirResourceMarkedToBeDeleted {
		// If the fhirStore is empty or not ready there is no API calls needed to delete the fhir resource
		if fhirStore == nil || fhirStore.Status.Status != api.CREATED {
			err = removeFhirResourceFinalizer(fhirResource, r, ctx)
			if err != nil {
				return ctrl.Result{}, err
			}
			return ctrl.Result{}, err
		}
		// This means the fhirstore is in the CREATED state remove resource from store
		fhirResourceDeleteCall, err := api.BuildFHIRStoreResourceDeleteCall(configFhirResource.GCPProject, configFhirResource.GCPLocation, fhirStore.Spec.DatasetID, fhirStore.Spec.FhirStoreID, fhirResource.Spec.ResourceType, fhirResourceID)
		if err != nil {
			fhirResourceLogger.V(1).Error(err, "Something went wrong building delete call", "fhirResource", fhirResource.Name, "namesapce", fhirResource.Namespace)
			return ctrl.Result{}, err
		}
		fhirResourceGetCall, err := api.BuildFHIRStoreResourceGetCall(configFhirResource.GCPProject, configFhirResource.GCPLocation, fhirStore.Spec.DatasetID, fhirStore.Spec.FhirStoreID, fhirResource.Spec.ResourceType, fhirResourceID)
		if err != nil {
			fhirResourceLogger.V(1).Error(err, "Something went wrong building get call", "fhirResource", fhirResource.Name, "namesapce", fhirResource.Namespace)
			return ctrl.Result{}, err
		}
		err = api.DeleteAndReadFHIRStoreResource(fhirResourceDeleteCall, fhirResourceGetCall, fhirResource)
		r.Status().Update(ctx, fhirResource)
		if err != nil {
			fhirResourceLogger.V(1).Error(err, "Something went wrong during delete", "fhirResource", fhirResource.Name, "namesapce", fhirResource.Namespace)
			return ctrl.Result{}, nil
		}
		err = removeFhirResourceFinalizer(fhirResource, r, ctx)
		if err != nil {
			return ctrl.Result{}, err
		}
		return ctrl.Result{}, nil
	}
	err = addFhirResourceFinalizer(fhirResource, r, ctx)
	if err != nil {
		return ctrl.Result{}, err
	}
	// means the resource can not be added due to the store not being fully up
	if fhirStore == nil || fhirStore.Status.Status != api.CREATED {
		// requeue the event but wait 3 seconds
		fhirResource.Status.Status = api.PENDING
		fhirResource.Status.Message = api.FHIRStoreResourcePendingOnFhirStoreStatus(fhirResource.Spec.Selector.Name)
		r.Status().Update(ctx, fhirResource)
		return ctrl.Result{RequeueAfter: pendingResourceDuration}, nil
	}
	// fhir store is ready add the resource
	fhirResourceCreateOrUpdateCall, err := api.BuildFHIRStoreResourceUpdateCall(configFhirResource.GCPProject, configFhirResource.GCPLocation, fhirStore.Spec.DatasetID, fhirStore.Spec.FhirStoreID, fhirResource.Spec.Representation, fhirResource.Spec.ResourceType, fhirResourceID)
	if err != nil {
		fhirResourceLogger.V(1).Error(err, "Something went wrong building update call", "fhirResource", fhirResource.Name, "namesapce", fhirResource.Namespace)
		return ctrl.Result{}, nil
	}
	enqueu, err := api.CreateOrUpdateFHIRResource(fhirResourceCreateOrUpdateCall, fhirResource)
	r.Status().Update(ctx, fhirResource)
	if err != nil {
		fhirResourceLogger.Error(err, "Something went wrong during creation", "fhirResource", fhirResource.Name, "namesapce", fhirResource.Namespace)
		r.Status().Update(ctx, fhirResource)
		return ctrl.Result{}, nil
	} else if !enqueu {
		return ctrl.Result{}, nil
	} else {
		return ctrl.Result{RequeueAfter: pendingResourceDuration}, nil
	}
}

// SetupWithManager sets up the controller with the Manager.
func (r *FhirResourceReconciler) SetupWithManager(mgr ctrl.Manager, conf *api.ConfigVars) error {
	logger.V(1).Info("Starting reconcile loop for fhirresource_controller.go")
	configFhirResource = conf
	return ctrl.NewControllerManagedBy(mgr).
		For(&fhirv1alpha1.FhirResource{}).
		Complete(r)
}

// Add a finalizer to the kubernetes fhirresourc's metadata
//   finalizers:
//  - fhir.imaware.com/finalizer
// This allows for full cleanup of the fhirresource and any other dependencies
func addFhirResourceFinalizer(fhirResource *v1alpha1.FhirResource, r *FhirResourceReconciler, ctx context.Context) error {
	if !controllerutil.ContainsFinalizer(fhirResource, FHISTORE_FINALIZER) {
		fhirResourceLogger.V(1).Info(fmt.Sprintf("Adding finalizer to resource %v", fhirResource.Name))
		controllerutil.AddFinalizer(fhirResource, FHISTORE_FINALIZER)
		err := r.Update(ctx, fhirResource)
		if err != nil {
			fhirResourceLogger.V(1).Error(err, fmt.Sprintf("Failed to add finalizer for fhirstore resource %v", fhirResource.Name))
			return err
		}
		fhirResourceLogger.V(1).Info(fmt.Sprintf("Added finalizer to resource %v in namespace %v", fhirResource.Name, fhirResource.Namespace))
	} else {
		fhirResourceLogger.V(1).Info(fmt.Sprintf("Finalizer already present on resource %v in namespace %v", fhirResource.Name, fhirResource.Namespace))
	}
	return nil
}

// Remove a finalizer to the kubernetes fhirresourc's metadata
//   finalizers:
//  - fhir.imaware.com/finalizer
// This allows for full removal of the resource from the kubernetes ETCD
func removeFhirResourceFinalizer(fhirResource *v1alpha1.FhirResource, r *FhirResourceReconciler, ctx context.Context) error {
	fhirResourceLogger.V(1).Info(fmt.Sprintf("Removing finalizer from resource %v in namespace %v", fhirResource.Name, fhirResource.Namespace))
	controllerutil.RemoveFinalizer(fhirResource, FHISTORE_FINALIZER)
	err := r.Update(ctx, fhirResource)
	if err != nil {
		fhirResourceLogger.V(1).Error(err, fmt.Sprintf("Failed to update finalizer for fhirstore resource %v in namespace %v", fhirResource.Name, fhirResource.Namespace))
		return err
	}
	fhirResourceLogger.V(1).Info(fmt.Sprintf("Removed finalizer from resource %v in namespace %v", fhirResource.Name, fhirResource.Namespace))
	return nil
}

// Get the namesapce FhirStore from the kubernetes cluster
func getNamespacedFhirStore(fhirStoreRequest ctrl.Request, r *FhirResourceReconciler, ctx context.Context) *v1alpha1.FhirStore {
	fhirStore := &v1alpha1.FhirStore{}
	err := r.Get(context.TODO(), fhirStoreRequest.NamespacedName, fhirStore)
	if err != nil {
		if errors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile request.
			// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
			// Return and don't requeue
			fhirResourceLogger.V(1).Info(fmt.Sprintf("FhirStore %v resource in namespace %v not found", fhirStoreRequest.Name, fhirStoreRequest.Namespace))
			return nil
		}
		// Error reading the object - requeue the request.
		fhirResourceLogger.V(1).Error(err, "Failed to get FhirStore %v resource in namespace %v", fhirStore.Name, fhirStore.Namespace)
		return nil
	}
	fhirResourceLogger.V(1).Info(fmt.Sprintf("FhirStore %v resource in namespace %v found %v status", fhirStoreRequest.Name, fhirStoreRequest.Namespace, fhirStore.Status.Status))
	return fhirStore
}
