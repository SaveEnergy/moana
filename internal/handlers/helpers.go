package handlers

import (
	"fmt"
	"math"
	"net/http"
	"strings"
	"time"
	"unicode/utf8"

	"moana/internal/auth"
	"moana/internal/money"
	"moana/internal/store"
	"moana/internal/timeutil"
)

func templateFormatEUR(cents int64) string {
	return money.FormatEUR(cents)
}

func templateFormatLocal(t time.Time, tz string) string {
	loc := timeutil.LoadLocation(tz)
	return t.In(loc).Format("2006-01-02 15:04")
}

// templateFormatLocalTime formats clock time in the user's zone (e.g. for history rows).
func templateFormatLocalTime(t time.Time, tz string) string {
	loc := timeutil.LoadLocation(tz)
	return t.In(loc).Format("3:04 PM")
}

func templateFormatEURAbs(cents int64) string {
	if cents < 0 {
		cents = -cents
	}
	return money.FormatEUR(cents)
}

func templateUserInitial(email string) string {
	email = strings.TrimSpace(email)
	if email == "" {
		return "?"
	}
	r, _ := utf8.DecodeRuneInString(email)
	if r == utf8.RuneError {
		return "?"
	}
	return strings.ToUpper(string(r))
}

func templateIsNegFloat(x float64) bool {
	return !math.IsNaN(x) && !math.IsInf(x, 0) && x < 0
}

func templatePercentSigned(x float64) string {
	if math.IsNaN(x) || math.IsInf(x, 0) {
		return "—"
	}
	sign := ""
	if x >= 0 {
		sign = "+"
	}
	return fmt.Sprintf("%s%.1f%%", sign, x)
}

func templateFormatCompactEUR(cents int64) string {
	x := cents
	if x < 0 {
		x = -x
	}
	if x < 100_000 {
		return money.FormatEUR(cents)
	}
	v := float64(x) / 100.0 / 1000.0
	return fmt.Sprintf("€%.1fk", v)
}

func (a *App) currentUser(r *http.Request) (*store.User, error) {
	sess, err := auth.ReadSession(r, a.Config.SessionSecret)
	if err != nil || sess == nil {
		return nil, err
	}
	ctx := r.Context()
	u, err := a.Store.GetUserByID(ctx, sess.UserID)
	if err != nil {
		return nil, err
	}
	if u == nil {
		return nil, fmt.Errorf("user missing")
	}
	if u.Role != sess.Role {
		// role changed server-side; treat as logout
		return nil, fmt.Errorf("stale session")
	}
	return u, nil
}
