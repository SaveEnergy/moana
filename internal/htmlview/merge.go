package htmlview

import (
	"html/template"
	"maps"
)

// MergeFuncMaps overlays maps left-to-right; later maps override earlier keys.
func MergeFuncMaps(inputs ...template.FuncMap) template.FuncMap {
	switch len(inputs) {
	case 0:
		return make(template.FuncMap)
	case 1:
		return maps.Clone(inputs[0])
	case 2:
		m0, m1 := inputs[0], inputs[1]
		out := make(template.FuncMap, len(m0)+len(m1))
		for k, v := range m0 {
			out[k] = v
		}
		for k, v := range m1 {
			out[k] = v
		}
		return out
	}
	n := 0
	for _, m := range inputs {
		n += len(m)
	}
	out := make(template.FuncMap, n)
	for _, m := range inputs {
		for k, v := range m {
			out[k] = v
		}
	}
	return out
}
