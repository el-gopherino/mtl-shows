package main

import (
	"fmt"
	"os"
	"strings"
)

func saveTonightEvents(events []Event, filename string) error {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintln("# Shows & Events tonight"))
	sb.WriteString(fmt.Sprintln("----------------------------------------------\n"))

	count := 0
	for _, e := range events {
		if isToday(e.ParsedDate) {
			count++
			sb.WriteString(fmt.Sprintf("Event #%d\n", count))
			sb.WriteString(fmt.Sprintf("Name:      %s\n", e.Name))
			sb.WriteString(fmt.Sprintf("Venue:     %s\n", e.Venue))
			sb.WriteString(fmt.Sprintf("Date:      %s\n", e.Date))
			sb.WriteString(fmt.Sprintf("Address:   %s\n", e.Address))
			sb.WriteString(fmt.Sprintf("Time:      %s\n", e.Time))
			if e.Price != "" {
				sb.WriteString(fmt.Sprintf("Price:     %s\n", e.Price))
			} else {
				sb.WriteString(fmt.Sprintln("Price:     not available"))
			}
			sb.WriteString(fmt.Sprintf("Ticket Link:  %s\n\n", e.TicketURL))
			sb.WriteString(strings.Repeat("-", 90) + "\n\n")
		}
	}

	return os.WriteFile(filename, []byte(sb.String()), 0644)
}

func saveThisWeekEvents(events []Event, filename string) error {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintln("# Shows & Events this week"))
	sb.WriteString(fmt.Sprintln("----------------------------------------------\n"))

	count := 0
	for _, e := range events {
		if isThisWeek(e.ParsedDate) {
			count++
			sb.WriteString(fmt.Sprintf("Event #%d\n", count))
			sb.WriteString(fmt.Sprintf("Name:      %s\n", e.Name))
			sb.WriteString(fmt.Sprintf("Venue:     %s\n", e.Venue))
			sb.WriteString(fmt.Sprintf("Date:      %s\n", e.Date))
			sb.WriteString(fmt.Sprintf("Address:   %s\n", e.Address))
			sb.WriteString(fmt.Sprintf("Time:      %s\n", e.Time))
			if e.Price != "" {
				sb.WriteString(fmt.Sprintf("Price:     %s\n", e.Price))
			} else {
				sb.WriteString(fmt.Sprintln("Price:     not available"))
			}
			sb.WriteString(fmt.Sprintf("Ticket Link:  %s\n\n", e.TicketURL))
			sb.WriteString(strings.Repeat("-", 90) + "\n\n")
		}
	}

	return os.WriteFile(filename, []byte(sb.String()), 0644)
}
