{{- define "injector.NamespaceSelector" }}
{{- if .Values.injectionNamespaces }}
namespaceSelector:
  matchExpressions:
    - key: kubernetes.io/metadata.name
      operator: In
      values:
{{ toYaml .Values.injectionNamespaces | indent 8 }}
{{- end }}
{{- end }}
