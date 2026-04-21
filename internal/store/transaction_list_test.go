package store

import (
	"context"
	"fmt"
	"testing"
	"time"

	"moana/internal/passwordtest"
)

func TestListTransactions_cancelledContext(t *testing.T) {
	t.Parallel()
	st := testStore(t)
	_, err := st.ListTransactions(alreadyCancelledContext(t), 1, TransactionFilter{})
	assertErrIsContextCanceled(t, err)
}

func TestListTransactions_recentLimitOnly_cancelledContext(t *testing.T) {
	t.Parallel()
	st := testStore(t)
	_, err := st.ListTransactions(alreadyCancelledContext(t), 1, TransactionFilter{Limit: 3})
	assertErrIsContextCanceled(t, err)
}

func TestListTransactions_kindIncomeRecentLimit_cancelledContext(t *testing.T) {
	t.Parallel()
	st := testStore(t)
	_, err := st.ListTransactions(alreadyCancelledContext(t), 1, TransactionFilter{Kind: "income", Limit: 5})
	assertErrIsContextCanceled(t, err)
}

func TestListTransactions_kindExpenseWithLimit_newestFirst(t *testing.T) {
	t.Parallel()
	st := testStore(t)
	ctx := context.Background()
	hash := passwordtest.MustHash(t, "x")
	uid, err := st.CreateUser(ctx, "kind-exp-limit@example.com", hash, "user")
	if err != nil {
		t.Fatal(err)
	}
	u, err := st.GetUserByID(ctx, uid)
	if err != nil || u == nil {
		t.Fatal(err)
	}
	hid := u.HouseholdID
	base := time.Date(2026, 11, 1, 12, 0, 0, 0, time.UTC)
	if _, err := st.CreateTransaction(ctx, uid, hid, -100, base, "old", nil); err != nil {
		t.Fatal(err)
	}
	if _, err := st.CreateTransaction(ctx, uid, hid, -200, base.Add(time.Hour), "mid", nil); err != nil {
		t.Fatal(err)
	}
	if _, err := st.CreateTransaction(ctx, uid, hid, -300, base.Add(2*time.Hour), "new", nil); err != nil {
		t.Fatal(err)
	}
	txs, err := st.ListTransactions(ctx, hid, TransactionFilter{Kind: "expense", Limit: 2})
	if err != nil {
		t.Fatal(err)
	}
	if len(txs) != 2 {
		t.Fatalf("want 2, got %d %+v", len(txs), txs)
	}
	if txs[0].Description != "new" || txs[1].Description != "mid" {
		t.Fatalf("want newest first: %+v %+v", txs[0], txs[1])
	}
}

func TestStaticListTransactionsQuery(t *testing.T) {
	t.Parallel()
	cases := []struct {
		old    bool
		kind   string
		wantOK bool
	}{
		{false, "", true},
		{false, "income", true},
		{false, "expense", true},
		{false, "bogus", false},
		{true, "", true},
		{true, "income", true},
		{true, "expense", true},
		{true, "bogus", false},
	}
	for _, tc := range cases {
		tc := tc
		t.Run(fmt.Sprintf("old_%v_kind_%q", tc.old, tc.kind), func(t *testing.T) {
			t.Parallel()
			_, ok := staticListTransactionsQuery(tc.old, tc.kind)
			if ok != tc.wantOK {
				t.Fatalf("ok=%v want %v", ok, tc.wantOK)
			}
		})
	}
}

func TestListTransactions_oldestFirstIncomeWithLimit_oldestFirst(t *testing.T) {
	t.Parallel()
	st := testStore(t)
	ctx := context.Background()
	hash := passwordtest.MustHash(t, "x")
	uid, err := st.CreateUser(ctx, "old-inc@example.com", hash, "user")
	if err != nil {
		t.Fatal(err)
	}
	u, err := st.GetUserByID(ctx, uid)
	if err != nil || u == nil {
		t.Fatal(err)
	}
	hid := u.HouseholdID
	base := time.Date(2026, 12, 1, 12, 0, 0, 0, time.UTC)
	if _, err := st.CreateTransaction(ctx, uid, hid, 100, base, "first", nil); err != nil {
		t.Fatal(err)
	}
	if _, err := st.CreateTransaction(ctx, uid, hid, 200, base.Add(time.Hour), "second", nil); err != nil {
		t.Fatal(err)
	}
	if _, err := st.CreateTransaction(ctx, uid, hid, 300, base.Add(2*time.Hour), "third", nil); err != nil {
		t.Fatal(err)
	}
	txs, err := st.ListTransactions(ctx, hid, TransactionFilter{Kind: "income", Limit: 2, OldestFirst: true})
	if err != nil {
		t.Fatal(err)
	}
	if len(txs) != 2 {
		t.Fatalf("want 2 got %d %+v", len(txs), txs)
	}
	if txs[0].Description != "first" || txs[1].Description != "second" {
		t.Fatalf("want oldest first: %+v %+v", txs[0], txs[1])
	}
}

