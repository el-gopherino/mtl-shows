package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"regexp"
	"strings"
	"time"
)

var h2Regex = regexp.MustCompile(`<h2[^>]*>(.*?)</h2>`)
var tagStripper = regexp.MustCompile(`<[^>]*>`)

func scrapeBarLeRitzJSON() EventList {
	client := &http.Client{Timeout: 15 * time.Second}
	resp, err := client.Get(BarLeRitzURL)
	if err != nil {
		log.Printf("[bar-le-ritz] HTTP request failed: %v", err)
		return nil
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("[bar-le-ritz] unexpected status: %d", resp.StatusCode)
		return nil
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("[bar-le-ritz] failed to read body: %v", err)
		return nil
	}

	var sqResp squarespaceResponse
	if err = json.Unmarshal(body, &sqResp); err != nil {
		log.Printf("[bar-le-ritz] JSON parse failed: %v", err)
		return nil
	}

	allItems := mergeSquarespaceItems(sqResp)

	events := make(EventList, 0, len(allItems))
	for _, item := range allItems {
		e := convertBarLeRitzItem(item)
		if e.AlreadyHappened {
			continue
		}
		events = append(events, e)
	}

	fmt.Printf("Scraping for Bar Le Ritz PDB finished.\n")
	return events
}

func convertBarLeRitzItem(item squarespaceItem) Event {
	startTime := time.UnixMilli(item.StartDate).In(loc)

	dateStr := fmt.Sprintf("%s %d, %d",
		startTime.Month().String(),
		startTime.Day(),
		startTime.Year(),
	)
	timeStr := startTime.Format("15:04")

	name := extractEventName(item.Body)

	e := Event{
		VenueKey:   "bar-le-ritz",
		Name:       name,
		Venue:      "Bar Le Ritz PDB",
		Address:    "179 Rue Jean-Talon Ouest",
		Date:       dateStr,
		Time:       timeStr,
		TicketURL:  "https://www.barleritzpdb.com" + item.FullURL,
		EventImage: item.AssetURL,
	}

	e.enrichEvent()
	return e
}

// extractEventName pulls the artist/event name from the first <h2> in the Squarespace body HTML.
func extractEventName(body string) string {
	match := h2Regex.FindStringSubmatch(body)
	if len(match) < 2 {
		return ""
	}
	name := tagStripper.ReplaceAllString(match[1], "")
	return strings.TrimSpace(name)
}
