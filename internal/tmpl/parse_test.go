package tmpl

import "testing"

func TestParse_embeddedTemplates(t *testing.T) {
	t.Parallel()
	tmpl, err := Parse()
	if err != nil {
		t.Fatal(err)
	}
	if tmpl == nil {
		t.Fatal("nil template root")
	}
	// Spot-check: every template name from ParseFS is registered.
	for _, name := range tmpl.Templates() {
		if name.Name() == "" {
			continue
		}
		if tmpl.Lookup(name.Name()) == nil {
			t.Fatalf("missing %q", name.Name())
		}
	}
}
