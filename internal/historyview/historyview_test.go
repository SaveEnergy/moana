package historyview

import (
	"net/url"
	"testing"
	"time"
)

func TestBuildNav(t *testing.T) {
	t.Parallel()
	u, _ := url.Parse("/history?kind=expense&q=a")
	nav := BuildNav(u)
	if nav.LinkAll == "" || nav.SortOldest == "" {
		t.Fatal(nav)
	}
}

func TestFormatDayLabel_today(t *testing.T) {
	t.Parallel()
	loc := time.Local
	now := time.Now().In(loc)
	d0 := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, loc)
	s := FormatDayLabel(d0, loc)
	if len(s) < 8 {
		t.Fatal(s)
	}
}
