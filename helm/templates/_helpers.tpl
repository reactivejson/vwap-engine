{{/* vim: set filetype=mustache: */}}
{{/*
Expand the name of the chart.
*/}}
{{- define "vwap-engine.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{/*
Create a default fully qualified app name.
We truncate at 63 chars because some Kubernetes name fields are limited to this (by the DNS naming spec).
*/}}
{{- define "vwap-enginefullname" -}}
{{- $name := default .Chart.Name .Values.nameOverride -}}
{{- printf "%s-%s" .Release.Name $name | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{- /*
vwap-engine.chartref prints a chart name and version.
It does minimal escaping for use in Kubernetes labels.
*/ -}}
{{- define "vwap-engine.chartref" -}}
  {{- replace "+" "_" .Chart.Version | printf "%s-%s" .Chart.Name -}}
{{- end -}}

{{/*
vwap-engine.labels.standard prints the standard Helm labels.
The standard labels are frequently used in metadata.
*/}}
{{- define "vwap-engine.labels.standard" -}}
app: {{template "vwap-engine.name" .}}
chart: {{template "vwap-engine.chartref" . }}
app.kubernetes.io/name: {{template "vwap-engine.name" .}}
helm.sh/chart: {{template "vwap-engine.chartref" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
app.kubernetes.io/version: {{ .Chart.Version }}
com.nokia.neo/commitId: ${COMMIT_ID}
{{- end -}}

{{/*
vwap-engine.template.labels prints the template metadata labels.
*/}}
{{- define "vwap-engine.template.labels" -}}
app: {{template "vwap-engine.name" .}}
{{- end -}}

{{- define "vwap-engine.app" -}}
app: {{template "vwap-engine.name" .}}
{{- end -}}

{{- define "annotateResources" -}}
# Preserve the workingsetpluginregistrations.ws.nokia.com crd for changes;
kubectl annotate --overwrite crd workingsetpluginregistrations.ws.nokia.com "helm.sh/resource-policy"=keep;
sleep 30;
{{- end -}}
