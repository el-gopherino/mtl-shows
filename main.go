package main

import (
	"fmt"
	"strings"
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

func scrapeVenue(venueKey string, venue Venue) (events []Event) {

	c := colly.NewCollector(
		colly.AllowedDomains(venue.AllowedDomains...),
	)

	c.OnHTML(venue.Selector, func(h *colly.HTMLElement) {
		event, missing := parseEvent(h, venueKey, venue)
		if event.AlreadyHappened {
			return
		}
		if len(missing) > 0 {
			fmt.Printf("\t[%s]: missing: %v\n", venueKey, strings.Join(missing, ", "))
		}
		events = append(events, event)
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Printf("Website: %s\n", r.URL.Host)
	})

	c.OnError(func(r *colly.Response, e error) {
		fmt.Printf("Error: %s\n", e.Error())
	})

	for _, link := range venue.Links {
		c.Visit(link)
	}

	return events
}
