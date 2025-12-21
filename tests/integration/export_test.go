package integration

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/epy0n0ff/epgtimer-cli/internal/client"
	"github.com/epy0n0ff/epgtimer-cli/internal/formatters"
	"github.com/epy0n0ff/epgtimer-cli/internal/models"
	"github.com/epy0n0ff/epgtimer-cli/tests/testdata"
)

// TestJSONFormatter tests JSON export format
func TestJSONFormatter(t *testing.T) {
	// Create mock server
	mock := testdata.NewMockEMWUIServer()
	defer mock.Close()

	// Create client and retrieve rules
	apiClient := client.NewClient(mock.URL())
	response, err := apiClient.EnumAutoAdd()
	if err != nil {
		t.Fatalf("EnumAutoAdd() failed: %v", err)
	}

	// Format as JSON
	formatter := &formatters.JSONFormatter{}
	output, err := formatter.Format(response.Items)
	if err != nil {
		t.Fatalf("JSON Format() failed: %v", err)
	}

	// Verify output is valid JSON
	var parsed []models.AutoAddRule
	if err := json.Unmarshal([]byte(output), &parsed); err != nil {
		t.Fatalf("Output is not valid JSON: %v\nOutput: %s", err, output)
	}

	// Verify count matches
	if len(parsed) != len(response.Items) {
		t.Errorf("Expected %d rules in JSON, got %d", len(response.Items), len(parsed))
	}

	// Verify first rule data
	if len(parsed) > 0 {
		if parsed[0].ID != 1 {
			t.Errorf("Expected first rule ID=1, got %d", parsed[0].ID)
		}
		if parsed[0].SearchSettings.AndKey != "サイエンスZERO" {
			t.Errorf("Expected first rule AndKey='サイエンスZERO', got '%s'", parsed[0].SearchSettings.AndKey)
		}
	}

	// Verify JSON is pretty-printed (indented)
	if !strings.Contains(output, "\n") || !strings.Contains(output, "  ") {
		t.Error("JSON output should be pretty-printed with indentation")
	}
}

// TestCSVFormatter tests CSV export format
func TestCSVFormatter(t *testing.T) {
	// Create mock server
	mock := testdata.NewMockEMWUIServer()
	defer mock.Close()

	// Create client and retrieve rules
	apiClient := client.NewClient(mock.URL())
	response, err := apiClient.EnumAutoAdd()
	if err != nil {
		t.Fatalf("EnumAutoAdd() failed: %v", err)
	}

	// Format as CSV
	formatter := &formatters.CSVFormatter{}
	output, err := formatter.Format(response.Items)
	if err != nil {
		t.Fatalf("CSV Format() failed: %v", err)
	}

	// Verify CSV structure
	lines := strings.Split(strings.TrimSpace(output), "\n")
	expectedLines := len(response.Items) + 1 // +1 for header

	if len(lines) != expectedLines {
		t.Errorf("Expected %d lines in CSV (1 header + %d data), got %d", expectedLines, len(response.Items), len(lines))
	}

	// Verify header line
	header := lines[0]
	requiredHeaders := []string{"ID", "Enabled", "AndKey", "NotKey", "RegExp", "Channels", "ChannelCount"}
	for _, h := range requiredHeaders {
		if !strings.Contains(header, h) {
			t.Errorf("CSV header missing required field: %s", h)
		}
	}

	// Verify first data line contains expected data
	if len(lines) > 1 {
		firstLine := lines[1]
		// Should contain rule ID 1 and keyword サイエンスZERO
		if !strings.Contains(firstLine, "1") {
			t.Error("First CSV line should contain ID=1")
		}
		if !strings.Contains(firstLine, "サイエンスZERO") {
			t.Error("First CSV line should contain 'サイエンスZERO'")
		}
	}

	// Verify CSV uses commas
	if !strings.Contains(output, ",") {
		t.Error("CSV output should contain comma separators")
	}
}

