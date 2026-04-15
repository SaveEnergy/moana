package app_test

import (
	"testing"

	moanaapp "moana/internal/app"
	"moana/internal/dbutil"
	"moana/internal/testutil"
)

func TestNew_parsesTemplates(t *testing.T) {
	t.Parallel()
	st, db, err := dbutil.OpenStore(":memory:")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	a, err := moanaapp.New(testutil.DefaultTestConfig(), st)
	if err != nil {
		t.Fatal(err)
	}
	if a == nil || a.Render == nil || a.Render.Templates == nil {
		t.Fatal("expected app with render engine")
	}
}

func TestHTTPHandler_wiresRouter(t *testing.T) {
	t.Parallel()
	st, db, err := dbutil.OpenStore(":memory:")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	h, err := moanaapp.HTTPHandler(testutil.DefaultTestConfig(), st)
	if err != nil {
		t.Fatal(err)
	}
	if h == nil {
		t.Fatal("nil handler")
	}
}
