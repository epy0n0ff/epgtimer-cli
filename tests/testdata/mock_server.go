package testdata

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
)

// MockEMWUIServer creates a test HTTP server that simulates EMWUI SetAutoAdd API
type MockEMWUIServer struct {
	Server *httptest.Server
	// Callbacks for custom behavior
	OnSetAutoAdd func(values map[string][]string) (success bool, message string)
	// CSRF token to return in HTML page
	CToken string
}

// NewMockEMWUIServer creates a new mock EMWUI server
func NewMockEMWUIServer() *MockEMWUIServer {
	mock := &MockEMWUIServer{
		CToken: "test-csrf-token-12345", // Default test token
	}

	// Default handler
	mock.Server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Handle HTML page request for ctok
		if r.URL.Path == "/EMWUI/autoaddepg.html" {
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			html := fmt.Sprintf(`<!DOCTYPE html>
<html>
<head><title>AutoAdd EPG</title></head>
<body>
<form method="post" action="/api/SetAutoAdd">
<input type="hidden" name="ctok" value="%s" />
<input type="text" name="andKey" />
</form>
</body>
</html>`, mock.CToken)
			fmt.Fprint(w, html)
			return
		}

		// Only handle SetAutoAdd endpoint
		if !strings.HasPrefix(r.URL.Path, "/api/SetAutoAdd") {
			http.Error(w, "Not found", http.StatusNotFound)
			return
		}

		// Check method
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// Parse form data
		if err := r.ParseForm(); err != nil {
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}

		// Extract values
		values := make(map[string][]string)
		for key, vals := range r.Form {
			values[key] = vals
		}

		// Call custom handler if set
		var success bool
		var message string
		if mock.OnSetAutoAdd != nil {
			success, message = mock.OnSetAutoAdd(values)
		} else {
			// Default: success if andKey and serviceList are present
			success = len(values["andKey"]) > 0 && values["andKey"][0] != "" &&
				len(values["serviceList"]) > 0
			if success {
				message = "Automatic recording rule created successfully"
			} else {
				message = "Missing required parameters"
			}
		}

		// Send XML response (matching EMWUI format)
		w.Header().Set("Content-Type", "text/xml; charset=utf-8")
		if success {
			fmt.Fprintf(w, `<?xml version="1.0" encoding="UTF-8" ?><entry><success>%s</success></entry>`, message)
		} else {
			fmt.Fprintf(w, `<?xml version="1.0" encoding="UTF-8" ?><entry><err>%s</err></entry>`, message)
		}
	}))

	return mock
}

// Close shuts down the mock server
func (m *MockEMWUIServer) Close() {
	m.Server.Close()
}

// URL returns the base URL of the mock server
func (m *MockEMWUIServer) URL() string {
	return m.Server.URL
}

// SetAutoAddHandler sets a custom handler for SetAutoAdd requests
func (m *MockEMWUIServer) SetAutoAddHandler(handler func(values map[string][]string) (success bool, message string)) {
	m.OnSetAutoAdd = handler
}

// NewFailingServer creates a mock server that always returns errors
func NewFailingServer() *MockEMWUIServer {
	mock := NewMockEMWUIServer()
	mock.CToken = "test-csrf-token-failing"
	mock.OnSetAutoAdd = func(values map[string][]string) (bool, string) {
		return false, "Internal server error"
	}
	return mock
}

// NewHTMLServer creates a mock server that returns HTML instead of JSON (simulates wrong endpoint)
func NewHTMLServer() *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		fmt.Fprintf(w, "<html><body>EMWUI Web Interface</body></html>")
	}))
}
