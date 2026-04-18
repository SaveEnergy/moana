package handlers

import (
	"errors"

	"moana/internal/httperr"
	"moana/internal/store"
)

// userFacingStoreMessage maps store errors to short UI copy.
// Unknown errors use [httperr.InternalMessage] so driver/SQL text never reaches the HTML alert.
func userFacingStoreMessage(err error) string {
	for _, row := range storeUserMessages {
		if errors.Is(err, row.sentinel) {
			return row.msg
		}
	}
	return httperr.InternalMessage
}

// storeUserMessages is ordered: first matching sentinel wins ([errors.Is] chain).
var storeUserMessages = []struct {
	sentinel error
	msg      string
}{
	{store.ErrInvalidCategory, "That category is not valid for this household."},
	{store.ErrCategoryNotFound, "That category could not be found."},
	{store.ErrDuplicateCategoryName, "A category with that name already exists."},
	{store.ErrDuplicateUserEmail, "A user with that email already exists."},
	{store.ErrInvalidUserEmail, "That email address is not valid."},
}
