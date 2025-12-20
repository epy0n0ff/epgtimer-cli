# EMWUI API Contract: Add Reservation

**Feature**: `001-epgtimer-cli`
**Created**: 2025-12-20
**API Version**: EDCB EMWUI (xtne6f/EDCB fork)

---

## Overview

This document specifies the contract for the EMWUI AddReserveManual.lua endpoint, which is used to create manual time-based EPG recording reservations.

**Endpoint Purpose**: Add a new recording reservation by specifying channel, date/time, duration, and program title.

**Protocol**: HTTP GET with query string parameters
**Authentication**: None (trusted network)
**Response Format**: Plain text

---

## Endpoint Specification

### Base URL

```
http://<hostname>:<port>
```

**Default Port**: 5510
**Configuration**: Set via `EMWUI_ENDPOINT` environment variable

**Examples**:
- `http://localhost:5510`
- `http://192.168.1.100:5510`
- `http://epgtimer.local:5510`

### Endpoint Path

```
GET /api/AddReserveManual.lua
```

---

## Request Parameters

All parameters are passed as URL query string parameters.

| Parameter | Type | Required | Format | Description | Example |
|-----------|------|----------|--------|-------------|---------|
| `sid` | uint16 | Yes | Numeric string | Service ID (channel identifier) | `1024` |
| `startDate` | string | Yes | `YYYYMMDD` | Start date in JST | `20251220` |
| `startTime` | string | Yes | `HHMMSS` | Start time in JST (24-hour format) | `190000` |
| `duration` | uint32 | Yes | Seconds | Duration in seconds | `1800` (30 min) |
| `title` | string | Yes | URL-encoded UTF-8 | Program title | `NHK%E3%83%8B%E3%83%A5%E3%83%BC%E3%82%B9` |

### Parameter Details

#### 1. `sid` (Service ID)

- **Type**: Unsigned 16-bit integer as string
- **Range**: Typically 1000-9999 for Japanese digital TV
- **Purpose**: Identifies the broadcast channel/service
- **Validation**: Must be non-zero, no format validation by CLI (EMWUI will reject if invalid)
- **Examples**:
  - `1024` - NHK総合 (terrestrial)
  - `101` - BS1 (satellite)

#### 2. `startDate` (Start Date)

- **Type**: String
- **Format**: `YYYYMMDD` (8 digits)
- **Timezone**: Japan Standard Time (JST, UTC+9)
- **Validation**: Must be a valid date, typically in the future
- **Examples**:
  - `20251220` - December 20, 2025
  - `20260101` - January 1, 2026

#### 3. `startTime` (Start Time)

- **Type**: String
- **Format**: `HHMMSS` (6 digits, 24-hour format)
- **Timezone**: Japan Standard Time (JST, UTC+9)
- **Range**: `000000` to `235959`
- **Examples**:
  - `190000` - 7:00 PM
  - `083000` - 8:30 AM
  - `003000` - 12:30 AM

#### 4. `duration` (Duration)

- **Type**: Unsigned 32-bit integer as string
- **Unit**: Seconds
- **Range**: Practical range 60-43200 (1 minute to 12 hours)
- **Conversion**: User provides minutes, CLI converts to seconds
- **Examples**:
  - `1800` - 30 minutes
  - `3600` - 1 hour
  - `5400` - 90 minutes

#### 5. `title` (Program Title)

- **Type**: String
- **Encoding**: UTF-8, URL-encoded
- **Max Length**: 255 characters (before encoding)
- **Special Characters**: Allowed (Japanese, symbols, spaces)
- **Encoding Method**: Use `net/url.QueryEscape()` in Go
- **Examples**:
  - Raw: `NHKニュース7`
  - Encoded: `NHK%E3%83%8B%E3%83%A5%E3%83%BC%E3%82%B97`
  - Raw: `ドラマSP 特別編`
  - Encoded: `%E3%83%89%E3%83%A9%E3%83%9ESP%20%E7%89%B9%E5%88%A5%E7%B7%A8`

---

## Request Examples

### Example 1: Simple Reservation

**Scenario**: Record "NHKニュース7" on channel 1024, starting Dec 20, 2025 at 7:00 PM JST for 30 minutes.

