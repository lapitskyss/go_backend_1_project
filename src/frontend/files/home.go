package files

import (
	"embed"
	"html/template"
)

//go:embed templates
var temp embed.FS

type Templates struct {
	HomeTemplate *template.Template
}

func InitTemplates() (*Templates, error) {
	tpl, err := template.ParseFS(temp, "templates/index.html")
	if err != nil {
		return nil, err
	}

	return &Templates{
		HomeTemplate: tpl,
	}, nil
}
