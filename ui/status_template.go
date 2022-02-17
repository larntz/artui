package ui

const statusTemplate string = `# Application: {{ .Name }}

**Git Repository:** 

**Sync Revision:** {{ .Status.Sync.Revision }}

## Resources

| Name | Kind | Sync | Health |
|---|---|---|---|
{{- range $res := .Status.Resources }}
| {{$res.Name}} | {{$res.Kind}} | {{$res.Status}} | {{ if not $res.Health }}✓{{else}}{{$res.Health.Status}}{{end}} | 
{{- end }}


`
