package montecargo

// MonteCarloSimulation orchestrates the Monte Carlo simulation process.
func MonteCarloSimulation(events []Event, numSimulations int, dependencies map[string][]Dependency) SimulationResult {
	// Filter events into independent and dependent categories
	independentEvents := filterIndependentEvents(events, dependencies)
	dependentEvents := filterDependentEvents(events, dependencies)

	// Run simulation for independent events
	independentResults, independentEventStats := simulate(independentEvents, numSimulations, dependencies, make(map[string]EventStat))

	// Run simulation for dependent events
	// Use the independentEventStats as initial stats for dependent events
	dependentResults, dependentEventStats := simulateDependent(dependentEvents, numSimulations, dependencies, independentEventStats)

	// Combine results from independent and dependent simulations
	combinedResults := combineSimulationResults(independentResults, dependentResults)

	// Use the most updated event stats (from dependentEventStats)
	finalEventStats := combineEventStats(independentEventStats, dependentEventStats)

	// Convert the map of EventResults to a SimulationResult
	finalResults := SimulationResult{EventResults: combinedResults.EventResults, EventStats: finalEventStats}
	return finalResults
}
