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
		maps.Copy(out, m0)
		maps.Copy(out, m1)
		return out
	case 3:
		m0, m1, m2 := inputs[0], inputs[1], inputs[2]
		out := make(template.FuncMap, len(m0)+len(m1)+len(m2))
		maps.Copy(out, m0)
		maps.Copy(out, m1)
		maps.Copy(out, m2)
		return out
	case 4:
		m0, m1, m2, m3 := inputs[0], inputs[1], inputs[2], inputs[3]
		out := make(template.FuncMap, len(m0)+len(m1)+len(m2)+len(m3))
		maps.Copy(out, m0)
		maps.Copy(out, m1)
		maps.Copy(out, m2)
		maps.Copy(out, m3)
		return out
	}
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
