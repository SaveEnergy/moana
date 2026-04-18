package dashboard

// statsPeriodConfig normalizes ?period= for the trailing stats window and prior-period labels.
type statsPeriodConfig struct {
	Period        string // "30d" or "12m"
	InclusiveDays int
	PriorPhrase   string
}

func parseStatsPeriod(periodQuery string) statsPeriodConfig {
	switch periodQuery {
	case StatsPeriod12m:
		return statsPeriodConfig{
			Period:        StatsPeriod12m,
			InclusiveDays: 365,
			PriorPhrase:   "prior 12 months",
		}
	case StatsPeriod30d, "":
		return statsPeriodConfig{
			Period:        StatsPeriod30d,
			InclusiveDays: 30,
			PriorPhrase:   "prior 30 days",
		}
	default:
		return statsPeriodConfig{
			Period:        StatsPeriod30d,
			InclusiveDays: 30,
			PriorPhrase:   "prior 30 days",
		}
	}
}
