package main

import (
	"html/template"
	"time"
)

type PageData struct {
	Title         string
	Events        EventList
	EventCount    int
	GeneratedAt   string
	LastScrapedAt string
	VenueFilter   string
	VenuesJSON    template.JS // JSON array of venues with coordinates + their events
}

func newPageData(title string, events EventList) PageData {
	return PageData{
		Title:       title,
		Events:      events,
		EventCount:  len(events),
		GeneratedAt: time.Now().In(loc).Format("Monday, January 2 at 3:04PM"),
	}
}
