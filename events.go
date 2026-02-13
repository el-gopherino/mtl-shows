package main

import (
	"log"
	"sort"
	"time"
)

type Event struct {
	VenueKey        string // key for venue map
	Name            string
	Venue           string
	Date            string
	Address         string
	Time            string
	Price           string
	TicketURL       string
	DayOfWeek       string
	CalendarURL     string
	ICSData         string
	EventImage      string
	ParsedDate      time.Time
	PriceValue      float64
	DaysUntil       int
	AlreadyHappened bool
	IsFree          bool
	isRightNow      bool
	IsToday         bool
	IsThisWeekend   bool
	IsThisWeek      bool
}

func (e *Event) enrichEvent() {

	e.PriceValue = parsePrice(e.Price)
	e.IsFree = e.PriceValue == 0

	parsedDate, err := parseDate(e.Date)
	if err != nil {
		// if date parsing fails (which it shouldn't), set defaults for date-dependent fields
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
	if e.Time == "" {
		missing = append(missing, "Time")
	}
	if e.ParsedDate.IsZero() {
		missing = append(missing, "ParsedDate")
	}
	return missing
}

// EventList implements sort.Interface
type EventList []Event

func (el EventList) Len() int      { return len(el) }
func (el EventList) Swap(i, j int) { el[i], el[j] = el[j], el[i] }

// Less defaults to sort by date (soonest first)
func (el EventList) Less(i, j int) bool {
	return el[i].ParsedDate.Before(el[j].ParsedDate)
}

func (el EventList) SortByDate() {
	sort.Stable(el)
}

// SortByPrice sorts events by price (cheapest to most expensive) -> maybe use it, since venues don't always show price
func (el EventList) SortByPrice() {
	sort.SliceStable(el, func(i, j int) bool {
		return el[i].PriceValue < el[j].PriceValue
	})
}

// RightNow returns events still happening from 1 an hour ago, and events happening within 2 hours of current time
func (el EventList) RightNow() (result EventList) {
	now := time.Now().In(loc)

	for _, e := range el {
		eventStart := combineDateAndTime(e.ParsedDate, e.Time)
		if eventStart.IsZero() {
			continue // cant parse -> skip
		}

		// Event started up to 1 hour ago (likely still happening)
		startedRecently := eventStart.After(now.Add(-1*time.Hour)) && eventStart.Before(now)

		// Event starts within the next 2 hours
		startsSoon := eventStart.After(now) && eventStart.Before(now.Add(2*time.Hour))

		if startedRecently || startsSoon {
			result = append(result, e)
		}
	}
	return result
}

func (el EventList) Tonight() (result EventList) {
	for _, e := range el {
		if e.IsToday {
			result = append(result, e)
		}
	}
	return result
}

func (el EventList) ThisWeek() (result EventList) {
	for _, e := range el {
		if e.IsThisWeek {
			result = append(result, e)
		}
	}
	return result
}

func (el EventList) ThisWeekend() (result EventList) {
	for _, e := range el {
		if e.IsThisWeekend {
			result = append(result, e)
		}
	}
	return result
}

// Free returns events that are free (no money)
func (el EventList) Free() (result EventList) {
	for _, e := range el {
		if e.IsFree {
			result = append(result, e)
		}
	}
	return result
}

// ByDay return events by the specified day in the day argument
func (el EventList) ByDay(day time.Weekday) (result EventList) {
	for _, e := range el {
		if e.ParsedDate.Weekday() == day {
			result = append(result, e)
		}
	}
	return result
}
