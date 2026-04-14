package tmpl

import (
	"html/template"
	"strings"

	"moana/internal/category"
	"moana/internal/icons"
	"moana/internal/store"
)

func templateCategoryIconChoices() []icons.Def {
	return icons.CategoryPalette
}

// templateCategoryIconFormValue returns the normalized icon id for form matching (maps legacy emoji).
func templateCategoryIconFormValue(c store.Category) string {
	return category.NormalizeStoredIcon(c.Icon)
}

func categoryDisplayIcon(c store.Category, index int) template.HTML {
	id := resolveCategoryIconID(c.Icon, c.Name, index)
	return icons.SVG(id, "moana-icon moana-icon--sm")
}

func categoryNoneIconHTML() template.HTML {
	return icons.SVG("ban", "moana-icon moana-icon--sm")
}

func resolveCategoryIconID(stored, name string, index int) string {
	if x := category.NormalizeStoredIcon(stored); x != "" {
		return x
	}
	return category.PickerIcon(name, index)
}

// dashboardRecentIcon shows the category glyph when labeled; otherwise income/expense arrows.
func dashboardRecentIcon(tx store.Transaction, index int) template.HTML {
	name := strings.TrimSpace(tx.CategoryName)
	if name != "" {
		c := store.Category{Name: name, Icon: tx.CategoryIcon}
		return categoryDisplayIcon(c, index)
	}
	if tx.AmountCents < 0 {
		return icons.SVG("trending-down", "moana-icon moana-icon--sm")
	}
	return icons.SVG("trending-up", "moana-icon moana-icon--sm")
}
