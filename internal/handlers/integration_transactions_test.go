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

func TestTransactionCreateValidationErrorRendersForm(t *testing.T) {
	t.Parallel()
	app, srv, cleanup := testutil.NewAppServer(t)
	defer cleanup()
	testutil.MustCreateUser(t, app, "badamount@moana.test", "pw", "user")
	client := testutil.NewCookieClient(t)
	testutil.MustLogin(t, client, srv.URL, "badamount@moana.test", "pw")
	day := time.Now().UTC().Format("2006-01-02")
	resp, err := client.PostForm(srv.URL+"/transactions", url.Values{
		"amount":      {"not-a-number"},
		"kind":        {"expense"},
		"occurred_on": {day},
		"description": {"x"},
	})
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("status %d (expected 200 with form error)", resp.StatusCode)
	}
	body, _ := io.ReadAll(resp.Body)
	s := string(body)
	if !strings.Contains(s, `class="alert alert-error"`) {
		t.Fatalf("expected validation error in body: %s", s[:min(400, len(s))])
	}
}

func TestCreateExpenseStoresNegativeCents(t *testing.T) {
	t.Parallel()
	app, srv, cleanup := testutil.NewAppServer(t)
	defer cleanup()
	testutil.MustCreateUser(t, app, "tx@moana.test", "pw", "user")
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
	app, srv, cleanup := testutil.NewAppServer(t)
	defer cleanup()
	ctx := context.Background()
	uid := testutil.MustCreateUser(t, app, "edit@moana.test", "pw", "user")
	u, err := app.Store.GetUserByID(ctx, uid)
	if err != nil || u == nil {
		t.Fatal(err)
	}
	hid := u.HouseholdID
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

func TestTransactionCreate_invalidCategoryIDShowsMessage(t *testing.T) {
	t.Parallel()
	app, srv, cleanup := testutil.NewAppServer(t)
	defer cleanup()
	testutil.MustCreateUser(t, app, "badcat@moana.test", "pw", "user")
	client := testutil.NewCookieClient(t)
	testutil.MustLogin(t, client, srv.URL, "badcat@moana.test", "pw")
	day := time.Now().UTC().Format("2006-01-02")
	resp, err := client.PostForm(srv.URL+"/transactions", url.Values{
		"amount":      {"10.00"},
		"kind":        {"expense"},
		"occurred_on": {day},
		"description": {"x"},
		"category_id": {"99999999999999"},
	})
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("status %d want 200 with form error", resp.StatusCode)
	}
	body, _ := io.ReadAll(resp.Body)
	s := string(body)
	if !strings.Contains(s, "That category is not valid for this household.") {
		t.Fatalf("expected user-facing category error, got: %s", s[:min(600, len(s))])
	}
}

func TestTransactionEditNotFound(t *testing.T) {
	t.Parallel()
	app, srv, cleanup := testutil.NewAppServer(t)
	defer cleanup()
	testutil.MustCreateUser(t, app, "nf@moana.test", "pw", "user")
	client := testutil.NewCookieClient(t)
	testutil.MustLogin(t, client, srv.URL, "nf@moana.test", "pw")
	resp, err := client.Get(srv.URL + "/transactions/999999999/edit")
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusNotFound {
		t.Fatalf("status %d want 404", resp.StatusCode)
	}
}

func TestTransactionUpdate_validationErrorRendersForm(t *testing.T) {
	t.Parallel()
	app, srv, cleanup := testutil.NewAppServer(t)
	defer cleanup()
	ctx := context.Background()
	uid := testutil.MustCreateUser(t, app, "tx-upd-bad@moana.test", "pw", "user")
	u, err := app.Store.GetUserByID(ctx, uid)
	if err != nil || u == nil {
		t.Fatal(err)
	}
	client := testutil.NewCookieClient(t)
	testutil.MustLogin(t, client, srv.URL, "tx-upd-bad@moana.test", "pw")
	day := time.Now().UTC().Format("2006-01-02")
	resp, err := client.PostForm(srv.URL+"/transactions", url.Values{
		"amount":      {"10.00"},
		"kind":        {"expense"},
		"occurred_on": {day},
		"description": {"seed"},
	})
	if err != nil {
		t.Fatal(err)
	}
	resp.Body.Close()
	txs, err := app.Store.ListTransactions(ctx, u.HouseholdID, store.TransactionFilter{Limit: 1})
	if err != nil || len(txs) != 1 {
		t.Fatalf("list: %v len=%d", err, len(txs))
	}
	id := txs[0].ID

	resp2, err := client.PostForm(fmt.Sprintf("%s/transactions/%d", srv.URL, id), url.Values{
		"next":        {"/history"},
		"amount":      {"not-a-number"},
		"kind":        {"expense"},
		"occurred_on": {day},
		"description": {"x"},
		"category_id": {""},
	})
	if err != nil {
		t.Fatal(err)
	}
	defer resp2.Body.Close()
	if resp2.StatusCode != http.StatusOK {
		t.Fatalf("status %d want 200 with form error", resp2.StatusCode)
	}
	body, _ := io.ReadAll(resp2.Body)
	s := string(body)
	if !strings.Contains(s, `class="alert alert-error"`) {
		t.Fatalf("expected validation error on edit form: %s", s[:min(500, len(s))])
	}
}
