package testing

import (
	"fmt"
	"math"
	"testing"

	"github.com/bcdannyboy/montecargo/montecargo"
	testing_utils "github.com/bcdannyboy/montecargo/testing/testing_utils"
	"github.com/stretchr/testify/assert"
)

var events = []montecargo.Event{
	{
		Name:                            "Test Event A",
		LowerProb:                       0.2,
		UpperProb:                       0.4,
		LowerProbStdDev:                 testing_utils.Float64Pointer(0.05),
		UpperProbStdDev:                 testing_utils.Float64Pointer(0.05),
		Confidence:                      0.8,
		ConfidenceStdDev:                testing_utils.Float64Pointer(0.1),
		Timeframe:                       montecargo.Yearly,
		MinImpact:                       testing_utils.Float64Pointer(10000),
		MaxImpact:                       testing_utils.Float64Pointer(50000),
		MinImpactStdDev:                 testing_utils.Float64Pointer(5000),
		MaxImpactStdDev:                 testing_utils.Float64Pointer(10000),
		IsCostSaving:                    false,
		CostOfImplementationLower:       testing_utils.Float64Pointer(15000),
		CostOfImplementationUpper:       testing_utils.Float64Pointer(40000),
		CostOfImplementationLowerStdDev: testing_utils.Float64Pointer(5000),
		CostOfImplementationUpperStdDev: testing_utils.Float64Pointer(10000),
	},
	{
		Name:                            "Test Event B",
		LowerProb:                       0.2,
		UpperProb:                       0.4,
		LowerProbStdDev:                 testing_utils.Float64Pointer(0.05),
		UpperProbStdDev:                 testing_utils.Float64Pointer(0.05),
		Confidence:                      0.8,
		ConfidenceStdDev:                testing_utils.Float64Pointer(0.1),
		Timeframe:                       montecargo.Yearly,
		MinImpact:                       testing_utils.Float64Pointer(10000),
		MaxImpact:                       testing_utils.Float64Pointer(50000),
		MinImpactStdDev:                 testing_utils.Float64Pointer(5000),
		MaxImpactStdDev:                 testing_utils.Float64Pointer(10000),
		IsCostSaving:                    false,
		CostOfImplementationLower:       testing_utils.Float64Pointer(15000),
		CostOfImplementationUpper:       testing_utils.Float64Pointer(40000),
		CostOfImplementationLowerStdDev: testing_utils.Float64Pointer(5000),
		CostOfImplementationUpperStdDev: testing_utils.Float64Pointer(10000),
	},
	{
		Name:                            "Test Event C",
		LowerProb:                       0.3,
		UpperProb:                       0.6,
		LowerProbStdDev:                 testing_utils.Float64Pointer(0.1),
		UpperProbStdDev:                 testing_utils.Float64Pointer(0.1),
		Confidence:                      0.9,
		ConfidenceStdDev:                testing_utils.Float64Pointer(0.05),
		Timeframe:                       montecargo.Monthly,
		MinImpact:                       testing_utils.Float64Pointer(20000),
		MaxImpact:                       testing_utils.Float64Pointer(100000),
		MinImpactStdDev:                 testing_utils.Float64Pointer(10000),
		MaxImpactStdDev:                 testing_utils.Float64Pointer(20000),
		IsCostSaving:                    true,
		CostOfImplementationLower:       testing_utils.Float64Pointer(15000),
		CostOfImplementationUpper:       testing_utils.Float64Pointer(40000),
		CostOfImplementationLowerStdDev: testing_utils.Float64Pointer(5000),
		CostOfImplementationUpperStdDev: testing_utils.Float64Pointer(10000),
	},
	{
		Name:                            "Test Event D",
		LowerProb:                       0.3,
		UpperProb:                       0.6,
		LowerProbStdDev:                 testing_utils.Float64Pointer(0.1),
		UpperProbStdDev:                 testing_utils.Float64Pointer(0.1),
		Confidence:                      0.9,
		ConfidenceStdDev:                testing_utils.Float64Pointer(0.05),
		Timeframe:                       montecargo.Monthly,
		MinImpact:                       testing_utils.Float64Pointer(20000),
		MaxImpact:                       testing_utils.Float64Pointer(100000),
		MinImpactStdDev:                 testing_utils.Float64Pointer(10000),
		MaxImpactStdDev:                 testing_utils.Float64Pointer(20000),
		IsCostSaving:                    true,
		CostOfImplementationLower:       testing_utils.Float64Pointer(15000),
		CostOfImplementationUpper:       testing_utils.Float64Pointer(40000),
		CostOfImplementationLowerStdDev: testing_utils.Float64Pointer(5000),
		CostOfImplementationUpperStdDev: testing_utils.Float64Pointer(10000),
	},
	{
		Name:                            "Test Event E",
		LowerProb:                       0.45,
		UpperProb:                       0.55,
		LowerProbStdDev:                 testing_utils.Float64Pointer(0.02),
		UpperProbStdDev:                 testing_utils.Float64Pointer(0.02),
		Confidence:                      0.95,
		ConfidenceStdDev:                testing_utils.Float64Pointer(0.01),
		Timeframe:                       montecargo.Yearly,
		MinImpact:                       testing_utils.Float64Pointer(20000),
		MaxImpact:                       testing_utils.Float64Pointer(100000),
		MinImpactStdDev:                 testing_utils.Float64Pointer(10000),
		MaxImpactStdDev:                 testing_utils.Float64Pointer(20000),
		IsCostSaving:                    true,
		CostOfImplementationLower:       testing_utils.Float64Pointer(15000),
		CostOfImplementationUpper:       testing_utils.Float64Pointer(40000),
		CostOfImplementationLowerStdDev: testing_utils.Float64Pointer(5000),
		CostOfImplementationUpperStdDev: testing_utils.Float64Pointer(10000),
	},
	{
		Name:                            "Test Event F",
		LowerProb:                       0.45,
		UpperProb:                       0.55,
		LowerProbStdDev:                 testing_utils.Float64Pointer(0.02),
		UpperProbStdDev:                 testing_utils.Float64Pointer(0.02),
		Confidence:                      0.95,
		ConfidenceStdDev:                testing_utils.Float64Pointer(0.01),
		Timeframe:                       montecargo.Yearly,
		MinImpact:                       testing_utils.Float64Pointer(20000),
		MaxImpact:                       testing_utils.Float64Pointer(100000),
		MinImpactStdDev:                 testing_utils.Float64Pointer(10000),
		MaxImpactStdDev:                 testing_utils.Float64Pointer(20000),
		IsCostSaving:                    true,
		CostOfImplementationLower:       testing_utils.Float64Pointer(15000),
		CostOfImplementationUpper:       testing_utils.Float64Pointer(40000),
		CostOfImplementationLowerStdDev: testing_utils.Float64Pointer(5000),
		CostOfImplementationUpperStdDev: testing_utils.Float64Pointer(10000),
	},
	{
		Name:                            "Test Event G",
		LowerProb:                       0.1,
		UpperProb:                       0.9,
		LowerProbStdDev:                 testing_utils.Float64Pointer(0.15),
		UpperProbStdDev:                 testing_utils.Float64Pointer(0.15),
		Confidence:                      0.7,
		ConfidenceStdDev:                testing_utils.Float64Pointer(0.2),
		Timeframe:                       montecargo.Monthly,
		MinImpact:                       testing_utils.Float64Pointer(20000),
		MaxImpact:                       testing_utils.Float64Pointer(100000),
		MinImpactStdDev:                 testing_utils.Float64Pointer(10000),
		MaxImpactStdDev:                 testing_utils.Float64Pointer(20000),
		IsCostSaving:                    true,
		CostOfImplementationLower:       testing_utils.Float64Pointer(15000),
		CostOfImplementationUpper:       testing_utils.Float64Pointer(40000),
		CostOfImplementationLowerStdDev: testing_utils.Float64Pointer(5000),
		CostOfImplementationUpperStdDev: testing_utils.Float64Pointer(10000),
	},
	{
		Name:                            "Test Event H",
		LowerProb:                       0.1,
		UpperProb:                       0.9,
		LowerProbStdDev:                 testing_utils.Float64Pointer(0.15),
		UpperProbStdDev:                 testing_utils.Float64Pointer(0.15),
		Confidence:                      0.7,
		ConfidenceStdDev:                testing_utils.Float64Pointer(0.2),
		Timeframe:                       montecargo.Monthly,
		MinImpact:                       testing_utils.Float64Pointer(20000),
		MaxImpact:                       testing_utils.Float64Pointer(100000),
		MinImpactStdDev:                 testing_utils.Float64Pointer(10000),
		MaxImpactStdDev:                 testing_utils.Float64Pointer(20000),
		IsCostSaving:                    true,
		CostOfImplementationLower:       testing_utils.Float64Pointer(15000),
		CostOfImplementationUpper:       testing_utils.Float64Pointer(40000),
		CostOfImplementationLowerStdDev: testing_utils.Float64Pointer(5000),
		CostOfImplementationUpperStdDev: testing_utils.Float64Pointer(10000),
	},
	{
		Name:                            "Test Event I",
		LowerProb:                       0.01,
		UpperProb:                       0.05,
		LowerProbStdDev:                 testing_utils.Float64Pointer(0.005),
		UpperProbStdDev:                 testing_utils.Float64Pointer(0.005),
		Confidence:                      0.9,
		ConfidenceStdDev:                testing_utils.Float64Pointer(0.05),
		Timeframe:                       montecargo.EveryFiveYears,
		MinImpact:                       testing_utils.Float64Pointer(20000),
		MaxImpact:                       testing_utils.Float64Pointer(100000),
		MinImpactStdDev:                 testing_utils.Float64Pointer(10000),
		MaxImpactStdDev:                 testing_utils.Float64Pointer(20000),
		IsCostSaving:                    true,
		CostOfImplementationLower:       testing_utils.Float64Pointer(15000),
		CostOfImplementationUpper:       testing_utils.Float64Pointer(40000),
		CostOfImplementationLowerStdDev: testing_utils.Float64Pointer(5000),
		CostOfImplementationUpperStdDev: testing_utils.Float64Pointer(10000),
	},
	{
		Name:                            "Test Event J",
		LowerProb:                       0.01,
		UpperProb:                       0.05,
		LowerProbStdDev:                 testing_utils.Float64Pointer(0.005),
		UpperProbStdDev:                 testing_utils.Float64Pointer(0.005),
		Confidence:                      0.9,
		ConfidenceStdDev:                testing_utils.Float64Pointer(0.05),
		Timeframe:                       montecargo.EveryFiveYears,
		MinImpact:                       testing_utils.Float64Pointer(20000),
		MaxImpact:                       testing_utils.Float64Pointer(100000),
		MinImpactStdDev:                 testing_utils.Float64Pointer(10000),
		MaxImpactStdDev:                 testing_utils.Float64Pointer(20000),
		IsCostSaving:                    true,
		CostOfImplementationLower:       testing_utils.Float64Pointer(15000),
		CostOfImplementationUpper:       testing_utils.Float64Pointer(40000),
		CostOfImplementationLowerStdDev: testing_utils.Float64Pointer(5000),
		CostOfImplementationUpperStdDev: testing_utils.Float64Pointer(10000),
	},
	{
		Name:                            "Test Event K",
		LowerProb:                       0.9,
		UpperProb:                       0.99,
		LowerProbStdDev:                 testing_utils.Float64Pointer(0.02),
		UpperProbStdDev:                 testing_utils.Float64Pointer(0.02),
		Confidence:                      0.99,
		ConfidenceStdDev:                testing_utils.Float64Pointer(0.01),
		Timeframe:                       montecargo.Daily,
		MinImpact:                       testing_utils.Float64Pointer(20000),
		MaxImpact:                       testing_utils.Float64Pointer(100000),
		MinImpactStdDev:                 testing_utils.Float64Pointer(10000),
		MaxImpactStdDev:                 testing_utils.Float64Pointer(20000),
		IsCostSaving:                    true,
		CostOfImplementationLower:       testing_utils.Float64Pointer(15000),
		CostOfImplementationUpper:       testing_utils.Float64Pointer(40000),
		CostOfImplementationLowerStdDev: testing_utils.Float64Pointer(5000),
		CostOfImplementationUpperStdDev: testing_utils.Float64Pointer(10000),
	},
	{
		Name:                            "Test Event L",
		LowerProb:                       0.9,
		UpperProb:                       0.99,
		LowerProbStdDev:                 testing_utils.Float64Pointer(0.02),
		UpperProbStdDev:                 testing_utils.Float64Pointer(0.02),
		Confidence:                      0.99,
		ConfidenceStdDev:                testing_utils.Float64Pointer(0.01),
		Timeframe:                       montecargo.Daily,
		MinImpact:                       testing_utils.Float64Pointer(20000),
		MaxImpact:                       testing_utils.Float64Pointer(100000),
		MinImpactStdDev:                 testing_utils.Float64Pointer(10000),
		MaxImpactStdDev:                 testing_utils.Float64Pointer(20000),
		IsCostSaving:                    true,
		CostOfImplementationLower:       testing_utils.Float64Pointer(15000),
		CostOfImplementationUpper:       testing_utils.Float64Pointer(40000),
		CostOfImplementationLowerStdDev: testing_utils.Float64Pointer(5000),
		CostOfImplementationUpperStdDev: testing_utils.Float64Pointer(10000),
	},
}

