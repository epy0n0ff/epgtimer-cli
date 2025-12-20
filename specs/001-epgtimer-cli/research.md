# EMWUI API Research: EpgTimer EMWUI (EpgDataCap_Bon EDCB)

**Research Date**: 2025-12-20
**Feature**: `001-epgtimer-cli`
**Status**: Initial Research

---

## Executive Summary

EMWUI (EpgDataCap_Bon Multi-platform Web UI) is a web interface for EDCB (EpgDataCap_Bon), a Japanese digital TV recording system. The EMWUI API is a **custom HTTP-based API** that uses **Lua-based endpoints** with **query string parameters** and returns **plain text or custom-formatted responses** (not REST/JSON).

**Key Findings**:
- Protocol: HTTP with custom Lua CGI endpoints
- Format: NOT REST/JSON - uses query strings and plain text responses
- Base Path: `/api/` with Lua script endpoints
- Authentication: None (trusted network)
- Character Encoding: UTF-8 (Japanese support)
- Primary Implementation: xtne6f/EDCB fork on GitHub

---

## 1. EMWUI API Protocol & Architecture

### 1.1 Protocol Type

**EMWUI uses a custom HTTP API protocol**, not standard REST/JSON:

- **Base URL Pattern**: `http://<host>:<port>/api/<endpoint>.lua?<params>`
- **Request Method**: Primarily GET requests with query string parameters
- **Response Format**: Plain text, custom delimited format (not JSON/XML)
- **Default Port**: Typically 5510 (configurable)

### 1.2 Technology Stack

- **Web Server**: Built-in Lua-based HTTP server (part of EDCB)
- **API Scripts**: Lua CGI scripts in `EDCB/EpgTimerSrv/Setting/HttpPublic/api/`
- **Data Format**: Custom pipe-delimited or tab-delimited text format
- **Character Encoding**: UTF-8 (Shift-JIS in older versions)

### 1.3 Connection Requirements

```
Base URL: http://localhost:5510
Authentication: None (assumes trusted local network)
Headers: None required
Timeout: 5-10 seconds recommended
```

---

## 2. Core API Endpoints

### 2.1 List Reservations Endpoint

**Endpoint**: `/api/EnumReserve.lua`

**Purpose**: Retrieve list of all EPG recording reservations

**Request**:
```
GET /api/EnumReserve.lua HTTP/1.1
Host: localhost:5510
```

**Query Parameters**:
- None required for basic list
- `type=0` - All reservations (default)
- `type=1` - Only enabled reservations
- `type=2` - Only disabled reservations

**Response Format**:
```
Plain text, tab-delimited format
Each line represents one reservation with fields separated by tabs

Format (approximate):
<ReserveID>\t<Title>\t<StartTime>\t<Duration>\t<ServiceName>\t<NetworkName>\t<RecMode>\t<RecSetting>\t<Comment>\t...

Example:
1\tNHKニュース7\t2025-12-20 19:00:00\t1800\tNHK総合\tNHK\t1\t0\t\t...
2\tドラマSP\t2025-12-20 21:00:00\t3600\tTBS\tTBS\t1\t0\t\t...
```

**Response Status**:
- Success: HTTP 200 with text data
- No reservations: HTTP 200 with empty body
- Error: HTTP 500 or empty response

### 2.2 Add Reservation Endpoint

**Endpoint**: `/api/AddReserve.lua`

**Purpose**: Create a new EPG recording reservation

**Request**:
```
GET /api/AddReserve.lua?onid=<NetworkID>&tsid=<TSID>&sid=<ServiceID>&eid=<EventID>&tuijyu=<追従設定> HTTP/1.1
Host: localhost:5510
```

**Query Parameters** (Critical):
- `onid` - Original Network ID (required)
- `tsid` - Transport Stream ID (required)
- `sid` - Service ID (channel identifier) (required)
- `eid` - Event ID (program identifier from EPG) (required)
- `tuijyu` - 追従設定 (follow/auto-update mode): 0=none, 1=auto-adjust

