package historyview

import (
	"context"
	"net/url"
	"testing"
	"time"

	"moana/internal/dbutil"
)

func TestBuildPage_invalidDateRange_returnsErrorPayload(t *testing.T) {
	t.Parallel()
	st := dbutil.MustOpenMemStore(t)
	ctx := context.Background()
	u, err := url.Parse(RoutePath + "?from=not-a-date&to=2026-01-31")
	if err != nil {
		t.Fatal(err)
	}
	data, err := BuildPage(ctx, st, 1, time.UTC, u, RoutePath)
	if err != nil {
		t.Fatal(err)
	}
	if data.Error != "Invalid date range." {
		t.Fatalf("Error=%q want invalid date message", data.Error)
	}
	if !data.FilterActive {
		t.Fatal("FilterActive want true so UI keeps date fields")
	}
	if data.Groups != nil {
		t.Fatalf("Groups=%v want nil", data.Groups)
	}
}

func TestBuildPage_partialDateOnly_returnsErrorPayload(t *testing.T) {
	t.Parallel()
	st := dbutil.MustOpenMemStore(t)
	ctx := context.Background()
	u, err := url.Parse(RoutePath + "?from=2026-01-01")
	if err != nil {
		t.Fatal(err)
	}
	data, err := BuildPage(ctx, st, 1, time.UTC, u, RoutePath)
	if err != nil {
		t.Fatal(err)
	}
	if data.Error != "Invalid date range." {
		t.Fatalf("Error=%q", data.Error)
	}
}
