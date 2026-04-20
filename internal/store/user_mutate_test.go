package store

import (
	"context"
	"database/sql"
	"errors"
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

func TestUpdateUserProfile_noSuchUser(t *testing.T) {
	t.Parallel()
	st := testStore(t)
	ctx := context.Background()
	err := st.UpdateUserProfile(ctx, 999999999999, "a", "b")
	if !errors.Is(err, sql.ErrNoRows) {
		t.Fatalf("got %v want %v", err, sql.ErrNoRows)
	}
}

func TestUpdateUserPassword_noSuchUser(t *testing.T) {
	t.Parallel()
	st := testStore(t)
	ctx := context.Background()
	err := st.UpdateUserPassword(ctx, 999999999999, []byte("irrelevant"))
	if !errors.Is(err, sql.ErrNoRows) {
		t.Fatalf("got %v want %v", err, sql.ErrNoRows)
	}
}

func TestCreateUser_cancelledContext(t *testing.T) {
	t.Parallel()
	st := testStore(t)
	hash := passwordtest.MustHash(t, "x")
	_, err := st.CreateUser(alreadyCancelledContext(t), "ctx-create-user@example.com", hash, "user")
	assertErrIsContextCanceled(t, err)
}

func TestUpdateUserProfile_cancelledContext(t *testing.T) {
	t.Parallel()
	st := testStore(t)
	ctx := context.Background()
	hash := passwordtest.MustHash(t, "x")
	uid, err := st.CreateUser(ctx, "ctx-prof@example.com", hash, "user")
	if err != nil {
		t.Fatal(err)
	}
	err = st.UpdateUserProfile(alreadyCancelledContext(t), uid, "A", "B")
	assertErrIsContextCanceled(t, err)
}

func TestUpdateUserPassword_cancelledContext(t *testing.T) {
	t.Parallel()
	st := testStore(t)
	ctx := context.Background()
	hash := passwordtest.MustHash(t, "x")
	uid, err := st.CreateUser(ctx, "ctx-pw@example.com", hash, "user")
	if err != nil {
		t.Fatal(err)
	}
	newHash := passwordtest.MustHash(t, "new-secret-long")
	err = st.UpdateUserPassword(alreadyCancelledContext(t), uid, newHash)
	assertErrIsContextCanceled(t, err)
}
