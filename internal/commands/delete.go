package commands

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/epy0n0ff/epgtimer-cli/internal/client"
	"github.com/spf13/cobra"
)

// deleteCmd represents the delete command
var deleteCmd = &cobra.Command{
	Use:   "delete [rule-id]",
	Short: "Delete an automatic recording rule",
	Long: `Delete an automatic recording rule by its ID.

To find the rule ID, use the 'list' command to see all automatic recording rules.

Example:
  # List all rules to find the ID
  epgtimer list

  # Delete rule with ID 334
  epgtimer delete 334

  # Using --id flag
  epgtimer delete --id 334`,
	Args: cobra.MaximumNArgs(1),
	RunE: runDeleteCommand,
}

var deleteRuleID int

func init() {
	rootCmd.AddCommand(deleteCmd)

	// Define flags
	deleteCmd.Flags().IntVar(&deleteRuleID, "id", 0, "Rule ID to delete")
}

func runDeleteCommand(cmd *cobra.Command, args []string) error {
	// Get EMWUI endpoint
	endpoint, err := GetEMWUIEndpoint(cmd)
	if err != nil {
		return fmt.Errorf("configuration error: %w\n\nPlease set EMWUI_ENDPOINT environment variable:\n  export EMWUI_ENDPOINT=http://localhost:5510\n\nOr use --endpoint flag:\n  epgtimer delete --endpoint http://localhost:5510 334", err)
	}

	// Determine rule ID from args or flag
	var ruleID int
	if len(args) == 1 {
		// Parse from positional argument
		id, err := strconv.Atoi(args[0])
		if err != nil {
			return fmt.Errorf("invalid rule ID '%s': must be a number", args[0])
		}
		ruleID = id
	} else if deleteRuleID > 0 {
		// Use flag value
		ruleID = deleteRuleID
	} else {
		return fmt.Errorf("rule ID is required\n\nUsage:\n  epgtimer delete [rule-id]\n  epgtimer delete --id [rule-id]\n\nTo find rule IDs, run:\n  epgtimer list")
	}

	// Validate rule ID
	if ruleID <= 0 {
		return fmt.Errorf("invalid rule ID: must be greater than 0")
	}

	// Create client
	c := client.NewClient(endpoint)

	// Call API
	_, err = c.DeleteAutoAdd(ruleID)
	if err != nil {
		// Provide helpful error messages based on error type
		errMsg := err.Error()

		if strings.Contains(errMsg, "connection refused") {
			return fmt.Errorf("connection failed: %w\n\nPlease check:\n  1. EMWUI service is running\n  2. EMWUI_ENDPOINT is correct (current: %s)\n  3. Network connectivity", err, endpoint)
		}

		if strings.Contains(errMsg, "timeout") {
			return fmt.Errorf("connection timeout: %w\n\nEMWUI server did not respond in time (current endpoint: %s)", err, endpoint)
		}

		if strings.Contains(errMsg, "invalid rule ID") {
			return fmt.Errorf("validation error: %w", err)
		}

		if strings.Contains(errMsg, "HTML response") {
			return fmt.Errorf("unexpected response: %w\n\nPossible causes:\n  1. Incorrect endpoint URL\n  2. API path has changed\n  3. EMWUI version incompatibility", err)
		}

		return fmt.Errorf("failed to delete recording rule: %w", err)
	}

	// Success
	fmt.Printf("✓ Automatic recording rule (ID: %d) deleted successfully\n", ruleID)

	return nil
}
