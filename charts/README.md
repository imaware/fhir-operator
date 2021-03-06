# fhir-operator

![Version: 0.1.0](https://img.shields.io/badge/Version-0.1.0-informational?style=flat-square) ![Type: application](https://img.shields.io/badge/Type-application-informational?style=flat-square) ![AppVersion: 1.16.0](https://img.shields.io/badge/AppVersion-1.16.0-informational?style=flat-square)

A Helm chart for Kubernetes

## Values

| Key | Type | Default | Description |
|-----|------|---------|-------------|
| affinity | object | `{}` | Requirements for pods to run on node |
| commandFlags | list | `[]` | Additional command arguments that can be passed to controller |
| env.DEBUG_ENABLED | string | `"false"` | env var for enabling debug logging |
| env.ENVIRONMENT | string | `""` | env var for specifying the environment the fhir operator is deployed in |
| env.GCP_LOCATION | string | `""` | (REQUIRED) env var that points to the GCP location |
| env.GCP_PROJECT | string | `""` | (REQUIRED) env var that points to the project in GCP |
| env.SENTRY_DSN | string | `""` | env var for setting sentry DSN for traces |
| env.SENTRY_ENABLED | string | `"fasle"` | env var for enabling sentry |
| env.SENTRY_SAMPLE_RATE | string | `"1.0"` | env var to specify the sample rate for sentry |
| fullnameOverride | string | `""` |  |
| iamServiceAccount.create | bool | `true` | If running in a cluster with Anthos config connector will create GCP IAM resources |
| iamServiceAccount.polices[0] | object | `{"project":"","role":"roles/healthcare.fhirResourceEditor"}` | (REQUIRED) role for the fhir operator to edit fhir resources in the project. Project must be set to the GCP project in use. |
| iamServiceAccount.polices[1] | object | `{"project":"","role":"roles/healthcare.fhirStoreAdmin"}` | (REQUIRED) role for the fhir operator to edit manage fhir stores in the project. Project must be set to the GCP project in use. |
| iamServiceAccount.polices[2] | object | `{"role":"roles/iam.workloadIdentityUser"}` | (REQUIRED) role for the fhir operator use workload identity |
| iamServiceAccount.polices[3].project | string | `""` |  |
| iamServiceAccount.polices[3].role | string | `"roles/storage.admin"` |  |
| iamServiceAccount.polices[4].project | string | `""` |  |
| iamServiceAccount.polices[4].role | string | `"roles/pubsub.admin"` |  |
| iamServiceAccount.project | string | `""` | (REQUIRED) GCP project |
| image.pullPolicy | string | `"Always"` | ImagePullPolicy settings |
| image.repository | string | `"gcr.io/imaware-test/fhir-operator"` | Repository hosting the image |
| image.tag | string | `"1.0.0"` | Overrides the image tag whose default is the chart appVersion. |
| imagePullSecrets | list | `[]` | Image pull secrets for the container registry |
| nameOverride | string | `""` |  |
| namespace | string | `""` | What namesapce to deploy namespaced resources too |
| nodeSelector | object | `{}` | Which node pods should run on |
| podAnnotations | object | `{}` | Additional pod annotations |
| podSecurityContext.runAsNonRoot | bool | `true` | Allow pods to run as root |
| pspCreate | bool | `true` | To create the PSP or not |
| replicaCount | int | `1` | Number of replicas |
| resources.limits.cpu | string | `"250m"` |  |
| resources.limits.memory | string | `"512Mi"` |  |
| resources.requests.cpu | string | `"100m"` |  |
| resources.requests.memory | string | `"256Mi"` |  |
| securityContext | object | `{}` | Additional container secuirty context settings |
| serviceAccount.annotations | object | `{}` |  |
| serviceAccount.create | bool | `true` | Specifies whether a service account should be created |
| serviceAccount.name | string | `"fhir-operator"` | The name of the service account to use. If not set and create is true, a name is generated using the fullname template |
| tolerations | list | `[]` | Which nodes pods should tolerate |

----------------------------------------------
Autogenerated from chart metadata using [helm-docs v1.5.0](https://github.com/norwoodj/helm-docs/releases/v1.5.0)
