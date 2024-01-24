package montecargo

// MonteCarloSimulation orchestrates the Monte Carlo simulation process.
func MonteCarloSimulation(events []Event, numSimulations int, dependencies map[string][]Dependency) SimulationResult {
	// Filter events into independent and dependent categories
	independentEvents := filterIndependentEvents(events, dependencies)
	dependentEvents := filterDependentEvents(events, dependencies)

	// Run simulation for independent events
	independentResults := simulate(independentEvents, numSimulations, map[string]EventStat{}, dependencies)

	// Run simulation for dependent events
	dependentResults := simulateDependent(dependentEvents, numSimulations, dependencies, map[string]EventStat{})

	// Combine results from independent and dependent simulations
	combinedResults := combineSimulationResults(independentResults, dependentResults)

	// Calculate statistics for all events based on the combined results
	// allEventStats := CalculateEventStats(combinedResults.EventResults, numSimulations)

	// Convert the map of EventResults to a SimulationResult
	finalResults := SimulationResult{EventResults: combinedResults.EventResults}
	return finalResults
}
