package store

import (
	"context"
	"crypto/sha256"
	"errors"
	"testing"
	"time"

	"moana/internal/auth"
	"moana/internal/passwordtest"
)

func TestRedeemPasswordResetToken_success(t *testing.T) {
	t.Parallel()
	st := testStore(t)
	ctx := context.Background()
	hash0 := passwordtest.MustHash(t, "old-secret")
	uid, err := st.CreateUser(ctx, "reset-ok@example.com", hash0, "user")
	if err != nil {
		t.Fatal(err)
	}
	raw := "one-time-raw-token-for-test"
	h := sha256.Sum256([]byte(raw))
	now := time.Now().UTC()
	if err := st.ReplacePasswordResetToken(ctx, uid, h[:], now.Add(30*time.Minute), now); err != nil {
		t.Fatal(err)
	}
	newHash, err := auth.HashPassword("new-secret-99")
	if err != nil {
		t.Fatal(err)
	}
	if err := st.RedeemPasswordResetToken(ctx, h[:], newHash); err != nil {
		t.Fatal(err)
	}
	u, err := st.GetUserByEmail(ctx, "reset-ok@example.com")
	if err != nil {
		t.Fatal(err)
	}
	if err := auth.CheckPassword(u.PasswordHash, "new-secret-99"); err != nil {
		t.Fatalf("new password: %v", err)
	}
	if err := st.RedeemPasswordResetToken(ctx, h[:], newHash); !errors.Is(err, ErrPasswordResetInvalid) {
		t.Fatalf("second redeem: got %v want ErrPasswordResetInvalid", err)
	}
}

func TestRedeemPasswordResetToken_expired(t *testing.T) {
	t.Parallel()
	st := testStore(t)
	ctx := context.Background()
	hash0 := passwordtest.MustHash(t, "pw")
	uid, err := st.CreateUser(ctx, "reset-exp@example.com", hash0, "user")
	if err != nil {
		t.Fatal(err)
	}
	raw := "expired-token-raw"
	h := sha256.Sum256([]byte(raw))
	past := time.Now().UTC().Add(-1 * time.Hour)
	now := time.Now().UTC()
	if err := st.ReplacePasswordResetToken(ctx, uid, h[:], past, now); err != nil {
		t.Fatal(err)
	}
	newHash, err := auth.HashPassword("nope")
	if err != nil {
		t.Fatal(err)
	}
	if err := st.RedeemPasswordResetToken(ctx, h[:], newHash); !errors.Is(err, ErrPasswordResetInvalid) {
		t.Fatalf("got %v want ErrPasswordResetInvalid", err)
	}
}

func TestRedeemPasswordResetToken_unknown(t *testing.T) {
	t.Parallel()
	st := testStore(t)
	ctx := context.Background()
	var h [32]byte
	h[0] = 0xff
	nh, err := auth.HashPassword("x")
	if err != nil {
		t.Fatal(err)
	}
	if err := st.RedeemPasswordResetToken(ctx, h[:], nh); !errors.Is(err, ErrPasswordResetInvalid) {
		t.Fatalf("got %v want ErrPasswordResetInvalid", err)
	}
}