**Additional Parameters** (Optional):
- `recMode` - Recording mode (0-9, encoding quality)
- `priority` - Priority level (1-5)
- `marginStart` - Start margin in seconds (pre-recording)
- `marginEnd` - End margin in seconds (post-recording)

**Response Format**:
```
Plain text status message

Success: "OK" or "追加しました" (Added)
Error: "NG" or error message in Japanese
```

**Response Status**:
- Success: HTTP 200 with "OK"
- Duplicate: HTTP 200 with error message
- Invalid params: HTTP 400 or error message
- Server error: HTTP 500

### 2.3 Alternative: Manual Reservation (Time-based)

**Endpoint**: `/api/AddReserveManual.lua`

**Purpose**: Create manual time-based reservation (not EPG-based)

**Request**:
```
GET /api/AddReserveManual.lua?sid=<ServiceID>&startDate=<YYYYMMDD>&startTime=<HHMMSS>&duration=<seconds>&title=<Title> HTTP/1.1
```

**Query Parameters**:
- `sid` - Service ID (channel)
- `startDate` - Start date (YYYYMMDD format)
- `startTime` - Start time (HHMMSS format)
- `duration` - Duration in seconds
- `title` - Program title (URL-encoded UTF-8)

**Note**: This is more suitable for CLI usage when EPG event IDs are not available.

---

## 3. Data Models & Entity Structure

### 3.1 Reservation Entity

Based on EDCB reservation structure:

```go
type Reservation struct {
    ReserveID       uint32    // Unique reservation ID
    Title           string    // Program title
    StartTime       time.Time // Scheduled start time
    DurationSec     uint32    // Duration in seconds
    ServiceName     string    // Channel name (e.g., "NHK総合")
    NetworkName     string    // Network name
    OriginalNetworkID uint16  // ONID (for EPG identification)
    TransportStreamID uint16  // TSID
    ServiceID       uint16    // SID (channel ID)
    EventID         uint16    // EID (program ID from EPG)
    RecMode         uint8     // Recording mode (quality)
    RecSetting      uint32    // Recording settings flags
    Priority        uint8     // Priority (1-5)
    Status          ReserveStatus // Current status
    Comment         string    // User comment
    IsEnabled       bool      // Is reservation active
}

type ReserveStatus uint8

const (
    StatusScheduled ReserveStatus = iota
    StatusRecording
    StatusCompleted
    StatusFailed
    StatusDisabled
)
```

### 3.2 Service (Channel) Entity

```go
type Service struct {
    OriginalNetworkID uint16 // ONID
    TransportStreamID uint16 // TSID
    ServiceID         uint16 // SID
    ServiceName       string // Display name
    NetworkName       string // Network/broadcaster
    ServiceType       uint16 // Service type (0x01=TV, 0x02=Radio)
}
```

---

## 4. Error Handling

### 4.1 Error Response Format

EMWUI does not use standard HTTP error codes consistently:

**Common Error Patterns**:
- Connection refused: TCP connection error (server not running)
- HTTP 500: Internal Lua script error
- HTTP 200 + "NG": Operation failed (application-level error)
- HTTP 200 + Empty: No data available or error
- HTTP 200 + Japanese error message: Specific error description

### 4.2 Error Detection Strategy

Since responses are not structured JSON, errors must be detected by:

1. **HTTP Status Code**: Check for non-200 status
2. **Response Body Parsing**: Check for "NG" or error keywords
3. **Empty Response**: Treat as potential error or no data
4. **Timeout**: Network/server issues
5. **Malformed Response**: Parsing errors indicate API changes

### 4.3 Common Error Scenarios

