# montecargo

`montecargo` is a Monte Carlo simulation tool implemented in Go (Golang), designed for assessing and quantifying risks in various scenarios, particularly useful in cybersecurity risk analysis. This package allows users to simulate different events with varying probabilities and impacts, providing a statistical approach to risk management.

## Features

- **Event-Based Simulations:** Simulate a wide range of events with customizable probabilities and impacts.
- **Survey Support:** Import probabilities from survey data with standard deviations and confidence intervals.
- **Multivariate Event Dependencies:** Define multiple dependencies between events to simulate cascading effects. ('happens' or 'not happens' conditions)
- **Timeframe Adjustments:** Adjust event probabilities based on different timeframes (e.g., daily, yearly).
- **Impact Analysis:** Calculate financial impacts of events, including mean and standard deviation.
- **Concurrency Support:** Leverages Go's concurrency features for efficient simulation over multiple CPU cores.

## Installation

To install `montecargo`, use the following `go get` command:

```$ go get github.com/bcdannyboy/montecargo/montecargo```

## The `Event` Type

The core of montecargo is the Event type, which represents a potential event that can occur in a simulation. Each Event includes several fields:

- *Name*: A descriptive name of the event.
- *LowerProb* and UpperProb: Define the lower and upper bounds of the event's probability.
- *LowerProbStdDev* and UpperProbStdDev: (Optional) Standard deviations for the lower and upper probability bounds.
- *Confidence*: A confidence level for the event's probability.
- *ConfidenceStdDev*: (Optional) Standard deviation for the confidence level.
- *Timeframe*: The timeframe over which the event probability is considered (e.g., Yearly, Monthly).
- *MinImpact* and *MaxImpact*: (Optional) Define the minimum and maximum financial impacts of the event.
- *MinImpactStdDev* and *MaxImpactStdDev*: (Optional) Standard deviations for the minimum and maximum impacts.

##  TimeFrames

montecargo supports various timeframes for event probability calculations. The available timeframes are:

- *Daily*: Event probability adjusted for daily occurrence.
- *Weekly*: Event probability adjusted for weekly occurrence.
- *Monthly*: Event probability adjusted for monthly occurrence.
- *Yearly*: Event probability considered on a yearly basis.
- *EveryTwoYears*: Event probability adjusted for occurrence every two years.
- *EveryFiveYears*: Event probability adjusted for occurrence every five years.
- *EveryTenYears*: Event probability adjusted for occurrence every ten years.

# Usage

## Basic Usage

1. import the package

    `import "github.com/bcdannyboy/montecargo/montecargo"`

2. Define Events

    ```
    events := []montecargo.Event{
        {
            Name:       "Phishing Attack",
            LowerProb:  0.70,
            UpperProb:  0.90,
            Confidence: 0.90,
            Timeframe:  montecargo.Yearly,
        },
        // ... other events ...
    }
    ```

3. Define Dependencies

    ```
    dependencies := map[string][]montecargo.Dependency{
		"Ransomware Attack": {
			{EventName: "Data Breach", Condition: "happens"},          // ransomware attack depends on data breach happening
			{EventName: "Breach Detection", Condition: "not happens"}, // ransomware attack depends on breach detection not happening
		},
		// ... other dependencies ...
	}
    ```
4. Run the simulation

    ```
    numSimulations := 100000
	montecargo.MonteCarloSimulation(&events, numSimulations, dependencies)
    ```

5. Analyze Results

    ```
    for _, event := range events {
        probability, _, impactMean, _ := montecargo.MeanSTD(event, numSimulations)
        fmt.Printf("Event: %s, Probability: %.2f, Mean Impact: $%.2f\n", event.Name, probability, impactMean)
    }
    ```

## Advanced Usage

For more advanced scenarios, including events with standard deviations for probabilities and impacts, refer to the provided example in the main package.

