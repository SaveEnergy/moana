package handlers

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"moana/internal/store"
)

// Categories lists categories and handles create/delete.
func (a *App) Categories(w http.ResponseWriter, r *http.Request, u *store.User) {
	ctx := r.Context()
	cats, err := a.Store.ListCategories(ctx, u.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	data := struct {
		Error      string
		Categories []store.Category
	}{
		Categories: cats,
	}
	ld := LayoutData{
		Title:  "Categories",
		User:   u,
		Year:   time.Now().UTC().Year(),
		Active: "cat",
	}
	a.renderShell(w, "categories_inner.html", data, ld)
}

// CategoryCreate handles POST /categories.
func (a *App) CategoryCreate(w http.ResponseWriter, r *http.Request, u *store.User) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "bad form", http.StatusBadRequest)
		return
	}
	name := strings.TrimSpace(r.FormValue("name"))
	if name == "" {
		a.categoriesWithError(w, r, u, "Name is required.")
		return
	}
	ctx := r.Context()
	if _, err := a.Store.CreateCategory(ctx, u.ID, name); err != nil {
		a.categoriesWithError(w, r, u, "Could not create category (duplicate name?).")
		return
	}
	http.Redirect(w, r, "/categories", http.StatusSeeOther)
}

// CategoryDelete handles POST /categories/delete.
func (a *App) CategoryDelete(w http.ResponseWriter, r *http.Request, u *store.User) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "bad form", http.StatusBadRequest)
		return
	}
	id, err := strconv.ParseInt(r.FormValue("id"), 10, 64)
	if err != nil || id <= 0 {
		http.Redirect(w, r, "/categories", http.StatusSeeOther)
		return
	}
	ctx := r.Context()
	if err := a.Store.DeleteCategory(ctx, u.ID, id); err != nil {
		a.categoriesWithError(w, r, u, err.Error())
		return
	}
	http.Redirect(w, r, "/categories", http.StatusSeeOther)
}

func (a *App) categoriesWithError(w http.ResponseWriter, r *http.Request, u *store.User, msg string) {
	ctx := r.Context()
	cats, err := a.Store.ListCategories(ctx, u.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	data := struct {
		Error      string
		Categories []store.Category
	}{
		Error:      msg,
		Categories: cats,
	}
	ld := LayoutData{
		Title:  "Categories",
		User:   u,
		Year:   time.Now().UTC().Year(),
		Active: "cat",
	}
	a.renderShell(w, "categories_inner.html", data, ld)
}