| Scenario | Detection | User Message |
|----------|-----------|--------------|
| Server not running | Connection refused | "Cannot connect to EMWUI service at {URL}. Ensure EpgTimer is running." |
| Invalid endpoint | HTTP 404 | "EMWUI API endpoint not found. Check EMWUI version compatibility." |
| Invalid parameters | "NG" response or HTTP 400 | "Invalid parameters: {details}" |
| Duplicate reservation | "NG" + message | "Reservation already exists for this program" |
| EPG data not available | Empty or error response | "EPG data not available for specified program" |

---

## 5. API Versioning

### 5.1 Current Status

**EMWUI does NOT have formal API versioning**:

- No version numbers in URL paths
- No version headers
- API changes occur with EDCB releases
- Breaking changes possible between versions

### 5.2 Version Detection Strategy

**Recommended approach**:
1. Attempt API call to known endpoint
2. Parse response format
3. If parsing fails, log warning about possible version mismatch
4. Provide clear error message to user

**EDCB Version References**:
- Original: xtne6f/EDCB (active development)
- Fork history: EpgDataCap_Bon → xtne6f/EDCB → various forks
- EMWUI versions vary by fork

### 5.3 Compatibility Considerations

```go
// Recommended version detection
type APIVersion struct {
    ServerVersion string // From response headers or detection
    Compatible    bool   // Can this client work with server?
}

// Check compatibility by attempting simple API call
func DetectAPIVersion(baseURL string) (*APIVersion, error) {
    // Try EnumReserve endpoint
    resp, err := http.Get(baseURL + "/api/EnumReserve.lua")
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    // If we get 200 OK, assume compatible
    if resp.StatusCode == 200 {
        return &APIVersion{Compatible: true}, nil
    }

    return &APIVersion{Compatible: false}, fmt.Errorf("incompatible API")
}
```

---

## 6. Official Resources & Community

### 6.1 GitHub Repositories

**Primary Repository**:
- **URL**: https://github.com/xtne6f/EDCB
- **Description**: Main active fork of EpgDataCap_Bon (EDCB)
- **Components**:
  - `EDCB/EpgTimerSrv/` - Main EPG timer service
  - `EDCB/EpgTimerSrv/Setting/HttpPublic/` - Web UI files
  - `EDCB/EpgTimerSrv/Setting/HttpPublic/api/` - Lua API scripts
- **API Scripts**: Located in `HttpPublic/api/*.lua`
  - `EnumReserve.lua` - List reservations
  - `AddReserve.lua` - Add EPG-based reservation
  - `AddReserveManual.lua` - Add manual reservation
  - `DelReserve.lua` - Delete reservation
  - `ChgReserve.lua` - Modify reservation
  - `EnumService.lua` - List channels/services
  - `EnumEventInfo.lua` - List EPG programs

**Alternative/Related Repositories**:
- Original (archived): https://github.com/tkmsst/EDCB (historical reference)
- EMWUI-specific: Look for `/Setting/HttpPublic/` directory in EDCB forks

### 6.2 Documentation Sources

**Official Documentation**:
- **Location**: Primarily in Japanese within EDCB repository
- **README**: `EDCB/Document/` directory (mostly Japanese)
- **API Docs**: Limited - primarily read Lua source code
- **Wiki**: Some EDCB forks have community wikis (Japanese)

**Community Resources**:
- 5ch (2ch) boards: Japanese BBS discussions on EDCB/EpgTimer
- Personal blogs: Japanese tech blogs with EDCB setup guides
- GitHub Issues: xtne6f/EDCB issues contain API usage examples

### 6.3 Example Client Implementations

**Known Clients**:
1. **EMWUI Web Interface** (JavaScript)
   - Location: `EDCB/EpgTimerSrv/Setting/HttpPublic/`
   - Files: `default.html`, `js/reserve.js`
   - Shows actual API usage patterns

2. **EpgTimer (Windows GUI)**
   - C++ implementation within EDCB
   - Uses native APIs, not HTTP

