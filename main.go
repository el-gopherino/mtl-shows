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
	serve := flag.Bool("serve", false, "scrapes concurrently first, then runs API server with 1 hour scheduler")
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
		ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
		defer cancel()

		go runOnSchedule(ctx, 1*time.Hour) // run scraper once every hour

		mux := http.NewServeMux()

		mux.HandleFunc("/", handlePage("All Events", func(el EventList) EventList { return el }))
		mux.HandleFunc("/tonight", handlePage("Tonight", EventList.Tonight))
		mux.HandleFunc("/tomorrow", handlePage("Tomorrow", EventList.Tomorrow))
		mux.HandleFunc("/this-week", handlePage("This Week", EventList.ThisWeek))
		mux.HandleFunc("/this-weekend", handlePage("This Weekend", EventList.ThisWeekend))
		mux.HandleFunc("/right-now", handlePage("Right Now", EventList.RightNow))

		//mux.HandleFunc("/events", handleAllEvents)
		//mux.HandleFunc("/events/right-now", handleRightNow)
		//mux.HandleFunc("/events/tonight", handleTonight)
		//mux.HandleFunc("/events/tomorrow", handleTomorrow)
		//mux.HandleFunc("/events/this-week", handleThisWeek)
		//mux.HandleFunc("/events/this-weekend", handleThisWeekend)

		srv := &http.Server{
			Addr:    ":8080",
			Handler: corsMiddleware(mux),
		}

		// allows to kill scheduler and port with ctrl+C
		go func() {
			<-ctx.Done()
			srv.Shutdown(context.Background())
		}()

		fmt.Println("API server running on port :8080")
		srv.ListenAndServe()
		return
	}
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
	now := time.Now()
	fmt.Println("running scraper in concurrent mode...")
	allEvents := make(map[string]EventList)

	var scrapeMu sync.Mutex
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
			scrapeMu.Lock()
			allEvents[k] = events
			scrapeMu.Unlock()
		}(key, venue)
	}

	wg.Wait()
	saveAllEvents(allEvents)

	fmt.Println("\nScraping of all venues complete.")
	fmt.Printf("Scraping took %v\n", time.Since(now))
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
