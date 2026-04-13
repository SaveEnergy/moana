package handlers

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	"moana/internal/auth"
	"moana/internal/config"
	"moana/internal/db"
	"moana/internal/store"
)

func testApp(t *testing.T) (*App, func()) {
	t.Helper()
	database, err := db.Open(":memory:")
	if err != nil {
		t.Fatal(err)
	}
	st := store.New(database)
	tmpl, err := ParseTemplates()
	if err != nil {
		database.Close()
		t.Fatal(err)
	}
	cfg := &config.Config{
		Listen:         ":0",
		DBPath:         ":memory:",
		SessionSecret:  []byte("integration-test-session-secret-32b!"),
		SecureCookies:  false,
		SessionMaxAge:  time.Hour,
		RequestTimeout: 30 * time.Second,
	}
	app := &App{Config: cfg, Store: st, Templates: tmpl}
	cleanup := func() { database.Close() }
	return app, cleanup
}

func TestHealth(t *testing.T) {
	t.Parallel()
	app, cleanup := testApp(t)
	defer cleanup()
	srv := httptest.NewServer(NewRouterWithRouterOptions(&RouterOptions{DisableRequestLogging: true}, app))
	defer srv.Close()
	resp, err := http.Get(srv.URL + "/health")
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("status %d", resp.StatusCode)
	}
	b, _ := io.ReadAll(resp.Body)
	if string(b) != "ok" {
		t.Fatalf("body %q", b)
	}
}

func TestLoginAndOverview(t *testing.T) {
	t.Parallel()
	app, cleanup := testApp(t)
	defer cleanup()
	ctx := context.Background()
	hash, err := auth.HashPassword("correct-password")
	if err != nil {
		t.Fatal(err)
	}
	_, err = app.Store.CreateUser(ctx, "user@integration.test", hash, "user", "UTC")
	if err != nil {
		t.Fatal(err)
	}

	srv := httptest.NewServer(NewRouterWithRouterOptions(&RouterOptions{DisableRequestLogging: true}, app))
	defer srv.Close()

	jar, err := cookiejar.New(nil)
	if err != nil {
		t.Fatal(err)
	}
	client := &http.Client{Jar: jar}

	resp, err := client.PostForm(srv.URL+"/login", url.Values{
		"email":    {"user@integration.test"},
		"password": {"correct-password"},
	})
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("login status %d", resp.StatusCode)
	}
	body, _ := io.ReadAll(resp.Body)
	if !strings.Contains(string(body), "Current portfolio") {
		snippet := string(body)
		if len(snippet) > 300 {
			snippet = snippet[:300] + "…"
		}
		t.Fatalf("expected overview content, got: %s", snippet)
	}
}

func TestAdminUsersForbiddenForNonAdmin(t *testing.T) {
	t.Parallel()
	app, cleanup := testApp(t)
	defer cleanup()
	ctx := context.Background()
	hash, err := auth.HashPassword("pw")
	if err != nil {
		t.Fatal(err)
	}
	_, err = app.Store.CreateUser(ctx, "plain@moana.test", hash, "user", "UTC")
	if err != nil {
		t.Fatal(err)
	}
	srv := httptest.NewServer(NewRouterWithRouterOptions(&RouterOptions{DisableRequestLogging: true}, app))
	defer srv.Close()
	jar, _ := cookiejar.New(nil)
	client := &http.Client{Jar: jar}
	_, err = client.PostForm(srv.URL+"/login", url.Values{
		"email":    {"plain@moana.test"},
		"password": {"pw"},
	})
	if err != nil {
		t.Fatal(err)
	}
	resp, err := client.Get(srv.URL + "/admin/users")
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusForbidden {
		t.Fatalf("status %d", resp.StatusCode)
	}
}

func TestAdminUsersOKForAdmin(t *testing.T) {
	t.Parallel()
	app, cleanup := testApp(t)
	defer cleanup()
	ctx := context.Background()
	hash, err := auth.HashPassword("pw")
	if err != nil {
		t.Fatal(err)
	}
	_, err = app.Store.CreateUser(ctx, "root@moana.test", hash, "admin", "UTC")
	if err != nil {
		t.Fatal(err)
	}
	srv := httptest.NewServer(NewRouterWithRouterOptions(&RouterOptions{DisableRequestLogging: true}, app))
	defer srv.Close()
	jar, _ := cookiejar.New(nil)
	client := &http.Client{Jar: jar}
	_, err = client.PostForm(srv.URL+"/login", url.Values{
		"email":    {"root@moana.test"},
		"password": {"pw"},
	})
	if err != nil {
		t.Fatal(err)
	}
	resp, err := client.Get(srv.URL + "/admin/users")
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("status %d", resp.StatusCode)
	}
	body, _ := io.ReadAll(resp.Body)
	if !strings.Contains(string(body), "User management") {
		t.Fatalf("expected admin page")
	}
}

