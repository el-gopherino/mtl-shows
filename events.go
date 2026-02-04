package main

import (
	"image"
	"regexp"
	"strconv"
	"strings"
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
	EventImage  *image.Image

	ParsedDate time.Time
	DaysUntil  int
	PriceValue float64

	IsFree        bool
	IsToday       bool
	IsThisWeekend bool
	IsThisWeek    bool
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

		// log.Printf("warning: could not parse date %q for event: %v", e.Date, err)
		return
	}

	e.ParsedDate = parsedDate
	e.DaysUntil = daysUntil(e.ParsedDate)
	e.DayOfWeek = e.ParsedDate.Weekday().String()
	e.IsToday = isToday(e.ParsedDate)
	e.IsThisWeekend = isThisWeekend(e.ParsedDate)
	e.IsThisWeek = e.DaysUntil >= 0 && e.DaysUntil <= 7
}

func validateEvent(e Event) (missing []string) {
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

// parsePrice parses a price string to a float64 value
func parsePrice(priceStr string) float64 {
	priceStr = strings.ToLower(priceStr)
	if strings.Contains(priceStr, "free") || strings.Contains(priceStr, "gratuit") {
		return 0
	}

	re := regexp.MustCompile(`[\d.]+`)
	match := re.FindString(priceStr)
	if match == "" {
		return 0
	}

	price, _ := strconv.ParseFloat(match, 64)
	return price
}
