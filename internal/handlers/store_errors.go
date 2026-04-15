package handlers

import (
	"errors"

	"moana/internal/httperr"
	"moana/internal/store"
)

// userFacingStoreMessage maps store errors to short UI copy.
// Unknown errors use [httperr.InternalMessage] so driver/SQL text never reaches the HTML alert.
func userFacingStoreMessage(err error) string {
	if errors.Is(err, store.ErrInvalidCategory) {
		return "That category is not valid for this household."
	}
	if errors.Is(err, store.ErrCategoryNotFound) {
		return "That category could not be found."
	}
	if errors.Is(err, store.ErrDuplicateCategoryName) {
		return "A category with that name already exists."
	}
	return httperr.InternalMessage
}
