package main

import (
	"encoding/json"
	"html/template"
	"net/http"
	"sync"
)

var (
	cachedEvents EventList
	mu           sync.RWMutex
)

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET")
		if r.Method == http.MethodOptions {
			return
		}
		next.ServeHTTP(w, r)
	})
}

func handlePage(title string, filter func(list EventList) EventList) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		mu.RLock()
		events := filter(cachedEvents)
		mu.RUnlock()

		venueFilter := r.URL.Query().Get("venue")
		if venueFilter != "" {
			events = events.ByVenue(venueFilter)
		}

		tmpl := template.Must(template.ParseFiles("templates/base.html"))
		data := newPageData(title, events)
		data.VenueFilter = venueFilter
		tmpl.ExecuteTemplate(w, "base", data)
	}
}

func handleAllEvents(w http.ResponseWriter, r *http.Request) {
	mu.RLock()
	events := cachedEvents
	mu.RUnlock()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(newJSONEnvelope(events, "All Events"))
}

func handleRightNow(w http.ResponseWriter, r *http.Request) {
	mu.RLock()
	events := cachedEvents.RightNow()
	mu.RUnlock()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(newJSONEnvelope(events, "Right Now"))
}

func handleTonight(w http.ResponseWriter, r *http.Request) {
	mu.RLock()
	events := cachedEvents.Tonight()
	mu.RUnlock()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(newJSONEnvelope(events, "Tonight"))
}

func handleTomorrow(w http.ResponseWriter, r *http.Request) {
	mu.RLock()
	events := cachedEvents.Tomorrow()
	mu.RUnlock()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(newJSONEnvelope(events, "Tomorrow"))
}

func handleThisWeek(w http.ResponseWriter, r *http.Request) {
	mu.RLock()
	events := cachedEvents.ThisWeek()
	mu.RUnlock()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(newJSONEnvelope(events, "This Week"))
}

func handleThisWeekend(w http.ResponseWriter, r *http.Request) {
	mu.RLock()
	events := cachedEvents.ThisWeekend()
	mu.RUnlock()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(newJSONEnvelope(events, "This Weekend"))
}
