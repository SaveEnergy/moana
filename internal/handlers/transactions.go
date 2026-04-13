package handlers

import (
	"database/sql"
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	"moana/internal/money"
	"moana/internal/store"
	"moana/internal/timeutil"
)

func absCents(n int64) int64 {
	if n < 0 {
		return -n
	}
	return n
}

// txFormData is the new-entry form (GET /transactions, form errors).
type txFormData struct {
	Error      string
	Categories []store.Category
	UserTZ     string
	Today      string
}

// txEditFormData is the edit form (GET/POST /transactions/{id}/edit).
type txEditFormData struct {
	Error         string
	Categories    []store.Category
	UserTZ        string
	TxID          int64
	Kind          string
	Amount        string
	OccurredOn    string
	Description   string
	SelectedCatID int64
	Next          string
}

func safeRedirectTarget(next string) string {
	next = strings.TrimSpace(next)
	if next == "" || !strings.HasPrefix(next, "/") || strings.HasPrefix(next, "//") {
		return "/history"
	}
	return next
}

// Transactions shows the new income entry form only.
func (a *App) Transactions(w http.ResponseWriter, r *http.Request, u *store.User) {
	ctx := r.Context()
	loc := timeutil.LoadLocation(u.Timezone)
	today := time.Now().In(loc).Format("2006-01-02")
	cats, err := a.Store.ListCategories(ctx, u.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	data := txFormData{
		Error:      "",
		Categories: cats,
		UserTZ:     u.Timezone,
		Today:      today,
	}
	ld := LayoutData{
		Title:  "New entry",
		User:   u,
		Year:   time.Now().UTC().Year(),
		Active: "tx",
	}
	a.renderShell(w, "transactions_new_inner.html", data, ld)
}

// TransactionCreate handles POST /transactions.
func (a *App) TransactionCreate(w http.ResponseWriter, r *http.Request, u *store.User) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "bad form", http.StatusBadRequest)
		return
	}
	amountStr := r.FormValue("amount")
	dateStr := r.FormValue("occurred_on")
	desc := r.FormValue("description")
	catStr := r.FormValue("category_id")

	cents, err := money.ParseEURToCents(amountStr)
	if err != nil {
		a.transactionsError(w, r, u, err.Error())
		return
	}
	cents = absCents(cents)
	if cents == 0 {
		a.transactionsError(w, r, u, "Amount must be greater than zero.")
		return
	}
	if r.FormValue("kind") == "expense" {
		cents = -cents
	}
	if dateStr == "" {
		a.transactionsError(w, r, u, "Date is required.")
		return
	}
	loc := timeutil.LoadLocation(u.Timezone)
	dayStart, err := time.ParseInLocation("2006-01-02", dateStr, loc)
	if err != nil {
		a.transactionsError(w, r, u, "Invalid date.")
		return
	}
	occurred := dayStart.UTC()

	var catID *int64
	if catStr != "" {
		id, err := strconv.ParseInt(catStr, 10, 64)
		if err != nil {
			a.transactionsError(w, r, u, "Invalid category.")
			return
		}
		catID = &id
	}

	ctx := r.Context()
	if _, err := a.Store.CreateTransaction(ctx, u.ID, cents, occurred, desc, catID); err != nil {
		a.transactionsError(w, r, u, err.Error())
		return
	}
	http.Redirect(w, r, "/history", http.StatusSeeOther)
}

