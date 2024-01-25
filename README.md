# montecargo

`montecargo` is a Monte Carlo simulation tool implemented in Go (Golang), designed for assessing and quantifying risks in various scenarios, particularly useful in cybersecurity risk analysis. This package allows users to simulate different events with varying probabilities and impacts, providing a statistical approach to risk management.

## Features

- **Event-Based Simulations:** Simulate a wide range of events with customizable probabilities and impacts.
- **Survey Support:** Import probabilities from survey data with standard deviations and confidence intervals.
- **Multivariate Event Dependencies:** Define multiple dependencies between events to simulate cascading effects. ('happens' or 'not happens' conditions)
- **Timeframe Adjustments:** Adjust event probabilities based on different timeframes (e.g., daily, yearly).
- **Impact Analysis:** Calculate financial impacts of events, including mean and standard deviation.
- **Implementation Cost Analysis:** Calculate the cost of implementing preventive measures and cost-saving events.
- **Cost Savings Analysis:** Evaluate the financial benefits of preventive measures and cost-saving events. This number is impacted by the cost of implementation for cost saving events.
- **Concurrency Support:** Leverages Go's concurrency features for efficient simulation over multiple CPU cores.

## Installation

To install `montecargo`, use the following `go get` command:

```$ go get github.com/bcdannyboy/montecargo/montecargo```

## The `Event` Type

The core of montecargo is the Event type, which represents a potential event that can occur in a simulation. Each Event includes several fields:

    - *Name*: A descriptive name of the event.
    - *LowerProb* and *UpperProb*: Define the lower and upper bounds of the event's probability.
    - *LowerProbStdDev* and *UpperProbStdDev*: (Optional) Standard deviations for the lower and upper probability bounds.
    - *Confidence*: A confidence level for the event's probability.
    - *ConfidenceStdDev*: (Optional) Standard deviation for the confidence level.
    - *Timeframe*: The timeframe over which the event probability is considered (e.g., Yearly, Monthly).
    - *MinImpact* and *MaxImpact*: (Optional) Define the minimum and maximum financial impacts of the event.
    - *MinImpactStdDev* and *MaxImpactStdDev*: (Optional) Standard deviations for the minimum and maximum impacts.
    - *IsCostSaving*: Indicates if the event is a cost-saving measure.
    - *CostOfImplementationLower* and *CostOfImplementationUpper* (Optional): The lower and upper bounds of cost of implementing the event (e.g., cost of a security control, will offset the control's overall cost savings).
    - *CostOfImplementationLowerStdDev* and *CostOfImplementationUpperStdDev* (Optional): Standard deviation for the cost of implementation.

##  TimeFrames

montecargo supports various timeframes for event probability calculations. The available timeframes are:

- *Daily*: Event probability adjusted for daily occurrence.
- *Weekly*: Event probability adjusted for weekly occurrence.
- *Monthly*: Event probability adjusted for monthly occurrence.
- *Yearly*: Event probability considered on a yearly basis.
- *EveryTwoYears*: Event probability adjusted for occurrence every two years.
- *EveryFiveYears*: Event probability adjusted for occurrence every five years.
- *EveryTenYears*: Event probability adjusted for occurrence every ten years.

## Dependency Types

In `montecargo`, dependencies between events are a crucial aspect of the simulation. They allow for the modeling of complex scenarios where the occurrence of one event can influence the likelihood of another. There are two primary types of dependencies that can be defined:

- **Happens Dependency**: This type of dependency indicates that the occurrence of one event increases the likelihood of another event happening. For example, if a "Data Breach" event happens, it might increase the probability of a "Ransomware Attack".

- **Not Happens Dependency**: Conversely, this dependency type suggests that the non-occurrence of one event affects the likelihood of another event. For instance, if a "Breach Detection" event does not happen, it could increase the chances of a successful "Ransomware Attack".

### Defining Dependencies

Dependencies are defined in a map where the key is the name of the dependent event, and the value is a slice of `montecargo.Dependency` structs. Each `Dependency` struct includes the `EventName` and the `Condition` (either "happens" or "not happens"). Here's an example of defining dependencies:

    ```
    dependencies := map[string][]montecargo.Dependency{
        "Ransomware Attack": {
            {EventName: "Data Breach", Condition: "happens"},          // Ransomware attack depends on data breach happening
            {EventName: "Breach Detection", Condition: "not happens"}, // Ransomware attack depends on breach detection not happening
        },
        // ... other dependencies ...
    }
    ```

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

