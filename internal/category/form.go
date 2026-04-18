package category

import (
	"net/http"
	"strings"
)

// ParseColorFromForm reads color / color_custom from a POST form (category create/update modal).
func ParseColorFromForm(r *http.Request) string {
	c := strings.TrimSpace(r.FormValue(FieldColor))
	if c == ColorPickCustom {
		return NormalizeColor(r.FormValue(FieldColorCustom))
	}
	return NormalizeColor(c)
}
