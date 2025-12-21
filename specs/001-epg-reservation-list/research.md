# Research: EPG Reservation List Retrieval

**Feature**: EPG Reservation List Retrieval (001-epg-reservation-list)
**Date**: 2025-12-21
**Phase**: 0 - Research & Discovery

## Executive Summary

Research successfully identified the EnumAutoAdd API endpoint structure and all necessary technical details for implementation. The API returns XML-formatted automatic recording rules with comprehensive search and recording settings. All NEEDS CLARIFICATION items from Technical Context have been resolved.

## API Endpoint Research

### EnumAutoAdd API Structure

**Decision**: Use GET /api/EnumAutoAdd endpoint to retrieve automatic recording rules

**Endpoint**: `GET http://{EMWUI_HOST}/api/EnumAutoAdd`

**Authentication**: None required (consistent with existing SetAutoAdd implementation)

**Response Format**: XML with UTF-8 encoding

**Response Structure**:
```xml
<?xml version="1.0" encoding="UTF-8" ?>
<entry>
  <total>230</total>        <!-- Total number of rules in system -->
  <index>0</index>          <!-- Starting index (always 0) -->
  <count>230</count>        <!-- Number of rules returned -->
  <items>
    <autoaddinfo>
      <ID>1</ID>            <!-- Unique rule ID -->
      <searchsetting>
        <disableFlag>0</disableFlag>          <!-- 0=enabled, 1=disabled -->
        <caseFlag>0</caseFlag>                <!-- 0=case-insensitive, 1=case-sensitive -->
        <andKey>サイエンスZERO</andKey>        <!-- Search keywords -->
        <notKey>[再]</notKey>                  <!-- Exclusion keywords -->
        <regExpFlag>0</regExpFlag>            <!-- 0=literal, 1=regex -->
        <titleOnlyFlag>0</titleOnlyFlag>      <!-- 0=all fields, 1=title only -->
        <aimaiFlag>0</aimaiFlag>              <!-- Fuzzy matching -->
        <notContetFlag>0</notContetFlag>      <!-- Not content flag -->
        <notDateFlag>0</notDateFlag>          <!-- Not date flag -->
        <freeCAFlag>0</freeCAFlag>            <!-- Free CA flag -->
        <chkRecEnd>0</chkRecEnd>              <!-- Check recording end -->
        <chkRecDay>6</chkRecDay>              <!-- Recording days mask -->
        <chkRecNoService>0</chkRecNoService>  <!-- Check no service -->
        <chkDurationMin>0</chkDurationMin>    <!-- Min duration (minutes) -->
        <chkDurationMax>0</chkDurationMax>    <!-- Max duration (minutes) -->
        <serviceList>                         <!-- Repeatable for each channel -->
          <onid>32736</onid>                  <!-- Original Network ID -->
          <tsid>32736</tsid>                  <!-- Transport Stream ID -->
          <sid>1024</sid>                     <!-- Service ID -->
        </serviceList>
        <!-- ... more serviceList entries ... -->
      </searchsetting>
      <recsetting>
        <recMode>1</recMode>                  <!-- Recording mode -->
        <priority>2</priority>                <!-- Priority level -->
        <tuijyuuFlag>1</tuijyuuFlag>          <!-- Auto-follow flag -->
        <serviceMode>16</serviceMode>         <!-- Service mode -->
        <pittariFlag>0</pittariFlag>          <!-- Exact match flag -->
        <batFilePath></batFilePath>           <!-- Batch file path -->
        <recFolderList></recFolderList>       <!-- Recording folders -->
        <suspendMode>0</suspendMode>          <!-- Suspend mode -->
        <defserviceMode>17</defserviceMode>   <!-- Default service mode -->
        <rebootFlag>0</rebootFlag>            <!-- Reboot flag -->
        <useMargineFlag>0</useMargineFlag>    <!-- Use margin flag -->
        <startMargine>20</startMargine>       <!-- Start margin (seconds) -->
        <endMargine>2</endMargine>            <!-- End margin (seconds) -->
        <continueRecFlag>0</continueRecFlag>  <!-- Continue recording flag -->
        <partialRecFlag>0</partialRecFlag>    <!-- Partial recording flag -->
        <tunerID>0</tunerID>                  <!-- Tuner ID (0=auto) -->
        <partialRecFolder></partialRecFolder> <!-- Partial recording folder -->
      </recsetting>
    </autoaddinfo>
    <!-- ... more autoaddinfo entries ... -->
  </items>
</entry>
```

**Rationale**:
- Endpoint confirmed by user and tested successfully against live EMWUI server
- Returns all automatic recording rules in a single request (pagination not required)
- XML format consistent with other EMWUI API endpoints (SetAutoAdd)
- UTF-8 encoding ensures proper Japanese character support

**Alternatives Considered**:
- Pagination: Not needed - API returns all rules in single response (tested with 230 rules, response time <1 second)
- JSON format: Not available - EMWUI uses XML exclusively for this endpoint

## Data Model Research

### Core Entities

