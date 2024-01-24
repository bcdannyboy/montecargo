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

	// Define cybersecurity events with various timeframes
	events := []montecargo.Event{
		{
			Name:       "Phishing Attack",
			LowerProb:  0.70,
			UpperProb:  0.90,
			Confidence: 0.90,
			Timeframe:  montecargo.Yearly,
		},
		{
			Name:       "Ransomware Attack",
			LowerProb:  0.15,
			UpperProb:  0.45,
			Confidence: 0.85,
			Timeframe:  montecargo.Yearly,
			MinImpact:  float64Pointer(500000),
			MaxImpact:  float64Pointer(2000000),
		},
		{
			Name:       "Data Breach",
			LowerProb:  0.05,
			UpperProb:  0.25,
			Confidence: 0.80,
			Timeframe:  montecargo.EveryTwoYears,
			MinImpact:  float64Pointer(1000000),
			MaxImpact:  float64Pointer(5000000),
		},
		{
			Name:       "Breach Detection",
			LowerProb:  0.50,
			UpperProb:  0.70,
			Confidence: 0.95,
			Timeframe:  montecargo.Monthly,
		},
		{
			Name:       "System Compromise",
			LowerProb:  0.10,
			UpperProb:  0.30,
			Confidence: 0.75,
			Timeframe:  montecargo.EveryFiveYears,
			MinImpact:  float64Pointer(1500000),
			MaxImpact:  float64Pointer(3000000),
		},
		{
			Name:             "Insider Threat",
			LowerProb:        0.05,
			UpperProb:        0.20,
			LowerProbStdDev:  float64Pointer(0.02),
			UpperProbStdDev:  float64Pointer(0.03),
			Confidence:       0.70,
			ConfidenceStdDev: float64Pointer(0.05),
			Timeframe:        montecargo.Yearly,
			MinImpact:        float64Pointer(300000),
			MaxImpact:        float64Pointer(1500000),
			MinImpactStdDev:  float64Pointer(50000),
			MaxImpactStdDev:  float64Pointer(200000),
		},
		{
			Name:             "IT System Failure",
			LowerProb:        0.10,
			UpperProb:        0.40,
			LowerProbStdDev:  float64Pointer(0.05),
			UpperProbStdDev:  float64Pointer(0.10),
			Confidence:       0.60,
			ConfidenceStdDev: float64Pointer(0.07),
			Timeframe:        montecargo.EveryTwoYears,
			MinImpact:        float64Pointer(200000),
			MaxImpact:        float64Pointer(1000000),
			MinImpactStdDev:  float64Pointer(100000),
			MaxImpactStdDev:  float64Pointer(300000),
		},
		// ... other events ...
	}

	// Define dependencies
	dependencies := map[string]string{
		"Ransomware Attack": "Data Breach", // Ransomware Attack depends on Data Breach
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

		// Print the results
		fmt.Printf("Event: %s\n", event.Name)
		fmt.Printf("  Probability: %.2f%%\n", probability*100)
		fmt.Printf("  Standard Deviation: %.2f%%\n", probStdDev*100)
		fmt.Printf("  Chance within timeframe: %.2f%% to %.2f%%\n", probLowerBound*100, probUpperBound*100)

		// Print impact information if applicable
		if event.MinImpact != nil || event.MaxImpact != nil {
			impactLowerBound := impactMean - impactStdDev
			impactUpperBound := impactMean + impactStdDev

			fmt.Printf("  Mean Financial Impact: $%.2f\n", impactMean)
			fmt.Printf("  Financial Impact Standard Deviation: $%.2f\n", impactStdDev)
			fmt.Printf("  Expected Financial Impact Range: $%.2f to $%.2f\n", impactLowerBound, impactUpperBound)
		}
	}
}
