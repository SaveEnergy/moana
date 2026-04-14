package htmlview

import (
	"html/template"

	"moana/internal/icons"
)

// IconFuncs returns template helpers that emit Lucide SVG markup.
func IconFuncs() template.FuncMap {
	return template.FuncMap{
		"moanaIcon":      MoanaIconGrid,
		"moanaIconNav":   MoanaIconNav,
		"moanaIconPlain": MoanaIconPlain,
	}
}

// MoanaIconPlain renders a Lucide icon with the default class.
func MoanaIconPlain(id string) template.HTML {
	return icons.SVG(id, "moana-icon")
}

// MoanaIconGrid renders an icon sized for category grids.
func MoanaIconGrid(id string) template.HTML {
	return icons.SVG(id, "moana-icon moana-icon--grid")
}

// MoanaIconNav renders an icon for navigation rows.
func MoanaIconNav(id string) template.HTML {
	return icons.SVG(id, "moana-icon moana-icon--nav")
}
