package main

import (
	"strings"

	"github.com/gocolly/colly/v2"
)

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

	e := Event{
		VenueKey:  "cafe-campus",
		Name:      strings.TrimSpace(h.ChildText("h4 a")),
		Date:      date,
		Venue:     "Cafe Campus",
		Address:   "57 Rue Prince-Arthur Est",
		Time:      eventTime,
		TicketURL: h.ChildAttr("span.sh-address a", "href"),
	}

	e.enrichEvent()
	return e
}

func parseQuaiDesBrumes(h *colly.HTMLElement) Event {

	e := Event{
		VenueKey:  "quai-des-brumes",
		Name:      h.ChildText("h3.mec-event-title a"),
		Date:      h.ChildText("span.mec-start-date-label"),
		Venue:     "Quai des Brumes",
		Address:   "4481 Rue Saint-Denis",
		Time:      h.ChildText("span.mec-start-time"),
		TicketURL: h.ChildAttr("a.mec-booking-button", "href"),
	}

	e.enrichEvent()
	return e
}

func parseGeneric(h *colly.HTMLElement) Event {
	return Event{
		Name: strings.TrimSpace(h.Text),
	}
}
