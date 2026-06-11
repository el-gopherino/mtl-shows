package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gocolly/colly/v2"
)

// TODO: update pour chaque nouvelle venue d'ajoutée dans venues.go
func parseEvent(h *colly.HTMLElement, venueKey string) (Event, []string) {
	var e Event

	switch venueKey {
	case "casa-del-popolo", "la-sala-rossa", "la-sotterenea", "ptit-ours", "la-toscadura":
		e = parseCasaGroup(h, venueKey)
	case "cafe-campus":
		e = parseCafeCampus(h)
	case "quai-des-brumes":
		e = parseQuaiDesBrumes(h)
	case "hemisphere-gauche":
		e = parseHemisphereGauche(h)
	case "verre-bouteille":
		e = parseVerreBouteille(h)
	case "piranha-bar":
		e = parsePiranhaBar(h)
	case "club-soda":
		e = parseClubSoda(h)
	case "le-ministere":
		e = parseLeMinistere(h)
	case "fairmount-theatre":
		e = parseFairmountTheatre(h)
	default:
		fmt.Printf("ERROR: no parser for [%q]. Returning empty event...\n", venueKey)
		return Event{}, nil
	}

	return e, e.validateEvent()
}

func parseCasaGroup(h *colly.HTMLElement, venueKey string) Event {
	children := h.DOM.Children().Filter("div")
	eventImage := h.DOM.Parent().Find("img.object-cover").AttrOr("src", "")

	e := Event{
		VenueKey:   venueKey,
		Name:       strings.TrimSpace(children.Eq(1).Text()),
		Date:       strings.TrimSpace(children.Eq(0).Text()),
		Venue:      strings.TrimSpace(children.Eq(2).Find("div").First().Text()),
		Address:    strings.TrimSpace(children.Eq(2).Find("div").Last().Text()),
		Time:       strings.TrimSpace(children.Eq(3).Text()),
		Price:      strings.TrimSpace(children.Eq(4).Text()),
		TicketURL:  h.ChildAttr("a.btn-inverse", "href"),
		EventImage: eventImage,
	}
	e.enrichEvent()
	return e
}

func parseCafeCampus(h *colly.HTMLElement) Event {

	rawDateTime := strings.TrimSpace(h.ChildText("span.sh-date"))
	date, eventTime := splitDateTime(rawDateTime)
	priceText := h.ChildText("div.sh-excerpt")
	price := extractAdvancePrice(priceText)

	imageStyle, _ := h.DOM.Parent().Find("div.noo-thumbnail").Attr("style")
	eventImage := extractBackgroundURL(imageStyle)

	e := Event{
		VenueKey:   "cafe-campus",
		Name:       strings.TrimSpace(h.ChildText("h4 a")),
		Date:       date,
		Venue:      "Cafe Campus",
		Address:    "57 Rue Prince-Arthur Est",
		Time:       eventTime,
		Price:      price,
		TicketURL:  h.ChildAttr("span.sh-address a", "href"),
		EventImage: eventImage,
	}

	e.enrichEvent()
	return e
}

func parseQuaiDesBrumes(h *colly.HTMLElement) Event {
	// prices not shown on page
	e := Event{
		VenueKey:   "quai-des-brumes",
		Name:       h.ChildText("h3.mec-event-title a"),
		Date:       h.ChildText("span.mec-start-date-label"),
		Venue:      "Quai des Brumes",
		Address:    "4481 Rue Saint-Denis",
		Time:       h.ChildText("span.mec-start-time"),
		TicketURL:  h.ChildAttr("a.mec-color-hover", "href"),
		EventImage: h.ChildAttr("div.mec-event-image img", "src"),
	}

	e.enrichEvent()
	return e
}

func parsePiranhaBar(h *colly.HTMLElement) Event {

	ticketPath := h.ChildAttr("a.eventlist-button", "href")
	ticketURL := ""
	if ticketPath != "" {
		ticketURL = "https://www.piranhabar.ca" + ticketPath
	}

	e := Event{
		VenueKey:   "piranha-bar",
		Name:       h.ChildText("h1.eventlist-title a"),
		Date:       strings.TrimSpace(h.DOM.Find("time.event-date").First().Text()),
		Venue:      "Piranha Bar",
		Address:    "680 Rue Sainte-Catherine Ouest",
		Time:       strings.TrimSpace(h.DOM.Find("time.event-time-localized").First().Text()),
		TicketURL:  ticketURL,
		EventImage: h.ChildAttr("a.eventlist-column-thumbnail img", "src"),
	}

	e.enrichEvent()
	return e
}

