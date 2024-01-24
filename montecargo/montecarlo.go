package montecargo

import (
	"runtime"
	"sync"
)

func MonteCarloSimulation(events *[]Event, numSimulations int, dependencies map[string][]Dependency) {
	var wg sync.WaitGroup
	numCPU := runtime.NumCPU()
	resultsChan := make(chan [][3]int, numSimulations*numCPU)

	// Simulate independent events
	independentEvents := filterIndependentEvents(*events, dependencies)
	wg.Add(1)
	go simulate(&independentEvents, numSimulations, &wg, resultsChan)

	// Calculate probabilities and standard deviations for independent events
	var mutex sync.Mutex
	eventStats := calculateEventStats(independentEvents, numSimulations)

	// Simulate dependent events
	dependentEvents := filterDependentEvents(*events, dependencies)
	wg.Add(1)
	go simulateDependent(&dependentEvents, numSimulations, &wg, resultsChan, dependencies, eventStats, &mutex)

	// Wait for all simulations to finish
	wg.Wait()

	// Close the channel when all goroutines are done
	close(resultsChan)

	// Collect results and update sums
	collectResults(events, resultsChan)
}
