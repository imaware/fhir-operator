package api

import (
	"fmt"
	"os"
	"strconv"
)

const GCP_PROJECT = "GCP_PROJECT"
const GCP_LOCATION = "GCP_LOCATION"
const GOOGLE_APPLICATION_CREDENTIALS = "GOOGLE_APPLICATION_CREDENTIALS"
const DEBUG_ENABLED = "DEBUG_ENABLED"
const SENTRY_ENABLED = "SENTRY_ENABLED"
const ENVIRONMENT = "ENVIRONMENT"
const RELEASE_TAG = "RELEASE_TAG"
const SENTRY_DSN = "SENTRY_DSN"
const SENTRY_SAMPLE_RATE = "SENTRY_SAMPLE_RATE"

type ConfigVars struct {
	GCPProject       string
	GCPLocation      string
	DebugEnabled     bool
	SentryEnabled    bool
	SentryDSN        string
	Environment      string
	ReleaseTag       string
	SentrySampleRate float64
}

func GetConfig() (*ConfigVars, error) {
	gcpProject := os.Getenv(GCP_PROJECT)
	gcpLocation := os.Getenv(GCP_LOCATION)
	gcpCredentials := os.Getenv(GOOGLE_APPLICATION_CREDENTIALS)
	debugEnabled, err := strconv.ParseBool(os.Getenv(DEBUG_ENABLED))
	if err != nil {
		debugEnabled = false
	}
	sentrySampleRate, err := strconv.ParseFloat(os.Getenv(SENTRY_SAMPLE_RATE), 64)
	if err != nil {
		sentrySampleRate = 1.0
	}
	sentryDsn := os.Getenv(SENTRY_DSN)
	environment := os.Getenv(ENVIRONMENT)
	releaseTag := os.Getenv(RELEASE_TAG)
	sentryEnabled, err := strconv.ParseBool(os.Getenv(SENTRY_ENABLED))
	if err != nil {
		sentryEnabled = false
	}
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
		GCPProject:       gcpProject,
		GCPLocation:      gcpLocation,
		DebugEnabled:     debugEnabled,
		SentryEnabled:    sentryEnabled,
		SentryDSN:        sentryDsn,
		Environment:      environment,
		ReleaseTag:       releaseTag,
		SentrySampleRate: sentrySampleRate,
	}, nil
}
