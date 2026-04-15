package handlers

import (
	"net/http"
	"strconv"
	"strings"

	"moana/internal/category"
	"moana/internal/httperr"
	"moana/internal/store"
)

// Categories lists categories and handles create/delete.
func (a *App) Categories(w http.ResponseWriter, r *http.Request, u *store.User) {
	ctx := r.Context()
	data, err := category.BuildCategoriesList(ctx, a.Store, u.HouseholdID, "")
	if err != nil {
		httperr.Internal(w, r, err)
		return
	}
	a.renderShell(w, "categories.html", data, layoutShell("Categories", "cat", u))
}

// CategoryCreate handles POST /categories.
func (a *App) CategoryCreate(w http.ResponseWriter, r *http.Request, u *store.User) {
	if !requireParseForm(w, r) {
		return
	}
	name := strings.TrimSpace(r.FormValue("name"))
	if name == "" {
		a.categoriesWithError(w, r, u, "Name is required.")
		return
	}
	icon := category.NormalizeStoredIcon(r.FormValue("icon"))
	color := category.ParseColorFromForm(r)
	ctx := r.Context()
	if _, err := a.Store.CreateCategory(ctx, u.HouseholdID, name, icon, color); err != nil {
		a.categoriesWithError(w, r, u, userFacingStoreMessage(err))
		return
	}
	http.Redirect(w, r, "/categories", http.StatusSeeOther)
}

// CategoryDelete handles POST /categories/delete.
func (a *App) CategoryDelete(w http.ResponseWriter, r *http.Request, u *store.User) {
	if !requireParseForm(w, r) {
		return
	}
	id, err := strconv.ParseInt(r.FormValue("id"), 10, 64)
	if err != nil || id <= 0 {
		http.Redirect(w, r, "/categories", http.StatusSeeOther)
		return
	}
	ctx := r.Context()
	if err := a.Store.DeleteCategory(ctx, u.HouseholdID, id); err != nil {
		a.categoriesWithError(w, r, u, userFacingStoreMessage(err))
		return
	}
	http.Redirect(w, r, "/categories", http.StatusSeeOther)
}

// CategoryUpdate handles POST /categories/update (name + icon).
func (a *App) CategoryUpdate(w http.ResponseWriter, r *http.Request, u *store.User) {
	if !requireParseForm(w, r) {
		return
	}
	id, err := strconv.ParseInt(r.FormValue("id"), 10, 64)
	if err != nil || id <= 0 {
		a.categoriesWithError(w, r, u, "Invalid category.")
		return
	}
	name := strings.TrimSpace(r.FormValue("name"))
	if name == "" {
		a.categoriesWithError(w, r, u, "Name is required.")
		return
	}
	icon := category.NormalizeStoredIcon(r.FormValue("icon"))
	color := category.ParseColorFromForm(r)
	ctx := r.Context()
	if err := a.Store.UpdateCategory(ctx, u.HouseholdID, id, name, icon, color); err != nil {
		a.categoriesWithError(w, r, u, userFacingStoreMessage(err))
		return
	}
	http.Redirect(w, r, "/categories", http.StatusSeeOther)
}
