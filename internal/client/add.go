package client

import (
	"encoding/xml"
	"fmt"

	"github.com/epy0n0ff/epgtimer-cli/internal/models"
)

// SetAutoAdd creates a new automatic recording rule via the SetAutoAdd API
func (c *Client) SetAutoAdd(req *models.AutoAddRuleRequest) (*models.AutoAddRuleResponse, error) {
	// Validate request
	if err := req.Validate(); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	// Fetch CSRF token from HTML page
	ctok, err := c.GetCToken()
	if err != nil {
		return nil, fmt.Errorf("failed to get CSRF token: %w", err)
	}
	req.CToken = ctok

	// Convert to form data
	formData := req.ToFormData()

	// Send POST request
	body, err := c.Post("/api/SetAutoAdd?id=0", formData)
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
