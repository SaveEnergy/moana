package htmlview

import "html/template"

// MergeFuncMaps overlays maps left-to-right; later maps override earlier keys.
func MergeFuncMaps(maps ...template.FuncMap) template.FuncMap {
	n := 0
	for _, m := range maps {
		n += len(m)
	}
	out := make(template.FuncMap, n)
	for _, m := range maps {
		for k, v := range m {
			out[k] = v
		}
	}
	return out
}
