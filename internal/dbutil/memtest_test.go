package dbutil

import (
	"context"
	"testing"

	"moana/internal/auth"
)

func TestMustOpenMemStore(t *testing.T) {
	t.Parallel()
	st := MustOpenMemStore(t)
	ctx := context.Background()
	hash, err := auth.HashPassword("x")
	if err != nil {
		t.Fatal(err)
	}
	_, err = st.CreateUser(ctx, "mem@example.com", hash, "user")
	if err != nil {
		t.Fatal(err)
	}
}
