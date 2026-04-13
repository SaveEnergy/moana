package handlers

import (
	"fmt"
	"html/template"
	"math"
	"net/http"
	"strings"
	"time"

	"moana/internal/money"
	"moana/internal/store"
	"moana/internal/timeutil"
)

// Default savings goal for the progress bar (no user settings yet).
const savingsGoalCents = 12_000_00

type homeCategorySlice struct {
	Name string
	Pct  float64
}

// Home shows portfolio-style overview: balances, trends, charts, recent activity.
func (a *App) Home(w http.ResponseWriter, r *http.Request, u *store.User) {
	ctx := r.Context()
	loc := timeutil.LoadLocation(u.Timezone)
	now := time.Now().UTC()

	cmStart, cmEnd := timeutil.CalendarMonthRangeUTC(loc, now, 0)
	pmStart, pmEnd := timeutil.CalendarMonthRangeUTC(loc, now, 1)
	ppmStart, ppmEnd := timeutil.CalendarMonthRangeUTC(loc, now, 2)
	ytdStart, ytdEnd := timeutil.CurrentCalendarYearToDateRangeUTC(loc, now)

	running, err := a.Store.SumAmountCents(ctx, u.ID, nil, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	monthIncome, err := a.Store.SumAmountCentsByKind(ctx, u.ID, &cmStart, &cmEnd, "income")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	monthExpense, err := a.Store.SumAmountCentsByKind(ctx, u.ID, &cmStart, &cmEnd, "expense")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	monthNet := monthIncome + monthExpense

	prevMonthNet, err := a.Store.SumAmountCents(ctx, u.ID, &pmStart, &pmEnd)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	prevPrevMonthNet, err := a.Store.SumAmountCents(ctx, u.ID, &ppmStart, &ppmEnd)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	prevMonthExp, err := a.Store.SumAmountCentsByKind(ctx, u.ID, &pmStart, &pmEnd, "expense")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	prevPrevMonthExp, err := a.Store.SumAmountCentsByKind(ctx, u.ID, &ppmStart, &ppmEnd, "expense")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	ytdNet, err := a.Store.SumAmountCents(ctx, u.ID, &ytdStart, &ytdEnd)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	momNetPct := pctChangeNet(prevMonthNet, prevPrevMonthNet)
	expenseTrendPct := pctChangeExpense(prevMonthExp, prevPrevMonthExp)

	savingsProgress := 0
	if savingsGoalCents > 0 && monthNet > 0 {
		savingsProgress = int(math.Min(100, float64(monthNet)/float64(savingsGoalCents)*100))
	}

	incomeRows, err := a.Store.ListCategoryAmountsInRange(ctx, u.ID, &cmStart, &cmEnd, "income")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	expenseRows, err := a.Store.ListCategoryAmountsInRange(ctx, u.ID, &cmStart, &cmEnd, "expense")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	incomeMerged := mergeCategoryTopN(incomeRows, cashflowMergeLimit)
	expenseMerged := mergeCategoryTopN(expenseRows, cashflowMergeLimit)
	teAbs := abs64(monthExpense)
	var cashflowSVG template.HTML
	if monthIncome > 0 {
		cashflowSVG = buildCashflowSVG(incomeMerged, expenseMerged, monthIncome, teAbs)
	}

	cats, err := a.Store.ListTopExpenseCategories(ctx, u.ID, &cmStart, &cmEnd, 5)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	totalAbs := abs64(monthExpense)
	var slices []homeCategorySlice
	var pcts []float64
	sumCat := int64(0)
	for _, c := range cats {
		a := abs64(c.TotalCents)
		sumCat += a
		if totalAbs > 0 {
			p := float64(a) / float64(totalAbs) * 100
			slices = append(slices, homeCategorySlice{Name: c.CategoryName, Pct: p})
			pcts = append(pcts, p)
		}
	}
	if rest := totalAbs - sumCat; rest > 0 && totalAbs > 0 {
		p := float64(rest) / float64(totalAbs) * 100
		slices = append(slices, homeCategorySlice{Name: "Other", Pct: p})
		pcts = append(pcts, p)
	}
	donutGrad := donutGradient(pcts)

	var insight string
	if len(cats) > 0 && totalAbs > 0 {
		suggest := abs64(cats[0].TotalCents) / 10
		if suggest < 1 {
			suggest = 1
		}
		insight = fmt.Sprintf("Reduce “%s” by about %s to free room in your monthly budget.",
			cats[0].CategoryName, money.FormatEUR(suggest))
	}

	recent, err := a.Store.ListTransactions(ctx, u.ID, store.TransactionFilter{Limit: 5})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data := struct {
		RunningTotal     int64
		Timezone         string
		MonthIncome      int64
		MonthExpense     int64
		MonthNet         int64
		MoMNetPct        float64
		ExpenseTrendPct  float64
		YTDNet           int64
		SavingsGoalCents int64
		SavingsProgress  int
		CashflowSVG      template.HTML
		CategorySlices   []homeCategorySlice
		DonutGradient    string
		DonutTotalAbs    int64
		Insight          string
		Recent           []store.Transaction
	}{
		RunningTotal:     running,
		Timezone:         u.Timezone,
		MonthIncome:      monthIncome,
		MonthExpense:     monthExpense,
		MonthNet:         monthNet,
		MoMNetPct:        momNetPct,
		ExpenseTrendPct:  expenseTrendPct,
		YTDNet:           ytdNet,
		SavingsGoalCents: savingsGoalCents,
		SavingsProgress:  savingsProgress,
		CashflowSVG:      cashflowSVG,
		CategorySlices:   slices,
		DonutGradient:    donutGrad,
		DonutTotalAbs:    totalAbs,
		Insight:          insight,
		Recent:           recent,
	}
	ld := LayoutData{
		Title:     "Dashboard",
		User:      u,
		Year:      time.Now().UTC().Year(),
		Active:    "home",
		MainClass: "layer-stack layer-stack--wide",
	}
	a.renderShell(w, "home_inner.html", data, ld)
}

func abs64(x int64) int64 {
	if x < 0 {
		return -x
	}
	return x
}

func pctChangeNet(prev, prevPrev int64) float64 {
	if prevPrev == 0 {
		if prev == 0 {
			return 0
		}
		return 100
	}
	return float64(prev-prevPrev) / float64(abs64(prevPrev)) * 100
}

func pctChangeExpense(prev, prevPrev int64) float64 {
	a, b := abs64(prev), abs64(prevPrev)
	if b == 0 {
		if a == 0 {
			return 0
		}
		return 100
	}
	return float64(a-b) / float64(b) * 100
}

func donutGradient(pcts []float64) string {
	if len(pcts) == 0 {
		return ""
	}
	colors := []string{"#306369", "#4a7d82", "#678a92", "#8aa3a8", "#b5c4c8"}
	var b strings.Builder
	b.WriteString("conic-gradient(from -90deg, ")
	cum := 0.0
	for i, p := range pcts {
		if i > 0 {
			b.WriteString(", ")
		}
		next := cum + p
		if next > 100.01 {
			next = 100
		}
		fmt.Fprintf(&b, "%s %.3f%% %.3f%%", colors[i%len(colors)], cum, next)
		cum = next
	}
	b.WriteString(")")
	return b.String()
}
