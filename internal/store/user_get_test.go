package store

import (
	"context"
	"testing"

	"moana/internal/passwordtest"
)

func TestGetUserByID_notFound(t *testing.T) {
	t.Parallel()
	st := testStore(t)
	ctx := context.Background()
	u, err := st.GetUserByID(ctx, 999999999)
	if err != nil {
		t.Fatal(err)
	}
	if u != nil {
		t.Fatalf("want nil, got %+v", u)
	}
}

func TestGetUserByID_found(t *testing.T) {
	t.Parallel()
	st := testStore(t)
	ctx := context.Background()
	hash := passwordtest.MustHash(t, "x")
	uid, err := st.CreateUser(ctx, "gid-found@example.com", hash, "user")
	if err != nil {
		t.Fatal(err)
	}
	u, err := st.GetUserByID(ctx, uid)
	if err != nil {
		t.Fatal(err)
	}
	if u == nil || u.ID != uid || u.Email != "gid-found@example.com" {
		t.Fatalf("user %+v", u)
	}
}
