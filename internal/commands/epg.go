package commands

import (
	"fmt"
	"os"
	"strings"

	"github.com/epy0n0ff/epgtimer-cli/internal/client"
	"github.com/epy0n0ff/epgtimer-cli/internal/formatters"
	"github.com/epy0n0ff/epgtimer-cli/internal/models"
	"github.com/spf13/cobra"
)

var epgCmd = &cobra.Command{
	Use:   "epg",
	Short: "List EPG (Electronic Program Guide) for channels",
	Long: `Retrieve and display EPG (Electronic Program Guide) data from EpgTimer's EMWUI service.

The epg command retrieves program schedule information for specified channels
and displays them in a human-readable table format. Each program shows its date,
time, duration, title, and genre.

Channel Selection:
  --channel       Specific channel in ONID-TSID-SID format
  --all-channels  Retrieve EPG for all channels from serviceList_without_local.txt

Output Formats:
  table  - Human-readable table format (default)
           Best for: Quick viewing in terminal
  json   - JSON format with full event structure
           Best for: Programmatic processing, piping to jq
  csv    - Comma-separated values
           Best for: Excel/spreadsheet analysis, data import
  tsv    - Tab-separated values
           Best for: Data exchange, simple parsing

Examples:
  # List EPG for a specific channel
  epgtimer epg --channel "32736-32736-1024"

  # List EPG for all channels in serviceList_without_local.txt
  epgtimer epg --all-channels

  # Filter by title
  epgtimer epg --channel "32736-32736-1024" --title "ニュース"

  # Filter by genre
  epgtimer epg --channel "32736-32736-1024" --genre "ドラマ"

  # Export to JSON file
  epgtimer epg --channel "32736-32736-1024" --format json --output epg.json

  # Export all channels to CSV
  epgtimer epg --all-channels --format csv -o epg.csv
`,
	RunE: runEPG,
}

func init() {
	// Channel selection flags
	epgCmd.Flags().String("channel", "", "Channel ID in ONID-TSID-SID format (e.g., 32736-32736-1024)")
	epgCmd.Flags().Bool("all-channels", false, "Retrieve EPG for all channels from serviceList_without_local.txt")

	// Filter flags
	epgCmd.Flags().String("title", "", "Filter by program title (substring match, case-insensitive)")
	epgCmd.Flags().String("genre", "", "Filter by genre (substring match, case-insensitive)")

	// Export flags
	epgCmd.Flags().String("format", "table", "Output format: table, json, csv, tsv")
	epgCmd.Flags().StringP("output", "o", "", "Output file path (default: stdout)")
}

func runEPG(cmd *cobra.Command, args []string) error {
	// Get EMWUI endpoint
	endpoint, err := cmd.Flags().GetString("endpoint")
	if err != nil {
		return fmt.Errorf("failed to get endpoint flag: %w", err)
	}

	if endpoint == "" {
		endpoint = os.Getenv("EMWUI_ENDPOINT")
	}

	if endpoint == "" {
		return fmt.Errorf("EMWUI endpoint not configured\n\nPlease set the endpoint using:\n  1. --endpoint flag: epgtimer epg --endpoint http://192.168.1.10:5510\n  2. EMWUI_ENDPOINT environment variable: export EMWUI_ENDPOINT=http://192.168.1.10:5510")
	}

	// Get channel selection flags
	channel, _ := cmd.Flags().GetString("channel")
	allChannels, _ := cmd.Flags().GetBool("all-channels")

	// Validate channel selection
	if !allChannels && channel == "" {
		return fmt.Errorf("must specify either --channel or --all-channels")
	}

	if allChannels && channel != "" {
		return fmt.Errorf("cannot use both --channel and --all-channels")
	}

	// Create API client
	apiClient := client.NewClient(endpoint)

	var allEvents []models.EventInfo

	if allChannels {
		// Read channel list from serviceList_without_local.txt
		channels, err := readChannelList("serviceList_without_local.txt")
		if err != nil {
			return fmt.Errorf("failed to read channel list: %w", err)
		}

		fmt.Printf("Retrieving EPG for %d channels...\n", len(channels))

		// Retrieve EPG for each channel
		for i, ch := range channels {
			fmt.Printf("  [%d/%d] Retrieving %s...\r", i+1, len(channels), ch.String())

			response, err := apiClient.EnumEventInfo(ch.ONID, ch.TSID, ch.SID)
			if err != nil {
				fmt.Printf("\n  Warning: Failed to retrieve EPG for %s: %v\n", ch.String(), err)
				continue
			}

			allEvents = append(allEvents, response.Items...)
		}

		fmt.Printf("\nRetrieved %d programs from %d channels\n\n", len(allEvents), len(channels))
	} else {
		// Parse single channel
		ch, err := models.ParseServiceListEntry(channel)
		if err != nil {
			return fmt.Errorf("invalid channel format: %w", err)
		}

		// Retrieve EPG for single channel
		response, err := apiClient.EnumEventInfo(ch.ONID, ch.TSID, ch.SID)
		if err != nil {
			return formatConnectionError(err, endpoint)
		}

		allEvents = response.Items
	}

	// Apply filters
	filteredEvents := applyEPGFilters(cmd, allEvents)

	// Handle empty results
	if len(filteredEvents) == 0 {
		fmt.Println("No programs match the specified filters.")
		return nil
	}

	// Get format flag
	format, _ := cmd.Flags().GetString("format")

	// Format events
	var output string
	switch format {
	case "table":
		formatter := &formatters.EPGTableFormatter{}
		output, err = formatter.Format(filteredEvents)
	case "json":
		output = formatEPGAsJSON(filteredEvents)
	case "csv":
		output = formatEPGAsCSV(filteredEvents)
	case "tsv":
		output = formatEPGAsTSV(filteredEvents)
	default:
		return fmt.Errorf("unsupported format '%s'. Supported formats: table, json, csv, tsv", format)
	}

	if err != nil {
		return fmt.Errorf("failed to format output: %w", err)
	}

	// Get output flag
	outputPath, _ := cmd.Flags().GetString("output")

	// Write to file or stdout
	if outputPath != "" {
		if err := os.WriteFile(outputPath, []byte(output), 0644); err != nil {
			return fmt.Errorf("failed to write output to file '%s': %w", outputPath, err)
		}
		fmt.Printf("Successfully exported %d programs to %s\n", len(filteredEvents), outputPath)
	} else {
		fmt.Print(output)
	}

	return nil
}

