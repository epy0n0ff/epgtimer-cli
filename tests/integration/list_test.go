package integration

import (
	"strings"
	"testing"

	"github.com/epy0n0ff/epgtimer-cli/internal/client"
	"github.com/epy0n0ff/epgtimer-cli/tests/testdata"
)

// TestEnumAutoAdd_Success tests successful rule retrieval
func TestEnumAutoAdd_Success(t *testing.T) {
	// Create mock server
	mock := testdata.NewMockEMWUIServer()
	defer mock.Close()

	// Create client
	apiClient := client.NewClient(mock.URL())

	// Call EnumAutoAdd
	response, err := apiClient.EnumAutoAdd()
	if err != nil {
		t.Fatalf("EnumAutoAdd() failed: %v", err)
	}

	// Verify response structure
	if response.Total != 3 {
		t.Errorf("Expected Total=3, got %d", response.Total)
	}

	if response.Count != 3 {
		t.Errorf("Expected Count=3, got %d", response.Count)
	}

	if len(response.Items) != 3 {
		t.Fatalf("Expected 3 rules, got %d", len(response.Items))
	}

	// Verify first rule
	rule1 := response.Items[0]
	if rule1.ID != 1 {
		t.Errorf("Expected ID=1, got %d", rule1.ID)
	}

	if rule1.SearchSettings.AndKey != "サイエンスZERO" {
		t.Errorf("Expected AndKey='サイエンスZERO', got '%s'", rule1.SearchSettings.AndKey)
	}

	if rule1.SearchSettings.NotKey != "[再]" {
		t.Errorf("Expected NotKey='[再]', got '%s'", rule1.SearchSettings.NotKey)
	}

	if !rule1.SearchSettings.IsEnabled() {
		t.Error("Expected rule1 to be enabled")
	}

	if rule1.SearchSettings.ChannelCount() != 2 {
		t.Errorf("Expected 2 channels, got %d", rule1.SearchSettings.ChannelCount())
	}

	// Verify second rule
	rule2 := response.Items[1]
	if rule2.SearchSettings.AndKey != "ブラタモリ" {
		t.Errorf("Expected AndKey='ブラタモリ', got '%s'", rule2.SearchSettings.AndKey)
	}

	// Verify third rule (disabled, regex)
	rule3 := response.Items[2]
	if rule3.SearchSettings.IsEnabled() {
		t.Error("Expected rule3 to be disabled")
	}

	if !rule3.SearchSettings.IsRegex() {
		t.Error("Expected rule3 to have regex enabled")
	}

	if rule3.SearchSettings.AndKey != "^NHKニュース" {
		t.Errorf("Expected AndKey='^NHKニュース', got '%s'", rule3.SearchSettings.AndKey)
	}
}

// TestEnumAutoAdd_Empty tests handling of empty response
func TestEnumAutoAdd_Empty(t *testing.T) {
	// Create mock server configured for empty response
	mock := testdata.NewEmptyEnumAutoAddServer()
	defer mock.Close()

	// Create client
	apiClient := client.NewClient(mock.URL())

	// Call EnumAutoAdd
	response, err := apiClient.EnumAutoAdd()
	if err != nil {
		t.Fatalf("EnumAutoAdd() failed: %v", err)
	}

	// Verify empty response
	if response.Total != 0 {
		t.Errorf("Expected Total=0, got %d", response.Total)
	}

	if response.Count != 0 {
		t.Errorf("Expected Count=0, got %d", response.Count)
	}

	if len(response.Items) != 0 {
		t.Errorf("Expected 0 rules, got %d", len(response.Items))
	}
}

// TestEnumAutoAdd_ConnectionError tests connection error handling
func TestEnumAutoAdd_ConnectionError(t *testing.T) {
	// Create client with invalid endpoint
	apiClient := client.NewClient("http://localhost:99999")

	// Call EnumAutoAdd (should fail)
	_, err := apiClient.EnumAutoAdd()
	if err == nil {
		t.Fatal("Expected connection error, got nil")
	}

	// Verify error message contains helpful context
	errMsg := err.Error()
	if !strings.Contains(errMsg, "failed to connect") {
		t.Errorf("Error message should mention connection failure, got: %s", errMsg)
	}
}

// TestEnumAutoAdd_MalformedXML tests handling of malformed XML response
func TestEnumAutoAdd_MalformedXML(t *testing.T) {
	// Create mock server that returns invalid XML
	mock := testdata.NewMockEMWUIServer()
	mock.SetEnumAutoAddHandler(func() (string, int) {
		return "<invalid><xml>", 200
	})
	defer mock.Close()

	// Create client
	apiClient := client.NewClient(mock.URL())

	// Call EnumAutoAdd (should fail to parse)
	_, err := apiClient.EnumAutoAdd()
	if err == nil {
		t.Fatal("Expected XML parse error, got nil")
	}

	// Verify error message mentions parsing failure
	errMsg := err.Error()
	if !strings.Contains(errMsg, "failed to parse XML") {
		t.Errorf("Error message should mention XML parsing failure, got: %s", errMsg)
	}
}

// TestEnumAutoAdd_JapaneseEncoding tests proper handling of Japanese characters
func TestEnumAutoAdd_JapaneseEncoding(t *testing.T) {
	// Create mock server
	mock := testdata.NewMockEMWUIServer()
	defer mock.Close()

	// Create client
	apiClient := client.NewClient(mock.URL())

	// Call EnumAutoAdd
	response, err := apiClient.EnumAutoAdd()
	if err != nil {
		t.Fatalf("EnumAutoAdd() failed: %v", err)
	}

	// Verify Japanese characters are preserved correctly
	japaneseKeywords := []string{"サイエンスZERO", "ブラタモリ", "NHKニュース"}

	for i, expected := range japaneseKeywords {
		if i >= len(response.Items) {
			break
		}
		actual := response.Items[i].SearchSettings.AndKey
		// Check if it contains expected Japanese characters (may have prefix like ^)
		if !strings.Contains(actual, strings.TrimPrefix(expected, "^")) {
			t.Errorf("Rule %d: Japanese characters not preserved correctly. Expected to contain '%s', got '%s'",
				i+1, expected, actual)
		}
	}
}

// TestServiceInfo_String tests the String() method for channel formatting
func TestServiceInfo_String(t *testing.T) {
	// Create mock server
	mock := testdata.NewMockEMWUIServer()
	defer mock.Close()

	// Create client
	apiClient := client.NewClient(mock.URL())

	// Call EnumAutoAdd
	response, err := apiClient.EnumAutoAdd()
	if err != nil {
		t.Fatalf("EnumAutoAdd() failed: %v", err)
	}

	// Test String() formatting
	if len(response.Items) > 0 && len(response.Items[0].SearchSettings.ServiceList) > 0 {
		channel := response.Items[0].SearchSettings.ServiceList[0]
		channelStr := channel.String()

		expectedFormat := "32736-32736-1024"
		if channelStr != expectedFormat {
			t.Errorf("Expected channel format '%s', got '%s'", expectedFormat, channelStr)
		}
	}
}
