# Data Model: EPG Reservation List Retrieval

**Feature**: EPG Reservation List Retrieval (001-epg-reservation-list)
**Date**: 2025-12-21
**Phase**: 1 - Design

## Overview

This document defines the data structures for automatic recording rules retrieved from the EMWUI EnumAutoAdd API. The model maps XML response elements to Go structs with support for both XML parsing (input) and multiple export formats (output: JSON, CSV, TSV).

## Core Entities

### EnumAutoAddResponse

**Purpose**: Root response structure from GET /api/EnumAutoAdd

**Fields**:
| Field | Type | XML Tag | Description | Validation |
|-------|------|---------|-------------|------------|
| XMLName | xml.Name | `xml:"entry"` | XML root element name | - |
| Total | int | `xml:"total"` | Total count of rules in system | >= 0 |
| Index | int | `xml:"index"` | Starting index (always 0) | >= 0 |
| Count | int | `xml:"count"` | Number of rules returned | >= 0, <= Total |
| Items | []AutoAddRule | `xml:"items>autoaddinfo"` | Array of automatic recording rules | - |

**Relationships**:
- Contains many AutoAddRule entities
- Count should equal len(Items)

**State Transitions**: N/A (read-only response)

**Example**:
```go
type EnumAutoAddResponse struct {
    XMLName xml.Name      `xml:"entry"`
    Total   int           `xml:"total"`
    Index   int           `xml:"index"`
    Count   int           `xml:"count"`
    Items   []AutoAddRule `xml:"items>autoaddinfo"`
}
```

---

### AutoAddRule

**Purpose**: Represents a single automatic recording rule configuration

**Fields**:
| Field | Type | XML/JSON Tag | Description | Validation |
|-------|------|--------------|-------------|------------|
| ID | int | `xml:"ID" json:"id"` | Unique rule identifier | > 0 |
| SearchSettings | SearchSettings | `xml:"searchsetting" json:"search"` | Keyword search criteria | Required |
| RecordingSettings | RecordingSettings | `xml:"recsetting" json:"recording"` | Recording configuration | Required |

**Relationships**:
- Has one SearchSettings
- Has one RecordingSettings

**Display Representation**:
- Default (table): ID, Enabled status, AndKey, NotKey, Channel count
- Detailed: All fields from SearchSettings and RecordingSettings
- Export: Full structure with nested objects

**Example**:
```go
type AutoAddRule struct {
    ID                int               `xml:"ID" json:"id"`
    SearchSettings    SearchSettings    `xml:"searchsetting" json:"search"`
    RecordingSettings RecordingSettings `xml:"recsetting" json:"recording"`
}
```

---

### SearchSettings

**Purpose**: Defines keyword-based search criteria for automatic recording

**Fields**:
| Field | Type | XML/JSON Tag | Description | Validation | Default |
|-------|------|--------------|-------------|------------|---------|
| DisableFlag | int | `xml:"disableFlag" json:"disabled"` | 0=enabled, 1=disabled | 0 or 1 | 0 |
| CaseFlag | int | `xml:"caseFlag" json:"case_sensitive"` | 0=case-insensitive, 1=case-sensitive | 0 or 1 | 0 |
| AndKey | string | `xml:"andKey" json:"and_key"` | Required search keywords (UTF-8) | - | "" |
| NotKey | string | `xml:"notKey" json:"not_key"` | Exclusion keywords (UTF-8) | - | "" |
| RegExpFlag | int | `xml:"regExpFlag" json:"regex"` | 0=literal, 1=regex pattern | 0 or 1 | 0 |
| TitleOnlyFlag | int | `xml:"titleOnlyFlag" json:"title_only"` | 0=all fields, 1=title only | 0 or 1 | 0 |
| AimaiFlag | int | `xml:"aimaiFlag" json:"fuzzy"` | Fuzzy matching enabled | 0 or 1 | 0 |
| NotContetFlag | int | `xml:"notContetFlag" json:"not_content"` | Not content flag | 0 or 1 | 0 |
| NotDateFlag | int | `xml:"notDateFlag" json:"not_date"` | Not date flag | 0 or 1 | 0 |
| FreeCAFlag | int | `xml:"freeCAFlag" json:"free_ca"` | 0=include pay channels, 1=free only | 0 or 1 | 0 |
| ChkRecEnd | int | `xml:"chkRecEnd" json:"check_rec_end"` | Check recording end | >= 0 | 0 |
| ChkRecDay | int | `xml:"chkRecDay" json:"rec_day_mask"` | Recording days bitmask | >= 0 | 6 |
| ChkRecNoService | int | `xml:"chkRecNoService" json:"check_no_service"` | Check no service | 0 or 1 | 0 |
| ChkDurationMin | int | `xml:"chkDurationMin" json:"duration_min"` | Minimum duration (minutes, 0=no limit) | >= 0 | 0 |
| ChkDurationMax | int | `xml:"chkDurationMax" json:"duration_max"` | Maximum duration (minutes, 0=no limit) | >= 0 | 0 |
| ServiceList | []ServiceInfo | `xml:"serviceList" json:"channels"` | Array of channel identifiers | - | [] |

