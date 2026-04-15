package dashboard

import (
	"context"
	"fmt"
	"math"
	"time"

	"moana/internal/money"
	"moana/internal/store"
	"moana/internal/timeutil"
)

// BuildPageData loads aggregates and layout data for the dashboard (no HTTP).
func BuildPageData(ctx context.Context, st *store.Store, householdID int64, loc *time.Location, now time.Time, periodQuery string) (PageData, error) {
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

	running, err := st.SumAmountCents(ctx, householdID, nil, nil)
	if err != nil {
		return PageData{}, err
	}

	periodIncome, err := st.SumAmountCentsByKind(ctx, householdID, &curStart, &curEnd, "income")
	if err != nil {
		return PageData{}, err
	}
	periodExpense, err := st.SumAmountCentsByKind(ctx, householdID, &curStart, &curEnd, "expense")
	if err != nil {
		return PageData{}, err
	}
	periodNet := periodIncome + periodExpense

	prevPeriodNet, err := st.SumAmountCents(ctx, householdID, &prevStart, &prevEnd)
	if err != nil {
		return PageData{}, err
	}

	prevPeriodExp, err := st.SumAmountCentsByKind(ctx, householdID, &prevStart, &prevEnd, "expense")
	if err != nil {
		return PageData{}, err
	}
	prevPeriodIncome, err := st.SumAmountCentsByKind(ctx, householdID, &prevStart, &prevEnd, "income")
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

	outflowRows, outflowDonut, totalAbs, err := buildOutflowSection(ctx, st, householdID, curStart, curEnd, periodExpense)
	if err != nil {
		return PageData{}, err
	}

	heatmapRangeLabel, heatmapCells, heatmapCols, err := buildHeatmapSection(ctx, st, householdID, loc, now)
	if err != nil {
		return PageData{}, err
	}

	recent, err := st.ListTransactions(ctx, householdID, store.TransactionFilter{Limit: 5})
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
