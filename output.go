package main

import (
	"fmt"
	"os"
	"strings"
)

func saveAllEvents(allEvents map[string]EventList) {

	if err := os.MkdirAll("right_now", 0755); err != nil {
		fmt.Printf("Failed to create right_now directory: %v\n", err)
		return
	}
	if err := os.MkdirAll("tonight", 0755); err != nil {
		fmt.Printf("Failed to create tonight directory: %v\n", err)
		return
	}
	if err := os.MkdirAll("tomorrow", 0755); err != nil {
		fmt.Printf("Failed to create tomorrow directory: %v\n", err)
		return
	}
	if err := os.MkdirAll("this_week", 0755); err != nil {
		fmt.Printf("Failed to create this_week directory: %v\n", err)
		return
	}
	if err := os.MkdirAll("this_weekend", 0755); err != nil {
		fmt.Printf("Failed to create this_weekend directory: %v\n", err)
		return
	}
	if err := os.MkdirAll("output", 0755); err != nil {
		fmt.Printf("Failed to create output directory: %v\n", err)
		return
	}

	// Create group subdirectories
	for _, venue := range allVenues {
		if venue.Group != "" {
			if err := os.MkdirAll(fmt.Sprintf("output/%s", venue.Group), 0755); err != nil {
				fmt.Printf("Failed to create group directory %s: %v\n", venue.Group, err)
			}
		}
	}

	var allEventsList EventList
	for venueKey, event := range allEvents {
		venue := allVenues[venueKey]

		// Use group subdir if venue belongs to a group (rare)
		path := "output"
		if venue.Group != "" {
			path = fmt.Sprintf("output/%s", venue.Group)
		}
		saveAllEventsToTextFile(event, fmt.Sprintf("%s/%s.txt", path, venueKey), venue.Name)
		saveVenueEventsJSON(event, venueKey, venue.Name, path)

		allEventsList = append(allEventsList, event...)
	}

	allEventsList.SortByDate()

	// all Jason output
	saveAllEventsJSON(allEventsList)
	saveRightNowJSON(allEventsList)
	saveTonightJSON(allEventsList)
	saveTomorrowJSON(allEventsList)
	saveThisWeekJSON(allEventsList)
	saveThisWeekendJSON(allEventsList)

	// all filtered output
	saveHappeningRightNowEvents(allEventsList.RightNow(), "right_now/right_now.txt")
	saveTonightEvents(allEventsList.Tonight(), "tonight/tonight.txt")
	saveTomorrowEvents(allEventsList.Tomorrow(), "tomorrow/tomorrow.txt")
	saveThisWeekEvents(allEventsList.ThisWeek(), "this_week/this_week.txt")
	saveThisWeekendEvents(allEventsList.ThisWeekend(), "this_weekend/this_weekend.txt")
}

func saveAllEventsToTextFile(events EventList, filename, venueName string) error {
	var sb strings.Builder
	sb.Grow(len(events) * 512)

	sb.WriteString(fmt.Sprintf("%s Events\n", venueName))
	sb.WriteString(strings.Repeat("=", 90) + "\n\n")

	if len(events) == 0 {
		sb.WriteString("No upcoming events.\n")
	}

	for i, e := range events {
		sb.WriteString(fmt.Sprintf("Event #%d\n", i+1))
		sb.WriteString(fmt.Sprintf("Name:      %s\n", e.Name))
		sb.WriteString(fmt.Sprintf("Venue:     %s\n", e.Venue))
		sb.WriteString(fmt.Sprintf("Address:   %s\n", e.Address))
		sb.WriteString(fmt.Sprintf("Date:      %s\n", e.Date))
		if e.Time != "" {
			sb.WriteString(fmt.Sprintf("Time:      %s\n", e.Time))
		} else {
			sb.WriteString(fmt.Sprintln("Time:      not available"))
		}
		if e.Price != "" {
			sb.WriteString(fmt.Sprintf("Price:     %s\n", e.Price))
		} else {
			sb.WriteString(fmt.Sprintln("Price:     not available."))
		}
		if e.TicketURL != "" {
			sb.WriteString(fmt.Sprintf("Ticket Link:  %s\n\n", e.TicketURL))
		} else {
			sb.WriteString(fmt.Sprintln("Ticket Link:  not available\n\n"))
		}

		// ---------------------------- detailed info ------------------------------
		sb.WriteString(fmt.Sprintf("Parsed Date:       %s\n", e.ParsedDate))
		sb.WriteString(fmt.Sprintf("Day of week:       %s\n", e.DayOfWeek))
		if e.Price != "" {
			sb.WriteString(fmt.Sprintf("Price value:       %.2f\n", e.PriceValue))
			sb.WriteString(fmt.Sprintf("is free:           %t\n", e.IsFree))
		}
		sb.WriteString(fmt.Sprintf("Days Until:        %d\n", e.DaysUntil))
		sb.WriteString(fmt.Sprintf("is today:          %t\n", e.IsToday))
		sb.WriteString(fmt.Sprintf("is this week:      %t\n", e.IsThisWeek))
		sb.WriteString(fmt.Sprintf("is this weekend:   %t\n\n", e.IsThisWeekend))
		sb.WriteString(strings.Repeat("-", 90) + "\n\n")
	}

	return os.WriteFile(filename, []byte(sb.String()), 0644)
}

