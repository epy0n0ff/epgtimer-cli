# API Contract: EMWUI EnumAutoAdd Endpoint

**Feature**: EPG Reservation List Retrieval (001-epg-reservation-list)
**API Endpoint**: GET /api/EnumAutoAdd
**Version**: EMWUI (xtne6f/EDCB fork, 2024-2025)

## Endpoint Overview

**Base URL**: Configured via EMWUI_ENDPOINT environment variable (e.g., `http://192.168.1.10:5510`)

**Full URL**: `{EMWUI_ENDPOINT}/api/EnumAutoAdd`

**HTTP Method**: GET

**Authentication**: None required

**Content-Type**: application/xml; charset=UTF-8

## Request

### HTTP Method

```
GET /api/EnumAutoAdd HTTP/1.1
Host: 192.168.1.10:5510
```

### Query Parameters

None. API returns all automatic recording rules in a single request.

### Request Headers

No special headers required. Standard HTTP GET request.

### Authentication

No authentication required. EMWUI service runs on trusted local network without access control.

## Response

### Success Response (HTTP 200)

**Content-Type**: `text/xml; charset=UTF-8`

**Status Code**: 200 OK

**Body**: XML document containing all automatic recording rules

**Structure**:
```xml
<?xml version="1.0" encoding="UTF-8" ?>
<entry>
  <total>230</total>
  <index>0</index>
  <count>230</count>
  <items>
    <autoaddinfo>
      <ID>1</ID>
      <searchsetting>
        <disableFlag>0</disableFlag>
        <caseFlag>0</caseFlag>
        <andKey>サイエンスZERO</andKey>
        <notKey>[再]</notKey>
        <regExpFlag>0</regExpFlag>
        <titleOnlyFlag>0</titleOnlyFlag>
        <aimaiFlag>0</aimaiFlag>
        <notContetFlag>0</notContetFlag>
        <notDateFlag>0</notDateFlag>
        <freeCAFlag>0</freeCAFlag>
        <chkRecEnd>0</chkRecEnd>
        <chkRecDay>6</chkRecDay>
        <chkRecNoService>0</chkRecNoService>
        <chkDurationMin>0</chkDurationMin>
        <chkDurationMax>0</chkDurationMax>
        <serviceList>
          <onid>32736</onid>
          <tsid>32736</tsid>
          <sid>1024</sid>
        </serviceList>
        <!-- ... more serviceList entries ... -->
      </searchsetting>
      <recsetting>
        <recMode>1</recMode>
        <priority>2</priority>
        <tuijyuuFlag>1</tuijyuuFlag>
        <serviceMode>16</serviceMode>
        <pittariFlag>0</pittariFlag>
        <batFilePath></batFilePath>
        <recFolderList></recFolderList>
        <suspendMode>0</suspendMode>
        <defserviceMode>17</defserviceMode>
        <rebootFlag>0</rebootFlag>
        <useMargineFlag>0</useMargineFlag>
        <startMargine>20</startMargine>
        <endMargine>2</endMargine>
        <continueRecFlag>0</continueRecFlag>
        <partialRecFlag>0</partialRecFlag>
        <tunerID>0</tunerID>
        <partialRecFolder></partialRecFolder>
      </recsetting>
    </autoaddinfo>
    <!-- ... more autoaddinfo entries ... -->
  </items>
</entry>
```

### Empty Response (HTTP 200)

**Scenario**: No automatic recording rules configured

**Body**:
```xml
<?xml version="1.0" encoding="UTF-8" ?>
<entry>
  <total>0</total>
  <index>0</index>
  <count>0</count>
  <items>
  </items>
</entry>
```

### Error Responses

#### Service Unavailable (HTTP 503)

**Scenario**: EMWUI service is starting up or temporarily unavailable

**Status Code**: 503 Service Unavailable

**Body**: May be empty or contain HTML error page

#### Not Found (HTTP 404)

**Scenario**: API endpoint does not exist (incompatible EMWUI version)

**Status Code**: 404 Not Found

**Body**: HTML error page or empty

#### Internal Server Error (HTTP 500)

**Scenario**: EMWUI encountered an internal error

**Status Code**: 500 Internal Server Error

**Body**: May contain error message or HTML error page

## Data Model

### Root Element: `<entry>`

