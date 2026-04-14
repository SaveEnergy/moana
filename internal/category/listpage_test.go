package category

import (
	"context"
	"testing"

	"moana/internal/dbutil"
)

func TestBuildCategoriesList(t *testing.T) {
	t.Parallel()
	st, db, err := dbutil.OpenStore(":memory:")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	ctx := context.Background()
	d, err := BuildCategoriesList(ctx, st, 1, "")
	if err != nil {
		t.Fatal(err)
	}
	if d.Error != "" || d.Categories != nil {
		t.Fatalf("empty user: %+v", d)
	}
}
