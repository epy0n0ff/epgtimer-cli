package formatters

import (
	"fmt"
	"strings"

	"github.com/epy0n0ff/epgtimer-cli/internal/models"
)

// ChannelsTableFormatter formats channels as a human-readable table
type ChannelsTableFormatter struct{}

// Format converts channels to table format
func (t *ChannelsTableFormatter) Format(channels []models.ChannelInfo) (string, error) {
	if len(channels) == 0 {
		return "No channels found.\n", nil
	}

	var output strings.Builder

	// Header
	output.WriteString(fmt.Sprintf("%-20s %-12s %-6s %-40s %-20s\n",
		"Channel ID", "Type", "Key", "Channel Name", "Network"))
	output.WriteString(strings.Repeat("-", 100) + "\n")

	// Data rows
	for _, ch := range channels {
		channelID := ch.ChannelID()
		serviceType := ch.ServiceTypeString()
		keyID := "-"
		if ch.RemoteControlKeyID > 0 {
			keyID = fmt.Sprintf("%d", ch.RemoteControlKeyID)
		}
		channelName := truncate(ch.ServiceName, 40)
		networkName := truncate(ch.NetworkName, 20)

		output.WriteString(fmt.Sprintf("%-20s %-12s %-6s %-40s %-20s\n",
			channelID, serviceType, keyID, channelName, networkName))
	}

	output.WriteString(fmt.Sprintf("\nTotal: %d channels\n", len(channels)))

	return output.String(), nil
}

// truncate limits string length and adds ellipsis if needed
func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	if maxLen <= 3 {
		return s[:maxLen]
	}
	return s[:maxLen-3] + "..."
}
