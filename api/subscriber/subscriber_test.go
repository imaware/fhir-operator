package subscriber

import (
	"testing"

	"cloud.google.com/go/pubsub"
	"github.com/imaware/fhir-operator/api/utils"
	"github.com/imaware/fhir-operator/mocks"
)

func Test_create_fhir_resource(t *testing.T) {
	gcsEvent := generateFakeGCSEvent()
	fakeSubscription := &pubsub.Subscription{}
	consumerConfig := ConsumerConfig{
		Namespace:         "test",
		SubscriptionName:  "for",
		TopicID:           "sf_updates",
		FhirStoreSelector: "test",
		GCSClient:         &mocks.MockGCSClientGood{StorageClient: nil},
		Subscription:      fakeSubscription,
		K8sClient:         nil}
	fhirRepresentation := `
	{
		"resourceType": "ActivityDefinition",
		"id": "0062fbbc-e09b-4819-80b4-9c02724d1f11",
		"title": "Ferritinnn",
		"status": "active",
		"observationResultRequirement": [
		  {
			"reference": "ObservationDefinition/788075ad-1f44-4b72-8642-6a2af6c8b1da",
			"display": "Ferritin"
		  }
		]
	  }
	  `
	_, err := buildFhirResource(&consumerConfig, gcsEvent, fhirRepresentation)
	if err != nil {
		t.Errorf("expected no error but got error %v", err)
	}

}

func Test_process_event_resource_does_not_exist(t *testing.T) {
	pubsubMessage := generateFakePubsubMessageFinalize()
	fakeSubscription := &pubsub.Subscription{}
	consumerConfig := ConsumerConfig{
		Namespace:         "test",
		SubscriptionName:  "for",
		TopicID:           "sf_updates",
		FhirStoreSelector: "test",
		GCSClient:         &mocks.MockGCSClientGood{StorageClient: nil},
		Subscription:      fakeSubscription,
		K8sClient:         &mocks.MockK8sClientCreatePassed{}}
	processEvent(&consumerConfig, pubsubMessage)
}

func Test_process_event_resource_does_exist(t *testing.T) {
	pubsubMessage := generateFakePubsubMessageFinalize()
	fakeSubscription := &pubsub.Subscription{}
	consumerConfig := ConsumerConfig{
		Namespace:         "test",
		SubscriptionName:  "for",
		TopicID:           "sf_updates",
		FhirStoreSelector: "test",
		GCSClient:         &mocks.MockGCSClientGood{StorageClient: nil},
		Subscription:      fakeSubscription,
		K8sClient:         &mocks.MockK8sClientCreateAlreadyExists{}}
	processEvent(&consumerConfig, pubsubMessage)
}

func Test_process_event_delete_resource(t *testing.T) {
	pubsubMessage := generateFakePubsubMessageDelete()
	fakeSubscription := &pubsub.Subscription{}
	consumerConfig := ConsumerConfig{
		Namespace:         "test",
		SubscriptionName:  "for",
		TopicID:           "sf_updates",
		FhirStoreSelector: "test",
		GCSClient:         &mocks.MockGCSClientGood{StorageClient: nil},
		Subscription:      fakeSubscription,
		K8sClient:         &mocks.MockK8sClientCreatePassed{}}
	processEvent(&consumerConfig, pubsubMessage)
}

func generateFakePubsubMessageFinalize() *pubsub.Message {
	notification := "notification"
	eventtype := utils.EVENT_TYPE_FINALIZE
	payloadformat := "json"
	bucketid := "id"
	objectid := "object"
	objectgeneration := "gen"
	eventtime := "time"
	attributes := make(map[string]string)
	attributes["notificationConfig"] = notification
	attributes["eventtype"] = eventtype
	attributes["payloadformat"] = payloadformat
	attributes["bucketid"] = bucketid
	attributes["objectid"] = objectid
	attributes["objectgeneration"] = objectgeneration
	attributes["eventtime"] = eventtime
	pubsubMessage := &pubsub.Message{
		Attributes: attributes,
	}
	return pubsubMessage
}

func generateFakePubsubMessageDelete() *pubsub.Message {
	notification := "notification"
	eventtype := utils.EVENT_TYPE_DELETE
	payloadformat := "json"
	bucketid := "id"
	objectid := "object"
	objectgeneration := "gen"
	eventtime := "time"
	attributes := make(map[string]string)
	attributes["notificationConfig"] = notification
	attributes["eventtype"] = eventtype
	attributes["payloadformat"] = payloadformat
	attributes["bucketid"] = bucketid
	attributes["objectid"] = objectid
	attributes["objectgeneration"] = objectgeneration
	attributes["eventtime"] = eventtime
	pubsubMessage := &pubsub.Message{
		Attributes: attributes,
	}
	return pubsubMessage
}

func generateFakeGCSEvent() *utils.GCSEvent {
	notification := "notification"
	eventtype := "event"
	payloadformat := "json"
	bucketid := "id"
	objectid := "object"
	objectgeneration := "gen"
	eventtime := "time"
	gcsEvent := &utils.GCSEvent{
		NotificationConfig: notification,
		EventType:          eventtype,
		PayloadFormat:      payloadformat,
		BucketId:           bucketid,
		ObjectId:           objectid,
		ObjectGeneration:   objectgeneration,
		EventTime:          eventtime,
	}
	return gcsEvent

}
