package handlers

import (
	"net/http"

	"moana/internal/category"
	"moana/internal/store"
)

func (a *App) categoriesWithError(w http.ResponseWriter, r *http.Request, u *store.User, msg string) {
	ctx := r.Context()
	data, err := category.BuildCategoriesList(ctx, a.Store, u.HouseholdID, msg)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	a.renderShell(w, "categories.html", data, layoutShell("Categories", "cat", u))
}
