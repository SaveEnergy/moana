package store

import (
	"context"
	"errors"
	"testing"

	"moana/internal/auth"
)

func TestGetUserByEmail_ignoresWhitespaceInQuery(t *testing.T) {
	t.Parallel()
	st := testStore(t)
	ctx := context.Background()
	hash, err := auth.HashPassword("x")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := st.CreateUser(ctx, "findme@example.com", hash, "user"); err != nil {
		t.Fatal(err)
	}
	u, err := st.GetUserByEmail(ctx, "  findme@example.com  ")
	if err != nil {
		t.Fatal(err)
	}
	if u == nil || u.Email != "findme@example.com" {
		t.Fatalf("got %+v", u)
	}
}

func TestGetUserByEmail_emptyAfterTrimReturnsNil(t *testing.T) {
	t.Parallel()
	st := testStore(t)
	ctx := context.Background()
	u, err := st.GetUserByEmail(ctx, "   ")
	if err != nil {
		t.Fatal(err)
	}
	if u != nil {
		t.Fatalf("want nil user, got %+v", u)
	}
}

func TestCreateUser_duplicateAfterTrimmingWhitespace(t *testing.T) {
	t.Parallel()
	st := testStore(t)
	ctx := context.Background()
	hash, err := auth.HashPassword("x")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := st.CreateUser(ctx, "  duptrim@example.com  ", hash, "user"); err != nil {
		t.Fatal(err)
	}
	_, err = st.CreateUser(ctx, "duptrim@example.com", hash, "user")
	if !errors.Is(err, ErrDuplicateUserEmail) {
		t.Fatalf("got %v want %v", err, ErrDuplicateUserEmail)
	}
}

func TestCreateUser_whitespaceOnlyEmail(t *testing.T) {
	t.Parallel()
	st := testStore(t)
	ctx := context.Background()
	hash, err := auth.HashPassword("x")
	if err != nil {
		t.Fatal(err)
	}
	_, err = st.CreateUser(ctx, "  \t  ", hash, "user")
	if !errors.Is(err, ErrInvalidUserEmail) {
		t.Fatalf("got %v want %v", err, ErrInvalidUserEmail)
	}
}

func TestCreateHouseholdMember_whitespaceOnlyEmail(t *testing.T) {
	t.Parallel()
	st := testStore(t)
	ctx := context.Background()
	hash, err := auth.HashPassword("x")
	if err != nil {
		t.Fatal(err)
	}
	uid, err := st.CreateUser(ctx, "owner-ws@example.com", hash, "user")
	if err != nil {
		t.Fatal(err)
	}
	owner, err := st.GetUserByID(ctx, uid)
	if err != nil || owner == nil {
		t.Fatal(err)
	}
	_, err = st.CreateHouseholdMember(ctx, owner.HouseholdID, "  ", hash)
	if !errors.Is(err, ErrInvalidUserEmail) {
		t.Fatalf("got %v want %v", err, ErrInvalidUserEmail)
	}
}

func TestCreateHouseholdMember_trimsStoredEmail(t *testing.T) {
	t.Parallel()
	st := testStore(t)
	ctx := context.Background()
	hash, err := auth.HashPassword("x")
	if err != nil {
		t.Fatal(err)
	}
	uid, err := st.CreateUser(ctx, "owner-trim@example.com", hash, "user")
	if err != nil {
		t.Fatal(err)
	}
	owner, err := st.GetUserByID(ctx, uid)
	if err != nil || owner == nil {
		t.Fatal(err)
	}
	hash2, err := auth.HashPassword("y")
	if err != nil {
		t.Fatal(err)
	}
	mid, err := st.CreateHouseholdMember(ctx, owner.HouseholdID, "  member-trim@example.com  ", hash2)
	if err != nil {
		t.Fatal(err)
	}
	m, err := st.GetUserByID(ctx, mid)
	if err != nil || m == nil || m.Email != "member-trim@example.com" {
		t.Fatalf("member %+v err=%v", m, err)
	}
}
