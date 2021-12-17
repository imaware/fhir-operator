package utils

import "testing"

const RESOURCE_TYPE = "resourceType"

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