| Field | Type | Description | Constraints |
|-------|------|-------------|-------------|
| `<total>` | integer | Total number of rules in system | >= 0 |
| `<index>` | integer | Starting index (always 0) | == 0 |
| `<count>` | integer | Number of rules returned | >= 0, <= total |
| `<items>` | container | Container for autoaddinfo elements | Contains 0..* autoaddinfo |

### AutoAddInfo Element: `<autoaddinfo>`

| Field | Type | Description | Constraints |
|-------|------|-------------|-------------|
| `<ID>` | integer | Unique rule identifier | > 0 |
| `<searchsetting>` | container | Search criteria | Required |
| `<recsetting>` | container | Recording settings | Required |

### SearchSetting Element: `<searchsetting>`

| Field | Type | Description | Constraints | Default |
|-------|------|-------------|-------------|---------|
| `<disableFlag>` | integer | 0=enabled, 1=disabled | 0 or 1 | 0 |
| `<caseFlag>` | integer | 0=case-insensitive, 1=case-sensitive | 0 or 1 | 0 |
| `<andKey>` | string | Required search keywords (UTF-8) | Any string | "" |
| `<notKey>` | string | Exclusion keywords (UTF-8) | Any string | "" |
| `<regExpFlag>` | integer | 0=literal, 1=regex pattern | 0 or 1 | 0 |
| `<titleOnlyFlag>` | integer | 0=all fields, 1=title only | 0 or 1 | 0 |
| `<aimaiFlag>` | integer | Fuzzy matching enabled | 0 or 1 | 0 |
| `<notContetFlag>` | integer | Not content flag | 0 or 1 | 0 |
| `<notDateFlag>` | integer | Not date flag | 0 or 1 | 0 |
| `<freeCAFlag>` | integer | 0=include pay, 1=free only | 0 or 1 | 0 |
| `<chkRecEnd>` | integer | Check recording end | >= 0 | 0 |
| `<chkRecDay>` | integer | Recording days bitmask | >= 0 | 6 |
| `<chkRecNoService>` | integer | Check no service | 0 or 1 | 0 |
| `<chkDurationMin>` | integer | Min duration (minutes, 0=none) | >= 0 | 0 |
| `<chkDurationMax>` | integer | Max duration (minutes, 0=none) | >= 0 | 0 |
| `<serviceList>` | container | Channel identifier (repeatable) | 0..* serviceList | [] |

### ServiceList Element: `<serviceList>` (Repeatable)

| Field | Type | Description | Constraints |
|-------|------|-------------|-------------|
| `<onid>` | integer | Original Network ID | > 0 |
| `<tsid>` | integer | Transport Stream ID | > 0 |
| `<sid>` | integer | Service ID | > 0 |

**Note**: `<serviceList>` can appear 0 or more times within `<searchsetting>`. Each occurrence represents one channel.

### RecSetting Element: `<recsetting>`

| Field | Type | Description | Constraints | Default |
|-------|------|-------------|-------------|---------|
| `<recMode>` | integer | Recording mode (1=standard) | >= 0 | 1 |
| `<priority>` | integer | Recording priority | >= 0 | 2 |
| `<tuijyuuFlag>` | integer | 1=auto-follow enabled | 0 or 1 | 1 |
| `<serviceMode>` | integer | Service mode | >= 0 | 16 |
| `<pittariFlag>` | integer | Exact time match | 0 or 1 | 0 |
| `<batFilePath>` | string | Post-recording batch file | Any string | "" |
| `<recFolderList>` | string | Recording folder list | Any string | "" |
| `<suspendMode>` | integer | System suspend mode | >= 0 | 0 |
| `<defserviceMode>` | integer | Default service mode | >= 0 | 17 |
| `<rebootFlag>` | integer | Reboot after recording | 0 or 1 | 0 |
| `<useMargineFlag>` | integer | Use recording margins | 0 or 1 | 0 |
| `<startMargine>` | integer | Start margin (seconds) | >= 0 | 20 |
| `<endMargine>` | integer | End margin (seconds) | >= 0 | 2 |
| `<continueRecFlag>` | integer | Continue recording flag | 0 or 1 | 0 |
| `<partialRecFlag>` | integer | Partial recording flag | 0 or 1 | 0 |
| `<tunerID>` | integer | Tuner ID (0=auto) | >= 0 | 0 |
| `<partialRecFolder>` | string | Partial recording folder | Any string | "" |

## Character Encoding

**Encoding**: UTF-8