**Relationships**:
- Contains many ServiceInfo (channels)

**Validation Rules**:
- AndKey or NotKey should be non-empty (at least one search criterion)
- If ChkDurationMin > 0 and ChkDurationMax > 0, then ChkDurationMin < ChkDurationMax
- ServiceList can be empty (searches all channels)

**Business Logic**:
- IsEnabled(): Returns true if DisableFlag == 0
- IsRegex(): Returns true if RegExpFlag == 1
- HasDurationFilter(): Returns true if ChkDurationMin > 0 or ChkDurationMax > 0
- ChannelCount(): Returns len(ServiceList)

**Example**:
```go
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

// Helper methods
func (s *SearchSettings) IsEnabled() bool {
    return s.DisableFlag == 0
}

func (s *SearchSettings) IsRegex() bool {
    return s.RegExpFlag == 1
}

func (s *SearchSettings) HasDurationFilter() bool {
    return s.ChkDurationMin > 0 || s.ChkDurationMax > 0
}

func (s *SearchSettings) ChannelCount() int {
    return len(s.ServiceList)
}
```

---

### ServiceInfo

**Purpose**: Represents a broadcast channel identifier in ONID-TSID-SID format

**Fields**:
| Field | Type | XML/JSON Tag | Description | Validation |
|-------|------|--------------|-------------|------------|
| ONID | int | `xml:"onid" json:"onid"` | Original Network ID | > 0 |
| TSID | int | `xml:"tsid" json:"tsid"` | Transport Stream ID | > 0 |
| SID | int | `xml:"sid" json:"sid"` | Service ID | > 0 |

**Relationships**:
- Belongs to SearchSettings
- Can be referenced by existing channel.go utilities

**Display Representation**:
- String format: "ONID-TSID-SID" (e.g., "32736-32736-1024")
- Used in CSV/TSV output
- Filterable via --channel flag

**Validation Rules**:
- All IDs must be positive integers
- Represents a unique channel in Japanese digital broadcasting

**Business Logic**:
- String(): Returns "ONID-TSID-SID" formatted string
- Matches(onid, tsid, sid): Returns true if all IDs match

**Example**:
```go
type ServiceInfo struct {
    ONID int `xml:"onid" json:"onid"`
    TSID int `xml:"tsid" json:"tsid"`
    SID  int `xml:"sid" json:"sid"`
}

// Helper methods
func (s *ServiceInfo) String() string {
    return fmt.Sprintf("%d-%d-%d", s.ONID, s.TSID, s.SID)
}

func (s *ServiceInfo) Matches(onid, tsid, sid int) bool {
    return s.ONID == onid && s.TSID == tsid && s.SID == sid
}
```

---

### RecordingSettings

**Purpose**: Defines recording behavior and post-processing options

