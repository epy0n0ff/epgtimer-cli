package client

import (
	"encoding/xml"
	"fmt"
	"io"

	"github.com/epy0n0ff/epgtimer-cli/internal/models"
)

// EnumAutoAdd retrieves all automatic recording rules from EMWUI API
// GET /api/EnumAutoAdd
func (c *Client) EnumAutoAdd() (*models.EnumAutoAddResponse, error) {
	// Build request URL
	url := c.BaseURL + "/api/EnumAutoAdd"

	// Send GET request
	resp, err := c.HTTPClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to EMWUI service at %s: %w", c.BaseURL, err)
	}
	defer resp.Body.Close()

	// Check status code
	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API returned status %d: %s", resp.StatusCode, string(body))
	}

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	// Parse XML response
	var response models.EnumAutoAddResponse
	if err := xml.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("failed to parse XML response: %w\nResponse body: %s", err, string(body))
	}

	return &response, nil
}