**XML Declaration**: `<?xml version="1.0" encoding="UTF-8" ?>`

**Japanese Characters**: Fully supported in `<andKey>`, `<notKey>`, and other string fields

**Examples**:
- `<andKey>サイエンスZERO</andKey>`
- `<andKey>ブラタモリ</andKey>`
- `<andKey>NHKスペシャル</andKey>`

**Client Requirements**:
- Parse XML as UTF-8
- Preserve UTF-8 encoding in output (JSON, CSV, TSV)
- Terminal must support UTF-8 for proper display

## Performance Characteristics

### Response Time

- **Typical** (100-500 rules): 100-300ms
- **Large** (500-1000 rules): 300-800ms
- **Very large** (1000+ rules): 1-3 seconds

**Note**: Time measured at API server, does not include network latency.

### Response Size

- **Typical** (100-500 rules, 20-60 channels each): 50-200 KB
- **Large** (1000+ rules): 500 KB - 1 MB

**Compression**: Not supported (no Content-Encoding: gzip)

### Concurrency

- **Read-only**: Safe for concurrent requests
- **No rate limiting**: API does not enforce rate limits
- **Best practice**: Single request per operation (no need for polling)

## Compatibility

### EMWUI Version

- **Tested**: xtne6f/EDCB fork (2024-2025)
- **API Version**: No explicit version in endpoint or response
- **Stability**: Schema assumed stable (no breaking changes expected)

### API Changes

**Known Stable Elements**:
- Endpoint path: `/api/EnumAutoAdd`
- Root structure: `<entry>`, `<total>`, `<count>`, `<items>`
- Core fields: `<ID>`, `<andKey>`, `<notKey>`, `<serviceList>`

**Potential Future Changes** (not guaranteed):
- Additional fields in `<searchsetting>` or `<recsetting>` (backward compatible if added)
- New query parameters for filtering (not currently supported)

**Handling Unknown Fields**:
- Client should ignore unknown XML elements (forward compatibility)
- Go's `encoding/xml` automatically ignores unmapped fields

## Error Handling

### Network Errors

**Connection Refused**:
- EMWUI service not running
- Incorrect endpoint URL
- Firewall blocking connection

**Timeout**:
- Network latency too high
- EMWUI service hung or overloaded

**DNS Resolution Failure**:
- Incorrect hostname in EMWUI_ENDPOINT
- DNS server unavailable

### HTTP Errors

**404 Not Found**:
- API endpoint not available (incompatible EMWUI version)
- Typo in endpoint path

**500 Internal Server Error**:
- EMWUI bug or crash
- Database corruption

**503 Service Unavailable**:
- EMWUI starting up
- Maintenance mode

### XML Parsing Errors

**Malformed XML**:
- Syntax error in response
- Truncated response (network interruption)

**Invalid Data**:
- Negative IDs or counts
- Invalid UTF-8 sequences

**Client Handling**:
- Validate XML structure before parsing
- Return descriptive error with context (line number, partial content)
- Log full response for debugging

## Testing

### Test Scenarios

1. **Success Case**: Retrieve 3 rules with Japanese keywords
2. **Empty Case**: Retrieve from system with no rules (total=0, count=0)
3. **Large Case**: Retrieve 230+ rules (verify performance)
4. **Connection Error**: Request with invalid endpoint
5. **Timeout**: Request with very slow network
6. **Malformed XML**: Inject truncated response

### Test Fixtures

Located in `tests/testdata/responses/`:

- `enumautoadd_success.xml`: 3 rules with varied configurations
- `enumautoadd_empty.xml`: Empty items list (total=0, count=0)
- `enumautoadd_large.xml`: 100+ rules for performance testing

### Mock Server

Extend `tests/testdata/mock_server.go` to handle `/api/EnumAutoAdd`:

```go
// In mock server handler
if r.URL.Path == "/api/EnumAutoAdd" {
    w.Header().Set("Content-Type", "text/xml; charset=utf-8")
    if mock.ReturnEmpty {
        fmt.Fprint(w, enumAutoAddEmptyXML)
    } else {
        fmt.Fprint(w, enumAutoAddSuccessXML)
    }
    return
}
```

## Security Considerations

### Network Security

- **Local Network Only**: EMWUI intended for trusted local network
- **No TLS**: HTTP only (no HTTPS support)
- **Exposure Risk**: Do not expose endpoint to public internet

### Authentication

