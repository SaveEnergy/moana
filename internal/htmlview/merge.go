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
	default:
		n := 0
		for _, m := range inputs {
			n += len(m)
		}
		out := make(template.FuncMap, n)
		for _, m := range inputs {
			maps.Copy(out, m)
		}
		return out
	}
}
