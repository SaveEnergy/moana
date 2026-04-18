package historyview

import (
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
