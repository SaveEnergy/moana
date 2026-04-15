package dashboard

// statsPeriodConfig normalizes ?period= for the trailing stats window and prior-period labels.
type statsPeriodConfig struct {
	Period        string // "30d" or "12m"
	InclusiveDays int
	PriorPhrase   string
}

func parseStatsPeriod(periodQuery string) statsPeriodConfig {
	switch periodQuery {
	case "12m":
		return statsPeriodConfig{
			Period:        "12m",
			InclusiveDays: 365,
			PriorPhrase:   "prior 12 months",
		}
	case "30d", "":
		return statsPeriodConfig{
			Period:        "30d",
			InclusiveDays: 30,
			PriorPhrase:   "prior 30 days",
		}
	default:
		return statsPeriodConfig{
			Period:        "30d",
			InclusiveDays: 30,
			PriorPhrase:   "prior 30 days",
		}
	}
}
