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

var reservationsCmd = &cobra.Command{
	Use:   "reservations",
	Short: "List manual reservations",
	Long: `Retrieve and display manual reservations from EpgTimer's EMWUI service.

The reservations command retrieves all manual recording reservations configured in EpgTimer
and displays them in a human-readable table format. Each reservation shows its ID, date,
time, program title, and station name.

Output Formats:
  table  - Human-readable table format (default)
           Best for: Quick viewing in terminal
  json   - JSON format with full reservation structure
           Best for: Programmatic processing, piping to jq
  csv    - Comma-separated values
           Best for: Excel/spreadsheet analysis, data import
  tsv    - Tab-separated values
           Best for: Data exchange, simple parsing

Examples:
  # List all reservations
  epgtimer reservations

  # List reservations from a specific EMWUI server
  epgtimer reservations --endpoint http://192.168.1.10:5510

  # Filter by title
  epgtimer reservations --title "ニュース"

  # Filter by station
  epgtimer reservations --station "NHK"

  # Export to JSON file
  epgtimer reservations --format json --output reservations.json

  # Export to CSV
  epgtimer reservations --format csv -o reservations.csv
`,
	RunE: runReservations,
}

func init() {
	// Filter flags
	reservationsCmd.Flags().String("title", "", "Filter by title (substring match, case-insensitive)")
	reservationsCmd.Flags().String("station", "", "Filter by station name (substring match, case-insensitive)")
	reservationsCmd.Flags().String("channel", "", "Filter by channel ID (exact match, format: ONID-TSID-SID)")

	// Export flags
	reservationsCmd.Flags().String("format", "table", "Output format: table, json, csv, tsv")
	reservationsCmd.Flags().StringP("output", "o", "", "Output file path (default: stdout)")
}

func runReservations(cmd *cobra.Command, args []string) error {
	// Get EMWUI endpoint
	endpoint, err := cmd.Flags().GetString("endpoint")
	if err != nil {
		return fmt.Errorf("failed to get endpoint flag: %w", err)
	}

	if endpoint == "" {
		endpoint = os.Getenv("EMWUI_ENDPOINT")
	}

	if endpoint == "" {
		return fmt.Errorf("EMWUI endpoint not configured\n\nPlease set the endpoint using:\n  1. --endpoint flag: epgtimer reservations --endpoint http://192.168.1.10:5510\n  2. EMWUI_ENDPOINT environment variable: export EMWUI_ENDPOINT=http://192.168.1.10:5510")
	}

	// Create API client
	apiClient := client.NewClient(endpoint)

	// Retrieve reservations
	response, err := apiClient.EnumReserveInfo()
	if err != nil {
		return formatConnectionError(err, endpoint)
	}

	// Apply filters
	filteredReservations := applyReservationFilters(cmd, response.Items)

	// Handle empty results
	if len(filteredReservations) == 0 {
		fmt.Println("No reservations match the specified filters.")
		return nil
	}

	// Get format flag
	format, _ := cmd.Flags().GetString("format")

	// Format reservations
	var output string
	switch format {
	case "table":
		formatter := &formatters.ReservationsTableFormatter{}
		output, err = formatter.Format(filteredReservations)
	case "json":
		output = formatReservationsAsJSON(filteredReservations)
	case "csv":
		output = formatReservationsAsCSV(filteredReservations)
	case "tsv":
		output = formatReservationsAsTSV(filteredReservations)
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
		fmt.Printf("Successfully exported %d reservations to %s\n", len(filteredReservations), outputPath)
	} else {
		fmt.Print(output)
	}

	return nil
}

