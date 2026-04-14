package category

import (
	"regexp"
	"strings"

	"moana/internal/store"
)

// ColorPreset is a curated accent for category cards (reference-style muted palette).
type ColorPreset struct {
	Hex   string
	Label string
}

// Presets is the ordered list shown in the category picker.
var Presets = []ColorPreset{
	{Hex: "#ea580c", Label: "Orange"},
	{Hex: "#6d8b7a", Label: "Sage"},
	{Hex: "#818cf8", Label: "Periwinkle"},
	{Hex: "#c45c26", Label: "Terracotta"},
	{Hex: "#db7eb1", Label: "Rose"},
	{Hex: "#7c3aed", Label: "Violet"},
	{Hex: "#9a3412", Label: "Rust"},
	{Hex: "#0ea5e9", Label: "Sky"},
	{Hex: "#a78bfa", Label: "Lilac"},
	{Hex: "#14b8a6", Label: "Mint"},
}

var hexColorRe = regexp.MustCompile(`(?i)^#[0-9a-f]{6}$`)

// TemplateColorPresets exposes presets to html/template (same shape as before).
func TemplateColorPresets() []ColorPreset {
	return Presets
}

// NormalizeColor returns canonical hex, or "" for auto / unset.
func NormalizeColor(s string) string {
	s = strings.TrimSpace(s)
	if s == "" {
		return ""
	}
	if !hexColorRe.MatchString(s) {
		return ""
	}
	return strings.ToLower(s)
}

func isPresetColor(hex string) bool {
	h := strings.ToLower(strings.TrimSpace(hex))
	for _, p := range Presets {
		if p.Hex == h {
			return true
		}
	}
	return false
}

// IsCustomColor reports whether the stored color is non-empty and not one of the preset swatches.
func IsCustomColor(c store.Category) bool {
	col := strings.TrimSpace(c.Color)
	if col == "" {
		return false
	}
	nc := NormalizeColor(col)
	if nc == "" {
		return true
	}
	return !isPresetColor(nc)
}

// HexOrDefault returns a #rrggbb string for templates (fallback when unset/invalid).
func HexOrDefault(c store.Category) string {
	x := NormalizeColor(c.Color)
	if x != "" {
		return x
	}
	return "#818cf8"
}