- **None**: No authentication required
- **Authorization**: No role-based access control
- **Trust Model**: All clients on network have full read access

### Data Sensitivity

- **Recording Preferences**: May reveal viewing habits
- **Keywords**: Search terms may be sensitive
- **Channels**: Subscription information exposed

**Recommendations**:
- Restrict network access to trusted devices
- Use VPN if accessing remotely
- Protect exported files with file system permissions

### Input Validation

- **Read-only**: GET request, no user input in request
- **XML Parsing**: Use safe standard library (no XXE vulnerabilities)
- **No Code Execution**: Response data treated as data only (no eval or exec)

## API Usage Example

### cURL Example

```bash
# Basic request
curl -X GET "http://192.168.1.10:5510/api/EnumAutoAdd"

# With timeout
curl -X GET --max-time 10 "http://192.168.1.10:5510/api/EnumAutoAdd"

# Save to file
curl -X GET "http://192.168.1.10:5510/api/EnumAutoAdd" -o response.xml

# Pretty-print XML
curl -X GET "http://192.168.1.10:5510/api/EnumAutoAdd" | xmllint --format -
```

### Go Client Example

```go
import (
    "encoding/xml"
    "io"
    "net/http"
    "time"
)

client := &http.Client{Timeout: 10 * time.Second}
url := "http://192.168.1.10:5510/api/EnumAutoAdd"

resp, err := client.Get(url)
if err != nil {
    return fmt.Errorf("request failed: %w", err)
}
defer resp.Body.Close()

if resp.StatusCode != http.StatusOK {
    return fmt.Errorf("unexpected status: %d", resp.StatusCode)
}

body, err := io.ReadAll(resp.Body)
if err != nil {
    return fmt.Errorf("read failed: %w", err)
}

var response EnumAutoAddResponse
if err := xml.Unmarshal(body, &response); err != nil {
    return fmt.Errorf("parse failed: %w", err)
}

// Use response.Items
```

## Related APIs

### SetAutoAdd (POST)

- **Purpose**: Create or update automatic recording rule
- **Method**: POST /api/SetAutoAdd?id=0
- **Auth**: Requires CSRF token (ctok) from HTML page
- **Note**: Different from EnumAutoAdd (GET, no auth)

### DeleteAutoAdd (POST /api/SetAutoAdd with del=1)

- **Purpose**: Delete an existing automatic recording rule
- **Method**: POST /api/SetAutoAdd?id={rule_id}
- **Auth**: Requires CSRF token (ctok) from HTML page
- **Request Body**: `del=1&ctok={csrf_token}`
- **Example Request**:
  ```bash
  curl 'http://192.168.1.10:5510/api/SetAutoAdd?id=334' \
    -H 'Content-Type: application/x-www-form-urlencoded; charset=UTF-8' \
    --data-raw 'del=1&ctok=26aa1285b9c974a95d31718a1d0f17bb'
  ```
- **Success Response**:
  ```xml
  <?xml version="1.0" encoding="UTF-8" ?><entry><success>EPG自動予約を削除しました</success></entry>
  ```
- **Error Response**:
  ```xml
  <?xml version="1.0" encoding="UTF-8" ?><entry><err>エラーメッセージ</err></entry>
  ```
- **Note**: Uses the same SetAutoAdd endpoint as create/update, but with `del=1` parameter

### Other EMWUI APIs

- GET /api/EnumService: List available channels
- GET /api/EnumReserveInfo: List manual reservations
- GET /api/EnumRecInfo: List recorded programs
- GET /api/EnumEventInfo: Retrieve EPG (Electronic Program Guide) data

**Note**: This feature focuses only on EnumAutoAdd for automatic recording rules.

## EnumEventInfo API (GET /api/EnumEventInfo)

### Endpoint Overview

**Purpose**: Retrieve EPG (Electronic Program Guide) data for a specific channel

**Full URL**: `{EMWUI_ENDPOINT}/api/EnumEventInfo?ONID={ONID}&TSID={TSID}&SID={SID}&basic=0&count=1000`

**HTTP Method**: GET

**Authentication**: None required

**Content-Type**: text/xml; charset=UTF-8

### Request Parameters

| Parameter | Type | Required | Description | Example |
|-----------|------|----------|-------------|---------|
| ONID | integer | Yes | Original Network ID | 32736 |
| TSID | integer | Yes | Transport Stream ID | 32736 |
| SID | integer | Yes | Service ID | 1024 |
| basic | integer | No | Basic mode flag (0=detailed) | 0 |
| count | integer | No | Maximum events to return | 1000 |

