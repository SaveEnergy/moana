package store

import (
	"context"
	"testing"
	"time"

	"moana/internal/auth"
)

func TestSumIncomeExpenseCentsInRange_matchesSplitQueries(t *testing.T) {
	t.Parallel()
	st := testStore(t)
	ctx := context.Background()
	hash, err := auth.HashPassword("x")
	if err != nil {
		t.Fatal(err)
	}
	uid, err := st.CreateUser(ctx, "agg@example.com", hash, "user")
	if err != nil {
		t.Fatal(err)
	}
	u, err := st.GetUserByID(ctx, uid)
	if err != nil || u == nil {
		t.Fatal(err)
	}
	hid := u.HouseholdID

	day := time.Date(2026, 7, 1, 12, 0, 0, 0, time.UTC)
	if _, err := st.CreateTransaction(ctx, uid, hid, 10000, day, "in", nil); err != nil {
		t.Fatal(err)
	}
	if _, err := st.CreateTransaction(ctx, uid, hid, -3000, day, "out", nil); err != nil {
		t.Fatal(err)
	}

	from := day.Add(-time.Hour)
	to := day.Add(time.Hour)

	inc, exp, err := st.SumIncomeExpenseCentsInRange(ctx, hid, &from, &to)
	if err != nil {
		t.Fatal(err)
	}
	wantInc, err := st.SumAmountCentsByKind(ctx, hid, &from, &to, "income")
	if err != nil {
		t.Fatal(err)
	}
	wantExp, err := st.SumAmountCentsByKind(ctx, hid, &from, &to, "expense")
	if err != nil {
		t.Fatal(err)
	}
	if inc != wantInc || exp != wantExp {
		t.Fatalf("combined %+v %+v want income=%d expense=%d", inc, exp, wantInc, wantExp)
	}
	if inc+exp != wantInc+wantExp {
		t.Fatal("net mismatch")
	}
}
