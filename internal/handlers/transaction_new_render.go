package handlers

import (
	"net/http"

	"moana/internal/httperr"
	"moana/internal/store"
	"moana/internal/tz"
)

func (a *App) transactionsError(w http.ResponseWriter, r *http.Request, u *store.User, msg string) {
	ctx := r.Context()
	cats, err := a.Store.ListCategories(ctx, u.HouseholdID)
	if err != nil {
		httperr.Internal(w, r, err)
		return
	}
	loc := tz.DisplayLocation(r)
	today := todayLocalCalendarDate(loc)
	sel := categoryIDFromForm(r.FormValue("category_id"))
	data := txFormData{
		Error:         msg,
		Categories:    cats,
		Today:         today,
		SelectedCatID: sel,
	}
	a.renderShell(w, "transactions_new.html", data, layoutShell("New entry", "tx", u))
}