func applyReservationFilters(cmd *cobra.Command, reservations []models.ReservationInfo) []models.ReservationInfo {
	title, _ := cmd.Flags().GetString("title")
	station, _ := cmd.Flags().GetString("station")
	channel, _ := cmd.Flags().GetString("channel")

	var filtered []models.ReservationInfo
	for _, res := range reservations {
		// Title filter
		if title != "" {
			titleLower := strings.ToLower(res.Title)
			filterLower := strings.ToLower(title)
			if !strings.Contains(titleLower, filterLower) {
				continue
			}
		}

		// Station filter
		if station != "" {
			stationLower := strings.ToLower(res.StationName)
			filterLower := strings.ToLower(station)
			if !strings.Contains(stationLower, filterLower) {
				continue
			}
		}

		// Channel filter
		if channel != "" {
			if res.ChannelID() != channel {
				continue
			}
		}

		filtered = append(filtered, res)
	}

	return filtered
}

// Temporary JSON formatter for reservations
func formatReservationsAsJSON(reservations []models.ReservationInfo) string {
	var sb strings.Builder
	sb.WriteString("[\n")
	for i, res := range reservations {
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
    "comment": "%s"
  }`,
			res.ID,
			escapeJSON(res.Title),
			res.StartDate,
			res.StartTime,
			res.DurationSecond,
			res.DurationMinutes(),
			escapeJSON(res.StationName),
			res.ChannelID(),
			res.ONID,
			res.TSID,
			res.SID,
			res.EventID,
			escapeJSON(res.Comment),
		))
		if i < len(reservations)-1 {
			sb.WriteString(",\n")
		} else {
			sb.WriteString("\n")
		}
	}
	sb.WriteString("]\n")
	return sb.String()
}

// Temporary CSV formatter for reservations
func formatReservationsAsCSV(reservations []models.ReservationInfo) string {
	var sb strings.Builder
	sb.WriteString("ID,Title,StartDate,StartTime,DurationMinutes,StationName,ChannelID,ONID,TSID,SID,EventID,Comment\n")
	for _, res := range reservations {
		sb.WriteString(fmt.Sprintf("%d,%s,%s,%s,%d,%s,%s,%d,%d,%d,%d,%s\n",
			res.ID,
			escapeCSV(res.Title),
			res.StartDate,
			res.StartTime,
			res.DurationMinutes(),
			escapeCSV(res.StationName),
			res.ChannelID(),
			res.ONID,
			res.TSID,
			res.SID,
			res.EventID,
			escapeCSV(res.Comment),
		))
	}
	return sb.String()
}

// Temporary TSV formatter for reservations
func formatReservationsAsTSV(reservations []models.ReservationInfo) string {
	var sb strings.Builder
	sb.WriteString("ID\tTitle\tStartDate\tStartTime\tDurationMinutes\tStationName\tChannelID\tONID\tTSID\tSID\tEventID\tComment\n")
	for _, res := range reservations {
		sb.WriteString(fmt.Sprintf("%d\t%s\t%s\t%s\t%d\t%s\t%s\t%d\t%d\t%d\t%d\t%s\n",
			res.ID,
			escapeTSV(res.Title),
			res.StartDate,
			res.StartTime,
			res.DurationMinutes(),
			escapeTSV(res.StationName),
			res.ChannelID(),
			res.ONID,
			res.TSID,
			res.SID,
			res.EventID,
			escapeTSV(res.Comment),
		))
	}
	return sb.String()
}

// Helper functions for escaping special characters
func escapeJSON(s string) string {
	s = strings.ReplaceAll(s, "\\", "\\\\")
	s = strings.ReplaceAll(s, "\"", "\\\"")
	s = strings.ReplaceAll(s, "\n", "\\n")
	s = strings.ReplaceAll(s, "\r", "\\r")
	s = strings.ReplaceAll(s, "\t", "\\t")
	return s
}

func escapeCSV(s string) string {
	if strings.ContainsAny(s, ",\"\n\r") {
		s = strings.ReplaceAll(s, "\"", "\"\"")
		return "\"" + s + "\""
	}
	return s
}

func escapeTSV(s string) string {
	s = strings.ReplaceAll(s, "\t", " ")
	s = strings.ReplaceAll(s, "\n", " ")
	s = strings.ReplaceAll(s, "\r", " ")
	return s
}
