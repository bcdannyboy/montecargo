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
	Name               string
	LowerProb          float64
	UpperProb          float64
	LowerProbStdDev    *float64 // Optional standard deviation for LowerProb
	UpperProbStdDev    *float64 // Optional standard deviation for UpperProb
	Confidence         float64
	ConfidenceStdDev   *float64 // Optional standard deviation for Confidence
	Timeframe          Timeframe
	Results            []int
	Sum                int
	SumOfSquares       float64
	MinImpact          *float64 // Optional minimum financial impact
	MaxImpact          *float64 // Optional maximum financial impact
	MinImpactStdDev    *float64 // Optional standard deviation for MinImpact
	MaxImpactStdDev    *float64 // Optional standard deviation for MaxImpact
	ImpactSum          float64
	ImpactSumOfSquares float64
}
