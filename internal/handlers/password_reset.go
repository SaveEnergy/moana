package handlers

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"moana/internal/auth"
	"moana/internal/httperr"
	"moana/internal/mail"
)

const minResetPasswordLen = 8

// forgotPasswordData is used by forgot_password.html.
type forgotPasswordData struct {
	Title              string
	Error              string
	Info               string
	Year               int
	RepoURL            string
	Unavailable        bool
	LoginPath          string
	ForgotPasswordPath string
}

// resetPasswordData is used by reset_password.html. Token is the hidden one-time value; empty hides the form.
type resetPasswordData struct {
	Title             string
	Error             string
	Year              int
	RepoURL           string
	Token             string
	LoginPath         string
	ResetPasswordPath string
}

// ForgotPasswordPage shows the request form (GET /forgot-password).
func (a *App) ForgotPasswordPage(w http.ResponseWriter, r *http.Request) {
	_, err := a.CurrentUser(r)
	if err == nil {
		http.Redirect(w, r, DashboardPath, http.StatusSeeOther)
		return
	}
	if !errors.Is(err, ErrAuthRequired) {
		httperr.Internal(w, r, err)
		return
	}
	if a.PasswordReset == nil {
		a.renderSimple(w, r, "forgot_password.html", a.forgotTemplate("", "", true, r.URL.Query().Get(PasswordResetQuerySent) == "1"))
		return
	}
	a.renderSimple(w, r, "forgot_password.html", a.forgotTemplate("", "", false, r.URL.Query().Get(PasswordResetQuerySent) == "1"))
}

// ForgotPasswordSubmit issues a reset link when outbound mail is enabled (POST /forgot-password).
func (a *App) ForgotPasswordSubmit(w http.ResponseWriter, r *http.Request) {
	_, err := a.CurrentUser(r)
	if err == nil {
		http.Redirect(w, r, DashboardPath, http.StatusSeeOther)
		return
	}
	if !errors.Is(err, ErrAuthRequired) {
		httperr.Internal(w, r, err)
		return
	}
	if !requireParseForm(w, r) {
		return
	}
	if a.PasswordReset == nil {
		a.renderSimple(w, r, "forgot_password.html", a.forgotTemplate("", "", true, false))
		return
	}
	email := strings.TrimSpace(r.FormValue(LoginFieldEmail))
	if email == "" {
		redir := ForgotPasswordPath + "?" + PasswordResetQuerySent + "=1"
		http.Redirect(w, r, redir, http.StatusSeeOther)
		return
	}
	ctx := r.Context()
	u, err := a.Store.GetUserByEmail(ctx, email)
	if err != nil {
		httperr.Internal(w, r, err)
		return
	}
	if u != nil {
		raw, hash, err := newPasswordResetRawAndHash()
		if err != nil {
			httperr.Internal(w, r, err)
			return
		}
		now := time.Now().UTC()
		exp := now.Add(a.Config.PasswordResetTTL)
		if a.Config.PasswordResetTTL <= 0 {
			exp = now.Add(60 * time.Minute)
		}
		if err := a.Store.ReplacePasswordResetToken(ctx, u.ID, hash[:], exp, now); err != nil {
			httperr.Internal(w, r, err)
			return
		}
		q := url.Values{}
		q.Set(ResetFieldToken, raw)
		resetURL := mail.PublicResetURL(a.Config.PublicBaseURL, ResetPasswordPath+"?"+q.Encode())
		if err := a.PasswordReset.SendPasswordReset(ctx, u.Email, resetURL); err != nil {
			log.Printf("password reset email: %v", err)
		}
	}
	redir := ForgotPasswordPath + "?" + PasswordResetQuerySent + "=1"
	http.Redirect(w, r, redir, http.StatusSeeOther)
}

