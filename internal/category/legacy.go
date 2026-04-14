package category

import (
	"strings"

	"moana/internal/icons"
)

// legacyEmojiToIcon maps old DB emoji values to Lucide icon ids.
var legacyEmojiToIcon = map[string]string{
	"◫": "circle-dollar-sign", "🍽": "utensils", "🚗": "car", "⚕": "stethoscope",
	"🛍": "shopping-bag", "⌂": "house", "🎬": "film", "✈": "plane", "💰": "wallet",
	"📱": "smartphone", "☕": "coffee", "🛒": "shopping-cart", "🏥": "stethoscope",
	"💼": "briefcase", "🎓": "graduation-cap", "🐕": "dog", "✏": "square-pen",
	"◇": "gem", "○": "circle-dollar-sign", "◈": "gem", "△": "percent", "□": "calculator",
	"⬡": "snowflake", "🛠": "hammer", "🏡": "house", "⚡": "flame", "🌿": "sun",
	"🎁": "gift", "🏦": "landmark", "💳": "credit-card", "📊": "percent", "🧾": "receipt",
}

// NormalizeStoredIcon returns a Lucide icon id if s is a known id or legacy emoji, else "" (auto / unknown).
func NormalizeStoredIcon(s string) string {
	s = strings.TrimSpace(s)
	if s == "" {
		return ""
	}
	if icons.ValidID(s) {
		return s
	}
	if x := legacyEmojiToIcon[s]; x != "" {
		return x
	}
	return ""
}
