module github.com/imaware/fhir-operator

go 1.16

require (
	cloud.google.com/go/pubsub v1.18.0
	cloud.google.com/go/storage v1.21.0
	github.com/googleapis/gax-go/v2 v2.1.1
	github.com/onsi/ginkgo v1.16.5
	github.com/onsi/gomega v1.19.0
	golang.org/x/tools v0.1.10 // indirect
	google.golang.org/api v0.69.0
	k8s.io/apimachinery v0.23.4
	k8s.io/client-go v0.23.4
	sigs.k8s.io/controller-runtime v0.11.1
)
