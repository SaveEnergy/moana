package handlers

import (
	"net/http"

	"moana/internal/httperr"
	"moana/internal/store"
	"moana/internal/txform"
	"moana/internal/tz"
)

// Transactions shows the new income entry form only.
func (a *App) Transactions(w http.ResponseWriter, r *http.Request, u *store.User) {
	ctx := r.Context()
	loc := tz.DisplayLocation(r)
	today := todayLocalCalendarDate(loc)
	cats, err := a.Store.ListCategories(ctx, u.HouseholdID)
	if err != nil {
		httperr.Internal(w, r, err)
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
		r.FormValue(txform.FieldAmount),
		r.FormValue(txform.FieldOccurredOn),
		r.FormValue(txform.FieldDescription),
		r.FormValue(txform.FieldCategoryID),
		r.FormValue(txform.FieldKind),
		loc,
	)
	if errMsg != "" {
		a.transactionsError(w, r, u, errMsg)
		return
	}

	ctx := r.Context()
	if _, err := a.Store.CreateTransaction(ctx, u.ID, u.HouseholdID, p.AmountCents, p.OccurredUTC, p.Description, p.CategoryID); err != nil {
		a.transactionsError(w, r, u, userFacingStoreMessage(err))
		return
	}
	http.Redirect(w, r, HistoryPath, http.StatusSeeOther)
}
