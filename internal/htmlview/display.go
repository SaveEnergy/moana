package htmlview

import (
	"strings"
	"unicode"
	"unicode/utf8"
)

// UserInitial returns one uppercase letter from email, or "?".
func UserInitial(email string) string {
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

// AdminDisplayName derives a short display name from an email local part.
func AdminDisplayName(email string) string {
	email = strings.TrimSpace(strings.ToLower(email))
	at := strings.Index(email, "@")
	local := email
	if at > 0 {
		local = email[:at]
	}
	if local == "" {
		return strings.TrimSpace(email)
	}
	local = strings.ReplaceAll(local, ".", " ")
	local = strings.ReplaceAll(local, "_", " ")
	words := strings.Fields(local)
	for i, w := range words {
		if w == "" {
			continue
		}
		runes := []rune(w)
		runes[0] = unicode.ToUpper(runes[0])
		for j := 1; j < len(runes); j++ {
			runes[j] = unicode.ToLower(runes[j])
		}
		words[i] = string(runes)
	}
	return strings.Join(words, " ")
}

// AdminRoleLabel maps a role string to a display label.
func AdminRoleLabel(role string) string {
	switch strings.ToLower(strings.TrimSpace(role)) {
	case "admin":
		return "Admin"
	default:
		return "Member"
	}
}

// ProfileDisplayName prefers first+last, else falls back to email-based display.
func ProfileDisplayName(first, last, email string) string {
	f := strings.TrimSpace(first)
	l := strings.TrimSpace(last)
	if f != "" || l != "" {
		return strings.TrimSpace(f + " " + l)
	}
	return AdminDisplayName(email)
}

// ProfileInitial prefers first name initial, then last, then email initial.
func ProfileInitial(first, last, email string) string {
	f := strings.TrimSpace(first)
	l := strings.TrimSpace(last)
	if f != "" {
		r, _ := utf8.DecodeRuneInString(strings.ToUpper(f))
		return string(r)
	}
	if l != "" {
		r, _ := utf8.DecodeRuneInString(strings.ToUpper(l))
		return string(r)
	}
	return UserInitial(email)
}

// HouseholdRoleLabel maps household membership role to a display label.
func HouseholdRoleLabel(role string) string {
	switch strings.ToLower(strings.TrimSpace(role)) {
	case "owner":
		return "Owner"
	case "admin":
		return "Admin"
	default:
		return "Member"
	}
}
