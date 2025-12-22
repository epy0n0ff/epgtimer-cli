package formatters

import (
	"fmt"
	"strings"

	"github.com/epy0n0ff/epgtimer-cli/internal/models"
)

// RecordingsTableFormatter formats recordings as a human-readable table
type RecordingsTableFormatter struct{}

// Format converts recordings to table format
func (t *RecordingsTableFormatter) Format(recordings []models.RecordingInfo) (string, error) {
	if len(recordings) == 0 {
		return "No recordings found.\n", nil
	}

	var output strings.Builder

	// Header
	output.WriteString(fmt.Sprintf("%-6s %-12s %-6s %-50s %-20s\n",
		"ID", "Date", "Time", "Title", "Station"))
	output.WriteString(strings.Repeat("-", 100) + "\n")

	// Data rows
	for _, rec := range recordings {
		// Format date: 2025/12/22 -> 12/22
		dateParts := strings.Split(rec.StartDate, "/")
		shortDate := fmt.Sprintf("%s/%s", dateParts[1], dateParts[2])

		// Format time: 22:30:00 -> 22:30
		timeParts := strings.Split(rec.StartTime, ":")
		shortTime := fmt.Sprintf("%s:%s", timeParts[0], timeParts[1])

		title := truncate(rec.Title, 50)
		station := truncate(rec.StationName, 20)

		output.WriteString(fmt.Sprintf("%-6d %-12s %-6s %-50s %-20s\n",
			rec.ID, shortDate, shortTime, title, station))
	}

	output.WriteString(fmt.Sprintf("\nTotal: %d recordings\n", len(recordings)))

	return output.String(), nil
}
