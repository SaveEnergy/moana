package category

import "testing"

func TestFormFieldNames_nonEmpty(t *testing.T) {
	t.Parallel()
	for _, name := range []string{
		FieldName, FieldIcon, FieldID, FieldColor, FieldColorCustom,
	} {
		if name == "" {
			t.Fatal("empty form field name constant")
		}
	}
	if ColorPickCustom == "" {
		t.Fatal("empty ColorPickCustom")
	}
}