3. **Third-party scripts**
   - Community PowerShell/Python scripts
   - Search GitHub for "EDCB API" or "EpgTimer script"

**Recommended Reference**:
```
Read Lua API scripts directly:
https://github.com/xtne6f/EDCB/tree/master/EDCB/EpgTimerSrv/Setting/HttpPublic/api

These are the definitive source of truth for API behavior.
```

---

## 7. Parsing Strategy & Implementation Notes

### 7.1 Response Parsing Challenges

**Challenge 1: Tab-Delimited Format**
- Fields separated by `\t` (tab character)
- Field order may vary by EDCB version
- Some fields may be empty (consecutive tabs)
- Last fields may be omitted

**Solution**:
```go
// Use robust tab-split parser
fields := strings.Split(line, "\t")
// Check field count before accessing
if len(fields) < expectedCount {
    return fmt.Errorf("unexpected field count: got %d, want >= %d", len(fields), expectedCount)
}
```

**Challenge 2: Date/Time Parsing**
- Format: `YYYY-MM-DD HH:MM:SS` (UTF-8)
- Timezone: Typically Japan Standard Time (JST, UTC+9)
- May need to handle local timezone conversion

**Solution**:
```go
// Parse with location
loc, _ := time.LoadLocation("Asia/Tokyo")
startTime, err := time.ParseInLocation("2006-01-02 15:04:05", timeStr, loc)
```

**Challenge 3: Japanese Character Handling**
- UTF-8 encoding required
- Terminal width calculation with multi-byte chars
- Use `golang.org/x/text/width` for East Asian Width

### 7.2 HTTP Client Configuration

```go
// Recommended HTTP client setup
client := &http.Client{
    Timeout: 10 * time.Second,
    Transport: &http.Transport{
        MaxIdleConns:        10,
        IdleConnTimeout:     30 * time.Second,
        DisableCompression:  false,
    },
}

// Always check Content-Type (may be text/plain or text/html)
resp, err := client.Get(url)
if err != nil {
    return fmt.Errorf("connection error: %w", err)
}
defer resp.Body.Close()

// Read response body
body, err := io.ReadAll(resp.Body)
if err != nil {
    return fmt.Errorf("read error: %w", err)
}

// Convert to UTF-8 string
content := string(body)
```

### 7.3 URL Encoding for Japanese Text

When adding reservations with Japanese titles:

```go
import "net/url"

// Encode title parameter
title := "ニュース番組"
encodedTitle := url.QueryEscape(title)

// Build URL
apiURL := fmt.Sprintf("%s/api/AddReserveManual.lua?title=%s&...", baseURL, encodedTitle)
```

---

## 8. Alternative Approaches Considered

### 8.1 TCP/IP Native Protocol

EDCB also has a **native TCP/IP protocol** on port 4510 (default):

**Pros**:
- More efficient binary protocol
- Richer data structures
- Bidirectional communication

**Cons**:
- Complex binary format (requires packet parsing)
- Not well documented
- Requires understanding of EDCB internal structures

**Decision**: Use HTTP/EMWUI API for simplicity and maintainability.

### 8.2 Direct Database Access

EDCB stores data in:
- `EpgTimer.db` (SQLite database)
- Various `.dat` files

**Pros**:
- Direct data access
- No network dependency

**Cons**:
- Bypasses EDCB business logic
- Risk of data corruption
- No validation or constraint enforcement
- Requires file system access

**Decision**: Use API layer to maintain data integrity.

---

## 9. Recommended Technology Choices

### 9.1 CLI Framework

**Recommendation**: **Cobra** (github.com/spf13/cobra)

**Rationale**:
- Industry standard for Go CLIs
- Excellent subcommand support
- Auto-generated help text
- Shell completion support
- Wide community adoption
- Small dependency footprint

**Alternative Considered**: Standard `flag` package
- Too basic for multi-command CLI
- Poor help text generation
- Manual subcommand routing

### 9.2 Terminal Output Library

