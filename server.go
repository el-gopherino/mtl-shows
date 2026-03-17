package main

import (
	"encoding/json"
	"net/http"
	"sync"
)

var (
	cachedEvents EventList
	mu           sync.RWMutex
)

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
