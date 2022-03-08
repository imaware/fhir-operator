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

	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/predicate"

	"github.com/imaware/fhir-operator/api"
	"github.com/imaware/fhir-operator/api/subscriber"
	"github.com/imaware/fhir-operator/api/utils"
	"github.com/imaware/fhir-operator/api/v1alpha1"
	fhirv1alpha1 "github.com/imaware/fhir-operator/api/v1alpha1"
)

var (
	fhirGCSconnectorLogger = ctrl.Log.WithName("fhirgcsconnector_controller.go")
)

// FhirGCSConnectorReconciler reconciles a FhirGCSConnector object
type FhirGCSConnectorReconciler struct {
	client.Client
	Scheme        *runtime.Scheme
	PubSubClient  *api.GCPPUBClient
	StorageClient *api.GCSClient
}

//+kubebuilder:rbac:groups=fhir.imaware.com,resources=fhirgcsconnectors,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=fhir.imaware.com,resources=fhirgcsconnectors/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=fhir.imaware.com,resources=fhirgcsconnectors/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the FhirGCSConnector object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.10.0/pkg/reconcile
func (r *FhirGCSConnectorReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	var result = ctrl.Result{}
	var err error = nil
	// Get all valuable metadata for object
	fhirGCSConnector := &fhirv1alpha1.FhirGCSConnector{}
	err = r.Get(context.TODO(), req.NamespacedName, fhirGCSConnector)
	if err != nil {
		if errors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile request.
			// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
			// Return and don't requeue
			fhirGCSconnectorLogger.V(1).Info("FhirGCSConnector resource not found.")
			// set err to nil to not requeue the request
			err = nil
		} else {
			fhirGCSconnectorLogger.V(1).Error(err, "Failed to get FhirGCSConnector")
		}
	} else {
		// delete the resource and pub sub subscription
		isFhirGCSMarkedForDeletion := fhirGCSConnector.GetDeletionTimestamp() != nil
		if isFhirGCSMarkedForDeletion {
			err = r.PubSubClient.DeleteSubscription(r.PubSubClient.GetSubscription(fhirGCSConnector.Spec.SubscriptionName))
			if err != nil {
				fhirGCSconnectorLogger.Info(fmt.Sprintf("Failed to delete subscription %s due to: [ %v ]", fhirGCSConnector.Spec.SubscriptionName, err))
				err = nil
				utils.RemoveFinalizer(fhirGCSConnector, FHISTORE_FINALIZER)
				err := r.Update(ctx, fhirGCSConnector)
				if err != nil {
					logger.V(1).Error(err, fmt.Sprintf("Failed to remove finalizer for fhirGCSConnector resource %v", fhirGCSConnector.Name))
				}
			}
		} else {
			// start consuming from topic
			utils.AddFinalizer(fhirGCSConnector, FHISTORE_FINALIZER)
			err := r.Update(ctx, fhirGCSConnector)
			if err == nil {
				r.initConsumer(fhirGCSConnector)
			}
		}
	}
	// update status reflecting resource
	r.Status().Update(ctx, fhirGCSConnector)
	return result, err
}

// SetupWithManager sets up the controller with the Manager.
func (r *FhirGCSConnectorReconciler) SetupWithManager(mgr ctrl.Manager) error {
	pred := predicate.GenerationChangedPredicate{}
	return ctrl.NewControllerManagedBy(mgr).
		For(&fhirv1alpha1.FhirGCSConnector{}).WithEventFilter(pred).
		Complete(r)
}

// Initialize the consumer and or start the goroutine to consume from the topic
func (r *FhirGCSConnectorReconciler) initConsumer(fhirGCSConnector *v1alpha1.FhirGCSConnector) {
	topic, err := r.PubSubClient.GetTopic(fhirGCSConnector.Spec.Topic)
	if err != nil {
		fhirGCSConnector.Status.Status = api.FHIRGCSConnectorFailedToGetTopicStatus(fhirGCSConnector.Spec.SubscriptionName, err.Error())
		return
	}
	subscription := r.PubSubClient.GetSubscription(fhirGCSConnector.Spec.SubscriptionName)
	exists, err := subscription.Exists(context.Background())
	if err != nil {
		fhirGCSConnector.Status.Status = api.FHIRGCSConnectorFailedSubscriptionCreationStatus(fhirGCSConnector.Spec.SubscriptionName, err.Error())
		return
	}
	if !exists {
		err = r.PubSubClient.CreateSubscription(topic, fhirGCSConnector.Spec.Filter, fhirGCSConnector.Spec.SubscriptionName)
		if err != nil {
			fhirGCSConnector.Status.Status = api.FHIRGCSConnectorFailedSubscriptionCreationStatus(fhirGCSConnector.Spec.SubscriptionName, err.Error())
			return
		}
	}
	// TODO Implement update logic if needed
	consumerConfig := subscriber.ConsumerConfig{
		Namespace:         fhirGCSConnector.Namespace,
		SubscriptionName:  fhirGCSConnector.Spec.SubscriptionName,
		TopicID:           fhirGCSConnector.Spec.Topic,
		FhirStoreSelector: fhirGCSConnector.Spec.FhirStoreSelector.Name,
		GCSClient:         r.StorageClient,
		Subscription:      subscription,
		K8sClient:         r.Client,
	}
	go subscriber.Consume(&consumerConfig)
	fhirGCSConnector.Status.Status = api.FHIRGCSConnectorConsumingFromTopicStatus(fhirGCSConnector.Spec.SubscriptionName, fhirGCSConnector.Spec.Topic)
	fhirGCSconnectorLogger.Info(fmt.Sprintf("Started consumer for resource %s in namesapce %s", fhirGCSConnector.Name, consumerConfig.Namespace))
}
