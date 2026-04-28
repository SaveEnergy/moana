package store

import (
	"reflect"
	"testing"
	"time"
)

func TestTransactionListQueryPlanSQL(t *testing.T) {
	t.Parallel()
	from := time.Date(2026, 4, 1, 12, 0, 0, 0, time.UTC)
	to := time.Date(2026, 4, 2, 12, 0, 0, 0, time.UTC)

	cases := []struct {
		name string
		plan transactionListQueryPlan
		want string
	}{
		{
			name: "default newest first no limit",
			plan: newTransactionListQueryPlan(nil, nil, "", "", false, 0),
			want: sqlTransactionListFromHousehold + sqlTransactionListOrderDesc,
		},
		{
			name: "oldest first income with limit",
			plan: newTransactionListQueryPlan(nil, nil, " income ", "", true, 10),
			want: sqlTransactionListFromHousehold + sqlFilterAmountIncome + sqlTransactionListOrderAsc + sqlTransactionListLimitSuffix,
		},
		{
			name: "both bounds search expense newest first limit",
			plan: newTransactionListQueryPlan(&from, &to, "expense", "coffee", false, 50),
			want: sqlTransactionListFromHousehold + sqlFilterOccurredAtFrom + sqlFilterOccurredAtTo + sqlFilterAmountExpense + sqlTransactionListSearchLike + sqlTransactionListOrderDesc + sqlTransactionListLimitSuffix,
		},
		{
			name: "from bound search unknown kind oldest first no limit",
			plan: newTransactionListQueryPlan(&from, nil, "bogus", "rent", true, 0),
			want: sqlTransactionListFromHousehold + sqlFilterOccurredAtFrom + sqlTransactionListSearchLike + sqlTransactionListOrderAsc,
		},
		{
			name: "to bound whitespace search is ignored",
			plan: newTransactionListQueryPlan(nil, &to, "expense", " \t ", false, 5),
			want: sqlTransactionListFromHousehold + sqlFilterOccurredAtTo + sqlFilterAmountExpense + sqlTransactionListOrderDesc + sqlTransactionListLimitSuffix,
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			if got := tc.plan.SQL(); got != tc.want {
				t.Fatalf("query drift:\ngot  %q\nwant %q", got, tc.want)
			}
		})
	}
}

func TestTransactionListQueryArgs(t *testing.T) {
	t.Parallel()
	from := time.Date(2026, 4, 1, 12, 30, 0, 0, time.FixedZone("offset", 2*60*60))
	to := time.Date(2026, 4, 2, 12, 30, 0, 0, time.UTC)

	got := transactionListQueryArgs(42, &from, &to, "50% off_!", 25)
	want := []any{
		int64(42),
		"2026-04-01T10:30:00Z",
		"2026-04-02T12:30:00Z",
		"%50!% off!_!!%",
		"%50!% off!_!!%",
		25,
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("args mismatch:\ngot  %#v\nwant %#v", got, want)
	}
}
