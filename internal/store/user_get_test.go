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

func TestGetUserByID_cancelledContext(t *testing.T) {
	t.Parallel()
	st := testStore(t)
	_, err := st.GetUserByID(alreadyCancelledContext(t), 1)
	assertErrIsContextCanceled(t, err)
}

func TestGetUserByIDWithUnreadNotificationCount_notFound(t *testing.T) {
	t.Parallel()
	st := testStore(t)
	ctx := context.Background()
	u, unread, err := st.GetUserByIDWithUnreadNotificationCount(ctx, 999999999)
	if err != nil {
		t.Fatal(err)
	}
	if u != nil || unread != 0 {
		t.Fatalf("want nil,0 got u=%v unread=%d", u, unread)
	}
}

func TestGetUserByIDWithUnreadNotificationCount_matchesCountUnread(t *testing.T) {
	t.Parallel()
	st := testStore(t)
	ctx := context.Background()
	hash := passwordtest.MustHash(t, "x")
	uid, err := st.CreateUser(ctx, "unread-combo@example.com", hash, "user")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := st.InsertNotification(ctx, uid, "a"); err != nil {
		t.Fatal(err)
	}
	if _, err := st.InsertNotification(ctx, uid, "b"); err != nil {
		t.Fatal(err)
	}
	n1, err := st.CountUnreadNotificationsForUser(ctx, uid)
	if err != nil || n1 != 2 {
		t.Fatalf("count unread: %d err=%v", n1, err)
	}
	u, unread, err := st.GetUserByIDWithUnreadNotificationCount(ctx, uid)
	if err != nil {
		t.Fatal(err)
	}
	if u == nil || u.ID != uid {
		t.Fatalf("user %+v", u)
	}
	if unread != n1 {
		t.Fatalf("combo unread %d want %d", unread, n1)
	}
}

func TestGetUserByIDWithUnreadNotificationCount_cancelledContext(t *testing.T) {
	t.Parallel()
	st := testStore(t)
	_, _, err := st.GetUserByIDWithUnreadNotificationCount(alreadyCancelledContext(t), 1)
	assertErrIsContextCanceled(t, err)
}

func TestGetUserByEmail_cancelledContext(t *testing.T) {
	t.Parallel()
	st := testStore(t)
	_, err := st.GetUserByEmail(alreadyCancelledContext(t), "cancel-ctx@example.com")
	assertErrIsContextCanceled(t, err)
}
