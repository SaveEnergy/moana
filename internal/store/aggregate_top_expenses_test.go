package store

import (
	"context"
	"strconv"
	"testing"
	"time"
)

func TestListTopExpenseCategories_fromOnly_toOnly_matchBothBounds(t *testing.T) {
	t.Parallel()
	st := testStore(t)
	ctx := context.Background()
	uid, hid := mustCreateUserWithHousehold(t, st, ctx, "top-exp-bounds@example.com")
	a, err := st.CreateCategory(ctx, hid, "A", "", "")
	if err != nil {
		t.Fatal(err)
	}
	b, err := st.CreateCategory(ctx, hid, "B", "", "")
	if err != nil {
		t.Fatal(err)
	}
	day := time.Date(2026, 4, 11, 12, 0, 0, 0, time.UTC)
	if _, err := st.CreateTransaction(ctx, uid, hid, -3000, day, "a", &a); err != nil {
		t.Fatal(err)
	}
	if _, err := st.CreateTransaction(ctx, uid, hid, -8000, day, "b", &b); err != nil {
		t.Fatal(err)
	}
	from := day.Add(-time.Hour)
	to := day.Add(time.Hour)
	limit := 2
	gotBoth, err := st.ListTopExpenseCategories(ctx, hid, &from, &to, limit)
	if err != nil {
		t.Fatal(err)
	}
	gotFrom, err := st.ListTopExpenseCategories(ctx, hid, &from, nil, limit)
	if err != nil {
		t.Fatal(err)
	}
	gotTo, err := st.ListTopExpenseCategories(ctx, hid, nil, &to, limit)
	if err != nil {
		t.Fatal(err)
	}
	if len(gotBoth) != 2 || len(gotFrom) != 2 || len(gotTo) != 2 {
		t.Fatalf("both=%+v from=%+v to=%+v", gotBoth, gotFrom, gotTo)
	}
	for i := range gotBoth {
		if gotBoth[i].CategoryName != gotFrom[i].CategoryName || gotBoth[i].TotalCents != gotFrom[i].TotalCents {
			t.Fatalf("row %d both vs from: %+v %+v", i, gotBoth[i], gotFrom[i])
		}
		if gotBoth[i].CategoryName != gotTo[i].CategoryName || gotBoth[i].TotalCents != gotTo[i].TotalCents {
			t.Fatalf("row %d both vs to: %+v %+v", i, gotBoth[i], gotTo[i])
		}
	}
}

func TestListTopExpenseCategories_mostNegativeFirst(t *testing.T) {
	t.Parallel()
	st := testStore(t)
	ctx := context.Background()
	uid, hid := mustCreateUserWithHousehold(t, st, ctx, "top-exp@example.com")
	food, err := st.CreateCategory(ctx, hid, "Food", "", "")
	if err != nil {
		t.Fatal(err)
	}
	side, err := st.CreateCategory(ctx, hid, "Side", "", "")
	if err != nil {
		t.Fatal(err)
	}
	day := time.Date(2026, 4, 1, 12, 0, 0, 0, time.UTC)
	if _, err := st.CreateTransaction(ctx, uid, hid, -10000, day, "groceries", &food); err != nil {
		t.Fatal(err)
	}
	if _, err := st.CreateTransaction(ctx, uid, hid, -5000, day, "misc", &side); err != nil {
		t.Fatal(err)
	}
	from := day.Add(-time.Hour)
	to := day.Add(time.Hour)
	rows, err := st.ListTopExpenseCategories(ctx, hid, &from, &to, 2)
	if err != nil {
		t.Fatal(err)
	}
	if len(rows) != 2 {
		t.Fatalf("len=%d %+v", len(rows), rows)
	}
	if rows[0].TotalCents != -10000 || rows[0].CategoryName != "Food" {
		t.Fatalf("row0 %+v", rows[0])
	}
	if rows[1].TotalCents != -5000 || rows[1].CategoryName != "Side" {
		t.Fatalf("row1 %+v", rows[1])
	}
}

