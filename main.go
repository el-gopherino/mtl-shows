package main

import (
	"fmt"
	"time"

	"github.com/gocolly/colly/v2"
)

func main() {
	now := time.Now()

	allEvents := make(map[string][]Event)
	for key, venue := range allVenues {

		fmt.Printf("Scraping %s...\n", venue.Name)
		events := scrapeVenue(key, venue)
		allEvents[key] = events

		fmt.Printf("\tFound %d events\n\n", len(events))
	}
	saveAllEvents(allEvents)

	fmt.Println("runtime duration: ", time.Since(now))
}

// TODO: updater pour chaque nouvelle venue
func parseEvent(h *colly.HTMLElement, venueKey string, venue Venue) (Event, []string) {
	var e Event

	switch venueKey {
	case "casa-del-popolo":
		e = parseCasaDelPopolo(h)
	case "la-sala-rossa":
		e = parseSalaRossa(h)
	case "la-sotterenea":
		e = parseLaSotterenea(h)
	case "ptit-ours":
		e = parsePtitOurs(h)
	case "la-toscadura":
		e = parseLaToscadura(h)

		// -------------------------------------

	case "cafe-campus":
		e = parseCafeCampus(h)
	case "quai-des-brumes":
		e = parseQuaiDesBrumes(h)
	default:
		e = parseGeneric(h)
	}
	return e, validateEvent(e)
}

func scrapeVenue(venueKey string, venue Venue) (events []Event) {

	c := colly.NewCollector(
		colly.AllowedDomains(venue.AllowedDomains...),
	)

	c.OnHTML(venue.Selector, func(h *colly.HTMLElement) {
		event, missing := parseEvent(h, venueKey, venue)
		if len(missing) > 0 {
			fmt.Printf("\t[%s] Skipping event, missing: %v\n", venueKey, missing)
		}
		events = append(events, event)
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Printf("Scanning: %s\n", r.URL)
	})

	c.OnError(func(r *colly.Response, e error) {
		fmt.Printf("Error: %s\n", e.Error())
	})

	for _, link := range venue.Links {
		c.Visit(link)
	}

	return events
}
