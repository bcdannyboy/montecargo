package testing_utils

import (
	"math"
	"math/rand"
	"time"

	"github.com/bcdannyboy/montecargo/montecargo"
	"gonum.org/v1/gonum/stat"
)

// CalculateStats calculates the average and standard deviation of the outcomes
func CalculateStats(outcomes []bool) (avg float64, stdDev float64) {
	var sum int
	for _, outcome := range outcomes {
		if outcome {
			sum++
		}
	}

	avg = float64(sum) / float64(len(outcomes))

	// Calculate standard deviation
	var sumOfSquares float64
	for _, outcome := range outcomes {
		if outcome {
			sumOfSquares += math.Pow(1-avg, 2)
		} else {
			sumOfSquares += math.Pow(0-avg, 2)
		}
	}
	stdDev = math.Sqrt(sumOfSquares / float64(len(outcomes)))

	return avg, stdDev
}

// CalculateAdjustedConfidenceScore calculates the confidence score based on the error rate,
// standard deviation, and using statistical validation methods.
func CalculateAdjustedConfidenceScore(errorRate, stdDev float64, event montecargo.Event) float64 {

	// Base confidence from error rate
	baseConfidence := 1.0 - errorRate

	// Statistical validation
	errorDistribution := makeNormalDistribution(errorRate, stdDev)
	cdf := cumulativeDistributionFunc(errorDistribution)
	validatedConfidence := 1.0 - cdf(0.0)

	// Monte Carlo simulation of possible confidence scores
	simulations := 1_000_000
	confScores := make([]float64, simulations)
	weights := make([]float64, simulations)
	for i := 0; i < simulations; i++ {
		adjustment := rand.NormFloat64() * stdDev
		score := baseConfidence + adjustment
		confScores[i] = math.Max(0.0, math.Min(1.0, score))
		weights[i] = 1.0
	}

	simulatedConfidence := stat.Mean(confScores, weights)

	// Take weighted average of scores
	finalConfidence := 0.7*validatedConfidence + 0.3*simulatedConfidence

	return finalConfidence
}

// DetermineConfidenceThreshold determines the threshold for the confidence score
// based on the event's characteristics and using statistical methods.
func DetermineConfidenceThreshold(event montecargo.Event) float64 {

	// Base threshold starts at 0.7
	baseThreshold := 0.7

	// Calculate the variability ratio
	variabilityRatio := 1.0
	if event.UpperProbStdDev != nil && event.LowerProbStdDev != nil {
		avgStdDev := (*event.UpperProbStdDev + *event.LowerProbStdDev) / 2
		variabilityRatio = avgStdDev / (event.UpperProb - event.LowerProb)
	}

	// Reduce base threshold for events with higher variability
	// Less aggressive reduction to avoid overly tight thresholds
	baseThreshold -= variabilityRatio * 0.1

	// Adjust threshold based on event timeframe
	switch event.Timeframe {
	case montecargo.Daily, montecargo.Weekly:
		baseThreshold -= 0.01
	case montecargo.Monthly:
		baseThreshold -= 0.005
	case montecargo.Yearly:
		// No adjustment
	case montecargo.EveryTwoYears, montecargo.EveryFiveYears:
		baseThreshold += 0.01
	case montecargo.EveryTenYears:
		baseThreshold += 0.02
	}

	// Actuarial adjustment
	occurrences := montecargo.GetOccurrencesPerYear(event.Timeframe)
	actuarialAdjustment := 1 - math.Pow(1-baseThreshold, 1/occurrences)

	baseThreshold = math.Min(baseThreshold+actuarialAdjustment*0.2, 0.9)

	// Monte Carlo simulation to refine threshold
	simulations := 1_000_000
	thresholdSim := make([]float64, simulations)
	weights := make([]float64, simulations)
	for i := 0; i < simulations; i++ {
		adjustment := rand.NormFloat64() * 0.02
		thresholdSim[i] = baseThreshold + adjustment
		weights[i] = 1.0
	}

	adjustedThreshold := stat.Mean(thresholdSim, weights)

	return math.Max(0.5, math.Min(0.9, adjustedThreshold))
}

