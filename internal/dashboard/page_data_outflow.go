package dashboard

import (
	"context"
	"time"

	"moana/internal/category"
	"moana/internal/money"
	"moana/internal/store"
)

func buildOutflowSection(ctx context.Context, st *store.Store, householdID int64, curStart, curEnd time.Time, periodExpense int64) ([]OutflowRow, string, int64, error) {
	expenseRows, err := st.ListCategoryAmountsInRange(ctx, householdID, &curStart, &curEnd, "expense")
	if err != nil {
		return nil, "", 0, err
	}
	outflowMerged := MergeCategoryTopN(expenseRows, defaultOutflowMergeLimit)
	totalAbs := money.AbsCents(periodExpense)
	var outflowRows []OutflowRow
	var pcts []float64
	var hexes []string
	for _, ca := range outflowMerged {
		if totalAbs <= 0 {
			break
		}
		p := float64(ca.AmountCents) / float64(totalAbs) * 100
		cat := store.Category{Name: ca.Name, Icon: ca.Icon, Color: ca.Color}
		outflowRows = append(outflowRows, OutflowRow{
			Category:    cat,
			AmountCents: ca.AmountCents,
			Pct:         p,
		})
		pcts = append(pcts, p)
		hexes = append(hexes, category.HexOrDefault(cat))
	}
	outflowDonut := DonutConicGradient(pcts, hexes)
	return outflowRows, outflowDonut, totalAbs, nil
}