// ResetPasswordPage shows the new-password form (GET /reset-password?token=).
func (a *App) ResetPasswordPage(w http.ResponseWriter, r *http.Request) {
	_, err := a.CurrentUser(r)
	if err == nil {
		http.Redirect(w, r, DashboardPath, http.StatusSeeOther)
		return
	}
	if !errors.Is(err, ErrAuthRequired) {
		httperr.Internal(w, r, err)
		return
	}
	token := strings.TrimSpace(r.URL.Query().Get(ResetFieldToken))
	if token == "" {
		a.renderSimple(w, r, "reset_password.html", a.resetTemplate("This reset link is missing the token. Open the full link from your email.", ""))
		return
	}
	if a.PasswordReset == nil {
		a.renderSimple(w, r, "reset_password.html", a.resetTemplate("Password reset is not available on this server (outbound email is not configured).", ""))
		return
	}
	a.renderSimple(w, r, "reset_password.html", a.resetTemplate("", token))
}

// ResetPasswordSubmit applies a new password (POST /reset-password).
func (a *App) ResetPasswordSubmit(w http.ResponseWriter, r *http.Request) {
	_, err := a.CurrentUser(r)
	if err == nil {
		http.Redirect(w, r, DashboardPath, http.StatusSeeOther)
		return
	}
	if !errors.Is(err, ErrAuthRequired) {
		httperr.Internal(w, r, err)
		return
	}
	if !requireParseForm(w, r) {
		return
	}
	if a.PasswordReset == nil {
		a.renderSimple(w, r, "reset_password.html", a.resetTemplate("Password reset is not available on this server (outbound email is not configured).", r.FormValue(ResetFieldToken)))
		return
	}
	token := strings.TrimSpace(r.FormValue(ResetFieldToken))
	pw := r.FormValue(LoginFieldPassword)
	confirm := r.FormValue(ResetFieldPasswordConfirm)
	if token == "" {
		a.renderSimple(w, r, "reset_password.html", a.resetTemplate("Missing token. Open the full link from your email.", ""))
		return
	}
	if len(pw) < minResetPasswordLen {
		a.renderSimple(w, r, "reset_password.html", a.resetTemplate("Password must be at least 8 characters.", token))
		return
	}
	if pw != confirm {
		a.renderSimple(w, r, "reset_password.html", a.resetTemplate("Passwords do not match.", token))
		return
	}
	h := sha256.Sum256([]byte(token))
	newHash, err := auth.HashPassword(pw)
	if err != nil {
		a.renderSimple(w, r, "reset_password.html", a.resetTemplate("Could not set that password. Try a different one.", token))
		return
	}
	ctx := r.Context()
	if err := a.Store.RedeemPasswordResetToken(ctx, h[:], newHash); err != nil {
		a.renderSimple(w, r, "reset_password.html", a.resetTemplate(userFacingStoreMessage(err), token))
		return
	}
	http.Redirect(w, r, LoginPath, http.StatusSeeOther)
}

func (a *App) forgotTemplate(err string, info string, unavailable, sent bool) forgotPasswordData {
	if info == "" && sent {
		info = "If that address is registered, a reset link is on the way. Check your inbox and spam."
	}
	return forgotPasswordData{
		Title:              "Forgot password",
		Error:              err,
		Info:               info,
		Year:               time.Now().UTC().Year(),
		RepoURL:            a.Config.RepoURL,
		Unavailable:        unavailable,
		LoginPath:          LoginPath,
		ForgotPasswordPath: ForgotPasswordPath,
	}
}

func (a *App) resetTemplate(err, token string) resetPasswordData {
	return resetPasswordData{
		Title:             "Set a new password",
		Error:             err,
		Year:              time.Now().UTC().Year(),
		RepoURL:           a.Config.RepoURL,
		Token:             token,
		LoginPath:         LoginPath,
		ResetPasswordPath: ResetPasswordPath,
	}
}

func newPasswordResetRawAndHash() (raw string, sum [32]byte, err error) {
	b := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		return "", [32]byte{}, err
	}
	raw = base64.RawURLEncoding.EncodeToString(b)
	sum = sha256.Sum256([]byte(raw))
	return raw, sum, nil
}
