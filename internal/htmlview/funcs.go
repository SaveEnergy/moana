package htmlview

import "html/template"

// TemplateFuncs returns FuncMap entries safe to merge with handler-specific helpers.
func TemplateFuncs() template.FuncMap {
	return template.FuncMap{
		"formatEUR":           FormatEUR,
		"formatEURAbs":        FormatEURAbs,
		"formatCompactEUR":    FormatCompactEUR,
		"formatRFC3339UTC":    FormatRFC3339UTC,
		"formatPercentSigned": FormatPercentSigned,
		"isNegFloat":          IsNegFloat,
		"userInitial":         UserInitial,
		"adminDisplayName":    AdminDisplayName,
		"adminRoleLabel":      AdminRoleLabel,
		"profileDisplayName":  ProfileDisplayName,
		"profileInitial":      ProfileInitial,
		"householdRoleLabel":  HouseholdRoleLabel,
		"attr":                Attr,
	}
}