// AdjustProbabilityForTimeframe adjusts the probability of an event based on its timeframe.
func AdjustProbabilityForTimeframe(probability float64, timeframe montecargo.Timeframe) float64 {

	switch timeframe {
	case montecargo.Daily:
		return probability / 365
	case montecargo.Weekly:
		return probability / 52
	case montecargo.Monthly:
		return probability / 12
	case montecargo.Yearly:
		return probability
	case montecargo.EveryTwoYears:
		return probability * 2
	case montecargo.EveryFiveYears:
		return probability * 5
	case montecargo.EveryTenYears:
		return probability * 10
	default:
		return probability
	}

}

// AdjustProbabilityWithConfidenceStdDev adjusts the probability based on a normal distribution
// centered around the original probability and a standard deviation defined by the confidence standard deviation.
func AdjustProbabilityWithConfidenceStdDev(probability, confidenceStdDev float64, localRand *rand.Rand) float64 {
	// Adjust the probability based on a normal distribution
	// The adjustment reflects the uncertainty in the probability estimate
	adjustedProbability := probability + localRand.NormFloat64()*confidenceStdDev

	// Ensure the adjusted probability is within the range [0, 1]
	// This step is crucial to maintain the validity of the probability value
	if adjustedProbability < 0 {
		adjustedProbability = 0
	} else if adjustedProbability > 1 {
		adjustedProbability = 1
	}

	return adjustedProbability
}

// CalculateExpectedAverage computes the expected average probability of an event,
// considering adjustments made in the MonteCarloSimulation.
func CalculateExpectedAverage(event montecargo.Event, dependencies map[string][]montecargo.Dependency, eventStats map[string]montecargo.EventStat) float64 {
	localRand := rand.New(rand.NewSource(time.Now().UnixNano()))

	// Start with the basic average probability
	probRange := event.UpperProb - event.LowerProb
	avgProb := event.LowerProb + probRange/2

	// Adjust for the event's timeframe
	avgProb = AdjustProbabilityForTimeframe(avgProb, event.Timeframe)

	// Adjust for confidence standard deviation
	if event.ConfidenceStdDev != nil {
		avgProb = AdjustProbabilityWithConfidenceStdDev(avgProb, *event.ConfidenceStdDev, localRand)
	}

	// Adjust for dependencies if any
	if dependentConditions, exists := dependencies[event.Name]; exists {
		for _, condition := range dependentConditions {
			dependencyStats := eventStats[condition.EventName]
			if condition.Condition == "not happens" {
				avgProb *= (1 - dependencyStats.Probability)
			} else {
				avgProb *= dependencyStats.Probability
			}
		}
	}

	// Ensure the probability is within the range [0, 1]
	if avgProb < 0 {
		avgProb = 0
	} else if avgProb > 1 {
		avgProb = 1
	}

	return avgProb
}

// makeNormalDistribution creates a normal distribution with given mean and std dev
func makeNormalDistribution(mean, stdDev float64) func(x float64) float64 {
	variance := stdDev * stdDev
	normFactor := 1.0 / (stdDev * math.Sqrt(2*math.Pi))
	return func(x float64) float64 {
		exponent := -((x - mean) * (x - mean)) / (2 * variance)
		return normFactor * math.Exp(exponent)
	}
}

// cumulativeDistributionFunc computes the CDF of a given distribution
func cumulativeDistributionFunc(dist func(x float64) float64) func(x float64) float64 {
	return func(x float64) float64 {
		sum := 0.0
		var i float64
		numSteps := 100_000
		stepSize := (x - float64(-3)) / float64(numSteps)

		for i = -3; i < x; i += stepSize {
			sum += dist(i) * stepSize
		}
		return sum
	}
}
