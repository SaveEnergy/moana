package store

import (
	"context"
	"errors"
	"math"
	"strconv"
	"testing"
	"time"

	"moana/internal/passwordtest"
)

func TestListCategoryAmountsInRange_trimmedKindAcceptsExpense(t *testing.T) {
	t.Parallel()
	st := testStore(t)
	ctx := context.Background()
	hash := passwordtest.MustHash(t, "x")
	uid, err := st.CreateUser(ctx, "cat-trim-kind@example.com", hash, "user")
	if err != nil {
		t.Fatal(err)
	}
	u, err := st.GetUserByID(ctx, uid)
	if err != nil || u == nil {
		t.Fatal(err)
	}
	hid := u.HouseholdID
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

func TestListCategoryAmountsInRange_invalidKind(t *testing.T) {
	t.Parallel()
	st := testStore(t)
	ctx := context.Background()
	hash := passwordtest.MustHash(t, "x")
	uid, err := st.CreateUser(ctx, "cat-kind@example.com", hash, "user")
	if err != nil {
		t.Fatal(err)
	}
	u, err := st.GetUserByID(ctx, uid)
	if err != nil || u == nil {
		t.Fatal(err)
	}
	hid := u.HouseholdID
	from := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	to := time.Date(2026, 1, 31, 23, 59, 59, 0, time.UTC)

	_, err = st.ListCategoryAmountsInRange(ctx, hid, &from, &to, "net")
	if !errors.Is(err, ErrInvalidCategoryAmountKind) {
		t.Fatalf("got %v want %v", err, ErrInvalidCategoryAmountKind)
	}
}

func TestListCategoryAmountsInRange_expenseUncategorizedAndCategory(t *testing.T) {
	t.Parallel()
	st := testStore(t)
	ctx := context.Background()
	hash := passwordtest.MustHash(t, "x")
	uid, err := st.CreateUser(ctx, "cat-exp@example.com", hash, "user")
	if err != nil {
		t.Fatal(err)
	}
	u, err := st.GetUserByID(ctx, uid)
	if err != nil || u == nil {
		t.Fatal(err)
	}
	hid := u.HouseholdID
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
	hash := passwordtest.MustHash(t, "x")
	uid, err := st.CreateUser(ctx, "cat-inc@example.com", hash, "user")
	if err != nil {
		t.Fatal(err)
	}
	u, err := st.GetUserByID(ctx, uid)
	if err != nil || u == nil {
		t.Fatal(err)
	}
	hid := u.HouseholdID
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
	hash := passwordtest.MustHash(t, "x")
	uid, err := st.CreateUser(ctx, "cat-exp-min64@example.com", hash, "user")
	if err != nil {
		t.Fatal(err)
	}
	u, err := st.GetUserByID(ctx, uid)
	if err != nil || u == nil {
		t.Fatal(err)
	}
	hid := u.HouseholdID
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

func TestListTopExpenseCategories_mostNegativeFirst(t *testing.T) {
	t.Parallel()
	st := testStore(t)
	ctx := context.Background()
	hash := passwordtest.MustHash(t, "x")
	uid, err := st.CreateUser(ctx, "top-exp@example.com", hash, "user")
	if err != nil {
		t.Fatal(err)
	}
	u, err := st.GetUserByID(ctx, uid)
	if err != nil || u == nil {
		t.Fatal(err)
	}
	hid := u.HouseholdID
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
	hash := passwordtest.MustHash(t, "x")
	uid, err := st.CreateUser(ctx, "top-exp-limit@example.com", hash, "user")
	if err != nil {
		t.Fatal(err)
	}
	u, err := st.GetUserByID(ctx, uid)
	if err != nil || u == nil {
		t.Fatal(err)
	}
	hid := u.HouseholdID
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
	hash := passwordtest.MustHash(t, "x")
	uid, err := st.CreateUser(ctx, "cat-nil-wide@example.com", hash, "user")
	if err != nil {
		t.Fatal(err)
	}
	u, err := st.GetUserByID(ctx, uid)
	if err != nil || u == nil {
		t.Fatal(err)
	}
	hid := u.HouseholdID
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
	hash := passwordtest.MustHash(t, "x")
	uid, err := st.CreateUser(ctx, "top-nil-wide@example.com", hash, "user")
	if err != nil {
		t.Fatal(err)
	}
	u, err := st.GetUserByID(ctx, uid)
	if err != nil || u == nil {
		t.Fatal(err)
	}
	hid := u.HouseholdID
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
