package historyview

import (
	"net/url"
	"testing"

	"moana/internal/safepath"
)

func TestRoutePath_matchesSafepathDefault(t *testing.T) {
	t.Parallel()
	if RoutePath != safepath.Default {
		t.Fatalf("RoutePath=%q must equal safepath.Default (%q); handlers redirect and history nav must agree",
			RoutePath, safepath.Default)
	}
}

func TestKindAndSortConstants_alignWithParseHistoryURL(t *testing.T) {
	t.Parallel()
	u, err := url.Parse(RoutePath + "?" + QueryKind + "=" + KindIncome + "&" + QuerySort + "=" + SortOldestValue)
	if err != nil {
		t.Fatal(err)
	}
	p := ParseHistoryURL(u)
	if p.kind != KindIncome || p.filterKind != KindIncome || p.sortLabel != SortOldestValue || !p.oldestFirst {
		t.Fatalf("ParseHistoryURL: %+v", p)
	}
}
