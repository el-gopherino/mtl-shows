package main

import (
	"fmt"
	"os"
	"time"
)

func saveAllEvents(allEvents map[string]EventList) {

	// Create group subdirectories
	for _, venue := range allVenues {
		if venue.Group != "" {
			if err := os.MkdirAll(fmt.Sprintf("output/%s", venue.Group), 0755); err != nil {
				fmt.Printf("Failed to create group directory %s: %v\n", venue.Group, err)
			}
		}
	}

	totalEvents := 0
	for _, events := range allEvents {
		totalEvents += len(events)
	}
	allEventsList := make(EventList, 0, totalEvents)

	for venueKey, event := range allEvents {
		venue := allVenues[venueKey]

		// Use group subdir if venue belongs to a group (rare)
		path := "output"
		if venue.Group != "" {
			path = fmt.Sprintf("output/%s", venue.Group)
		}
		// saveAllEventsToTextFile(event, fmt.Sprintf("%s/%s.txt", path, venueKey), venue.Name)
		saveVenueEventsJSON(event, venueKey, venue.Name, path)

		allEventsList = append(allEventsList, event...)
	}

	allEventsList.SortByDate()

	// update cached events for the API
	mu.Lock()
	cachedEvents = allEventsList
	cachedMarkers = buildMarkers(allEventsList)
	lastScrapedAt = time.Now().In(loc)
	mu.Unlock()
}
