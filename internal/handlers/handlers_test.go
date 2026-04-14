package handlers_test

import (
	"testing"

	"moana/internal/handlers"
	"moana/internal/testutil"
)

func testApp(t *testing.T) (*handlers.App, func()) {
	t.Helper()
	return testutil.NewApp(t)
}
