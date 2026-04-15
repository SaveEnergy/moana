package tmpl

import (
	"html/template"
	"net/url"

	"moana/internal/category"
	"moana/internal/household"
	"moana/internal/htmlview"
)

// TemplateFuncMap merges htmlview helpers with app-specific template functions.
// Call from [Parse] or tests that need the same registration without loading templates.
func TemplateFuncMap() template.FuncMap {
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
	return htmlview.MergeFuncMaps(
		htmlview.TemplateFuncs(),
		htmlview.IconFuncs(),
		handlerFuncs,
	)
}
