# Data Model: EpgTimer CLI (SetAutoAdd)

**Feature**: `001-epgtimer-cli`
**Created**: 2025-12-20
**Purpose**: Define data structures and validation rules for creating automatic recording rules via EMWUI SetAutoAdd API

---

## Overview

This document defines the domain models and API request/response types for the EpgTimer CLI. The CLI uses the SetAutoAdd endpoint to create keyword-based automatic recording rules that will trigger recordings for any program matching the specified criteria.

**API Endpoint**: `POST /api/SetAutoAdd?id=0`
**Content-Type**: `application/x-www-form-urlencoded`

---

## Core Domain Models

### 1. AutoAddRuleRequest

Represents a request to create a new automatic recording rule.

```go
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
	AddChg           int    `json:"addchg"`             // 1 = add/change mode
	TitleOnlyFlag    int    `json:"title_only_flag"`    // 1 = search title only
	DayList          string `json:"day_list"`           // "on" = all days
	StartTime        string `json:"start_time"`         // Time filter start "HH:MM"
	EndTime          string `json:"end_time"`           // Time filter end "HH:MM"
	DateList         string `json:"date_list"`          // Date filter (empty = all dates)
	FreeCAFlag       int    `json:"free_ca_flag"`       // 0 = include pay channels
	ChkDurationMin   int    `json:"chk_duration_min"`   // Minimum duration filter (0 = no limit)
	ChkDurationMax   int    `json:"chk_duration_max"`   // Maximum duration filter (0 = no limit)
	ChkRecDay        int    `json:"chk_rec_day"`        // Recording days mask
	PresetID         int    `json:"preset_id"`          // Preset ID (65535 = custom)
	ONID             string `json:"onid"`               // Original Network ID filter (empty = all)
	TSID             string `json:"tsid"`               // Transport Stream ID filter (empty = all)
	SID              string `json:"sid"`                // Service ID filter (empty = all)
	EID              string `json:"eid"`                // Event ID filter (empty = all)
	CToken           string `json:"ctok"`               // CSRF token (fetched from HTML page)
	RecMode          int    `json:"rec_mode"`           // Recording mode (1 = standard)
	TuijyuuFlag      int    `json:"tuijyuu_flag"`       // Auto-follow flag (1 = enabled)
	Priority         int    `json:"priority"`           // Recording priority (2 = normal)
	UseDefMarginFlag int    `json:"use_def_margin_flag"` // Use default margins (1 = yes)
	ServiceMode      int    `json:"service_mode"`       // Service mode (1 = enabled)
	TunerID          int    `json:"tuner_id"`           // Tuner ID (0 = auto)
	SuspendMode      int    `json:"suspend_mode"`       // Suspend mode (0 = disabled)
	BatFilePath      string `json:"bat_file_path"`      // Batch file path (empty)
	BatFileTag       string `json:"bat_file_tag"`       // Batch file tag (empty)
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
		PresetID:         65535,
		ONID:             "",
		TSID:             "",
		SID:              "",
		EID:              "",
		CToken:           "98357b8eedf096855c1cb636303ab2af", // From curl sample
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
func (r *AutoAddRuleRequest) ToFormData() string {
	v := url.Values{}

	// Add user parameters
	v.Set("andKey", r.AndKey)
	v.Set("notKey", r.NotKey) // Can be empty

	// Add serviceList (multiple values)
	for _, service := range r.ServiceList {
		v.Add("serviceList", service)
	}

	// Add default parameters
	v.Set("addchg", fmt.Sprintf("%d", r.AddChg))
	v.Set("titleOnlyFlag", fmt.Sprintf("%d", r.TitleOnlyFlag))
	v.Set("dayList", r.DayList)
	v.Set("startTime", r.StartTime)
	v.Set("endTime", r.EndTime)
	v.Set("dateList", r.DateList)
	v.Set("freeCAFlag", fmt.Sprintf("%d", r.FreeCAFlag))
	v.Set("chkDurationMin", fmt.Sprintf("%d", r.ChkDurationMin))
	v.Set("chkDurationMax", fmt.Sprintf("%d", r.ChkDurationMax))
	v.Set("chkRecDay", fmt.Sprintf("%d", r.ChkRecDay))
	v.Set("presetID", fmt.Sprintf("%d", r.PresetID))
	v.Set("onid", r.ONID)
	v.Set("tsid", r.TSID)
	v.Set("sid", r.SID)
	v.Set("eid", r.EID)
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
```

---

### 2. AutoAddRuleResponse

