package integration

import (
	"strings"
	"testing"

	"github.com/epy0n0ff/epgtimer-cli/internal/client"
	"github.com/epy0n0ff/epgtimer-cli/internal/models"
	"github.com/epy0n0ff/epgtimer-cli/tests/testdata"
)

// TestGetCToken tests CSRF token retrieval from HTML page
func TestGetCToken(t *testing.T) {
	// Create mock server
	mock := testdata.NewMockEMWUIServer()
	defer mock.Close()

	// Create client
	c := client.NewClient(mock.URL())

	// Get ctok
	ctok, err := c.GetCToken()

	// Verify
	if err != nil {
		t.Fatalf("Expected success, got error: %v", err)
	}

	if ctok == "" {
		t.Error("Expected non-empty ctok, got empty string")
	}

	if ctok != mock.CToken {
		t.Errorf("Expected ctok '%s', got '%s'", mock.CToken, ctok)
	}
}

// TestGetCToken_InvalidHTML tests error when HTML doesn't contain ctok
func TestGetCToken_InvalidHTML(t *testing.T) {
	// Create HTML server without ctok
	htmlServer := testdata.NewHTMLServer()
	defer htmlServer.Close()

	// Create client
	c := client.NewClient(htmlServer.URL)

	// Try to get ctok
	_, err := c.GetCToken()

	// Verify error
	if err == nil {
		t.Fatal("Expected error for missing ctok, got nil")
	}

	if !strings.Contains(err.Error(), "ctok not found") {
		t.Errorf("Expected error message to mention 'ctok not found', got: %v", err)
	}
}

// TestSetAutoAdd_Success tests successful automatic recording rule creation
func TestSetAutoAdd_Success(t *testing.T) {
	// Create mock server
	mock := testdata.NewMockEMWUIServer()
	defer mock.Close()

	// Create client
	c := client.NewClient(mock.URL())

	// Create request
	req := models.NewAutoAddRuleRequest(
		"ニュース",                    // andKey
		"",                         // notKey
		[]string{"32736-32736-1024"}, // serviceList
	)

	// Call API
	resp, err := c.SetAutoAdd(req)

	// Verify
	if err != nil {
		t.Fatalf("Expected success, got error: %v", err)
	}

	if !resp.IsSuccess() {
		t.Errorf("Expected success=true, got success=false")
	}
}

// TestSetAutoAdd_MissingAndKey tests error when andKey is missing
func TestSetAutoAdd_MissingAndKey(t *testing.T) {
	// Create mock server
	mock := testdata.NewMockEMWUIServer()
	defer mock.Close()

	// Create client
	c := client.NewClient(mock.URL())

	// Create request with empty andKey
	req := models.NewAutoAddRuleRequest(
		"",                         // andKey (empty)
		"",                         // notKey
		[]string{"32736-32736-1024"}, // serviceList
	)

	// Call API
	_, err := c.SetAutoAdd(req)

	// Verify error
	if err == nil {
		t.Fatal("Expected validation error for missing andKey, got nil")
	}

	if !strings.Contains(err.Error(), "andKey") {
		t.Errorf("Expected error message to mention 'andKey', got: %v", err)
	}
}

// TestSetAutoAdd_MissingServiceList tests error when serviceList is missing
func TestSetAutoAdd_MissingServiceList(t *testing.T) {
	// Create mock server
	mock := testdata.NewMockEMWUIServer()
	defer mock.Close()

	// Create client
	c := client.NewClient(mock.URL())

	// Create request with empty serviceList
	req := models.NewAutoAddRuleRequest(
		"ニュース", // andKey
		"",      // notKey
		[]string{}, // serviceList (empty)
	)

	// Call API
	_, err := c.SetAutoAdd(req)

	// Verify error
	if err == nil {
		t.Fatal("Expected validation error for missing serviceList, got nil")
	}

	if !strings.Contains(err.Error(), "serviceList") {
		t.Errorf("Expected error message to mention 'serviceList', got: %v", err)
	}
}

// TestSetAutoAdd_InvalidServiceListFormat tests error when serviceList has invalid format
func TestSetAutoAdd_InvalidServiceListFormat(t *testing.T) {
	// Create mock server
	mock := testdata.NewMockEMWUIServer()
	defer mock.Close()

	// Create client
	c := client.NewClient(mock.URL())

	// Create request with invalid serviceList format
	req := models.NewAutoAddRuleRequest(
		"ニュース",           // andKey
		"",               // notKey
		[]string{"invalid"}, // serviceList (invalid format)
	)

	// Call API
	_, err := c.SetAutoAdd(req)

	// Verify error
	if err == nil {
		t.Fatal("Expected validation error for invalid serviceList format, got nil")
	}

	if !strings.Contains(err.Error(), "ONID-TSID-SID") {
		t.Errorf("Expected error message to mention format 'ONID-TSID-SID', got: %v", err)
	}
}

