package handlers

import (
	"embed"
	"html/template"
	"io/fs"
	"net/url"
)

//go:embed templates/*.html
var templatesFS embed.FS

// ParseTemplates loads HTML templates with shared helpers.
func ParseTemplates() (*template.Template, error) {
	funcs := template.FuncMap{
		"formatEUR":         templateFormatEUR,
		"formatEURAbs":      templateFormatEURAbs,
		"formatCompactEUR":  templateFormatCompactEUR,
		"formatLocal":       templateFormatLocal,
		"formatLocalTime":   templateFormatLocalTime,
		"formatPercentSigned": templatePercentSigned,
		"isNegFloat":          templateIsNegFloat,
		"userInitial":       templateUserInitial,
		"urlQuery":          func(s string) string { return url.QueryEscape(s) },
	}
	return template.New("").Funcs(funcs).ParseFS(templatesFS, "templates/*.html")
}

// StaticFS returns embedded static assets.
func StaticFS() (fs.FS, error) {
	return fs.Sub(staticFS, "static")
}

//go:embed static
var staticFS embed.FS
