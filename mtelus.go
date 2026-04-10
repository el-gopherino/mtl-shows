package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

const mtelusAPIBase = "https://mtelus.com/api/algolia/search?query="

type mtelusResponse struct {
	Hits    []mtelusHit `json:"hits"`
	NbHits  int         `json:"nbHits"`
	NbPages int         `json:"nbPages"`
	Page    int         `json:"page"`
}

type mtelusHit struct {
	Title     string        `json:"title"`
	Slug      string        `json:"slug"`
	ShowTime  int64         `json:"show_time"`
	ShowDate  int64         `json:"show_date"`
	DoorTime  int64         `json:"door_time"`
	Venue     mtelusVenue   `json:"venue"`
	Thumbnail string        `json:"thumbnail"`
	EventTag  string        `json:"event_tag"`
	Genre     []mtelusGenre `json:"genre"`
}

type mtelusVenue struct {
	Code string `json:"code"`
	Slug string `json:"slug"`
	Name string `json:"name"`
}

type mtelusGenre struct {
	Name string `json:"name"`
	Slug string `json:"slug"`
}

type mtelusQuery struct {
	Filters   mtelusFilters `json:"filters"`
	Options   mtelusOpts    `json:"options"`
	IndexName string        `json:"indexName"`
	VenueSlug string        `json:"venueSlug"`
}

type mtelusFilters struct {
	DisplayMode string   `json:"displayMode"`
	Type        []string `json:"type"`
	Search      string   `json:"search"`
}

type mtelusOpts struct {
	HitsPerPage int `json:"hitsPerPage"`
	Page        int `json:"page"`
}

func buildMTelusURL(page int) string {
	q := mtelusQuery{
		Filters: mtelusFilters{
			DisplayMode: "list",
			Type:        []string{"evenko_show", "show"},
			Search:      "",
		},
		Options: mtelusOpts{
			HitsPerPage: 20,
			Page:        page,
		},
		IndexName: "master_evenko_en-CA",
		VenueSlug: "mtelus",
	}

	jsonBytes, _ := json.Marshal(q)
	encoded := base64.StdEncoding.EncodeToString(jsonBytes)
	return mtelusAPIBase + encoded
}

func scrapeMTelusJSON() EventList {
	client := &http.Client{Timeout: 15 * time.Second}
	var allHits []mtelusHit

	// fetch first page to get total pages
	for page := 0; ; page++ {
		url := buildMTelusURL(page)

		resp, err := client.Get(url)
		if err != nil {
			log.Printf("[mtelus] HTTP request failed (page %d): %v", page, err)
			break
		}

		body, err := io.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			log.Printf("[mtelus] failed to read body (page %d): %v", page, err)
			break
		}

		if resp.StatusCode != http.StatusOK {
			log.Printf("[mtelus] unexpected status (page %d): %d", page, resp.StatusCode)
			break
		}

		var mResp mtelusResponse
		if err = json.Unmarshal(body, &mResp); err != nil {
			log.Printf("[mtelus] JSON parse failed (page %d): %v", page, err)
			break
		}

		allHits = append(allHits, mResp.Hits...)

		if page+1 >= mResp.NbPages {
			break
		}
	}

	events := make(EventList, 0, len(allHits))
	for _, hit := range allHits {
		e := convertMTelusHit(hit)
		if e.AlreadyHappened {
			continue
		}
		events = append(events, e)
	}

	fmt.Printf("Scraping for MTelus finished.\n")
	return events
}

func convertMTelusHit(hit mtelusHit) Event {
	showTime := time.Unix(hit.ShowTime, 0).In(loc)

	dateStr := fmt.Sprintf("%s %d, %d",
		showTime.Month().String(),
		showTime.Day(),
		showTime.Year(),
	)
	timeStr := showTime.Format("3:04 PM")

	thumbnail := hit.Thumbnail
	if thumbnail != "" && thumbnail[:2] == "//" {
		thumbnail = "https:" + thumbnail
	}

	ticketURL := fmt.Sprintf("https://mtelus.com/en/events/mtelus/%s", hit.Slug)

	e := Event{
		VenueKey:   "mtelus",
		Name:       hit.Title,
		Venue:      "MTelus",
		Address:    "59 Rue Sainte-Catherine Est",
		Date:       dateStr,
		Time:       timeStr,
		TicketURL:  ticketURL,
		EventImage: thumbnail,
	}

	e.enrichEvent()
	return e
}
