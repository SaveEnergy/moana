package handlers

import (
	"net/http"

	"moana/internal/category"
	"moana/internal/httperr"
	"moana/internal/store"
)

func (a *App) categoriesWithError(w http.ResponseWriter, r *http.Request, u *store.User, msg string) {
	ctx := r.Context()
	data, err := category.BuildCategoriesList(ctx, a.Store, u.HouseholdID, msg)
	if err != nil {
		httperr.Internal(w, r, err)
		return
	}
	a.renderShell(w, "categories.html", data, layoutShell("Categories", "cat", u))
}
