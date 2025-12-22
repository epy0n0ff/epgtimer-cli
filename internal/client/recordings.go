package client

import (
	"encoding/xml"
	"fmt"
	"io"

	"github.com/epy0n0ff/epgtimer-cli/internal/models"
)

// EnumRecInfo retrieves all recorded programs from EMWUI
func (c *Client) EnumRecInfo() (*models.EnumRecInfoResponse, error) {
	url := c.BaseURL + "/api/EnumRecInfo"

	resp, err := c.HTTPClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to EMWUI service: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("API returned status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var response models.EnumRecInfoResponse
	if err := xml.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("failed to parse XML response: %w\nResponse body: %s", err, string(body))
	}

	return &response, nil
}
