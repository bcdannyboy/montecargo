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
	resultsChan := make(chan [][2]int, cpuCores)

	for i := 0; i < cpuCores; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			localRand := rand.New(rand.NewSource(time.Now().UnixNano()))
			for j := 0; j < numSimulations/cpuCores; j++ {
				batchResults := make([][2]int, len(*events))
				for k, event := range *events {
					adjustedProb := adjustProbabilityForTimeframe(event)
					result := 0
					if localRand.Float64() < adjustedProb {
						result = 1
					}
					batchResults[k] = [2]int{k, result}
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
			(*events)[eventIndex].Results = append((*events)[eventIndex].Results, outcome)
			(*events)[eventIndex].Sum += outcome
			(*events)[eventIndex].SumOfSquares += float64(outcome * outcome)
		}
	}
}
