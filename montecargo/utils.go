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

func filterIndependentEvents(events []Event, dependencies map[string][]Dependency) []Event {
	var independentEvents []Event
	for _, event := range events {
		if _, exists := dependencies[event.Name]; !exists {
			independentEvents = append(independentEvents, event)
		}
	}
	return independentEvents
}

func filterDependentEvents(events []Event, dependencies map[string][]Dependency) []Event {
	var dependentEvents []Event
	for _, event := range events {
		if _, exists := dependencies[event.Name]; exists {
			dependentEvents = append(dependentEvents, event)
		}
	}
	return dependentEvents
}

func collectResults(events *[]Event, resultsChan chan [][3]int) {
	for batch := range resultsChan {
		for _, result := range batch {
			eventIndex := result[0]
			outcome := result[1]
			impact := result[2] // Extract impact
			(*events)[eventIndex].Results = append((*events)[eventIndex].Results, outcome)
			(*events)[eventIndex].Sum += outcome
			(*events)[eventIndex].SumOfSquares += float64(outcome * outcome)

			// Accumulate impact values
			if outcome == 1 {
				(*events)[eventIndex].ImpactSum += float64(impact)
				(*events)[eventIndex].ImpactSumOfSquares += float64(impact * impact)
			}
		}
	}
}

func TimeframeToString(tf Timeframe) string {
	switch tf {
	case Daily:
		return "1 day"
	case Weekly:
		return "1 week"
	case Monthly:
		return "1 month"
	case Yearly:
		return "1 year"
	case EveryTwoYears:
		return "2 years"
	case EveryFiveYears:
		return "5 years"
	case EveryTenYears:
		return "10 years"
	default:
		return "unknown timeframe"
	}
}
