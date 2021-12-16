package api

import (
	"fmt"
	"os"
)

const GCP_PROJECT = "GCP_PROJECT"
const GCP_LOCATION = "GCP_LOCATION"
const GOOGLE_APPLICATION_CREDENTIALS = "GOOGLE_APPLICATION_CREDENTIALS"

type ConfigVars struct {
	GCPProject  string
	GCPLocation string
}

func GetConfig() (*ConfigVars, error) {
	gcpProject := os.Getenv(GCP_PROJECT)
	gcpLocation := os.Getenv(GCP_LOCATION)
	gcpCredentials := os.Getenv(GOOGLE_APPLICATION_CREDENTIALS)
	if len(gcpProject) == 0 {
		return nil, fmt.Errorf("environment variable not set: %v", GCP_PROJECT)
	}
	if len(gcpLocation) == 0 {
		return nil, fmt.Errorf("environment variable not set: %v", GCP_LOCATION)
	}
	if len(gcpCredentials) == 0 {
		return nil, fmt.Errorf("environment variable not set: %v", GOOGLE_APPLICATION_CREDENTIALS)
	}
	return &ConfigVars{
		GCPProject:  gcpProject,
		GCPLocation: gcpLocation,
	}, nil
}
