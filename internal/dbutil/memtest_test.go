package dbutil

import (
	"context"
	"testing"

	"moana/internal/passwordtest"
)

func TestMustOpenMemStore(t *testing.T) {
	t.Parallel()
	st := MustOpenMemStore(t)
	ctx := context.Background()
	hash := passwordtest.MustHash(t, "x")
	_, err := st.CreateUser(ctx, "mem@example.com", hash, "user")
	if err != nil {
		t.Fatal(err)
	}
}
