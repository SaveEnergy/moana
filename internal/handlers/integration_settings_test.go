package handlers_test

import (
	"context"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"testing"

	"moana/internal/auth"
	"moana/internal/testutil"
)

func TestSettingsProfileUpdate_firstNameShowsSuccess(t *testing.T) {
	t.Parallel()
	app, srv, cleanup := testutil.NewAppServer(t)
	defer cleanup()
	testutil.MustCreateUser(t, app, "profile@integration.test", "pw", "user")
	client := testutil.NewCookieClient(t)
	testutil.MustLogin(t, client, srv.URL, "profile@integration.test", "pw")

	resp, err := client.PostForm(srv.URL+"/settings/profile", url.Values{
		"first_name": {"Pat"},
		"last_name":  {""},
	})
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("status %d want 200 after redirect to settings", resp.StatusCode)
	}
	body, _ := io.ReadAll(resp.Body)
	s := string(body)
	if !strings.Contains(s, "Profile saved.") {
		t.Fatalf("expected success banner, got: %s", s[:min(800, len(s))])
	}
	if !strings.Contains(s, `class="alert alert-success`) {
		t.Fatal("expected success alert class")
	}
}

func TestSettingsHouseholdUpdate_nameShowsSuccess(t *testing.T) {
	t.Parallel()
	app, srv, cleanup := testutil.NewAppServer(t)
	defer cleanup()
	testutil.MustCreateUser(t, app, "hhname@integration.test", "pw", "user")
	client := testutil.NewCookieClient(t)
	testutil.MustLogin(t, client, srv.URL, "hhname@integration.test", "pw")

	resp, err := client.PostForm(srv.URL+"/settings/household", url.Values{
		"household_name": {"The Casa"},
	})
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("status %d want 200 after redirect to settings", resp.StatusCode)
	}
	body, _ := io.ReadAll(resp.Body)
	s := string(body)
	if !strings.Contains(s, "Household name saved.") {
		t.Fatalf("expected household success banner, got: %s", s[:min(800, len(s))])
	}
}

func TestSettingsHouseholdMemberAdd_showsSuccess(t *testing.T) {
	t.Parallel()
	app, srv, cleanup := testutil.NewAppServer(t)
	defer cleanup()
	testutil.MustCreateUser(t, app, "owner-mem@integration.test", "pw", "user")
	client := testutil.NewCookieClient(t)
	testutil.MustLogin(t, client, srv.URL, "owner-mem@integration.test", "pw")

	resp, err := client.PostForm(srv.URL+"/settings/household/members", url.Values{
		"email":    {"newmember@integration.test"},
		"password": {"member-initial-secret-99"},
	})
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("status %d want 200 after redirect to settings", resp.StatusCode)
	}
	body, _ := io.ReadAll(resp.Body)
	s := string(body)
	if !strings.Contains(s, "Member added.") {
		t.Fatalf("expected member success banner, got: %s", s[:min(900, len(s))])
	}
}

func TestSettingsHouseholdMemberRemove_ownerRemovesMember(t *testing.T) {
	t.Parallel()
	app, srv, cleanup := testutil.NewAppServer(t)
	defer cleanup()
	testutil.MustCreateUser(t, app, "rem-owner@integration.test", "pw", "user")
	client := testutil.NewCookieClient(t)
	testutil.MustLogin(t, client, srv.URL, "rem-owner@integration.test", "pw")

	addResp, err := client.PostForm(srv.URL+"/settings/household/members", url.Values{
		"email":    {"rem-member@integration.test"},
		"password": {"rem-member-secret-9"},
	})
	if err != nil {
		t.Fatal(err)
	}
	defer addResp.Body.Close()
	if addResp.StatusCode != http.StatusOK {
		t.Fatalf("add member status %d", addResp.StatusCode)
	}

	ctx := context.Background()
	member, err := app.Store.GetUserByEmail(ctx, "rem-member@integration.test")
	if err != nil || member == nil {
		t.Fatalf("member user: %+v err=%v", member, err)
	}

	rmResp, err := client.PostForm(srv.URL+"/settings/household/members/remove", url.Values{
		"user_id": {strconv.FormatInt(member.ID, 10)},
	})
	if err != nil {
		t.Fatal(err)
	}
	defer rmResp.Body.Close()
	if rmResp.StatusCode != http.StatusOK {
		t.Fatalf("remove status %d", rmResp.StatusCode)
	}
	body, _ := io.ReadAll(rmResp.Body)
	s := string(body)
	if !strings.Contains(s, "Member removed from the household.") {
		t.Fatalf("expected removed banner, got: %s", s[:min(900, len(s))])
	}
}

func TestSettingsHouseholdMemberRemove_invalidUserIDShowsError(t *testing.T) {
	t.Parallel()
	app, srv, cleanup := testutil.NewAppServer(t)
	defer cleanup()
	testutil.MustCreateUser(t, app, "rem-invalid@integration.test", "pw", "user")
	client := testutil.NewCookieClient(t)
	testutil.MustLogin(t, client, srv.URL, "rem-invalid@integration.test", "pw")

	resp, err := client.PostForm(srv.URL+"/settings/household/members/remove", url.Values{
		"user_id": {"not-a-number"},
	})
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("status %d", resp.StatusCode)
	}
	body, _ := io.ReadAll(resp.Body)
	s := string(body)
	if !strings.Contains(s, "Invalid member.") {
		t.Fatalf("expected validation error, got: %s", s[:min(900, len(s))])
	}
	if !strings.Contains(s, `class="alert alert-error`) {
		t.Fatal("expected error alert class")
	}
}