func TestListTransactions_oldestFirstLimitOnly_cancelledContext(t *testing.T) {
	t.Parallel()
	st := testStore(t)
	_, err := st.ListTransactions(alreadyCancelledContext(t), 1, TransactionFilter{Limit: 3, OldestFirst: true})
	assertErrIsContextCanceled(t, err)
}

func TestListTransactions_respectsLimit(t *testing.T) {
	t.Parallel()
	st := testStore(t)
	ctx := context.Background()

	hash := passwordtest.MustHash(t, "pw-limit-test")
	uid, err := st.CreateUser(ctx, "limit-list@example.com", hash, "user")
	if err != nil {
		t.Fatal(err)
	}
	u, err := st.GetUserByEmail(ctx, "limit-list@example.com")
	if err != nil || u == nil {
		t.Fatalf("user: %+v err=%v", u, err)
	}
	hid := u.HouseholdID

	base := time.Date(2026, 4, 1, 12, 0, 0, 0, time.UTC)
	for i := range 5 {
		occ := base.Add(time.Duration(i) * time.Hour)
		if _, err := st.CreateTransaction(ctx, uid, hid, int64(100+i), occ, "tx", nil); err != nil {
			t.Fatal(err)
		}
	}

	txs, err := st.ListTransactions(ctx, hid, TransactionFilter{Limit: 3})
	if err != nil {
		t.Fatal(err)
	}
	if len(txs) != 3 {
		t.Fatalf("want 3 rows, got %d", len(txs))
	}
	// Newest first (default): last inserted hour wins
	if txs[0].AmountCents != 104 || txs[2].AmountCents != 102 {
		t.Fatalf("unexpected order: %+v %+v", txs[0], txs[2])
	}

	txsAll, err := st.ListTransactions(ctx, hid, TransactionFilter{})
	if err != nil {
		t.Fatal(err)
	}
	if len(txsAll) != 5 {
		t.Fatalf("unlimited list: want 5, got %d", len(txsAll))
	}

	txsZero, err := st.ListTransactions(ctx, hid, TransactionFilter{Limit: 0})
	if err != nil {
		t.Fatal(err)
	}
	if len(txsZero) != len(txsAll) {
		t.Fatalf("Limit 0: want %d rows like empty filter, got %d", len(txsAll), len(txsZero))
	}

	txsNeg, err := st.ListTransactions(ctx, hid, TransactionFilter{Limit: -1})
	if err != nil {
		t.Fatal(err)
	}
	if len(txsNeg) != 5 {
		t.Fatalf("negative limit treated as no cap: want 5, got %d", len(txsNeg))
	}
}

func TestListTransactions_searchLiteralPercentDoesNotMatchAllRows(t *testing.T) {
	t.Parallel()
	st := testStore(t)
	ctx := context.Background()
	hash := passwordtest.MustHash(t, "x")
	uid, err := st.CreateUser(ctx, "like-pct@example.com", hash, "user")
	if err != nil {
		t.Fatal(err)
	}
	u, err := st.GetUserByID(ctx, uid)
	if err != nil || u == nil {
		t.Fatal(err)
	}
	hid := u.HouseholdID
	day := time.Date(2026, 5, 1, 12, 0, 0, 0, time.UTC)
	if _, err := st.CreateTransaction(ctx, uid, hid, -100, day, "coffee", nil); err != nil {
		t.Fatal(err)
	}
	if _, err := st.CreateTransaction(ctx, uid, hid, -200, day, "a%b", nil); err != nil {
		t.Fatal(err)
	}
	txs, err := st.ListTransactions(ctx, hid, TransactionFilter{Search: "%"})
	if err != nil {
		t.Fatal(err)
	}
	if len(txs) != 1 || txs[0].Description != "a%b" {
		t.Fatalf("want 1 row with literal %%, got %+v", txs)
	}
}

func TestListTransactions_whitespaceOnlySearchSkipped(t *testing.T) {
	t.Parallel()
	st := testStore(t)
	ctx := context.Background()
	hash := passwordtest.MustHash(t, "x")
	uid, err := st.CreateUser(ctx, "ws-search@example.com", hash, "user")
	if err != nil {
		t.Fatal(err)
	}
	u, err := st.GetUserByID(ctx, uid)
	if err != nil || u == nil {
		t.Fatal(err)
	}
	hid := u.HouseholdID
	day := time.Date(2026, 6, 1, 12, 0, 0, 0, time.UTC)
	if _, err := st.CreateTransaction(ctx, uid, hid, -100, day, "coffee", nil); err != nil {
		t.Fatal(err)
	}
	txs, err := st.ListTransactions(ctx, hid, TransactionFilter{Search: " \t  "})
	if err != nil {
		t.Fatal(err)
	}
	if len(txs) != 1 {
		t.Fatalf("whitespace-only search should not filter, got %d rows", len(txs))
	}
}

