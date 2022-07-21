module github.com/imaware/fhir-operator

go 1.16

require (
	cloud.google.com/go/pubsub v1.18.0
	cloud.google.com/go/storage v1.21.0
	github.com/getsentry/sentry-go v0.13.0
	github.com/googleapis/gax-go/v2 v2.1.1
	github.com/googleapis/gnostic v0.5.5 // indirect
	google.golang.org/api v0.69.0
	k8s.io/apimachinery v0.24.2
	k8s.io/client-go v0.24.2
	sigs.k8s.io/controller-runtime v0.12.3
)
