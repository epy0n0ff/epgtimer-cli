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

var channelsCmd = &cobra.Command{
	Use:   "channels",
	Short: "List available channels",
	Long: `Retrieve and display available channels from EpgTimer's EMWUI service.

The channels command retrieves all available channels (services) configured in EpgTimer
and displays them in a human-readable table format. Each channel shows its identifier,
type (TV/Radio/Data), remote control key ID, name, and network.

Output Formats:
  table  - Human-readable table format (default)
           Best for: Quick viewing in terminal
  json   - JSON format with full channel structure
           Best for: Programmatic processing, piping to jq
  csv    - Comma-separated values
           Best for: Excel/spreadsheet analysis, data import
  tsv    - Tab-separated values
           Best for: Data exchange, simple parsing

Examples:
  # List all channels
  epgtimer channels

  # List channels from a specific EMWUI server
  epgtimer channels --endpoint http://192.168.1.10:5510

  # Show only TV channels
  epgtimer channels --tv

  # Show only radio channels
  epgtimer channels --radio

  # Filter by network name
  epgtimer channels --network "BS Digital"

  # Export to JSON file
  epgtimer channels --format json --output channels.json

  # Export TV channels to CSV
  epgtimer channels --tv --format csv -o tv_channels.csv
`,
	RunE: runChannels,
}

func init() {
	// Filter flags
	channelsCmd.Flags().Bool("tv", false, "Show only TV channels (service_type=1)")
	channelsCmd.Flags().Bool("radio", false, "Show only radio channels (service_type=2)")
	channelsCmd.Flags().Bool("data", false, "Show only data channels (service_type=192)")
	channelsCmd.Flags().String("network", "", "Filter by network name (substring match, case-insensitive)")
	channelsCmd.Flags().String("name", "", "Filter by channel name (substring match, case-insensitive)")

	// Export flags
	channelsCmd.Flags().String("format", "table", "Output format: table, json, csv, tsv")
	channelsCmd.Flags().StringP("output", "o", "", "Output file path (default: stdout)")
}

func runChannels(cmd *cobra.Command, args []string) error {
	// Get EMWUI endpoint
	endpoint, err := cmd.Flags().GetString("endpoint")
	if err != nil {
		return fmt.Errorf("failed to get endpoint flag: %w", err)
	}

	if endpoint == "" {
		endpoint = os.Getenv("EMWUI_ENDPOINT")
	}

	if endpoint == "" {
		return fmt.Errorf("EMWUI endpoint not configured\n\nPlease set the endpoint using:\n  1. --endpoint flag: epgtimer channels --endpoint http://192.168.1.10:5510\n  2. EMWUI_ENDPOINT environment variable: export EMWUI_ENDPOINT=http://192.168.1.10:5510")
	}

	// Create API client
	apiClient := client.NewClient(endpoint)

	// Retrieve channels
	response, err := apiClient.EnumService()
	if err != nil {
		return formatConnectionError(err, endpoint)
	}

	// Apply filters
	filteredChannels := applyChannelFilters(cmd, response.Items)

	// Handle empty results
	if len(filteredChannels) == 0 {
		fmt.Println("No channels match the specified filters.")
		return nil
	}

	// Get format flag
	format, _ := cmd.Flags().GetString("format")

	// Format channels
	var output string
	switch format {
	case "table":
		formatter := &formatters.ChannelsTableFormatter{}
		output, err = formatter.Format(filteredChannels)
	case "json":
		// Use custom JSON formatter for channels
		output = formatChannelsAsJSON(filteredChannels)
	case "csv":
		output = formatChannelsAsCSV(filteredChannels)
	case "tsv":
		output = formatChannelsAsTSV(filteredChannels)
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
		fmt.Printf("Successfully exported %d channels to %s\n", len(filteredChannels), outputPath)
	} else {
		fmt.Print(output)
	}

	return nil
}