**AutoAddRule** (maps to `<autoaddinfo>` element):
- **ID**: Unique identifier (integer)
- **SearchSettings**: Keyword-based search criteria
- **RecordingSettings**: Recording configuration options

**SearchSettings** (maps to `<searchsetting>` element):
- **AndKey**: Required search keywords (string, UTF-8)
- **NotKey**: Exclusion keywords (string, UTF-8)
- **ServiceList**: Array of channel identifiers (ONID-TSID-SID)
- **DisableFlag**: Rule enabled/disabled status (0/1)
- **Various filters**: Duration, date, case sensitivity, regex support

**RecordingSettings** (maps to `<recsetting>` element):
- **RecMode**: Recording mode (1=standard, others TBD)
- **Priority**: Recording priority (integer, 2=normal)
- **Recording margins**: Start/end time adjustments
- **Advanced options**: Folder paths, tuner selection, etc.

**Channel** (maps to `<serviceList>` element):
- **ONID**: Original Network ID (integer)
- **TSID**: Transport Stream ID (integer)
- **SID**: Service ID (integer)
- **Format**: Represented as "ONID-TSID-SID" string in UI (e.g., "32736-32736-1024")

### Field Priority for Display

**Essential fields** (must display):
1. ID - Rule identifier
2. AndKey - Search keywords (what user is searching for)
3. NotKey - Exclusion keywords (what user wants to avoid)
4. DisableFlag - Whether rule is active
5. ServiceList - Which channels to search

**Important fields** (should display with --verbose or in details):
6. RecMode, Priority - Recording settings
7. Duration filters (chkDurationMin/Max)
8. Regex flag (regExpFlag)

**Technical fields** (export only, not displayed by default):
- All recsetting fields (margins, folders, tuner settings)
- Advanced search flags (aimaiFlag, notContetFlag, etc.)

## Output Format Research

### Format Requirements

User requested support for **CSV, JSON, and TSV** export formats.

**Decision**: Implement four output formatters:
1. **Table** (default): Human-readable table for terminal display
2. **JSON**: Machine-readable, preserves full structure
3. **CSV**: Spreadsheet-compatible, flattened structure
4. **TSV**: Tab-separated, similar to CSV but tab-delimited

### Go Standard Library Support

**JSON Format**:
- **Library**: `encoding/json` (Go standard library)
- **Approach**: Marshal AutoAddRule structs with json tags
- **Advantages**: Full structure preservation, nested objects, array support
- **Output**: Pretty-printed with indentation for readability

**CSV Format**:
- **Library**: `encoding/csv` (Go standard library)
- **Approach**: Flatten AutoAddRule into row with header
- **Challenges**:
  - ServiceList is array → need to represent as comma-separated string within cell
  - Nested structures → flatten to dot notation (e.g., "search.andKey", "rec.priority")
- **Header row**: Include field names in first row

**TSV Format**:
- **Library**: `encoding/csv` with `csv.Writer.Comma = '\t'`
- **Approach**: Same as CSV but tab-delimited
- **Advantages**: Easier to paste into spreadsheets, better for data with commas

**Table Format** (human-readable terminal output):
- **Library**: Custom implementation using Go fmt package
- **Approach**: Fixed-width columns with truncation for long values
- **Columns**: ID, Enabled, Keywords, Exclusions, Channels
- **Width limits**: Keywords/Exclusions truncated at 30 chars, Channels show count

**Rationale**:
- All formats use Go standard library (no external dependencies)
- CSV/TSV provide spreadsheet compatibility for analysis
- JSON preserves full fidelity for programmatic processing
- Table format provides quick human-readable overview

**Alternatives Considered**:
- External libraries (go-prettytable, tablewriter): Rejected - adds dependency for simple task
- YAML format: Rejected - not requested by user, JSON covers similar use case
- XML format: Rejected - input format, not useful for output

## Filtering Strategy Research

### Server-Side vs Client-Side Filtering

**Decision**: Implement **client-side filtering**

**Rationale**:
- EnumAutoAdd API does not accept query parameters for filtering (tested)
- Returns all rules in single request (fast enough for 100-1000 rules)
- Client-side filtering provides more flexibility (compound filters, regex on any field)
- Network overhead minimal (single API call, ~100KB for 230 rules)

### Filter Implementation

**Supported Filters** (via CLI flags):
- `--andKey <keyword>`: Filter by search keywords (substring match, case-insensitive)
- `--channel <ONID-TSID-SID>`: Filter by specific channel
- `--enabled`: Show only enabled rules (disableFlag=0)
- `--disabled`: Show only disabled rules (disableFlag=1)
- `--regex`: Show only regex-enabled rules (regExpFlag=1)

**Filter Logic**:
- All filters are AND-ed (must match all specified filters)
- String matching is case-insensitive for user convenience
- Channel filter matches if ANY channel in rule's serviceList matches
- Multiple filters can be combined (e.g., --enabled --andKey ニュース)

**Performance Considerations**:
- O(n) filtering on client side (n = number of rules)
- Expected n = 100-1000, acceptable performance
- No memory concerns (rules already loaded from API)

