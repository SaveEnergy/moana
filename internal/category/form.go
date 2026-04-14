package category

import (
	"net/http"
	"strings"
)

// ParseColorFromForm reads color / color_custom from a POST form (category create/update modal).
func ParseColorFromForm(r *http.Request) string {
	c := strings.TrimSpace(r.FormValue("color"))
	if c == "custom" {
		return NormalizeColor(r.FormValue("color_custom"))
	}
	return NormalizeColor(c)
}
