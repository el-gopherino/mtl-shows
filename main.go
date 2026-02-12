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

		var events []Event
		if key == "turbo-haus" {
			events = scrapeTurboHausJSON()
		} else {
			events = scrapeVenue(key, venue)
		}

		allEvents[key] = events
	}

	saveAllEvents(allEvents)
	fmt.Println("runtime duration: ", time.Since(now))
}

func scrapeVenue(venueKey string, venue Venue) (events []Event) {
	events = make([]Event, 0, 20)

	c := colly.NewCollector(
		colly.AllowedDomains(venue.AllowedDomains...),
	)
	c.OnHTML(venue.Selector, func(h *colly.HTMLElement) {
		event, missing := parseEvent(h, venueKey)
		if event.AlreadyHappened {
			return
		}
		if len(missing) > 0 {
			fmt.Printf("\t[%s]: missing: %v\n", venueKey, strings.Join(missing, ", "))
		}
		events = append(events, event)
	})
	c.OnError(func(r *colly.Response, e error) {
		fmt.Printf("Error: %s\n", e.Error())
	})

	for _, link := range venue.Links {
		c.Visit(link)
	}

	return events
}
