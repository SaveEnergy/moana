package store

import (
	"context"
	"testing"

	"moana/internal/auth"
	"moana/internal/passwordtest"
)

func TestUpdateUserProfile_roundTrip(t *testing.T) {
	t.Parallel()
	st := testStore(t)
	ctx := context.Background()
	hash := passwordtest.MustHash(t, "x")
	uid, err := st.CreateUser(ctx, "prof-up@example.com", hash, "user")
	if err != nil {
		t.Fatal(err)
	}
	if err := st.UpdateUserProfile(ctx, uid, "Pat", "Doe"); err != nil {
		t.Fatal(err)
	}
	u, err := st.GetUserByID(ctx, uid)
	if err != nil || u == nil {
		t.Fatal(err)
	}
	if u.FirstName != "Pat" || u.LastName != "Doe" {
		t.Fatalf("names %+q %+q", u.FirstName, u.LastName)
	}
}

func TestUpdateUserPassword_roundTrip(t *testing.T) {
	t.Parallel()
	st := testStore(t)
	ctx := context.Background()
	hash := passwordtest.MustHash(t, "old-secret")
	uid, err := st.CreateUser(ctx, "pw-up@example.com", hash, "user")
	if err != nil {
		t.Fatal(err)
	}
	newHash := passwordtest.MustHash(t, "new-secret-long")
	if err := st.UpdateUserPassword(ctx, uid, newHash); err != nil {
		t.Fatal(err)
	}
	u, err := st.GetUserByID(ctx, uid)
	if err != nil || u == nil {
		t.Fatal(err)
	}
	if err := auth.CheckPassword(u.PasswordHash, "new-secret-long"); err != nil {
		t.Fatalf("check password: %v", err)
	}
}
