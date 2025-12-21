package commands

import (
	"fmt"
	"os"

	"github.com/epy0n0ff/epgtimer-cli/internal/client"
	"github.com/epy0n0ff/epgtimer-cli/internal/formatters"
	"github.com/epy0n0ff/epgtimer-cli/internal/models"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List automatic recording rules",
	Long: `Retrieve and display automatic recording rules from EpgTimer's EMWUI service.

The list command retrieves all automatic recording rules configured in EpgTimer
and displays them in a human-readable table format. Each rule shows its ID,
enabled status, search keywords, exclusion keywords, and number of channels.

Output Formats:
  table  - Human-readable table format (default)
           Best for: Quick viewing in terminal
  json   - JSON format with full rule structure
           Best for: Backups, programmatic processing, piping to jq
  csv    - Comma-separated values
           Best for: Excel/spreadsheet analysis, data import
  tsv    - Tab-separated values
           Best for: Data exchange, simple parsing

Examples:
  # List all rules
  epgtimer list

  # List rules from a specific EMWUI server
  epgtimer list --endpoint http://192.168.1.10:5510

  # Filter by keyword
  epgtimer list --andKey ニュース

  # Filter by channel
  epgtimer list --channel 32736-32736-1024

  # Show only enabled rules
  epgtimer list --enabled

  # Combine filters
  epgtimer list --enabled --andKey ドラマ

  # Export to JSON file (full structure)
  epgtimer list --format json --output rules.json

  # Export filtered rules to CSV (for Excel)
  epgtimer list --enabled --format csv -o enabled_rules.csv

  # Output as TSV to stdout (pipe to other commands)
  epgtimer list --format tsv

  # Backup all rules to JSON
  epgtimer list --format json -o backup-$(date +%Y%m%d).json
`,
	RunE: runList,
}

func init() {
	// Filter flags (Phase 4)
	listCmd.Flags().String("andKey", "", "Filter by search keyword (substring match, case-insensitive)")
	listCmd.Flags().String("channel", "", "Filter by channel (ONID-TSID-SID format, e.g., 32736-32736-1024)")
	listCmd.Flags().Bool("enabled", false, "Show only enabled rules")
	listCmd.Flags().Bool("disabled", false, "Show only disabled rules")
	listCmd.Flags().Bool("regex", false, "Show only regex-enabled rules")

	// Export flags (Phase 5)
	listCmd.Flags().String("format", "table", "Output format: table, json, csv, tsv")
	listCmd.Flags().StringP("output", "o", "", "Output file path (default: stdout)")
}

func runList(cmd *cobra.Command, args []string) error {
	// Get EMWUI endpoint from root command flag or environment variable
	endpoint, err := cmd.Flags().GetString("endpoint")
	if err != nil {
		return fmt.Errorf("failed to get endpoint flag: %w", err)
	}

	if endpoint == "" {
		endpoint = os.Getenv("EMWUI_ENDPOINT")
	}

	if endpoint == "" {
		return fmt.Errorf("EMWUI endpoint not configured\n\nPlease set the endpoint using:\n  1. --endpoint flag: epgtimer list --endpoint http://192.168.1.10:5510\n  2. EMWUI_ENDPOINT environment variable: export EMWUI_ENDPOINT=http://192.168.1.10:5510")
	}

	// Create API client
	apiClient := client.NewClient(endpoint)

	// Retrieve rules
	response, err := apiClient.EnumAutoAdd()
	if err != nil {
		return formatConnectionError(err, endpoint)
	}

	// Build filter options from flags
	filterOpts := models.FilterOptions{}

	andKey, _ := cmd.Flags().GetString("andKey")
	filterOpts.AndKeyFilter = andKey

	channel, _ := cmd.Flags().GetString("channel")
	filterOpts.ChannelFilter = channel

	enabled, _ := cmd.Flags().GetBool("enabled")
	filterOpts.EnabledOnly = enabled

	disabled, _ := cmd.Flags().GetBool("disabled")
	filterOpts.DisabledOnly = disabled

	regex, _ := cmd.Flags().GetBool("regex")
	filterOpts.RegexOnly = regex

	// Apply filters if any are active
	filteredRules := response.Items
	if filterOpts.HasFilters() {
		filteredRules = make([]models.AutoAddRule, 0)
		for _, rule := range response.Items {
			if filterOpts.Matches(&rule) {
				filteredRules = append(filteredRules, rule)
			}
		}

		// Handle empty results
		if len(filteredRules) == 0 {
			fmt.Println("No automatic recording rules match the specified filters.")
			return nil
		}
	}

	// Get format flag
	format, _ := cmd.Flags().GetString("format")

	// Select appropriate formatter
	var formatter interface {
		Format([]models.AutoAddRule) (string, error)
	}

	switch format {
	case "json":
		formatter = &formatters.JSONFormatter{}
	case "csv":
		formatter = &formatters.CSVFormatter{}
	case "tsv":
		formatter = &formatters.TSVFormatter{}
	case "table":
		formatter = &formatters.TableFormatter{}
	default:
		return fmt.Errorf("unsupported format '%s'. Supported formats: table, json, csv, tsv", format)
	}

	// Format rules
	output, err := formatter.Format(filteredRules)
	if err != nil {
		return fmt.Errorf("failed to format output: %w", err)
	}

	// Get output flag
	outputPath, _ := cmd.Flags().GetString("output")

	// Write to file or stdout
	if outputPath != "" {
		// Write to file
		if err := os.WriteFile(outputPath, []byte(output), 0644); err != nil {
			return fmt.Errorf("failed to write output to file '%s': %w", outputPath, err)
		}
		fmt.Printf("Successfully exported %d rules to %s\n", len(filteredRules), outputPath)
	} else {
		// Write to stdout
		fmt.Print(output)
	}

	return nil
}

// formatConnectionError formats connection errors with troubleshooting guidance
func formatConnectionError(err error, endpoint string) error {
	return fmt.Errorf(`Failed to connect to EMWUI service at %s
Error: %v

Troubleshooting:
1. Check that EpgTimer is running
2. Verify EMWUI_ENDPOINT environment variable or --endpoint flag
3. Confirm network connectivity to the EMWUI server
4. Ensure the EMWUI web interface is accessible`, endpoint, err)
}