// readChannelList reads channel IDs from serviceList_without_local.txt
func readChannelList(filename string) ([]*models.ServiceListEntry, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read file '%s': %w", filename, err)
	}

	lines := strings.Split(string(data), "\n")
	var channels []*models.ServiceListEntry

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		ch, err := models.ParseServiceListEntry(line)
		if err != nil {
			return nil, fmt.Errorf("invalid channel format '%s': %w", line, err)
		}

		channels = append(channels, ch)
	}

	return channels, nil
}

func applyEPGFilters(cmd *cobra.Command, events []models.EventInfo) []models.EventInfo {
	title, _ := cmd.Flags().GetString("title")
	genre, _ := cmd.Flags().GetString("genre")

	var filtered []models.EventInfo
	for _, event := range events {
		// Title filter
		if title != "" {
			titleLower := strings.ToLower(event.EventName)
			filterLower := strings.ToLower(title)
			if !strings.Contains(titleLower, filterLower) {
				continue
			}
		}

		// Genre filter
		if genre != "" {
			genreLower := strings.ToLower(event.GenreString())
			filterLower := strings.ToLower(genre)
			if !strings.Contains(genreLower, filterLower) {
				continue
			}
		}

		filtered = append(filtered, event)
	}

	return filtered
}

// Temporary JSON formatter for EPG
func formatEPGAsJSON(events []models.EventInfo) string {
	var sb strings.Builder
	sb.WriteString("[\n")
	for i, event := range events {
		sb.WriteString(fmt.Sprintf(`  {
    "channel_id": "%s",
    "onid": %d,
    "tsid": %d,
    "sid": %d,
    "event_id": %d,
    "service_name": "%s",
    "start_date": "%s",
    "start_time": "%s",
    "duration_minutes": %d,
    "event_name": "%s",
    "event_text": "%s",
    "genre": "%s"
  }`,
			event.ChannelID(),
			event.ONID,
			event.TSID,
			event.SID,
			event.EventID,
			escapeJSON(event.ServiceName),
			event.StartDate,
			event.StartTime,
			event.DurationMinutes(),
			escapeJSON(event.EventName),
			escapeJSON(event.EventText),
			escapeJSON(event.GenreString()),
		))
		if i < len(events)-1 {
			sb.WriteString(",\n")
		} else {
			sb.WriteString("\n")
		}
	}
	sb.WriteString("]\n")
	return sb.String()
}

// Temporary CSV formatter for EPG
func formatEPGAsCSV(events []models.EventInfo) string {
	var sb strings.Builder
	sb.WriteString("ChannelID,ONID,TSID,SID,EventID,ServiceName,StartDate,StartTime,DurationMinutes,EventName,EventText,Genre\n")
	for _, event := range events {
		sb.WriteString(fmt.Sprintf("%s,%d,%d,%d,%d,%s,%s,%s,%d,%s,%s,%s\n",
			event.ChannelID(),
			event.ONID,
			event.TSID,
			event.SID,
			event.EventID,
			escapeCSV(event.ServiceName),
			event.StartDate,
			event.StartTime,
			event.DurationMinutes(),
			escapeCSV(event.EventName),
			escapeCSV(event.EventText),
			escapeCSV(event.GenreString()),
		))
	}
	return sb.String()
}

// Temporary TSV formatter for EPG
func formatEPGAsTSV(events []models.EventInfo) string {
	var sb strings.Builder
	sb.WriteString("ChannelID\tONID\tTSID\tSID\tEventID\tServiceName\tStartDate\tStartTime\tDurationMinutes\tEventName\tEventText\tGenre\n")
	for _, event := range events {
		sb.WriteString(fmt.Sprintf("%s\t%d\t%d\t%d\t%d\t%s\t%s\t%s\t%d\t%s\t%s\t%s\n",
			event.ChannelID(),
			event.ONID,
			event.TSID,
			event.SID,
			event.EventID,
			escapeTSV(event.ServiceName),
			event.StartDate,
			event.StartTime,
			event.DurationMinutes(),
			escapeTSV(event.EventName),
			escapeTSV(event.EventText),
			escapeTSV(event.GenreString()),
		))
	}
	return sb.String()
}
