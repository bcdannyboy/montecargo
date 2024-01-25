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
	Name                            string
	LowerProb                       float64
	UpperProb                       float64
	LowerProbStdDev                 *float64 // Optional standard deviation for LowerProb
	UpperProbStdDev                 *float64 // Optional standard deviation for UpperProb
	Confidence                      float64
	ConfidenceStdDev                *float64 // Optional standard deviation for Confidence
	Timeframe                       Timeframe
	Results                         []int
	MinImpact                       *float64 // Optional minimum financial impact
	MaxImpact                       *float64 // Optional maximum financial impact
	MinImpactStdDev                 *float64 // Optional standard deviation for MinImpact
	MaxImpactStdDev                 *float64 // Optional standard deviation for MaxImpact
	IsCostSaving                    bool
	CostOfImplementationLower       *float64 // Optional lower bound of cost of implementation
	CostOfImplementationUpper       *float64 // Optional upper bound of cost of implementation
	CostOfImplementationLowerStdDev *float64 // Optional standard deviation for lower bound of cost
	CostOfImplementationUpperStdDev *float64 // Optional standard deviation for upper bound of cost

}

type SimulationResult struct {
	EventResults map[string]EventResult
	EventStats   map[string]EventStat // Added field to store event statistics
}

type EventResult struct {
	Sum                int
	SumOfSquares       float64
	ImpactSum          float64
	ImpactSumOfSquares float64
	MinCostLowerBound  float64 // Minimum of the lower bound of cost of implementation
	MaxCostUpperBound  float64 // Maximum of the upper bound of cost of implementation
}

type EventStat struct {
	Probability             float64
	StdDev                  float64
	MinCostOfImplementation float64 // Minimum estimated cost of implementation
	MaxCostOfImplementation float64 // Maximum estimated cost of implementation
}

type Dependency struct {
	EventName string
	Condition string // "happens" or "not happens"
}
