package main

import (
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
	default:
		e = parseGeneric(h)
	}

	return e, e.validateEvent()
}

func parseCasaDelPopolo(h *colly.HTMLElement) Event {
	children := h.DOM.Children().Filter("div")

	e := Event{
		VenueKey:  "casa-del-popolo",
		Name:      strings.TrimSpace(children.Eq(1).Text()),
		Date:      strings.TrimSpace(children.Eq(0).Text()),
		Venue:     strings.TrimSpace(children.Eq(2).Find("div").First().Text()),
		Address:   strings.TrimSpace(children.Eq(2).Find("div").Last().Text()),
		Time:      strings.TrimSpace(children.Eq(3).Text()),
		Price:     strings.TrimSpace(children.Eq(4).Text()),
		TicketURL: h.ChildAttr("a.btn-inverse", "href"),
	}
	e.enrichEvent()

	return e
}

func parseSalaRossa(h *colly.HTMLElement) Event {
	children := h.DOM.Children().Filter("div")

	e := Event{
		VenueKey:  "la-sala-rossa",
		Name:      strings.TrimSpace(children.Eq(1).Text()),
		Date:      strings.TrimSpace(children.Eq(0).Text()),
		Venue:     strings.TrimSpace(children.Eq(2).Find("div").First().Text()),
		Address:   strings.TrimSpace(children.Eq(2).Find("div").Last().Text()),
		Time:      strings.TrimSpace(children.Eq(3).Text()),
		Price:     strings.TrimSpace(children.Eq(4).Text()),
		TicketURL: h.ChildAttr("a.btn-inverse", "href"),
	}
	e.enrichEvent()

	return e
}

func parseLaSotterenea(h *colly.HTMLElement) Event {
	children := h.DOM.Children().Filter("div")

	e := Event{
		VenueKey:  "la-sotterenea",
		Name:      strings.TrimSpace(children.Eq(1).Text()),
		Date:      strings.TrimSpace(children.Eq(0).Text()),
		Venue:     strings.TrimSpace(children.Eq(2).Find("div").First().Text()),
		Address:   strings.TrimSpace(children.Eq(2).Find("div").Last().Text()),
		Time:      strings.TrimSpace(children.Eq(3).Text()),
		Price:     strings.TrimSpace(children.Eq(4).Text()),
		TicketURL: h.ChildAttr("a.btn-inverse", "href"),
	}
	e.enrichEvent()

	return e
}

func parsePtitOurs(h *colly.HTMLElement) Event {
	children := h.DOM.Children().Filter("div")

	e := Event{
		VenueKey:  "ptit-ours",
		Name:      strings.TrimSpace(children.Eq(1).Text()),
		Date:      strings.TrimSpace(children.Eq(0).Text()),
		Venue:     strings.TrimSpace(children.Eq(2).Find("div").First().Text()),
		Address:   strings.TrimSpace(children.Eq(2).Find("div").Last().Text()),
		Time:      strings.TrimSpace(children.Eq(3).Text()),
		Price:     strings.TrimSpace(children.Eq(4).Text()),
		TicketURL: h.ChildAttr("a.btn-inverse", "href"),
	}
	e.enrichEvent()

	return e
}

func parseLaToscadura(h *colly.HTMLElement) Event {
	children := h.DOM.Children().Filter("div")

	e := Event{
		VenueKey:  "la-toscadura",
		Name:      strings.TrimSpace(children.Eq(1).Text()),
		Date:      strings.TrimSpace(children.Eq(0).Text()),
		Venue:     strings.TrimSpace(children.Eq(2).Find("div").First().Text()),
		Address:   strings.TrimSpace(children.Eq(2).Find("div").Last().Text()),
		Time:      strings.TrimSpace(children.Eq(3).Text()),
		Price:     strings.TrimSpace(children.Eq(4).Text()),
		TicketURL: h.ChildAttr("a.btn-inverse", "href"),
	}
	e.enrichEvent()

	return e
}

func parseCafeCampus(h *colly.HTMLElement) Event {

	rawDateTime := strings.TrimSpace(h.ChildText("span.sh-date"))
	date, eventTime := splitDateTime(rawDateTime)

	priceText := h.ChildText("div.sh-excerpt")
	price := extractAdvancePrice(priceText)

	e := Event{
		VenueKey:  "cafe-campus",
		Name:      strings.TrimSpace(h.ChildText("h4 a")),
		Date:      date,
		Venue:     "Cafe Campus",
		Address:   "57 Rue Prince-Arthur Est",
		Time:      eventTime,
		Price:     price,
		TicketURL: h.ChildAttr("span.sh-address a", "href"),
	}

	e.enrichEvent()
	return e
}

func parseQuaiDesBrumes(h *colly.HTMLElement) Event {
	// prices not shown on page
	e := Event{
		VenueKey:  "quai-des-brumes",
		Name:      h.ChildText("h3.mec-event-title a"),
		Date:      h.ChildText("span.mec-start-date-label"),
		Venue:     "Quai des Brumes",
		Address:   "4481 Rue Saint-Denis",
		Time:      h.ChildText("span.mec-start-time"),
		TicketURL: h.ChildAttr("a.mec-color-hover", "href"),
	}

	e.enrichEvent()
	return e
}

func parseHemisphereGauche(h *colly.HTMLElement) Event {
	e := Event{
		VenueKey: "hemisphere-gauche",
		Name:     h.ChildText("a.WFgzOI"),
		Date:     h.ChildText("span.GiNWmM"),
		Venue:    "L'Hémisphere Gauche",
		Address:  "221 Beaubien Est",
		// skipping time for this one
		TicketURL: h.ChildAttr("a.DjQEyU m022zm aUkG34", "href"),
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

	// Ticket URL from any link to showInfo.php
	ticketURL := h.ChildAttr("a[href*='showInfo']", "href")

	// Event image is in the inline style: background: url("https://...")
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

// parseGeneric is the default for the switch case. Should never have to run.
func parseGeneric(h *colly.HTMLElement) Event {
	return Event{
		Name: strings.TrimSpace(h.Text),
	}
}
