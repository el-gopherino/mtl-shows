package main

import (
	"sort"
	"strings"
	"time"
)

type EventList []Event

func (el EventList) Len() int           { return len(el) }
func (el EventList) Swap(i, j int)      { el[i], el[j] = el[j], el[i] }
func (el EventList) Less(i, j int) bool { return el[i].ParsedDate.Before(el[j].ParsedDate) } // defaults to sort by date soonest to latest

func (el EventList) SortByDate() {
	sort.Stable(el)
}

// SortByPrice sorts events by price (cheapest to most expensive) -> maybe use it, since venues don't always show price
func (el EventList) SortByPrice() {
	sort.SliceStable(el, func(i, j int) bool {
		return el[i].PriceValue < el[j].PriceValue
	})
}

// RightNow returns events still happening from 1 hour ago, and events happening within 2 hours of current time
func (el EventList) RightNow() (result EventList) {
	now := time.Now().In(loc)

	for _, e := range el {
		eventStart := combineDateAndTime(e.ParsedDate, e.Time)
		if eventStart.IsZero() {
			continue // cant parse, skip.
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

func (el EventList) Tomorrow() (result EventList) {
	for _, e := range el {
		if e.IsTomorrow {
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

// todo: remplacer par ByPrice
func (el EventList) Free() (result EventList) {
	for _, e := range el {
		if e.IsFree {
			result = append(result, e)
		}
	}
	return result
}

// ByWeekday return events by the specified day in the day argument
func (el EventList) ByWeekday(day time.Weekday) (result EventList) {
	for _, e := range el {
		if e.ParsedDate.Weekday() == day {
			result = append(result, e)
		}
	}
	return result
}

// ByVenue filters by venue
func (el EventList) ByVenue(venueKey string) (result EventList) {
	for _, e := range el {
		if e.VenueKey == venueKey {
			result = append(result, e)
		}
	}
	return result
}

// ByGenre filters to events whose comma-separated Genre field contains the
// given genre as one of its own tokens, case-insensitively -- same shape as
// ByVenue above (a plain scan, no separate index), but a substring/token
// match rather than exact-equality: Event.Genre can genuinely hold more
// than one genre (see that field's own comment), so "techno" has to match
// an event tagged "Techno, House" without also matching "Technoise" as a
// false positive -- splitting on comma and comparing whole tokens handles
// both. Added for the sidebar's new genre tag filter (spec item 6) -- same
// ?genre= query param / handlePage wiring as ByVenue's own ?venue=, see
// server.go.
func (el EventList) ByGenre(genre string) (result EventList) {
	g := strings.ToLower(strings.TrimSpace(genre))
	for _, e := range el {
		for tok := range strings.SplitSeq(e.Genre, ",") {
			if strings.ToLower(strings.TrimSpace(tok)) == g {
				result = append(result, e)
				break
			}
		}
	}
	return result
}

// DistinctGenres collects every distinct genre token actually present
// across el's own Event.Genre fields (each one comma-separated, see that
// field's own comment), case-insensitively deduplicated but keeping the
// first-seen casing, sorted alphabetically -- this is .Genres on PageData
// (server.go), the sidebar's own genre tag list (frontend/base.html,
// underneath "Add a venue"). Built fresh from whatever's actually been
// scraped, same as VenueOption's own venueList (output.go) -- an event
// source that exposes no genre info just contributes nothing here, never
// a hardcoded/fake tag.
func DistinctGenres(el EventList) []string {
	seen := make(map[string]string) // lowercase token -> first-seen casing
	for _, e := range el {
		for tok := range strings.SplitSeq(e.Genre, ",") {
			t := strings.TrimSpace(tok)
			if t == "" {
				continue
			}
			key := strings.ToLower(t)
			if _, ok := seen[key]; !ok {
				seen[key] = t
			}
		}
	}
	genres := make([]string, 0, len(seen))
	for _, v := range seen {
		genres = append(genres, v)
	}
	sort.Strings(genres)
	return genres
}

// ByName filters to events whose Name (the artist/show name) contains query (case-insensitive).
func (el EventList) ByName(query string) (result EventList) {
	q := strings.ToLower(strings.TrimSpace(query))
	for _, e := range el {
		if strings.Contains(strings.ToLower(e.Name), q) {
			result = append(result, e)
		}
	}
	return result
}

// GroupByVenue indexes events by VenueKey in a single pass.
func (el EventList) GroupByVenue() map[string]EventList {
	result := make(map[string]EventList, len(allVenues))
	for _, e := range el {
		result[e.VenueKey] = append(result[e.VenueKey], e)
	}
	return result
}
