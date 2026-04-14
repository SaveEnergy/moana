package tmpl

import (
	"html/template"
	"net/url"

	"moana/internal/assets"
	"moana/internal/category"
	"moana/internal/household"
	"moana/internal/htmlview"
)

// Parse loads *.html from the embedded template FS with shared helpers.
func Parse() (*template.Template, error) {
	handlerFuncs := template.FuncMap{
		"categoryPickerAccent":      category.PickerAccent,
		"categoryPickerHint":        category.PickerHint,
		"categoryIconChoices":       templateCategoryIconChoices,
		"categoryColorPresets":      category.TemplateColorPresets,
		"categoryDisplayIcon":       categoryDisplayIcon,
		"dashboardRecentIcon":       dashboardRecentIcon,
		"categoryNoneIcon":          categoryNoneIconHTML,
		"categoryIsCustomColor":     category.IsCustomColor,
		"categoryColorHex":          category.HexOrDefault,
		"categoryIconFormValue":     templateCategoryIconFormValue,
		"settingsCanRemoveMember":   household.CanRemoveMemberRow,
		"settingsCanLeaveHousehold": household.CanLeave,
		"urlQuery":                  func(s string) string { return url.QueryEscape(s) },
	}
	funcs := htmlview.MergeFuncMaps(
		htmlview.TemplateFuncs(),
		htmlview.IconFuncs(),
		handlerFuncs,
	)
	return template.New("").Funcs(funcs).ParseFS(assets.TemplateFS(), "templates/*.html")
}
