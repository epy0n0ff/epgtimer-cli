package formatters

import (
	"fmt"
	"strings"

	"github.com/epy0n0ff/epgtimer-cli/internal/models"
)

// ReservationsTableFormatter formats reservations as a human-readable table
type ReservationsTableFormatter struct{}

// Format converts reservations to table format
func (t *ReservationsTableFormatter) Format(reservations []models.ReservationInfo) (string, error) {
	if len(reservations) == 0 {
		return "No reservations found.\n", nil
	}

	var output strings.Builder

	// Header
	output.WriteString(fmt.Sprintf("%-6s %-12s %-6s %-50s %-20s\n",
		"ID", "Date", "Time", "Title", "Station"))
	output.WriteString(strings.Repeat("-", 100) + "\n")

	// Data rows
	for _, res := range reservations {
		// Format date: 2025/12/22 -> 12/22
		dateParts := strings.Split(res.StartDate, "/")
		shortDate := fmt.Sprintf("%s/%s", dateParts[1], dateParts[2])

		// Format time: 22:30:00 -> 22:30
		timeParts := strings.Split(res.StartTime, ":")
		shortTime := fmt.Sprintf("%s:%s", timeParts[0], timeParts[1])

		title := truncate(res.Title, 50)
		station := truncate(res.StationName, 20)

		output.WriteString(fmt.Sprintf("%-6d %-12s %-6s %-50s %-20s\n",
			res.ID, shortDate, shortTime, title, station))
	}

	output.WriteString(fmt.Sprintf("\nTotal: %d reservations\n", len(reservations)))

	return output.String(), nil
}