// TestSetAutoAdd_ConnectionError tests error when connection fails
func TestSetAutoAdd_ConnectionError(t *testing.T) {
	// Create client with invalid endpoint (server not running)
	c := client.NewClient("http://localhost:59999")

	// Create request
	req := models.NewAutoAddRuleRequest(
		"ニュース",                    // andKey
		"",                         // notKey
		[]string{"32736-32736-1024"}, // serviceList
	)

	// Call API
	_, err := c.SetAutoAdd(req)

	// Verify error
	if err == nil {
		t.Fatal("Expected connection error, got nil")
	}

	if !strings.Contains(err.Error(), "failed") {
		t.Errorf("Expected error message to indicate failure, got: %v", err)
	}
}

// TestSetAutoAdd_JapaneseEncoding tests Japanese keyword encoding
func TestSetAutoAdd_JapaneseEncoding(t *testing.T) {
	// Create mock server with custom handler to verify encoding
	mock := testdata.NewMockEMWUIServer()
	defer mock.Close()

	var receivedAndKey string
	var receivedNotKey string

	mock.SetAutoAddHandler(func(values map[string][]string) (bool, string) {
		if len(values["andKey"]) > 0 {
			receivedAndKey = values["andKey"][0]
		}
		if len(values["notKey"]) > 0 {
			receivedNotKey = values["notKey"][0]
		}
		return true, "Success"
	})

	// Create client
	c := client.NewClient(mock.URL())

	// Create request with Japanese keywords
	expectedAndKey := "ニュース番組"
	expectedNotKey := "再放送"
	req := models.NewAutoAddRuleRequest(
		expectedAndKey,             // andKey (Japanese)
		expectedNotKey,             // notKey (Japanese)
		[]string{"32736-32736-1024"}, // serviceList
	)

	// Call API
	resp, err := c.SetAutoAdd(req)

	// Verify
	if err != nil {
		t.Fatalf("Expected success, got error: %v", err)
	}

	if !resp.IsSuccess() {
		t.Errorf("Expected success=true, got success=false")
	}

	// Verify Japanese characters were received correctly
	if receivedAndKey != expectedAndKey {
		t.Errorf("Expected andKey '%s', got '%s'", expectedAndKey, receivedAndKey)
	}

	if receivedNotKey != expectedNotKey {
		t.Errorf("Expected notKey '%s', got '%s'", expectedNotKey, receivedNotKey)
	}
}

// TestSetAutoAdd_HTMLResponse tests error when server returns HTML instead of JSON
func TestSetAutoAdd_HTMLResponse(t *testing.T) {
	// Create HTML server (simulates wrong endpoint)
	htmlServer := testdata.NewHTMLServer()
	defer htmlServer.Close()

	// Create client
	c := client.NewClient(htmlServer.URL)

	// Create request
	req := models.NewAutoAddRuleRequest(
		"ニュース",                    // andKey
		"",                         // notKey
		[]string{"32736-32736-1024"}, // serviceList
	)

	// Call API
	_, err := c.SetAutoAdd(req)

	// Verify error
	if err == nil {
		t.Fatal("Expected error for HTML response, got nil")
	}

	if !strings.Contains(err.Error(), "HTML") {
		t.Errorf("Expected error message to mention 'HTML', got: %v", err)
	}
}

// TestSetAutoAdd_APIError tests error when API returns error response
func TestSetAutoAdd_APIError(t *testing.T) {
	// Create failing server
	mock := testdata.NewFailingServer()
	defer mock.Close()

	// Create client
	c := client.NewClient(mock.URL())

	// Create request
	req := models.NewAutoAddRuleRequest(
		"ニュース",                    // andKey
		"",                         // notKey
		[]string{"32736-32736-1024"}, // serviceList
	)

	// Call API
	_, err := c.SetAutoAdd(req)

	// Verify error
	if err == nil {
		t.Fatal("Expected API error, got nil")
	}

	if !strings.Contains(err.Error(), "error") {
		t.Errorf("Expected error message to mention 'error', got: %v", err)
	}
}