// TestTSVFormatter tests TSV export format
func TestTSVFormatter(t *testing.T) {
	// Create mock server
	mock := testdata.NewMockEMWUIServer()
	defer mock.Close()

	// Create client and retrieve rules
	apiClient := client.NewClient(mock.URL())
	response, err := apiClient.EnumAutoAdd()
	if err != nil {
		t.Fatalf("EnumAutoAdd() failed: %v", err)
	}

	// Format as TSV
	formatter := &formatters.TSVFormatter{}
	output, err := formatter.Format(response.Items)
	if err != nil {
		t.Fatalf("TSV Format() failed: %v", err)
	}

	// Verify TSV structure
	lines := strings.Split(strings.TrimSpace(output), "\n")
	expectedLines := len(response.Items) + 1 // +1 for header

	if len(lines) != expectedLines {
		t.Errorf("Expected %d lines in TSV (1 header + %d data), got %d", expectedLines, len(response.Items), len(lines))
	}

	// Verify header line
	header := lines[0]
	requiredHeaders := []string{"ID", "Enabled", "AndKey", "NotKey", "RegExp", "Channels", "ChannelCount"}
	for _, h := range requiredHeaders {
		if !strings.Contains(header, h) {
			t.Errorf("TSV header missing required field: %s", h)
		}
	}

	// Verify first data line contains expected data
	if len(lines) > 1 {
		firstLine := lines[1]
		// Should contain rule ID 1 and keyword サイエンスZERO
		if !strings.Contains(firstLine, "1") {
			t.Error("First TSV line should contain ID=1")
		}
		if !strings.Contains(firstLine, "サイエンスZERO") {
			t.Error("First TSV line should contain 'サイエンスZERO'")
		}
	}

	// Verify TSV uses tabs (not commas)
	if !strings.Contains(output, "\t") {
		t.Error("TSV output should contain tab separators")
	}

	// TSV should not use commas as primary separator
	// (may contain commas in data, but tabs should be more frequent in header)
	tabCount := strings.Count(header, "\t")
	commaCount := strings.Count(header, ",")
	if tabCount == 0 || tabCount < commaCount {
		t.Error("TSV should use tabs as primary separator, not commas")
	}
}

// TestExportEmptyList tests export formats with empty rule list
func TestExportEmptyList(t *testing.T) {
	emptyRules := []models.AutoAddRule{}

	// Test JSON with empty list
	jsonFormatter := &formatters.JSONFormatter{}
	jsonOutput, err := jsonFormatter.Format(emptyRules)
	if err != nil {
		t.Errorf("JSON formatter should handle empty list: %v", err)
	}
	if !strings.Contains(jsonOutput, "[") || !strings.Contains(jsonOutput, "]") {
		t.Error("JSON output for empty list should contain empty array []")
	}

	// Test CSV with empty list (should have header only)
	csvFormatter := &formatters.CSVFormatter{}
	csvOutput, err := csvFormatter.Format(emptyRules)
	if err != nil {
		t.Errorf("CSV formatter should handle empty list: %v", err)
	}
	csvLines := strings.Split(strings.TrimSpace(csvOutput), "\n")
	if len(csvLines) != 1 {
		t.Errorf("CSV with empty data should have 1 line (header only), got %d", len(csvLines))
	}

	// Test TSV with empty list (should have header only)
	tsvFormatter := &formatters.TSVFormatter{}
	tsvOutput, err := tsvFormatter.Format(emptyRules)
	if err != nil {
		t.Errorf("TSV formatter should handle empty list: %v", err)
	}
	tsvLines := strings.Split(strings.TrimSpace(tsvOutput), "\n")
	if len(tsvLines) != 1 {
		t.Errorf("TSV with empty data should have 1 line (header only), got %d", len(tsvLines))
	}
}

// TestExportWithSpecialCharacters tests handling of special characters in export
func TestExportWithSpecialCharacters(t *testing.T) {
	// Create mock server
	mock := testdata.NewMockEMWUIServer()
	defer mock.Close()

	// Create client and retrieve rules
	apiClient := client.NewClient(mock.URL())
	response, err := apiClient.EnumAutoAdd()
	if err != nil {
		t.Fatalf("EnumAutoAdd() failed: %v", err)
	}

	// Test JSON handles Japanese characters
	jsonFormatter := &formatters.JSONFormatter{}
	jsonOutput, err := jsonFormatter.Format(response.Items)
	if err != nil {
		t.Fatalf("JSON Format() failed: %v", err)
	}
	if !strings.Contains(jsonOutput, "サイエンスZERO") {
		t.Error("JSON should preserve Japanese characters")
	}

	// Test CSV handles Japanese characters
	csvFormatter := &formatters.CSVFormatter{}
	csvOutput, err := csvFormatter.Format(response.Items)
	if err != nil {
		t.Fatalf("CSV Format() failed: %v", err)
	}
	if !strings.Contains(csvOutput, "サイエンスZERO") {
		t.Error("CSV should preserve Japanese characters")
	}

	// Test TSV handles Japanese characters
	tsvFormatter := &formatters.TSVFormatter{}
	tsvOutput, err := tsvFormatter.Format(response.Items)
	if err != nil {
		t.Fatalf("TSV Format() failed: %v", err)
	}
	if !strings.Contains(tsvOutput, "サイエンスZERO") {
		t.Error("TSV should preserve Japanese characters")
	}

	// Test CSV/TSV handle special characters in NotKey (brackets)
	if !strings.Contains(csvOutput, "[再]") {
		t.Error("CSV should preserve brackets in NotKey")
	}
	if !strings.Contains(tsvOutput, "[再]") {
		t.Error("TSV should preserve brackets in NotKey")
	}
}
