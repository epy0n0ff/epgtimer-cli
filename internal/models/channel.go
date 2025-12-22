package models

import (
	"encoding/xml"
	"fmt"
	"strconv"
	"strings"
)

// ServiceListEntry represents a channel in ONID-TSID-SID format
type ServiceListEntry struct {
	ONID int // Original Network ID
	TSID int // Transport Stream ID
	SID  int // Service ID
}

// EnumServiceResponse represents the response from EMWUI EnumService API
type EnumServiceResponse struct {
	XMLName xml.Name      `xml:"entry"`
	Total   int           `xml:"total"`
	Index   int           `xml:"index"`
	Count   int           `xml:"count"`
	Items   []ChannelInfo `xml:"items>serviceinfo"`
}

// ChannelInfo represents a single channel/service from EnumService API
type ChannelInfo struct {
	ONID                   int    `xml:"ONID" json:"onid"`
	TSID                   int    `xml:"TSID" json:"tsid"`
	SID                    int    `xml:"SID" json:"sid"`
	ServiceType            int    `xml:"service_type" json:"service_type"`
	PartialReceptionFlag   int    `xml:"partialReceptionFlag" json:"partial_reception_flag"`
	ServiceProviderName    string `xml:"service_provider_name" json:"service_provider_name"`
	ServiceName            string `xml:"service_name" json:"service_name"`
	NetworkName            string `xml:"network_name" json:"network_name"`
	TSName                 string `xml:"ts_name" json:"ts_name"`
	RemoteControlKeyID     int    `xml:"remote_control_key_id" json:"remote_control_key_id"`
}

// ChannelID returns the channel identifier in ONID-TSID-SID format
func (c *ChannelInfo) ChannelID() string {
	return fmt.Sprintf("%d-%d-%d", c.ONID, c.TSID, c.SID)
}

// IsTV returns true if this is a TV service (service_type == 1)
func (c *ChannelInfo) IsTV() bool {
	return c.ServiceType == 1
}

// IsRadio returns true if this is a radio service (service_type == 2)
func (c *ChannelInfo) IsRadio() bool {
	return c.ServiceType == 2
}

// IsData returns true if this is a data service (service_type == 192)
func (c *ChannelInfo) IsData() bool {
	return c.ServiceType == 192
}

// ServiceTypeString returns a human-readable service type
func (c *ChannelInfo) ServiceTypeString() string {
	switch c.ServiceType {
	case 1:
		return "TV"
	case 2:
		return "Radio"
	case 192:
		return "Data"
	default:
		return fmt.Sprintf("Type%d", c.ServiceType)
	}
}

// String returns the string representation in "ONID-TSID-SID" format
func (s *ServiceListEntry) String() string {
	return fmt.Sprintf("%d-%d-%d", s.ONID, s.TSID, s.SID)
}

// ParseServiceListEntry parses a string in "ONID-TSID-SID" format
// Example: "32736-32736-1024" -> ServiceListEntry{ONID: 32736, TSID: 32736, SID: 1024}
func ParseServiceListEntry(s string) (*ServiceListEntry, error) {
	parts := strings.Split(s, "-")
	if len(parts) != 3 {
		return nil, fmt.Errorf("invalid format: expected 'ONID-TSID-SID', got '%s'", s)
	}

	onid, err := strconv.Atoi(parts[0])
	if err != nil {
		return nil, fmt.Errorf("invalid ONID '%s': %w", parts[0], err)
	}

	tsid, err := strconv.Atoi(parts[1])
	if err != nil {
		return nil, fmt.Errorf("invalid TSID '%s': %w", parts[1], err)
	}

	sid, err := strconv.Atoi(parts[2])
	if err != nil {
		return nil, fmt.Errorf("invalid SID '%s': %w", parts[2], err)
	}

	return &ServiceListEntry{
		ONID: onid,
		TSID: tsid,
		SID:  sid,
	}, nil
}
