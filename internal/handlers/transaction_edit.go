package handlers

import (
	"database/sql"
	"errors"
	"net/http"

	"moana/internal/httperr"
	"moana/internal/money"
	"moana/internal/safepath"
	"moana/internal/store"
	"moana/internal/txform"
	"moana/internal/tz"
)

// TransactionEdit shows the edit form for a transaction (GET /transactions/{id}/edit).
func (a *App) TransactionEdit(w http.ResponseWriter, r *http.Request, u *store.User) {
	id, ok := pathPositiveInt64(r, "id")
	if !ok {
		http.NotFound(w, r)
		return
	}
	ctx := r.Context()
	tx, err := a.Store.GetTransactionByID(ctx, u.HouseholdID, id)
	if err != nil {
		httperr.Internal(w, r, err)
		return
	}
	if tx == nil {
		http.NotFound(w, r)
		return
	}
	cats, err := a.Store.ListCategories(ctx, u.HouseholdID)
	if err != nil {
		httperr.Internal(w, r, err)
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
	id, ok := pathPositiveInt64(r, "id")
	if !ok {
		http.NotFound(w, r)
		return
	}
	next := safepath.Internal(r.FormValue("next"))
	ctx := r.Context()
	existing, err := a.Store.GetTransactionByID(ctx, u.HouseholdID, id)
	if err != nil {
		httperr.Internal(w, r, err)
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

	if err := a.Store.UpdateTransaction(ctx, u.HouseholdID, u.ID, id, p.AmountCents, p.OccurredUTC, p.Description, p.CategoryID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.NotFound(w, r)
			return
		}
		a.renderTransactionEditFailed(w, r, u, id, next, userFacingStoreMessage(err))
		return
	}
	http.Redirect(w, r, next, http.StatusSeeOther)
}
