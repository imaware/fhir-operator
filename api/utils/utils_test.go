package utils

import (
	"fmt"
	"testing"

	"github.com/imaware/fhir-operator/api/v1alpha1"
)

func Test_get_fhir_json_map(t *testing.T) {
	var expected = "ActivityDefinition"
	var jsonString = `
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
	jsonMap, err := JsonStringToMap(jsonString)
	if err != nil {
		t.Errorf("Expected no error but returned error %v", err)
	}
	actual := jsonMap[RESOURCE_TYPE]
	if actual != expected {
		t.Errorf("Expected string %v but got %v", expected, actual)
	}

}

func Test_Attributes_to_map(t *testing.T) {
	notification := "notification"
	eventtype := "event"
	payloadformat := "json"
	bucketid := "id"
	objectid := "object"
	objectgeneration := "gen"
	eventtime := "time"
	attributes := make(map[string]string)
	expectedGcsEvent := GCSEvent{
		NotificationConfig: notification,
		EventType:          eventtype,
		PayloadFormat:      payloadformat,
		BucketId:           bucketid,
		ObjectId:           objectid,
		ObjectGeneration:   objectgeneration,
		EventTime:          eventtime,
	}
	attributes["notificationConfig"] = notification
	attributes["eventtype"] = eventtype
	attributes["payloadformat"] = payloadformat
	attributes["bucketid"] = bucketid
	attributes["objectid"] = objectid
	attributes["objectgeneration"] = objectgeneration
	attributes["eventtime"] = eventtime

	gcsEvent, err := AttributeMapToStruct(attributes)
	if err != nil {
		t.Errorf("Expected no error but returned error %v", err)
	}

	if *gcsEvent == expectedGcsEvent {
	} else {
		t.Errorf("Generated event %v struct is not equqal to expected %v", gcsEvent, expectedGcsEvent)
	}

}

func Test_generate_fhirResourceManifest(t *testing.T) {
	name := "test"
	namespace := "namespace"
	representation := `
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
	selector := "select"
	resourceType := "ActivityDefinition"
	fhirResource := GenerateFhirResourceManifest(selector, namespace, name, representation, resourceType)
	if fhirResource.Name != name {
		t.Errorf("Expected name %s but got name %s", name, fhirResource.Name)
	}
	if fhirResource.Namespace != namespace {
		t.Errorf("Expected namespace %s but got namespace %s", namespace, fhirResource.Namespace)
	}
	if fhirResource.Spec.Selector.Name != selector {
		t.Errorf("Expected selector.name %s but got selector.name %s", selector, fhirResource.Spec.Selector.Name)
	}
	if fhirResource.Spec.Representation != representation {
		t.Errorf("Expected representation %s but got representation %s", representation, fhirResource.Spec.Representation)
	}
	if fhirResource.Name != name {
		t.Errorf("Expected resource type %s but got resource type %s", resourceType, fhirResource.Spec.ResourceType)
	}
}

func Test_get_object_file_name_nestedfolders(t *testing.T) {
	expectedObject := "here"
	objectId := fmt.Sprintf("test/file/is/%s", expectedObject)
	actualObject := GetBucketObjectFileName(objectId)
	if actualObject != expectedObject {
		t.Errorf("Expected object name %s but got object name %s", expectedObject, actualObject)
	}
}

func Test_get_object_file_name_no_folders(t *testing.T) {
	expectedObject := "here"
	objectId := expectedObject
	actualObject := GetBucketObjectFileName(objectId)
	if actualObject != expectedObject {
		t.Errorf("Expected object name %s but got object name %s", expectedObject, actualObject)
	}
}

func Test_add_finalizer(t *testing.T) {
	resource := &v1alpha1.FhirResource{}
	AddFinalizer(resource, "test")
	if len(resource.ObjectMeta.Finalizers) == 0 {
		t.Errorf("Expected finalizer to be applied to object")
	}
}

func Test_remove_finalizer(t *testing.T) {
	resource := &v1alpha1.FhirResource{}
	finalizers := make([]string, 1)
	finalizers[0] = "test"
	resource.ObjectMeta.Finalizers = finalizers
	RemoveFinalizer(resource, "test")
	if len(resource.ObjectMeta.Finalizers) != 0 {
		t.Errorf("Expected finalizer to be applied to object")
	}
}
