package category

import (
	"context"
	"errors"
	"testing"
	"time"

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

func TestBuildCategoriesList_expiredContext(t *testing.T) {
	t.Parallel()
	st := dbutil.MustOpenMemStore(t)
	deadlineCtx, cancel := context.WithDeadline(context.Background(), time.Now().Add(-time.Hour))
	defer cancel()
	_, err := BuildCategoriesList(deadlineCtx, st, 1, "")
	if err == nil {
		t.Fatal("expected error from expired context")
	}
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("want context.DeadlineExceeded, got %v", err)
	}
}
