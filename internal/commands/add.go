package commands

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/epy0n0ff/epgtimer-cli/internal/client"
	"github.com/epy0n0ff/epgtimer-cli/internal/models"
	"github.com/spf13/cobra"
)

var (
	andKey          string
	notKey          string
	serviceList     []string
	serviceListFile string
)

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a new automatic recording rule",
	Long: `Add a new automatic recording rule based on keywords.

The rule will automatically record programs that match the search keywords
on the specified channels.

Example:
  epgtimer add --andKey "ニュース" --serviceList "32736-32736-1024"
  epgtimer add --andKey "ドラマ" --notKey "再放送" --serviceList "32736-32736-1024,32736-32736-1025"
  epgtimer add --andKey "映画" --serviceListFile channels.txt

Channel format: ONID-TSID-SID (e.g., "32736-32736-1024" for NHK総合)

serviceListFile format (one channel per line):
  # Tokyo channels
  32736-32736-1024
  32736-32736-1025
  32737-32737-1032`,
	RunE: runAddCommand,
}

func init() {
	rootCmd.AddCommand(addCmd)

	// Define flags
	addCmd.Flags().StringVar(&andKey, "andKey", "", "Search keywords (required, title must contain these keywords)")
	addCmd.Flags().StringVar(&notKey, "notKey", "", "Exclusion keywords (optional, title must not contain these keywords)")
	addCmd.Flags().StringSliceVar(&serviceList, "serviceList", []string{}, "Channel list in ONID-TSID-SID format (comma-separated)")
	addCmd.Flags().StringVar(&serviceListFile, "serviceListFile", "", "File containing channel list (one channel per line in ONID-TSID-SID format)")

	// Mark required flags - andKey is always required, serviceList or serviceListFile must be provided
	addCmd.MarkFlagRequired("andKey")
}

func runAddCommand(cmd *cobra.Command, args []string) error {
	// Get EMWUI endpoint
	endpoint, err := GetEMWUIEndpoint(cmd)
	if err != nil {
		return fmt.Errorf("configuration error: %w\n\nPlease set EMWUI_ENDPOINT environment variable:\n  export EMWUI_ENDPOINT=http://localhost:5510\n\nOr use --endpoint flag:\n  epgtimer add --endpoint http://localhost:5510 --andKey \"...\" --serviceList \"...\"", err)
	}

	// Validate andKey
	if strings.TrimSpace(andKey) == "" {
		return fmt.Errorf("--andKey cannot be empty")
	}

	// Read serviceList from file if specified
	if serviceListFile != "" {
		channels, err := readServiceListFromFile(serviceListFile)
		if err != nil {
			return fmt.Errorf("failed to read serviceListFile: %w", err)
		}
		serviceList = append(serviceList, channels...)
	}

	// Validate serviceList
	if len(serviceList) == 0 {
		return fmt.Errorf("--serviceList or --serviceListFile must contain at least one channel")
	}

	// Validate serviceList format
	for i, service := range serviceList {
		if _, err := models.ParseServiceListEntry(service); err != nil {
			return fmt.Errorf("invalid channel format in --serviceList[%d]: %w\n\nExpected format: ONID-TSID-SID (e.g., \"32736-32736-1024\")", i, err)
		}
	}

	// Create request
	req := models.NewAutoAddRuleRequest(andKey, notKey, serviceList)

	// Create client
	c := client.NewClient(endpoint)

	// Call API
	_, err = c.SetAutoAdd(req)
	if err != nil {
		// Provide helpful error messages based on error type
		errMsg := err.Error()

		if strings.Contains(errMsg, "connection refused") {
			return fmt.Errorf("connection failed: %w\n\nPlease check:\n  1. EMWUI service is running\n  2. EMWUI_ENDPOINT is correct (current: %s)\n  3. Network connectivity", err, endpoint)
		}

		if strings.Contains(errMsg, "timeout") {
			return fmt.Errorf("connection timeout: %w\n\nEMWUI server did not respond in time (current endpoint: %s)", err, endpoint)
		}

		if strings.Contains(errMsg, "validation failed") {
			return fmt.Errorf("validation error: %w", err)
		}

		if strings.Contains(errMsg, "HTML response") {
			return fmt.Errorf("unexpected response: %w\n\nPossible causes:\n  1. Incorrect endpoint URL\n  2. API path has changed\n  3. EMWUI version incompatibility", err)
		}

		return fmt.Errorf("failed to add recording rule: %w", err)
	}

	// Success
	fmt.Println("✓ Automatic recording rule created successfully")
	fmt.Printf("\nSearch keywords: %s\n", andKey)
	if notKey != "" {
		fmt.Printf("Exclusion keywords: %s\n", notKey)
	}
	fmt.Printf("Channels: %d channels\n", len(serviceList))
	if len(serviceList) <= 10 {
		fmt.Printf("  %s\n", strings.Join(serviceList, ", "))
	} else {
		fmt.Printf("  %s, ... (%d more)\n", strings.Join(serviceList[:10], ", "), len(serviceList)-10)
	}

	return nil
}

// readServiceListFromFile reads channel list from a file
// File format: one channel per line in ONID-TSID-SID format
// Empty lines and lines starting with # are ignored
func readServiceListFromFile(filename string) ([]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	var channels []string
	scanner := bufio.NewScanner(file)
	lineNum := 0

	for scanner.Scan() {
		lineNum++
		line := strings.TrimSpace(scanner.Text())

		// Skip empty lines and comments
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// Validate format
		if _, err := models.ParseServiceListEntry(line); err != nil {
			return nil, fmt.Errorf("invalid format at line %d: %w", lineNum, err)
		}

		channels = append(channels, line)
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading file: %w", err)
	}

	if len(channels) == 0 {
		return nil, fmt.Errorf("file contains no valid channels")
	}

	return channels, nil
}