### Response Structure

```xml
<?xml version="1.0" encoding="UTF-8" ?>
<entry>
  <total>539</total>
  <index>0</index>
  <count>10</count>
  <items>
    <eventinfo>
      <ONID>32736</ONID>
      <TSID>32736</TSID>
      <SID>1024</SID>
      <eventID>7323</eventID>
      <service_name>ＮＨＫ総合１・東京</service_name>
      <startDate>2025/12/23</startDate>
      <startTime>06:00:00</startTime>
      <startDayOfWeek>2</startDayOfWeek>
      <duration>1800</duration>
      <event_name>ＮＨＫニュース　おはよう日本</event_name>
      <event_text>番組の概要説明</event_text>
      <contentInfo>
        <nibble1>0</nibble1>
        <nibble2>0</nibble2>
        <component_type_name>ニュース／報道 - 定時・総合</component_type_name>
      </contentInfo>
      <freeCAFlag>0</freeCAFlag>
      <event_ext_text>番組の詳細説明</event_ext_text>
    </eventinfo>
    <!-- ... more eventinfo entries ... -->
  </items>
</entry>
```

### Data Model

#### EventInfo Element: `<eventinfo>`

| Field | Type | Description | Constraints |
|-------|------|-------------|-------------|
| `<ONID>` | integer | Original Network ID | > 0 |
| `<TSID>` | integer | Transport Stream ID | > 0 |
| `<SID>` | integer | Service ID | > 0 |
| `<eventID>` | integer | Event identifier | > 0 |
| `<service_name>` | string | Channel/service name (UTF-8) | Any string |
| `<startDate>` | string | Start date (YYYY/MM/DD) | Date format |
| `<startTime>` | string | Start time (HH:MM:SS) | Time format |
| `<startDayOfWeek>` | integer | Day of week (0=Sun, 1=Mon, ...) | 0-6 |
| `<duration>` | integer | Duration in seconds | >= 0 |
| `<event_name>` | string | Program title (UTF-8) | Any string |
| `<event_text>` | string | Program description (UTF-8) | Any string |
| `<event_ext_text>` | string | Extended description (UTF-8) | Any string |
| `<freeCAFlag>` | integer | 0=free, 1=scrambled | 0 or 1 |
| `<contentInfo>` | container | Genre information (repeatable) | 0..* contentInfo |

#### ContentInfo Element: `<contentInfo>` (Repeatable)

| Field | Type | Description | Constraints |
|-------|------|-------------|-------------|
| `<nibble1>` | integer | Genre main category | 0-15 |
| `<nibble2>` | integer | Genre sub category | 0-15 |
| `<component_type_name>` | string | Human-readable genre name (UTF-8) | Any string |

### Character Encoding

**Encoding**: UTF-8

**Japanese Characters**: Fully supported in all string fields

**Examples**:
- `<event_name>ＮＨＫニュース　おはよう日本</event_name>`
- `<service_name>ＮＨＫ総合１・東京</service_name>`

### Performance Characteristics

#### Response Time

- **Typical** (30-100 events): 100-500ms
- **Large** (500-1000 events): 500ms-2s

#### Response Size

- **Typical** (100 events): 100-300 KB
- **Large** (1000 events): 1-3 MB

### Error Responses

#### EPG Data Not Available

**Scenario**: EPG data is still loading or not available for the channel

**Response**:
```xml
<?xml version="1.0" encoding="UTF-8" ?>
<entry><err>EPGデータを読み込み中、または存在しません</err></entry>
```

#### Invalid Channel

**Scenario**: ONID/TSID/SID combination does not exist

**Response**: Empty response (total=0, count=0)

### Usage Example

```bash
# Retrieve EPG for NHK (ONID=32736, TSID=32736, SID=1024)
curl "http://192.168.1.10:5510/api/EnumEventInfo?ONID=32736&TSID=32736&SID=1024&basic=0&count=1000"
```

### CLI Integration

The `epg` command uses this API to retrieve program schedule data:

```bash
# View EPG for a specific channel
epgtimer epg --channel "32736-32736-1024"

# View EPG for all channels
epgtimer epg --all-channels

# Filter by title
epgtimer epg --channel "32736-32736-1024" --title "ニュース"
```