func (a *App) transactionsError(w http.ResponseWriter, r *http.Request, u *store.User, msg string) {
	ctx := r.Context()
	cats, err := a.Store.ListCategories(ctx, u.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	loc := timeutil.LoadLocation(u.Timezone)
	today := time.Now().In(loc).Format("2006-01-02")
	data := txFormData{
		Error:      msg,
		Categories: cats,
		UserTZ:     u.Timezone,
		Today:      today,
	}
	ld := LayoutData{
		Title:  "New entry",
		User:   u,
		Year:   time.Now().UTC().Year(),
		Active: "tx",
	}
	a.renderShell(w, "transactions_new_inner.html", data, ld)
}

// TransactionEdit shows the edit form for a transaction (GET /transactions/{id}/edit).
func (a *App) TransactionEdit(w http.ResponseWriter, r *http.Request, u *store.User) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil || id < 1 {
		http.NotFound(w, r)
		return
	}
	ctx := r.Context()
	tx, err := a.Store.GetTransactionByID(ctx, u.ID, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if tx == nil {
		http.NotFound(w, r)
		return
	}
	cats, err := a.Store.ListCategories(ctx, u.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	loc := timeutil.LoadLocation(u.Timezone)
	kind := "income"
	if tx.AmountCents < 0 {
		kind = "expense"
	}
	sel := int64(0)
	if tx.CategoryID.Valid {
		sel = tx.CategoryID.Int64
	}
	data := txEditFormData{
		Error:         "",
		Categories:    cats,
		UserTZ:        u.Timezone,
		TxID:          tx.ID,
		Kind:          kind,
		Amount:        money.FormatDecimalEURAbs(tx.AmountCents),
		OccurredOn:    tx.OccurredAt.In(loc).Format("2006-01-02"),
		Description:   tx.Description,
		SelectedCatID: sel,
		Next:          safeRedirectTarget(r.URL.Query().Get("next")),
	}
	a.renderTransactionEdit(w, u, data)
}

// TransactionUpdate applies edits (POST /transactions/{id}).
func (a *App) TransactionUpdate(w http.ResponseWriter, r *http.Request, u *store.User) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "bad form", http.StatusBadRequest)
		return
	}
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil || id < 1 {
		http.NotFound(w, r)
		return
	}
	next := safeRedirectTarget(r.FormValue("next"))
	ctx := r.Context()
	existing, err := a.Store.GetTransactionByID(ctx, u.ID, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if existing == nil {
		http.NotFound(w, r)
		return
	}

	amountStr := r.FormValue("amount")
	dateStr := r.FormValue("occurred_on")
	desc := r.FormValue("description")
	catStr := r.FormValue("category_id")

	cents, err := money.ParseEURToCents(amountStr)
	if err != nil {
		a.renderTransactionEditFailed(w, r, u, id, next, err.Error())
		return
	}
	cents = absCents(cents)
	if cents == 0 {
		a.renderTransactionEditFailed(w, r, u, id, next, "Amount must be greater than zero.")
		return
	}
	if r.FormValue("kind") == "expense" {
		cents = -cents
	}
	if dateStr == "" {
		a.renderTransactionEditFailed(w, r, u, id, next, "Date is required.")
		return
	}
	loc := timeutil.LoadLocation(u.Timezone)
	dayStart, err := time.ParseInLocation("2006-01-02", dateStr, loc)
	if err != nil {
		a.renderTransactionEditFailed(w, r, u, id, next, "Invalid date.")
		return
	}
	occurred := dayStart.UTC()

	var catID *int64
	if catStr != "" {
		cid, err := strconv.ParseInt(catStr, 10, 64)
		if err != nil {
			a.renderTransactionEditFailed(w, r, u, id, next, "Invalid category.")
			return
		}
		catID = &cid
	}

	if err := a.Store.UpdateTransaction(ctx, u.ID, id, cents, occurred, desc, catID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.NotFound(w, r)
			return
		}
		a.renderTransactionEditFailed(w, r, u, id, next, err.Error())
		return
	}
	http.Redirect(w, r, next, http.StatusSeeOther)
}

func (a *App) renderTransactionEdit(w http.ResponseWriter, u *store.User, data txEditFormData) {
	ld := LayoutData{
		Title:  "Edit entry",
		User:   u,
		Year:   time.Now().UTC().Year(),
		Active: "history",
	}
	a.renderShell(w, "transactions_edit_inner.html", data, ld)
}

// renderTransactionEditFailed re-renders the edit form after POST validation failure (keeps user input).
func (a *App) renderTransactionEditFailed(w http.ResponseWriter, r *http.Request, u *store.User, id int64, next, msg string) {
	ctx := r.Context()
	cats, err := a.Store.ListCategories(ctx, u.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	sel := int64(0)
	if c := r.FormValue("category_id"); c != "" {
		sel, _ = strconv.ParseInt(c, 10, 64)
	}
	kind := r.FormValue("kind")
	if kind != "income" && kind != "expense" {
		kind = "income"
	}
	data := txEditFormData{
		Error:         msg,
		Categories:    cats,
		UserTZ:        u.Timezone,
		TxID:          id,
		Kind:          kind,
		Amount:        r.FormValue("amount"),
		OccurredOn:    r.FormValue("occurred_on"),
		Description:   r.FormValue("description"),
		SelectedCatID: sel,
		Next:          next,
	}
	a.renderTransactionEdit(w, u, data)
}
