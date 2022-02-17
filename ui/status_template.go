package ui

const statusTemplate string = `# Application: {{ .Name }}

` + "```" + `yaml
{{ .LongStatus }}
` + "```\n"
