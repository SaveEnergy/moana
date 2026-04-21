package store

import "testing"

func TestSqlAmountKindFilter(t *testing.T) {
	t.Parallel()
	if got, want := sqlAmountKindFilter("income"), sqlFilterAmountIncome; got != want {
		t.Fatalf("income: got %q want %q", got, want)
	}
	if got, want := sqlAmountKindFilter("expense"), sqlFilterAmountExpense; got != want {
		t.Fatalf("expense: got %q want %q", got, want)
	}
	for _, kind := range []string{"", "net", " ", "INCOME"} {
		if sqlAmountKindFilter(kind) != "" {
			t.Fatalf("kind %q: want empty fragment", kind)
		}
	}
}
