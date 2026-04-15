package category

import (
	"context"

	"moana/internal/store"
)

// CategoriesListData is the template payload for the categories list page.
type CategoriesListData struct {
	Error      string
	Categories []store.Category
}

// BuildCategoriesList loads categories; errMsg is set when re-rendering after a validation or store error.
func BuildCategoriesList(ctx context.Context, st *store.Store, householdID int64, errMsg string) (CategoriesListData, error) {
	cats, err := st.ListCategories(ctx, householdID)
	if err != nil {
		return CategoriesListData{}, err
	}
	return CategoriesListData{Error: errMsg, Categories: cats}, nil
}
