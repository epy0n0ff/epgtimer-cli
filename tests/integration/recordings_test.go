package integration

import (
	"strings"
	"testing"

	"github.com/epy0n0ff/epgtimer-cli/internal/client"
	"github.com/epy0n0ff/epgtimer-cli/tests/testdata"
)

// TestEnumRecInfo_Success tests successful recording retrieval
func TestEnumRecInfo_Success(t *testing.T) {
	// Create mock server
	mock := testdata.NewMockEMWUIServer()
	defer mock.Close()

	// Create client
	apiClient := client.NewClient(mock.URL())

	// Call EnumRecInfo
	response, err := apiClient.EnumRecInfo()
	if err != nil {
		t.Fatalf("EnumRecInfo() failed: %v", err)
	}

	// Verify response structure
	if response.Total != 2 {
		t.Errorf("Expected Total=2, got %d", response.Total)
	}

	if response.Count != 2 {
		t.Errorf("Expected Count=2, got %d", response.Count)
	}

	if len(response.Items) != 2 {
		t.Fatalf("Expected 2 recordings, got %d", len(response.Items))
	}

	// Verify first recording
	rec1 := response.Items[0]
	if rec1.ID != 2001 {
		t.Errorf("Expected ID=2001, got %d", rec1.ID)
	}

	if rec1.Title != "サイエンスZERO　過去回" {
		t.Errorf("Expected Title='サイエンスZERO　過去回', got '%s'", rec1.Title)
	}

	if rec1.StartDate != "2025/12/15" {
		t.Errorf("Expected StartDate='2025/12/15', got '%s'", rec1.StartDate)
	}

	if rec1.StartTime != "23:30:00" {
		t.Errorf("Expected StartTime='23:30:00', got '%s'", rec1.StartTime)
	}

	if rec1.DurationSecond != 1800 {
		t.Errorf("Expected DurationSecond=1800, got %d", rec1.DurationSecond)
	}

	if rec1.DurationMinutes() != 30 {
		t.Errorf("Expected DurationMinutes=30, got %d", rec1.DurationMinutes())
	}

	if rec1.ChannelID() != "32736-32736-1024" {
		t.Errorf("Expected ChannelID='32736-32736-1024', got '%s'", rec1.ChannelID())
	}

	if rec1.RecFilePath != `C:\Recorded\サイエンスZERO_20251215_233000.ts` {
		t.Errorf("Expected RecFilePath='C:\\Recorded\\サイエンスZERO_20251215_233000.ts', got '%s'", rec1.RecFilePath)
	}

	if rec1.IsProtected() {
		t.Error("Expected rec1 to not be protected")
	}

	// Verify second recording
	rec2 := response.Items[1]
	if rec2.Title != "ブラタモリ　過去回" {
		t.Errorf("Expected Title='ブラタモリ　過去回', got '%s'", rec2.Title)
	}

	if rec2.DurationMinutes() != 45 {
		t.Errorf("Expected DurationMinutes=45, got %d", rec2.DurationMinutes())
	}

	if !rec2.IsProtected() {
		t.Error("Expected rec2 to be protected")
	}
}

// TestEnumRecInfo_Empty tests handling of empty response
func TestEnumRecInfo_Empty(t *testing.T) {
	// Create mock server configured for empty response
	mock := testdata.NewEmptyEnumRecInfoServer()
	defer mock.Close()

	// Create client
	apiClient := client.NewClient(mock.URL())

	// Call EnumRecInfo
	response, err := apiClient.EnumRecInfo()
	if err != nil {
		t.Fatalf("EnumRecInfo() failed: %v", err)
	}

	// Verify empty response
	if response.Total != 0 {
		t.Errorf("Expected Total=0, got %d", response.Total)
	}

	if response.Count != 0 {
		t.Errorf("Expected Count=0, got %d", response.Count)
	}

	if len(response.Items) != 0 {
		t.Errorf("Expected 0 recordings, got %d", len(response.Items))
	}
}

// TestEnumRecInfo_ConnectionError tests connection error handling
func TestEnumRecInfo_ConnectionError(t *testing.T) {
	// Create client with invalid endpoint
	apiClient := client.NewClient("http://localhost:99999")

	// Call EnumRecInfo (should fail)
	_, err := apiClient.EnumRecInfo()
	if err == nil {
		t.Fatal("Expected connection error, got nil")
	}

	// Verify error message contains helpful context
	errMsg := err.Error()
	if !strings.Contains(errMsg, "failed to connect") {
		t.Errorf("Error message should mention connection failure, got: %s", errMsg)
	}
}

// TestEnumRecInfo_MalformedXML tests handling of malformed XML response
func TestEnumRecInfo_MalformedXML(t *testing.T) {
	// Create mock server that returns invalid XML
	mock := testdata.NewMockEMWUIServer()
	mock.SetEnumRecInfoHandler(func() (string, int) {
		return "<invalid><xml>", 200
	})
	defer mock.Close()

	// Create client
	apiClient := client.NewClient(mock.URL())

	// Call EnumRecInfo (should fail to parse)
	_, err := apiClient.EnumRecInfo()
	if err == nil {
		t.Fatal("Expected XML parse error, got nil")
	}

	// Verify error message mentions parsing failure
	errMsg := err.Error()
	if !strings.Contains(errMsg, "failed to parse XML") {
		t.Errorf("Error message should mention XML parsing failure, got: %s", errMsg)
	}
}

// TestEnumRecInfo_JapaneseEncoding tests proper handling of Japanese characters
func TestEnumRecInfo_JapaneseEncoding(t *testing.T) {
	// Create mock server
	mock := testdata.NewMockEMWUIServer()
	defer mock.Close()

	// Create client
	apiClient := client.NewClient(mock.URL())

	// Call EnumRecInfo
	response, err := apiClient.EnumRecInfo()
	if err != nil {
		t.Fatalf("EnumRecInfo() failed: %v", err)
	}

	// Verify Japanese characters are preserved correctly
	japaneseTitles := []string{"サイエンスZERO　過去回", "ブラタモリ　過去回"}

	for i, expected := range japaneseTitles {
		if i >= len(response.Items) {
			break
		}
		actual := response.Items[i].Title
		if actual != expected {
			t.Errorf("Recording %d: Japanese characters not preserved correctly. Expected '%s', got '%s'",
				i+1, expected, actual)
		}
	}

	// Verify Japanese characters in file path
	if len(response.Items) > 0 {
		filePath := response.Items[0].RecFilePath
		if !strings.Contains(filePath, "サイエンスZERO") {
			t.Errorf("Expected file path to contain 'サイエンスZERO', got '%s'", filePath)
		}
	}
}

// TestRecordingInfo_IsProtected tests the IsProtected() method
func TestRecordingInfo_IsProtected(t *testing.T) {
	// Create mock server
	mock := testdata.NewMockEMWUIServer()
	defer mock.Close()

	// Create client
	apiClient := client.NewClient(mock.URL())

	// Call EnumRecInfo
	response, err := apiClient.EnumRecInfo()
	if err != nil {
		t.Fatalf("EnumRecInfo() failed: %v", err)
	}

	// Test IsProtected() method
	if len(response.Items) >= 2 {
		// First recording should not be protected
		if response.Items[0].IsProtected() {
			t.Error("Expected first recording to not be protected")
		}

		// Second recording should be protected
		if !response.Items[1].IsProtected() {
			t.Error("Expected second recording to be protected")
		}
	}
}
