package mocks

import "cloud.google.com/go/storage"

type MockGCSClientGood struct {
	StorageClient *storage.Client
}

func (m *MockGCSClientGood) DownLoadBucketObject(objectName string, bucketHandle *storage.BucketHandle) ([]byte, error) {
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
	fakeBytes := []byte(fhirRepresentation)
	return fakeBytes, nil
}

func (m *MockGCSClientGood) GetBucketHandle(bucketName string) *storage.BucketHandle {
	return &storage.BucketHandle{}
}
