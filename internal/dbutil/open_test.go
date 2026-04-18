package dbutil

import (
	"context"
	"testing"

	"moana/internal/passwordtest"
)

func TestOpenStore_memory(t *testing.T) {
	t.Parallel()
	st, db, err := OpenStore(":memory:")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	if st == nil {
		t.Fatal("nil store")
	}
	ctx := context.Background()
	hash := passwordtest.MustHash(t, "x")
	if _, err := st.CreateUser(ctx, "openstore@example.com", hash, "user"); err != nil {
		t.Fatal(err)
	}
}
