package main

import (
	"context"
	"flag"
	"fmt"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

func main() {
	now := time.Now()

	sequential := flag.Bool("seq", false, "sequential scraping")
	concurrent := flag.Bool("conc", false, "concurrent scraping")
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
