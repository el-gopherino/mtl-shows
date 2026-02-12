package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

type squarespaceResponse struct {
	Upcoming []squarespaceItem `json:"upcoming"`
	Past     []squarespaceItem `json:"past"`
	Items    []squarespaceItem `json:"items"`
}

type squarespaceItem struct {
	Title     string               `json:"title"`
	FullURL   string               `json:"fullUrl"`
	StartDate int64                `json:"startDate"` // Unix ms
	EndDate   int64                `json:"endDate"`   // Unix ms
	AssetURL  string               `json:"assetUrl"`  // poster image
	Location  *squarespaceLocation `json:"location"`
	Excerpt   string               `json:"excerpt"`
	Tags      []string             `json:"tags"`
}

type squarespaceLocation struct {
	AddressTitle   string `json:"addressTitle"`
	AddressLine1   string `json:"addressLine1"`
	AddressLine2   string `json:"addressLine2"`
	AddressCountry string `json:"addressCountry"`
}

func scrapeTurboHausJSON() []Event {
	url := TurboHausURL

	client := &http.Client{Timeout: 15 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		log.Printf("[turbo-haus] HTTP request failed: %v", err)
		return nil
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("[turbo-haus] unexpected status: %d", resp.StatusCode)
		return nil
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("[turbo-haus] failed to read body: %v", err)
		return nil
	}

	var sqResp squarespaceResponse
	if err = json.Unmarshal(body, &sqResp); err != nil {
		log.Printf("[turbo-haus] JSON parse failed: %v", err)
		return nil
	}

	// Squarespace may return events in "upcoming", "items", or both.
	// Merge them and deduplicate by URL.
	allItems := mergeSquarespaceItems(sqResp)

	events := make([]Event, 0, len(allItems))
	for _, item := range allItems {
		e := convertSquarespaceItem(item)
		if e.AlreadyHappened {
			continue
		}
		events = append(events, e)
	}

	fmt.Printf("Scraping for Turbo Haus finished.\n")
	return events
}

func convertSquarespaceItem(item squarespaceItem) Event {
	// Squarespace timestamps are Unix milliseconds
	startTime := time.UnixMilli(item.StartDate).In(loc)
	endTime := time.UnixMilli(item.EndDate).In(loc)

	dateStr := fmt.Sprintf("%s %d, %d",
		startTime.Month().String(),
		startTime.Day(),
		startTime.Year(),
	)
	timeStr := startTime.Format("15:04")
	endTimeStr := endTime.Format("15:04")
	_ = endTimeStr // for end time

	venue := "Turbo Haüs"
	address := "2040 Rue Saint-Denis"
	if item.Location != nil {
		if item.Location.AddressTitle != "" {
			venue = item.Location.AddressTitle
		}
		if item.Location.AddressLine1 != "" {
			address = item.Location.AddressLine1
		}
	}

	e := Event{
		VenueKey:   "turbo-haus",
		Name:       item.Title,
		Venue:      venue,
		Address:    address,
		Date:       dateStr,
		Time:       timeStr,
		TicketURL:  "https://www.turbohaus.ca" + item.FullURL,
		EventImage: item.AssetURL,
	}

	e.enrichEvent()
	return e
}

func mergeSquarespaceItems(resp squarespaceResponse) []squarespaceItem {
	seen := make(map[string]struct{})
	var merged []squarespaceItem

	addUnique := func(items []squarespaceItem) {
		for _, item := range items {
			if _, exists := seen[item.FullURL]; !exists {
				seen[item.FullURL] = struct{}{}
				merged = append(merged, item)
			}
		}
	}
	addUnique(resp.Upcoming)
	addUnique(resp.Items)

	return merged
}
