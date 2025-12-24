package models

import (
	"encoding/xml"
	"fmt"
	"time"
)

// EnumEventInfoResponse represents the response from EMWUI EnumEventInfo API
type EnumEventInfoResponse struct {
	XMLName xml.Name    `xml:"entry"`
	Total   int         `xml:"total"`
	Index   int         `xml:"index"`
	Count   int         `xml:"count"`
	Items   []EventInfo `xml:"items>eventinfo"`
}

// EventInfo represents a single program event from EnumEventInfo API
type EventInfo struct {
	ONID          int           `xml:"ONID" json:"onid"`
	TSID          int           `xml:"TSID" json:"tsid"`
	SID           int           `xml:"SID" json:"sid"`
	EventID       int           `xml:"eventID" json:"event_id"`
	ServiceName   string        `xml:"service_name" json:"service_name"`
	StartDate     string        `xml:"startDate" json:"start_date"`      // Format: 2025/12/23
	StartTime     string        `xml:"startTime" json:"start_time"`      // Format: 06:00:00
	StartDayOfWeek int          `xml:"startDayOfWeek" json:"start_day_of_week"`
	Duration      int           `xml:"duration" json:"duration"`         // Duration in seconds
	EventName     string        `xml:"event_name" json:"event_name"`
	EventText     string        `xml:"event_text" json:"event_text"`
	EventExtText  string        `xml:"event_ext_text" json:"event_ext_text"`
	FreeCAFlag    int           `xml:"freeCAFlag" json:"free_ca_flag"`
	ContentInfo   []ContentInfo `xml:"contentInfo" json:"content_info"`
}

// ContentInfo represents content genre information
type ContentInfo struct {
	Nibble1            int    `xml:"nibble1" json:"nibble1"`
	Nibble2            int    `xml:"nibble2" json:"nibble2"`
	ComponentTypeName  string `xml:"component_type_name" json:"component_type_name"`
}

// ChannelID returns the channel identifier in ONID-TSID-SID format
func (e *EventInfo) ChannelID() string {
	return fmt.Sprintf("%d-%d-%d", e.ONID, e.TSID, e.SID)
}

// StartDateTime parses and returns the start date and time as a time.Time
func (e *EventInfo) StartDateTime() (time.Time, error) {
	// Combine date and time: "2025/12/22 22:30:00"
	dateTimeStr := e.StartDate + " " + e.StartTime
	return time.Parse("2006/01/02 15:04:05", dateTimeStr)
}

// EndDateTime calculates and returns the end date and time
func (e *EventInfo) EndDateTime() (time.Time, error) {
	start, err := e.StartDateTime()
	if err != nil {
		return time.Time{}, err
	}
	return start.Add(time.Duration(e.Duration) * time.Second), nil
}

// DurationMinutes returns the duration in minutes
func (e *EventInfo) DurationMinutes() int {
	return e.Duration / 60
}

// IsFreeCA returns true if the program is free (not scrambled)
func (e *EventInfo) IsFreeCA() bool {
	return e.FreeCAFlag == 0
}

// GenreString returns a comma-separated list of genre names
func (e *EventInfo) GenreString() string {
	if len(e.ContentInfo) == 0 {
		return ""
	}

	genres := make([]string, 0, len(e.ContentInfo))
	for _, content := range e.ContentInfo {
		if content.ComponentTypeName != "" {
			genres = append(genres, content.ComponentTypeName)
		}
	}

	if len(genres) == 0 {
		return ""
	}

	// Return first genre for brevity in table display
	return genres[0]
}
