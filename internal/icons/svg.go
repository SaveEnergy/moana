package icons

import (
	"fmt"
	"html"
	"html/template"
)

// ValidID reports whether id is a known icon (including nav/UI extras).
func ValidID(id string) bool {
	if id == "" {
		return false
	}
	_, ok := byID[id]
	return ok
}

// SVG renders a full <svg> element with the given CSS class, or falls back to the first palette icon.
func SVG(id, class string) template.HTML {
	def, ok := byID[id]
	if !ok || def.Inner == "" {
		if len(CategoryPalette) > 0 {
			def = CategoryPalette[0]
		} else {
			return template.HTML("")
		}
	}
	s := fmt.Sprintf(
		`<svg class="%s" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true">%s</svg>`,
		html.EscapeString(class), def.Inner)
	return template.HTML(s)
}

// Inner returns raw SVG path children for embedding (e.g. chart labels), or "".
func Inner(id string) string {
	def, ok := byID[id]
	if !ok {
		return ""
	}
	return def.Inner
}
