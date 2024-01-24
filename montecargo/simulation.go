package montecargo

import (
	"math/rand"
	"runtime"
	"sync"
	"time"
)

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

					// Aggregate results for each event
					eventResult := localResult.EventResults[event.Name]
					eventResult.Sum += result
					eventResult.SumOfSquares += float64(result * result)
					eventResult.ImpactSum += float64(impact)
					eventResult.ImpactSumOfSquares += float64(impact * impact)

					// Handle cost-saving items
					if event.IsCostSaving {
						mitigatedImpact := calculateTotalMitigatedImpact(event, events, eventStats, dependencies)
						eventResult.ImpactSum += mitigatedImpact // Add as positive value
					}

					localResult.EventResults[event.Name] = eventResult
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

				// Adjust probability based on dependencies
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

					// Aggregate results for each event
					eventResult := localResult.EventResults[event.Name]
					eventResult.Sum += result
					eventResult.SumOfSquares += float64(result * result)
					eventResult.ImpactSum += float64(impact)
					eventResult.ImpactSumOfSquares += float64(impact * impact)

					// Handle cost-saving items
					if event.IsCostSaving {
						mitigatedImpact := calculateTotalMitigatedImpact(event, events, eventStats, dependencies)
						eventResult.ImpactSum += mitigatedImpact // Add as positive value
					}

					localResult.EventResults[event.Name] = eventResult
				}
			}
		}
	}

	return localResult
}

func simulate(events []Event, numSimulations int, eventStats map[string]EventStat, dependencies map[string][]Dependency) SimulationResult {
	finalResult := SimulationResult{EventResults: make(map[string]EventResult)}
	var wg sync.WaitGroup
	resultsChan := make(chan SimulationResult, runtime.NumCPU())

	for i := 0; i < runtime.NumCPU(); i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			// Pass the additional parameters to simulateEvents
			result := simulateEvents(events, numSimulations/runtime.NumCPU(), eventStats, dependencies)
			resultsChan <- result
		}()
	}

	wg.Wait()
	close(resultsChan)

	// Aggregate results from all goroutines
	for result := range resultsChan {
		for eventName, eventResult := range result.EventResults {
			// Aggregate results for each event
			finalResult.EventResults[eventName] = aggregateEventResults(finalResult.EventResults[eventName], eventResult)
		}
	}
	return finalResult
}

func simulateDependent(events []Event, numSimulations int, dependencies map[string][]Dependency, eventStats map[string]EventStat) SimulationResult {
	finalResult := SimulationResult{EventResults: make(map[string]EventResult)}
	var wg sync.WaitGroup
	resultsChan := make(chan SimulationResult, runtime.NumCPU())

	for i := 0; i < runtime.NumCPU(); i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			result := simulateDependentEvents(events, numSimulations/runtime.NumCPU(), dependencies, eventStats)
			resultsChan <- result
		}()
	}

	wg.Wait()
	close(resultsChan)

	// Aggregate results from all goroutines
	for result := range resultsChan {
		for eventName, eventResult := range result.EventResults {
			// Aggregate results for each event
			finalResult.EventResults[eventName] = aggregateEventResults(finalResult.EventResults[eventName], eventResult)
		}
	}

	return finalResult
}
