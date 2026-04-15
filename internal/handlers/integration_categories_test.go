package handlers_test

import (
	"context"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"testing"

	"moana/internal/testutil"
)

func TestCategoryDelete_removesExisting(t *testing.T) {
	t.Parallel()
	app, srv, cleanup := testutil.NewAppServer(t)
	defer cleanup()
	testutil.MustCreateUser(t, app, "catdelok@moana.test", "pw", "user")
	ctx := context.Background()
	u, err := app.Store.GetUserByEmail(ctx, "catdelok@moana.test")
	if err != nil || u == nil {
		t.Fatal(err)
	}
	id, err := app.Store.CreateCategory(ctx, u.HouseholdID, "ToDeleteMe", "", "")
	if err != nil {
		t.Fatal(err)
	}
	client := testutil.NewCookieClient(t)
	testutil.MustLogin(t, client, srv.URL, "catdelok@moana.test", "pw")
	resp, err := client.PostForm(srv.URL+"/categories/delete", url.Values{
		"id": {strconv.FormatInt(id, 10)},
	})
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("status %d want 200 after redirect to categories", resp.StatusCode)
	}
	body, _ := io.ReadAll(resp.Body)
	s := string(body)
	if strings.Contains(s, "ToDeleteMe") {
		t.Fatalf("deleted category name should not appear on page: %s", s[:min(600, len(s))])
	}
}

func TestCategoryDelete_notFoundShowsMessage(t *testing.T) {
	t.Parallel()
	app, srv, cleanup := testutil.NewAppServer(t)
	defer cleanup()
	testutil.MustCreateUser(t, app, "catdel@moana.test", "pw", "user")
	client := testutil.NewCookieClient(t)
	testutil.MustLogin(t, client, srv.URL, "catdel@moana.test", "pw")
	resp, err := client.PostForm(srv.URL+"/categories/delete", url.Values{
		"id": {"99999999999999"},
	})
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("status %d want 200 with error banner", resp.StatusCode)
	}
	body, _ := io.ReadAll(resp.Body)
	s := string(body)
	if !strings.Contains(s, "That category could not be found.") {
		t.Fatalf("expected not-found copy, got: %s", s[:min(500, len(s))])
	}
}

func TestCategoryCreate_duplicateNameShowsMessage(t *testing.T) {
	t.Parallel()
	app, srv, cleanup := testutil.NewAppServer(t)
	defer cleanup()
	testutil.MustCreateUser(t, app, "dupcat@moana.test", "pw", "user")
	client := testutil.NewCookieClient(t)
	testutil.MustLogin(t, client, srv.URL, "dupcat@moana.test", "pw")
	form := url.Values{
		"name":  {"SameName"},
		"icon":  {""},
		"color": {""},
	}
	resp, err := client.PostForm(srv.URL+"/categories", form)
	if err != nil {
		t.Fatal(err)
	}
	_, _ = io.Copy(io.Discard, resp.Body)
	_ = resp.Body.Close()
	// Client follows redirect to GET /categories → 200.
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("first create final status %d want 200", resp.StatusCode)
	}
	resp2, err := client.PostForm(srv.URL+"/categories", form)
	if err != nil {
		t.Fatal(err)
	}
	defer resp2.Body.Close()
	if resp2.StatusCode != http.StatusOK {
		t.Fatalf("second create status %d want 200 with error", resp2.StatusCode)
	}
	body, _ := io.ReadAll(resp2.Body)
	s := string(body)
	if !strings.Contains(s, "A category with that name already exists.") {
		t.Fatalf("expected duplicate copy, got: %s", s[:min(600, len(s))])
	}
}

func TestCategoryUpdate_renamesCategory(t *testing.T) {
	t.Parallel()
	app, srv, cleanup := testutil.NewAppServer(t)
	defer cleanup()
	testutil.MustCreateUser(t, app, "catupd@moana.test", "pw", "user")
	ctx := context.Background()
	u, err := app.Store.GetUserByEmail(ctx, "catupd@moana.test")
	if err != nil || u == nil {
		t.Fatalf("user: %+v err=%v", u, err)
	}
	id, err := app.Store.CreateCategory(ctx, u.HouseholdID, "BeforeName", "", "")
	if err != nil {
		t.Fatal(err)
	}
	client := testutil.NewCookieClient(t)
	testutil.MustLogin(t, client, srv.URL, "catupd@moana.test", "pw")
	resp, err := client.PostForm(srv.URL+"/categories/update", url.Values{
		"id":    {strconv.FormatInt(id, 10)},
		"name":  {"AfterName"},
		"icon":  {""},
		"color": {""},
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
	if !strings.Contains(s, "AfterName") {
		t.Fatalf("expected updated name on page: %s", s[:min(600, len(s))])
	}
	if strings.Contains(s, "BeforeName") {
		t.Fatal("old category name should not appear after rename")
	}
}

func TestCategoryUpdate_invalidIDShowsMessage(t *testing.T) {
	t.Parallel()
	app, srv, cleanup := testutil.NewAppServer(t)
	defer cleanup()
	testutil.MustCreateUser(t, app, "catbadid@moana.test", "pw", "user")
	client := testutil.NewCookieClient(t)
	testutil.MustLogin(t, client, srv.URL, "catbadid@moana.test", "pw")
	resp, err := client.PostForm(srv.URL+"/categories/update", url.Values{
		"id":    {"0"},
		"name":  {"X"},
		"icon":  {""},
		"color": {""},
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
	if !strings.Contains(s, "Invalid category.") {
		t.Fatalf("expected invalid id copy, got: %s", s[:min(600, len(s))])
	}
}