func TestListTransactions_combinedFilterOldestFirstWithSearchAndBounds(t *testing.T) {
	t.Parallel()
	st := testStore(t)
	ctx := context.Background()
	hash := passwordtest.MustHash(t, "x")
	uid, err := st.CreateUser(ctx, "combo-list@example.com", hash, "user")
	if err != nil {
		t.Fatal(err)
	}
	u, err := st.GetUserByID(ctx, uid)
	if err != nil || u == nil {
		t.Fatal(err)
	}
	hid := u.HouseholdID

	day1 := time.Date(2026, 8, 10, 15, 0, 0, 0, time.UTC)
	day2 := time.Date(2026, 8, 11, 15, 0, 0, 0, time.UTC)
	from := time.Date(2026, 8, 10, 0, 0, 0, 0, time.UTC)
	to := time.Date(2026, 8, 10, 23, 59, 59, 0, time.UTC)

	if _, err := st.CreateTransaction(ctx, uid, hid, 500, day1, "salary august", nil); err != nil {
		t.Fatal(err)
	}
	if _, err := st.CreateTransaction(ctx, uid, hid, -100, day1, "coffee shop", nil); err != nil {
		t.Fatal(err)
	}
	if _, err := st.CreateTransaction(ctx, uid, hid, 300, day2, "salary next day", nil); err != nil {
		t.Fatal(err)
	}

	txs, err := st.ListTransactions(ctx, hid, TransactionFilter{
		FromUTC:     &from,
		ToUTC:       &to,
		Kind:        "income",
		Search:      "sal",
		OldestFirst: true,
		Limit:       10,
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(txs) != 1 || txs[0].Description != "salary august" || txs[0].AmountCents != 500 {
		t.Fatalf("want 1 income row on day1 matching search, got %+v", txs)
	}
}

func TestListTransactions_kindIgnoresSurroundingSpace(t *testing.T) {
	t.Parallel()
	st := testStore(t)
	ctx := context.Background()
	hash := passwordtest.MustHash(t, "x")
	uid, err := st.CreateUser(ctx, "kind-trim@example.com", hash, "user")
	if err != nil {
		t.Fatal(err)
	}
	u, err := st.GetUserByID(ctx, uid)
	if err != nil || u == nil {
		t.Fatal(err)
	}
	hid := u.HouseholdID
	day := time.Date(2026, 7, 1, 12, 0, 0, 0, time.UTC)
	if _, err := st.CreateTransaction(ctx, uid, hid, 100, day, "in", nil); err != nil {
		t.Fatal(err)
	}
	if _, err := st.CreateTransaction(ctx, uid, hid, -50, day, "out", nil); err != nil {
		t.Fatal(err)
	}
	txs, err := st.ListTransactions(ctx, hid, TransactionFilter{Kind: "  income  "})
	if err != nil {
		t.Fatal(err)
	}
	if len(txs) != 1 || txs[0].AmountCents != 100 {
		t.Fatalf("want 1 income row, got %+v", txs)
	}
}

func TestListTransactions_kindExpenseOnlyNegativeAmounts(t *testing.T) {
	t.Parallel()
	st := testStore(t)
	ctx := context.Background()
	hash := passwordtest.MustHash(t, "x")
	uid, err := st.CreateUser(ctx, "kind-expense@example.com", hash, "user")
	if err != nil {
		t.Fatal(err)
	}
	u, err := st.GetUserByID(ctx, uid)
	if err != nil || u == nil {
		t.Fatal(err)
	}
	hid := u.HouseholdID
	day := time.Date(2026, 9, 1, 12, 0, 0, 0, time.UTC)
	if _, err := st.CreateTransaction(ctx, uid, hid, 100, day, "in", nil); err != nil {
		t.Fatal(err)
	}
	if _, err := st.CreateTransaction(ctx, uid, hid, -50, day, "coffee", nil); err != nil {
		t.Fatal(err)
	}
	txs, err := st.ListTransactions(ctx, hid, TransactionFilter{Kind: "expense"})
	if err != nil {
		t.Fatal(err)
	}
	if len(txs) != 1 || txs[0].AmountCents != -50 || txs[0].Description != "coffee" {
		t.Fatalf("want 1 expense row, got %+v", txs)
	}
}

func TestListTransactions_unknownKindDoesNotFilterBySign(t *testing.T) {
	t.Parallel()
	st := testStore(t)
	ctx := context.Background()
	hash := passwordtest.MustHash(t, "x")
	uid, err := st.CreateUser(ctx, "kind-unknown@example.com", hash, "user")
	if err != nil {
		t.Fatal(err)
	}
	u, err := st.GetUserByID(ctx, uid)
	if err != nil || u == nil {
		t.Fatal(err)
	}
	hid := u.HouseholdID
	day := time.Date(2026, 10, 1, 12, 0, 0, 0, time.UTC)
	if _, err := st.CreateTransaction(ctx, uid, hid, 100, day, "in", nil); err != nil {
		t.Fatal(err)
	}
	if _, err := st.CreateTransaction(ctx, uid, hid, -50, day, "out", nil); err != nil {
		t.Fatal(err)
	}
	txs, err := st.ListTransactions(ctx, hid, TransactionFilter{Kind: "not-a-valid-kind"})
	if err != nil {
		t.Fatal(err)
	}
	if len(txs) != 2 {
		t.Fatalf("unknown kind must not apply income/expense SQL filter: want 2 rows, got %d", len(txs))
	}
}
