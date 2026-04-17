package auth

import (
	"testing"

	"golang.org/x/crypto/bcrypt"
)

func TestHashPasswordArgon2idRoundTrip(t *testing.T) {
	hash, err := HashPassword("correct horse battery staple")
	if err != nil {
		t.Fatal(err)
	}
	if err := CheckPassword(hash, "correct horse battery staple"); err != nil {
		t.Fatalf("verify: %v", err)
	}
	if err := CheckPassword(hash, "wrong"); err == nil {
		t.Fatal("expected error for wrong password")
	}
}

func TestCheckPassword_emptyHash(t *testing.T) {
	t.Parallel()
	if err := CheckPassword(nil, "x"); err == nil {
		t.Fatal("nil hash must not verify")
	}
	if err := CheckPassword([]byte{}, "x"); err == nil {
		t.Fatal("empty hash must not verify")
	}
}

func TestCheckPasswordBcryptLegacy(t *testing.T) {
	legacy, err := bcrypt.GenerateFromPassword([]byte("legacy-secret"), 10)
	if err != nil {
		t.Fatal(err)
	}
	if err := CheckPassword(legacy, "legacy-secret"); err != nil {
		t.Fatalf("bcrypt verify: %v", err)
	}
	if err := CheckPassword(legacy, "nope"); err == nil {
		t.Fatal("expected mismatch")
	}
}
