{{ define "config/customconfig.toml.tpl" }}
# Helm Generated Config

[global]
port = ":{{ .Values.service.port }}"

[log]
level = "{{ .Values.kcmConfig.log_level }}"

[kafka]
default_refresh_rate = {{ .Values.kcmConfig.default_refresh_rate }}
min_kafka_version = "{{ .Values.kcmConfig.min_kafka_version }}"
admin_timeout = {{ .Values.kcmConfig.admin_timeout }}

[clusters]
{{- range $k, $v := .Values.kcmConfig.clusters }}
  [clusters.{{ $k }}]
  brokers = ["{{ $v.brokers }}"]
  topicfilter="{{ $v.topicfilter }}"
{{ end }}
{{ end }}