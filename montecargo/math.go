package montecargo

import (
	"math"
	"math/rand"
)

func MeanSTD(event Event, numSimulations int) (probability, probStdDev, impactMean, impactStdDev float64) {
	probability = float64(event.Sum) / float64(numSimulations)
	probMean := probability
	probVariance := (event.SumOfSquares / float64(numSimulations)) - (probMean * probMean)
	probStdDev = math.Sqrt(probVariance)

	// Check if impacts are used (i.e., MinImpact and MaxImpact are not nil)
	if event.MinImpact != nil && event.MaxImpact != nil {
		// Calculating mean and standard deviation for impact
		// Note: Impact calculations are only relevant when the event occurs (event.Sum > 0)
		if event.Sum > 0 {
			impactMean = event.ImpactSum / float64(event.Sum)
			impactVariance := (event.ImpactSumOfSquares / float64(event.Sum)) - (impactMean * impactMean)
			impactStdDev = math.Sqrt(impactVariance)
		}
	}

	return
}

func adjustProbabilityWithConfidenceStdDev(probability, confidenceStdDev float64, localRand *rand.Rand) float64 {
	// Adjust the probability based on a normal distribution centered around the original probability
	// and a standard deviation defined by the confidence standard deviation.
	adjustedProbability := probability + localRand.NormFloat64()*confidenceStdDev

	// Ensure the adjusted probability is within the range [0, 1]
	if adjustedProbability < 0 {
		adjustedProbability = 0
	} else if adjustedProbability > 1 {
		adjustedProbability = 1
	}

	return adjustedProbability
}
