package montecargo

type Timeframe int

const (
	Daily Timeframe = iota
	Weekly
	Monthly
	Yearly
	EveryTwoYears
	EveryFiveYears
	EveryTenYears
)

type Event struct {
	Name            string
	LowerProb       float64
	UpperProb       float64
	LowerProbStdDev *float64 // Optional standard deviation for LowerProb
	UpperProbStdDev *float64 // Optional standard deviation for UpperProb
	Confidence      float64
	Timeframe       Timeframe
	Results         []int
	Sum             int
	SumOfSquares    float64
}
