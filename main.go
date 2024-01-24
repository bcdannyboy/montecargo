package main

import (
	"fmt"
	"math/rand"
	"runtime"
	"time"

	"github.com/bcdannyboy/montecargo/montecargo"
)

func float64Pointer(value float64) *float64 {
	return &value
}

func main() {
	startTime := time.Now()
	rand.Seed(time.Now().UnixNano())

	events := []montecargo.Event{
		{
			Name:             "Host-Level Breach Detected",
			LowerProb:        0.80,
			LowerProbStdDev:  float64Pointer(0.10),
			UpperProb:        0.95,
			UpperProbStdDev:  float64Pointer(0.05),
			Confidence:       0.6,
			ConfidenceStdDev: float64Pointer(0.10),
			MinImpact:        float64Pointer(5000),
			MaxImpact:        float64Pointer(10000),
			MinImpactStdDev:  float64Pointer(10000),
			MaxImpactStdDev:  float64Pointer(10000),
			Timeframe:        montecargo.Monthly,
			IsCostSaving:     true,
		},
		{
			Name:       "Ransomware Attack",
			LowerProb:  0.25,
			UpperProb:  0.50,
			Confidence: 0.70,
			Timeframe:  montecargo.Yearly,
			MinImpact:  float64Pointer(400000),
			MaxImpact:  float64Pointer(1500000),
		},
		{
			Name:       "Data Breach",
			LowerProb:  0.10,
			UpperProb:  0.25,
			Confidence: 0.80,
			Timeframe:  montecargo.Yearly,
			MinImpact:  float64Pointer(800000),
			MaxImpact:  float64Pointer(2500000),
		},
		{
			Name:       "System Compromise",
			LowerProb:  0.15,
			UpperProb:  0.35,
			Confidence: 0.65,
			Timeframe:  montecargo.EveryFiveYears,
			MinImpact:  float64Pointer(1000000),
			MaxImpact:  float64Pointer(2000000),
		},
		{
			Name:       "Insider Threat",
			LowerProb:  0.05,
			UpperProb:  0.15,
			Confidence: 0.75,
			Timeframe:  montecargo.Yearly,
			MinImpact:  float64Pointer(200000),
			MaxImpact:  float64Pointer(1000000),
		},
		{
			Name:       "IT System Failure",
			LowerProb:  0.05,
			UpperProb:  0.20,
			Confidence: 0.60,
			Timeframe:  montecargo.EveryTwoYears,
			MinImpact:  float64Pointer(100000),
			MaxImpact:  float64Pointer(500000),
		},
		// ... other events ...
	}

	dependencies := map[string][]montecargo.Dependency{
		"Ransomware Attack": {
			{EventName: "Data Breach", Condition: "happens"},
			{EventName: "Host-Level Breach Detected", Condition: "not happens"},
		},
		"Insider Threat": {
			{EventName: "Host-Level Breach Detected", Condition: "not happens"},
		},
		// ... other dependencies ...
	}

	// Perform Monte Carlo Simulation
	numSimulations := 10_000_000
	simulationResult := montecargo.MonteCarloSimulation(events, numSimulations, dependencies)

	// Calculate event statistics
	eventStats := montecargo.CalculateEventStats(simulationResult.EventResults, numSimulations)

	// Calculate expected loss range
	totalMinLoss, totalMaxLoss, totalAvgLoss, probExceedTotalMin, probExceedTotalMax, lossBreakdown := montecargo.CalculateExpectedLossRange(events, eventStats)

	fmt.Printf("Total Expected Minimum Loss: $%.2f\n", totalMinLoss)
	fmt.Printf("Total Expected Maximum Loss: $%.2f\n", totalMaxLoss)
	fmt.Printf("Total Expected Average Loss: $%.2f\n", totalAvgLoss)
	fmt.Printf("Probability of Exceeding Total Min Loss: %.2f%%\n", probExceedTotalMin*100)
	fmt.Printf("Probability of Exceeding Total Max Loss: %.2f%%\n", probExceedTotalMax*100)

	fmt.Println()
	// Output the results and calculate standard deviation
	for _, event := range events {
		eventResult := simulationResult.EventResults[event.Name]
		probability, probStdDev, impactMean, impactStdDev := montecargo.MeanSTD(eventResult, numSimulations)

		// Calculate the upper and lower bounds for probability
		probLowerBound := probability - probStdDev
		if probLowerBound < 0 {
			probLowerBound = 0 // Ensure the lower bound is not negative
		}
		probUpperBound := probability + probStdDev
		if probUpperBound > 1 {
			probUpperBound = 1 // Ensure the upper bound does not exceed 100%
		}

		// Convert Timeframe to a human-readable string
		timeframeStr := montecargo.TimeframeToString(event.Timeframe)

		// Print the results
		fmt.Printf("Event: %s\n", event.Name)
		fmt.Printf("  Probability: %.2f%%\n", probability*100)
		fmt.Printf("  Standard Deviation: %.2f%%\n", probStdDev*100)
		fmt.Printf("  Chance within %s: %.2f%% to %.2f%%\n", timeframeStr, probLowerBound*100, probUpperBound*100)

		// Print impact information if applicable
		if (event.MinImpact != nil || event.MaxImpact != nil) && !event.IsCostSaving {
			impactLowerBound := impactMean - impactStdDev
			impactUpperBound := impactMean + impactStdDev
			// Output for cost-incurring events
			fmt.Printf("  Mean Financial Impact: $%.2f\n", impactMean)
			fmt.Printf("  Financial Impact Standard Deviation: $%.2f\n", impactStdDev)
			fmt.Printf("  Expected Financial Impact Range within %s: $%.2f to $%.2f\n", timeframeStr, impactLowerBound, impactUpperBound)

		}

		if event.IsCostSaving {
			lowerBound := impactMean - impactStdDev
			upperBound := impactMean + impactStdDev

			// Output for cost-saving events
			fmt.Printf("  Estimated Savings: $%.2f\n", -impactMean)
			fmt.Printf("  Savings Standard Deviation: $%.2f\n", impactStdDev)
			fmt.Printf("  Expected Savings Range within %s: $%.2f to $%.2f\n", timeframeStr, -upperBound, -lowerBound)

		} else {

			lossInfo := lossBreakdown[event.Name]
			fmt.Printf("  Min Loss: $%.2f\n", lossInfo.MinLoss)
			fmt.Printf("  Max Loss: $%.2f\n", lossInfo.MaxLoss)
			fmt.Printf("  Probability of Exceeding Min Loss: %.2f%%\n", lossInfo.ProbabilityExceedMin*100)
			fmt.Printf("  Probability of Exceeding Max Loss: %.2f%%\n", lossInfo.ProbabilityExceedMax*100)
			fmt.Printf("  Probability of Exceeding Avg Loss: %.2f%%\n", lossInfo.ProbabilityExceedAvg*100)
			fmt.Println()
		}
	}

	depCount := 0
	for _, event := range events {
		if deps, exists := dependencies[event.Name]; exists {
			depCount += len(deps)
		}
	}

	numCPU := runtime.NumCPU()
	fmt.Printf("performed %d simulations with %d items and %d dependencies in %s on %d CPU cores\n", numSimulations, len(events), depCount, time.Since(startTime), numCPU)
}
