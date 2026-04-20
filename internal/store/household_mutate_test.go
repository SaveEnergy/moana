package store

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"moana/internal/passwordtest"
)

func TestUpdateHouseholdName(t *testing.T) {
	t.Parallel()
	st := testStore(t)
	ctx := context.Background()
	hash := passwordtest.MustHash(t, "x")
	uid, err := st.CreateUser(ctx, "uhname@example.com", hash, "user")
	if err != nil {
		t.Fatal(err)
	}
	u, err := st.GetUserByID(ctx, uid)
	if err != nil || u == nil {
		t.Fatal(err)
	}
	hid := u.HouseholdID

	if err := st.UpdateHouseholdName(ctx, hid, "Renamed Household"); err != nil {
		t.Fatal(err)
	}
	h, err := st.GetHousehold(ctx, hid)
	if err != nil {
		t.Fatal(err)
	}
	if h == nil || h.Name != "Renamed Household" {
		t.Fatalf("household %+v", h)
	}
}

func TestUpdateHouseholdName_noSuchHousehold(t *testing.T) {
	t.Parallel()
	st := testStore(t)
	ctx := context.Background()
	err := st.UpdateHouseholdName(ctx, 999999999999, "x")
	if !errors.Is(err, sql.ErrNoRows) {
		t.Fatalf("got %v want %v", err, sql.ErrNoRows)
	}
}

func TestCreateHouseholdMember_unknownHouseholdID(t *testing.T) {
	t.Parallel()
	st := testStore(t)
	ctx := context.Background()
	hash := passwordtest.MustHash(t, "x")
	_, err := st.CreateHouseholdMember(ctx, 999999999999, "orphan@example.com", hash)
	if err == nil {
		t.Fatal("expected error when household_id does not exist (FK)")
	}
}

func TestDetachUserToSoloHousehold_noSuchUser(t *testing.T) {
	t.Parallel()
	st := testStore(t)
	ctx := context.Background()
	var before int64
	if err := st.DB.QueryRowContext(ctx, `SELECT COUNT(*) FROM households`).Scan(&before); err != nil {
		t.Fatal(err)
	}
	err := st.DetachUserToSoloHousehold(ctx, 999999999999)
	if !errors.Is(err, sql.ErrNoRows) {
		t.Fatalf("got %v want %v", err, sql.ErrNoRows)
	}
	var after int64
	if err := st.DB.QueryRowContext(ctx, `SELECT COUNT(*) FROM households`).Scan(&after); err != nil {
		t.Fatal(err)
	}
	if after != before {
		t.Fatalf("household count %d want %d (failed detach must roll back new household row)", after, before)
	}
}

func TestDetachUserToSoloHousehold(t *testing.T) {
	t.Parallel()
	st := testStore(t)
	ctx := context.Background()
	hash := passwordtest.MustHash(t, "pw1")
	ownerID, err := st.CreateUser(ctx, "detach-owner@example.com", hash, "user")
	if err != nil {
		t.Fatal(err)
	}
	owner, err := st.GetUserByID(ctx, ownerID)
	if err != nil || owner == nil {
		t.Fatal(err)
	}
	hid := owner.HouseholdID

	hash2 := passwordtest.MustHash(t, "pw2")
	memberID, err := st.CreateHouseholdMember(ctx, hid, "detach-member@example.com", hash2)
	if err != nil {
		t.Fatal(err)
	}

	if err := st.DetachUserToSoloHousehold(ctx, memberID); err != nil {
		t.Fatal(err)
	}
	member, err := st.GetUserByID(ctx, memberID)
	if err != nil || member == nil {
		t.Fatal(err)
	}
	if member.HouseholdID == hid {
		t.Fatalf("user still in old household %d", hid)
	}
	if member.HouseholdRole != "owner" {
		t.Fatalf("HouseholdRole %q want owner", member.HouseholdRole)
	}

	nOld, err := st.CountHouseholdMembers(ctx, hid)
	if err != nil {
		t.Fatal(err)
	}
	if nOld != 1 {
		t.Fatalf("old household members %d want 1", nOld)
	}
	nSolo, err := st.CountHouseholdMembers(ctx, member.HouseholdID)
	if err != nil {
		t.Fatal(err)
	}
	if nSolo != 1 {
		t.Fatalf("new household members %d want 1", nSolo)
	}
}
