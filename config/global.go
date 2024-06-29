package config

var (
	TargetMaskMain      = 1
	TargetMaskTest      = 2
	TargetMaskModerator = 4

	TargetMaskAll = TargetMaskMain | TargetMaskTest | TargetMaskModerator
)

const (
	PriceCacheKey = "PriceCacheKey"
)

var HelpCommandTemplate string = `<table>{{range .}}<tr><td>{{ .Name }}</td><td>{{ .Desc }}</td></tr>{{end}}</table>`
