# Feature Specification: EPG Reservation List Retrieval

**Feature Branch**: `001-epg-reservation-list`
**Created**: 2025-12-21
**Status**: Draft
**Input**: User description: "EPG予約一覧取得機能を実装"

## User Scenarios & Testing *(mandatory)*

### User Story 1 - View All Active Reservations (Priority: P1)

A user wants to see all currently scheduled EPG recording reservations from the command line. The system will retrieve and display all active reservations including program titles, channels, scheduled times, and recording status.

**Why this priority**: This is the core read functionality that enables users to monitor their recording schedule, verify that desired programs are queued for recording, and identify any scheduling conflicts or gaps. Essential for users to have visibility into what content will be recorded.

**Independent Test**: Can be fully tested by executing a list command that retrieves all reservations from EMWUI and displays them in a readable format. Success means the command completes without errors and shows all reservations with key details (title, channel, time, status).

**Acceptance Scenarios**:

1. **Given** the user has active recording reservations in EpgTimer, **When** the user runs the list command, **Then** all reservations are displayed with program title, channel name/ID, scheduled date/time, duration, and recording status
2. **Given** the user has no active reservations, **When** the user runs the list command, **Then** a message is displayed indicating "No active reservations found"
3. **Given** the EMWUI service is unavailable, **When** the user runs the list command, **Then** a clear error message is displayed with connection details and troubleshooting guidance
4. **Given** the user has 100+ reservations, **When** the user runs the list command, **Then** all reservations are displayed in a paginated or scrollable format
5. **Given** reservations contain Japanese program titles, **When** the user runs the list command, **Then** Japanese characters are displayed correctly without encoding errors
6. **Given** reservations have various statuses (pending, recording, completed, failed), **When** the user runs the list command, **Then** each reservation shows its current status clearly

---

### User Story 2 - Filter Reservations by Criteria (Priority: P2)

A user wants to filter the reservation list by specific criteria such as channel, date range, or status to quickly find relevant recordings without viewing the entire list.

**Why this priority**: As users accumulate many reservations, filtering becomes important for quick navigation and management. This enhances usability but is not essential for basic visibility of reservations.

**Independent Test**: Can be tested by executing list commands with filter flags (e.g., --channel, --date, --status) and verifying that only matching reservations are displayed. The core list functionality (User Story 1) must work first.

**Acceptance Scenarios**:

1. **Given** the user has reservations across multiple channels, **When** the user runs the list command with --channel filter, **Then** only reservations for the specified channel are displayed
2. **Given** the user has reservations spanning multiple dates, **When** the user runs the list command with --date or --from/--to filters, **Then** only reservations within the specified date range are displayed
3. **Given** the user has reservations with different statuses, **When** the user runs the list command with --status filter, **Then** only reservations matching the specified status are displayed
4. **Given** the user applies multiple filters simultaneously, **When** the user runs the list command, **Then** only reservations matching all filter criteria are displayed

---

### User Story 3 - Export Reservations to File (Priority: P3)

A user wants to export the reservation list to a file (CSV, JSON, or text) for backup, analysis, or sharing with other tools.

**Why this priority**: This is a convenience feature for advanced users who want to process reservation data externally or keep records. Not essential for core functionality.

**Independent Test**: Can be tested by executing list commands with --output or --format flags and verifying that a properly formatted file is created with all reservation data. Depends on User Story 1 working.

**Acceptance Scenarios**:

1. **Given** the user has active reservations, **When** the user runs the list command with --output flag specifying a file path, **Then** a file is created containing all reservations in the specified format
2. **Given** the user requests CSV format, **When** the user runs the list command with --format csv, **Then** reservations are exported with proper CSV structure and headers
3. **Given** the user requests JSON format, **When** the user runs the list command with --format json, **Then** reservations are exported as a valid JSON array with all fields
4. **Given** the output file already exists, **When** the user runs the list command with --output, **Then** the user is prompted to confirm overwrite or the file is automatically overwritten with a warning

---

### Edge Cases

