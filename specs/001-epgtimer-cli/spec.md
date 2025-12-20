# Feature Specification: EpgTimer EMWUI CLI Interface

**Feature Branch**: `001-epgtimer-cli`
**Created**: 2025-12-20
**Status**: Draft
**Input**: User description: "EpgTimerのEMWUIのインターフェイスに対応したEPG予約の新規追加ができるCLIを作りたい｡"
**Implementation**: SetAutoAdd API - keyword-based automatic recording rules

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Add Automatic Recording Rule (Priority: P1)

A user wants to create an automatic recording rule from the command line by specifying search keywords and target channels. The system will automatically record all programs matching the criteria without manual intervention for each broadcast.

**Why this priority**: This is the core automation functionality that enables users to set up recurring recording schedules based on program content rather than specific times. Essential for capturing series, regular news programs, or any content matching specific criteria across multiple channels.

**Independent Test**: Can be fully tested by executing a command with search keywords (andKey), exclusion keywords (notKey), and channel list (serviceList), then verifying the rule is created successfully with a clear confirmation message.

**Acceptance Scenarios**:

1. **Given** valid search parameters (andKey, serviceList), **When** the user runs the add command with these parameters, **Then** a new automatic recording rule is created in EpgTimer and a success confirmation is displayed
2. **Given** the user provides both andKey and notKey, **When** the user runs the add command, **Then** the rule is created to include programs matching andKey but exclude those matching notKey
3. **Given** the user provides incomplete parameters (missing andKey or serviceList), **When** the user runs the add command, **Then** a clear validation error message is displayed indicating which parameters are missing
4. **Given** an identical rule already exists, **When** the user runs the add command, **Then** the system rejects the request with a clear error message indicating a duplicate rule exists
5. **Given** the EMWUI service is unavailable, **When** the user runs the add command, **Then** a clear error message is displayed with connection details and troubleshooting guidance
6. **Given** the user provides Japanese keywords with special characters, **When** the user runs the add command, **Then** the keywords are properly encoded and the rule is created successfully

---

### Edge Cases

- What happens when the EpgTimer EMWUI service endpoint URL is not configured or is invalid?
- How does the system handle network timeouts or temporary connection failures?
- How does the CLI handle character encoding for Japanese keywords (UTF-8)?
- What happens when the EMWUI API returns unexpected response formats or errors?
- How are very long keyword strings (>100 characters) handled?
- What happens if the channel ID format in serviceList is invalid?
- How are special characters in keywords (spaces, punctuation, symbols) encoded?
- What happens when the user provides an empty andKey?
- How does the system handle multiple keywords in andKey/notKey (space-separated? comma-separated?)?

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: CLI MUST connect to the EpgTimer EMWUI service endpoint to create automatic recording rules
- **FR-002**: CLI MUST provide an add command that creates automatic recording rules by accepting andKey (search keywords), notKey (exclusion keywords), and serviceList (target channels)
- **FR-003**: CLI MUST validate that andKey and serviceList are provided and non-empty before sending requests
- **FR-004**: CLI MUST properly encode Japanese characters (UTF-8) and special characters in keywords for URL form data
- **FR-005**: CLI MUST handle connection errors gracefully and display user-friendly error messages indicating the nature of the failure
- **FR-006**: CLI MUST support configuration of the EMWUI service endpoint via environment variable (EMWUI_ENDPOINT)
- **FR-007**: CLI MUST display appropriate success or failure messages after attempting to add a rule
- **FR-008**: CLI MUST connect to the EMWUI service without authentication (service runs on trusted network)
- **FR-009**: CLI MUST use default values from the provided curl sample for all parameters except andKey, notKey, and serviceList
- **FR-010**: CLI MUST format serviceList in the format "ONID-TSID-SID" as shown in the curl sample (e.g., "32736-32736-1024")

### Key Entities

- **AutoAddRule**: Represents an automatic recording rule, including search criteria (andKey, notKey), target channels (serviceList), and recording settings
- **SearchKeywords**: Keywords for matching program titles (andKey for inclusion, notKey for exclusion)
- **ChannelService**: Represents a channel in serviceList format (ONID-TSID-SID triplet)

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: Users can successfully create a new automatic recording rule with valid parameters on the first attempt with clear confirmation within 2 seconds
- **SC-002**: 95% of invalid commands or parameters result in clear, actionable error messages that indicate what went wrong and how to fix it
- **SC-003**: CLI operates reliably with the target EMWUI service with 99% success rate for valid add operations under normal network conditions
- **SC-004**: Users can add a rule without consulting documentation after seeing help output once
- **SC-005**: CLI handles Japanese keywords without character encoding errors in 100% of test cases
- **SC-006**: CLI properly formats serviceList parameters matching EMWUI API expectations in 100% of test cases
