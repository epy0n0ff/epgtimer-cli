package integration

import (
	"strings"
	"testing"

	"github.com/epy0n0ff/epgtimer-cli/internal/client"
	"github.com/epy0n0ff/epgtimer-cli/tests/testdata"
)

// TestDeleteAutoAdd_Success tests successful deletion of automatic recording rule
func TestDeleteAutoAdd_Success(t *testing.T) {
	// Create mock server
	mock := testdata.NewMockEMWUIServer()
	defer mock.Close()

	// Create client
	apiClient := client.NewClient(mock.URL())

	// Call DeleteAutoAdd with valid ID
	response, err := apiClient.DeleteAutoAdd(334)
	if err != nil {
		t.Fatalf("DeleteAutoAdd() failed: %v", err)
	}

	// Verify response
	if !response.IsSuccess() {
		t.Errorf("Expected success response, got error: %s", response.GetError())
	}

	message := response.GetMessage()
	if message == "" {
		t.Error("Expected success message, got empty string")
	}
}

// TestDeleteAutoAdd_InvalidID tests deletion with invalid ID (0)
func TestDeleteAutoAdd_InvalidID(t *testing.T) {
	// Create mock server
	mock := testdata.NewMockEMWUIServer()
	defer mock.Close()

	// Create client
	apiClient := client.NewClient(mock.URL())

	// Call DeleteAutoAdd with invalid ID (0)
	_, err := apiClient.DeleteAutoAdd(0)
	if err == nil {
		t.Fatal("Expected error for invalid ID, got nil")
	}

	// Verify error message
	errMsg := err.Error()
	if !strings.Contains(errMsg, "invalid rule ID") {
		t.Errorf("Error message should mention invalid ID, got: %s", errMsg)
	}
}

// TestDeleteAutoAdd_NegativeID tests deletion with negative ID
func TestDeleteAutoAdd_NegativeID(t *testing.T) {
	// Create mock server
	mock := testdata.NewMockEMWUIServer()
	defer mock.Close()

	// Create client
	apiClient := client.NewClient(mock.URL())

	// Call DeleteAutoAdd with negative ID
	_, err := apiClient.DeleteAutoAdd(-1)
	if err == nil {
		t.Fatal("Expected error for negative ID, got nil")
	}

	// Verify error message
	errMsg := err.Error()
	if !strings.Contains(errMsg, "invalid rule ID") {
		t.Errorf("Error message should mention invalid ID, got: %s", errMsg)
	}
}

// TestDeleteAutoAdd_ServerError tests handling of server error response
func TestDeleteAutoAdd_ServerError(t *testing.T) {
	// Create mock server with custom delete handler that returns error
	mock := testdata.NewMockEMWUIServer()
	mock.SetDeleteAutoAddHandler(func(id int) (bool, string) {
		return false, "Rule not found"
	})
	defer mock.Close()

	// Create client
	apiClient := client.NewClient(mock.URL())

	// Call DeleteAutoAdd
	_, err := apiClient.DeleteAutoAdd(999)
	if err == nil {
		t.Fatal("Expected error from server, got nil")
	}

	// Verify error message contains server error
	errMsg := err.Error()
	if !strings.Contains(errMsg, "Rule not found") {
		t.Errorf("Error message should contain server error, got: %s", errMsg)
	}
}

// TestDeleteAutoAdd_ConnectionError tests connection error handling
func TestDeleteAutoAdd_ConnectionError(t *testing.T) {
	// Create client with invalid endpoint
	apiClient := client.NewClient("http://localhost:99999")

	// Call DeleteAutoAdd (should fail)
	_, err := apiClient.DeleteAutoAdd(334)
	if err == nil {
		t.Fatal("Expected connection error, got nil")
	}

	// Verify error message contains helpful context
	errMsg := err.Error()
	if !strings.Contains(errMsg, "failed to") {
		t.Errorf("Error message should mention failure, got: %s", errMsg)
	}
}

// TestDeleteAutoAdd_MalformedXML tests handling of malformed XML response
func TestDeleteAutoAdd_MalformedXML(t *testing.T) {
	// Create mock server that returns invalid XML for delete
	mock := testdata.NewMockEMWUIServer()
	mock.SetDeleteAutoAddHandler(func(id int) (bool, string) {
		// This won't be used; we'll override the response in the handler
		return true, "success"
	})
	// Override with custom handler that returns malformed XML
	mock.OnDeleteAutoAdd = nil
	mock.SetAutoAddHandler(func(values map[string][]string) (bool, string) {
		// Check if it's a delete request
		if len(values["del"]) > 0 && values["del"][0] == "1" {
			// Return invalid response that will cause XML parsing to fail
			// This is a bit of a hack, but the mock server framework doesn't
			// support returning raw XML directly
			return false, "" // Empty error message will cause parsing issues
		}
		return true, "success"
	})
	defer mock.Close()

	// Actually, the mock server always returns valid XML format
	// Let's test a different scenario - when ctok cannot be fetched
	// Skip this test as the mock always returns valid XML
	t.Skip("Mock server always returns valid XML format")
}

// TestDeleteAutoAdd_ValidIDs tests deletion with various valid IDs
func TestDeleteAutoAdd_ValidIDs(t *testing.T) {
	testCases := []struct {
		name string
		id   int
	}{
		{"Small ID", 1},
		{"Medium ID", 334},
		{"Large ID", 9999},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create mock server
			mock := testdata.NewMockEMWUIServer()
			defer mock.Close()

			// Track that correct ID was received
			var receivedID int
			mock.SetDeleteAutoAddHandler(func(id int) (bool, string) {
				receivedID = id
				return true, "Deleted successfully"
			})

			// Create client
			apiClient := client.NewClient(mock.URL())

			// Call DeleteAutoAdd
			_, err := apiClient.DeleteAutoAdd(tc.id)
			if err != nil {
				t.Fatalf("DeleteAutoAdd(%d) failed: %v", tc.id, err)
			}

			// Verify correct ID was sent to server
			if receivedID != tc.id {
				t.Errorf("Expected ID %d to be sent to server, got %d", tc.id, receivedID)
			}
		})
	}
}
