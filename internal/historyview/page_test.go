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

func TestBuildPage_partialDateRange_requiresBothBounds(t *testing.T) {
	t.Parallel()
	st := dbutil.MustOpenMemStore(t)
	ctx := context.Background()
	loc := time.UTC

	for _, raw := range []string{
		"/history?from=2026-01-01",
		"/history?to=2026-01-31",
	} {
		u, err := url.Parse(raw)
		if err != nil {
			t.Fatal(err)
		}
		d, err := BuildPage(ctx, st, 1, loc, u, u.String())
		if err != nil {
			t.Fatal(err)
		}
		if d.Error != "Invalid date range." {
			t.Fatalf("%s: got %q", raw, d.Error)
		}
	}
}
