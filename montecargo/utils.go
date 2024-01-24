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
