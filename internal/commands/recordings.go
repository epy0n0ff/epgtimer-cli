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

var recordingsCmd = &cobra.Command{
	Use:   "recordings",
	Short: "List recorded programs",
	Long: `Retrieve and display recorded programs from EpgTimer's EMWUI service.

The recordings command retrieves all recorded programs stored in EpgTimer
and displays them in a human-readable table format. Each recording shows its ID, date,
time, program title, and station name.

Note: The API returns recordings in paginated batches (200 items per request).
This command retrieves only the first batch by default.

Output Formats:
  table  - Human-readable table format (default)
           Best for: Quick viewing in terminal
  json   - JSON format with full recording structure
           Best for: Programmatic processing, piping to jq
  csv    - Comma-separated values
           Best for: Excel/spreadsheet analysis, data import
  tsv    - Tab-separated values
           Best for: Data exchange, simple parsing

Examples:
  # List recordings
  epgtimer recordings

  # List recordings from a specific EMWUI server
  epgtimer recordings --endpoint http://192.168.1.10:5510

  # Filter by title
  epgtimer recordings --title "ニュース"

  # Filter by station
  epgtimer recordings --station "NHK"

  # Show only protected recordings
  epgtimer recordings --protected

  # Export to JSON file
  epgtimer recordings --format json --output recordings.json

  # Export to CSV
  epgtimer recordings --format csv -o recordings.csv
`,
	RunE: runRecordings,
}

func init() {
	// Filter flags
	recordingsCmd.Flags().String("title", "", "Filter by title (substring match, case-insensitive)")
	recordingsCmd.Flags().String("station", "", "Filter by station name (substring match, case-insensitive)")
	recordingsCmd.Flags().String("channel", "", "Filter by channel ID (exact match, format: ONID-TSID-SID)")
	recordingsCmd.Flags().Bool("protected", false, "Show only protected recordings")

	// Export flags
	recordingsCmd.Flags().String("format", "table", "Output format: table, json, csv, tsv")
	recordingsCmd.Flags().StringP("output", "o", "", "Output file path (default: stdout)")
}

func runRecordings(cmd *cobra.Command, args []string) error {
	// Get EMWUI endpoint
	endpoint, err := cmd.Flags().GetString("endpoint")
	if err != nil {
		return fmt.Errorf("failed to get endpoint flag: %w", err)
	}

	if endpoint == "" {
		endpoint = os.Getenv("EMWUI_ENDPOINT")
	}

	if endpoint == "" {
		return fmt.Errorf("EMWUI endpoint not configured\n\nPlease set the endpoint using:\n  1. --endpoint flag: epgtimer recordings --endpoint http://192.168.1.10:5510\n  2. EMWUI_ENDPOINT environment variable: export EMWUI_ENDPOINT=http://192.168.1.10:5510")
	}

	// Create API client
	apiClient := client.NewClient(endpoint)

	// Retrieve recordings
	response, err := apiClient.EnumRecInfo()
	if err != nil {
		return formatConnectionError(err, endpoint)
	}

	// Apply filters
	filteredRecordings := applyRecordingFilters(cmd, response.Items)

	// Handle empty results
	if len(filteredRecordings) == 0 {
		fmt.Println("No recordings match the specified filters.")
		return nil
	}

	// Get format flag
	format, _ := cmd.Flags().GetString("format")

	// Format recordings
	var output string
	switch format {
	case "table":
		formatter := &formatters.RecordingsTableFormatter{}
		output, err = formatter.Format(filteredRecordings)
	case "json":
		output = formatRecordingsAsJSON(filteredRecordings)
	case "csv":
		output = formatRecordingsAsCSV(filteredRecordings)
	case "tsv":
		output = formatRecordingsAsTSV(filteredRecordings)
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
		fmt.Printf("Successfully exported %d recordings to %s\n", len(filteredRecordings), outputPath)
	} else {
		fmt.Print(output)
	}

	return nil
}

