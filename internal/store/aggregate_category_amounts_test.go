package store

import (
	"context"
	"errors"
	"math"
	"testing"
	"time"
)

func TestListCategoryAmountsInRange_trimmedKindAcceptsExpense(t *testing.T) {
	t.Parallel()
	st := testStore(t)
	ctx := context.Background()
	uid, hid := mustCreateUserWithHousehold(t, st, ctx, "cat-trim-kind@example.com")
	from := time.Date(2026, 2, 1, 0, 0, 0, 0, time.UTC)
	to := time.Date(2026, 2, 28, 23, 59, 59, 0, time.UTC)
	if _, err := st.CreateTransaction(ctx, uid, hid, -5000, time.Date(2026, 2, 1, 12, 0, 0, 0, time.UTC), "x", nil); err != nil {
		t.Fatal(err)
	}
	rows, err := st.ListCategoryAmountsInRange(ctx, hid, &from, &to, " expense ")
	if err != nil {
		t.Fatal(err)
	}
	if len(rows) != 1 || rows[0].AmountCents != 5000 {
		t.Fatalf("got %+v", rows)
	}
}

func TestListCategoryAmountsInRange_fromOnly_toOnly_matchBothBounds(t *testing.T) {
	t.Parallel()
	st := testStore(t)
	ctx := context.Background()
	uid, hid := mustCreateUserWithHousehold(t, st, ctx, "cat-one-bound@example.com")
	catID, err := st.CreateCategory(ctx, hid, "Food", "", "")
	if err != nil {
		t.Fatal(err)
	}
	day := time.Date(2026, 4, 10, 12, 0, 0, 0, time.UTC)
	if _, err := st.CreateTransaction(ctx, uid, hid, -9000, day, "lunch", &catID); err != nil {
		t.Fatal(err)
	}
	from := day.Add(-time.Hour)
	to := day.Add(time.Hour)
	gotBoth, err := st.ListCategoryAmountsInRange(ctx, hid, &from, &to, "expense")
	if err != nil {
		t.Fatal(err)
	}
	gotFrom, err := st.ListCategoryAmountsInRange(ctx, hid, &from, nil, "expense")
	if err != nil {
		t.Fatal(err)
	}
	gotTo, err := st.ListCategoryAmountsInRange(ctx, hid, nil, &to, "expense")
	if err != nil {
		t.Fatal(err)
	}
	if len(gotBoth) != 1 || len(gotFrom) != 1 || len(gotTo) != 1 {
		t.Fatalf("both=%+v from=%+v to=%+v", gotBoth, gotFrom, gotTo)
	}
	if gotBoth[0].AmountCents != gotFrom[0].AmountCents || gotBoth[0].AmountCents != gotTo[0].AmountCents {
		t.Fatalf("amount mismatch both=%d from=%d to=%d", gotBoth[0].AmountCents, gotFrom[0].AmountCents, gotTo[0].AmountCents)
	}
}

func TestListCategoryAmountsInRange_invalidKind(t *testing.T) {
	t.Parallel()
	st := testStore(t)
	ctx := context.Background()
	_, hid := mustCreateUserWithHousehold(t, st, ctx, "cat-kind@example.com")
	from := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	to := time.Date(2026, 1, 31, 23, 59, 59, 0, time.UTC)

	_, err := st.ListCategoryAmountsInRange(ctx, hid, &from, &to, "net")
	if !errors.Is(err, ErrInvalidCategoryAmountKind) {
		t.Fatalf("got %v want %v", err, ErrInvalidCategoryAmountKind)
	}
}

func TestListCategoryAmountsInRange_expenseUncategorizedAndCategory(t *testing.T) {
	t.Parallel()
	st := testStore(t)
	ctx := context.Background()
	uid, hid := mustCreateUserWithHousehold(t, st, ctx, "cat-exp@example.com")
	catID, err := st.CreateCategory(ctx, hid, "Food", "utensils", "#112233")
	if err != nil {
		t.Fatal(err)
	}
	day := time.Date(2026, 6, 15, 12, 0, 0, 0, time.UTC)
	if _, err := st.CreateTransaction(ctx, uid, hid, -10000, day, "groceries", &catID); err != nil {
		t.Fatal(err)
	}
	if _, err := st.CreateTransaction(ctx, uid, hid, -4000, day, "cash", nil); err != nil {
		t.Fatal(err)
	}
	from := day.Add(-time.Hour)
	to := day.Add(time.Hour)

	rows, err := st.ListCategoryAmountsInRange(ctx, hid, &from, &to, "expense")
	if err != nil {
		t.Fatal(err)
	}
	if len(rows) != 2 {
		t.Fatalf("len=%d %+v", len(rows), rows)
	}
	// Largest expense (most negative sum) first → Food 10000 then uncategorized 4000.
	if rows[0].Name != "Food" || rows[0].AmountCents != 10000 || rows[0].Icon != "utensils" {
		t.Fatalf("row0 %+v", rows[0])
	}
	if rows[1].Name != "Uncategorized" || rows[1].AmountCents != 4000 {
		t.Fatalf("row1 %+v", rows[1])
	}
}

