package montecargo

import (
	"fmt"
	"math"
	"math/rand"
)

func MeanSTD(eventResult EventResult, numSimulations int) (probability, probStdDev, impactMean, impactStdDev float64) {
	probability = float64(eventResult.Sum) / float64(numSimulations)
	probMean := probability
	probVariance := (eventResult.SumOfSquares / float64(numSimulations)) - (probMean * probMean)
	probStdDev = math.Sqrt(probVariance)

	// Calculating mean and standard deviation for impact
	// Note: Impact calculations are only relevant when the event occurs (eventResult.Sum > 0)
	if eventResult.Sum > 0 {
		impactMean = eventResult.ImpactSum / float64(eventResult.Sum)
		impactVariance := (eventResult.ImpactSumOfSquares / float64(eventResult.Sum)) - (impactMean * impactMean)
		impactStdDev = math.Sqrt(impactVariance)
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

func adjustDependentProbability(probability, dependencyProb, dependencyStdDev float64, localRand *rand.Rand) float64 {
	// scale the probability by the dependency's probability
	// and apply a random adjustment based on the dependency's standard deviation
	adjustedProbability := probability * dependencyProb
	adjustmentFactor := localRand.NormFloat64() * dependencyStdDev
	adjustedProbability += adjustmentFactor

	// Ensure the adjusted probability is within the range [0, 1]
	if adjustedProbability < 0 {
		adjustedProbability = 0
	} else if adjustedProbability > 1 {
		adjustedProbability = 1
	}

	return adjustedProbability
}

func CalculateEventStats(simulationResults map[string]EventResult, numSimulations int, events []Event) map[string]EventStat {
	eventStats := make(map[string]EventStat)

	for _, event := range events {
		eventResult, exists := simulationResults[event.Name]
		if !exists {
			fmt.Printf("Event result not found for %s\n", event.Name)
			continue
		}

		probability, stdDev, _, _ := MeanSTD(eventResult, numSimulations)
		stat := EventStat{
			Probability: probability,
			StdDev:      stdDev,
		}

		// Calculate the bounds for the cost of implementation if applicable
		if event.CostOfImplementationLower != nil && event.CostOfImplementationUpper != nil {
			lowerCost := *event.CostOfImplementationLower
			upperCost := *event.CostOfImplementationUpper

			if event.CostOfImplementationLowerStdDev != nil {
				lowerCost -= *event.CostOfImplementationLowerStdDev // Minimum
				upperCost += *event.CostOfImplementationUpperStdDev // Maximum
			}

			stat.MinCostOfImplementation = math.Max(0, lowerCost) // Ensure it's not negative
			stat.MaxCostOfImplementation = upperCost
		}

		eventStats[event.Name] = stat
	}

	return eventStats
}

func NormalCDF(x, mean, stdDev float64) float64 {
	return 0.5 * (1 + math.Erf((x-mean)/(stdDev*math.Sqrt2)))
}

func CalculateExpectedLossRange(events []Event, eventStats map[string]EventStat) (float64, float64, float64, float64, float64, map[string]struct {
	MinLoss              float64
	MaxLoss              float64
	AvgLoss              float64
	ProbabilityExceedMin float64
	ProbabilityExceedMax float64
	ProbabilityExceedAvg float64
}) {
	totalMinLoss := 0.0
	totalMaxLoss := 0.0
	totalAvgLoss := 0.0
	var totalVariance float64
	lossBreakdown := make(map[string]struct {
		MinLoss              float64
		MaxLoss              float64
		AvgLoss              float64
		ProbabilityExceedMin float64
		ProbabilityExceedMax float64
		ProbabilityExceedAvg float64
	})

	for _, event := range events {
		stat, exists := eventStats[event.Name]
		if exists && event.MinImpact != nil && event.MaxImpact != nil {
			eventMinLoss := *event.MinImpact * stat.Probability
			eventMaxLoss := *event.MaxImpact * stat.Probability
			eventAvgLoss := (eventMinLoss + eventMaxLoss) / 2
			totalMinLoss += eventMinLoss
			totalMaxLoss += eventMaxLoss
			totalAvgLoss += eventAvgLoss
			totalVariance += stat.StdDev * stat.StdDev

			// Calculate probabilities of exceeding min, max, and avg loss
			probExceedMin := 1 - NormalCDF(eventMinLoss, eventAvgLoss, stat.StdDev)
			probExceedMax := 1 - NormalCDF(eventMaxLoss, eventAvgLoss, stat.StdDev)
			probExceedAvg := 1 - NormalCDF(eventAvgLoss, eventAvgLoss, stat.StdDev)

			lossBreakdown[event.Name] = struct {
				MinLoss              float64
				MaxLoss              float64
				AvgLoss              float64
				ProbabilityExceedMin float64
				ProbabilityExceedMax float64
				ProbabilityExceedAvg float64
			}{
				MinLoss:              eventMinLoss,
				MaxLoss:              eventMaxLoss,
				AvgLoss:              eventAvgLoss,
				ProbabilityExceedMin: probExceedMin,
				ProbabilityExceedMax: probExceedMax,
				ProbabilityExceedAvg: probExceedAvg,
			}
		}
	}

	totalStdDev := math.Sqrt(totalVariance)
	probExceedTotalMin := 1 - NormalCDF(totalMinLoss, totalAvgLoss, totalStdDev)
	probExceedTotalMax := 1 - NormalCDF(totalMaxLoss, totalAvgLoss, totalStdDev)

	return totalMinLoss, totalMaxLoss, totalAvgLoss, probExceedTotalMin, probExceedTotalMax, lossBreakdown
}

func calculateImpact(event Event, probability float64, localRand *rand.Rand) int {
	if event.MinImpact == nil || event.MaxImpact == nil {
		return 0
	}

	impactRange := *event.MaxImpact - *event.MinImpact
	impact := *event.MinImpact + impactRange*localRand.Float64()

	if event.ConfidenceStdDev != nil {
		impactAdjustment := localRand.NormFloat64() * *event.ConfidenceStdDev
		impact += impactAdjustment
	}

	if event.IsCostSaving {
		return -int(impact * probability) // Negative impact for cost savings
	}
	return int(impact * probability) // Positive impact for losses
}

func calculateTotalMitigatedImpact(costSavingEvent Event, events []Event, eventStats map[string]EventStat, dependencies map[string][]Dependency, localRand *rand.Rand) (float64, float64, float64) {
	totalSavings := 0.0
	var lowerCost, upperCost float64

	costSavingEventStat, exists := eventStats[costSavingEvent.Name]
	if !exists {
		fmt.Printf("Cost-saving event stat not found for %s\n", costSavingEvent.Name)
		return 0.0, 0.0, 0.0
	}

	occurrences := GetOccurrencesPerYear(costSavingEvent.Timeframe)

	for _, dep := range dependencies[costSavingEvent.Name] {
		if dep.Condition == "not happens" {
			dependentEventStat, exists := eventStats[dep.EventName]
			if exists {
				dependentEvent, found := findEventByName(events, dep.EventName)
				if found && dependentEvent.MinImpact != nil && dependentEvent.MaxImpact != nil {
					avgImpact := (*dependentEvent.MinImpact + *dependentEvent.MaxImpact) / 2
					savings := avgImpact * (1 - dependentEventStat.Probability) * costSavingEventStat.Probability
					totalSavings += savings * occurrences
				}
			}
		}
	}

	if costSavingEvent.CostOfImplementationLower != nil {
		lowerCost = *costSavingEvent.CostOfImplementationLower
		if costSavingEvent.CostOfImplementationLowerStdDev != nil {
			lowerCostAdjustment := localRand.NormFloat64() * *costSavingEvent.CostOfImplementationLowerStdDev
			lowerCost += lowerCostAdjustment
		}
	}

	if costSavingEvent.CostOfImplementationUpper != nil {
		upperCost = *costSavingEvent.CostOfImplementationUpper
		if costSavingEvent.CostOfImplementationUpperStdDev != nil {
			upperCostAdjustment := localRand.NormFloat64() * *costSavingEvent.CostOfImplementationUpperStdDev
			upperCost += upperCostAdjustment
		}
	}

	return totalSavings, lowerCost, upperCost
}
