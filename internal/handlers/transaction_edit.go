package handlers

import (
	"database/sql"
	"errors"
	"net/http"
	"strconv"

	"moana/internal/money"
	"moana/internal/safepath"
	"moana/internal/store"
	"moana/internal/txform"
	"moana/internal/tz"
)

// txEditFormData is the edit form (GET/POST /transactions/{id}/edit).
type txEditFormData struct {
	Error         string
	Categories    []store.Category
	TxID          int64
	Kind          string
	Amount        string
	OccurredOn    string
	Description   string
	SelectedCatID int64
	Next          string
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
	loc := tz.DisplayLocation(r)
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
		TxID:          tx.ID,
		Kind:          kind,
		Amount:        money.FormatDecimalEURAbs(tx.AmountCents),
		OccurredOn:    tx.OccurredAt.In(loc).Format("2006-01-02"),
		Description:   tx.Description,
		SelectedCatID: sel,
		Next:          safepath.Internal(r.URL.Query().Get("next")),
	}
	a.renderTransactionEdit(w, u, data)
}

// TransactionUpdate applies edits (POST /transactions/{id}).
func (a *App) TransactionUpdate(w http.ResponseWriter, r *http.Request, u *store.User) {
	if !requireParseForm(w, r) {
		return
	}
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil || id < 1 {
		http.NotFound(w, r)
		return
	}
	next := safepath.Internal(r.FormValue("next"))
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
		a.renderTransactionEditFailed(w, r, u, id, next, errMsg)
		return
	}

	if err := a.Store.UpdateTransaction(ctx, u.ID, id, p.AmountCents, p.OccurredUTC, p.Description, p.CategoryID); err != nil {
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
	a.renderShell(w, "transactions_edit.html", data, layoutShell("Edit entry", "history", u))
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
