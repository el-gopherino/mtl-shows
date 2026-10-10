package main

import (
	"log"
	"time"
)

type Event struct {
	VenueKey   string `json:"venue_key"`
	Name       string `json:"name"`
	Venue      string `json:"venue"`
	Date       string `json:"date"`
	Address    string `json:"address"`
	Time       string `json:"time,omitempty"`
	Price      string `json:"price,omitempty"`
	TicketURL  string `json:"ticket_url,omitempty"`
	EventImage string `json:"event_image,omitempty"`
	DayOfWeek  string `json:"day_of_week"`
	Genre      string `json:"genre,omitempty"`

	IsFree        bool `json:"is_free"`
	IsToday       bool `json:"is_today"`
	IsTomorrow    bool `json:"is_tomorrow"`
	IsThisWeekend bool `json:"is_this_weekend"`
	IsThisWeek    bool `json:"is_this_week"`

	PriceValue      float64   `json:"-"`
	ParsedDate      time.Time `json:"-"`
	ParsedEndDate   time.Time `json:"-"`
	DaysUntil       int       `json:"-"`
	AlreadyHappened bool      `json:"-"`
}

func (e *Event) enrichEvent() {
	e.PriceValue = parsePrice(e.Price)

	if e.PriceValue == 0 {
		e.IsFree = true
	}

	start, end, err := parseDateRange(e.Date)
	if err != nil {
		// if date parsing fails (which it shouldn't), set defaults for date-dependent fields
		e.ParsedDate = time.Time{}
		e.ParsedEndDate = time.Time{}
		e.DaysUntil = -1
		e.DayOfWeek = ""
		e.IsToday = false
		e.IsTomorrow = false
		e.IsThisWeekend = false
		e.IsThisWeek = false

		log.Printf("WARNING: could not parse date %q for event: %v", e.Date, err)
		return
	}

	e.ParsedDate = start
	e.ParsedEndDate = end
	e.AlreadyHappened = isPast(end)
	e.DaysUntil = daysUntil(e.ParsedDate)
	e.DayOfWeek = start.Weekday().String()

	daysUntilStart := e.DaysUntil
	daysUntilEnd := daysUntil(end)

	e.IsToday = daysUntilStart <= 0 && daysUntilEnd >= 0
	e.IsTomorrow = daysUntilStart <= 1 && daysUntilEnd >= 1
	e.IsThisWeek = daysUntilStart <= 7 && daysUntilEnd >= 0

	e.IsThisWeekend = isThisWeekend(start) || isThisWeekend(end) ||
		(e.IsToday && isThisWeekend(time.Now().In(loc)))
}

// finalizeEvent is the single post-parse step every venue goes through:
// it computes the derived fields and reports which required fields are missing.
func (e *Event) finalizeEvent() (missing []string) {
	e.enrichEvent()

	if e.Name == "" {
		missing = append(missing, "Name")
	}
	if e.Venue == "" {
		missing = append(missing, "Venue")
	}
	if e.Price == "" {
		missing = append(missing, "Price Value")
	}
	if e.Time == "" {
		missing = append(missing, "Time")
	}
	if e.ParsedDate.IsZero() {
		missing = append(missing, "ParsedDate")
	}
	return missing
}
