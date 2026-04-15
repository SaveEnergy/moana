// Package htmlview registers shared html/template helpers: formatting and escaping in
// [TemplateFuncs] (format.go, format_money.go; RFC3339 uses moana/internal/timeutil for SQLite alignment), Lucide SVG in [IconFuncs]; merge with [MergeFuncMaps] for handler-specific funcs.
// Name and avatar text helpers live in display_identity.go; role labels in display_roles.go.
package htmlview