**Recommendation**: **Custom table formatter** (no external dependency)

**Rationale**:
- Simple tab-aligned output sufficient
- Avoid unnecessary dependencies
- Easy to customize for Japanese text
- Standard library `text/tabwriter` adequate

**Alternative Considered**: `tablewriter`
- Adds external dependency
- Overkill for simple CLI output
- Not needed for basic formatting

**Fallback**: Use `golang.org/x/text/width` for accurate width calculation with CJK characters

### 9.3 HTTP Client

**Recommendation**: Standard library `net/http`

**Rationale**:
- No dependencies needed
- Sufficient for simple HTTP GET requests
- Good timeout and error handling
- Well-documented and stable

---

## 10. Testing Strategy

### 10.1 Mock EMWUI Server

Create test HTTP server that mimics EMWUI responses:

```go
// tests/testserver/emwui_mock.go
type MockEMWUI struct {
    Reservations []Reservation
}

func (m *MockEMWUI) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    switch r.URL.Path {
    case "/api/EnumReserve.lua":
        m.handleEnumReserve(w, r)
    case "/api/AddReserve.lua":
        m.handleAddReserve(w, r)
    default:
        http.NotFound(w, r)
    }
}
```

### 10.2 Test Data

Create fixtures with:
- Japanese program titles (UTF-8)
- Various channel types (BS, CS, terrestrial)
- Different reservation statuses
- Edge cases (long titles, special characters)

### 10.3 Integration Tests

```go
// tests/integration/list_test.go
func TestListReservations(t *testing.T) {
    // Start mock server
    mock := testserver.NewMockEMWUI()
    server := httptest.NewServer(mock)
    defer server.Close()

    // Run CLI command
    cmd := exec.Command("./epgtimer", "list", "--endpoint", server.URL)
    output, err := cmd.CombinedOutput()

    // Assert
    assert.NoError(t, err)
    assert.Contains(t, string(output), "NHKニュース")
}
```

---

## 11. Configuration Management

### 11.1 Environment Variables

**Primary Configuration**:
```bash
EMWUI_ENDPOINT=http://localhost:5510
```

**Optional Configuration**:
```bash
EMWUI_TIMEOUT=10          # HTTP timeout in seconds
EMWUI_DEBUG=true          # Enable debug logging
```

### 11.2 Configuration Priority

1. Command-line flags (highest priority)
2. Environment variables
3. Config file (if implemented)
4. Default values (lowest priority)

**Example**:
```bash
# Override endpoint via flag
epgtimer list --endpoint http://192.168.1.100:5510

# Or use environment variable
export EMWUI_ENDPOINT=http://192.168.1.100:5510
epgtimer list
```

---

## 12. Open Questions & Risks

### 12.1 Unresolved Questions

1. **Field Order Stability**: Does tab-delimited field order change between EDCB versions?
   - **Mitigation**: Parse by position, document expected format, add version detection

2. **Status Field Values**: What are exact numeric values for reservation status?
   - **Mitigation**: Read Lua scripts directly, test with live system

3. **Character Encoding Edge Cases**: Are all responses guaranteed UTF-8?
   - **Mitigation**: Detect encoding, add fallback to Shift-JIS if needed

4. **Concurrent Modification**: What happens if reservation is modified during API call?
   - **Mitigation**: Accept eventual consistency, add retry logic

### 12.2 Implementation Risks

| Risk | Impact | Likelihood | Mitigation |
|------|--------|-----------|------------|
| EMWUI API changes without notice | High | Medium | Version detection, graceful degradation |
| Japanese text display issues | Medium | Low | Test on multiple terminals, use width library |
| Network timeout in slow environments | Low | Medium | Configurable timeout, clear error messages |
| Parsing errors from unexpected formats | High | Medium | Robust error handling, log raw responses |

### 12.3 Next Steps

