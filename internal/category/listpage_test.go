package category

import (
	"context"
	"testing"

	"moana/internal/dbutil"
)

func TestBuildCategoriesList(t *testing.T) {
	t.Parallel()
	st := dbutil.MustOpenMemStore(t)
	ctx := context.Background()
	d, err := BuildCategoriesList(ctx, st, 1, "")
	if err != nil {
		t.Fatal(err)
	}
	if d.Error != "" || len(d.Categories) != 0 {
		t.Fatalf("empty household: %+v", d)
	}
}