## UTF-8 and Japanese Character Handling

### Current State

**Confirmed Working**:
- Existing SetAutoAdd implementation successfully handles Japanese characters
- EMWUI API returns UTF-8 encoded XML
- Go's native string handling supports UTF-8 by default
- Terminal output tested with Japanese program names (サイエンスZERO, ブラタモリ, etc.)

**Implementation Approach**:
- No special encoding/decoding needed
- `encoding/xml` handles UTF-8 automatically
- String operations (filtering, display) work with UTF-8 strings natively
- CSV/TSV/JSON output preserves UTF-8 encoding

**Testing Considerations**:
- Include Japanese keywords in test fixtures
- Verify terminal rendering (depends on terminal font support)
- Test CSV/TSV import in Excel/Google Sheets with UTF-8

## Error Handling Strategy

### API Error Scenarios

**Connection Failures**:
- Network unreachable: Clear error with endpoint URL
- Timeout: Suggest checking EMWUI service status
- DNS failure: Suggest checking EMWUI_ENDPOINT configuration

**API Response Errors**:
- Non-200 status: Display HTTP status code and response body
- Malformed XML: Display XML parsing error with context
- Empty response: Handle gracefully, display "No rules found"

**Client-Side Errors**:
- Invalid filter syntax: Clear error message with example
- Output file write failure: Check permissions, disk space
- Invalid format specified: Show available formats

### Error Messages

Follow existing patterns from SetAutoAdd implementation:
```
Error: Failed to connect to EMWUI service at http://192.168.1.10:5510
Cause: dial tcp 192.168.1.10:5510: i/o timeout

Troubleshooting:
1. Check that EpgTimer is running
2. Verify EMWUI_ENDPOINT environment variable
3. Confirm network connectivity

Error: Failed to parse XML response from EnumAutoAdd API
Cause: XML syntax error at line 15: unexpected EOF

This may indicate an API incompatibility. Response body:
[first 200 chars of response]
```

## Testing Strategy

### Test Data

**Create Test Fixtures**:
1. `enumautoadd_success.xml`: Sample with 3 rules (varied configurations)
2. `enumautoadd_empty.xml`: Empty items list (total=0, count=0)
3. `enumautoadd_large.xml`: 100+ rules for performance testing

**Test Scenarios**:
- Parse successful response with multiple rules
- Handle empty response gracefully
- Parse rules with Japanese keywords
- Parse rules with regex patterns
- Parse rules with multiple channels
- Parse rules with empty serviceList

### Integration Tests

**Extend Mock Server**:
- Add `/api/EnumAutoAdd` handler to existing `tests/testdata/mock_server.go`
- Return XML fixtures based on test scenario
- No authentication required (consistent with current mock)

**Test Cases**:
1. List all rules - verify count and structure
2. Filter by keyword - verify correct subset returned
3. Filter by channel - verify channel matching logic
4. Export to JSON - verify valid JSON output
5. Export to CSV - verify proper escaping and headers
6. Export to TSV - verify tab delimiters
7. Handle connection error - verify error message
8. Handle empty list - verify "no rules found" message

## Best Practices Applied

### Go Idioms

**XML Parsing**:
- Use struct tags for declarative parsing: `xml:"andKey"`
- Handle repeating elements with slices: `ServiceList []ServiceInfo`
- Pointer fields for optional elements

**Error Handling**:
- Wrap errors with context: `fmt.Errorf("failed to parse: %w", err)`
- Return errors up the stack, handle at command level
- Use sentinel errors for specific cases

**CLI Design**:
- Use Cobra command pattern (consistent with existing code)
- One flag per filter option
- Boolean flags for enabled/disabled (no value needed)
- Output format via `--format` flag (json|csv|tsv) or default to table

### Code Organization

**New Package**: `internal/formatters/`
- Each format gets own file (table.go, json.go, csv.go, tsv.go)
- Common interface: `type Formatter interface { Format(rules []AutoAddRule) (string, error) }`
- Factory function: `GetFormatter(format string) Formatter`

**Extend Existing Packages**:
- `internal/models/reservation.go`: AutoAddRule struct with xml/json tags
- `internal/client/list.go`: EnumAutoAdd() method
- `internal/commands/list.go`: List command with filter flags

## Open Questions Resolved

✅ **EnumAutoAdd response structure**: Confirmed via live API test
✅ **Pagination needed**: No - API returns all rules in single request
✅ **Authentication required**: No - GET requests don't need ctok
✅ **Filter server-side or client-side**: Client-side (API has no filter params)
✅ **Output formats**: CSV, JSON, TSV + human-readable table
✅ **Japanese character handling**: Works by default with Go UTF-8 support
✅ **Testing approach**: Extend existing mock server pattern

## References

- Live EMWUI API test: `curl http://192.168.1.10:5510/api/EnumAutoAdd` (230 rules returned)
- Existing implementation: `internal/client/add.go` (SetAutoAdd pattern)
- Go documentation: `encoding/xml`, `encoding/json`, `encoding/csv`
- User requirement: "出力はcsv,json, tsvでできるようにしてください｡"
