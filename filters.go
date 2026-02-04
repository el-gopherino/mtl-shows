package main

import (
	"sort"
	"strings"
)

func FilterByVenueAlphabetically(allEvents []Event) {
	sort.Slice(allEvents, func(i, j int) bool {
		return strings.ToLower(allEvents[i].Venue) < strings.ToLower(allEvents[j].Venue)
	})
}

// FilterByPrice sorts events by ascending price
func FilterByPrice(allEvents []Event) {
	sort.Slice(allEvents, func(i, j int) bool {
		return allEvents[i].PriceValue < (allEvents[j].PriceValue)
	})
}

// todo
// Filter by day (i.e. Saturdays, Thursdays, etc.)

// FilterByDay return a slice of events by the desired day of the month(i.e. Saturdays, Thursdays, etc.)
//func FilterByDay(allEvents []Event, day time.Weekday) []Event {
//
//}
