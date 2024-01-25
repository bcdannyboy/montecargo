package montecargo

import (
	"math/rand"
	"runtime"
	"sync"
	"time"
)

var mutex sync.Mutex

// Simulate events without dependencies

func simulateEvents(events []Event, numSimulations int, eventStats map[string]EventStat, dependencies map[string][]Dependency) SimulationResult {
	localResult := SimulationResult{EventResults: make(map[string]EventResult)}
	cpuCores := runtime.NumCPU()

	for i := 0; i < cpuCores; i++ {
		localRand := rand.New(rand.NewSource(time.Now().UnixNano()))
		for j := 0; j < numSimulations/cpuCores; j++ {
			for _, event := range events {
				adjustedProb := adjustProbabilityForTimeframe(event)
				if event.ConfidenceStdDev != nil {
					adjustedProb = adjustProbabilityWithConfidenceStdDev(adjustedProb, *event.ConfidenceStdDev, localRand)
				}

				result := 0
				if localRand.Float64() < adjustedProb {
					result = 1
					impact := calculateImpact(event, adjustedProb, localRand)

					mutex.Lock()
					eventResult := localResult.EventResults[event.Name]
					eventResult.Sum += result
					eventResult.SumOfSquares += float64(result * result)
					eventResult.ImpactSum += float64(impact)
					eventResult.ImpactSumOfSquares += float64(impact * impact)
					localResult.EventResults[event.Name] = eventResult
					mutex.Unlock()
				}
			}
		}
	}

	return localResult
}

// Simulate events with dependencies
func simulateDependentEvents(events []Event, numSimulations int, dependencies map[string][]Dependency, eventStats map[string]EventStat) SimulationResult {
	localResult := SimulationResult{EventResults: make(map[string]EventResult)}
	cpuCores := runtime.NumCPU()

	for i := 0; i < cpuCores; i++ {
		localRand := rand.New(rand.NewSource(time.Now().UnixNano()))
		for j := 0; j < numSimulations/cpuCores; j++ {
			for _, event := range events {
				adjustedProb := adjustProbabilityForTimeframe(event)

				if dependentConditions, exists := dependencies[event.Name]; exists {
					for _, condition := range dependentConditions {
						dependencyStats := eventStats[condition.EventName]
						if condition.Condition == "not happens" {
							adjustedProb *= (1 - dependencyStats.Probability)
						} else {
							adjustedProb *= dependencyStats.Probability
						}
					}
				}

				result := 0
				if localRand.Float64() < adjustedProb {
					result = 1
					impact := calculateImpact(event, adjustedProb, localRand)

					mutex.Lock()
					eventResult := localResult.EventResults[event.Name]
					eventResult.Sum += result
					eventResult.SumOfSquares += float64(result * result)
					eventResult.ImpactSum += float64(impact)
					eventResult.ImpactSumOfSquares += float64(impact * impact)
					localResult.EventResults[event.Name] = eventResult
					mutex.Unlock()
				}
			}
		}
	}

	return localResult
}

func simulate(events []Event, numSimulations int, dependencies map[string][]Dependency, initialEventStats map[string]EventStat) (SimulationResult, map[string]EventStat) {
	finalResult := SimulationResult{EventResults: make(map[string]EventResult)}
	var wg sync.WaitGroup
	resultsChan := make(chan SimulationResult, runtime.NumCPU())

	// Use the provided initialEventStats if available; otherwise, create a new map.
	eventStats := initialEventStats
	if len(eventStats) == 0 {
		eventStats = make(map[string]EventStat)
	}

	for i := 0; i < runtime.NumCPU(); i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			result := simulateEvents(events, numSimulations/runtime.NumCPU(), eventStats, dependencies)
			resultsChan <- result
		}()
	}

	wg.Wait()
	close(resultsChan)

	// Aggregate results from all goroutines
	for result := range resultsChan {
		for eventName, eventResult := range result.EventResults {
			finalResult.EventResults[eventName] = aggregateEventResults(finalResult.EventResults[eventName], eventResult)
		}
	}

	// Calculate event stats after simulation
	eventStats = CalculateEventStats(finalResult.EventResults, numSimulations, events)

	return finalResult, eventStats
}

func simulateDependent(events []Event, numSimulations int, dependencies map[string][]Dependency, updatedEventStats map[string]EventStat) (SimulationResult, map[string]EventStat) {
	finalResult := SimulationResult{EventResults: make(map[string]EventResult)}
	var wg sync.WaitGroup
	resultsChan := make(chan SimulationResult, runtime.NumCPU())

	for i := 0; i < runtime.NumCPU(); i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			result := simulateDependentEvents(events, numSimulations/runtime.NumCPU(), dependencies, updatedEventStats)
			resultsChan <- result
		}()
	}

	wg.Wait()
	close(resultsChan)

	// Aggregate results from all goroutines
	for result := range resultsChan {
		for eventName, eventResult := range result.EventResults {
			finalResult.EventResults[eventName] = aggregateEventResults(finalResult.EventResults[eventName], eventResult)
		}
	}

	// Update event stats after dependent simulation
	eventStats := CalculateEventStats(finalResult.EventResults, numSimulations, events)
	return finalResult, eventStats
}
