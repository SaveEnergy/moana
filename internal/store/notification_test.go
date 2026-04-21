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

func TestMarkNotificationRead_setsReadAt(t *testing.T) {
	t.Parallel()
	st := testStore(t)
	ctx := context.Background()
	hash := passwordtest.MustHash(t, "x")
	uid, err := st.CreateUser(ctx, "notif-read@example.com", hash, "user")
	if err != nil {
		t.Fatal(err)
	}
	nid, err := st.InsertNotification(ctx, uid, "read me")
	if err != nil {
		t.Fatal(err)
	}
	if err := st.MarkNotificationRead(ctx, uid, nid); err != nil {
		t.Fatal(err)
	}
	list, err := st.ListNotificationsForUser(ctx, uid, 10)
	if err != nil {
		t.Fatal(err)
	}
	if len(list) != 1 || list[0].ReadAt == nil {
		t.Fatalf("expected read: %+v", list)
	}
}

func TestMarkNotificationRead_wrongUser(t *testing.T) {
	t.Parallel()
	st := testStore(t)
	ctx := context.Background()
	hash := passwordtest.MustHash(t, "x")
	uidA, err := st.CreateUser(ctx, "notif-a@example.com", hash, "user")
	if err != nil {
		t.Fatal(err)
	}
	uidB, err := st.CreateUser(ctx, "notif-b@example.com", hash, "user")
	if err != nil {
		t.Fatal(err)
	}
	nid, err := st.InsertNotification(ctx, uidA, "private")
	if err != nil {
		t.Fatal(err)
	}
	err = st.MarkNotificationRead(ctx, uidB, nid)
	if !errors.Is(err, ErrNotificationNotFound) {
		t.Fatalf("got %v want %v", err, ErrNotificationNotFound)
	}
}

func TestMarkNotificationRead_cancelledContext(t *testing.T) {
	t.Parallel()
	st := testStore(t)
	err := st.MarkNotificationRead(alreadyCancelledContext(t), 1, 1)
	assertErrIsContextCanceled(t, err)
}
