package commands

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	// Version information
	Version = "0.1.0"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "epgtimer",
	Short: "CLI for EpgTimer EMWUI interface",
	Long: `epgtimer is a command-line interface for EpgTimer's EMWUI.
It allows you to manage automatic recording rules via keyword-based searches.

Set the EMWUI_ENDPOINT environment variable to your EMWUI server URL:
  export EMWUI_ENDPOINT=http://localhost:5510

Example:
  epgtimer add --andKey "ニュース" --serviceList "32736-32736-1024"`,
	Version: Version,
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	// Global flags can be added here
	rootCmd.PersistentFlags().StringP("endpoint", "e", "", "EMWUI server endpoint (overrides EMWUI_ENDPOINT env var)")

	// Register subcommands
	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(channelsCmd)
	rootCmd.AddCommand(reservationsCmd)
	rootCmd.AddCommand(recordingsCmd)
	rootCmd.AddCommand(epgCmd)
}

// GetEMWUIEndpoint returns the EMWUI endpoint from flag or environment variable
func GetEMWUIEndpoint(cmd *cobra.Command) (string, error) {
	// First try flag
	endpoint, err := cmd.Flags().GetString("endpoint")
	if err == nil && endpoint != "" {
		return endpoint, nil
	}

	// Then try environment variable
	endpoint = os.Getenv("EMWUI_ENDPOINT")
	if endpoint == "" {
		return "", fmt.Errorf("EMWUI_ENDPOINT is not set. Please set the environment variable or use --endpoint flag")
	}

	return endpoint, nil
}
