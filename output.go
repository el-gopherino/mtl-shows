package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

func saveAllEvents(allEvents map[string][]Event) {

	if err := os.MkdirAll("tonight", 0755); err != nil {
		fmt.Printf("Failed to create tonight directory: %v\n", err)
		return
	}
	if err := os.MkdirAll("this_week", 0755); err != nil {
		fmt.Printf("Failed to create tonight directory: %v\n", err)
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

	// events slice for all filtered events (tonight, thisWeek, etc.)
	var allEventsList []Event
	for venueKey, events := range allEvents {
		venue := allVenues[venueKey]

		// Use group subdir if venue belongs to a group (rare)
		path := "output"
		if venue.Group != "" {
			path = fmt.Sprintf("output/%s", venue.Group)
		}

		saveAllEventsToTextFile(events, fmt.Sprintf("%s/%s_events.txt", path, venueKey), venue.Name)
		saveAllEventsToMarkdown(events, fmt.Sprintf("%s/%s_events.md", path, venueKey), venue.Name)

		allEventsList = append(allEventsList, events...)
	}

	saveTonightEvents(allEventsList, "tonight/tonight_events.txt")
	saveThisWeekEvents(allEventsList, "this_week/this_week_events.txt")
}

func saveAllEventsToTextFile(events []Event, filename, venueName string) error {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("%s Events\n", venueName))
	sb.WriteString(strings.Repeat("=", 90) + "\n\n")

	if len(events) == 0 {
		sb.WriteString("No upcoming events.\n")
	}

	for i, e := range events {

		sb.WriteString(fmt.Sprintf("Event #%d\n", i+1))
		sb.WriteString(fmt.Sprintf("Name:      %s\n", e.Name))
		sb.WriteString(fmt.Sprintf("Venue:     %s\n", e.Venue))
		sb.WriteString(fmt.Sprintf("Date:      %s\n", e.Date))
		sb.WriteString(fmt.Sprintf("Address:   %s\n", e.Address))
		sb.WriteString(fmt.Sprintf("Time:      %s\n", e.Time))
		sb.WriteString(fmt.Sprintf("Price:     %s\n", e.Price))
		sb.WriteString(fmt.Sprintf("Tickets :  %s\n", e.TicketURL))

		sb.WriteString(fmt.Sprintf("\nMore detailed info:\n"))

		sb.WriteString(fmt.Sprintf("Parsed Date:       %s\n", e.ParsedDate)) // TODO: fix pour quai des brumes + voir si heure local EST est possible
		sb.WriteString(fmt.Sprintf("Day of week:       %s\n", e.DayOfWeek))  // ✅
		sb.WriteString(fmt.Sprintf("\n"))
		sb.WriteString(fmt.Sprintf("Price value:       %.2f\n", e.PriceValue)) // ✅
		sb.WriteString(fmt.Sprintf("is free:           %t\n", e.IsFree))       // ✅
		sb.WriteString(fmt.Sprintf("\n"))
		sb.WriteString(fmt.Sprintf("Days Until:        %d\n", e.DaysUntil))       // ✅
		sb.WriteString(fmt.Sprintf("is today:          %t\n", e.IsToday))         // ✅
		sb.WriteString(fmt.Sprintf("is this week:      %t\n", e.IsThisWeek))      // ✅
		sb.WriteString(fmt.Sprintf("is this weekend:   %t\n\n", e.IsThisWeekend)) // ✅

		sb.WriteString(strings.Repeat("-", 90) + "\n\n")
	}

	return os.WriteFile(filename, []byte(sb.String()), 0644)
}

func saveAllEventsToMarkdown(events []Event, filename, venueName string) error {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("# %s\n\n", venueName))
	sb.WriteString(fmt.Sprintf("*%d upcoming events*\n\n", len(events)))
	sb.WriteString("---\n\n")

	for _, e := range events {
		sb.WriteString(fmt.Sprintf("## %s\n\n", e.Name))
		sb.WriteString(fmt.Sprintf("- **Date:** %s @ %s\n", e.Date, e.Time))
		if e.Venue != "" {
			sb.WriteString(fmt.Sprintf("- **Venue:** %s\n", e.Venue))
		}
		if e.Address != "" {
			sb.WriteString(fmt.Sprintf("- **Address:** %s\n", e.Address))
		}
		if e.Price != "" {
			sb.WriteString(fmt.Sprintf("- **Price:** %s\n", e.Price))
		}
		if e.TicketURL != "" {
			sb.WriteString(fmt.Sprintf("- **Tickets:** [Link here](%s)\n", e.TicketURL))
		}
		sb.WriteString("\n---\n\n")
	}

	return os.WriteFile(filename, []byte(sb.String()), 0644)
}

func saveAllEventsToJson(events []Event, fileName string) error {
	data, err := json.MarshalIndent(events, "", "\t")
	if err != nil {
		return fmt.Errorf("error parsing venues to JSON: %w", err)
	}
	return os.WriteFile(fileName, data, 0644)
}

//func saveAllEventsToHtml(events []Event, fileName string) error {
//
//}