**Fields**:
| Field | Type | XML/JSON Tag | Description | Validation | Default |
|-------|------|--------------|-------------|------------|---------|
| RecMode | int | `xml:"recMode" json:"rec_mode"` | Recording mode (1=standard) | >= 0 | 1 |
| Priority | int | `xml:"priority" json:"priority"` | Recording priority | >= 0 | 2 |
| TuijyuuFlag | int | `xml:"tuijyuuFlag" json:"auto_follow"` | 1=auto-follow enabled | 0 or 1 | 1 |
| ServiceMode | int | `xml:"serviceMode" json:"service_mode"` | Service mode | >= 0 | 16 |
| PittariFlag | int | `xml:"pittariFlag" json:"exact_match"` | Exact time match | 0 or 1 | 0 |
| BatFilePath | string | `xml:"batFilePath" json:"bat_file"` | Post-recording batch file path | - | "" |
| RecFolderList | string | `xml:"recFolderList" json:"rec_folders"` | Recording folder list | - | "" |
| SuspendMode | int | `xml:"suspendMode" json:"suspend_mode"` | System suspend mode | >= 0 | 0 |
| DefServiceMode | int | `xml:"defserviceMode" json:"def_service_mode"` | Default service mode | >= 0 | 17 |
| RebootFlag | int | `xml:"rebootFlag" json:"reboot"` | Reboot after recording | 0 or 1 | 0 |
| UseMargineFlag | int | `xml:"useMargineFlag" json:"use_margin"` | Use recording margins | 0 or 1 | 0 |
| StartMargine | int | `xml:"startMargine" json:"start_margin"` | Start margin (seconds) | >= 0 | 20 |
| EndMargine | int | `xml:"endMargine" json:"end_margin"` | End margin (seconds) | >= 0 | 2 |
| ContinueRecFlag | int | `xml:"continueRecFlag" json:"continue_rec"` | Continue recording flag | 0 or 1 | 0 |
| PartialRecFlag | int | `xml:"partialRecFlag" json:"partial_rec"` | Partial recording flag | 0 or 1 | 0 |
| TunerID | int | `xml:"tunerID" json:"tuner_id"` | Tuner ID (0=auto) | >= 0 | 0 |
| PartialRecFolder | string | `xml:"partialRecFolder" json:"partial_rec_folder"` | Partial recording folder | - | "" |

**Relationships**:
- Belongs to AutoAddRule

**Validation Rules**:
- Priority typically ranges 1-5 (higher = more important)
- Margins should be reasonable (< 600 seconds / 10 minutes)
- If UseMargineFlag == 0, StartMargine and EndMargine are ignored

**Business Logic**:
- IsAutoFollow(): Returns true if TuijyuuFlag == 1
- HasMargins(): Returns true if UseMargineFlag == 1
- UsesCustomTuner(): Returns true if TunerID > 0

**Example**:
```go
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

// Helper methods
func (r *RecordingSettings) IsAutoFollow() bool {
    return r.TuijyuuFlag == 1
}

func (r *RecordingSettings) HasMargins() bool {
    return r.UseMargineFlag == 1
}

func (r *RecordingSettings) UsesCustomTuner() bool {
    return r.TunerID > 0
}
```

---

## Filter Criteria (Client-Side)

### FilterOptions

**Purpose**: User-specified filtering criteria for displaying rules

**Fields**:
| Field | Type | Description |
|-------|------|-------------|
| AndKeyFilter | string | Filter by AndKey substring (case-insensitive) |
| ChannelFilter | string | Filter by channel "ONID-TSID-SID" format |
| EnabledOnly | bool | Show only enabled rules (DisableFlag==0) |
| DisabledOnly | bool | Show only disabled rules (DisableFlag==1) |
| RegexOnly | bool | Show only regex rules (RegExpFlag==1) |

**Validation Rules**:
- EnabledOnly and DisabledOnly cannot both be true
- ChannelFilter must match "ONID-TSID-SID" format if provided

**Business Logic**:
- Matches(rule): Returns true if rule passes all active filters
- All filters are AND-ed together