- What happens when the EMWUI API endpoint is not configured or is invalid?
- How does the system handle network timeouts or temporary connection failures during retrieval?
- How are very large reservation lists (1000+ entries) handled for performance and display?
- What happens if the EMWUI API returns malformed or incomplete reservation data?
- How does the CLI handle reservations with missing or null fields (e.g., no program title)?
- What happens when Japanese program titles contain special characters or emoji?
- How are date/time values displayed to users (JST timezone, user-friendly format)?
- What happens if the EMWUI API response format changes or is inconsistent?
- How does the CLI handle concurrent modifications to reservations while listing?

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: CLI MUST connect to the EpgTimer EMWUI service endpoint to retrieve recording reservations
- **FR-002**: CLI MUST provide a list command that retrieves and displays all active EPG recording reservations
- **FR-003**: CLI MUST display key information for each reservation: program title, channel identifier, scheduled date/time, duration, and recording status
- **FR-004**: CLI MUST properly decode and display Japanese characters (UTF-8) in program titles and channel names
- **FR-005**: CLI MUST handle empty result sets gracefully with a user-friendly message indicating no reservations are found
- **FR-006**: CLI MUST handle connection errors gracefully and display user-friendly error messages indicating the nature of the failure
- **FR-007**: CLI MUST support configuration of the EMWUI service endpoint via environment variable (EMWUI_ENDPOINT), consistent with existing CLI functionality
- **FR-008**: CLI MUST connect to the EMWUI service without authentication (service runs on trusted network), consistent with existing CLI behavior
- **FR-009**: CLI MUST handle large result sets (100+ reservations) without performance degradation or display issues
- **FR-010**: CLI MUST display date/time values in Japan Standard Time (JST, UTC+9) in a user-friendly format (e.g., "2025-12-21 19:00")
- **FR-011**: CLI MUST provide filtering options via command-line flags (--channel, --date/--from/--to, --status) for User Story 2
- **FR-012**: CLI MUST support exporting reservation list to file with --output flag and --format option (csv, json, text) for User Story 3
- **FR-013**: CLI MUST use the EMWUI API endpoint GET /api/EnumAutoAdd to retrieve the list of automatic recording rules

### Key Entities

- **Reservation**: Represents a scheduled EPG recording, including program title, channel, scheduled date/time, duration, and status
- **Channel**: Represents a broadcast channel, identified by ONID-TSID-SID format or channel name
- **RecordingStatus**: Represents the current state of a reservation (pending, recording, completed, failed, cancelled)
- **DateTimeRange**: Represents a time period for filtering reservations (start date/time, end date/time)

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: Users can successfully retrieve and view their reservation list on the first attempt within 3 seconds for typical result sets (1-100 reservations)
- **SC-002**: CLI displays all reservation details correctly with proper Japanese character encoding in 100% of test cases
- **SC-003**: Users can identify specific reservations among 50+ entries using filter options within 10 seconds
- **SC-004**: CLI handles EMWUI connection failures with clear, actionable error messages in 100% of failure scenarios
- **SC-005**: Exported files (CSV, JSON) are valid and can be opened/parsed by standard tools in 100% of export operations
- **SC-006**: Users can understand the list command functionality without consulting documentation after seeing help output once
- **SC-007**: CLI handles edge cases (empty lists, network errors, malformed data) gracefully without crashing in 100% of test cases

## Assumptions *(include if applicable)*

- The EMWUI API provides an endpoint for retrieving recording reservation lists (need to identify exact endpoint and response format)
- Reservations retrieved include essential fields: program title, channel identifier, scheduled start time, duration, and status
- The EMWUI API response format is consistent with other EMWUI endpoints (likely XML based on existing SetAutoAdd implementation)
- Users primarily want to see active/upcoming reservations rather than historical completed recordings
- Date/time values in EMWUI API are in Japan Standard Time (JST, UTC+9)
- The existing EMWUI_ENDPOINT configuration will be reused for this feature
- No authentication is required to read reservation data (consistent with existing CLI behavior)
- Reservation statuses follow EDCB/EMWUI standard status codes

## Dependencies *(include if applicable)*

### External Dependencies

- EpgTimer EMWUI service must be running and accessible via configured endpoint
- EMWUI API must provide a reservation list/query endpoint (exact endpoint needs verification)

### Internal Dependencies

- Reuses existing EMWUI client configuration and connection handling from epgtimer-cli (001-epgtimer-cli)
- Reuses existing error handling patterns and user messaging from epgtimer-cli
- Reuses existing UTF-8 encoding/decoding logic for Japanese text

## Out of Scope *(include if applicable)*

- Modifying or deleting reservations (separate feature)
- Real-time updates or live monitoring of reservation status changes
- Detailed program metadata beyond basic EPG information (description, genre, ratings)
- Integration with external calendar systems
- Notification or alerting when recordings start/complete
- Sorting options for displayed reservations (beyond filtering)
- Graphical or terminal UI (TUI) - output is text-based only
