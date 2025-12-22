package integration

import (
	"strings"
	"testing"

	"github.com/epy0n0ff/epgtimer-cli/internal/client"
	"github.com/epy0n0ff/epgtimer-cli/tests/testdata"
)

// TestEnumService_Success tests successful channel retrieval
func TestEnumService_Success(t *testing.T) {
	// Create mock server
	mock := testdata.NewMockEMWUIServer()
	defer mock.Close()

	// Create client
	apiClient := client.NewClient(mock.URL())

	// Call EnumService
	response, err := apiClient.EnumService()
	if err != nil {
		t.Fatalf("EnumService() failed: %v", err)
	}

	// Verify response structure
	if response.Total != 3 {
		t.Errorf("Expected Total=3, got %d", response.Total)
	}

	if response.Count != 3 {
		t.Errorf("Expected Count=3, got %d", response.Count)
	}

	if len(response.Items) != 3 {
		t.Fatalf("Expected 3 channels, got %d", len(response.Items))
	}

	// Verify first channel (NHK総合)
	ch1 := response.Items[0]
	if ch1.ONID != 32736 {
		t.Errorf("Expected ONID=32736, got %d", ch1.ONID)
	}

	if ch1.TSID != 32736 {
		t.Errorf("Expected TSID=32736, got %d", ch1.TSID)
	}

	if ch1.SID != 1024 {
		t.Errorf("Expected SID=1024, got %d", ch1.SID)
	}

	if ch1.ServiceName != "NHK総合・東京" {
		t.Errorf("Expected ServiceName='NHK総合・東京', got '%s'", ch1.ServiceName)
	}

	if !ch1.IsTV() {
		t.Error("Expected ch1 to be a TV channel")
	}

	if ch1.ServiceTypeString() != "TV" {
		t.Errorf("Expected ServiceTypeString='TV', got '%s'", ch1.ServiceTypeString())
	}

	if ch1.ChannelID() != "32736-32736-1024" {
		t.Errorf("Expected ChannelID='32736-32736-1024', got '%s'", ch1.ChannelID())
	}

	if ch1.RemoteControlKeyID != 1 {
		t.Errorf("Expected RemoteControlKeyID=1, got %d", ch1.RemoteControlKeyID)
	}

	// Verify second channel (BS朝日)
	ch2 := response.Items[1]
	if ch2.ServiceName != "BS朝日" {
		t.Errorf("Expected ServiceName='BS朝日', got '%s'", ch2.ServiceName)
	}

	if ch2.NetworkName != "BS Digital" {
		t.Errorf("Expected NetworkName='BS Digital', got '%s'", ch2.NetworkName)
	}

	// Verify third channel (NHKラジオ第1 - Radio)
	ch3 := response.Items[2]
	if !ch3.IsRadio() {
		t.Error("Expected ch3 to be a radio channel")
	}

	if ch3.ServiceTypeString() != "Radio" {
		t.Errorf("Expected ServiceTypeString='Radio', got '%s'", ch3.ServiceTypeString())
	}
}

// TestEnumService_Empty tests handling of empty response
func TestEnumService_Empty(t *testing.T) {
	// Create mock server configured for empty response
	mock := testdata.NewEmptyEnumServiceServer()
	defer mock.Close()

	// Create client
	apiClient := client.NewClient(mock.URL())

	// Call EnumService
	response, err := apiClient.EnumService()
	if err != nil {
		t.Fatalf("EnumService() failed: %v", err)
	}

	// Verify empty response
	if response.Total != 0 {
		t.Errorf("Expected Total=0, got %d", response.Total)
	}

	if response.Count != 0 {
		t.Errorf("Expected Count=0, got %d", response.Count)
	}

	if len(response.Items) != 0 {
		t.Errorf("Expected 0 channels, got %d", len(response.Items))
	}
}

// TestEnumService_ConnectionError tests connection error handling
func TestEnumService_ConnectionError(t *testing.T) {
	// Create client with invalid endpoint
	apiClient := client.NewClient("http://localhost:99999")

	// Call EnumService (should fail)
	_, err := apiClient.EnumService()
	if err == nil {
		t.Fatal("Expected connection error, got nil")
	}

	// Verify error message contains helpful context
	errMsg := err.Error()
	if !strings.Contains(errMsg, "failed to connect") {
		t.Errorf("Error message should mention connection failure, got: %s", errMsg)
	}
}

// TestEnumService_MalformedXML tests handling of malformed XML response
func TestEnumService_MalformedXML(t *testing.T) {
	// Create mock server that returns invalid XML
	mock := testdata.NewMockEMWUIServer()
	mock.SetEnumServiceHandler(func() (string, int) {
		return "<invalid><xml>", 200
	})
	defer mock.Close()

	// Create client
	apiClient := client.NewClient(mock.URL())

	// Call EnumService (should fail to parse)
	_, err := apiClient.EnumService()
	if err == nil {
		t.Fatal("Expected XML parse error, got nil")
	}

	// Verify error message mentions parsing failure
	errMsg := err.Error()
	if !strings.Contains(errMsg, "failed to parse XML") {
		t.Errorf("Error message should mention XML parsing failure, got: %s", errMsg)
	}
}

// TestEnumService_JapaneseEncoding tests proper handling of Japanese characters
func TestEnumService_JapaneseEncoding(t *testing.T) {
	// Create mock server
	mock := testdata.NewMockEMWUIServer()
	defer mock.Close()

	// Create client
	apiClient := client.NewClient(mock.URL())

	// Call EnumService
	response, err := apiClient.EnumService()
	if err != nil {
		t.Fatalf("EnumService() failed: %v", err)
	}

	// Verify Japanese characters are preserved correctly
	japaneseNames := []string{"NHK総合・東京", "BS朝日", "NHKラジオ第1"}

	for i, expected := range japaneseNames {
		if i >= len(response.Items) {
			break
		}
		actual := response.Items[i].ServiceName
		if actual != expected {
			t.Errorf("Channel %d: Japanese characters not preserved correctly. Expected '%s', got '%s'",
				i+1, expected, actual)
		}
	}
}

// TestChannelInfo_ChannelID tests the ChannelID() method
func TestChannelInfo_ChannelID(t *testing.T) {
	// Create mock server
	mock := testdata.NewMockEMWUIServer()
	defer mock.Close()

	// Create client
	apiClient := client.NewClient(mock.URL())

	// Call EnumService
	response, err := apiClient.EnumService()
	if err != nil {
		t.Fatalf("EnumService() failed: %v", err)
	}

	// Test ChannelID() formatting
	if len(response.Items) > 0 {
		channel := response.Items[0]
		channelID := channel.ChannelID()

		expectedFormat := "32736-32736-1024"
		if channelID != expectedFormat {
			t.Errorf("Expected channel ID '%s', got '%s'", expectedFormat, channelID)
		}
	}
}
