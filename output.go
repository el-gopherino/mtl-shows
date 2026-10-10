package main

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

func saveAllEvents(allEvents map[string]EventList) {
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
		path := "output"
		if venue.Group != "" {
			path = fmt.Sprintf("output/%s", venue.Group)
		}
		_ = saveVenueEventsJSON(event, venueKey, venue.Name, path)

		allEventsList = append(allEventsList, event...)
	}

	allEventsList.SortByDate()

	// update cached events
	mu.Lock()
	defer mu.Unlock()
	cachedEvents = allEventsList
	cachedMarkers = buildMarkers(allEventsList)
	lastScrapedAt = time.Now().In(loc)
}

type JSONEnvelope struct {
	Title       string    `json:"title"`
	EventCount  int       `json:"event_count"`
	GeneratedAt string    `json:"generated_at"`
	Events      EventList `json:"events"`
}

func newJSONEnvelope(events EventList, title string) JSONEnvelope {
	return JSONEnvelope{
		Title:       title,
		EventCount:  len(events),
		GeneratedAt: time.Now().In(loc).Format(time.RFC3339),
		Events:      events,
	}
}

func saveJSON(events EventList, title, fileName string) error {
	envelope := newJSONEnvelope(events, title)

	data, err := json.MarshalIndent(envelope, "", "\t")
	if err != nil {
		return fmt.Errorf("error marshaling %s to JSON: %w", title, err)
	}
	return os.WriteFile(fileName, data, 0644)
}

func saveVenueEventsJSON(events EventList, venueKey, venueName, path string) error {
	return saveJSON(events, venueName+" Events", fmt.Sprintf("%s/%s.json", path, venueKey))
}
