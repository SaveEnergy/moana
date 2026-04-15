// Package tmpl parses embedded HTML templates and registers template.FuncMap helpers
// ([TemplateFuncMap], [Parse]; parse_test.go ensures embedded templates load). It sits between embedded assets and HTTP handlers so the HTTP layer stays focused on request/response.
package tmpl
