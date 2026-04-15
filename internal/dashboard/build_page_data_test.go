package dashboard

import (
	"context"
	"testing"
	"time"

	"moana/internal/auth"
	"moana/internal/dbutil"
)

func TestBuildPageData_smoke(t *testing.T) {
	t.Parallel()
	st := dbutil.MustOpenMemStore(t)
	ctx := context.Background()
	hash, err := auth.HashPassword("x")
	if err != nil {
		t.Fatal(err)
	}
	uid, err := st.CreateUser(ctx, "dash-smoke@example.com", hash, "user")
	if err != nil {
		t.Fatal(err)
	}
	u, err := st.GetUserByEmail(ctx, "dash-smoke@example.com")
	if err != nil || u == nil {
		t.Fatal(err)
	}
	day := time.Date(2026, 8, 1, 12, 0, 0, 0, time.UTC)
	if _, err := st.CreateTransaction(ctx, uid, u.HouseholdID, -5000, day, "coffee", nil); err != nil {
		t.Fatal(err)
	}
	loc := time.UTC
	now := day
	data, err := BuildPageData(ctx, st, u.HouseholdID, loc, now, "30d")
	if err != nil {
		t.Fatal(err)
	}
	if data.StatsPeriod != "30d" {
		t.Fatalf("StatsPeriod %q", data.StatsPeriod)
	}
	if len(data.Recent) != 1 {
		t.Fatalf("recent count %d", len(data.Recent))
	}
	data12, err := BuildPageData(ctx, st, u.HouseholdID, loc, now, "12m")
	if err != nil {
		t.Fatal(err)
	}
	if data12.StatsPeriod != "12m" {
		t.Fatalf("12m period %q", data12.StatsPeriod)
	}
}
