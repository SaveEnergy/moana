package dashboard

import (
	"moana/internal/store"
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
