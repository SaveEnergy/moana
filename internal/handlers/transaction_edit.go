package handlers

import (
	"database/sql"
	"errors"
	"net/http"

	"golang.org/x/sync/errgroup"

	"moana/internal/httperr"
	"moana/internal/money"
	"moana/internal/safepath"
	"moana/internal/store"
	"moana/internal/txform"
	"moana/internal/tz"
)

// TransactionEdit shows the edit form for a transaction (GET /transactions/{id}/edit).
// It loads the transaction row and the household category list concurrently (independent reads).
func (a *App) TransactionEdit(w http.ResponseWriter, r *http.Request, u *store.User) {
	id, ok := pathPositiveInt64(r, "id")
	if !ok {
		http.NotFound(w, r)
		return
	}
	ctx := r.Context()
	var (
		tx   *store.Transaction
		cats []store.Category
	)
	g, gctx := errgroup.WithContext(ctx)
	g.Go(func() error {
		var err error
		tx, err = a.Store.GetTransactionByID(gctx, u.HouseholdID, id)
		return err
	})
	g.Go(func() error {
		var err error
		cats, err = a.Store.ListCategories(gctx, u.HouseholdID)
		return err
	})
	if err := g.Wait(); err != nil {
		httperr.Internal(w, r, err)
		return
	}
	if tx == nil {
		http.NotFound(w, r)
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
		OccurredOn:    formatLocalCalendarDate(tx.OccurredAt, loc),
		Description:   tx.Description,
		SelectedCatID: sel,
		Next:          safepath.Internal(r.URL.Query().Get(TransactionNextQueryParam)),
	}
	a.renderTransactionEdit(w, r, u, data)
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
	next := safepath.Internal(r.FormValue(TransactionNextQueryParam))
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
		r.FormValue(txform.FieldAmount),
		r.FormValue(txform.FieldOccurredOn),
		r.FormValue(txform.FieldDescription),
		r.FormValue(txform.FieldCategoryID),
		r.FormValue(txform.FieldKind),
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