Represents the response from the EMWUI API after attempting to add a rule.

```go
package models

// AutoAddRuleResponse contains the result of an add rule request
type AutoAddRuleResponse struct {
	// Success indicates if the rule was created successfully
	Success bool `json:"success"`

	// Message contains either a success confirmation or error details
	Message string `json:"message"`

	// HTTPStatus is the HTTP status code from the EMWUI API
	HTTPStatus int `json:"http_status"`
}

// IsSuccess checks if the response indicates success
func (r *AutoAddRuleResponse) IsSuccess() bool {
	return r.Success && r.HTTPStatus == 200
}

// ErrorMessage returns a user-friendly error message
func (r *AutoAddRuleResponse) ErrorMessage() string {
	if r.IsSuccess() {
		return ""
	}

	// Map common error patterns to user-friendly messages
	switch {
	case r.HTTPStatus == 0:
		return "Failed to connect to EMWUI service. Check EMWUI_ENDPOINT configuration."
	case r.HTTPStatus >= 500:
		return fmt.Sprintf("EMWUI service error (HTTP %d): %s", r.HTTPStatus, r.Message)
	case r.HTTPStatus == 404:
		return "EMWUI API endpoint not found. Check EMWUI version compatibility."
	case strings.Contains(r.Message, "duplicate"):
		return "A similar automatic recording rule already exists"
	default:
		return fmt.Sprintf("Failed to add rule: %s (HTTP %d)", r.Message, r.HTTPStatus)
	}
}
```

---

### 3. ServiceListEntry

Represents a channel in the serviceList format.

```go
package models

import (
	"fmt"
	"strconv"
	"strings"
)

// ServiceListEntry represents a channel in ONID-TSID-SID format
type ServiceListEntry struct {
	ONID uint16 // Original Network ID
	TSID uint16 // Transport Stream ID
	SID  uint16 // Service ID (channel ID)
}

// String returns the serviceList format "ONID-TSID-SID"
func (s *ServiceListEntry) String() string {
	return fmt.Sprintf("%d-%d-%d", s.ONID, s.TSID, s.SID)
}

// ParseServiceListEntry parses a serviceList entry from "ONID-TSID-SID" format
func ParseServiceListEntry(entry string) (*ServiceListEntry, error) {
	parts := strings.Split(entry, "-")
	if len(parts) != 3 {
		return nil, fmt.Errorf("invalid serviceList format: expected 'ONID-TSID-SID', got '%s'", entry)
	}

	onid, err := strconv.ParseUint(parts[0], 10, 16)
	if err != nil {
		return nil, fmt.Errorf("invalid ONID: %w", err)
	}

	tsid, err := strconv.ParseUint(parts[1], 10, 16)
	if err != nil {
		return nil, fmt.Errorf("invalid TSID: %w", err)
	}

	sid, err := strconv.ParseUint(parts[2], 10, 16)
	if err != nil {
		return nil, fmt.Errorf("invalid SID: %w", err)
	}

	return &ServiceListEntry{
		ONID: uint16(onid),
		TSID: uint16(tsid),
		SID:  uint16(sid),
	}, nil
}

// Common channel definitions for user convenience
var (
	// Terrestrial channels (Tokyo area example)
	NHKSogo        = &ServiceListEntry{ONID: 32736, TSID: 32736, SID: 1024}  // NHK総合
	NHKEtv         = &ServiceListEntry{ONID: 32736, TSID: 32736, SID: 1025}  // NHK教育
	NTV            = &ServiceListEntry{ONID: 32737, TSID: 32737, SID: 1032}  // 日本テレビ
	TBS            = &ServiceListEntry{ONID: 32738, TSID: 32738, SID: 1040}  // TBS
	Fuji           = &ServiceListEntry{ONID: 32739, TSID: 32739, SID: 1048}  // フジテレビ
	TVAsahi        = &ServiceListEntry{ONID: 32740, TSID: 32740, SID: 1056}  // テレビ朝日
	TVTokyo        = &ServiceListEntry{ONID: 32741, TSID: 32741, SID: 1064}  // テレビ東京

	// BS channels (examples)
	BS1            = &ServiceListEntry{ONID: 4, TSID: 16625, SID: 101}  // BS1
	BSPremium      = &ServiceListEntry{ONID: 4, TSID: 16625, SID: 103}  // BSプレミアム
)
```

---

## Validation Rules

### Field-Level Validation

