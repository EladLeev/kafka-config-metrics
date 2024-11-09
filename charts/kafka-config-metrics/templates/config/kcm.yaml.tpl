{{- define "config/kcm.yaml.tpl" }}
global:
  port: ":{{ .Values.service.port }}"
  timeout: {{ .Values.kcmConfig.global.timeout }}

log:
  format: "{{ .Values.kcmConfig.log.format }}"
  level: "{{ .Values.kcmConfig.log.level }}"

kafka:
  defaultRefreshRate: {{ .Values.kcmConfig.kafka.defaultRefreshRate }}
  minKafkaVersion: "{{ .Values.kcmConfig.kafka.minKafkaVersion }}"
  adminTimeout: {{ .Values.kcmConfig.kafka.adminTimeout }}

clusters:
  {{- range $name, $cluster := .Values.kcmConfig.clusters }}
  {{ $name }}:
    brokers:
      {{- range $broker := $cluster.brokers }}
      - {{ $broker | quote }}
      {{- end }}
    topicFilter: {{ $cluster.topicFilter | quote }}
    {{- if $cluster.tls }}
    tls:
      enabled: {{ $cluster.tls.enabled }}
      {{- if $cluster.tls.caCert }}
      caCert: {{ $cluster.tls.caCert | quote }}
      {{- end }}
      {{- if $cluster.tls.clientCert }}
      clientCert: {{ $cluster.tls.clientCert | quote }}
      {{- end }}
      {{- if $cluster.tls.clientKey }}
      clientKey: {{ $cluster.tls.clientKey | quote }}
      {{- end }}
    {{- end }}
    {{- if $cluster.sasl }}
    sasl:
      enabled: {{ $cluster.sasl.enabled }}
      usernameEnv: {{ $cluster.sasl.usernameEnv | quote }}
      passwordEnv: {{ $cluster.sasl.passwordEnv | quote }}
      mechanism: {{ $cluster.sasl.mechanism | quote }}
    {{- end }}
  {{- end }}
{{- end }}