package assets

import (
	"embed"
	"io/fs"
)

//go:embed templates/*.html
var templatesFS embed.FS

//go:embed static
var staticRoot embed.FS

// TemplateFS is the embedded filesystem of *.html templates.
func TemplateFS() embed.FS {
	return templatesFS
}

// StaticFS is the filesystem rooted at the static/ directory (css/, js/).
func StaticFS() (fs.FS, error) {
	return fs.Sub(staticRoot, "static")
}
