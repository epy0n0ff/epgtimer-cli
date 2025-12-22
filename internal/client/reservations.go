package client

import (
	"encoding/xml"
	"fmt"
	"io"

	"github.com/epy0n0ff/epgtimer-cli/internal/models"
)

// EnumReserveInfo retrieves all manual reservations from EMWUI
func (c *Client) EnumReserveInfo() (*models.EnumReserveInfoResponse, error) {
	url := c.BaseURL + "/api/EnumReserveInfo"

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

	var response models.EnumReserveInfoResponse
	if err := xml.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("failed to parse XML response: %w\nResponse body: %s", err, string(body))
	}

	return &response, nil
}
