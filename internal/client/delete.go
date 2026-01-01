package client

import (
	"encoding/xml"
	"fmt"
	"net/url"

	"github.com/epy0n0ff/epgtimer-cli/internal/models"
)

// DeleteAutoAdd deletes an automatic recording rule via the SetAutoAdd API
func (c *Client) DeleteAutoAdd(id int) (*models.AutoAddRuleResponse, error) {
	if id <= 0 {
		return nil, fmt.Errorf("invalid rule ID: must be greater than 0")
	}

	// Fetch CSRF token from HTML page
	ctok, err := c.GetCToken()
	if err != nil {
		return nil, fmt.Errorf("failed to get CSRF token: %w", err)
	}

	// Create form data with del=1 and ctok
	formData := url.Values{}
	formData.Set("del", "1")
	formData.Set("ctok", ctok)

	// Send POST request with id in query parameter
	endpoint := fmt.Sprintf("/api/SetAutoAdd?id=%d", id)
	body, err := c.Post(endpoint, formData.Encode())
	if err != nil {
		return nil, fmt.Errorf("API call failed: %w", err)
	}

	// Parse XML response
	var response models.AutoAddRuleResponse
	if err := xml.Unmarshal(body, &response); err != nil {
		// If XML parsing fails, show the response body for debugging
		return nil, fmt.Errorf("failed to parse XML response: %w\nResponse body: %s", err, string(body))
	}

	// Check if request was successful
	if !response.IsSuccess() {
		errMsg := response.GetError()
		if errMsg == "" {
			errMsg = "unknown error"
		}
		return nil, fmt.Errorf("API returned error: %s", errMsg)
	}

	return &response, nil
}
