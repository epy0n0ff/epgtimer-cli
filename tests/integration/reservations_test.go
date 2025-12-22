package integration

import (
	"strings"
	"testing"

	"github.com/epy0n0ff/epgtimer-cli/internal/client"
	"github.com/epy0n0ff/epgtimer-cli/tests/testdata"
)

// TestEnumReserveInfo_Success tests successful reservation retrieval
func TestEnumReserveInfo_Success(t *testing.T) {
	// Create mock server
	mock := testdata.NewMockEMWUIServer()
	defer mock.Close()

	// Create client
	apiClient := client.NewClient(mock.URL())

	// Call EnumReserveInfo
	response, err := apiClient.EnumReserveInfo()
	if err != nil {
		t.Fatalf("EnumReserveInfo() failed: %v", err)
	}

	// Verify response structure
	if response.Total != 2 {
		t.Errorf("Expected Total=2, got %d", response.Total)
	}

	if response.Count != 2 {
		t.Errorf("Expected Count=2, got %d", response.Count)
	}

	if len(response.Items) != 2 {
		t.Fatalf("Expected 2 reservations, got %d", len(response.Items))
	}

	// Verify first reservation
	res1 := response.Items[0]
	if res1.ID != 1001 {
		t.Errorf("Expected ID=1001, got %d", res1.ID)
	}

	if res1.Title != "サイエンスZERO　テスト放送" {
		t.Errorf("Expected Title='サイエンスZERO　テスト放送', got '%s'", res1.Title)
	}

	if res1.StartDate != "2025/12/22" {
		t.Errorf("Expected StartDate='2025/12/22', got '%s'", res1.StartDate)
	}

	if res1.StartTime != "23:30:00" {
		t.Errorf("Expected StartTime='23:30:00', got '%s'", res1.StartTime)
	}

	if res1.DurationSecond != 1800 {
		t.Errorf("Expected DurationSecond=1800, got %d", res1.DurationSecond)
	}

	if res1.DurationMinutes() != 30 {
		t.Errorf("Expected DurationMinutes=30, got %d", res1.DurationMinutes())
	}

	if res1.ChannelID() != "32736-32736-1024" {
		t.Errorf("Expected ChannelID='32736-32736-1024', got '%s'", res1.ChannelID())
	}

	// Verify recording settings
	if res1.RecSetting.Priority != 5 {
		t.Errorf("Expected Priority=5, got %d", res1.RecSetting.Priority)
	}

	if res1.RecSetting.TuijyuuFlag != 1 {
		t.Errorf("Expected TuijyuuFlag=1, got %d", res1.RecSetting.TuijyuuFlag)
	}

	// Verify second reservation
	res2 := response.Items[1]
	if res2.Title != "ブラタモリ　テスト回" {
		t.Errorf("Expected Title='ブラタモリ　テスト回', got '%s'", res2.Title)
	}

	if res2.DurationMinutes() != 45 {
		t.Errorf("Expected DurationMinutes=45, got %d", res2.DurationMinutes())
	}
}

// TestEnumReserveInfo_Empty tests handling of empty response
func TestEnumReserveInfo_Empty(t *testing.T) {
	// Create mock server configured for empty response
	mock := testdata.NewEmptyEnumReserveInfoServer()
	defer mock.Close()

	// Create client
	apiClient := client.NewClient(mock.URL())

	// Call EnumReserveInfo
	response, err := apiClient.EnumReserveInfo()
	if err != nil {
		t.Fatalf("EnumReserveInfo() failed: %v", err)
	}

	// Verify empty response
	if response.Total != 0 {
		t.Errorf("Expected Total=0, got %d", response.Total)
	}

	if response.Count != 0 {
		t.Errorf("Expected Count=0, got %d", response.Count)
	}

	if len(response.Items) != 0 {
		t.Errorf("Expected 0 reservations, got %d", len(response.Items))
	}
}

// TestEnumReserveInfo_ConnectionError tests connection error handling
func TestEnumReserveInfo_ConnectionError(t *testing.T) {
	// Create client with invalid endpoint
	apiClient := client.NewClient("http://localhost:99999")

	// Call EnumReserveInfo (should fail)
	_, err := apiClient.EnumReserveInfo()
	if err == nil {
		t.Fatal("Expected connection error, got nil")
	}

	// Verify error message contains helpful context
	errMsg := err.Error()
	if !strings.Contains(errMsg, "failed to connect") {
		t.Errorf("Error message should mention connection failure, got: %s", errMsg)
	}
}

// TestEnumReserveInfo_MalformedXML tests handling of malformed XML response
func TestEnumReserveInfo_MalformedXML(t *testing.T) {
	// Create mock server that returns invalid XML
	mock := testdata.NewMockEMWUIServer()
	mock.SetEnumReserveInfoHandler(func() (string, int) {
		return "<invalid><xml>", 200
	})
	defer mock.Close()

	// Create client
	apiClient := client.NewClient(mock.URL())

	// Call EnumReserveInfo (should fail to parse)
	_, err := apiClient.EnumReserveInfo()
	if err == nil {
		t.Fatal("Expected XML parse error, got nil")
	}

	// Verify error message mentions parsing failure
	errMsg := err.Error()
	if !strings.Contains(errMsg, "failed to parse XML") {
		t.Errorf("Error message should mention XML parsing failure, got: %s", errMsg)
	}
}

// TestEnumReserveInfo_JapaneseEncoding tests proper handling of Japanese characters
func TestEnumReserveInfo_JapaneseEncoding(t *testing.T) {
	// Create mock server
	mock := testdata.NewMockEMWUIServer()
	defer mock.Close()

	// Create client
	apiClient := client.NewClient(mock.URL())

	// Call EnumReserveInfo
	response, err := apiClient.EnumReserveInfo()
	if err != nil {
		t.Fatalf("EnumReserveInfo() failed: %v", err)
	}

	// Verify Japanese characters are preserved correctly
	japaneseTitles := []string{"サイエンスZERO　テスト放送", "ブラタモリ　テスト回"}

	for i, expected := range japaneseTitles {
		if i >= len(response.Items) {
			break
		}
		actual := response.Items[i].Title
		if actual != expected {
			t.Errorf("Reservation %d: Japanese characters not preserved correctly. Expected '%s', got '%s'",
				i+1, expected, actual)
		}
	}
}

// TestReservationInfo_RecModeString tests the RecModeString() method
func TestReservationInfo_RecModeString(t *testing.T) {
	// Create mock server
	mock := testdata.NewMockEMWUIServer()
	defer mock.Close()

	// Create client
	apiClient := client.NewClient(mock.URL())

	// Call EnumReserveInfo
	response, err := apiClient.EnumReserveInfo()
	if err != nil {
		t.Fatalf("EnumReserveInfo() failed: %v", err)
	}

	// Test RecModeString() formatting
	if len(response.Items) > 0 {
		reservation := response.Items[0]
		recMode := reservation.RecModeString()

		expectedMode := "All"
		if recMode != expectedMode {
			t.Errorf("Expected RecModeString='%s', got '%s'", expectedMode, recMode)
		}
	}
}
