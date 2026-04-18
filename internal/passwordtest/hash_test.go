package passwordtest

import (
	"testing"

	"moana/internal/auth"

	"golang.org/x/crypto/bcrypt"
)

func TestMustHash_verifiesWithBcryptAndAuthCheckPassword(t *testing.T) {
	t.Parallel()
	h := MustHash(t, "integration-secret")
	if err := bcrypt.CompareHashAndPassword(h, []byte("integration-secret")); err != nil {
		t.Fatalf("bcrypt compare: %v", err)
	}
	if err := auth.CheckPassword(h, "integration-secret"); err != nil {
		t.Fatalf("auth.CheckPassword: %v", err)
	}
	if err := auth.CheckPassword(h, "wrong"); err == nil {
		t.Fatal("expected mismatch")
	}
}