func applyChannelFilters(cmd *cobra.Command, channels []models.ChannelInfo) []models.ChannelInfo {
	tv, _ := cmd.Flags().GetBool("tv")
	radio, _ := cmd.Flags().GetBool("radio")
	data, _ := cmd.Flags().GetBool("data")
	network, _ := cmd.Flags().GetString("network")
	name, _ := cmd.Flags().GetString("name")

	var filtered []models.ChannelInfo
	for _, ch := range channels {
		// Service type filter
		if tv && !ch.IsTV() {
			continue
		}
		if radio && !ch.IsRadio() {
			continue
		}
		if data && !ch.IsData() {
			continue
		}

		// Network filter
		if network != "" {
			networkLower := strings.ToLower(ch.NetworkName)
			filterLower := strings.ToLower(network)
			if !strings.Contains(networkLower, filterLower) {
				continue
			}
		}

		// Name filter
		if name != "" {
			nameLower := strings.ToLower(ch.ServiceName)
			filterLower := strings.ToLower(name)
			if !strings.Contains(nameLower, filterLower) {
				continue
			}
		}

		filtered = append(filtered, ch)
	}

	return filtered
}

// Temporary JSON formatter for channels
func formatChannelsAsJSON(channels []models.ChannelInfo) string {
	// Simple JSON formatting
	var sb strings.Builder
	sb.WriteString("[\n")
	for i, ch := range channels {
		sb.WriteString(fmt.Sprintf(`  {
    "channel_id": "%s",
    "onid": %d,
    "tsid": %d,
    "sid": %d,
    "service_type": %d,
    "service_type_name": "%s",
    "service_name": "%s",
    "service_provider_name": "%s",
    "network_name": "%s",
    "ts_name": "%s",
    "remote_control_key_id": %d
  }`,
			ch.ChannelID(),
			ch.ONID,
			ch.TSID,
			ch.SID,
			ch.ServiceType,
			ch.ServiceTypeString(),
			ch.ServiceName,
			ch.ServiceProviderName,
			ch.NetworkName,
			ch.TSName,
			ch.RemoteControlKeyID,
		))
		if i < len(channels)-1 {
			sb.WriteString(",\n")
		} else {
			sb.WriteString("\n")
		}
	}
	sb.WriteString("]\n")
	return sb.String()
}

// Temporary CSV formatter for channels
func formatChannelsAsCSV(channels []models.ChannelInfo) string {
	var sb strings.Builder
	sb.WriteString("ChannelID,ONID,TSID,SID,ServiceType,ServiceTypeName,ServiceName,ServiceProviderName,NetworkName,TSName,RemoteControlKeyID\n")
	for _, ch := range channels {
		sb.WriteString(fmt.Sprintf("%s,%d,%d,%d,%d,%s,%s,%s,%s,%s,%d\n",
			ch.ChannelID(),
			ch.ONID,
			ch.TSID,
			ch.SID,
			ch.ServiceType,
			ch.ServiceTypeString(),
			ch.ServiceName,
			ch.ServiceProviderName,
			ch.NetworkName,
			ch.TSName,
			ch.RemoteControlKeyID,
		))
	}
	return sb.String()
}

// Temporary TSV formatter for channels
func formatChannelsAsTSV(channels []models.ChannelInfo) string {
	var sb strings.Builder
	sb.WriteString("ChannelID\tONID\tTSID\tSID\tServiceType\tServiceTypeName\tServiceName\tServiceProviderName\tNetworkName\tTSName\tRemoteControlKeyID\n")
	for _, ch := range channels {
		sb.WriteString(fmt.Sprintf("%s\t%d\t%d\t%d\t%d\t%s\t%s\t%s\t%s\t%s\t%d\n",
			ch.ChannelID(),
			ch.ONID,
			ch.TSID,
			ch.SID,
			ch.ServiceType,
			ch.ServiceTypeString(),
			ch.ServiceName,
			ch.ServiceProviderName,
			ch.NetworkName,
			ch.TSName,
			ch.RemoteControlKeyID,
		))
	}
	return sb.String()
}
