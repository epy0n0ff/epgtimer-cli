package formatters

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"strconv"
	"strings"

	"github.com/epy0n0ff/epgtimer-cli/internal/models"
)

// TSVFormatter formats rules as TSV (Tab-Separated Values)
type TSVFormatter struct{}

// Format converts rules to TSV format
func (t *TSVFormatter) Format(rules []models.AutoAddRule) (string, error) {
	var buf bytes.Buffer
	writer := csv.NewWriter(&buf)
	writer.Comma = '\t' // Use tab as delimiter

	// Write header
	header := []string{
		"ID",
		"Enabled",
		"AndKey",
		"NotKey",
		"RegExp",
		"Channels",
		"ChannelCount",
		"Priority",
		"RecMode",
	}
	if err := writer.Write(header); err != nil {
		return "", fmt.Errorf("failed to write TSV header: %w", err)
	}

	// Write data rows
	for _, rule := range rules {
		// Build channel list
		channelList := make([]string, 0, len(rule.SearchSettings.ServiceList))
		for _, ch := range rule.SearchSettings.ServiceList {
			channelList = append(channelList, ch.String())
		}
		channelsStr := strings.Join(channelList, ";")

		row := []string{
			strconv.Itoa(rule.ID),
			strconv.FormatBool(rule.SearchSettings.IsEnabled()),
			rule.SearchSettings.AndKey,
			rule.SearchSettings.NotKey,
			strconv.FormatBool(rule.SearchSettings.IsRegex()),
			channelsStr,
			strconv.Itoa(rule.SearchSettings.ChannelCount()),
			strconv.Itoa(rule.RecordingSettings.Priority),
			strconv.Itoa(rule.RecordingSettings.RecMode),
		}
		if err := writer.Write(row); err != nil {
			return "", fmt.Errorf("failed to write TSV row: %w", err)
		}
	}

	writer.Flush()
	if err := writer.Error(); err != nil {
		return "", fmt.Errorf("TSV writer error: %w", err)
	}

	return buf.String(), nil
}
