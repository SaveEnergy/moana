package handlers_test

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"testing"
	"time"

	"moana/internal/store"
	"moana/internal/testutil"
)

func TestCreateExpenseStoresNegativeCents(t *testing.T) {
	t.Parallel()
	app, cleanup := testApp(t)
	defer cleanup()
	testutil.MustCreateUser(t, app, "tx@moana.test", "pw", "user")
	srv := testutil.NewServer(t, app)
	defer srv.Close()
	client := testutil.NewCookieClient(t)
	testutil.MustLogin(t, client, srv.URL, "tx@moana.test", "pw")
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
	uid := testutil.MustCreateUser(t, app, "edit@moana.test", "pw", "user")
	u, err := app.Store.GetUserByID(ctx, uid)
	if err != nil || u == nil {
		t.Fatal(err)
	}
	hid := u.HouseholdID
	srv := testutil.NewServer(t, app)
	defer srv.Close()
	client := testutil.NewCookieClient(t)
	testutil.MustLogin(t, client, srv.URL, "edit@moana.test", "pw")
	day := time.Now().UTC().Format("2006-01-02")
	resp, err := client.PostForm(srv.URL+"/transactions", url.Values{
		"amount":      {"10.00"},
		"kind":        {"expense"},
		"occurred_on": {day},
		"description": {"coffee"},
		"category_id": {""},
	})
	if err != nil {
		t.Fatal(err)
	}
	resp.Body.Close()
	txs, err := app.Store.ListTransactions(ctx, hid, store.TransactionFilter{Limit: 1})
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
	b := string(body)
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("edit get status %d: %s", resp.StatusCode, b[:min(200, len(b))])
	}
	if !strings.Contains(b, "Edit entry") || !strings.Contains(b, "10.00") {
		t.Fatal("expected edit form")
	}
	resp2, err := client.PostForm(fmt.Sprintf("%s/transactions/%d", srv.URL, id), url.Values{
		"next":        {"/history"},
		"amount":      {"20.00"},
		"kind":        {"expense"},
		"occurred_on": {day},
		"description": {"coffee fixed"},
		"category_id": {""},
	})
	if err != nil {
		t.Fatal(err)
	}
	defer resp2.Body.Close()
	if resp2.StatusCode != http.StatusOK {
		t.Fatalf("update status %d (expected 200 after redirect to history)", resp2.StatusCode)
	}
	tx, err := app.Store.GetTransactionByID(ctx, hid, id)
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
