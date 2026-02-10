package main

import (
	"log"
	"time"
)

type Event struct {
	VenueKey    string // for maps
	Name        string
	Venue       string
	Date        string
	Address     string
	Time        string
	Price       string
	TicketURL   string
	DayOfWeek   string
	CalendarURL string
	ICSData     string
	EventImage  string

	ParsedDate time.Time
	DaysUntil  int
	PriceValue float64

	AlreadyHappened bool
	IsFree          bool
	IsToday         bool
	IsThisWeekend   bool
	IsThisWeek      bool
}

func (e *Event) enrichEvent() {

	e.PriceValue = parsePrice(e.Price)
	e.IsFree = e.PriceValue == 0

	parsedDate, err := parseDate(e.Date)
	if err != nil {
		// set defaults for date-dependent fields
		e.ParsedDate = time.Time{}
		e.DaysUntil = -1
		e.DayOfWeek = ""
		e.IsToday = false
		e.IsThisWeekend = false
		e.IsThisWeek = false

		log.Printf("warning: could not parse date %q for event: %v", e.Date, err)
		return
	}

	e.ParsedDate = parsedDate
	e.AlreadyHappened = isPast(parsedDate)
	e.DaysUntil = daysUntil(e.ParsedDate)
	e.DayOfWeek = e.ParsedDate.Weekday().String()
	e.IsToday = isToday(e.ParsedDate)
	e.IsThisWeekend = isThisWeekend(e.ParsedDate)
	e.IsThisWeek = e.DaysUntil >= 0 && e.DaysUntil <= 7
}

func (e *Event) validateEvent() (missing []string) {
	if e.Name == "" {
		missing = append(missing, "Name")
	}
	if e.Venue == "" {
		missing = append(missing, "Venue")
	}
	if e.PriceValue == 0 {
		missing = append(missing, "Price Value")
	}
	if e.ParsedDate.IsZero() {
		missing = append(missing, "ParsedDate")
	}
	return missing
}
