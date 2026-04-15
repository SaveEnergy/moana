package store

import (
	"context"
	"errors"
	"testing"
	"time"

	"moana/internal/auth"
)

func TestListCategoryAmountsInRange_invalidKind(t *testing.T) {
	t.Parallel()
	st := testStore(t)
	ctx := context.Background()
	hash, err := auth.HashPassword("x")
	if err != nil {
		t.Fatal(err)
	}
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
	hash, err := auth.HashPassword("x")
	if err != nil {
		t.Fatal(err)
	}
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
	hash, err := auth.HashPassword("x")
	if err != nil {
		t.Fatal(err)
	}
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
