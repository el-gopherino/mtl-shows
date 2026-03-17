package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

func main() {
	now := time.Now()

	sequential := flag.Bool("seq", false, "sequential scraping")
	concurrent := flag.Bool("conc", false, "concurrent scraping")
	serve := flag.Bool("serve", false, "run API server")
	flag.Parse()

	if *sequential {
		runSequential()
		fmt.Println("\nruntime duration: ", time.Since(now))
		return
	}
	if *concurrent {
		runConcurrent()
		fmt.Println("\nruntime duration: ", time.Since(now))
		return
	}

	if *serve {
		runConcurrent()
		http.HandleFunc("/events", handleAllEvents)
		http.HandleFunc("/events/right-now", handleRightNow)
		http.HandleFunc("/events/tonight", handleTonight)
		http.HandleFunc("/events/tomorrow", handleTomorrow)
		http.HandleFunc("/events/this-week", handleThisWeek)
		http.HandleFunc("/events/this-weekend", handleThisWeekend)

		fmt.Println("API server running on port :8080")
		http.ListenAndServe(":8080", nil)
		return
	}

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	runOnSchedule(ctx, 1*time.Hour)
}

func runSequential() {
	fmt.Println("running scraper in sequential mode...")
	allEvents := make(map[string]EventList)

	for key, venue := range allVenues {
		fmt.Printf("Scraping %s...\n", venue.Name)

		var events EventList
		if key == "turbo-haus" {
			events = scrapeTurboHausJSON()
		} else {
			events = scrapeVenue(key, venue)
		}
		allEvents[key] = events
	}
	saveAllEvents(allEvents)
}

func runConcurrent() {
	fmt.Println("running scraper in concurrent mode...")
	allEvents := make(map[string]EventList)

	var mu sync.Mutex
	var wg sync.WaitGroup

	for key, venue := range allVenues {
		wg.Add(1)

		go func(k string, v Venue) {
			defer wg.Done()
			fmt.Printf("Scraping %s...\n", v.Name)
			var events EventList
			if k == "turbo-haus" {
				events = scrapeTurboHausJSON()
			} else {
				events = scrapeVenue(k, v)
			}
			mu.Lock()
			allEvents[k] = events
			mu.Unlock()
		}(key, venue)
	}

	wg.Wait()
	saveAllEvents(allEvents)
}

func runOnSchedule(ctx context.Context, interval time.Duration) {
	fmt.Println("running scraper on schedule...")
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	// scrape immediately on startup
	runConcurrent()

	for {
		select {
		case <-ticker.C:
			fmt.Println("Scheduled scrape starting...")
			runConcurrent()
		case <-ctx.Done():
			fmt.Println("Scheduler stopped.")
			return
		}
	}
}
