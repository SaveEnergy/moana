package historyview

import (
	"net/url"
	"strconv"
	"testing"
	"time"
)

func TestParseHistoryPageSize(t *testing.T) {
	t.Parallel()
	t.Run("empty defaults to 100", func(t *testing.T) {
		t.Parallel()
		if g := parseHistoryPageSize(url.Values{}); g != defaultHistoryPageSize {
			t.Fatalf("got %d want %d", g, defaultHistoryPageSize)
		}
	})
	t.Run("allowed sizes", func(t *testing.T) {
		t.Parallel()
		for _, want := range DefaultHistoryPageSizes {
			v := url.Values{}
			v.Set(QueryRows, strconv.Itoa(want))
			if g := parseHistoryPageSize(v); g != want {
				t.Fatalf("n=%d got %d want %d", want, g, want)
			}
		}
	})
	t.Run("invalid string defaults", func(t *testing.T) {
		t.Parallel()
		v := url.Values{}
		v.Set(QueryRows, "xx")
		if g := parseHistoryPageSize(v); g != defaultHistoryPageSize {
			t.Fatalf("got %d", g)
		}
	})
	t.Run("unknown positive below max snaps to default", func(t *testing.T) {
		t.Parallel()
		v := url.Values{}
		v.Set(QueryRows, "99")
		if g := parseHistoryPageSize(v); g != defaultHistoryPageSize {
			t.Fatalf("got %d want %d", g, defaultHistoryPageSize)
		}
	})
	t.Run("over hard max snaps to 2000", func(t *testing.T) {
		t.Parallel()
		v := url.Values{}
		v.Set(QueryRows, "99999")
		if g := parseHistoryPageSize(v); g != hardMaxHistoryPageSize {
			t.Fatalf("got %d want %d", g, hardMaxHistoryPageSize)
		}
	})
}

func TestBuildQuickRangeLinks(t *testing.T) {
	t.Parallel()
	loc, err := time.LoadLocation("Europe/Berlin")
	if err != nil {
		t.Fatal(err)
	}
	v := url.Values{}
	v.Set(QueryKind, KindExpense)
	v.Set(QueryRows, "250")
	v.Set(QuerySort, SortOldestValue)
	u := buildQuickRangeLinks(v, loc, 7)
	if u == "" || u[0] != '/' {
		t.Fatalf("unexpected %q", u)
	}
	parsed, err := url.Parse(u)
	if err != nil {
		t.Fatal(err)
	}
	q := parsed.Query()
	if q.Get(QueryKind) != KindExpense {
		t.Fatalf("kind=%q", q.Get(QueryKind))
	}
	if q.Get(QueryRows) != "250" {
		t.Fatalf("n=%q", q.Get(QueryRows))
	}
	if q.Get(QuerySort) != SortOldestValue {
		t.Fatalf("sort=%q", q.Get(QuerySort))
	}
	if q.Get(QueryFrom) == "" || q.Get(QueryTo) == "" {
		t.Fatalf("from=%q to=%q", q.Get(QueryFrom), q.Get(QueryTo))
	}
}
