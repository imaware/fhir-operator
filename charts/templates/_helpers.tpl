{{/*
Expand the name of the chart.
*/}}
{{- define "fhir-operator.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Create a default fully qualified app name.
We truncate at 63 chars because some Kubernetes name fields are limited to this (by the DNS naming spec).
If release name contains chart name it will be used as a full name.
*/}}
{{- define "fhir-operator.fullname" -}}
{{- if .Values.fullnameOverride }}
{{- .Values.fullnameOverride | trunc 63 | trimSuffix "-" }}
{{- else }}
{{- $name := default .Chart.Name .Values.nameOverride }}
{{- if contains $name .Release.Name }}
{{- .Release.Name | trunc 63 | trimSuffix "-" }}
{{- else }}
{{- printf "%s-%s" .Release.Name $name | trunc 63 | trimSuffix "-" }}
{{- end }}
{{- end }}
{{- end }}

{{/*
Create chart name and version as used by the chart label.
*/}}
{{- define "fhir-operator.chart" -}}
{{- printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Common labels
*/}}
{{- define "fhir-operator.labels" -}}
helm.sh/chart: {{ include "fhir-operator.chart" . }}
{{ include "fhir-operator.selectorLabels" . }}
{{- if .Chart.AppVersion }}
app.kubernetes.io/version: {{ .Chart.AppVersion | quote }}
{{- end }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
{{- end }}

{{/*
Selector labels
*/}}
{{- define "fhir-operator.selectorLabels" -}}
app.kubernetes.io/name: {{ include "fhir-operator.name" . }}
app.kubernetes.io/app: {{ include "fhir-operator.name" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end }}

{{/*
Create the name of the service account to use
*/}}
{{- define "fhir-operator.serviceAccountName" -}}
{{- if .Values.serviceAccount.create }}
{{- default (include "fhir-operator.fullname" .) .Values.serviceAccount.name }}
{{- else }}
{{- default "default" .Values.serviceAccount.name }}
{{- end }}
{{- end }}

{{- define "fhir-operator.iamAccountName" -}}
{{ include "fhir-operator.serviceAccountName" . | trunc 25 | trimSuffix "-" }}
{{- end }}

{{- define "fhir-operator.gcpAccountName" -}}
{{ include "fhir-operator.iamAccountName" . }}@{{ $.Values.iamServiceAccount.project }}.iam.gserviceaccount.com
{{- end }}
