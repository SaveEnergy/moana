package handlers

import (
	"net/http"
	"strconv"
	"time"

	"moana/internal/store"
	"moana/internal/txform"
	"moana/internal/tz"
)

// txFormData is the new-entry form (GET /transactions, form errors).
type txFormData struct {
	Error         string
	Categories    []store.Category
	Today         string
	SelectedCatID int64
}

// Transactions shows the new income entry form only.
func (a *App) Transactions(w http.ResponseWriter, r *http.Request, u *store.User) {
	ctx := r.Context()
	loc := tz.DisplayLocation(r)
	today := time.Now().In(loc).Format("2006-01-02")
	cats, err := a.Store.ListCategories(ctx, u.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	data := txFormData{
		Error:         "",
		Categories:    cats,
		Today:         today,
		SelectedCatID: 0,
	}
	a.renderShell(w, "transactions_new.html", data, layoutShell("New entry", "tx", u))
}

// TransactionCreate handles POST /transactions.
func (a *App) TransactionCreate(w http.ResponseWriter, r *http.Request, u *store.User) {
	if !requireParseForm(w, r) {
		return
	}
	loc := tz.DisplayLocation(r)
	p, errMsg := txform.Parse(
		r.FormValue("amount"),
		r.FormValue("occurred_on"),
		r.FormValue("description"),
		r.FormValue("category_id"),
		r.FormValue("kind"),
		loc,
	)
	if errMsg != "" {
		a.transactionsError(w, r, u, errMsg)
		return
	}

	ctx := r.Context()
	if _, err := a.Store.CreateTransaction(ctx, u.ID, p.AmountCents, p.OccurredUTC, p.Description, p.CategoryID); err != nil {
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
	loc := tz.DisplayLocation(r)
	today := time.Now().In(loc).Format("2006-01-02")
	sel := int64(0)
	if c := r.FormValue("category_id"); c != "" {
		sel, _ = strconv.ParseInt(c, 10, 64)
	}
	data := txFormData{
		Error:         msg,
		Categories:    cats,
		Today:         today,
		SelectedCatID: sel,
	}
	a.renderShell(w, "transactions_new.html", data, layoutShell("New entry", "tx", u))
}
