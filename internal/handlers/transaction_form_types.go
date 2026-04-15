package handlers

import "moana/internal/store"

// txFormData is the new-entry form (GET /transactions, form errors).
type txFormData struct {
	Error         string
	Categories    []store.Category
	Today         string
	SelectedCatID int64
}

// txEditFormData is the edit form (GET/POST /transactions/{id}/edit).
type txEditFormData struct {
	Error         string
	Categories    []store.Category
	TxID          int64
	Kind          string
	Amount        string
	OccurredOn    string
	Description   string
	SelectedCatID int64
	Next          string
}
