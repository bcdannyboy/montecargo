package montecargo

import (
	"math/rand"
	"runtime"
	"sync"
	"time"
)

func MonteCarloSimulation(events *[]Event, numSimulations int) {
	var wg sync.WaitGroup
	cpuCores := runtime.NumCPU()
	resultsChan := make(chan [][3]int, cpuCores) // Updated to handle impact

	for i := 0; i < cpuCores; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			localRand := rand.New(rand.NewSource(time.Now().UnixNano()))
			for j := 0; j < numSimulations/cpuCores; j++ {
				batchResults := make([][3]int, len(*events)) // Updated to handle impact
				for k, event := range *events {
					adjustedProb := adjustProbabilityForTimeframe(event)

					// Adjust probability based on confidence standard deviation if applicable
					if event.ConfidenceStdDev != nil {
						if event.LowerProbStdDev != nil && event.UpperProbStdDev != nil {
							// Adjust probability bounds based on confidence standard deviation
							adjustedProb = adjustProbabilityWithConfidenceStdDev(adjustedProb, *event.ConfidenceStdDev, localRand)
						} else {
							// Apply a random adjustment to probability
							adjustedProb += localRand.NormFloat64() * *event.ConfidenceStdDev
							if adjustedProb < 0 {
								adjustedProb = 0
							} else if adjustedProb > 1 {
								adjustedProb = 1
							}
						}
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
		}()
	}

	// Close the channel when all goroutines are done
	go func() {
		wg.Wait()
		close(resultsChan)
	}()

	// Collect results and update sums
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
