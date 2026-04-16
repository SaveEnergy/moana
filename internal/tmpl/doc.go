// Package tmpl parses embedded HTML templates and registers template.FuncMap helpers
// ([TemplateFuncMap], [Parse]; parse_test.go loads embedded templates; funcs_test.go checks merged func keys). It sits between embedded assets and HTTP handlers so the HTTP layer stays focused on request/response.
package tmpl
