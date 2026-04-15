package store

import (
	"context"
	"errors"
	"testing"

	"moana/internal/auth"
)

func TestDeleteCategory_notFound(t *testing.T) {
	t.Parallel()
	st := testStore(t)
	ctx := context.Background()

	hash, err := auth.HashPassword("pw-cat-del")
	if err != nil {
		t.Fatal(err)
	}
	_, err = st.CreateUser(ctx, "cat-del-nf@example.com", hash, "user")
	if err != nil {
		t.Fatal(err)
	}
	u, err := st.GetUserByEmail(ctx, "cat-del-nf@example.com")
	if err != nil || u == nil {
		t.Fatal(err)
	}

	err = st.DeleteCategory(ctx, u.HouseholdID, 999999999999)
	if !errors.Is(err, ErrCategoryNotFound) {
		t.Fatalf("got %v want %v", err, ErrCategoryNotFound)
	}
}

func TestUpdateCategory_notFound(t *testing.T) {
	t.Parallel()
	st := testStore(t)
	ctx := context.Background()

	hash, err := auth.HashPassword("pw-cat-up")
	if err != nil {
		t.Fatal(err)
	}
	_, err = st.CreateUser(ctx, "cat-up-nf@example.com", hash, "user")
	if err != nil {
		t.Fatal(err)
	}
	u, err := st.GetUserByEmail(ctx, "cat-up-nf@example.com")
	if err != nil || u == nil {
		t.Fatal(err)
	}

	err = st.UpdateCategory(ctx, u.HouseholdID, 999999999999, "x", "", "")
	if !errors.Is(err, ErrCategoryNotFound) {
		t.Fatalf("got %v want %v", err, ErrCategoryNotFound)
	}
}

func TestCreateCategory_duplicateName(t *testing.T) {
	t.Parallel()
	st := testStore(t)
	ctx := context.Background()

	hash, err := auth.HashPassword("pw-cat-dup")
	if err != nil {
		t.Fatal(err)
	}
	_, err = st.CreateUser(ctx, "cat-dup@example.com", hash, "user")
	if err != nil {
		t.Fatal(err)
	}
	u, err := st.GetUserByEmail(ctx, "cat-dup@example.com")
	if err != nil || u == nil {
		t.Fatal(err)
	}
	if _, err := st.CreateCategory(ctx, u.HouseholdID, "Groceries", "", ""); err != nil {
		t.Fatal(err)
	}
	_, err = st.CreateCategory(ctx, u.HouseholdID, "Groceries", "", "")
	if err == nil {
		t.Fatal("expected ErrDuplicateCategoryName")
	}
	if !errors.Is(err, ErrDuplicateCategoryName) {
		t.Fatalf("got %v want %v", err, ErrDuplicateCategoryName)
	}
}

func TestUpdateCategory_duplicateRename(t *testing.T) {
	t.Parallel()
	st := testStore(t)
	ctx := context.Background()

	hash, err := auth.HashPassword("pw-cat-dup2")
	if err != nil {
		t.Fatal(err)
	}
	_, err = st.CreateUser(ctx, "cat-dup2@example.com", hash, "user")
	if err != nil {
		t.Fatal(err)
	}
	u, err := st.GetUserByEmail(ctx, "cat-dup2@example.com")
	if err != nil || u == nil {
		t.Fatal(err)
	}
	hid := u.HouseholdID
	if _, err := st.CreateCategory(ctx, hid, "Alpha", "", ""); err != nil {
		t.Fatal(err)
	}
	idB, err := st.CreateCategory(ctx, hid, "Beta", "", "")
	if err != nil {
		t.Fatal(err)
	}
	err = st.UpdateCategory(ctx, hid, idB, "Alpha", "", "")
	if err == nil {
		t.Fatal("expected ErrDuplicateCategoryName")
	}
	if !errors.Is(err, ErrDuplicateCategoryName) {
		t.Fatalf("got %v want %v", err, ErrDuplicateCategoryName)
	}
}
