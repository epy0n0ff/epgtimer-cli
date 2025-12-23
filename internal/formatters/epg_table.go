package formatters

import (
	"fmt"
	"strings"

	"github.com/epy0n0ff/epgtimer-cli/internal/models"
)

// EPGTableFormatter formats EPG events as a human-readable table
type EPGTableFormatter struct{}

// Format converts EPG events to table format
func (t *EPGTableFormatter) Format(events []models.EventInfo) (string, error) {
	if len(events) == 0 {
		return "No events found.\n", nil
	}

	var output strings.Builder

	// Header
	output.WriteString(fmt.Sprintf("%-12s %-6s %-5s %-50s %-20s\n",
		"Date", "Time", "Mins", "Title", "Genre"))
	output.WriteString(strings.Repeat("-", 100) + "\n")

	// Data rows
	for _, event := range events {
		// Format date: 2025/12/22 -> 12/22
		dateParts := strings.Split(event.StartDate, "/")
		shortDate := fmt.Sprintf("%s/%s", dateParts[1], dateParts[2])

		// Format time: 22:30:00 -> 22:30
		timeParts := strings.Split(event.StartTime, ":")
		shortTime := fmt.Sprintf("%s:%s", timeParts[0], timeParts[1])

		title := truncate(event.EventName, 50)
		genre := truncate(event.GenreString(), 20)

		output.WriteString(fmt.Sprintf("%-12s %-6s %-5d %-50s %-20s\n",
			shortDate, shortTime, event.DurationMinutes(), title, genre))
	}

	output.WriteString(fmt.Sprintf("\nTotal: %d programs\n", len(events)))

	return output.String(), nil
}
