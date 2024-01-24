package montecargo

import "math"

func MeanSTD(event Event, numSimulations int) (probability float64, stdDev float64) {
	probability = float64(event.Sum) / float64(numSimulations)
	mean := probability
	variance := (event.SumOfSquares / float64(numSimulations)) - (mean * mean)
	stdDev = math.Sqrt(variance)
	return
}

func adjustProbabilityForTimeframe(event Event) float64 {
	// Adjust the probability based on the event's timeframe
	probRange := event.UpperProb - event.LowerProb
	avgProb := event.LowerProb + probRange/2

	switch event.Timeframe {
	case Daily:
		return avgProb / 365
	case Weekly:
		return avgProb / 52
	case Monthly:
		return avgProb / 12
	case Yearly:
		return avgProb
	case EveryTwoYears:
		return avgProb * 2
	case EveryFiveYears:
		return avgProb * 5
	case EveryTenYears:
		return avgProb * 10
	default:
		return avgProb
	}
}
