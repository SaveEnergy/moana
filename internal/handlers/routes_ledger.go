package handlers

import "net/http"

// registerLedgerRoutes wires transactions, history, and categories (shared money surface).
func registerLedgerRoutes(mux *http.ServeMux, app *App) {
	mux.Handle("GET /transactions", app.WithAuth(app.Transactions))
	mux.Handle("POST /transactions", app.WithAuth(app.TransactionCreate))
	mux.Handle("GET /transactions/{id}/edit", app.WithAuth(app.TransactionEdit))
	mux.Handle("POST /transactions/{id}", app.WithAuth(app.TransactionUpdate))
	mux.Handle("GET /history", app.WithAuth(app.History))
	mux.Handle("GET /categories", app.WithAuth(app.Categories))
	mux.Handle("POST /categories", app.WithAuth(app.CategoryCreate))
	mux.Handle("POST /categories/update", app.WithAuth(app.CategoryUpdate))
	mux.Handle("POST /categories/delete", app.WithAuth(app.CategoryDelete))
}
