package testing_utils

import (
	"math/rand"
	"time"

	"github.com/bcdannyboy/montecargo/montecargo"
)

func Float64Pointer(value float64) *float64 {
	return &value
}

func ShuffleArray(events []montecargo.Event) []montecargo.Event {
	rand.Seed(time.Now().UnixNano()) // Initialize the random number generator.

	shuffled := make([]montecargo.Event, len(events))
	copy(shuffled, events)

	for i := range shuffled {
		j := rand.Intn(i + 1)
		shuffled[i], shuffled[j] = shuffled[j], shuffled[i]
	}

	return shuffled
}