func TestProbabilityCalculations(t *testing.T) {

	// randomize events array
	shuffled_events := testing_utils.ShuffleArray(events)

	dependencies := map[string][]montecargo.Dependency{}
	numSimulations := 1_000_000
	fmt.Printf("Running %d simulations on %d events...\n", numSimulations, len(shuffled_events))
	simulationResult := montecargo.MonteCarloSimulation(shuffled_events, numSimulations, dependencies)

	for _, event := range shuffled_events {

		eventStat := simulationResult.EventStats[event.Name]

		expectedAvg := testing_utils.CalculateExpectedAverage(event, dependencies, simulationResult.EventStats)
		actualAvg := eventStat.Probability
		stdDev := eventStat.StdDev
		errorRate := math.Abs(actualAvg - expectedAvg)

		confidenceScore := testing_utils.CalculateAdjustedConfidenceScore(errorRate, stdDev, event)

		confidenceThreshold := testing_utils.DetermineConfidenceThreshold(event)

		// Calculate allowable range
		allowableLower := confidenceThreshold - 0.5*stdDev
		allowableUpper := confidenceThreshold + 0.5*stdDev

		fmt.Printf("Analyzing '%s':\n", event.Name)
		fmt.Printf("  - Expected Risk Occurrence: %.2f%%\n", expectedAvg*100)
		fmt.Printf("  - Observed Risk Occurrence: %.2f%%\n", actualAvg*100)
		fmt.Printf("  - Variability (Standard Deviation): %.2f%%\n", stdDev*100)
		fmt.Printf("  - Accuracy (Error Rate): %.2f%%\n", errorRate*100)
		fmt.Printf("  - Confidence in Model's Prediction: %.2f%%\n", confidenceScore*100)
		fmt.Printf("  - Confidence Threshold: %.2f%%\n", confidenceThreshold*100)
		fmt.Printf("  - Allowable Range: [%.2f%%, %.2f%%] based on 0.5*stdDev: %.2f%%, stdDev: %.2f%%\n", allowableLower*100, allowableUpper*100, 0.5*stdDev*100, stdDev*100)

		withinRange := confidenceScore > allowableLower &&
			confidenceScore < allowableUpper

		assert.True(t, withinRange,
			"Confidence score must be within range for event %s", event.Name)
	}
}

func TestProbabilityConvergence(t *testing.T) {

	dependencies := map[string][]montecargo.Dependency{}

	numSimulations := []int{10_000, 100_000, 1_000_000}

	for _, sims := range numSimulations {

		results := montecargo.MonteCarloSimulation(events, sims, dependencies)

		for _, event := range events {

			actualProb := results.EventStats[event.Name].Probability

			expectedProb := testing_utils.CalculateExpectedAverage(event, dependencies, results.EventStats)

			divergence := math.Abs(actualProb - expectedProb)

			threshold := 0.2 // Set threshold for acceptable divergence

			if event.Timeframe == montecargo.Daily {
				threshold = 0.1 // Tighter threshold for high frequency events
			}

			if divergence > threshold {
				t.Errorf("Probability for %s did not converge within "+
					"threshold after %d simulations", event.Name, sims)
			} else {
				fmt.Printf("Probability for %s converged within threshold after %d simulations\n", event.Name, sims)
			}
		}
	}
}
