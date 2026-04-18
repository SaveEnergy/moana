package passwordtest

import (
	"testing"

	"golang.org/x/crypto/bcrypt"
)

// MustHash returns a bcrypt hash suitable for [moana/internal/store.Store.CreateUser] and other
// tests. It uses [bcrypt.MinCost] so the suite stays fast; production code still uses Argon2id via
// [moana/internal/auth.HashPassword].
func MustHash(t testing.TB, password string) []byte {
	t.Helper()
	b, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	if err != nil {
		t.Fatal(err)
	}
	return b
}
