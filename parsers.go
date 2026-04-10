package main

import (
	"fmt"
	"strings"

	"github.com/gocolly/colly/v2"
)

// TODO: update pour chaque nouvelle venue d'ajoutée dans venues.go
func parseEvent(h *colly.HTMLElement, venueKey string) (Event, []string) {
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
	default:
		fmt.Printf("ERROR: no parser for [%q]. Returning empty event...\n", venueKey)
		return Event{}, nil
	}

	return e, e.validateEvent()
}

func parseCasaDelPopolo(h *colly.HTMLElement) Event {
	children := h.DOM.Children().Filter("div")
	eventImage := h.DOM.Parent().Find("img.object-cover").AttrOr("src", "")

	e := Event{
		VenueKey:   "casa-del-popolo",
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

func parseSalaRossa(h *colly.HTMLElement) Event {
	children := h.DOM.Children().Filter("div")
	eventImage := h.DOM.Parent().Find("img.object-cover").AttrOr("src", "")

	e := Event{
		VenueKey:   "la-sala-rossa",
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

func parseLaSotterenea(h *colly.HTMLElement) Event {
	children := h.DOM.Children().Filter("div")
	eventImage := h.DOM.Parent().Find("img.object-cover").AttrOr("src", "")

	e := Event{
		VenueKey:   "la-sotterenea",
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

func parsePtitOurs(h *colly.HTMLElement) Event {
	children := h.DOM.Children().Filter("div")
	eventImage := h.DOM.Parent().Find("img.object-cover").AttrOr("src", "")

	e := Event{
		VenueKey:   "ptit-ours",
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

func parseLaToscadura(h *colly.HTMLElement) Event {
	children := h.DOM.Children().Filter("div")
	eventImage := h.DOM.Parent().Find("img.object-cover").AttrOr("src", "")

	e := Event{
		VenueKey:   "la-toscadura",
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

	date, eventTime := "", ""
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