// TestSetAutoAdd_MultipleChannels tests successful creation with multiple channels
func TestSetAutoAdd_MultipleChannels(t *testing.T) {
	// Create mock server with custom handler
	mock := testdata.NewMockEMWUIServer()
	defer mock.Close()

	var receivedServiceList []string

	mock.SetAutoAddHandler(func(values map[string][]string) (bool, string) {
		if services, ok := values["serviceList"]; ok {
			receivedServiceList = services
		}
		return true, "Success"
	})

	// Create client
	c := client.NewClient(mock.URL())

	// Create request with multiple channels
	expectedChannels := []string{
		"32736-32736-1024",
		"32736-32736-1025",
		"32736-32736-1026",
	}
	req := models.NewAutoAddRuleRequest(
		"ニュース",         // andKey
		"",             // notKey
		expectedChannels, // serviceList (multiple)
	)

	// Call API
	resp, err := c.SetAutoAdd(req)

	// Verify
	if err != nil {
		t.Fatalf("Expected success, got error: %v", err)
	}

	if !resp.IsSuccess() {
		t.Errorf("Expected success=true, got success=false")
	}

	// Verify all channels were received (skip empty entries)
	var nonEmptyChannels []string
	for _, ch := range receivedServiceList {
		if ch != "" {
			nonEmptyChannels = append(nonEmptyChannels, ch)
		}
	}

	if len(nonEmptyChannels) != len(expectedChannels) {
		t.Errorf("Expected %d channels, got %d", len(expectedChannels), len(nonEmptyChannels))
	}

	for i, expected := range expectedChannels {
		if i >= len(nonEmptyChannels) || nonEmptyChannels[i] != expected {
			t.Errorf("Expected channel[%d] '%s', got '%s'", i, expected, nonEmptyChannels[i])
		}
	}
}

// TestSetAutoAdd_WithNotKey tests successful creation with exclusion keywords
func TestSetAutoAdd_WithNotKey(t *testing.T) {
	// Create mock server
	mock := testdata.NewMockEMWUIServer()
	defer mock.Close()

	var receivedNotKey string

	mock.SetAutoAddHandler(func(values map[string][]string) (bool, string) {
		if len(values["notKey"]) > 0 {
			receivedNotKey = values["notKey"][0]
		}
		return true, "Success"
	})

	// Create client
	c := client.NewClient(mock.URL())

	// Create request with notKey
	expectedNotKey := "再放送"
	req := models.NewAutoAddRuleRequest(
		"ドラマ",                      // andKey
		expectedNotKey,             // notKey
		[]string{"32736-32736-1024"}, // serviceList
	)

	// Call API
	resp, err := c.SetAutoAdd(req)

	// Verify
	if err != nil {
		t.Fatalf("Expected success, got error: %v", err)
	}

	if !resp.IsSuccess() {
		t.Errorf("Expected success=true, got success=false")
	}

	// Verify notKey was received
	if receivedNotKey != expectedNotKey {
		t.Errorf("Expected notKey '%s', got '%s'", expectedNotKey, receivedNotKey)
	}
}

// TestParseServiceListEntry tests the channel format parser
func TestParseServiceListEntry(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
		wantONID int
		wantTSID int
		wantSID  int
	}{
		{
			name:    "Valid NHK format",
			input:   "32736-32736-1024",
			wantErr: false,
			wantONID: 32736,
			wantTSID: 32736,
			wantSID:  1024,
		},
		{
			name:    "Invalid format - only two parts",
			input:   "32736-32736",
			wantErr: true,
		},
		{
			name:    "Invalid format - four parts",
			input:   "32736-32736-1024-extra",
			wantErr: true,
		},
		{
			name:    "Invalid format - non-numeric ONID",
			input:   "abc-32736-1024",
			wantErr: true,
		},
		{
			name:    "Invalid format - empty",
			input:   "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			entry, err := models.ParseServiceListEntry(tt.input)

			if tt.wantErr {
				if err == nil {
					t.Errorf("Expected error for input '%s', got nil", tt.input)
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error for input '%s', got: %v", tt.input, err)
				}
				if entry.ONID != tt.wantONID {
					t.Errorf("Expected ONID %d, got %d", tt.wantONID, entry.ONID)
				}
				if entry.TSID != tt.wantTSID {
					t.Errorf("Expected TSID %d, got %d", tt.wantTSID, entry.TSID)
				}
				if entry.SID != tt.wantSID {
					t.Errorf("Expected SID %d, got %d", tt.wantSID, entry.SID)
				}
			}
		})
	}
}

// TestAutoAddRuleRequest_ToFormData tests form data encoding
func TestAutoAddRuleRequest_ToFormData(t *testing.T) {
	req := models.NewAutoAddRuleRequest(
		"ニュース",                    // andKey
		"再放送",                      // notKey
		[]string{"32736-32736-1024"}, // serviceList
	)

	formData := req.ToFormData()

	// Verify form data contains required fields
	if !strings.Contains(formData, "andKey=") {
		t.Error("Form data missing andKey")
	}
	if !strings.Contains(formData, "notKey=") {
		t.Error("Form data missing notKey")
	}
	if !strings.Contains(formData, "serviceList=") {
		t.Error("Form data missing serviceList")
	}
	if !strings.Contains(formData, "addchg=1") {
		t.Error("Form data missing addchg default")
	}
	if !strings.Contains(formData, "ctok=") {
		t.Error("Form data missing ctok")
	}
}