func parseHemisphereGauche(h *colly.HTMLElement) Event {

	eventImage, _ := h.DOM.Parent().Find("div[data-hook=ev-list-image] img").First().Attr("src")

	e := Event{
		VenueKey: "hemisphere-gauche",
		Name:     h.ChildText("a.WFgzOI"),
		Date:     h.ChildText("span.GiNWmM"),
		Venue:    "L'Hémisphere Gauche",
		Address:  "221 Beaubien Est",
		// skipping time for this one
		TicketURL:  h.ChildAttr("a.DjQEyU m022zm aUkG34", "href"),
		EventImage: eventImage,
	}

	e.enrichEvent()
	return e
}

func parseVerreBouteille(h *colly.HTMLElement) Event {
	name := strings.TrimSpace(h.ChildText("h3.card-title"))

	// Date string: "12 Février à 20h"
	// The date h3 is inside card-content but has no class — it's the second h3
	rawDateTime := strings.TrimSpace(h.ChildText("div.card-content h3"))

	var date, eventTime string
	if strings.Contains(rawDateTime, " à ") {
		parts := strings.SplitN(rawDateTime, " à ", 2)
		date = strings.TrimSpace(parts[0])
		eventTime = convertFrenchTime(strings.TrimSpace(parts[1]))
	} else {
		date = rawDateTime
	}

	ticketURL := h.ChildAttr("a[href*='showInfo']", "href")

	imageStyle := h.ChildAttr("div.card", "style")
	eventImage := extractBackgroundURL(imageStyle)

	e := Event{
		VenueKey:   "verre-bouteille",
		Name:       name,
		Venue:      "Le Verre Bouteille",
		Address:    "2112 Avenue du Mont-Royal Est",
		Date:       date,
		Time:       eventTime,
		TicketURL:  ticketURL,
		EventImage: eventImage,
	}

	e.enrichEvent()
	return e
}

func parseClubSoda(h *colly.HTMLElement) Event {

	ticketURL := h.ChildAttr("a.stretched-link", "href")
	if ticketURL != "" && !strings.HasPrefix(ticketURL, "http") {
		ticketURL = "https://clubsoda.ca" + ticketURL
	}

	e := Event{
		VenueKey:   "club-soda",
		Name:       strings.TrimSpace(h.ChildText("h2.card-title")),
		Date:       strings.TrimSpace(h.ChildText("p.card-subtitle")),
		Venue:      "Club Soda",
		Address:    "1225 Boul. Saint-Laurent",
		TicketURL:  ticketURL,
		EventImage: h.ChildAttr("div.card-img-top img", "src"),
	}

	e.enrichEvent()
	return e
}

func parseLeMinistere(h *colly.HTMLElement) Event {

	ticketURL := h.ChildAttr("a.stretched-link", "href")
	if ticketURL != "" && !strings.HasPrefix(ticketURL, "http") {
		ticketURL = "https://leministere.ca" + ticketURL
	}

	e := Event{
		VenueKey:   "le-ministere",
		Name:       strings.TrimSpace(h.ChildText("h2.card-title")),
		Date:       strings.TrimSpace(h.ChildText("p.card-subtitle")),
		Venue:      "Le Ministère",
		Address:    "4521 Boul. Saint-Laurent",
		TicketURL:  ticketURL,
		EventImage: h.ChildAttr("div.card-img-top img", "src"),
	}

	e.enrichEvent()
	return e
}

func parseFairmountTheatre(h *colly.HTMLElement) Event {

	// artist name is in the description h2, not in h1.eventlist-title (which is just the French date)
	name := strings.TrimSpace(h.DOM.Find("div.eventlist-description h2").First().Text())

	// "8:00 p.m." -> "8:00 PM"
	showTime := strings.TrimSpace(h.DOM.Find("time.event-time-localized-start").First().Text())
	showTime = strings.ReplaceAll(showTime, "p.m.", "PM")
	showTime = strings.ReplaceAll(showTime, "a.m.", "AM")

	ticketURL := h.DOM.Find("div.sqs-block-button-container a.sqs-block-button-element").First().AttrOr("href", "")

	e := Event{
		VenueKey:   "fairmount-theatre",
		Name:       name,
		Date:       strings.TrimSpace(h.DOM.Find("time.event-date").First().Text()),
		Venue:      "Théâtre Fairmount",
		Address:    "5240 Avenue du Parc",
		Time:       showTime,
		TicketURL:  ticketURL,
		EventImage: h.DOM.Find("a.eventlist-column-thumbnail img").First().AttrOr("src", ""),
	}

	e.enrichEvent()
	return e
}

func scrapeTurboHausJSON() (events EventList) {
	client := &http.Client{Timeout: 15 * time.Second}
	resp, err := client.Get(TurboHausURL)
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

	events = make(EventList, 0, len(allItems))
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
