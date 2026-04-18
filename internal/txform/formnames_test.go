package txform

import "testing"

func TestFieldNames_nonEmpty(t *testing.T) {
	t.Parallel()
	for _, name := range []string{
		FieldAmount, FieldOccurredOn, FieldDescription, FieldCategoryID, FieldKind,
	} {
		if name == "" {
			t.Fatal("empty field name constant")
		}
	}
}