| Field | Required | Type | Constraints | Example |
|-------|----------|------|-------------|---------|
| andKey | Yes | string | 1-255 chars, UTF-8 | "ニュース" |
| notKey | No | string | 0-255 chars, UTF-8 | "再放送" |
| serviceList | Yes | []string | At least 1, format "ONID-TSID-SID" | ["32736-32736-1024"] |

### Business Logic Validation

1. **AndKey Validation**
   - Must not be empty or only whitespace
   - Maximum length: 255 characters (reasonable limit)
   - UTF-8 encoding required for Japanese text

2. **NotKey Validation**
   - Optional (can be empty)
   - Maximum length: 255 characters if provided
   - UTF-8 encoding required

3. **ServiceList Validation**
   - At least one entry required
   - Each entry must match format: "ONID-TSID-SID"
   - Components must be valid uint16 numbers
   - No duplicate entries (optional validation)

---

## Example Usage

### Creating an Automatic Recording Rule

```go
package main

import (
	"fmt"
	"internal/models"
)

func main() {
	// Create request with user parameters
	req := models.NewAutoAddRuleRequest(
		"わたしが恋人になれるわけ",  // andKey: search keyword
		"推しエンタ",                // notKey: exclusion keyword
		[]string{
			"32736-32736-1024",  // NHK総合
			"32736-32736-1025",  // NHK教育
			"32737-32737-1032",  // 日本テレビ
		},
	)

	// Validate
	if err := req.Validate(); err != nil {
		fmt.Printf("Validation error: %v\n", err)
		return
	}

	// Convert to form data for POST request
	formData := req.ToFormData()
	fmt.Printf("Form Data Length: %d bytes\n", len(formData))

	// formData is ready to send to EMWUI API
}
```

### Parsing ServiceList Entry

```go
// Parse from user input
entry, err := models.ParseServiceListEntry("32736-32736-1024")
if err != nil {
	fmt.Printf("Invalid format: %v\n", err)
	return
}

fmt.Printf("Channel: ONID=%d, TSID=%d, SID=%d\n",
	entry.ONID, entry.TSID, entry.SID)

// Use predefined channels
nhk := models.NHKSogo
fmt.Printf("NHK総合: %s\n", nhk.String())
// Output: NHK総合: 32736-32736-1024
```

---

## Testing Considerations

### Test Data

Create fixtures for:
1. **Valid requests**: Various keywords and channel combinations
2. **Invalid requests**: Missing andKey, empty serviceList, invalid format
3. **Edge cases**: Very long keywords, Japanese with special chars, many channels
4. **Form encoding**: Verify URL encoding of Japanese text

### Example Test Cases

```go
func TestAutoAddRuleRequest_Validate(t *testing.T) {
	tests := []struct {
		name    string
		req     *models.AutoAddRuleRequest
		wantErr bool
	}{
		{
			name: "valid request",
			req: models.NewAutoAddRuleRequest(
				"ニュース",
				"",
				[]string{"32736-32736-1024"},
			),
			wantErr: false,
		},
		{
			name: "missing andKey",
			req: models.NewAutoAddRuleRequest(
				"",
				"",
				[]string{"32736-32736-1024"},
			),
			wantErr: true,
		},
		{
			name: "empty serviceList",
			req: models.NewAutoAddRuleRequest(
				"ドラマ",
				"",
				[]string{},
			),
			wantErr: true,
		},
		{
			name: "invalid serviceList format",
			req: models.NewAutoAddRuleRequest(
				"映画",
				"",
				[]string{"invalid-format"},
			),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.req.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
```

---

## Default Parameters Reference

These parameters use values from the curl sample and should not normally be changed by users:

| Parameter | Default Value | Meaning |
|-----------|---------------|---------|
| addchg | 1 | Add/change mode |
| titleOnlyFlag | 1 | Search title only (not description) |
| dayList | "on" | All days enabled |
| startTime | "00:00" | Time filter start (midnight) |
| endTime | "01:00" | Time filter end (1 AM) |
| chkRecDay | 6 | Recording days mask |
| presetID | 65535 | Custom preset (not using predefined) |
| recMode | 1 | Standard recording mode |
| tuijyuuFlag | 1 | Auto-follow enabled |
| priority | 2 | Normal priority |
| ctok | (dynamic) | CSRF token fetched from HTML page |

**Note**: The ctok (CSRF token) value is dynamically retrieved from `/EMWUI/autoaddepg.html` before each request by parsing the hidden input field `<input type="hidden" name="ctok" value="..." />`.

---

**End of Data Model Document**

Last Updated: 2025-12-20
