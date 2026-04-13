package handlers

import (
	"fmt"
	"html"
	"html/template"
	"strings"

	"moana/internal/money"
	"moana/internal/store"
)

const cashflowMergeLimit = 6

var expenseRibbonColors = []string{
	"#b31b25", "#c45c26", "#2a6f97", "#6a4c93", "#457b9d", "#88b4c4", "#6c757d", "#8d99ae",
}

func mergeCategoryTopN(rows []store.CategoryAmount, limit int) []store.CategoryAmount {
	if len(rows) <= limit {
		return rows
	}
	out := make([]store.CategoryAmount, limit)
	copy(out, rows[:limit-1])
	var rest int64
	for _, r := range rows[limit-1:] {
		rest += r.AmountCents
	}
	out[limit-1] = store.CategoryAmount{Name: "Other", AmountCents: rest}
	return out
}

func truncLabel(s string, max int) string {
	r := []rune(s)
	if len(r) <= max {
		return s
	}
	if max < 2 {
		return "…"
	}
	return string(r[:max-1]) + "…"
}

func ribbonPath(x1, y0, y1, x2 float64) string {
	if y1-y0 < 0.5 {
		return ""
	}
	dx := 44.0
	return fmt.Sprintf("M %.2f %.2f C %.2f %.2f %.2f %.2f %.2f %.2f L %.2f %.2f C %.2f %.2f %.2f %.2f %.2f %.2f Z",
		x1, y0, x1+dx, y0, x2-dx, y0, x2, y0,
		x2, y1, x2-dx, y1, x1+dx, y1, x1, y1)
}

// buildCashflowSVG renders a Sankey-style cashflow diagram for the current month.
func buildCashflowSVG(income, expense []store.CategoryAmount, totalIncome, totalExpense int64) template.HTML {
	if totalIncome <= 0 {
		return ""
	}

	const (
		vbW      = 440.0
		vbH      = 200.0
		yTop     = 26.0
		flowH    = 132.0
		xLBar    = 38.0
		barW     = 26.0
		xLRight  = xLBar + barW
		xCenterL = 176.0
		xCenterR = 216.0
		xRLeft   = 310.0
		xRBar    = 338.0
		barWR    = 26.0
		xRRight  = xRBar + barWR
	)
	yBot := yTop + flowH
	te := totalExpense
	if te < 0 {
		te = -te
	}
	surplus := totalIncome - te
	overspend := surplus < 0
	if overspend {
		surplus = 0
	}

	type seg struct {
		name  string
		cents int64
		color string
	}
	var rightSegs []seg
	for i, e := range expense {
		c := expenseRibbonColors[i%len(expenseRibbonColors)]
		rightSegs = append(rightSegs, seg{name: e.Name, cents: e.AmountCents, color: c})
	}
	if surplus > 0 {
		rightSegs = append(rightSegs, seg{name: "Surplus", cents: surplus, color: "#1b5e20"})
	}

	var b strings.Builder
	fmt.Fprintf(&b, `<svg class="dashboard-cashflow-svg" viewBox="0 0 %.0f %.0f" preserveAspectRatio="xMidYMid meet" role="img" aria-label="Cashflow from income to expenses">`, vbW, vbH)

	// Left → center ribbons (income greens)
	var yOff float64
	for i, inc := range income {
		h := float64(inc.AmountCents) / float64(totalIncome) * flowH
		y0 := yTop + yOff
		y1 := y0 + h
		col := "#52b788"
		if i%2 == 1 {
			col = "#40916c"
		}
		p := ribbonPath(xLRight, y0, y1, xCenterL)
		if p != "" {
			fmt.Fprintf(&b, `<path fill="%s" fill-opacity="0.42" d="%s"/>`, col, p)
		}
		yOff += h
	}

	// Center → right ribbons
	scaleDenom := float64(totalIncome)
	if overspend {
		scaleDenom = float64(te)
		if scaleDenom < 1 {
			scaleDenom = 1
		}
	}
	yOff = 0
	for _, s := range rightSegs {
		h := float64(s.cents) / scaleDenom * flowH
		y0 := yTop + yOff
		y1 := y0 + h
		p := ribbonPath(xCenterR, y0, y1, xRLeft)
		if p != "" {
			fmt.Fprintf(&b, `<path fill="%s" fill-opacity="0.5" d="%s"/>`, s.color, p)
		}
		yOff += h
	}

	// Bars
	yOff = 0
	for _, inc := range income {
		h := float64(inc.AmountCents) / float64(totalIncome) * flowH
		y0 := yTop + yOff
		fmt.Fprintf(&b, `<rect x="%.1f" y="%.1f" width="%.1f" height="%.1f" rx="3" fill="#2d6a4f"/>`, xLBar, y0, barW, h)
		yOff += h
	}
	fmt.Fprintf(&b, `<rect x="%.1f" y="%.1f" width="%.1f" height="%.1f" rx="4" fill="#306369"/>`, xCenterL, yTop, xCenterR-xCenterL, flowH)

	yOff = 0
	for _, s := range rightSegs {
		h := float64(s.cents) / scaleDenom * flowH
		y0 := yTop + yOff
		fmt.Fprintf(&b, `<rect x="%.1f" y="%.1f" width="%.1f" height="%.1f" rx="3" fill="%s"/>`, xRBar, y0, barWR, h, s.color)
		yOff += h
	}

	// Labels left
	yOff = 0
	for _, inc := range income {
		h := float64(inc.AmountCents) / float64(totalIncome) * flowH
		y0 := yTop + yOff
		ym := y0 + h/2
		l := truncLabel(inc.Name, 14)
		fmt.Fprintf(&b, `<text x="32" y="%.1f" text-anchor="end" class="dashboard-cashflow-lbl">%s</text>`, ym-5, html.EscapeString(l))
		fmt.Fprintf(&b, `<text x="32" y="%.1f" text-anchor="end" class="dashboard-cashflow-sub">%s</text>`, ym+7, html.EscapeString(money.FormatEUR(inc.AmountCents)))
		yOff += h
	}

	// Center label
	cx := (xCenterL + xCenterR) / 2
	fmt.Fprintf(&b, `<text x="%.1f" y="%.1f" text-anchor="middle" class="dashboard-cashflow-center-lbl">%s</text>`, cx, yTop+flowH/2-10, html.EscapeString("Cash flow"))
	fmt.Fprintf(&b, `<text x="%.1f" y="%.1f" text-anchor="middle" class="dashboard-cashflow-center-amt">%s</text>`, cx, yTop+flowH/2+9, html.EscapeString(money.FormatEUR(totalIncome)))

	// Labels right
	yOff = 0
	for _, s := range rightSegs {
		h := float64(s.cents) / scaleDenom * flowH
		y0 := yTop + yOff
		ym := y0 + h/2
		l := truncLabel(s.name, 14)
		fmt.Fprintf(&b, `<text x="%.1f" y="%.1f" text-anchor="start" class="dashboard-cashflow-lbl">%s</text>`, xRRight+8, ym-5, html.EscapeString(l))
		fmt.Fprintf(&b, `<text x="%.1f" y="%.1f" text-anchor="start" class="dashboard-cashflow-sub">%s</text>`, xRRight+8, ym+7, html.EscapeString(money.FormatEUR(s.cents)))
		yOff += h
	}

	if overspend {
		fmt.Fprintf(&b, `<text x="%.1f" y="%.1f" text-anchor="middle" class="dashboard-cashflow-warn">Overspent vs income this month</text>`, vbW/2, yBot+18)
	}

	b.WriteString(`</svg>`)
	return template.HTML(b.String())
}
