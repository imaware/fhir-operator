package utils

import (
	"encoding/json"
	"strings"

	"github.com/imaware/fhir-operator/api/v1alpha1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

const RESOURCE_TYPE = "resourceType"

func JsonStringToMap(jsonString string) (map[string]interface{}, error) {
	// Declared an empty map interface
	var result map[string]interface{}
	// Unmarshal or Decode the JSON to the interface.
	err := json.Unmarshal([]byte(jsonString), &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func AttributeMapToStruct(attributes map[string]string) (*GCSEvent, error) {
	jsonbody, err := json.Marshal(attributes)
	if err != nil {
		return nil, err
	}

	gcsEvent := &GCSEvent{}
	if err := json.Unmarshal(jsonbody, &gcsEvent); err != nil {
		return nil, err
	}
	return gcsEvent, nil
}

// Generate a fhir resource object
func GenerateFhirResourceManifest(fhirStoreSelector string, namespace string, name string, representation string, resourceType string) *v1alpha1.FhirResource {
	fhirResource := &v1alpha1.FhirResource{}
	fhirResourceSelector := v1alpha1.FhirResourceSpecFhirStoreSelector{Name: fhirStoreSelector}
	fhirResource.Name = strings.ToLower(name)
	fhirResource.Namespace = namespace
	fhirResource.Spec.Selector = fhirResourceSelector
	fhirResource.Spec.Representation = representation
	fhirResource.Spec.ResourceType = resourceType
	return fhirResource
}

// Get the name of the object stored in GCS we split by the folder delimeter / to get the actual object name
// returns the name of the object
func GetBucketObjectFileName(objectId string) string {
	delimeter := "/"
	resourceName := objectId
	if strings.Contains(objectId, delimeter) {
		splits := strings.Split(objectId, delimeter)
		resourceName = splits[len(splits)-1]
	}
	return resourceName
}

// Adds a finalizer string to the objects metadata
func AddFinalizer(o client.Object, finalizer string) {
	if !controllerutil.ContainsFinalizer(o, finalizer) {
		controllerutil.AddFinalizer(o, finalizer)
	}
}

// Remove finalizer from ojbect metadata
func RemoveFinalizer(o client.Object, finalizer string) {
	if controllerutil.ContainsFinalizer(o, finalizer) {
		controllerutil.RemoveFinalizer(o, finalizer)
	}
}
