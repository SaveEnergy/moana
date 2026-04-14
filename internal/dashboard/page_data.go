package dashboard

import (
	"context"
	"fmt"
	"math"
	"time"

	"moana/internal/category"
	"moana/internal/money"
	"moana/internal/store"
	"moana/internal/timeutil"
)

// DefaultMonthlyExpenseBudgetCents is a display-only cap for the budget KPI until user settings exist.
const DefaultMonthlyExpenseBudgetCents = 5_000_00

const defaultOutflowMergeLimit = 10

// OutflowRow is one slice of the expense donut + table.
type OutflowRow struct {
	Category    store.Category
	AmountCents int64
	Pct         float64
}

// PageData is the template payload for the dashboard overview.
type PageData struct {
	StatsPeriod          string
	StatsPriorPhrase     string
	RunningTotal         int64
	MonthIncome          int64
	MonthExpense         int64
	MonthNet             int64
	NetVsPriorPct        float64
	IncomeTrendPct       float64
	ExpenseTrendPct      float64
	BudgetUsedPct        float64
	BudgetBarPct         float64
	BudgetMeta           string
	OutflowRows          []OutflowRow
	OutflowDonutGradient string
	OutflowTotalAbs      int64
	HeatmapRangeLabel    string
	HeatmapCells         []HeatmapCell
	HeatmapColCount      int
	Recent               []store.Transaction
}

// BuildPageData loads aggregates and layout data for the dashboard (no HTTP).
func BuildPageData(ctx context.Context, st *store.Store, userID int64, loc *time.Location, now time.Time, periodQuery string) (PageData, error) {
	statsPeriod := periodQuery
	inclusiveDays := 30
	switch statsPeriod {
	case "12m":
		inclusiveDays = 365
	case "30d", "":
		statsPeriod = "30d"
	default:
		statsPeriod = "30d"
	}

	statsPriorPhrase := "prior 30 days"
	if statsPeriod == "12m" {
		statsPriorPhrase = "prior 12 months"
	}

	curStart, curEnd := timeutil.TrailingLocalDaysInclusiveRangeUTC(loc, now, inclusiveDays)
	prevStart, prevEnd := timeutil.PriorTrailingLocalDaysInclusiveRangeUTC(loc, now, inclusiveDays)

	running, err := st.SumAmountCents(ctx, userID, nil, nil)
	if err != nil {
		return PageData{}, err
	}

	periodIncome, err := st.SumAmountCentsByKind(ctx, userID, &curStart, &curEnd, "income")
	if err != nil {
		return PageData{}, err
	}
	periodExpense, err := st.SumAmountCentsByKind(ctx, userID, &curStart, &curEnd, "expense")
	if err != nil {
		return PageData{}, err
	}
	periodNet := periodIncome + periodExpense

	prevPeriodNet, err := st.SumAmountCents(ctx, userID, &prevStart, &prevEnd)
	if err != nil {
		return PageData{}, err
	}

	prevPeriodExp, err := st.SumAmountCentsByKind(ctx, userID, &prevStart, &prevEnd, "expense")
	if err != nil {
		return PageData{}, err
	}
	prevPeriodIncome, err := st.SumAmountCentsByKind(ctx, userID, &prevStart, &prevEnd, "income")
	if err != nil {
		return PageData{}, err
	}

	netVsPriorPct := NetPctChange(periodNet, prevPeriodNet)
	incomeTrendPct := PctChangePositive(periodIncome, prevPeriodIncome)
	expenseTrendPct := PctChangePositive(money.AbsCents(periodExpense), money.AbsCents(prevPeriodExp))

	var budgetUsedPct float64
	var budgetCapCents int64
	var budgetMeta string
	switch statsPeriod {
	case "12m":
		budgetCapCents = DefaultMonthlyExpenseBudgetCents * 12
		budgetMeta = fmt.Sprintf("of %s annual budget (12× monthly) used", money.FormatEUR(budgetCapCents))
	default:
		budgetCapCents = DefaultMonthlyExpenseBudgetCents
		budgetMeta = fmt.Sprintf("of %s monthly budget used", money.FormatEUR(DefaultMonthlyExpenseBudgetCents))
	}
	if budgetCapCents > 0 {
		budgetUsedPct = float64(money.AbsCents(periodExpense)) / float64(budgetCapCents) * 100
	}
	budgetBarPct := math.Min(100, budgetUsedPct)

	expenseRows, err := st.ListCategoryAmountsInRange(ctx, userID, &curStart, &curEnd, "expense")
	if err != nil {
		return PageData{}, err
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

	todayLocal := time.Date(now.In(loc).Year(), now.In(loc).Month(), now.In(loc).Day(), 0, 0, 0, 0, loc)
	yearAgo := todayLocal.AddDate(0, 0, -364)
	startUTC := yearAgo.UTC()
	endOfToday := time.Date(todayLocal.Year(), todayLocal.Month(), todayLocal.Day(), 23, 59, 59, 999999999, todayLocal.Location())
	endUTC := endOfToday.UTC()
	heatmapRangeLabel := fmt.Sprintf("%s – %s", yearAgo.Format("Jan 2, 2006"), todayLocal.Format("Jan 2, 2006"))

	byDay, err := st.DailyAbsMovementByLocalDate(ctx, userID, startUTC, endUTC, loc)
	if err != nil {
		return PageData{}, err
	}
	heatmapCells := BuildHeatmapCellsRolling365(todayLocal, loc, byDay)
	heatmapCols := 1
	if n := len(heatmapCells); n > 0 {
		heatmapCols = (n + 6) / 7
	}

	recent, err := st.ListTransactions(ctx, userID, store.TransactionFilter{Limit: 5})
	if err != nil {
		return PageData{}, err
	}

	return PageData{
		StatsPeriod:          statsPeriod,
		StatsPriorPhrase:     statsPriorPhrase,
		RunningTotal:         running,
		MonthIncome:          periodIncome,
		MonthExpense:         periodExpense,
		MonthNet:             periodNet,
		NetVsPriorPct:        netVsPriorPct,
		IncomeTrendPct:       incomeTrendPct,
		ExpenseTrendPct:      expenseTrendPct,
		BudgetUsedPct:        budgetUsedPct,
		BudgetBarPct:         budgetBarPct,
		BudgetMeta:           budgetMeta,
		OutflowRows:          outflowRows,
		OutflowDonutGradient: outflowDonut,
		OutflowTotalAbs:      totalAbs,
		HeatmapRangeLabel:    heatmapRangeLabel,
		HeatmapCells:         heatmapCells,
		HeatmapColCount:      heatmapCols,
		Recent:               recent,
	}, nil
}
