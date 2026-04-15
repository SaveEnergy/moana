package auth

import (
	"strings"

	"golang.org/x/crypto/bcrypt"
)

// CheckPassword compares a stored hash (Argon2id or legacy bcrypt) with the plaintext password.
func CheckPassword(hash []byte, password string) error {
	if len(hash) == 0 {
		return bcrypt.ErrMismatchedHashAndPassword
	}
	h := string(hash)
	if strings.HasPrefix(h, "$argon2id$") {
		return checkArgon2id(h, password)
	}
	return bcrypt.CompareHashAndPassword(hash, []byte(password))
}
