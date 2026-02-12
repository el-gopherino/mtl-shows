package main

import "time"

type PageData struct {
	Title       string
	Events      []Event
	EventCount  int
	GeneratedAt string
}

func newPageData(title string, events []Event) PageData {
	return PageData{
		Title:       title,
		Events:      events,
		EventCount:  len(events),
		GeneratedAt: time.Now().In(loc).Format("Monday, January 2 at 3:04PM"),
	}
}

// todo: for event image parsing
func saveAllEventsHTML(allEvents []Event) error {
	return nil
}