func TestListTopExpenseCategories_limitBelowOneDefaultsToFive(t *testing.T) {
	t.Parallel()
	st := testStore(t)
	ctx := context.Background()
	uid, hid := mustCreateUserWithHousehold(t, st, ctx, "top-exp-limit@example.com")
	day := time.Date(2026, 4, 2, 12, 0, 0, 0, time.UTC)
	for i := range 6 {
		name := "C" + strconv.Itoa(i)
		cid, err := st.CreateCategory(ctx, hid, name, "", "")
		if err != nil {
			t.Fatal(err)
		}
		amt := int64(-1000 * (i + 1))
		if _, err := st.CreateTransaction(ctx, uid, hid, amt, day, "x", &cid); err != nil {
			t.Fatal(err)
		}
	}
	from := day.Add(-time.Hour)
	to := day.Add(time.Hour)
	rows, err := st.ListTopExpenseCategories(ctx, hid, &from, &to, 0)
	if err != nil {
		t.Fatal(err)
	}
	if len(rows) != 5 {
		t.Fatalf("limit<1 must default to 5 rows, got %d %+v", len(rows), rows)
	}
	// Most negative sums first: -6000 … -2000.
	if rows[0].TotalCents != -6000 || rows[4].TotalCents != -2000 {
		t.Fatalf("order or truncation: %+v", rows)
	}
}

func TestListTopExpenseCategories_cancelledContext(t *testing.T) {
	t.Parallel()
	st := testStore(t)
	from := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	to := time.Date(2026, 1, 31, 23, 59, 59, 0, time.UTC)
	_, err := st.ListTopExpenseCategories(alreadyCancelledContext(t), 1, &from, &to, 5)
	assertErrIsContextCanceled(t, err)
}

func TestListTopExpenseCategories_nilNil_cancelledContext(t *testing.T) {
	t.Parallel()
	st := testStore(t)
	_, err := st.ListTopExpenseCategories(alreadyCancelledContext(t), 1, nil, nil, 5)
	assertErrIsContextCanceled(t, err)
}

func TestListTopExpenseCategories_nilNil_matchesWideRange(t *testing.T) {
	t.Parallel()
	st := testStore(t)
	ctx := context.Background()
	uid, hid := mustCreateUserWithHousehold(t, st, ctx, "top-nil-wide@example.com")
	food, err := st.CreateCategory(ctx, hid, "Food", "", "")
	if err != nil {
		t.Fatal(err)
	}
	side, err := st.CreateCategory(ctx, hid, "Side", "", "")
	if err != nil {
		t.Fatal(err)
	}
	day := time.Date(2026, 5, 1, 12, 0, 0, 0, time.UTC)
	if _, err := st.CreateTransaction(ctx, uid, hid, -10000, day, "g", &food); err != nil {
		t.Fatal(err)
	}
	if _, err := st.CreateTransaction(ctx, uid, hid, -5000, day, "m", &side); err != nil {
		t.Fatal(err)
	}
	wideFrom := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
	wideTo := time.Date(2030, 1, 1, 0, 0, 0, 0, time.UTC)
	gotWide, err := st.ListTopExpenseCategories(ctx, hid, &wideFrom, &wideTo, 2)
	if err != nil {
		t.Fatal(err)
	}
	gotNil, err := st.ListTopExpenseCategories(ctx, hid, nil, nil, 2)
	if err != nil {
		t.Fatal(err)
	}
	if len(gotWide) != len(gotNil) {
		t.Fatalf("len wide=%d nil=%d wide=%+v nil=%+v", len(gotWide), len(gotNil), gotWide, gotNil)
	}
	for i := range gotWide {
		a, b := gotWide[i], gotNil[i]
		if a.CategoryName != b.CategoryName || a.TotalCents != b.TotalCents {
			t.Fatalf("row %d wide=%+v nil=%+v", i, a, b)
		}
	}
}
