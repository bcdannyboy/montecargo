package montecargo

func collectResults(eventResults map[string]*EventResult, events []Event, resultsChan chan [][3]int) {
	for batch := range resultsChan {
		for _, result := range batch {
			eventIndex := result[0]
			outcome := result[1]
			impact := result[2] // Extract impact

			eventName := events[eventIndex].Name // Assuming you can map index to name

			// Initialize EventResult if not already done
			if _, exists := eventResults[eventName]; !exists {
				eventResults[eventName] = &EventResult{}
			}

			// Aggregate results
			eventResults[eventName].Sum += outcome
			eventResults[eventName].SumOfSquares += float64(outcome * outcome)

			// Accumulate impact values
			if outcome == 1 {
				eventResults[eventName].ImpactSum += float64(impact)
				eventResults[eventName].ImpactSumOfSquares += float64(impact * impact)
			}
		}
	}
}

func aggregateEventResults(a, b EventResult) EventResult {
	// Logic to aggregate two EventResult instances
	return EventResult{
		Sum:                a.Sum + b.Sum,
		SumOfSquares:       a.SumOfSquares + b.SumOfSquares,
		ImpactSum:          a.ImpactSum + b.ImpactSum,
		ImpactSumOfSquares: a.ImpactSumOfSquares + b.ImpactSumOfSquares,
	}
}

func combineSimulationResults(independentResults, dependentResults SimulationResult) SimulationResult {
	combinedResults := SimulationResult{EventResults: make(map[string]EventResult)}
	for eventName, eventResult := range independentResults.EventResults {
		combinedResults.EventResults[eventName] = eventResult
	}
	for eventName, eventResult := range dependentResults.EventResults {
		if existingResult, exists := combinedResults.EventResults[eventName]; exists {
			combinedResults.EventResults[eventName] = aggregateEventResults(existingResult, eventResult)
		} else {
			combinedResults.EventResults[eventName] = eventResult
		}
	}
	return combinedResults
}
