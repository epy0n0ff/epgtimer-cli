package models

import (
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
