package client

import (
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"
	"time"
)

// Client wraps HTTP client for EMWUI API calls
type Client struct {
	BaseURL    string
	HTTPClient *http.Client
}

// NewClient creates a new EMWUI API client
func NewClient(baseURL string) *Client {
	return &Client{
		BaseURL: baseURL,
		HTTPClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// GetCToken fetches the CSRF token from EMWUI autoaddepg.html page
func (c *Client) GetCToken() (string, error) {
	url := c.BaseURL + "/EMWUI/autoaddepg.html"

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to fetch HTML page: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to fetch HTML page: status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read HTML response: %w", err)
	}

	// Extract ctok from hidden input field
	// Pattern: <input type="hidden" name="ctok" value="TOKEN_VALUE" />
	re := regexp.MustCompile(`<input[^>]*name="ctok"[^>]*value="([^"]+)"`)
	matches := re.FindSubmatch(body)

	if len(matches) < 2 {
		return "", fmt.Errorf("ctok not found in HTML page")
	}

	ctok := string(matches[1])
	if ctok == "" {
		return "", fmt.Errorf("ctok value is empty")
	}

	return ctok, nil
}

// Post sends a POST request with form data
func (c *Client) Post(endpoint string, formData string) ([]byte, error) {
	url := c.BaseURL + endpoint

	req, err := http.NewRequest("POST", url, strings.NewReader(formData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned status %d: %s", resp.StatusCode, string(body))
	}

	return body, nil
}
