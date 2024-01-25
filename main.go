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
			Name:             "Ransomware Attack",
			LowerProb:        0.1,
			LowerProbStdDev:  float64Pointer(0.1515),
			UpperProb:        0.625,
			UpperProbStdDev:  float64Pointer(0.1515),
			Confidence:       0.38825,
			ConfidenceStdDev: float64Pointer(0.194125),
			Timeframe:        montecargo.EveryFiveYears,
			MinImpact:        float64Pointer(275_000),
			MinImpactStdDev:  float64Pointer(137_500),
			MaxImpact:        float64Pointer(251_000_000),
			MaxImpactStdDev:  float64Pointer(100_400_000),
		},
		{
			Name:             "Data Breach",
			LowerProb:        0.15,
			LowerProbStdDev:  float64Pointer(0.1),
			UpperProb:        0.9,
			UpperProbStdDev:  float64Pointer(0.15),
			Confidence:       0.425,
			ConfidenceStdDev: float64Pointer(0.2125),
			Timeframe:        montecargo.EveryFiveYears,
			MinImpact:        float64Pointer(100_000),
			MinImpactStdDev:  float64Pointer(50_000),
			MaxImpact:        float64Pointer(300_000_000),
			MaxImpactStdDev:  float64Pointer(120_000_000),
		},
		{
			Name:             "System Compromise",
			LowerProb:        0.05,
			LowerProbStdDev:  float64Pointer(0.075),
			UpperProb:        0.55,
			UpperProbStdDev:  float64Pointer(0.165),
			Confidence:       0.3,
			ConfidenceStdDev: float64Pointer(0.15),
			Timeframe:        montecargo.EveryFiveYears,
			MinImpact:        float64Pointer(500_000),
			MinImpactStdDev:  float64Pointer(250_000),
			MaxImpact:        float64Pointer(200_000_000),
			MaxImpactStdDev:  float64Pointer(100_000_000),
		},

		// defenses
		{
			Name:                            "Host-Level Breach Detected",
			LowerProb:                       0.40,
			LowerProbStdDev:                 float64Pointer(0.08),
			UpperProb:                       0.70,
			UpperProbStdDev:                 float64Pointer(0.10),
			Confidence:                      0.55,
			ConfidenceStdDev:                float64Pointer(0.075),
			Timeframe:                       montecargo.EveryTwoYears,
			MinImpact:                       float64Pointer(10_000),
			MaxImpact:                       float64Pointer(50_000),
			MinImpactStdDev:                 float64Pointer(5_000),
			MaxImpactStdDev:                 float64Pointer(25_000),
			IsCostSaving:                    true,
			CostOfImplementationLower:       float64Pointer(20_000),
			CostOfImplementationLowerStdDev: float64Pointer(10_000),
			CostOfImplementationUpper:       float64Pointer(100_000),
			CostOfImplementationUpperStdDev: float64Pointer(50_000),
		},
		{
			Name:                            "Network-Level Breach Detected",
			LowerProb:                       0.45,
			LowerProbStdDev:                 float64Pointer(0.09),
			UpperProb:                       0.75,
			UpperProbStdDev:                 float64Pointer(0.12),
			Confidence:                      0.60,
			ConfidenceStdDev:                float64Pointer(0.08),
			Timeframe:                       montecargo.EveryTwoYears,
			MinImpact:                       float64Pointer(20_000),
			MaxImpact:                       float64Pointer(100_000),
			MinImpactStdDev:                 float64Pointer(10_000),
			MaxImpactStdDev:                 float64Pointer(50_000),
			IsCostSaving:                    true,
			CostOfImplementationLower:       float64Pointer(50_000),
			CostOfImplementationLowerStdDev: float64Pointer(25_000),
			CostOfImplementationUpper:       float64Pointer(200_000),
			CostOfImplementationUpperStdDev: float64Pointer(100_000),
		},
	}

	dependencies := map[string][]montecargo.Dependency{
		"Data Breach": {
			{EventName: "Network-Level Breach Detected", Condition: "not happens"},
		},
		"System Compromise": {
			{EventName: "Data Breach", Condition: "happens"},
			{EventName: "Host-Level Breach Detected", Condition: "not happens"},
		},
		"Ransomware Attack": {
			{EventName: "System Compromise", Condition: "happens"},
		},
		// ... other dependencies ...
	}

	// Perform Monte Carlo Simulation
	numSimulations := 1_000_000
	simulationResult := montecargo.MonteCarloSimulation(events, numSimulations, dependencies)

	// Use the updated event statistics
	eventStats := simulationResult.EventStats

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
		eventStat := eventStats[event.Name]
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
			fmt.Printf("  Expected Savings Range (considering cost of implementation) within %s: $%.2f to $%.2f\n", timeframeStr, -upperBound, -lowerBound)
			fmt.Printf("  Minimum Lower Bound of Cost of Implementation: $%.2f\n", eventStat.MinCostOfImplementation)
			fmt.Printf("  Maximum Upper Bound of Cost of Implementation: $%.2f\n", eventStat.MaxCostOfImplementation)
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