// -------------------------------------FILTERED---------------------------------------------------

func saveHappeningRightNowEvents(events EventList, filename string) error {
	var sb strings.Builder
	sb.Grow(len(events) * 512)

	sb.WriteString(fmt.Sprintln("# Shows & Events RIGHT NOW"))
	sb.WriteString(fmt.Sprintln("----------------------------------------------\n"))

	count := 0
	for _, e := range events {
		count++
		sb.WriteString(fmt.Sprintf("Event #%d\n", count))
		sb.WriteString(fmt.Sprintf("Name:      %s\n", e.Name))
		sb.WriteString(fmt.Sprintf("Venue:     %s\n", e.Venue))
		sb.WriteString(fmt.Sprintf("Date:      %s\n", e.Date))
		sb.WriteString(fmt.Sprintf("Address:   %s\n", e.Address))
		if e.Time != "" {
			sb.WriteString(fmt.Sprintf("Time:      %s\n", e.Time))
		} else {
			sb.WriteString(fmt.Sprintln("Time:      not available"))
		}
		if e.Price != "" {
			sb.WriteString(fmt.Sprintf("Price:     %s\n", e.Price))
		} else {
			sb.WriteString(fmt.Sprintln("Price:     not available"))
		}
		if e.TicketURL != "" {
			sb.WriteString(fmt.Sprintf("Ticket Link:  %s\n\n", e.TicketURL))
		} else {
			sb.WriteString(fmt.Sprintln("Ticket Link:  not available\n\n"))
		}
		sb.WriteString(strings.Repeat("-", 90) + "\n\n")

	}

	return os.WriteFile(filename, []byte(sb.String()), 0644)
}

func saveTonightEvents(events EventList, filename string) error {
	var sb strings.Builder
	sb.Grow(len(events) * 512)

	sb.WriteString(fmt.Sprintln("# Shows & Events tonight"))
	sb.WriteString(fmt.Sprintln("----------------------------------------------\n"))

	count := 0
	for _, e := range events {
		count++
		sb.WriteString(fmt.Sprintf("Event #%d\n", count))
		sb.WriteString(fmt.Sprintf("Name:      %s\n", e.Name))
		sb.WriteString(fmt.Sprintf("Venue:     %s\n", e.Venue))
		sb.WriteString(fmt.Sprintf("Date:      %s\n", e.Date))
		sb.WriteString(fmt.Sprintf("Address:   %s\n", e.Address))
		if e.Time != "" {
			sb.WriteString(fmt.Sprintf("Time:      %s\n", e.Time))
		} else {
			sb.WriteString(fmt.Sprintln("Time:      not available"))
		}
		if e.Price != "" {
			sb.WriteString(fmt.Sprintf("Price:     %s\n", e.Price))
		} else {
			sb.WriteString(fmt.Sprintln("Price:     not available"))
		}
		if e.TicketURL != "" {
			sb.WriteString(fmt.Sprintf("Ticket Link:  %s\n\n", e.TicketURL))
		} else {
			sb.WriteString(fmt.Sprintln("Ticket Link:  not available\n\n"))
		}
		sb.WriteString(strings.Repeat("-", 90) + "\n\n")

	}

	return os.WriteFile(filename, []byte(sb.String()), 0644)
}