func TestListCategoryAmountsInRange_incomeOrderedBySize(t *testing.T) {
	t.Parallel()
	st := testStore(t)
	ctx := context.Background()
	uid, hid := mustCreateUserWithHousehold(t, st, ctx, "cat-inc@example.com")
	small, err := st.CreateCategory(ctx, hid, "Side", "", "")
	if err != nil {
		t.Fatal(err)
	}
	large, err := st.CreateCategory(ctx, hid, "Salary", "", "")
	if err != nil {
		t.Fatal(err)
	}
	day := time.Date(2026, 3, 1, 10, 0, 0, 0, time.UTC)
	if _, err := st.CreateTransaction(ctx, uid, hid, 5000, day, "bonus", &small); err != nil {
		t.Fatal(err)
	}
	if _, err := st.CreateTransaction(ctx, uid, hid, 50000, day, "pay", &large); err != nil {
		t.Fatal(err)
	}
	from := day.Add(-time.Hour)
	to := day.Add(time.Hour)

	rows, err := st.ListCategoryAmountsInRange(ctx, hid, &from, &to, "income")
	if err != nil {
		t.Fatal(err)
	}
	if len(rows) != 2 {
		t.Fatalf("len=%d", len(rows))
	}
	if rows[0].Name != "Salary" || rows[0].AmountCents != 50000 {
		t.Fatalf("row0 %+v", rows[0])
	}
	if rows[1].Name != "Side" || rows[1].AmountCents != 5000 {
		t.Fatalf("row1 %+v", rows[1])
	}
}

func TestListCategoryAmountsInRange_expenseMinInt64SumUsesAbsCents(t *testing.T) {
	t.Parallel()
	st := testStore(t)
	ctx := context.Background()
	uid, hid := mustCreateUserWithHousehold(t, st, ctx, "cat-exp-min64@example.com")
	day := time.Date(2026, 6, 15, 12, 0, 0, 0, time.UTC)
	if _, err := st.CreateTransaction(ctx, uid, hid, math.MinInt64, day, "edge", nil); err != nil {
		t.Fatal(err)
	}
	from := day.Add(-time.Hour)
	to := day.Add(time.Hour)
	rows, err := st.ListCategoryAmountsInRange(ctx, hid, &from, &to, "expense")
	if err != nil {
		t.Fatal(err)
	}
	if len(rows) != 1 {
		t.Fatalf("len=%d %+v", len(rows), rows)
	}
	if rows[0].AmountCents != math.MaxInt64 {
		t.Fatalf("got %d want MaxInt64 (negating MinInt64 overflows)", rows[0].AmountCents)
	}
}

func TestListCategoryAmountsInRange_cancelledContext(t *testing.T) {
	t.Parallel()
	st := testStore(t)
	from := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	to := time.Date(2026, 1, 31, 23, 59, 59, 0, time.UTC)
	_, err := st.ListCategoryAmountsInRange(alreadyCancelledContext(t), 1, &from, &to, "expense")
	assertErrIsContextCanceled(t, err)
}

func TestListCategoryAmountsInRange_nilNil_cancelledContext(t *testing.T) {
	t.Parallel()
	st := testStore(t)
	_, err := st.ListCategoryAmountsInRange(alreadyCancelledContext(t), 1, nil, nil, "expense")
	assertErrIsContextCanceled(t, err)
}

func TestListCategoryAmountsInRange_nilNil_expenseMatchesWideRange(t *testing.T) {
	t.Parallel()
	st := testStore(t)
	ctx := context.Background()
	uid, hid := mustCreateUserWithHousehold(t, st, ctx, "cat-nil-wide@example.com")
	catID, err := st.CreateCategory(ctx, hid, "Food", "utensils", "#112233")
	if err != nil {
		t.Fatal(err)
	}
	day := time.Date(2026, 6, 15, 12, 0, 0, 0, time.UTC)
	if _, err := st.CreateTransaction(ctx, uid, hid, -10000, day, "groceries", &catID); err != nil {
		t.Fatal(err)
	}
	if _, err := st.CreateTransaction(ctx, uid, hid, -4000, day, "cash", nil); err != nil {
		t.Fatal(err)
	}
	wideFrom := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
	wideTo := time.Date(2030, 1, 1, 0, 0, 0, 0, time.UTC)
	gotWide, err := st.ListCategoryAmountsInRange(ctx, hid, &wideFrom, &wideTo, "expense")
	if err != nil {
		t.Fatal(err)
	}
	gotNil, err := st.ListCategoryAmountsInRange(ctx, hid, nil, nil, "expense")
	if err != nil {
		t.Fatal(err)
	}
	if len(gotWide) != len(gotNil) {
		t.Fatalf("len wide=%d nil=%d wide=%+v nil=%+v", len(gotWide), len(gotNil), gotWide, gotNil)
	}
	for i := range gotWide {
		a, b := gotWide[i], gotNil[i]
		if a.Name != b.Name || a.AmountCents != b.AmountCents || a.Icon != b.Icon || a.Color != b.Color {
			t.Fatalf("row %d wide=%+v nil=%+v", i, a, b)
		}
	}
}
