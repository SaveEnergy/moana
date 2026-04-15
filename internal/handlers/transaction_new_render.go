package handlers

import (
	"net/http"
	"strconv"
	"time"

	"moana/internal/store"
	"moana/internal/tz"
)

func (a *App) transactionsError(w http.ResponseWriter, r *http.Request, u *store.User, msg string) {
	ctx := r.Context()
	cats, err := a.Store.ListCategories(ctx, u.HouseholdID)
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
