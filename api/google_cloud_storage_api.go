package api

import (
	"context"
	"io/ioutil"
	"time"

	"cloud.google.com/go/storage"
)

type GCSClientCalls interface {
	DownLoadBucketObject(objectName string, bucketHandle *storage.BucketHandle) ([]byte, error)
	GetBucketHandle(bucketName string) *storage.BucketHandle
}

type GCSClient struct {
	StorageClient *storage.Client
}

// Download a bucket object from GCS
//
// bucketHandler the GCS handler for interacting with storage buckets
//
// objectName name of object to download
//
// returns a list of bytes based on whats downloaded or an error if
// downloading fails
func (c *GCSClient) DownLoadBucketObject(objectName string, bucketHandle *storage.BucketHandle) ([]byte, error) {
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, time.Second*50)
	defer cancel()
	rc, err := bucketHandle.Object(objectName).NewReader(ctx)
	if err != nil {
		return nil, err
	}
	defer rc.Close()
	fhirResourceData, err := ioutil.ReadAll(rc)
	if err != nil {
		return nil, err
	}
	return fhirResourceData, nil
}

// Returns a GCS storage bucket handler to interact with gcp
func (c *GCSClient) GetBucketHandle(bucketName string) *storage.BucketHandle {
	return c.StorageClient.Bucket(bucketName)
}
