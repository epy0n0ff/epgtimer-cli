package models

import (
	"encoding/xml"
	"fmt"
)

// EnumRecInfoResponse represents the response from EMWUI EnumRecInfo API
type EnumRecInfoResponse struct {
	XMLName xml.Name      `xml:"entry"`
	Total   int           `xml:"total"`
	Index   int           `xml:"index"`
	Count   int           `xml:"count"`
	Items   []RecordingInfo `xml:"items>recinfo"`
}

// RecordingInfo represents a single recorded program from EnumRecInfo API
type RecordingInfo struct {
	ID             int    `xml:"ID" json:"id"`
	Title          string `xml:"title" json:"title"`
	StartDate      string `xml:"startDate" json:"start_date"`      // Format: 2025/12/22
	StartTime      string `xml:"startTime" json:"start_time"`      // Format: 22:30:00
	DurationSecond int    `xml:"durationSecond" json:"duration_second"`
	StationName    string `xml:"stationName" json:"station_name"`
	ONID           int    `xml:"ONID" json:"onid"`
	TSID           int    `xml:"TSID" json:"tsid"`
	SID            int    `xml:"SID" json:"sid"`
	EventID        int    `xml:"eventID" json:"event_id"`
	Comment        string `xml:"comment" json:"comment"`
	RecFilePath    string `xml:"recFilePath" json:"rec_file_path"`
	ProtectFlag    int    `xml:"protectFlag" json:"protect_flag"`
}

// ChannelID returns the channel identifier in ONID-TSID-SID format
func (r *RecordingInfo) ChannelID() string {
	return fmt.Sprintf("%d-%d-%d", r.ONID, r.TSID, r.SID)
}

// DurationMinutes returns the duration in minutes
func (r *RecordingInfo) DurationMinutes() int {
	return r.DurationSecond / 60
}

// IsProtected returns true if the recording is protected from deletion
func (r *RecordingInfo) IsProtected() bool {
	return r.ProtectFlag == 1
}
