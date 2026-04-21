package store

import (
	"context"
	"errors"
	"testing"

	"moana/internal/passwordtest"
)

func TestInsertNotification_andListNotificationsForUser(t *testing.T) {
	t.Parallel()
	st := testStore(t)
	ctx := context.Background()
	hash := passwordtest.MustHash(t, "x")
	uid, err := st.CreateUser(ctx, "notif-store@example.com", hash, "user")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := st.InsertNotification(ctx, uid, "  hello inbox  "); err != nil {
		t.Fatal(err)
	}
	list, err := st.ListNotificationsForUser(ctx, uid, 10)
	if err != nil {
		t.Fatal(err)
	}
	if len(list) != 1 {
		t.Fatalf("len=%d %+v", len(list), list)
	}
	if list[0].Body != "hello inbox" {
		t.Fatalf("body %q", list[0].Body)
	}
	if list[0].ReadAt != nil {
		t.Fatal("expected unread")
	}
}

func TestInsertNotification_emptyBody(t *testing.T) {
	t.Parallel()
	st := testStore(t)
	ctx := context.Background()
	_, err := st.InsertNotification(ctx, 1, "   ")
	if !errors.Is(err, ErrInvalidNotificationBody) {
		t.Fatalf("got %v want %v", err, ErrInvalidNotificationBody)
	}
}

func TestListNotificationsForUser_cancelledContext(t *testing.T) {
	t.Parallel()
	st := testStore(t)
	_, err := st.ListNotificationsForUser(alreadyCancelledContext(t), 1, 10)
	assertErrIsContextCanceled(t, err)
}
