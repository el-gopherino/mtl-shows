package main

import (
	"sort"
	"time"
)

// FilterByPrice sorts events by ascending price
func FilterByPrice(allEvents []Event) {
	sort.Slice(allEvents, func(i, j int) bool {
		return allEvents[i].PriceValue < (allEvents[j].PriceValue)
	})
}

// FilterByDay filters by specific day (i.e. Saturdays, Thursdays, etc.).
func FilterByDay(allEvents []Event, day time.Weekday) {
	// todo
}