**HTTP Request**:
```http
GET /api/AddReserveManual.lua?sid=1024&startDate=20251220&startTime=190000&duration=1800&title=NHK%E3%83%8B%E3%83%A5%E3%83%BC%E3%82%B97 HTTP/1.1
Host: localhost:5510
User-Agent: epgtimer-cli/1.0
Accept: */*
```

**cURL Command**:
```bash
curl "http://localhost:5510/api/AddReserveManual.lua?sid=1024&startDate=20251220&startTime=190000&duration=1800&title=NHK%E3%83%8B%E3%83%A5%E3%83%BC%E3%82%B97"
```

### Example 2: Late Night Program

**Scenario**: Record a late-night show starting at 1:30 AM for 60 minutes.

**HTTP Request**:
```http
GET /api/AddReserveManual.lua?sid=1032&startDate=20251221&startTime=013000&duration=3600&title=%E6%B7%B1%E5%A4%9C%E3%83%89%E3%83%A9%E3%83%9E HTTP/1.1
Host: localhost:5510
```

### Example 3: Long Duration Recording

**Scenario**: Record a 2-hour special program.

**HTTP Request**:
```http
GET /api/AddReserveManual.lua?sid=1048&startDate=20251225&startTime=200000&duration=7200&title=%E7%89%B9%E5%88%A5%E7%95%AA%E7%B5%84 HTTP/1.1
Host: localhost:5510
```

---

## Response Specification

### Success Response

**HTTP Status**: `200 OK`
**Content-Type**: `text/plain; charset=UTF-8`
**Body**: Plain text success message

**Success Indicators**:
- Body contains: `OK` (ASCII)
- Or body contains: `追加しました` (Japanese: "Added")
- HTTP status is 200

**Example Success Response**:
```http
HTTP/1.1 200 OK
Content-Type: text/plain; charset=UTF-8
Content-Length: 2

OK
```

### Error Responses

EMWUI uses inconsistent error reporting. Errors may be indicated by:
1. HTTP status codes (500, 404)
2. `200 OK` with error message in body
3. Empty response body

#### Error Pattern 1: Application Error (NG)

**HTTP Status**: `200 OK`
**Body**: `NG` or Japanese error message

**Example**:
```http
HTTP/1.1 200 OK
Content-Type: text/plain; charset=UTF-8
Content-Length: 2

NG
```

**Meaning**: Request was rejected (invalid parameters, duplicate reservation, or business logic error)

####Error Pattern 2: Server Error

**HTTP Status**: `500 Internal Server Error`
**Body**: Error message or empty

**Example**:
```http
HTTP/1.1 500 Internal Server Error
Content-Type: text/plain; charset=UTF-8
Content-Length: 45

Lua script error: invalid date format
```

**Meaning**: Internal EMWUI/Lua script error

#### Error Pattern 3: Not Found

**HTTP Status**: `404 Not Found`
**Body**: HTML error page or empty

**Meaning**: Endpoint doesn't exist (wrong EMWUI version or URL)

#### Error Pattern 4: Connection Error

**HTTP Status**: N/A (no response)
**Error**: Network timeout or connection refused

**Meaning**: EMWUI service not running or wrong endpoint URL

---

## Error Detection Logic

Since EMWUI doesn't use standard error codes consistently, the CLI must implement custom error detection:

```go
func ParseEMWUIResponse(httpStatus int, body string) *AddReservationResponse {
	resp := &AddReservationResponse{
		HTTPStatus: httpStatus,
		Message:    strings.TrimSpace(body),
	}

	// Check HTTP status first
	if httpStatus == 0 {
		// Connection error (no response)
		resp.Success = false
		return resp
	}

	if httpStatus >= 400 {
		// HTTP error status
		resp.Success = false
		return resp
	}

	// HTTP 200 - check body content
	if httpStatus == 200 {
		// Success patterns
		if body == "OK" || strings.Contains(body, "追加しました") {
			resp.Success = true
			return resp
		}

		// Error patterns
		if body == "NG" || body == "" {
			resp.Success = false
			return resp
		}

		// Assume success if no clear error indicator
		resp.Success = true
		return resp
	}

	// Default to error for unknown status
	resp.Success = false
	return resp
}
```

---

## Common Error Scenarios