func TestSettingsHouseholdMemberRemove_memberCannotRemoveOwner(t *testing.T) {
	t.Parallel()
	app, srv, cleanup := testutil.NewAppServer(t)
	defer cleanup()
	testutil.MustCreateUser(t, app, "rem-o@integration.test", "pw", "user")
	ctx := context.Background()
	owner, err := app.Store.GetUserByEmail(ctx, "rem-o@integration.test")
	if err != nil || owner == nil {
		t.Fatalf("owner: %+v err=%v", owner, err)
	}
	hash, err := auth.HashPassword("mem-pw")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := app.Store.CreateHouseholdMember(ctx, owner.HouseholdID, "rem-m@integration.test", hash); err != nil {
		t.Fatal(err)
	}

	client := testutil.NewCookieClient(t)
	testutil.MustLogin(t, client, srv.URL, "rem-m@integration.test", "mem-pw")

	resp, err := client.PostForm(srv.URL+"/settings/household/members/remove", url.Values{
		"user_id": {strconv.FormatInt(owner.ID, 10)},
	})
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("status %d", resp.StatusCode)
	}
	body, _ := io.ReadAll(resp.Body)
	s := string(body)
	if !strings.Contains(s, "You cannot remove this member.") {
		t.Fatalf("expected permission error, got: %s", s[:min(900, len(s))])
	}
}

func TestSettingsProfileUpdate_passwordChangeShowsSuccess(t *testing.T) {
	t.Parallel()
	app, srv, cleanup := testutil.NewAppServer(t)
	defer cleanup()
	testutil.MustCreateUser(t, app, "profile-pw@integration.test", "original-secret-1", "user")
	client := testutil.NewCookieClient(t)
	testutil.MustLogin(t, client, srv.URL, "profile-pw@integration.test", "original-secret-1")

	resp, err := client.PostForm(srv.URL+"/settings/profile", url.Values{
		"first_name":           {"Pat"},
		"last_name":            {""},
		"current_password":     {"original-secret-1"},
		"new_password":         {"new-secret-very-long-2"},
		"new_password_confirm": {"new-secret-very-long-2"},
	})
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("status %d want 200 after redirect to settings", resp.StatusCode)
	}
	body, _ := io.ReadAll(resp.Body)
	s := string(body)
	if !strings.Contains(s, "Password updated.") {
		t.Fatalf("expected password success banner, got: %s", s[:min(900, len(s))])
	}
	if !strings.Contains(s, `class="alert alert-success`) {
		t.Fatal("expected success alert class")
	}

	client2 := testutil.NewCookieClient(t)
	testutil.MustLogin(t, client2, srv.URL, "profile-pw@integration.test", "new-secret-very-long-2")
}

func TestSettingsProfileUpdate_wrongCurrentPasswordShowsError(t *testing.T) {
	t.Parallel()
	app, srv, cleanup := testutil.NewAppServer(t)
	defer cleanup()
	testutil.MustCreateUser(t, app, "profile-badcur@integration.test", "good-secret-1", "user")
	client := testutil.NewCookieClient(t)
	testutil.MustLogin(t, client, srv.URL, "profile-badcur@integration.test", "good-secret-1")

	resp, err := client.PostForm(srv.URL+"/settings/profile", url.Values{
		"first_name":           {""},
		"last_name":            {""},
		"current_password":     {"wrong-password"},
		"new_password":         {"new-secret-very-long-2"},
		"new_password_confirm": {"new-secret-very-long-2"},
	})
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("status %d", resp.StatusCode)
	}
	body, _ := io.ReadAll(resp.Body)
	s := string(body)
	if !strings.Contains(s, "Current password is incorrect.") {
		t.Fatalf("expected wrong-current copy, got: %s", s[:min(900, len(s))])
	}
}

func TestSettingsProfileUpdate_newPasswordMismatchShowsError(t *testing.T) {
	t.Parallel()
	app, srv, cleanup := testutil.NewAppServer(t)
	defer cleanup()
	testutil.MustCreateUser(t, app, "profile-mismatch@integration.test", "good-secret-1", "user")
	client := testutil.NewCookieClient(t)
	testutil.MustLogin(t, client, srv.URL, "profile-mismatch@integration.test", "good-secret-1")

	resp, err := client.PostForm(srv.URL+"/settings/profile", url.Values{
		"first_name":           {""},
		"last_name":            {""},
		"current_password":     {"good-secret-1"},
		"new_password":         {"new-secret-very-long-2"},
		"new_password_confirm": {"other-secret-3"},
	})
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("status %d", resp.StatusCode)
	}
	body, _ := io.ReadAll(resp.Body)
	s := string(body)
	if !strings.Contains(s, "New passwords do not match.") {
		t.Fatalf("expected mismatch copy, got: %s", s[:min(900, len(s))])
	}
}
