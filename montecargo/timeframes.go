package montecargo

func ParseTimeframe(input string) Timeframe {
	switch input {
	case "daily":
		return Daily
	case "weekly":
		return Weekly
	case "monthly":
		return Monthly
	case "yearly":
		return Yearly
	case "2 years":
		return EveryTwoYears
	case "5 years":
		return EveryFiveYears
	case "10 years":
		return EveryTenYears
	default:
		return Yearly // Default or throw an error
	}
}

func adjustProbabilityForTimeframe(event Event) float64 {
	// Adjust the probability based on the event's timeframe
	probRange := event.UpperProb - event.LowerProb
	avgProb := event.LowerProb + probRange/2

	switch event.Timeframe {
	case Daily:
		return avgProb / 365
	case Weekly:
		return avgProb / 52
	case Monthly:
		return avgProb / 12
	case Yearly:
		return avgProb
	case EveryTwoYears:
		return avgProb * 2
	case EveryFiveYears:
		return avgProb * 5
	case EveryTenYears:
		return avgProb * 10
	default:
		return avgProb
	}
}

func GetOccurrencesPerYear(timeframe Timeframe) float64 {
	switch timeframe {
	case Daily:
		return 365
	case Weekly:
		return 52
	case Monthly:
		return 12
	case Yearly:
		return 1
	case EveryTwoYears:
		return 0.5
	case EveryFiveYears:
		return 0.2
	case EveryTenYears:
		return 0.1
	default:
		return 1 // Default to yearly if unknown
	}
}
