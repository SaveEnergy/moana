package historyview

import (
	"context"
	"net/url"
	"testing"
	"time"

	"moana/internal/dbutil"
)

func TestBuildPage_invalidDateRange(t *testing.T) {
	t.Parallel()
	st := dbutil.MustOpenMemStore(t)
	ctx := context.Background()
	loc := time.UTC
	u, err := url.Parse("/history?from=not-a-date&to=2020-01-02")
	if err != nil {
		t.Fatal(err)
	}
	d, err := BuildPage(ctx, st, 1, loc, u, u.String())
	if err != nil {
		t.Fatal(err)
	}
	if d.Error != "Invalid date range." {
		t.Fatalf("got %q", d.Error)
	}
}
