package dashboard

import "net/url"

// Query and value names for the dashboard stats window (?period=30d|12m).
// Used by [BuildPageData] and [PeriodQueryRelative]; keep in sync with templates via tmpl funcs.
const (
	PeriodQueryParam = "period"
	StatsPeriod30d   = "30d"
	StatsPeriod12m   = "12m"
)

// PeriodQueryRelative returns a root-relative URL that selects the stats window (?period=…).
func PeriodQueryRelative(period string) string {
	return "/?" + PeriodQueryParam + "=" + url.QueryEscape(period)
}
