package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/bcdannyboy/montecargo/montecargo"
)

func float64Pointer(value float64) *float64 {
	return &value
}

func main() {
	rand.Seed(time.Now().UnixNano())

	events := []montecargo.Event{
		{
			Name:       "Phishing Attack",
			LowerProb:  0.60,
			UpperProb:  0.80,
			Confidence: 0.85,
			Timeframe:  montecargo.Yearly,
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
			{EventName: "Phishing Attack", Condition: "not happens"},
		},
		// ... other dependencies ...
	}

	// Perform Monte Carlo Simulation
	numSimulations := 1_000_000
	montecargo.MonteCarloSimulation(&events, numSimulations, dependencies)

	// Output the results and calculate standard deviation
	for _, event := range events {
		probability, probStdDev, impactMean, impactStdDev := montecargo.MeanSTD(event, numSimulations)

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
		if event.MinImpact != nil || event.MaxImpact != nil {
			impactLowerBound := impactMean - impactStdDev
			impactUpperBound := impactMean + impactStdDev

			fmt.Printf("  Mean Financial Impact: $%.2f\n", impactMean)
			fmt.Printf("  Financial Impact Standard Deviation: $%.2f\n", impactStdDev)
			fmt.Printf("  Expected Financial Impact Range within %s: $%.2f to $%.2f\n", timeframeStr, impactLowerBound, impactUpperBound)
		}
	}
}
