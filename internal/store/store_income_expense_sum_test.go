package store

import (
	"context"
	"testing"
	"time"

	"moana/internal/passwordtest"
)

func TestSumIncomeExpenseCentsInRange_nilNil_matchesSplitQueries(t *testing.T) {
	t.Parallel()
	st := testStore(t)
	ctx := context.Background()
	hash := passwordtest.MustHash(t, "x")
	uid, err := st.CreateUser(ctx, "agg-nil@example.com", hash, "user")
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

	inc, exp, err := st.SumIncomeExpenseCentsInRange(ctx, hid, nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	wantInc, err := st.SumAmountCentsByKind(ctx, hid, nil, nil, "income")
	if err != nil {
		t.Fatal(err)
	}
	wantExp, err := st.SumAmountCentsByKind(ctx, hid, nil, nil, "expense")
	if err != nil {
		t.Fatal(err)
	}
	if inc != wantInc || exp != wantExp {
		t.Fatalf("nil,nil range got inc=%d exp=%d want inc=%d exp=%d", inc, exp, wantInc, wantExp)
	}
}

func TestSumIncomeExpenseCentsInRange_matchesSplitQueries(t *testing.T) {
	t.Parallel()
	st := testStore(t)
	ctx := context.Background()
	hash := passwordtest.MustHash(t, "x")
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

func TestSumIncomeExpenseCentsInRange_fromOnly_toOnly_matchSplitQueries(t *testing.T) {
	t.Parallel()
	st := testStore(t)
	ctx := context.Background()
	hash := passwordtest.MustHash(t, "x")
	uid, err := st.CreateUser(ctx, "agg-one-bound@example.com", hash, "user")
	if err != nil {
		t.Fatal(err)
	}
	u, err := st.GetUserByID(ctx, uid)
	if err != nil || u == nil {
		t.Fatal(err)
	}
	hid := u.HouseholdID
	day := time.Date(2026, 8, 1, 12, 0, 0, 0, time.UTC)
	if _, err := st.CreateTransaction(ctx, uid, hid, 7000, day, "in", nil); err != nil {
		t.Fatal(err)
	}
	if _, err := st.CreateTransaction(ctx, uid, hid, -2000, day, "out", nil); err != nil {
		t.Fatal(err)
	}

	from := day.Add(-time.Hour)
	to := day.Add(time.Hour)

	incFrom, expFrom, err := st.SumIncomeExpenseCentsInRange(ctx, hid, &from, nil)
	if err != nil {
		t.Fatal(err)
	}
	wantIncFrom, err := st.SumAmountCentsByKind(ctx, hid, &from, nil, "income")
	if err != nil {
		t.Fatal(err)
	}
	wantExpFrom, err := st.SumAmountCentsByKind(ctx, hid, &from, nil, "expense")
	if err != nil {
		t.Fatal(err)
	}
	if incFrom != wantIncFrom || expFrom != wantExpFrom {
		t.Fatalf("from-only: got inc=%d exp=%d want inc=%d exp=%d", incFrom, expFrom, wantIncFrom, wantExpFrom)
	}

	incTo, expTo, err := st.SumIncomeExpenseCentsInRange(ctx, hid, nil, &to)
	if err != nil {
		t.Fatal(err)
	}
	wantIncTo, err := st.SumAmountCentsByKind(ctx, hid, nil, &to, "income")
	if err != nil {
		t.Fatal(err)
	}
	wantExpTo, err := st.SumAmountCentsByKind(ctx, hid, nil, &to, "expense")
	if err != nil {
		t.Fatal(err)
	}
	if incTo != wantIncTo || expTo != wantExpTo {
		t.Fatalf("to-only: got inc=%d exp=%d want inc=%d exp=%d", incTo, expTo, wantIncTo, wantExpTo)
	}
}

func TestSumIncomeExpenseCentsInTwoRanges_matchesSeparateQueries(t *testing.T) {
	t.Parallel()
	st := testStore(t)
	ctx := context.Background()
	hash := passwordtest.MustHash(t, "x")
	uid, err := st.CreateUser(ctx, "two-range@example.com", hash, "user")
	if err != nil {
		t.Fatal(err)
	}
	u, err := st.GetUserByID(ctx, uid)
	if err != nil || u == nil {
		t.Fatal(err)
	}
	hid := u.HouseholdID

	// Range A: one income + one expense
	aMid := time.Date(2026, 7, 10, 12, 0, 0, 0, time.UTC)
	if _, err := st.CreateTransaction(ctx, uid, hid, 5000, aMid, "a-in", nil); err != nil {
		t.Fatal(err)
	}
	if _, err := st.CreateTransaction(ctx, uid, hid, -2000, aMid, "a-out", nil); err != nil {
		t.Fatal(err)
	}
	aFrom := aMid.Add(-time.Hour)
	aTo := aMid.Add(time.Hour)

	// Range B: different day, one income only
	bMid := time.Date(2026, 8, 1, 12, 0, 0, 0, time.UTC)
	if _, err := st.CreateTransaction(ctx, uid, hid, 3000, bMid, "b-in", nil); err != nil {
		t.Fatal(err)
	}
	bFrom := bMid.Add(-time.Hour)
	bTo := bMid.Add(time.Hour)

	wantAInc, wantAExp, err := st.SumIncomeExpenseCentsInRange(ctx, hid, &aFrom, &aTo)
	if err != nil {
		t.Fatal(err)
	}
	wantBInc, wantBExp, err := st.SumIncomeExpenseCentsInRange(ctx, hid, &bFrom, &bTo)
	if err != nil {
		t.Fatal(err)
	}

	gotAInc, gotAExp, gotBInc, gotBExp, err := st.SumIncomeExpenseCentsInTwoRanges(ctx, hid, aFrom, aTo, bFrom, bTo)
	if err != nil {
		t.Fatal(err)
	}
	if gotAInc != wantAInc || gotAExp != wantAExp || gotBInc != wantBInc || gotBExp != wantBExp {
		t.Fatalf("two-range %+v %+v %+v %+v want A inc=%d exp=%d B inc=%d exp=%d",
			gotAInc, gotAExp, gotBInc, gotBExp, wantAInc, wantAExp, wantBInc, wantBExp)
	}
}

func TestSumRunningTotalAndIncomeExpenseInTwoRanges_matchesSeparateQueries(t *testing.T) {
	t.Parallel()
	st := testStore(t)
	ctx := context.Background()
	hash := passwordtest.MustHash(t, "x")
	uid, err := st.CreateUser(ctx, "run-two@example.com", hash, "user")
	if err != nil {
		t.Fatal(err)
	}
	u, err := st.GetUserByID(ctx, uid)
	if err != nil || u == nil {
		t.Fatal(err)
	}
	hid := u.HouseholdID

	aMid := time.Date(2026, 7, 10, 12, 0, 0, 0, time.UTC)
	if _, err := st.CreateTransaction(ctx, uid, hid, 5000, aMid, "a-in", nil); err != nil {
		t.Fatal(err)
	}
	if _, err := st.CreateTransaction(ctx, uid, hid, -2000, aMid, "a-out", nil); err != nil {
		t.Fatal(err)
	}
	bMid := time.Date(2026, 8, 1, 12, 0, 0, 0, time.UTC)
	if _, err := st.CreateTransaction(ctx, uid, hid, 3000, bMid, "b-in", nil); err != nil {
		t.Fatal(err)
	}
	old := time.Date(2025, 1, 1, 12, 0, 0, 0, time.UTC)
	if _, err := st.CreateTransaction(ctx, uid, hid, -1000, old, "old", nil); err != nil {
		t.Fatal(err)
	}

	aFrom := aMid.Add(-time.Hour)
	aTo := aMid.Add(time.Hour)
	bFrom := bMid.Add(-time.Hour)
	bTo := bMid.Add(time.Hour)

	wantRun, err := st.SumAmountCents(ctx, hid, nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	wantAInc, wantAExp, err := st.SumIncomeExpenseCentsInRange(ctx, hid, &aFrom, &aTo)
	if err != nil {
		t.Fatal(err)
	}
	wantBInc, wantBExp, err := st.SumIncomeExpenseCentsInRange(ctx, hid, &bFrom, &bTo)
	if err != nil {
		t.Fatal(err)
	}

	gotRun, gotAInc, gotAExp, gotBInc, gotBExp, err := st.SumRunningTotalAndIncomeExpenseInTwoRanges(ctx, hid, aFrom, aTo, bFrom, bTo)
	if err != nil {
		t.Fatal(err)
	}
	if gotRun != wantRun {
		t.Fatalf("running got %d want %d", gotRun, wantRun)
	}
	if gotAInc != wantAInc || gotAExp != wantAExp || gotBInc != wantBInc || gotBExp != wantBExp {
		t.Fatalf("periods %+v %+v %+v %+v want A inc=%d exp=%d B inc=%d exp=%d",
			gotAInc, gotAExp, gotBInc, gotBExp, wantAInc, wantAExp, wantBInc, wantBExp)
	}
}

func TestSumRunningTotalAndIncomeExpenseInTwoRanges_cancelledContext(t *testing.T) {
	t.Parallel()
	st := testStore(t)
	aFrom := time.Date(2026, 7, 1, 0, 0, 0, 0, time.UTC)
	aTo := time.Date(2026, 7, 31, 23, 59, 59, 0, time.UTC)
	bFrom := time.Date(2026, 8, 1, 0, 0, 0, 0, time.UTC)
	bTo := time.Date(2026, 8, 31, 23, 59, 59, 0, time.UTC)
	_, _, _, _, _, err := st.SumRunningTotalAndIncomeExpenseInTwoRanges(alreadyCancelledContext(t), 1, aFrom, aTo, bFrom, bTo)
	assertErrIsContextCanceled(t, err)
}

func TestSumIncomeExpenseCentsInRange_cancelledContext(t *testing.T) {
	t.Parallel()
	st := testStore(t)
	from := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	to := time.Date(2026, 1, 31, 23, 59, 59, 0, time.UTC)
	_, _, err := st.SumIncomeExpenseCentsInRange(alreadyCancelledContext(t), 1, &from, &to)
	assertErrIsContextCanceled(t, err)
}

func TestSumAmountCentsByKind_cancelledContext(t *testing.T) {
	t.Parallel()
	st := testStore(t)
	from := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	to := time.Date(2026, 1, 31, 23, 59, 59, 0, time.UTC)
	_, err := st.SumAmountCentsByKind(alreadyCancelledContext(t), 1, &from, &to, "expense")
	assertErrIsContextCanceled(t, err)
}

func TestSumAmountCentsByKind_nilNil_cancelledContext(t *testing.T) {
	t.Parallel()
	st := testStore(t)
	_, err := st.SumAmountCentsByKind(alreadyCancelledContext(t), 1, nil, nil, "")
	assertErrIsContextCanceled(t, err)
}

func TestSumAmountCentsByKind_trimmedKindMatchesIncome(t *testing.T) {
	t.Parallel()
	st := testStore(t)
	ctx := context.Background()
	hash := passwordtest.MustHash(t, "x")
	uid, err := st.CreateUser(ctx, "sum-trim-kind@example.com", hash, "user")
	if err != nil {
		t.Fatal(err)
	}
	u, err := st.GetUserByID(ctx, uid)
	if err != nil || u == nil {
		t.Fatal(err)
	}
	hid := u.HouseholdID
	day := time.Date(2026, 10, 1, 12, 0, 0, 0, time.UTC)
	if _, err := st.CreateTransaction(ctx, uid, hid, 1000, day, "i", nil); err != nil {
		t.Fatal(err)
	}
	got, err := st.SumAmountCentsByKind(ctx, hid, nil, nil, "  income ")
	if err != nil {
		t.Fatal(err)
	}
	if got != 1000 {
		t.Fatalf("got %d want 1000", got)
	}
}

func TestSumAmountCentsByKind_nilNil_kindsMatchExpected(t *testing.T) {
	t.Parallel()
	st := testStore(t)
	ctx := context.Background()
	hash := passwordtest.MustHash(t, "x")
	uid, err := st.CreateUser(ctx, "sum-kind-nil@example.com", hash, "user")
	if err != nil {
		t.Fatal(err)
	}
	u, err := st.GetUserByID(ctx, uid)
	if err != nil || u == nil {
		t.Fatal(err)
	}
	hid := u.HouseholdID
	day := time.Date(2026, 9, 1, 12, 0, 0, 0, time.UTC)
	if _, err := st.CreateTransaction(ctx, uid, hid, 10000, day, "in", nil); err != nil {
		t.Fatal(err)
	}
	if _, err := st.CreateTransaction(ctx, uid, hid, -3000, day, "out", nil); err != nil {
		t.Fatal(err)
	}
	if _, err := st.CreateTransaction(ctx, uid, hid, 500, day, "small", nil); err != nil {
		t.Fatal(err)
	}

	gotAll, err := st.SumAmountCentsByKind(ctx, hid, nil, nil, "")
	if err != nil {
		t.Fatal(err)
	}
	if gotAll != 7500 {
		t.Fatalf("all kinds sum got %d want 7500", gotAll)
	}
	gotInc, err := st.SumAmountCentsByKind(ctx, hid, nil, nil, "income")
	if err != nil {
		t.Fatal(err)
	}
	if gotInc != 10500 {
		t.Fatalf("income sum got %d want 10500", gotInc)
	}
	gotExp, err := st.SumAmountCentsByKind(ctx, hid, nil, nil, "expense")
	if err != nil {
		t.Fatal(err)
	}
	if gotExp != -3000 {
		t.Fatalf("expense sum got %d want -3000", gotExp)
	}
	gotUnknownKindSameAsAll, err := st.SumAmountCentsByKind(ctx, hid, nil, nil, "not-a-kind")
	if err != nil {
		t.Fatal(err)
	}
	if gotUnknownKindSameAsAll != gotAll {
		t.Fatalf("unknown kind should match all-sum, got %d want %d", gotUnknownKindSameAsAll, gotAll)
	}
}
