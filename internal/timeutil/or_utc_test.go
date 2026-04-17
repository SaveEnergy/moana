package timeutil

import (
	"testing"
	"time"
)

func TestOrUTC(t *testing.T) {
	t.Parallel()
	if OrUTC(nil) != time.UTC {
		t.Fatal("nil must map to UTC")
	}
	loc := time.FixedZone("custom", 3600)
	if got := OrUTC(loc); got != loc {
		t.Fatalf("got %v want %v", got, loc)
	}
}
