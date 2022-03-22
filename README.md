# fhir-operator
Kubernetes operator to handle FHIR databases and resources.  

The Operator works off of FhirStore and FhirResource CRD. These resources give the power to create GCP FhirStores and FhirResources.

# Getting up and running

Configure the [values.yaml](charts/values.yaml) , make sure you have access to the kubernetes cluster, and run...

```
make deploy
```
Apply the example basic fhir store
```
kubectl apply -f examples/stores/fhirstore_basic.yaml
```
Describe resource and make sure it states CREATED in status  
Apply the example fhir resources
```
kubectl apply -f examples/resources
```
Describe the resources and make sure they state CREATED in status

## Resources
### *FhirStore
Resource used to create a GCP FhirStore within a dataset. Currently requires the dataset to exist. FhirStore's are namespace scoped and provide the ability to specify the name of the store and any IAM permission policies to apply to the API. Reference examples/stores for examples. If the fhir store already exists, it will simply bind the resource definition of the existing store.

- FhirStore.Spec.auth: each key represents a role to bind to a list of members
- FhirStore.Spec.options:
  - preventDelete: prevent the delete of the resource and fhirStore in case of resrouce deletion
### *FhirResource
Resource is used to create a Fhir resource in the FhirStore specified in the selector. *NOTE*: The selector points to the actual FhirStore resource and not the GCP Fhir store. A FhirResource can accomodate any FHIR representation. Reference examples/resources for examples

### *FhirGCSConnector
Resource is used to create a subscription to a pubsub topic that listenes for GCS bucket events that holds json data that represents a fhir resource. The name of the file stored in the bucket must be DNS compliant for kubernetes as each file will create a fhir resource that is bound to the same namespace as the FhirGCSConnector resource. Reference examples/fhirgcsconnector for examples. You must have a kubernetes fhirstore resource to use this

- FhirGCSConnector.Spec.topic: name of the topic to bind the subscription to
- FhirGCSConnector.Spec.subscriptionName: name of the subscription to create
- FhirGCSConnector.Spec.filter(optional): subscription filter that can be applied to the resource NOTE: this filter can not be changed once the resource is create. Must delete and re-create to update filter
- FhirGCSConnector.Spec.fhirStoreSelector:
  - name: name of the fhirstore to bind the fhir resoures too
## Testing
Make sure you have a GCP service account confgiured and pointing to GOOGLE_APPLICATION_CREDENTIALS envar
```
make test
```
