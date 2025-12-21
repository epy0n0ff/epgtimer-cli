package formatters

import (
	"fmt"
	"strings"

	"github.com/epy0n0ff/epgtimer-cli/internal/models"
)

// TableFormatter formats rules as a human-readable table
type TableFormatter struct{}

// Format converts rules to table format string
func (t *TableFormatter) Format(rules []models.AutoAddRule) (string, error) {
	if len(rules) == 0 {
		return "No automatic recording rules found.\n", nil
	}

	var sb strings.Builder

	// Header row
	sb.WriteString(fmt.Sprintf("%-4s  %-8s  %-30s  %-30s  %s\n",
		"ID", "Enabled", "Keywords", "Exclusions", "Channels"))

	// Data rows
	for _, rule := range rules {
		id := fmt.Sprintf("%d", rule.ID)
		enabled := "Yes"
		if !rule.SearchSettings.IsEnabled() {
			enabled = "No"
		}

		// Truncate keywords and exclusions to 30 characters
		keywords := truncateString(rule.SearchSettings.AndKey, 30)
		exclusions := truncateString(rule.SearchSettings.NotKey, 30)

		// Format channel count
		channelCount := fmt.Sprintf("%d channels", rule.SearchSettings.ChannelCount())

		sb.WriteString(fmt.Sprintf("%-4s  %-8s  %-30s  %-30s  %s\n",
			id, enabled, keywords, exclusions, channelCount))
	}

	return sb.String(), nil
}

// truncateString truncates a string to maxLen characters, adding "..." if truncated
func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}
