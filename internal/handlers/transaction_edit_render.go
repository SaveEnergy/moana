package handlers

import (
	"net/http"

	"moana/internal/httperr"
	"moana/internal/store"
	"moana/internal/txform"
)

func (a *App) renderTransactionEdit(w http.ResponseWriter, u *store.User, data txEditFormData) {
	a.renderShell(w, "transactions_edit.html", data, layoutShell("Edit entry", "history", u))
}

// renderTransactionEditFailed re-renders the edit form after POST validation failure (keeps user input).
func (a *App) renderTransactionEditFailed(w http.ResponseWriter, r *http.Request, u *store.User, id int64, next, msg string) {
	ctx := r.Context()
	cats, err := a.Store.ListCategories(ctx, u.HouseholdID)
	if err != nil {
		httperr.Internal(w, r, err)
		return
	}
	sel := categoryIDFromForm(r.FormValue(txform.FieldCategoryID))
	kind := r.FormValue(txform.FieldKind)
	if kind != "income" && kind != "expense" {
		kind = "income"
	}
	data := txEditFormData{
		Error:         msg,
		Categories:    cats,
		TxID:          id,
		Kind:          kind,
		Amount:        r.FormValue(txform.FieldAmount),
		OccurredOn:    r.FormValue(txform.FieldOccurredOn),
		Description:   r.FormValue(txform.FieldDescription),
		SelectedCatID: sel,
		Next:          next,
	}
	a.renderTransactionEdit(w, u, data)
}