1. **Validate endpoint URLs** against live EDCB instance
2. **Capture actual response examples** for parsing tests
3. **Test with different EDCB versions** to verify compatibility
4. **Document field mappings** from actual API responses
5. **Create comprehensive test fixtures** based on real data

---

## 13. References

### 13.1 Key URLs

- **Primary Repository**: https://github.com/xtne6f/EDCB
- **API Scripts**: https://github.com/xtne6f/EDCB/tree/master/EDCB/EpgTimerSrv/Setting/HttpPublic/api
- **Web UI Code**: https://github.com/xtne6f/EDCB/tree/master/EDCB/EpgTimerSrv/Setting/HttpPublic

### 13.2 Code References

```
Read these Lua files for API implementation details:
- api/EnumReserve.lua      (list reservations)
- api/AddReserve.lua        (add EPG reservation)
- api/AddReserveManual.lua  (add manual reservation)
- api/EnumService.lua       (list channels)
- api/EnumEventInfo.lua     (list EPG programs)
- default.html              (web UI usage examples)
- js/reserve.js             (JavaScript API client)
```

### 13.3 Search Keywords

For finding additional resources:
- "EDCB API"
- "EpgTimer EMWUI"
- "EpgDataCap_Bon API"
- "EDCB Lua script"
- "xtne6f EDCB"

---

## Appendix A: Sample API Call Flows

### A.1 List Reservations Flow

```
1. Client → Server: GET /api/EnumReserve.lua
2. Server processes request via Lua script
3. Server queries internal reservation database
4. Server formats response as tab-delimited text
5. Server → Client: HTTP 200 + text body
6. Client parses tab-delimited lines
7. Client formats output as table
8. Display to user
```

### A.2 Add Reservation Flow (EPG-based)

```
1. User provides: channel name, start time, title
2. Client → Server: GET /api/EnumService.lua (find channel SID)
3. Client → Server: GET /api/EnumEventInfo.lua?sid=<SID> (find event ID)
4. Client extracts ONID, TSID, SID, EID from EPG data
5. Client → Server: GET /api/AddReserve.lua?onid=<>&tsid=<>&sid=<>&eid=<>
6. Server validates parameters
7. Server creates reservation
8. Server → Client: HTTP 200 + "OK"
9. Client displays success message
```

### A.3 Add Reservation Flow (Manual/Time-based)

```
1. User provides: channel ID/name, start date/time, duration, title
2. Client validates parameters (time format, duration > 0)
3. Client → Server: GET /api/EnumService.lua (resolve channel name to SID)
4. Client formats parameters (URL encoding for title)
5. Client → Server: GET /api/AddReserveManual.lua?sid=<>&startDate=<>&...
6. Server validates time range
7. Server creates manual reservation
8. Server → Client: HTTP 200 + "OK" or error message
9. Client displays result
```

---

## Appendix B: Field Reference

### B.1 EnumReserve Response Fields

**Approximate field order** (may vary by version):

| Position | Field Name | Type | Description |
|----------|-----------|------|-------------|
| 0 | ReserveID | uint32 | Unique reservation ID |
| 1 | Title | string | Program title |
| 2 | StartTime | string | Start time (YYYY-MM-DD HH:MM:SS) |
| 3 | DurationSec | uint32 | Duration in seconds |
| 4 | ServiceName | string | Channel name |
| 5 | NetworkName | string | Network/broadcaster |
| 6 | ONID | uint16 | Original Network ID |
| 7 | TSID | uint16 | Transport Stream ID |
| 8 | SID | uint16 | Service ID |
| 9 | EID | uint16 | Event ID |
| 10 | RecMode | uint8 | Recording mode |
| 11 | Priority | uint8 | Priority (1-5) |
| 12 | Status | uint8 | Status code |
| ... | ... | ... | Additional fields vary |

**Note**: Actual field order should be verified by examining Lua script source or capturing real responses.

---

**End of Research Document**

Last Updated: 2025-12-20
