package historyview

import (
	"context"
	"net/url"
	"testing"
	"time"

	"moana/internal/dbutil"
	"moana/internal/passwordtest"
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
	if data.Error != InvalidDateRangeMessage {
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
	if data.Error != InvalidDateRangeMessage {
		t.Fatalf("Error=%q", data.Error)
	}
}

func TestBuildPage_smoke_listsTransactionsAndNav(t *testing.T) {
	t.Parallel()
	st := dbutil.MustOpenMemStore(t)
	ctx := context.Background()
	hash := passwordtest.MustHash(t, "x")
	uid, err := st.CreateUser(ctx, "hist-buildpage@example.com", hash, "user")
	if err != nil {
		t.Fatal(err)
	}
	u, err := st.GetUserByEmail(ctx, "hist-buildpage@example.com")
	if err != nil || u == nil {
		t.Fatal(err)
	}
	day := time.Date(2026, 8, 1, 12, 0, 0, 0, time.UTC)
	if _, err := st.CreateTransaction(ctx, uid, u.HouseholdID, -5000, day, "coffee", nil); err != nil {
		t.Fatal(err)
	}
	loc := time.UTC
	raw, err := url.Parse(RoutePath)
	if err != nil {
		t.Fatal(err)
	}
	data, err := BuildPage(ctx, st, u.HouseholdID, loc, raw, RoutePath)
	if err != nil {
		t.Fatal(err)
	}
	if data.Error != "" {
		t.Fatalf("Error=%q", data.Error)
	}
	if data.Truncated {
		t.Fatal("unexpected truncation")
	}
	if data.Nav.LinkAll == "" || data.Nav.SortNewest == "" {
		t.Fatalf("nav: %+v", data.Nav)
	}
	if len(data.Groups) == 0 {
		t.Fatal("expected at least one day group")
	}
}
