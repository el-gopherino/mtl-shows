package main

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/gocolly/colly/v2"
)

func main() {

	// runSequential()
	runConcurrent()
}

func runSequential() {
	now := time.Now()
	allEvents := make(map[string]EventList)

	for key, venue := range allVenues {
		fmt.Printf("Scraping %s...\n", venue.Name)

		var events EventList
		if key == "turbo-haus" {
			events = scrapeTurboHausJSON()
		} else {
			events = scrapeVenue(key, venue)
		}
		allEvents[key] = events
	}

	saveAllEvents(allEvents)
	fmt.Println("\nruntime duration: ", time.Since(now))
}

func runConcurrent() {

	now := time.Now()
	allEvents := make(map[string]EventList)

	var mu sync.Mutex
	var wg sync.WaitGroup

	for key, venue := range allVenues {
		wg.Add(1)

		go func(k string, v Venue) {
			defer wg.Done()
			fmt.Printf("Scraping %s...\n", v.Name)
			var events EventList
			if k == "turbo-haus" {
				events = scrapeTurboHausJSON()
			} else {
				events = scrapeVenue(k, v)
			}
			mu.Lock()
			allEvents[k] = events
			mu.Unlock()
		}(key, venue)
	}

	wg.Wait()
	saveAllEvents(allEvents)
	fmt.Println("\nruntime duration: ", time.Since(now))

}

// scrapeVenue scrapes the HTML content of a Venue, in order to parse it in parsers.go
func scrapeVenue(venueKey string, venue Venue) (events EventList) {
	events = make(EventList, 0, 20)

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
	c.OnScraped(func(r *colly.Response) {
		fmt.Printf("Scraping for %s finished.\n", venue.Name)
	})

	for _, link := range venue.Links {
		c.Visit(link)
	}

	return events
}
