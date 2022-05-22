#
# Subchart/helper chart calls mainchart.app from this file.

{{- define "mainchart.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{- define "mainchart.app" -}}
{{template "mainchart.name" .}}
{{- end -}}