func applyRecordingFilters(cmd *cobra.Command, recordings []models.RecordingInfo) []models.RecordingInfo {
	title, _ := cmd.Flags().GetString("title")
	station, _ := cmd.Flags().GetString("station")
	channel, _ := cmd.Flags().GetString("channel")
	protected, _ := cmd.Flags().GetBool("protected")

	var filtered []models.RecordingInfo
	for _, rec := range recordings {
		// Title filter
		if title != "" {
			titleLower := strings.ToLower(rec.Title)
			filterLower := strings.ToLower(title)
			if !strings.Contains(titleLower, filterLower) {
				continue
			}
		}

		// Station filter
		if station != "" {
			stationLower := strings.ToLower(rec.StationName)
			filterLower := strings.ToLower(station)
			if !strings.Contains(stationLower, filterLower) {
				continue
			}
		}

		// Channel filter
		if channel != "" {
			if rec.ChannelID() != channel {
				continue
			}
		}

		// Protected filter
		if protected && !rec.IsProtected() {
			continue
		}

		filtered = append(filtered, rec)
	}

	return filtered
}

// Temporary JSON formatter for recordings
func formatRecordingsAsJSON(recordings []models.RecordingInfo) string {
	var sb strings.Builder
	sb.WriteString("[\n")
	for i, rec := range recordings {
		sb.WriteString(fmt.Sprintf(`  {
    "id": %d,
    "title": "%s",
    "start_date": "%s",
    "start_time": "%s",
    "duration_second": %d,
    "duration_minutes": %d,
    "station_name": "%s",
    "channel_id": "%s",
    "onid": %d,
    "tsid": %d,
    "sid": %d,
    "event_id": %d,
    "comment": "%s",
    "rec_file_path": "%s",
    "protected": %t
  }`,
			rec.ID,
			escapeJSON(rec.Title),
			rec.StartDate,
			rec.StartTime,
			rec.DurationSecond,
			rec.DurationMinutes(),
			escapeJSON(rec.StationName),
			rec.ChannelID(),
			rec.ONID,
			rec.TSID,
			rec.SID,
			rec.EventID,
			escapeJSON(rec.Comment),
			escapeJSON(rec.RecFilePath),
			rec.IsProtected(),
		))
		if i < len(recordings)-1 {
			sb.WriteString(",\n")
		} else {
			sb.WriteString("\n")
		}
	}
	sb.WriteString("]\n")
	return sb.String()
}

// Temporary CSV formatter for recordings
func formatRecordingsAsCSV(recordings []models.RecordingInfo) string {
	var sb strings.Builder
	sb.WriteString("ID,Title,StartDate,StartTime,DurationMinutes,StationName,ChannelID,ONID,TSID,SID,EventID,Comment,RecFilePath,Protected\n")
	for _, rec := range recordings {
		sb.WriteString(fmt.Sprintf("%d,%s,%s,%s,%d,%s,%s,%d,%d,%d,%d,%s,%s,%t\n",
			rec.ID,
			escapeCSV(rec.Title),
			rec.StartDate,
			rec.StartTime,
			rec.DurationMinutes(),
			escapeCSV(rec.StationName),
			rec.ChannelID(),
			rec.ONID,
			rec.TSID,
			rec.SID,
			rec.EventID,
			escapeCSV(rec.Comment),
			escapeCSV(rec.RecFilePath),
			rec.IsProtected(),
		))
	}
	return sb.String()
}

// Temporary TSV formatter for recordings
func formatRecordingsAsTSV(recordings []models.RecordingInfo) string {
	var sb strings.Builder
	sb.WriteString("ID\tTitle\tStartDate\tStartTime\tDurationMinutes\tStationName\tChannelID\tONID\tTSID\tSID\tEventID\tComment\tRecFilePath\tProtected\n")
	for _, rec := range recordings {
		sb.WriteString(fmt.Sprintf("%d\t%s\t%s\t%s\t%d\t%s\t%s\t%d\t%d\t%d\t%d\t%s\t%s\t%t\n",
			rec.ID,
			escapeTSV(rec.Title),
			rec.StartDate,
			rec.StartTime,
			rec.DurationMinutes(),
			escapeTSV(rec.StationName),
			rec.ChannelID(),
			rec.ONID,
			rec.TSID,
			rec.SID,
			rec.EventID,
			escapeTSV(rec.Comment),
			escapeTSV(rec.RecFilePath),
			rec.IsProtected(),
		))
	}
	return sb.String()
}