**Example**:
```go
type FilterOptions struct {
    AndKeyFilter  string
    ChannelFilter string
    EnabledOnly   bool
    DisabledOnly  bool
    RegexOnly     bool
}

func (f *FilterOptions) Matches(rule *AutoAddRule) bool {
    // Check enabled/disabled
    if f.EnabledOnly && !rule.SearchSettings.IsEnabled() {
        return false
    }
    if f.DisabledOnly && rule.SearchSettings.IsEnabled() {
        return false
    }

    // Check regex
    if f.RegexOnly && !rule.SearchSettings.IsRegex() {
        return false
    }

    // Check AndKey substring
    if f.AndKeyFilter != "" {
        andKeyLower := strings.ToLower(rule.SearchSettings.AndKey)
        filterLower := strings.ToLower(f.AndKeyFilter)
        if !strings.Contains(andKeyLower, filterLower) {
            return false
        }
    }

    // Check channel
    if f.ChannelFilter != "" {
        channelMatch := false
        for _, channel := range rule.SearchSettings.ServiceList {
            if channel.String() == f.ChannelFilter {
                channelMatch = true
                break
            }
        }
        if !channelMatch {
            return false
        }
    }

    return true
}
```

---

## Export Formats

### CSV/TSV Flattening

**Approach**: Flatten nested structures into single row with dot notation

**Column Headers**:
```
id, enabled, and_key, not_key, regex, title_only, duration_min, duration_max, channel_count, channels, rec_mode, priority, auto_follow
```

**Channel Representation**: Comma-separated string (quoted in CSV, raw in TSV)
- Example: "32736-32736-1024,32736-32736-1025,32737-32737-1032"

**Boolean Representation**: "true"/"false" strings (easier for spreadsheets than 0/1)

**Example Row**:
```csv
1,true,"サイエンスZERO","[再]",false,false,0,0,62,"32736-32736-1024,32736-32736-1025,...",1,2,true
```

### JSON Structure

**Approach**: Preserve full nested structure with descriptive field names

**Example**:
```json
{
  "id": 1,
  "search": {
    "disabled": 0,
    "case_sensitive": 0,
    "and_key": "サイエンスZERO",
    "not_key": "[再]",
    "regex": 0,
    "title_only": 0,
    "fuzzy": 0,
    "not_content": 0,
    "not_date": 0,
    "free_ca": 0,
    "check_rec_end": 0,
    "rec_day_mask": 6,
    "check_no_service": 0,
    "duration_min": 0,
    "duration_max": 0,
    "channels": [
      {"onid": 32736, "tsid": 32736, "sid": 1024},
      {"onid": 32736, "tsid": 32736, "sid": 1025}
    ]
  },
  "recording": {
    "rec_mode": 1,
    "priority": 2,
    "auto_follow": 1,
    ...
  }
}
```

---

## Implementation Files

### New Files

- `internal/models/reservation.go`: EnumAutoAddResponse, AutoAddRule, SearchSettings, RecordingSettings, ServiceInfo
- `internal/models/filter.go`: FilterOptions with Matches() logic
- `internal/formatters/formatter.go`: Formatter interface
- `internal/formatters/table.go`: Table formatter implementation
- `internal/formatters/json.go`: JSON formatter implementation
- `internal/formatters/csv.go`: CSV formatter implementation
- `internal/formatters/tsv.go`: TSV formatter implementation

### Modified Files

- `internal/models/channel.go`: May reuse ParseServiceListEntry() for channel parsing (if applicable)

---

## Validation Summary

### Mandatory Validations

1. ✅ XML parsing: Verify structure matches EMWUI response
2. ✅ ID > 0: All rules must have positive IDs
3. ✅ Channel format: ONID, TSID, SID must be positive integers
4. ✅ Flag values: All flag fields must be 0 or 1
5. ✅ Duration logic: If both min/max set, min < max

### Optional Validations (Warnings)

- AndKey and NotKey both empty (rule may match everything)
- ServiceList empty (searches all channels, may be unintended)
- Very large channel counts (>100 channels, performance concern)

---

## Testing Considerations

### Unit Tests

- Parse valid XML response
- Parse empty response (total=0, count=0)
- Handle malformed XML gracefully
- Filter matching logic (all combinations)
- CSV/TSV flattening with edge cases (empty strings, commas in fields)
- JSON marshaling/unmarshaling round-trip

### Integration Tests

- Retrieve and parse real EMWUI response
- Filter and export to all formats
- Verify Japanese character preservation in all formats
- Verify large rule sets (230+ rules)

### Edge Cases

- Rule with no channels (ServiceList empty)
- Rule with 100+ channels
- AndKey/NotKey with special characters (quotes, commas, tabs)
- Zero-length AndKey/NotKey
- Negative or invalid flag values (how to handle?)
