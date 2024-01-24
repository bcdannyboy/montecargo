package montecargo

func findEventByName(events []Event, name string) (*Event, bool) {
	for _, event := range events {
		if event.Name == name {
			return &event, true
		}
	}
	return nil, false
}

func filterIndependentEvents(events []Event, dependencies map[string][]Dependency) []Event {
	var independentEvents []Event
	for _, event := range events {
		if _, exists := dependencies[event.Name]; !exists {
			independentEvents = append(independentEvents, event)
		}
	}
	return independentEvents
}

func filterDependentEvents(events []Event, dependencies map[string][]Dependency) []Event {
	var dependentEvents []Event
	for _, event := range events {
		if _, exists := dependencies[event.Name]; exists {
			dependentEvents = append(dependentEvents, event)
		}
	}
	return dependentEvents
}

func TimeframeToString(tf Timeframe) string {
	switch tf {
	case Daily:
		return "1 day"
	case Weekly:
		return "1 week"
	case Monthly:
		return "1 month"
	case Yearly:
		return "1 year"
	case EveryTwoYears:
		return "2 years"
	case EveryFiveYears:
		return "5 years"
	case EveryTenYears:
		return "10 years"
	default:
		return "unknown timeframe"
	}
}
