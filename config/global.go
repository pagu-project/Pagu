package config

const (
	PriceCacheKey = "PriceCacheKey"
)

var HelpCommandTemplate string = `<table>{{range .}}<tr><td>{{ .Name }}</td><td>{{ .Desc }}</td></tr>{{end}}</table>`
