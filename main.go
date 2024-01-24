package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/bcdannyboy/montecargo/montecargo"
)

func float64Pointer(value float64) *float64 {
	return &value
}

func main() {
	rand.Seed(time.Now().UnixNano())

	// Define cybersecurity events with various timeframes
	events := []montecargo.Event{
		{"Phishing Attack", 0.70, 0.90, nil, nil, 0.90, montecargo.Yearly, nil, 0, 0},
		{"Ransomware Attack", 0.15, 0.45, nil, nil, 0.85, montecargo.Yearly, nil, 0, 0},
		{"Data Breach", 0.05, 0.25, nil, nil, 0.80, montecargo.EveryTwoYears, nil, 0, 0},
		{"Breach Detection", 0.50, 0.70, nil, nil, 0.95, montecargo.Monthly, nil, 0, 0},
		{"System Compromise", 0.10, 0.30, nil, nil, 0.75, montecargo.EveryFiveYears, nil, 0, 0},
		{"Compliance Failure", 0.03, 0.20, nil, nil, 0.90, montecargo.EveryTenYears, nil, 0, 0},
		{"Reputation Damage", 0.20, 0.40, nil, nil, 0.85, montecargo.Yearly, nil, 0, 0},
		// Survey-based events with standard deviations
		{"Insider Threat", 0.05, 0.15, float64Pointer(0.01), float64Pointer(0.02), 0.80, montecargo.EveryTwoYears, nil, 0, 0},
		{"DDoS Attack", 0.30, 0.60, float64Pointer(0.05), float64Pointer(0.07), 0.90, montecargo.Monthly, nil, 0, 0},
		{"Regular Audit", 0.90, 0.95, nil, nil, 0.99, montecargo.Yearly, nil, 0, 0},
		{"Emergency Patching", 0.40, 0.70, nil, nil, 0.85, montecargo.Weekly, nil, 0, 0},
		{"User Training Failure", 0.10, 0.30, nil, nil, 0.70, montecargo.Yearly, nil, 0, 0},
	}

	// Perform Monte Carlo Simulation
	numSimulations := 10000000
	montecargo.MonteCarloSimulation(&events, numSimulations)

	// Output the results and calculate standard deviation
	for _, event := range events {
		probability, stdDev := montecargo.MeanSTD(event, numSimulations)

		// Calculate the upper and lower bounds
		lowerBound := probability - stdDev
		if lowerBound < 0 {
			lowerBound = 0 // Ensure the lower bound is not negative
		}
		upperBound := probability + stdDev
		if upperBound > 1 {
			upperBound = 1 // Ensure the upper bound does not exceed 100%
		}

		// Print the results
		fmt.Printf("Event: %s\n", event.Name)
		fmt.Printf("  Probability: %.2f%%\n", probability*100)
		fmt.Printf("  Standard Deviation: %.2f%%\n", stdDev*100)
		fmt.Printf("  Chance within timeframe: %.2f%% to %.2f%%\n", lowerBound*100, upperBound*100)
	}

}
