// Package tmpl parses embedded HTML templates and registers template.FuncMap helpers
// ([TemplateFuncMap], [Parse]). It sits between embedded assets and HTTP handlers so the HTTP layer stays focused on request/response.
package tmpl
