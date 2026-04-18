package server

import (
	"net/http"
	"testing"
)

func TestHealthPath_matchesGo122ServeMuxPattern(t *testing.T) {
	t.Parallel()
	const legacyGETPattern = "GET /health"
	if got := http.MethodGet + " " + HealthPath; got != legacyGETPattern {
		t.Fatalf("registered pattern %q must stay %q for liveness probes", got, legacyGETPattern)
	}
}
