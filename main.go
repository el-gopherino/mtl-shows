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

var (
	scrapeSchedule = 1 * time.Hour // scrape every hour

	// command-line flags
	sequential = flag.Bool("seq", false, "sequential scraping")
	concurrent = flag.Bool("conc", false, "concurrent scraping")
	serve      = flag.Bool("serve", false, "scrapes concurrently first, then runs API server with 1 hour scheduler")
)

func main() {
	now := time.Now()
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

		go runOnSchedule(ctx, scrapeSchedule)

		mux := http.NewServeMux()

		mux.HandleFunc("/", handlePage("All Events", func(el EventList) EventList { return el }))
		mux.HandleFunc("/right-now", handlePage("Right Now", EventList.RightNow))
		mux.HandleFunc("/tonight", handlePage("Tonight", EventList.Tonight))
		mux.HandleFunc("/tomorrow", handlePage("Tomorrow", EventList.Tomorrow))
		mux.HandleFunc("/this-week", handlePage("This Week", EventList.ThisWeek))
		mux.HandleFunc("/this-weekend", handlePage("This Weekend", EventList.ThisWeekend))

		mux.HandleFunc("/this-weekend/friday", handlePage("This Weekend — Friday",
			func(el EventList) EventList { return el.ThisWeekend().ByWeekday(time.Friday) }))
		mux.HandleFunc("/this-weekend/saturday", handlePage("This Weekend — Saturday",
			func(el EventList) EventList { return el.ThisWeekend().ByWeekday(time.Saturday) }))
		mux.HandleFunc("/this-weekend/sunday", handlePage("This Weekend — Sunday",
			func(el EventList) EventList { return el.ThisWeekend().ByWeekday(time.Sunday) }))

		srv := &http.Server{
			Addr:    ":8080",
			Handler: corsMiddleware(mux),
		}

		// allows to kill scheduler and port with ctrl+C
		go func() {
			<-ctx.Done()
			if err := srv.Shutdown(context.Background()); err != nil {
				fmt.Println("error shutting down application.")
				return
			}
		}()

		fmt.Println("API server running on port :8080")
		if err := srv.ListenAndServe(); err != nil {
			return
		}
		return
	}
}

func runOnSchedule(ctx context.Context, interval time.Duration) {
	fmt.Println("running scraper on schedule...")
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	// scrape immediately
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

			switch k {
			case "turbo-haus":
				events = scrapeTurboHausJSON()
			case "bar-le-ritz":
				events = scrapeBarLeRitzJSON()
			case "mtelus":
				events = scrapeMTelusJSON()
			default:
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

func runSequential() {
	fmt.Println("running scraper in sequential mode...")
	allEvents := make(map[string]EventList)

	for key, venue := range allVenues {
		fmt.Printf("Scraping %s...\n", venue.Name)

		var events EventList
		switch key {
		case "turbo-haus":
			events = scrapeTurboHausJSON()
		case "bar-le-ritz":
			events = scrapeBarLeRitzJSON()
		case "mtelus":
			events = scrapeMTelusJSON()
		default:
			events = scrapeVenue(key, venue)
		}
		allEvents[key] = events
	}
	saveAllEvents(allEvents)
}
