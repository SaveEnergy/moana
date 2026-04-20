package dashboard

import (
	"context"
	"fmt"
	"math"
	"time"

	"golang.org/x/sync/errgroup"

	"moana/internal/money"
	"moana/internal/store"
	"moana/internal/timeutil"
)

// BuildPageData loads aggregates and layout data for the dashboard (no HTTP).
// The running-total / period aggregate query runs concurrently with the heatmap and recent-transaction
// reads (independent work; [database/sql] pool allows overlapping reads under WAL). Outflow breakdown
// runs after aggregates complete because it needs periodExpense from the aggregate scan.
func BuildPageData(ctx context.Context, st *store.Store, householdID int64, loc *time.Location, now time.Time, periodQuery string) (PageData, error) {
	cfg := parseStatsPeriod(periodQuery)

	curStart, curEnd := timeutil.TrailingLocalDaysInclusiveRangeUTC(loc, now, cfg.InclusiveDays)
	prevStart, prevEnd := timeutil.PriorTrailingLocalDaysInclusiveRangeUTC(loc, now, cfg.InclusiveDays)

	var (
		running, periodIncome, periodExpense, prevPeriodIncome, prevPeriodExp int64
		heatmapRangeLabel                                                     string
		heatmapCells                                                          []HeatmapCell
		heatmapCols                                                           int
		recent                                                                []store.Transaction
	)
	g, gctx := errgroup.WithContext(ctx)
	g.Go(func() error {
		var err error
		running, periodIncome, periodExpense, prevPeriodIncome, prevPeriodExp, err = st.SumRunningTotalAndIncomeExpenseInTwoRanges(gctx, householdID, curStart, curEnd, prevStart, prevEnd)
		return err
	})
	g.Go(func() error {
		var err error
		heatmapRangeLabel, heatmapCells, heatmapCols, err = buildHeatmapSection(gctx, st, householdID, loc, now)
		return err
	})
	g.Go(func() error {
		var err error
		recent, err = st.ListTransactions(gctx, householdID, store.TransactionFilter{Limit: 5})
		return err
	})
	if err := g.Wait(); err != nil {
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
	case StatsPeriod12m:
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
