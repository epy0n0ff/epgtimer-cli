package models

import (
	"fmt"
	"net/url"
	"strings"
)

// AutoAddRuleRequest contains all parameters for creating an automatic recording rule
type AutoAddRuleRequest struct {
	// User-provided parameters (required by CLI)
	AndKey      string   `json:"and_key"`      // Search keywords (title must contain)
	NotKey      string   `json:"not_key"`      // Exclusion keywords (title must not contain)
	ServiceList []string `json:"service_list"` // Channel list in "ONID-TSID-SID" format

	// Default parameters (from curl sample - user should not modify these)
	AddChg           int    `json:"addchg"`              // 1 = add/change mode
	TitleOnlyFlag    int    `json:"title_only_flag"`     // 1 = search title only
	DayList          string `json:"day_list"`            // "on" = all days
	StartTime        string `json:"start_time"`          // Time filter start "HH:MM"
	EndTime          string `json:"end_time"`            // Time filter end "HH:MM"
	DateList         string `json:"date_list"`           // Date filter (empty = all dates)
	FreeCAFlag       int    `json:"free_ca_flag"`        // 0 = include pay channels
	ChkDurationMin   int    `json:"chk_duration_min"`    // Minimum duration filter (0 = no limit)
	ChkDurationMax   int    `json:"chk_duration_max"`    // Maximum duration filter (0 = no limit)
	ChkRecDay        int    `json:"chk_rec_day"`         // Recording days mask
	PresetID         int    `json:"preset_id"`           // Preset ID (65535 = custom)
	ONID             string `json:"onid"`                // Original Network ID filter (empty = all)
	TSID             string `json:"tsid"`                // Transport Stream ID filter (empty = all)
	SID              string `json:"sid"`                 // Service ID filter (empty = all)
	EID              string `json:"eid"`                 // Event ID filter (empty = all)
	CToken           string `json:"ctok"`                // CSRF token (fetched from HTML page)
	RecMode          int    `json:"rec_mode"`            // Recording mode (1 = standard)
	TuijyuuFlag      int    `json:"tuijyuu_flag"`        // Auto-follow flag (1 = enabled)
	Priority         int    `json:"priority"`            // Recording priority (2 = normal)
	UseDefMarginFlag int    `json:"use_def_margin_flag"` // Use default margins (1 = yes)
	ServiceMode      int    `json:"service_mode"`        // Service mode (1 = enabled)
	TunerID          int    `json:"tuner_id"`            // Tuner ID (0 = auto)
	SuspendMode      int    `json:"suspend_mode"`        // Suspend mode (0 = disabled)
	BatFilePath      string `json:"bat_file_path"`       // Batch file path (empty)
	BatFileTag       string `json:"bat_file_tag"`        // Batch file tag (empty)
}

// NewAutoAddRuleRequest creates a new request with default values from curl sample
func NewAutoAddRuleRequest(andKey string, notKey string, serviceList []string) *AutoAddRuleRequest {
	return &AutoAddRuleRequest{
		// User parameters
		AndKey:      andKey,
		NotKey:      notKey,
		ServiceList: serviceList,

		// Defaults from curl sample
		AddChg:           1,
		TitleOnlyFlag:    1,
		DayList:          "on",
		StartTime:        "00:00",
		EndTime:          "01:00",
		DateList:         "",
		FreeCAFlag:       0,
		ChkDurationMin:   0,
		ChkDurationMax:   0,
		ChkRecDay:        6,
		PresetID:         0, // 0 = use current settings (from curl sample)
		ONID:             "",
		TSID:             "",
		SID:              "",
		EID:              "",
		CToken:           "", // Will be fetched dynamically from HTML page
		RecMode:          1,
		TuijyuuFlag:      1,
		Priority:         2,
		UseDefMarginFlag: 1,
		ServiceMode:      1,
		TunerID:          0,
		SuspendMode:      0,
		BatFilePath:      "",
		BatFileTag:       "",
	}
}

// Validate checks if the request has valid parameters
func (r *AutoAddRuleRequest) Validate() error {
	// Validate AndKey (required)
	if strings.TrimSpace(r.AndKey) == "" {
		return fmt.Errorf("andKey is required (search keyword cannot be empty)")
	}

	// Validate ServiceList (required, at least one channel)
	if len(r.ServiceList) == 0 {
		return fmt.Errorf("serviceList is required (at least one channel must be specified)")
	}

	// Validate ServiceList format (ONID-TSID-SID)
	for i, service := range r.ServiceList {
		parts := strings.Split(service, "-")
		if len(parts) != 3 {
			return fmt.Errorf("serviceList[%d] has invalid format: expected 'ONID-TSID-SID', got '%s'", i, service)
		}
		// Could add numeric validation here if needed
	}

	// NotKey is optional (can be empty)

	return nil
}

// ToFormData converts the request to application/x-www-form-urlencoded format
// Order matches the curl sample format exactly
func (r *AutoAddRuleRequest) ToFormData() string {
	v := url.Values{}

	// Follow curl sample order
	v.Set("addchg", fmt.Sprintf("%d", r.AddChg))
	v.Set("andKey", r.AndKey)
	v.Set("notKey", r.NotKey) // Can be empty

	// Add empty serviceList first (as in curl sample), then actual values
	v.Add("serviceList", "")
	for _, service := range r.ServiceList {
		v.Add("serviceList", service)
	}

	v.Set("dayList", r.DayList)
	v.Set("startTime", r.StartTime)
	v.Set("endTime", r.EndTime)
	v.Set("dateList", r.DateList)
	v.Set("freeCAFlag", fmt.Sprintf("%d", r.FreeCAFlag))
	v.Set("chkDurationMin", fmt.Sprintf("%d", r.ChkDurationMin))
	v.Set("chkDurationMax", fmt.Sprintf("%d", r.ChkDurationMax))
	v.Set("chkRecDay", fmt.Sprintf("%d", r.ChkRecDay))

	// presetID appears twice in curl sample
	v.Add("presetID", fmt.Sprintf("%d", r.PresetID))
	v.Add("presetID", fmt.Sprintf("%d", r.PresetID))

	v.Set("ctok", r.CToken)
	v.Set("recMode", fmt.Sprintf("%d", r.RecMode))
	v.Set("tuijyuuFlag", fmt.Sprintf("%d", r.TuijyuuFlag))
	v.Set("priority", fmt.Sprintf("%d", r.Priority))
	v.Set("useDefMarginFlag", fmt.Sprintf("%d", r.UseDefMarginFlag))
	v.Set("serviceMode", fmt.Sprintf("%d", r.ServiceMode))
	v.Set("tunerID", fmt.Sprintf("%d", r.TunerID))
	v.Set("suspendMode", fmt.Sprintf("%d", r.SuspendMode))
	v.Set("batFilePath", r.BatFilePath)
	v.Set("batFileTag", r.BatFileTag)

	return v.Encode()
}
