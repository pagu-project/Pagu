package template

import (
	"bytes"
	"html/template"
)

func ExecuteHTML(tmpl string, keyValue any) (string, error) {
	b := bytes.Buffer{}
	tp, err := template.New("").Parse(tmpl)
	if err != nil {
		return "", err
	}

	if err := tp.Execute(&b, keyValue); err != nil {
		return "", err
	}

	return b.String(), nil
}
