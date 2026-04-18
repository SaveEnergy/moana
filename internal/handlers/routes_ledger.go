package handlers

import "net/http"

// registerLedgerRoutes wires transactions, history, and categories (shared money surface).
func registerLedgerRoutes(mux *http.ServeMux, app *App) {
	mux.Handle(http.MethodGet+" "+TransactionsPath, app.WithAuth(app.Transactions))
	mux.Handle(http.MethodPost+" "+TransactionsPath, app.WithAuth(app.TransactionCreate))
	mux.Handle(http.MethodGet+" "+TransactionsPath+"/{id}/edit", app.WithAuth(app.TransactionEdit))
	mux.Handle(http.MethodPost+" "+TransactionsPath+"/{id}", app.WithAuth(app.TransactionUpdate))
	mux.Handle(http.MethodGet+" "+HistoryPath, app.WithAuth(app.History))
	mux.Handle(http.MethodGet+" "+CategoriesPath, app.WithAuth(app.Categories))
	mux.Handle(http.MethodPost+" "+CategoriesPath, app.WithAuth(app.CategoryCreate))
	mux.Handle(http.MethodPost+" "+CategoriesUpdatePath, app.WithAuth(app.CategoryUpdate))
	mux.Handle(http.MethodPost+" "+CategoriesDeletePath, app.WithAuth(app.CategoryDelete))
}