func TestCreateExpenseStoresNegativeCents(t *testing.T) {
	t.Parallel()
	app, cleanup := testApp(t)
	defer cleanup()
	ctx := context.Background()
	hash, err := auth.HashPassword("pw")
	if err != nil {
		t.Fatal(err)
	}
	_, err = app.Store.CreateUser(ctx, "tx@moana.test", hash, "user", "UTC")
	if err != nil {
		t.Fatal(err)
	}
	srv := httptest.NewServer(NewRouterWithRouterOptions(&RouterOptions{DisableRequestLogging: true}, app))
	defer srv.Close()
	jar, _ := cookiejar.New(nil)
	client := &http.Client{Jar: jar}
	_, err = client.PostForm(srv.URL+"/login", url.Values{
		"email":    {"tx@moana.test"},
		"password": {"pw"},
	})
	if err != nil {
		t.Fatal(err)
	}
	day := time.Now().UTC().Format("2006-01-02")
	resp, err := client.PostForm(srv.URL+"/transactions", url.Values{
		"amount":      {"25.50"},
		"kind":        {"expense"},
		"occurred_on": {day},
		"description": {"test expense"},
	})
	if err != nil {
		t.Fatal(err)
	}
	resp.Body.Close()
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusSeeOther {
		t.Fatalf("post status %d", resp.StatusCode)
	}
	resp2, err := client.Get(srv.URL + "/history")
	if err != nil {
		t.Fatal(err)
	}
	defer resp2.Body.Close()
	body, _ := io.ReadAll(resp2.Body)
	s := string(body)
	if !strings.Contains(s, "-€25.50") {
		t.Fatalf("expected negative EUR in history: %s", s[:min(500, len(s))])
	}
	if !strings.Contains(s, "test expense") {
		t.Fatal("expected transaction description on history")
	}
}

func TestEditTransaction(t *testing.T) {
	t.Parallel()
	app, cleanup := testApp(t)
	defer cleanup()
	ctx := context.Background()
	hash, err := auth.HashPassword("pw")
	if err != nil {
		t.Fatal(err)
	}
	uid, err := app.Store.CreateUser(ctx, "edit@moana.test", hash, "user", "UTC")
	if err != nil {
		t.Fatal(err)
	}
	srv := httptest.NewServer(NewRouterWithRouterOptions(&RouterOptions{DisableRequestLogging: true}, app))
	defer srv.Close()
	jar, _ := cookiejar.New(nil)
	client := &http.Client{Jar: jar}
	_, err = client.PostForm(srv.URL+"/login", url.Values{
		"email":    {"edit@moana.test"},
		"password": {"pw"},
	})
	if err != nil {
		t.Fatal(err)
	}
	day := time.Now().UTC().Format("2006-01-02")
	resp, err := client.PostForm(srv.URL+"/transactions", url.Values{
		"amount":        {"10.00"},
		"kind":          {"expense"},
		"occurred_on":   {day},
		"description":   {"coffee"},
		"category_id":   {""},
	})
	if err != nil {
		t.Fatal(err)
	}
	resp.Body.Close()
	txs, err := app.Store.ListTransactions(ctx, uid, store.TransactionFilter{Limit: 1})
	if err != nil || len(txs) != 1 {
		t.Fatalf("list: %v", err)
	}
	id := txs[0].ID
	editURL := fmt.Sprintf("%s/transactions/%d/edit?next=/history", srv.URL, id)
	resp, err = client.Get(editURL)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("edit get status %d: %s", resp.StatusCode, string(body)[:min(200, len(body))])
	}
	if !strings.Contains(string(body), "Edit entry") || !strings.Contains(string(body), "10.00") {
		t.Fatal("expected edit form")
	}
	resp2, err := client.PostForm(fmt.Sprintf("%s/transactions/%d", srv.URL, id), url.Values{
		"next":          {"/history"},
		"amount":        {"20.00"},
		"kind":          {"expense"},
		"occurred_on":   {day},
		"description":   {"coffee fixed"},
		"category_id":   {""},
	})
	if err != nil {
		t.Fatal(err)
	}
	defer resp2.Body.Close()
	if resp2.StatusCode != http.StatusOK {
		t.Fatalf("update status %d (expected 200 after redirect to history)", resp2.StatusCode)
	}
	tx, err := app.Store.GetTransactionByID(ctx, uid, id)
	if err != nil || tx == nil {
		t.Fatal(err)
	}
	if tx.AmountCents != -2000 {
		t.Fatalf("amount %d", tx.AmountCents)
	}
	if tx.Description != "coffee fixed" {
		t.Fatalf("desc %q", tx.Description)
	}
}