func saveTomorrowEvents(events EventList, filename string) error {
	var sb strings.Builder
	sb.Grow(len(events) * 512)

	sb.WriteString(fmt.Sprintln("# Shows & Events tomorrow"))
	sb.WriteString(fmt.Sprintln("----------------------------------------------\n"))

	count := 0
	for _, e := range events {
		count++
		sb.WriteString(fmt.Sprintf("Event #%d\n", count))
		sb.WriteString(fmt.Sprintf("Name:      %s\n", e.Name))
		sb.WriteString(fmt.Sprintf("Venue:     %s\n", e.Venue))
		sb.WriteString(fmt.Sprintf("Date:      %s\n", e.Date))
		sb.WriteString(fmt.Sprintf("Address:   %s\n", e.Address))
		if e.Time != "" {
			sb.WriteString(fmt.Sprintf("Time:      %s\n", e.Time))
		} else {
			sb.WriteString(fmt.Sprintln("Time:      not available"))
		}
		if e.Price != "" {
			sb.WriteString(fmt.Sprintf("Price:     %s\n", e.Price))
		} else {
			sb.WriteString(fmt.Sprintln("Price:     not available"))
		}
		if e.TicketURL != "" {
			sb.WriteString(fmt.Sprintf("Ticket Link:  %s\n\n", e.TicketURL))
		} else {
			sb.WriteString(fmt.Sprintln("Ticket Link:  not available\n\n"))
		}
		sb.WriteString(strings.Repeat("-", 90) + "\n\n")

	}

	return os.WriteFile(filename, []byte(sb.String()), 0644)
}

func saveThisWeekEvents(events EventList, filename string) error {
	var sb strings.Builder
	sb.Grow(len(events) * 512)

	sb.WriteString(fmt.Sprintln("# Shows & Events this week"))
	sb.WriteString(fmt.Sprintln("----------------------------------------------\n"))

	count := 0
	for _, e := range events {
		count++
		sb.WriteString(fmt.Sprintf("Event #%d\n", count))
		sb.WriteString(fmt.Sprintf("Name:      %s\n", e.Name))
		sb.WriteString(fmt.Sprintf("Venue:     %s\n", e.Venue))
		sb.WriteString(fmt.Sprintf("Date:      %s\n", e.Date))
		sb.WriteString(fmt.Sprintf("Address:   %s\n", e.Address))
		if e.Time != "" {
			sb.WriteString(fmt.Sprintf("Time:      %s\n", e.Time))
		} else {
			sb.WriteString(fmt.Sprintln("Time:      not available"))
		}
		if e.Price != "" {
			sb.WriteString(fmt.Sprintf("Price:     %s\n", e.Price))
		} else {
			sb.WriteString(fmt.Sprintln("Price:     not available"))
		}
		if e.TicketURL != "" {
			sb.WriteString(fmt.Sprintf("Ticket Link:  %s\n\n", e.TicketURL))
		} else {
			sb.WriteString(fmt.Sprintln("Ticket Link:  not available\n\n"))
		}
		sb.WriteString(strings.Repeat("-", 90) + "\n\n")

	}

	return os.WriteFile(filename, []byte(sb.String()), 0644)
}

func saveThisWeekendEvents(events EventList, filename string) error {
	var sb strings.Builder
	sb.Grow(len(events) * 512)

	sb.WriteString(fmt.Sprintln("# Shows & Events this weekend"))
	sb.WriteString(fmt.Sprintln("----------------------------------------------\n"))

	count := 0
	for _, e := range events {
		count++
		sb.WriteString(fmt.Sprintf("Event #%d\n", count))
		sb.WriteString(fmt.Sprintf("Name:      %s\n", e.Name))
		sb.WriteString(fmt.Sprintf("Venue:     %s\n", e.Venue))
		sb.WriteString(fmt.Sprintf("Date:      %s\n", e.Date))
		sb.WriteString(fmt.Sprintf("Address:   %s\n", e.Address))
		if e.Time != "" {
			sb.WriteString(fmt.Sprintf("Time:      %s\n", e.Time))
		} else {
			sb.WriteString(fmt.Sprintln("Time:      not available"))
		}
		if e.Price != "" {
			sb.WriteString(fmt.Sprintf("Price:     %s\n", e.Price))
		} else {
			sb.WriteString(fmt.Sprintln("Price:     not available"))
		}
		if e.TicketURL != "" {
			sb.WriteString(fmt.Sprintf("Ticket Link:  %s\n\n", e.TicketURL))
		} else {
			sb.WriteString(fmt.Sprintln("Ticket Link:  not available\n\n"))
		}
		sb.WriteString(strings.Repeat("-", 90) + "\n\n")
	}

	return os.WriteFile(filename, []byte(sb.String()), 0644)
}
