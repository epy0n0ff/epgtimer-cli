package integration

import (
	"strings"
	"testing"

	"github.com/epy0n0ff/epgtimer-cli/internal/client"
	"github.com/epy0n0ff/epgtimer-cli/tests/testdata"
)

// TestEnumEventInfo_Success tests successful EPG retrieval
func TestEnumEventInfo_Success(t *testing.T) {
	// Create mock server
	mock := testdata.NewMockEMWUIServer()
	defer mock.Close()

	// Create client
	apiClient := client.NewClient(mock.URL())

	// Call EnumEventInfo
	response, err := apiClient.EnumEventInfo(32736, 32736, 1024)
	if err != nil {
		t.Fatalf("EnumEventInfo() failed: %v", err)
	}

	// Verify response structure
	if response.Total != 3 {
		t.Errorf("Expected Total=3, got %d", response.Total)
	}

	if response.Count != 3 {
		t.Errorf("Expected Count=3, got %d", response.Count)
	}

	if len(response.Items) != 3 {
		t.Fatalf("Expected 3 events, got %d", len(response.Items))
	}

	// Verify first event
	event1 := response.Items[0]
	if event1.ONID != 32736 {
		t.Errorf("Expected ONID=32736, got %d", event1.ONID)
	}

	if event1.TSID != 32736 {
		t.Errorf("Expected TSID=32736, got %d", event1.TSID)
	}

	if event1.SID != 1024 {
		t.Errorf("Expected SID=1024, got %d", event1.SID)
	}

	if event1.EventName != "ＮＨＫニュース　おはよう日本　テスト" {
		t.Errorf("Expected EventName='ＮＨＫニュース　おはよう日本　テスト', got '%s'", event1.EventName)
	}

	if event1.StartDate != "2025/12/23" {
		t.Errorf("Expected StartDate='2025/12/23', got '%s'", event1.StartDate)
	}

	if event1.StartTime != "06:00:00" {
		t.Errorf("Expected StartTime='06:00:00', got '%s'", event1.StartTime)
	}

	if event1.Duration != 1800 {
		t.Errorf("Expected Duration=1800, got %d", event1.Duration)
	}

	if event1.DurationMinutes() != 30 {
		t.Errorf("Expected DurationMinutes=30, got %d", event1.DurationMinutes())
	}

	if event1.ChannelID() != "32736-32736-1024" {
		t.Errorf("Expected ChannelID='32736-32736-1024', got '%s'", event1.ChannelID())
	}

	if !event1.IsFreeCA() {
		t.Error("Expected event1 to be free (not scrambled)")
	}

	// Verify second event (drama)
	event2 := response.Items[1]
	if event2.EventName != "【連続テレビ小説】ばけばけ（６２）" {
		t.Errorf("Expected EventName='【連続テレビ小説】ばけばけ（６２）', got '%s'", event2.EventName)
	}

	if event2.DurationMinutes() != 15 {
		t.Errorf("Expected DurationMinutes=15, got %d", event2.DurationMinutes())
	}

	// Verify third event
	event3 := response.Items[2]
	if event3.DurationMinutes() != 100 {
		t.Errorf("Expected DurationMinutes=100, got %d", event3.DurationMinutes())
	}
}

// TestEnumEventInfo_Empty tests handling of empty response
func TestEnumEventInfo_Empty(t *testing.T) {
	// Create mock server configured for empty response
	mock := testdata.NewEmptyEnumEventInfoServer()
	defer mock.Close()

	// Create client
	apiClient := client.NewClient(mock.URL())

	// Call EnumEventInfo
	response, err := apiClient.EnumEventInfo(32736, 32736, 1024)
	if err != nil {
		t.Fatalf("EnumEventInfo() failed: %v", err)
	}

	// Verify empty response
	if response.Total != 0 {
		t.Errorf("Expected Total=0, got %d", response.Total)
	}

	if response.Count != 0 {
		t.Errorf("Expected Count=0, got %d", response.Count)
	}

	if len(response.Items) != 0 {
		t.Errorf("Expected 0 events, got %d", len(response.Items))
	}
}

// TestEnumEventInfo_ConnectionError tests connection error handling
func TestEnumEventInfo_ConnectionError(t *testing.T) {
	// Create client with invalid endpoint
	apiClient := client.NewClient("http://localhost:99999")

	// Call EnumEventInfo (should fail)
	_, err := apiClient.EnumEventInfo(32736, 32736, 1024)
	if err == nil {
		t.Fatal("Expected connection error, got nil")
	}

	// Verify error message contains helpful context
	errMsg := err.Error()
	if !strings.Contains(errMsg, "failed to connect") {
		t.Errorf("Error message should mention connection failure, got: %s", errMsg)
	}
}

// TestEnumEventInfo_MalformedXML tests handling of malformed XML response
func TestEnumEventInfo_MalformedXML(t *testing.T) {
	// Create mock server that returns invalid XML
	mock := testdata.NewMockEMWUIServer()
	mock.SetEnumEventInfoHandler(func() (string, int) {
		return "<invalid><xml>", 200
	})
	defer mock.Close()

	// Create client
	apiClient := client.NewClient(mock.URL())

	// Call EnumEventInfo (should fail to parse)
	_, err := apiClient.EnumEventInfo(32736, 32736, 1024)
	if err == nil {
		t.Fatal("Expected XML parse error, got nil")
	}

	// Verify error message mentions parsing failure
	errMsg := err.Error()
	if !strings.Contains(errMsg, "failed to parse XML") {
		t.Errorf("Error message should mention XML parsing failure, got: %s", errMsg)
	}
}

// TestEnumEventInfo_JapaneseEncoding tests proper handling of Japanese characters
func TestEnumEventInfo_JapaneseEncoding(t *testing.T) {
	// Create mock server
	mock := testdata.NewMockEMWUIServer()
	defer mock.Close()

	// Create client
	apiClient := client.NewClient(mock.URL())

	// Call EnumEventInfo
	response, err := apiClient.EnumEventInfo(32736, 32736, 1024)
	if err != nil {
		t.Fatalf("EnumEventInfo() failed: %v", err)
	}

	// Verify Japanese characters are preserved correctly
	japaneseTitles := []string{
		"ＮＨＫニュース　おはよう日本　テスト",
		"【連続テレビ小説】ばけばけ（６２）",
		"あさイチ　テスト番組",
	}

	for i, expected := range japaneseTitles {
		if i >= len(response.Items) {
			break
		}
		actual := response.Items[i].EventName
		if actual != expected {
			t.Errorf("Event %d: Japanese characters not preserved correctly. Expected '%s', got '%s'",
				i+1, expected, actual)
		}
	}
}

// TestEventInfo_GenreString tests the GenreString() method
func TestEventInfo_GenreString(t *testing.T) {
	// Create mock server
	mock := testdata.NewMockEMWUIServer()
	defer mock.Close()

	// Create client
	apiClient := client.NewClient(mock.URL())

	// Call EnumEventInfo
	response, err := apiClient.EnumEventInfo(32736, 32736, 1024)
	if err != nil {
		t.Fatalf("EnumEventInfo() failed: %v", err)
	}

	// Test GenreString() method
	if len(response.Items) > 0 {
		event := response.Items[0]
		genre := event.GenreString()

		expectedGenre := "ニュース／報道 - 定時・総合"
		if genre != expectedGenre {
			t.Errorf("Expected GenreString='%s', got '%s'", expectedGenre, genre)
		}
	}
}
