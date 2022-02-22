package templates

// AppOverviewTemplate is the markdown template used to
// display Application overview.
const AppOverviewTemplate string = `# Application: {{ .Name }}

## Overview

| Setting | Value | 
|---|---|
|Git Repository  | ` + "`{{ .Status.Sync.ComparedTo.Source.RepoURL }}/{{ .Status.Sync.ComparedTo.Source.Path }}`" + ` | 
| Source Type | {{ .Status.SourceType }}|
|Dest Namespace  | {{ .Status.Sync.ComparedTo.Destination.Namespace }} | 
|Sync Revision   | {{ .Status.Sync.Revision }} ({{ .Status.Sync.Status }} & {{ .Status.Health.Status }}) |
|Target Revision | 
{{- if .Status.Sync.ComparedTo.Source.TargetRevision }}{{ .Status.Sync.ComparedTo.Source.TargetRevision }} |
{{- else -}}
HEAD |
{{- end }}
|Automated Prune | {{- if .Spec.SyncPolicy.Automated }}{{ .Spec.SyncPolicy.Automated.Prune }}{{ else }}✗{{ end -}} | 
|Automated Heal | {{- if .Spec.SyncPolicy.Automated }}{{ .Spec.SyncPolicy.Automated.SelfHeal}}{{ else }}✗{{ end -}} | 
{{- if .Spec.SyncPolicy.SyncOptions }}
{{- range $option := .Spec.SyncPolicy.SyncOptions }}
| Sync Option | {{ $option }} | 
{{- end }}
{{- else }}
| Sync Option | ✗ |
{{- end }}

{{- if .Status.Summary.ExternalURLs }}
## External Urls

{{- range $url := .Status.Summary.ExternalURLs }}
* {{ $url }}
{{- end }}

{{- end }}

{{- if .Status.Summary.Images }}
## Images

{{- range $image := .Status.Summary.Images }}
* {{ $image }}
{{- end }}

{{- end }}

## Resources

| Name | Kind | Sync | Health |
|---|---|---|---|
{{- range $res := .Status.Resources }}
| {{$res.Name}} | {{$res.Kind}} | {{$res.Status}} | {{ if not $res.Health }}✓{{else}}{{$res.Health.Status}}{{end}} | 
{{- end }}
`
