package models

import (
	"encoding/xml"
	"fmt"
)

// EnumAutoAddResponse represents the root response structure from GET /api/EnumAutoAdd
type EnumAutoAddResponse struct {
	XMLName xml.Name      `xml:"entry"`
	Total   int           `xml:"total"`
	Index   int           `xml:"index"`
	Count   int           `xml:"count"`
	Items   []AutoAddRule `xml:"items>autoaddinfo"`
}

// AutoAddRule represents a single automatic recording rule configuration
type AutoAddRule struct {
	ID                int               `xml:"ID" json:"id"`
	SearchSettings    SearchSettings    `xml:"searchsetting" json:"search"`
	RecordingSettings RecordingSettings `xml:"recsetting" json:"recording"`
}

// SearchSettings defines keyword-based search criteria for automatic recording
type SearchSettings struct {
	DisableFlag     int           `xml:"disableFlag" json:"disabled"`
	CaseFlag        int           `xml:"caseFlag" json:"case_sensitive"`
	AndKey          string        `xml:"andKey" json:"and_key"`
	NotKey          string        `xml:"notKey" json:"not_key"`
	RegExpFlag      int           `xml:"regExpFlag" json:"regex"`
	TitleOnlyFlag   int           `xml:"titleOnlyFlag" json:"title_only"`
	AimaiFlag       int           `xml:"aimaiFlag" json:"fuzzy"`
	NotContetFlag   int           `xml:"notContetFlag" json:"not_content"`
	NotDateFlag     int           `xml:"notDateFlag" json:"not_date"`
	FreeCAFlag      int           `xml:"freeCAFlag" json:"free_ca"`
	ChkRecEnd       int           `xml:"chkRecEnd" json:"check_rec_end"`
	ChkRecDay       int           `xml:"chkRecDay" json:"rec_day_mask"`
	ChkRecNoService int           `xml:"chkRecNoService" json:"check_no_service"`
	ChkDurationMin  int           `xml:"chkDurationMin" json:"duration_min"`
	ChkDurationMax  int           `xml:"chkDurationMax" json:"duration_max"`
	ServiceList     []ServiceInfo `xml:"serviceList" json:"channels"`
}

// IsEnabled returns true if the rule is enabled (DisableFlag == 0)
func (s *SearchSettings) IsEnabled() bool {
	return s.DisableFlag == 0
}

// IsRegex returns true if regex matching is enabled (RegExpFlag == 1)
func (s *SearchSettings) IsRegex() bool {
	return s.RegExpFlag == 1
}

// HasDurationFilter returns true if duration filtering is active
func (s *SearchSettings) HasDurationFilter() bool {
	return s.ChkDurationMin > 0 || s.ChkDurationMax > 0
}

// ChannelCount returns the number of channels in the service list
func (s *SearchSettings) ChannelCount() int {
	return len(s.ServiceList)
}

// ServiceInfo represents a broadcast channel identifier in ONID-TSID-SID format
type ServiceInfo struct {
	ONID int `xml:"onid" json:"onid"`
	TSID int `xml:"tsid" json:"tsid"`
	SID  int `xml:"sid" json:"sid"`
}

// String returns the channel identifier in "ONID-TSID-SID" format
func (s *ServiceInfo) String() string {
	return fmt.Sprintf("%d-%d-%d", s.ONID, s.TSID, s.SID)
}

// Matches returns true if the service info matches the given ONID/TSID/SID
func (s *ServiceInfo) Matches(onid, tsid, sid int) bool {
	return s.ONID == onid && s.TSID == tsid && s.SID == sid
}

// RecordingSettings defines recording behavior and post-processing options
type RecordingSettings struct {
	RecMode          int    `xml:"recMode" json:"rec_mode"`
	Priority         int    `xml:"priority" json:"priority"`
	TuijyuuFlag      int    `xml:"tuijyuuFlag" json:"auto_follow"`
	ServiceMode      int    `xml:"serviceMode" json:"service_mode"`
	PittariFlag      int    `xml:"pittariFlag" json:"exact_match"`
	BatFilePath      string `xml:"batFilePath" json:"bat_file"`
	RecFolderList    string `xml:"recFolderList" json:"rec_folders"`
	SuspendMode      int    `xml:"suspendMode" json:"suspend_mode"`
	DefServiceMode   int    `xml:"defserviceMode" json:"def_service_mode"`
	RebootFlag       int    `xml:"rebootFlag" json:"reboot"`
	UseMargineFlag   int    `xml:"useMargineFlag" json:"use_margin"`
	StartMargine     int    `xml:"startMargine" json:"start_margin"`
	EndMargine       int    `xml:"endMargine" json:"end_margin"`
	ContinueRecFlag  int    `xml:"continueRecFlag" json:"continue_rec"`
	PartialRecFlag   int    `xml:"partialRecFlag" json:"partial_rec"`
	TunerID          int    `xml:"tunerID" json:"tuner_id"`
	PartialRecFolder string `xml:"partialRecFolder" json:"partial_rec_folder"`
}

// IsAutoFollow returns true if auto-follow is enabled (TuijyuuFlag == 1)
func (r *RecordingSettings) IsAutoFollow() bool {
	return r.TuijyuuFlag == 1
}

// HasMargins returns true if recording margins are enabled (UseMargineFlag == 1)
func (r *RecordingSettings) HasMargins() bool {
	return r.UseMargineFlag == 1
}

// UsesCustomTuner returns true if a specific tuner is selected (TunerID > 0)
func (r *RecordingSettings) UsesCustomTuner() bool {
	return r.TunerID > 0
}