| Scenario | HTTP Status | Response Body | CLI Error Message |
|----------|-------------|---------------|-------------------|
| Connection refused | 0 | N/A | "Cannot connect to EMWUI service at {URL}. Ensure EpgTimer is running." |
| Endpoint not found | 404 | HTML or empty | "EMWUI API endpoint not found. Check EMWUI version compatibility." |
| Invalid parameters | 200 | "NG" | "Request rejected by EMWUI: invalid parameters or duplicate reservation" |
| Server error | 500 | Error message | "EMWUI service error (HTTP 500): {message}" |
| Duplicate reservation | 200 | "NG" or error msg | "Reservation already exists for this program" |
| Past start time | 200 | "NG" or error msg | "Start time must be in the future" |
| Invalid channel ID | 200 | "NG" or error msg | "Invalid channel ID: channel does not exist" |

---

## Client Implementation Guidelines

### 1. HTTP Client Configuration

```go
client := &http.Client{
	Timeout: 10 * time.Second,
	Transport: &http.Transport{
		MaxIdleConns:       10,
		IdleConnTimeout:    30 * time.Second,
		DisableCompression: false,
	},
}
```

### 2. Building Request URL

```go
func BuildAddReservationURL(baseURL string, req *AddReservationRequest) (string, error) {
	// Validate base URL
	if baseURL == "" {
		return "", fmt.Errorf("EMWUI_ENDPOINT not configured")
	}

	// Ensure base URL doesn't end with slash
	baseURL = strings.TrimSuffix(baseURL, "/")

	// Convert request to EMWUI parameters
	params := req.ToEMWUIParams()

	// Build query string
	v := url.Values{}
	v.Set("sid", params["sid"])
	v.Set("startDate", params["startDate"])
	v.Set("startTime", params["startTime"])
	v.Set("duration", params["duration"])
	v.Set("title", params["title"]) // QueryEscape is applied automatically

	// Build full URL
	fullURL := fmt.Sprintf("%s/api/AddReserveManual.lua?%s", baseURL, v.Encode())

	return fullURL, nil
}
```

### 3. Making Request

```go
func AddReservation(client *http.Client, url string) (*AddReservationResponse, error) {
	// Make HTTP GET request
	resp, err := client.Get(url)
	if err != nil {
		// Connection error
		return &AddReservationResponse{
			Success:    false,
			Message:    err.Error(),
			HTTPStatus: 0,
		}, fmt.Errorf("connection error: %w", err)
	}
	defer resp.Body.Close()

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	// Parse response
	return ParseEMWUIResponse(resp.StatusCode, string(body)), nil
}
```

---

## Testing

### Mock Server Responses

Create test fixtures for integration tests:

**Success Response**:
```go
func mockSuccessHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}
```

**Error Response (NG)**:
```go
func mockErrorHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("NG"))
}
```

**Server Error**:
```go
func mockServerErrorHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte("Internal error"))
}
```

### Test Cases

1. **Happy Path**: Valid parameters, successful reservation
2. **Invalid Channel**: Non-existent channel ID
3. **Duplicate Reservation**: Same program/time already reserved
4. **Past Time**: Start time in the past
5. **Connection Failure**: EMWUI service not running
6. **Server Error**: EMWUI returns 500
7. **Endpoint Not Found**: Wrong URL or EMWUI version
8. **Japanese Characters**: UTF-8 encoding and URL encoding
9. **Special Characters**: Titles with symbols, spaces, punctuation
10. **Boundary Values**: Min/max duration, edge case times (23:59, 00:00)

---

## Reference Implementation

See xtne6f/EDCB repository:
- Lua script: `EDCB/EpgTimerSrv/Setting/HttpPublic/api/AddReserveManual.lua`
- Web UI client: `EDCB/EpgTimerSrv/Setting/HttpPublic/js/reserve.js`

---

## Versioning & Compatibility

**Current Status**: EMWUI has no formal API versioning

**Compatibility Strategy**:
1. Attempt API call
2. Parse response
3. If parsing fails, log warning about potential version mismatch
4. Provide clear error message to user

**Known Variations**:
- Older versions may use Shift-JIS encoding instead of UTF-8
- Response messages may vary between EDCB forks
- Field order and formats are generally stable

---

**End of API Contract Document**

Last Updated: 2025-12-20
