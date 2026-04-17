package store

import (
	"context"
	"testing"
	"time"

	"moana/internal/auth"
)

func TestDailyAbsMovementByLocalDate_bucketingBerlin(t *testing.T) {
	t.Parallel()
	st := testStore(t)
	ctx := context.Background()

	hash, err := auth.HashPassword("pw-movement-test")
	if err != nil {
		t.Fatal(err)
	}
	uid, err := st.CreateUser(ctx, "movement-tz@example.com", hash, "user")
	if err != nil {
		t.Fatal(err)
	}
	u, err := st.GetUserByEmail(ctx, "movement-tz@example.com")
	if err != nil || u == nil {
		t.Fatal(err)
	}
	hid := u.HouseholdID

	loc, err := time.LoadLocation("Europe/Berlin")
	if err != nil {
		t.Fatal(err)
	}
	// 2026-01-15 23:30 UTC → 2026-01-16 00:30 CET (local calendar day is Jan 16).
	occ := time.Date(2026, 1, 15, 23, 30, 0, 0, time.UTC)
	if _, err := st.CreateTransaction(ctx, uid, hid, -4000, occ, "late", nil); err != nil {
		t.Fatal(err)
	}
	occ2 := time.Date(2026, 1, 16, 10, 0, 0, 0, time.UTC)
	if _, err := st.CreateTransaction(ctx, uid, hid, -2000, occ2, "noon", nil); err != nil {
		t.Fatal(err)
	}

	fromUTC := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	toUTC := time.Date(2026, 1, 31, 23, 59, 59, 999999999, time.UTC)

	byDay, err := st.DailyAbsMovementByLocalDate(ctx, hid, fromUTC, toUTC, loc)
	if err != nil {
		t.Fatal(err)
	}
	// Both fall on 2026-01-16 in Berlin.
	if got := byDay["2026-01-16"]; got != 6000 {
		t.Fatalf("2026-01-16 sum: got %d want 6000 (map=%v)", got, byDay)
	}
	if len(byDay) != 1 {
		t.Fatalf("expected one local day, got %v", byDay)
	}
}

func TestDailyAbsMovementByLocalDate_nilLocationMatchesUTC(t *testing.T) {
	t.Parallel()
	st := testStore(t)
	ctx := context.Background()
	hash, err := auth.HashPassword("pw-movement-nil")
	if err != nil {
		t.Fatal(err)
	}
	uid, err := st.CreateUser(ctx, "movement-nil@example.com", hash, "user")
	if err != nil {
		t.Fatal(err)
	}
	u, err := st.GetUserByEmail(ctx, "movement-nil@example.com")
	if err != nil || u == nil {
		t.Fatal(err)
	}
	hid := u.HouseholdID
	occ := time.Date(2026, 1, 15, 12, 0, 0, 0, time.UTC)
	if _, err := st.CreateTransaction(ctx, uid, hid, -1000, occ, "x", nil); err != nil {
		t.Fatal(err)
	}
	fromUTC := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	toUTC := time.Date(2026, 1, 31, 23, 59, 59, 999999999, time.UTC)
	mNil, err := st.DailyAbsMovementByLocalDate(ctx, hid, fromUTC, toUTC, nil)
	if err != nil {
		t.Fatal(err)
	}
	mUTC, err := st.DailyAbsMovementByLocalDate(ctx, hid, fromUTC, toUTC, time.UTC)
	if err != nil {
		t.Fatal(err)
	}
	if len(mNil) != len(mUTC) {
		t.Fatalf("nil=%v utc=%v", mNil, mUTC)
	}
	for k, v := range mNil {
		if mUTC[k] != v {
			t.Fatalf("key %q nil=%d utc=%d", k, v, mUTC[k])
		}
	}
}
