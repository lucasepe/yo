package template

import (
	"bytes"
	"text/template"
)

func ExecuteInline(data interface{}, s string) ([]byte, error) {
	// Build function map.
	funcMap := TxtFuncMap()

	// Build the template
	t := template.New("main")
	t.Funcs(funcMap)

	_, err := t.Parse(s)
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	if err := t.Execute(&buf, data); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
