package montecargo

import (
	"math/rand"
	"runtime"
	"sync"
	"time"
)

// Simulate events without dependencies
func simulateEvents(events *[]Event, numSimulations int, wg *sync.WaitGroup, resultsChan chan [][3]int) {
	defer wg.Done()
	cpuCores := runtime.NumCPU()
	for i := 0; i < cpuCores; i++ {
		wg.Add(1)
		go func(localEvents []Event) {
			defer wg.Done()
			localRand := rand.New(rand.NewSource(time.Now().UnixNano()))
			for j := 0; j < numSimulations/cpuCores; j++ {
				batchResults := make([][3]int, len(localEvents))
				for k, event := range localEvents {
					adjustedProb := adjustProbabilityForTimeframe(event)

					// Adjust probability based on confidence standard deviation if applicable
					if event.ConfidenceStdDev != nil {
						adjustedProb = adjustProbabilityWithConfidenceStdDev(adjustedProb, *event.ConfidenceStdDev, localRand)
					}

					result := 0
					impact := 0
					if localRand.Float64() < adjustedProb {
						result = 1
						// Calculate financial impact for the event
						if event.MinImpact != nil && event.MaxImpact != nil && *event.MaxImpact > *event.MinImpact {
							impact = int(localRand.Float64()*(*event.MaxImpact-*event.MinImpact) + *event.MinImpact)
							// Adjust impact based on confidence standard deviation if applicable
							if event.ConfidenceStdDev != nil {
								impactAdjustment := int(localRand.NormFloat64() * *event.ConfidenceStdDev)
								impact += impactAdjustment
							}
						}
					}
					batchResults[k] = [3]int{k, result, impact} // Include impact in results
				}
				resultsChan <- batchResults
			}
		}(*events)
	}
}

// Simulate events with dependencies
func simulateDependentEvents(events *[]Event, numSimulations int, wg *sync.WaitGroup, resultsChan chan [][3]int, dependencies map[string][]Dependency, eventStats map[string](struct {
	Probability float64
	StdDev      float64
}), mutex *sync.Mutex) {
	defer wg.Done()
	cpuCores := runtime.NumCPU()
	for i := 0; i < cpuCores; i++ {
		wg.Add(1)
		go func(localEvents []Event) {
			defer wg.Done()
			localRand := rand.New(rand.NewSource(time.Now().UnixNano()))
			for j := 0; j < numSimulations/cpuCores; j++ {
				batchResults := make([][3]int, len(localEvents))
				for k, event := range localEvents {
					adjustedProb := adjustProbabilityForTimeframe(event)

					// Adjust probability based on dependency outcomes
					if dependentConditions, exists := dependencies[event.Name]; exists {
						mutex.Lock() // Lock the mutex before accessing the shared map
						for _, condition := range dependentConditions {
							dependencyStats := eventStats[condition.EventName]
							if condition.Condition == "not happens" {
								adjustedProb *= (1 - dependencyStats.Probability)
							} else {
								adjustedProb *= dependencyStats.Probability
							}
						}
						mutex.Unlock() // Unlock the mutex after accessing the shared map
					}

					result := 0
					impact := 0
					if localRand.Float64() < adjustedProb {
						result = 1
						// Calculate financial impact for the event
						if event.MinImpact != nil && event.MaxImpact != nil && *event.MaxImpact > *event.MinImpact {
							impact = int(localRand.Float64()*(*event.MaxImpact-*event.MinImpact) + *event.MinImpact)
							// Adjust impact based on confidence standard deviation if applicable
							if event.ConfidenceStdDev != nil {
								impactAdjustment := int(localRand.NormFloat64() * *event.ConfidenceStdDev)
								impact += impactAdjustment
							}
						}
					}
					batchResults[k] = [3]int{k, result, impact} // Include impact in results
				}
				resultsChan <- batchResults
			}
		}(*events)
	}
}

func simulate(events *[]Event, numSimulations int, wg *sync.WaitGroup, resultsChan chan [][3]int) {
	defer wg.Done()
	for i := 0; i < runtime.NumCPU(); i++ {
		wg.Add(1)
		go simulateEvents(events, numSimulations/runtime.NumCPU(), wg, resultsChan)
	}
}

func simulateDependent(events *[]Event, numSimulations int, wg *sync.WaitGroup, resultsChan chan [][3]int, dependencies map[string][]Dependency, eventStats map[string](struct {
	Probability float64
	StdDev      float64
}), mutex *sync.Mutex) {
	defer wg.Done()
	for i := 0; i < runtime.NumCPU(); i++ {
		wg.Add(1)
		go simulateDependentEvents(events, numSimulations/runtime.NumCPU(), wg, resultsChan, dependencies, eventStats, mutex)
	}
}
