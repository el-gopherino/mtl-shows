package main

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

//type DateFormat string
//
//const (
//	ISO         DateFormat = "2024-11-30"
//	American    DateFormat = "11/30/2024"
//	European    DateFormat = "30/11/2024"
//	Long        DateFormat = "November 30, 2024"
//	LongOrdinal DateFormat = "November 30 2024"
//	Compact     DateFormat = "20241130"
//	Short       DateFormat = "Nov 30, 2024"
//)

var (
	// regex to strip ordinal suffixes for date (1st, 2nd, 3rd, 4th, etc.)
	ordinalRegex = regexp.MustCompile(`(\d+)(st|nd|rd|th)\b`)
	loc, _       = time.LoadLocation("America/Montreal")
)

var frenchMonthReplacer = strings.NewReplacer(
	"Janvier", "January", "janvier", "January",
	"Février", "February", "février", "February",
	"Mars", "March", "mars", "March",
	"Avril", "April", "avril", "April",
	"Mai", "May", "mai", "May",
	"Juin", "June", "juin", "June",
	"Juillet", "July", "juillet", "July",
	"Août", "August", "août", "August",
	"Septembre", "September", "septembre", "September",
	"Octobre", "October", "octobre", "October",
	"Novembre", "November", "novembre", "November",
	"Décembre", "December", "décembre", "December",

	"janv", "Jan", "Janv", "Jan",
	"Fév", "Feb", "fév", "Feb", "Févr", "Feb", "févr", "Feb",
	"Avr", "Apr", "avr", "Apr",
	"Juil", "Jul", "juil", "Jul",
	"Déc", "Dec", "déc", "Dec",
)

var frenchDayReplacer = strings.NewReplacer(
	"lun.", "Mon", "Lun.", "Mon",
	"mar.", "Tue", "Mar.", "Tue",
	"mer.", "Wed", "Mer.", "Wed",
	"jeu.", "Thu", "Jeu.", "Thu",
	"ven.", "Fri", "Ven.", "Fri",
	"sam.", "Sat", "Sam.", "Sat",
	"dim.", "Sun", "Dim.", "Sun",
)

func parseDate(date string) (time.Time, error) {

	normalized := translateMonth(date)
	normalized = frenchDayReplacer.Replace(normalized)
	normalized = ordinalRegex.ReplaceAllString(normalized, "$1")
	normalized = strings.TrimSpace(normalized)

	layouts := []struct {
		layout    string
		needsYear bool
	}{
		// Case Popolo : "Wednesday, January 28, 2026" (after removing "th")
		{"Monday, January 2, 2006", false},

		// Case Cafe campus: "30 janvier 2026, 20h"
		// date needs to be stripped and date needs to be translated
		{"2 January 2006", false},

		// Case quai des brumes : "10 Fév"
		{"2 Jan", true},

		// Case Hemisphere gauche : "sam. 07 févr"
		{"Mon 02 Jan", true},

		// default:
		{"January 2, 2006", false},
	}

	// strip time if it's in the date
	normalized = stripTime(normalized)

	for _, l := range layouts {
		t, err := time.ParseInLocation(l.layout, normalized, loc)
		if err == nil {
			if l.needsYear {
				t = time.Date(
					inferYear(t),
					t.Month(),
					t.Day(),
					0,
					0,
					0,
					0,
					t.Location(),
				)
			}
			return t, nil
		}
	}
	return time.Time{}, fmt.Errorf("unable to parse date: %q (normalized: %q)", date, normalized)
}

func stripTime(normalized string) string {
	index := strings.Index(normalized, ",")
	if index != -1 {
		afterComma := normalized[index+1:]
		if strings.Contains(afterComma, "h") || strings.Contains(afterComma, ":") {
			return strings.TrimSpace(normalized[:index])
		}
	}
	return normalized
}

func splitDateTime(raw string) (date, eventTime string) {
	parts := strings.Split(raw, ",")
	if len(parts) < 2 {
		return strings.TrimSpace(raw), ""
	}

	lastPart := strings.TrimSpace(parts[len(parts)-1])

	if strings.Contains(lastPart, "h") || strings.Contains(lastPart, ":") {
		date = strings.TrimSpace(strings.Join(parts[:len(parts)-1], ","))
		eventTime = convertFrenchTime(lastPart)
		return
	}

	return strings.TrimSpace(raw), ""
}

// translateMonth translates a month name from French to English, if it needs to be translated
func translateMonth(date string) string {
	return frenchMonthReplacer.Replace(date)
}

func convertFrenchTime(t string) string {
	t = strings.ToLower(strings.TrimSpace(t))
	t = strings.ReplaceAll(t, "h", ":")

	parts := strings.Split(t, ":")
	if len(parts) != 2 {
		return t
	}

	hour, err := strconv.Atoi(parts[0])
	if err != nil {
		return t
	}

	minute := parts[1]
	if minute == "" {
		minute = "00"
	}

	period := "AM"
	if hour >= 12 {
		period = "PM"
		if hour > 12 {
			hour -= 12
		}
	}
	if hour == 0 {
		hour = 12
	}

	return fmt.Sprintf("%d:%s %s", hour, minute, period)
}

func inferYear(t time.Time) int {
	now := time.Now()
	year := now.Year()

	today := time.Date(year, now.Month(), now.Day(), 0, 0, 0, 0, t.Location())
	thisYear := time.Date(year, t.Month(), t.Day(), 0, 0, 0, 0, t.Location())

	if thisYear.Before(today) {
		return year + 1
	}
	return year
}

func daysUntil(eventDate time.Time) int {
	now := time.Now()

	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	eventDay := time.Date(eventDate.Year(), eventDate.Month(), eventDate.Day(), 0, 0, 0, 0, now.Location())

	return int(eventDay.Sub(today).Hours() / 24)
}

func isPast(t time.Time) bool {
	return daysUntil(t) < 0
}

func eventAlreadyHappened(a, b time.Time) bool {
	ad := a.Day()
	bd := b.Day()

	if ad > bd {
		return true
	}
	return false

}

func isSameDay(a, b time.Time) bool {
	ay, am, ad := a.Date()
	by, bm, bd := b.Date()
	return ay == by && am == bm && ad == bd
}

func isToday(t time.Time) bool {
	return isSameDay(t, time.Now())
}

func isThisWeek(t time.Time) bool {
	nowYear, nowWeek := time.Now().ISOWeek()
	tYear, tWeek := t.ISOWeek()

	return nowYear == tYear && nowWeek == tWeek
}

// isThisWeekend includes Friday as the weekend
func isThisWeekend(t time.Time) bool {
	day := t.Weekday()
	isWeekendDay := day == time.Friday || day == time.Saturday || day == time.Sunday
	return isWeekendDay && isThisWeek(t)
}
