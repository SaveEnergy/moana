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
	cfg := parseStatsPeriod(periodQuery)

	curStart, curEnd := timeutil.TrailingLocalDaysInclusiveRangeUTC(loc, now, cfg.InclusiveDays)
	prevStart, prevEnd := timeutil.PriorTrailingLocalDaysInclusiveRangeUTC(loc, now, cfg.InclusiveDays)

	running, periodIncome, periodExpense, prevPeriodIncome, prevPeriodExp, err := st.SumRunningTotalAndIncomeExpenseInTwoRanges(ctx, householdID, curStart, curEnd, prevStart, prevEnd)
	if err != nil {
		return PageData{}, err
	}
	periodNet := periodIncome + periodExpense
	prevPeriodNet := prevPeriodIncome + prevPeriodExp

	netVsPriorPct := NetPctChange(periodNet, prevPeriodNet)
	incomeTrendPct := PctChangePositive(periodIncome, prevPeriodIncome)
	expenseTrendPct := PctChangePositive(money.AbsCents(periodExpense), money.AbsCents(prevPeriodExp))

	var budgetUsedPct float64
	var budgetCapCents int64
	var budgetMeta string
	switch cfg.Period {
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
		StatsPeriod:          cfg.Period,
		StatsPriorPhrase:     cfg.PriorPhrase,
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
