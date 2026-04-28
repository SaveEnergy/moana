package store

import (
	"strings"
	"testing"
	"time"

	"moana/internal/timeutil"
)

func TestAppendOccurredAtRange(t *testing.T) {
	t.Parallel()
	from := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	to := time.Date(2026, 1, 31, 23, 59, 59, 0, time.UTC)

	t.Run("neither", func(t *testing.T) {
		t.Parallel()
		var b strings.Builder
		b.WriteString("SELECT 1")
		args := appendOccurredAtRange(&b, []any{42}, nil, nil)
		if b.String() != "SELECT 1" {
			t.Fatalf("q %q", b.String())
		}
		if len(args) != 1 || args[0] != 42 {
			t.Fatalf("args %v", args)
		}
	})

	t.Run("from_only", func(t *testing.T) {
		t.Parallel()
		var b strings.Builder
		b.WriteString("SELECT 1")
		args := appendOccurredAtRange(&b, []any{1}, &from, nil)
		want := "SELECT 1" + sqlFilterOccurredAtFrom
		if b.String() != want {
			t.Fatalf("q %q want %q", b.String(), want)
		}
		if len(args) != 2 {
			t.Fatalf("args %v", args)
		}
		if args[1] != timeutil.FormatSQLiteUTC(from) {
			t.Fatalf("bound arg %v", args[1])
		}
	})

	t.Run("to_only", func(t *testing.T) {
		t.Parallel()
		var b strings.Builder
		b.WriteString("SELECT 1")
		args := appendOccurredAtRange(&b, []any{1}, nil, &to)
		want := "SELECT 1" + sqlFilterOccurredAtTo
		if b.String() != want {
			t.Fatalf("q %q want %q", b.String(), want)
		}
		if len(args) != 2 {
			t.Fatalf("args %v", args)
		}
		if args[1] != timeutil.FormatSQLiteUTC(to) {
			t.Fatalf("bound arg %v", args[1])
		}
	})

	t.Run("both", func(t *testing.T) {
		t.Parallel()
		var b strings.Builder
		b.WriteString("SELECT 1")
		args := appendOccurredAtRange(&b, []any{1}, &from, &to)
		want := "SELECT 1" + sqlFilterOccurredAtFrom + sqlFilterOccurredAtTo
		if b.String() != want {
			t.Fatalf("q %q want %q", b.String(), want)
		}
		if len(args) != 3 {
			t.Fatalf("args %v", args)
		}
		if args[1] != timeutil.FormatSQLiteUTC(from) || args[2] != timeutil.FormatSQLiteUTC(to) {
			t.Fatalf("bounds %v %v", args[1], args[2])
		}
	})
}

func TestSQLWithOccurredAtRange(t *testing.T) {
	t.Parallel()
	from := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	to := time.Date(2026, 1, 31, 23, 59, 59, 0, time.UTC)

	cases := []struct {
		name string
		from *time.Time
		to   *time.Time
		want string
	}{
		{
			name: "neither",
			want: "SELECT 1 ORDER BY 1",
		},
		{
			name: "from only",
			from: &from,
			want: "SELECT 1" + sqlFilterOccurredAtFrom + " ORDER BY 1",
		},
		{
			name: "to only",
			to:   &to,
			want: "SELECT 1" + sqlFilterOccurredAtTo + " ORDER BY 1",
		},
		{
			name: "both",
			from: &from,
			to:   &to,
			want: "SELECT 1" + sqlFilterOccurredAtFrom + sqlFilterOccurredAtTo + " ORDER BY 1",
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			if got := sqlWithOccurredAtRange("SELECT 1", tc.from, tc.to, " ORDER BY 1"); got != tc.want {
				t.Fatalf("query %q want %q", got, tc.want)
			}
		})
	}
}

func TestAggregateQueryShapeHelpers(t *testing.T) {
	t.Parallel()
	from := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	to := time.Date(2026, 1, 31, 23, 59, 59, 0, time.UTC)

	if got, want := sqlWithOccurredAtRange(sqlSumAmountAllHousehold, &from, &to, sqlFilterAmountIncome),
		sqlSumAmountAllHousehold+sqlFilterOccurredAtFrom+sqlFilterOccurredAtTo+sqlFilterAmountIncome; got != want {
		t.Fatalf("sum amount query %q want %q", got, want)
	}

	if got, want := sqlWithOccurredAtRange(sqlListTopExpenseCategoriesPrefix, &from, nil, sqlListTopExpenseCategoriesSuffix),
		sqlListTopExpenseCategoriesPrefix+sqlFilterOccurredAtFrom+sqlListTopExpenseCategoriesSuffix; got != want {
		t.Fatalf("top expense query %q want %q", got, want)
	}

	categorySuffix := sqlFilterAmountExpense + sqlListCategoryAmountsGroupBy + sqlListCategoryAmountsOrderExpense
	if got, want := sqlWithOccurredAtRange(sqlListCategoryAmountsSelectPrefix, nil, &to, categorySuffix),
		sqlListCategoryAmountsSelectPrefix+sqlFilterOccurredAtTo+categorySuffix; got != want {
		t.Fatalf("category amounts query %q want %q", got, want)
	}
}
