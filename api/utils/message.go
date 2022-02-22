package utils

const (
	EVENT_TYPE_FINALIZE = "OBJECT_FINALIZE"
	EVENT_TYPE_DELETE   = "OBJECT_DELETE"
)

type GCSEvent struct {
	NotificationConfig string
	EventType          string
	PayloadFormat      string
	BucketId           string
	ObjectId           string
	ObjectGeneration   string
	EventTime          string
}
