package handlers_test

import (
	"context"
	"io"
	"net/http"
	"net/url"
	"strings"
	"testing"
	"time"

	"moana/internal/handlers"
	"moana/internal/store"
	"moana/internal/testutil"
	"moana/internal/txform"
)

func TestTransactionCreateValidationErrorRendersForm(t *testing.T) {
	t.Parallel()
	app, srv, cleanup := testutil.NewAppServer(t)
	defer cleanup()
	testutil.MustCreateUser(t, app, "badamount@moana.test", "pw", "user")
	client := testutil.NewCookieClient(t)
	testutil.MustLogin(t, client, srv.URL, "badamount@moana.test", "pw")
	day := time.Now().UTC().Format("2006-01-02")
	resp, err := client.PostForm(srv.URL+handlers.TransactionsPath, url.Values{
		txform.FieldAmount:      {"not-a-number"},
		txform.FieldKind:        {"expense"},
		txform.FieldOccurredOn:  {day},
		txform.FieldDescription: {"x"},
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
	assertBodyHasErrorAlert(t, s)
}

func TestTransactionCreate_redirectsToHistory(t *testing.T) {
	t.Parallel()
	app, srv, cleanup := testutil.NewAppServer(t)
	defer cleanup()
	testutil.MustCreateUser(t, app, "redir-hist@moana.test", "pw", "user")
	client := testutil.NewCookieClient(t)
	testutil.MustLogin(t, client, srv.URL, "redir-hist@moana.test", "pw")
	day := time.Now().UTC().Format("2006-01-02")
	noFollow := *client
	noFollow.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}
	resp, err := noFollow.PostForm(srv.URL+handlers.TransactionsPath, url.Values{
		txform.FieldAmount:      {"7.00"},
		txform.FieldKind:        {"expense"},
		txform.FieldOccurredOn:  {day},
		txform.FieldDescription: {"redirect check"},
	})
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusSeeOther {
		t.Fatalf("status %d want %d", resp.StatusCode, http.StatusSeeOther)
	}
	loc := resp.Header.Get("Location")
	if !strings.Contains(loc, handlers.HistoryPath) {
		t.Fatalf("Location %q want redirect target %s", loc, handlers.HistoryPath)
	}
}

func TestTransactionCreate_zeroAmountShowsMessage(t *testing.T) {
	t.Parallel()
	app, srv, cleanup := testutil.NewAppServer(t)
	defer cleanup()
	testutil.MustCreateUser(t, app, "zero-amt@moana.test", "pw", "user")
	client := testutil.NewCookieClient(t)
	testutil.MustLogin(t, client, srv.URL, "zero-amt@moana.test", "pw")
	day := time.Now().UTC().Format("2006-01-02")
	resp, err := client.PostForm(srv.URL+handlers.TransactionsPath, url.Values{
		txform.FieldAmount:      {"0.00"},
		txform.FieldKind:        {"expense"},
		txform.FieldOccurredOn:  {day},
		txform.FieldDescription: {"x"},
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
	if !strings.Contains(s, "Amount must be greater than zero.") {
		t.Fatalf("expected zero-amount copy, got: %s", s[:min(500, len(s))])
	}
}

func TestTransactionCreate_emptyDateShowsMessage(t *testing.T) {
	t.Parallel()
	app, srv, cleanup := testutil.NewAppServer(t)
	defer cleanup()
	testutil.MustCreateUser(t, app, "nodate@moana.test", "pw", "user")
	client := testutil.NewCookieClient(t)
	testutil.MustLogin(t, client, srv.URL, "nodate@moana.test", "pw")
	resp, err := client.PostForm(srv.URL+handlers.TransactionsPath, url.Values{
		txform.FieldAmount:      {"10.00"},
		txform.FieldKind:        {"expense"},
		txform.FieldOccurredOn:  {""},
		txform.FieldDescription: {"x"},
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
	if !strings.Contains(s, "Date is required.") {
		t.Fatalf("expected date required copy, got: %s", s[:min(500, len(s))])
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
	resp, err := client.PostForm(srv.URL+handlers.TransactionsPath, url.Values{
		txform.FieldAmount:      {"25.50"},
		txform.FieldKind:        {"expense"},
		txform.FieldOccurredOn:  {day},
		txform.FieldDescription: {"test expense"},
	})
	if err != nil {
		t.Fatal(err)
	}
	resp.Body.Close()
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusSeeOther {
		t.Fatalf("post status %d", resp.StatusCode)
	}
	resp2, err := client.Get(srv.URL + handlers.HistoryPath)
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
	resp, err := client.PostForm(srv.URL+handlers.TransactionsPath, url.Values{
		txform.FieldAmount:      {"10.00"},
		txform.FieldKind:        {"expense"},
		txform.FieldOccurredOn:  {day},
		txform.FieldDescription: {"coffee"},
		txform.FieldCategoryID:  {""},
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
	editURL := srv.URL + handlers.TransactionEditURLPath(id) + "?" + handlers.TransactionNextQueryParam + "=" + url.QueryEscape(handlers.HistoryPath)
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
	resp2, err := client.PostForm(srv.URL+handlers.TransactionURLPath(id), url.Values{
		handlers.TransactionNextQueryParam: {handlers.HistoryPath},
		txform.FieldAmount:                 {"20.00"},
		txform.FieldKind:                   {"expense"},
		txform.FieldOccurredOn:             {day},
		txform.FieldDescription:            {"coffee fixed"},
		txform.FieldCategoryID:             {""},
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

func TestTransactionEdit_preservesSafeNextQuery(t *testing.T) {
	t.Parallel()
	app, srv, cleanup := testutil.NewAppServer(t)
	defer cleanup()
	ctx := context.Background()
	uid := testutil.MustCreateUser(t, app, "next-safe@moana.test", "pw", "user")
	u, err := app.Store.GetUserByID(ctx, uid)
	if err != nil || u == nil {
		t.Fatal(err)
	}
	client := testutil.NewCookieClient(t)
	testutil.MustLogin(t, client, srv.URL, "next-safe@moana.test", "pw")
	day := time.Now().UTC().Format("2006-01-02")
	resp, err := client.PostForm(srv.URL+handlers.TransactionsPath, url.Values{
		txform.FieldAmount:      {"5.00"},
		txform.FieldKind:        {"expense"},
		txform.FieldOccurredOn:  {day},
		txform.FieldDescription: {"n"},
	})
	if err != nil {
		t.Fatal(err)
	}
	resp.Body.Close()
	txs, err := app.Store.ListTransactions(ctx, u.HouseholdID, store.TransactionFilter{Limit: 1})
	if err != nil || len(txs) != 1 {
		t.Fatalf("list: %v", err)
	}
	id := txs[0].ID

	resp2, err := client.Get(srv.URL + handlers.TransactionEditURLPath(id) + "?" + handlers.TransactionNextQueryParam + "=" + url.QueryEscape(handlers.CategoriesPath))
	if err != nil {
		t.Fatal(err)
	}
	defer resp2.Body.Close()
	body, _ := io.ReadAll(resp2.Body)
	s := string(body)
	wantNext := `value="` + handlers.CategoriesPath + `"`
	nameNext := `name="` + handlers.TransactionNextQueryParam + `"`
	if !strings.Contains(s, nameNext) || !strings.Contains(s, wantNext) {
		t.Fatalf("expected hidden next=%s, got: %s", handlers.CategoriesPath, s[:min(600, len(s))])
	}
}

func TestTransactionEdit_unsafeNextQueryUsesDefaultInForm(t *testing.T) {
	t.Parallel()
	app, srv, cleanup := testutil.NewAppServer(t)
	defer cleanup()
	ctx := context.Background()
	uid := testutil.MustCreateUser(t, app, "next-unsafe@moana.test", "pw", "user")
	u, err := app.Store.GetUserByID(ctx, uid)
	if err != nil || u == nil {
		t.Fatal(err)
	}
	client := testutil.NewCookieClient(t)
	testutil.MustLogin(t, client, srv.URL, "next-unsafe@moana.test", "pw")
	day := time.Now().UTC().Format("2006-01-02")
	resp, err := client.PostForm(srv.URL+handlers.TransactionsPath, url.Values{
		txform.FieldAmount:      {"5.00"},
		txform.FieldKind:        {"expense"},
		txform.FieldOccurredOn:  {day},
		txform.FieldDescription: {"n"},
	})
	if err != nil {
		t.Fatal(err)
	}
	resp.Body.Close()
	txs, err := app.Store.ListTransactions(ctx, u.HouseholdID, store.TransactionFilter{Limit: 1})
	if err != nil || len(txs) != 1 {
		t.Fatalf("list: %v", err)
	}
	id := txs[0].ID

	resp2, err := client.Get(srv.URL + handlers.TransactionEditURLPath(id) + "?" + handlers.TransactionNextQueryParam + "=//evil.com/foo")
	if err != nil {
		t.Fatal(err)
	}
	defer resp2.Body.Close()
	body, _ := io.ReadAll(resp2.Body)
	s := string(body)
	if strings.Contains(s, "evil.com") {
		t.Fatalf("unsafe next host leaked into HTML")
	}
	idx := strings.Index(s, `name="`+handlers.TransactionNextQueryParam+`"`)
	if idx < 0 {
		t.Fatal("missing hidden next field")
	}
	snippet := s[idx:min(idx+120, len(s))]
	wantNext := `value="` + handlers.HistoryPath + `"`
	if !strings.Contains(snippet, wantNext) {
		t.Fatalf("expected sanitized default %s in hidden next, snippet: %q", handlers.HistoryPath, snippet)
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
	resp, err := client.PostForm(srv.URL+handlers.TransactionsPath, url.Values{
		txform.FieldAmount:      {"10.00"},
		txform.FieldKind:        {"expense"},
		txform.FieldOccurredOn:  {day},
		txform.FieldDescription: {"x"},
		txform.FieldCategoryID:  {"99999999999999"},
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
	resp, err := client.Get(srv.URL + handlers.TransactionEditURLPath(999999999))
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusNotFound {
		t.Fatalf("status %d want 404", resp.StatusCode)
	}
}

func TestTransactionEdit_invalidPathIDReturns404(t *testing.T) {
	t.Parallel()
	app, srv, cleanup := testutil.NewAppServer(t)
	defer cleanup()
	testutil.MustCreateUser(t, app, "bad-id-get@moana.test", "pw", "user")
	client := testutil.NewCookieClient(t)
	testutil.MustLogin(t, client, srv.URL, "bad-id-get@moana.test", "pw")
	resp, err := client.Get(srv.URL + handlers.TransactionsPath + "/not-a-number/edit")
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusNotFound {
		t.Fatalf("status %d want 404", resp.StatusCode)
	}
}

func TestTransactionUpdate_invalidPathIDReturns404(t *testing.T) {
	t.Parallel()
	app, srv, cleanup := testutil.NewAppServer(t)
	defer cleanup()
	testutil.MustCreateUser(t, app, "bad-id-post@moana.test", "pw", "user")
	client := testutil.NewCookieClient(t)
	testutil.MustLogin(t, client, srv.URL, "bad-id-post@moana.test", "pw")
	day := time.Now().UTC().Format("2006-01-02")
	resp, err := client.PostForm(srv.URL+handlers.TransactionsPath+"/not-a-number", url.Values{
		handlers.TransactionNextQueryParam: {handlers.HistoryPath},
		txform.FieldAmount:                 {"1.00"},
		txform.FieldKind:                   {"expense"},
		txform.FieldOccurredOn:             {day},
		txform.FieldDescription:            {"x"},
	})
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
	resp, err := client.PostForm(srv.URL+handlers.TransactionsPath, url.Values{
		txform.FieldAmount:      {"10.00"},
		txform.FieldKind:        {"expense"},
		txform.FieldOccurredOn:  {day},
		txform.FieldDescription: {"seed"},
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

	resp2, err := client.PostForm(srv.URL+handlers.TransactionURLPath(id), url.Values{
		handlers.TransactionNextQueryParam: {handlers.HistoryPath},
		txform.FieldAmount:                 {"not-a-number"},
		txform.FieldKind:                   {"expense"},
		txform.FieldOccurredOn:             {day},
		txform.FieldDescription:            {"x"},
		txform.FieldCategoryID:             {""},
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
	assertBodyHasErrorAlert(t, s)
}
