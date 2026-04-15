package tmpl

import (
	"html/template"

	"moana/internal/assets"
)

// Parse loads *.html from the embedded template FS with shared helpers.
func Parse() (*template.Template, error) {
	return template.New("").Funcs(TemplateFuncMap()).ParseFS(assets.TemplateFS(), "templates/*.html")
}
