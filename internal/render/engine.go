package render

import "html/template"

// Engine holds parsed templates and performs layout / simple renders.
type Engine struct {
	Templates *template.Template
}
