package client

import (
	"encoding/xml"
	"fmt"
	"io"

	"github.com/epy0n0ff/epgtimer-cli/internal/models"
)

// EnumEventInfo retrieves EPG (Electronic Program Guide) data for a specific channel
func (c *Client) EnumEventInfo(onid, tsid, sid int) (*models.EnumEventInfoResponse, error) {
	url := fmt.Sprintf("%s/api/EnumEventInfo?ONID=%d&TSID=%d&SID=%d&basic=0&count=1000",
		c.BaseURL, onid, tsid, sid)

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

	var response models.EnumEventInfoResponse
	if err := xml.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("failed to parse XML response: %w\nResponse body: %s", err, string(body))
	}

	return &response, nil
}
