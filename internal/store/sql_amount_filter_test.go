package store

import (
	"fmt"
	"testing"
)

func TestSqlAmountKindFilter(t *testing.T) {
	t.Parallel()
	cases := []struct {
		kind string
		want string
	}{
		{"income", sqlFilterAmountIncome},
		{"expense", sqlFilterAmountExpense},
		{"", ""},
		{"net", ""},
		{" ", ""},
		{"INCOME", ""},
		{" income ", sqlFilterAmountIncome},
		{"  expense\t", sqlFilterAmountExpense},
	}
	for _, tc := range cases {
		t.Run(fmt.Sprintf("kind_%q", tc.kind), func(t *testing.T) {
			t.Parallel()
			if got := sqlAmountKindFilter(tc.kind); got != tc.want {
				t.Fatalf("got %q want %q", got, tc.want)
			}
		})
	}
}
