package subscriber

import (
	"context"
	"fmt"

	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	ctrl "sigs.k8s.io/controller-runtime"

	"cloud.google.com/go/pubsub"
	"github.com/getsentry/sentry-go"
	"github.com/imaware/fhir-operator/api"
	"github.com/imaware/fhir-operator/api/utils"
	"github.com/imaware/fhir-operator/api/v1alpha1"
)

var (
	consumerLogger = ctrl.Log.WithName("consumer.go")
)

type ConsumerConfig struct {
	Namespace         string
	SubscriptionName  string
	TopicID           string
	FhirStoreSelector string
	GCSClient         api.GCSClientCalls
	Subscription      *pubsub.Subscription
	K8sClient         api.K8sClient
}

func Consume(consumerConfig *ConsumerConfig) {
	ctx := context.Background()
	// make it so we can stop the goroutine if an error happens
	// start the consume loop
	err := consumerConfig.Subscription.Receive(ctx, func(ctx context.Context, msg *pubsub.Message) {
		processEvent(consumerConfig, msg)
		msg.Ack()
	})
	if err != nil {
		consumerLogger.Info(fmt.Sprintf("Failed to subscribe to topic %s for subscription %s due to [ %v ]", consumerConfig.TopicID, consumerConfig.SubscriptionName, err))
	}
	consumerLogger.Info(fmt.Sprintf("Stopping thread for subscription %s", consumerConfig.SubscriptionName))
}

// Core logic to process a GCS pubsub event
//
// consumerConfig configuration to process event
//
// msg pubsub event from topic
func processEvent(consumerConfig *ConsumerConfig, msg *pubsub.Message) {
	gcsEvent, err := utils.AttributeMapToStruct(msg.Attributes)
	if err != nil {
		consumerLogger.Error(err, fmt.Sprintf("Failed serialize message to GCSEvent %v for subscription %s.", msg.Attributes, consumerConfig.SubscriptionName))
	}
	// object has been created or updated
	if gcsEvent.EventType == utils.EVENT_TYPE_FINALIZE {
		downloadedBytes, err := consumerConfig.GCSClient.DownLoadBucketObject(gcsEvent.ObjectId, consumerConfig.GCSClient.GetBucketHandle(gcsEvent.BucketId))
		// can be due to a race condition where the object is actually deleted prior to consuming the event
		if err != nil {
			consumerLogger.Error(err, fmt.Sprintf("Failed to download fhir resource from topic %s subscription %s", consumerConfig.TopicID, consumerConfig.SubscriptionName))
		} else {
			if len(downloadedBytes) > 0 {
				fhirResource, err := buildFhirResource(consumerConfig, gcsEvent, string(downloadedBytes))
				if err != nil {
					sentry.CaptureException(err)
					consumerLogger.Error(err, fmt.Sprintf("Failed to build fhir resource for subscription %s", consumerConfig.SubscriptionName))
				} else {
					// object has been created or updated
					// create the fhir resource manifest
					err = consumerConfig.K8sClient.Create(context.TODO(), fhirResource)
					if err != nil {
						// perform update instead
						if errors.IsAlreadyExists(err) {
							err = consumerConfig.K8sClient.Update(context.TODO(), fhirResource)
							if err != nil {
								sentry.CaptureException(err)
								consumerLogger.Error(err, fmt.Sprintf("Failed to update fhir resource %s for subscription %s in namespace %s", fhirResource.Name, consumerConfig.SubscriptionName, fhirResource.Namespace))
							}
						} else {
							sentry.CaptureException(err)
							consumerLogger.Error(err, fmt.Sprintf("Failed to create fhir resource %s for subscription %s in namesapce %s", fhirResource.Name, consumerConfig.SubscriptionName, fhirResource.Namespace))
						}
					}
				}
			}
		}
		// Event is to delete a resource
	} else if gcsEvent.EventType == utils.EVENT_TYPE_DELETE {
		// name of object in bucket correlates to name of resource
		resourceName := utils.GetBucketObjectFileName(gcsEvent.ObjectId)
		fhirResourceToDelete := &v1alpha1.FhirResource{
			ObjectMeta: metav1.ObjectMeta{
				Name:      resourceName,
				Namespace: consumerConfig.Namespace,
			},
		}
		err = consumerConfig.K8sClient.Delete(context.TODO(), fhirResourceToDelete)
		if err != nil {
			if !errors.IsNotFound(err) {
				consumerLogger.Info(fmt.Sprintf("Fhir resource %s in namespace %s for subscription %s does not exist", fhirResourceToDelete.Name, fhirResourceToDelete.Namespace, consumerConfig.SubscriptionName))
			}
		}
	}

}

// Build the fhir resource based on what was received from the pub sub event
//
// consumerConfig everything required to interact with the cloud
//
// gcsEvent the google cloud storage event metadata
//
// return a fhir resource manifest or an error
func buildFhirResource(consumerConfig *ConsumerConfig, gcsEvent *utils.GCSEvent, fhirRepresentation string) (*v1alpha1.FhirResource, error) {
	fhirResourceRepresentationMap, err := utils.JsonStringToMap(fhirRepresentation)
	if err != nil {
		return nil, fmt.Errorf("failed to convert representation: [ %s ] to a map", fhirRepresentation)
	}
	// name of object in bucket correlates to name of resource
	resourceName := utils.GetBucketObjectFileName(gcsEvent.ObjectId)
	// sprintf here is janky TODO something better
	fhirResource := utils.GenerateFhirResourceManifest(consumerConfig.FhirStoreSelector, consumerConfig.Namespace, resourceName, fhirRepresentation, fmt.Sprintf("%v", fhirResourceRepresentationMap[utils.RESOURCE_TYPE]))
	return fhirResource, nil

}
