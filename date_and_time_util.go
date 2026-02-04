package main

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// regex to strip ordinal suffixes for date (1st, 2nd, 3rd, 4th, etc.)
var ordinalRegex = regexp.MustCompile(`(\d+)(st|nd|rd|th)\b`)

var loc, _ = time.LoadLocation("America/Montreal")

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

	"Fév", "Feb", "fév", "Feb",
	"Avr", "Apr", "avr", "Apr",
	"Juil", "Jul", "juil", "Jul",
	"Déc", "Dec", "déc", "Dec",
)

func parseDate(date string) (time.Time, error) {

	normalized := translateMonth(date)
	normalized = ordinalRegex.ReplaceAllString(normalized, "$1")
	normalized = strings.TrimSpace(normalized)

	layouts := []struct {
		layout    string
		needsYear bool
	}{
		// Case Popolo : "Wednesday, January 28, 2026" (after removing "th")
		{"Monday, January 2, 2006", false},

		// Case Petit campus: "30 janvier 2026, 20h"
		{"2 January 2026", false},

		// Case quai des brumes : "10 Fév" -> needs year
		{"2 Jan", true},

		// Case Hemisphere gauche :
		// 		get rid of day (substring from after comma)
		//  	format date number?
		{"Fri, Feb 06", true},

		// default:
		{"January 2, 2006", false},
	}

	// strip time if it's in the date
	if idx := strings.Index(normalized, ","); idx != -1 {
		afterComma := normalized[idx+1:]
		if strings.Contains(afterComma, "h") || strings.Contains(afterComma, ":") {
			normalized = strings.TrimSpace(normalized[:idx])
		}
	}

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

func translateMonth(date string) string {
	return frenchMonthReplacer.Replace(date)
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

	thisYear := time.Date(year, t.Month(), t.Day(), 0, 0, 0, 0, t.Location())

	if thisYear.Before(now) {
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

func eventAlreadyHappened(a, b time.Time) bool {
	// todo
	//ad := a.Day()
	//bd := b.Day()
	//
	//if ad > bd {
	//	return true
	//}
	//return false

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

func isThisWeekend(t time.Time) bool {
	day := t.Weekday()
	isWeekendDay := day == time.Saturday || day == time.Sunday

	return isWeekendDay && isThisWeek(t)
}
