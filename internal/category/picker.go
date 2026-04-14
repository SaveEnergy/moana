package category

import (
	"strings"

	"moana/internal/store"
)

// pickerAccents when no custom color is set — warm/teal/coral/indigo.
var pickerAccents = []string{
	"#ea580c", "#0d9488", "#e11d48", "#6366f1", "#ca8a04", "#7c3aed",
	"#0ea5e9", "#059669", "#db2777", "#4f46e5",
}

// PickerAccent resolves the display accent for a category row/card.
func PickerAccent(c store.Category, index int) string {
	if col := strings.TrimSpace(c.Color); col != "" {
		if x := NormalizeColor(col); x != "" {
			return x
		}
	}
	if c.ID < 1 {
		return pickerAccents[index%len(pickerAccents)]
	}
	u := uint64(c.ID) % uint64(len(pickerAccents))
	return pickerAccents[u]
}

// PickerIcon chooses a Lucide icon id from a free-text category name (auto mode).
func PickerIcon(name string, index int) string {
	n := strings.ToLower(strings.TrimSpace(name))
	switch {
	case strings.Contains(n, "food") || strings.Contains(n, "dining") || strings.Contains(n, "grocery") || strings.Contains(n, "restaurant"):
		return "utensils"
	case strings.Contains(n, "transport") || strings.Contains(n, "commute") || strings.Contains(n, "fuel") || strings.Contains(n, "car") || strings.Contains(n, "parking"):
		return "car"
	case strings.Contains(n, "health") || strings.Contains(n, "medical") || strings.Contains(n, "pharmacy"):
		return "stethoscope"
	case strings.Contains(n, "shop") || strings.Contains(n, "retail") || strings.Contains(n, "clothes"):
		return "shopping-bag"
	case strings.Contains(n, "home") || strings.Contains(n, "rent") || strings.Contains(n, "housing") || strings.Contains(n, "utilities"):
		return "house"
	case strings.Contains(n, "entertain") || strings.Contains(n, "game") || strings.Contains(n, "stream"):
		return "film"
	case strings.Contains(n, "travel") || strings.Contains(n, "flight") || strings.Contains(n, "hotel"):
		return "plane"
	case strings.Contains(n, "salary") || strings.Contains(n, "income") || strings.Contains(n, "freelance") || strings.Contains(n, "payroll"):
		return "wallet"
	case strings.Contains(n, "save") || strings.Contains(n, "invest"):
		return "piggy-bank"
	default:
		fb := []string{"circle-dollar-sign", "gem", "percent", "landmark", "sparkles", "shopping-cart", "briefcase", "heart"}
		return fb[index%len(fb)]
	}
}

// PickerHint is a short subtitle derived from the category name (auto mode).
func PickerHint(name string) string {
	n := strings.ToLower(strings.TrimSpace(name))
	tests := []struct {
		kw, hint string
	}{
		{"food", "Dining & groceries"},
		{"grocery", "Dining & groceries"},
		{"dining", "Dining & groceries"},
		{"restaurant", "Dining & groceries"},
		{"transport", "Commute & fuel"},
		{"commute", "Commute & fuel"},
		{"fuel", "Commute & fuel"},
		{"car", "Commute & fuel"},
		{"health", "Medical & care"},
		{"medical", "Medical & care"},
		{"shop", "Purchases & orders"},
		{"retail", "Purchases & orders"},
		{"home", "Home & utilities"},
		{"rent", "Home & utilities"},
		{"housing", "Home & utilities"},
		{"entertain", "Media & leisure"},
		{"travel", "Trips & stays"},
		{"salary", "Pay & deposits"},
		{"income", "Pay & deposits"},
		{"freelance", "Pay & deposits"},
		{"save", "Savings & goals"},
		{"invest", "Savings & goals"},
	}
	for _, t := range tests {
		if strings.Contains(n, t.kw) {
			return t.hint
		}
	}
	return "Activity in this category"
}
