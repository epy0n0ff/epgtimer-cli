# Implementation Tasks: EPG Reservation List Retrieval

**Feature**: EPG Reservation List Retrieval (001-epg-reservation-list)
**Branch**: `001-epg-reservation-list`
**Generated**: 2025-12-21

## Task Summary

**Total Tasks**: 18
**User Stories**: 3 (P1, P2, P3)
**Parallel Opportunities**: 8 tasks marked [P]
**MVP Scope**: User Story 1 (Phase 3) - View All Active Reservations

## Implementation Strategy

**Incremental Delivery**:
1. **MVP** (Phase 1-3): Complete User Story 1 - basic list command with table output
2. **Enhancement 1** (Phase 4): Add User Story 2 - filtering capabilities
3. **Enhancement 2** (Phase 5): Add User Story 3 - export formats (CSV, JSON, TSV)
4. **Polish** (Phase 6): Cross-cutting concerns and documentation

**Key Principle**: Each user story phase is independently testable and can be deployed separately.

---

## Phase 1: Setup & Foundation

**Goal**: Prepare test infrastructure and create data models needed for all user stories.

**Prerequisites**: None (blocking for all user stories)

**Tasks**:

- [x] T001 Create test fixture files for EnumAutoAdd API responses in tests/testdata/responses/
- [x] T002 [P] Create enumautoadd_success.xml with 3 sample rules (Japanese keywords, multiple channels) in tests/testdata/responses/
- [x] T003 [P] Create enumautoadd_empty.xml with total=0, count=0 in tests/testdata/responses/
- [x] T004 Extend mock EMWUI server to handle GET /api/EnumAutoAdd endpoint in tests/testdata/mock_server.go

---

## Phase 2: Data Models (Foundational)

**Goal**: Implement all data structures needed to parse and represent automatic recording rules.

**Prerequisites**: Phase 1 complete

**Blocking For**: All user story implementations

**Tasks**:

- [x] T005 [P] Create EnumAutoAddResponse struct with XML tags in internal/models/reservation.go
- [x] T006 [P] Create AutoAddRule struct with XML and JSON tags in internal/models/reservation.go
- [x] T007 [P] Create SearchSettings struct with all fields and XML/JSON tags in internal/models/reservation.go
- [x] T008 [P] Create RecordingSettings struct with all fields and XML/JSON tags in internal/models/reservation.go
- [x] T009 [P] Create ServiceInfo struct with ONID/TSID/SID fields in internal/models/reservation.go
- [x] T010 [P] Add helper methods (IsEnabled, IsRegex, String) to SearchSettings and ServiceInfo in internal/models/reservation.go

---

## Phase 3: User Story 1 - View All Active Reservations (P1)

**User Story**: A user wants to see all currently scheduled EPG recording reservations from the command line.

**Goal**: Implement basic list command that retrieves and displays all rules in table format.

**Prerequisites**: Phase 1 and Phase 2 complete

**Independent Test**: Run `epgtimer list` and verify all rules are displayed in table format with correct encoding.

**Success Criteria**:
- Command retrieves all rules from EMWUI API
- Japanese characters display correctly
- Table shows ID, Enabled status, Keywords, Exclusions, Channel count
- Empty list handled gracefully
- Connection errors show helpful troubleshooting messages

**Tasks**:

- [x] T011 [US1] Implement GET /api/EnumAutoAdd in EnumAutoAdd() method in internal/client/list.go
- [x] T012 [US1] Add XML response parsing and error handling to EnumAutoAdd() in internal/client/list.go
- [x] T013 [P] [US1] Create table formatter with Format(rules) method in internal/formatters/table.go
- [x] T014 [P] [US1] Implement table column layout (ID, Enabled, Keywords, Exclusions, Channels) in internal/formatters/table.go
- [x] T015 [US1] Create list command with Cobra in internal/commands/list.go
- [x] T016 [US1] Add --endpoint flag (inherited from root) connection logic in internal/commands/list.go
- [x] T017 [US1] Call client.EnumAutoAdd() and format with table formatter in internal/commands/list.go
- [x] T018 [US1] Handle errors (connection, parse, empty list) with user-friendly messages in internal/commands/list.go
- [x] T019 [US1] Register list command with root command in internal/commands/root.go
- [x] T020 [US1] Write integration test for successful rule retrieval in tests/integration/list_test.go
- [x] T021 [US1] Write integration test for empty response handling in tests/integration/list_test.go
- [x] T022 [US1] Write integration test for connection error handling in tests/integration/list_test.go
- [x] T023 [US1] Write integration test for Japanese character encoding in tests/integration/list_test.go

**Acceptance Test Commands**:
```bash
# Test with live EMWUI server
epgtimer list --endpoint http://192.168.1.10:5510

# Expected output: Table with all rules
# ID  Enabled  Keywords             Exclusions  Channels
# 1   Yes      サイエンスZERO       [再]        62 channels
# ...
```

---

## Phase 4: User Story 2 - Filter Reservations by Criteria (P2)

**User Story**: A user wants to filter the reservation list by specific criteria to quickly find relevant recordings.

**Goal**: Add filtering flags to list command (--andKey, --channel, --enabled, --disabled, --regex).

**Prerequisites**: Phase 3 (User Story 1) complete

**Independent Test**: Run `epgtimer list --enabled --andKey ニュース` and verify only matching rules are displayed.

**Success Criteria**:
- --andKey filters by keyword substring (case-insensitive)
- --channel filters by ONID-TSID-SID format
- --enabled/--disabled filter by rule status
- --regex shows only regex-enabled rules
- Multiple filters work together (AND logic)
- Filter with no matches shows "No rules found"

**Tasks**:

- [ ] T024 [P] [US2] Create FilterOptions struct with filter fields in internal/models/filter.go
- [ ] T025 [P] [US2] Implement Matches(rule) method with filter logic in internal/models/filter.go
- [ ] T026 [US2] Add filter flags (--andKey, --channel, --enabled, --disabled, --regex) to list command in internal/commands/list.go
- [ ] T027 [US2] Build FilterOptions from flags and apply to retrieved rules in internal/commands/list.go
- [ ] T028 [US2] Handle empty filter results with "No rules match filter criteria" message in internal/commands/list.go
- [ ] T029 [US2] Write integration test for --andKey filter in tests/integration/list_test.go
- [ ] T030 [US2] Write integration test for --channel filter in tests/integration/list_test.go
- [ ] T031 [US2] Write integration test for --enabled/--disabled filters in tests/integration/list_test.go
- [ ] T032 [US2] Write integration test for combined filters in tests/integration/list_test.go

**Acceptance Test Commands**:
```bash
# Filter by keyword
epgtimer list --andKey ニュース

# Filter by channel
epgtimer list --channel 32736-32736-1024

# Filter by status
epgtimer list --enabled

# Combine filters
epgtimer list --enabled --andKey ドラマ
```

---

## Phase 5: User Story 3 - Export Reservations to File (P3)

**User Story**: A user wants to export the reservation list to a file (CSV, JSON, TSV) for backup or analysis.

**Goal**: Implement formatters for JSON, CSV, TSV and add --format and --output flags.

**Prerequisites**: Phase 3 (User Story 1) complete (independent of Phase 4)

**Independent Test**: Run `epgtimer list --format json --output rules.json` and verify file contains valid JSON with all rule data.

**Success Criteria**:
- JSON format preserves full nested structure
- CSV format has headers and proper escaping
- TSV format uses tab delimiters
- --output writes to specified file path
- File write errors show helpful messages
- UTF-8 encoding preserved in all formats

**Tasks**:

- [ ] T033 [P] [US3] Create JSON formatter with Format(rules) method in internal/formatters/json.go
- [ ] T034 [P] [US3] Implement JSON marshaling with indentation in internal/formatters/json.go
- [ ] T035 [P] [US3] Create CSV formatter with Format(rules) method in internal/formatters/csv.go
- [ ] T036 [P] [US3] Implement CSV header row and row flattening with proper escaping in internal/formatters/csv.go
- [ ] T037 [P] [US3] Create TSV formatter with Format(rules) method in internal/formatters/tsv.go
- [ ] T038 [P] [US3] Implement TSV output (tab-delimited, same structure as CSV) in internal/formatters/tsv.go
- [ ] T039 [US3] Add --format flag (table|json|csv|tsv, default=table) to list command in internal/commands/list.go
- [ ] T040 [US3] Add --output flag (file path, default=stdout) to list command in internal/commands/list.go
- [ ] T041 [US3] Implement format selection logic with formatter factory pattern in internal/commands/list.go
- [ ] T042 [US3] Implement file writing with error handling (permissions, disk space) in internal/commands/list.go
- [ ] T043 [US3] Write integration test for JSON export in tests/integration/list_test.go
- [ ] T044 [US3] Write integration test for CSV export with UTF-8 keywords in tests/integration/list_test.go
- [ ] T045 [US3] Write integration test for TSV export in tests/integration/list_test.go
- [ ] T046 [US3] Write integration test for file write errors in tests/integration/list_test.go

**Acceptance Test Commands**:
```bash
# Export to JSON
epgtimer list --format json --output rules.json

# Export to CSV
epgtimer list --format csv --output rules.csv

# Export to TSV
epgtimer list --format tsv --output rules.tsv

# Export filtered results
epgtimer list --enabled --format json --output enabled-rules.json
```

---

## Phase 6: Polish & Documentation

**Goal**: Final integration, edge case handling, and user documentation.

**Prerequisites**: Phases 3, 4, and 5 complete

**Tasks**:

- [ ] T047 Test with large dataset (230+ rules) and verify performance <3 seconds in tests/integration/list_test.go
- [ ] T048 Verify all error messages follow existing patterns from SetAutoAdd in internal/commands/list.go
- [ ] T049 Update README.md with list command usage examples in README.md
- [ ] T050 Update quickstart.md with real command outputs (if needed) in specs/001-epg-reservation-list/quickstart.md

---

## Dependency Graph

```
Phase 1 (Setup)
    ↓
Phase 2 (Data Models) ← BLOCKING for all user stories
    ↓
    ├─→ Phase 3 (US1: View All) ← MVP
    │       ↓
    │       ├─→ Phase 4 (US2: Filter) ← Depends on US1
    │       └─→ Phase 5 (US3: Export) ← Depends on US1 (independent of US2)
    │
    └─→ Phase 4 and Phase 5 can run in parallel after Phase 3
            ↓
        Phase 6 (Polish) ← Depends on all user stories
```

**Key Dependencies**:
- **Phase 2 blocks all user stories**: Must complete data models first
- **Phase 3 (US1) blocks Phase 4 and Phase 5**: Core list functionality required first
- **Phase 4 (US2) independent of Phase 5 (US3)**: Can implement in either order after US1
- **Phase 6 requires all phases**: Final integration and polish

**Suggested Order**:
1. Phase 1 → Phase 2 (foundational, must complete first)
2. Phase 3 (US1) - MVP delivery
3. Phase 4 (US2) OR Phase 5 (US3) - choose based on priority
4. Remaining user story
5. Phase 6 (polish)

---

## Parallel Execution Opportunities

### Phase 1 (Setup)
```bash
# Tasks T002, T003 can run in parallel (different files)
[T002] Create enumautoadd_success.xml
[T003] Create enumautoadd_empty.xml
```

### Phase 2 (Data Models)
```bash
# All struct definitions can be written in parallel (same file, different structs)
[T005] EnumAutoAddResponse struct
[T006] AutoAddRule struct
[T007] SearchSettings struct
[T008] RecordingSettings struct
[T009] ServiceInfo struct
[T010] Helper methods
```

### Phase 3 (US1)
```bash
# Formatters can be implemented while client code is being written
[T011-T012] API client (internal/client/list.go)
[T013-T014] Table formatter (internal/formatters/table.go) ← Parallel to client
```

### Phase 4 (US2)
```bash
# Filter model and filter logic can be written in parallel
[T024-T025] FilterOptions and Matches() (internal/models/filter.go)
[T026-T028] Apply filters in command (internal/commands/list.go)
```

### Phase 5 (US3)
```bash
# All formatters can be implemented in parallel
[T033-T034] JSON formatter (internal/formatters/json.go)
[T035-T036] CSV formatter (internal/formatters/csv.go)
[T037-T038] TSV formatter (internal/formatters/tsv.go)
```

---

## Testing Strategy

### Unit Tests
- Data model parsing (XML → structs)
- Filter matching logic (FilterOptions.Matches)
- Formatter output (table, JSON, CSV, TSV)

### Integration Tests
**User Story 1**:
- ✅ Successful retrieval and table display
- ✅ Empty response handling
- ✅ Connection error handling
- ✅ Japanese character encoding

**User Story 2**:
- ✅ Each filter type (--andKey, --channel, --enabled, --disabled, --regex)
- ✅ Combined filters (AND logic)
- ✅ Empty filter results

**User Story 3**:
- ✅ Each export format (JSON, CSV, TSV)
- ✅ File writing and error handling
- ✅ UTF-8 preservation in exports

### Manual Testing
```bash
# Test against live EMWUI server
export EMWUI_ENDPOINT="http://192.168.1.10:5510"

# Basic functionality
epgtimer list

# Filtering
epgtimer list --enabled --andKey ニュース

# Export
epgtimer list --format json --output backup.json
```

---

## Implementation Notes

### Existing Code Patterns to Follow

**Client Methods** (from internal/client/add.go):
- HTTP client with 10-second timeout
- POST method uses `client.Post(endpoint, formData)`
- GET method should use `client.Get(endpoint)` (needs implementation)
- Error wrapping: `fmt.Errorf("failed to X: %w", err)`

**Command Structure** (from internal/commands/add.go):
- Cobra command with `Use`, `Short`, `Long`
- Flags defined in `init()` with `addCmd.Flags()`
- RunE function for execution with error return
- Client initialization from root command context

**Error Messages** (from existing code):
- Connection errors: "Failed to connect to EMWUI service at {url}"
- Parse errors: "Failed to parse XML response"
- Include troubleshooting steps in error output

### New Patterns Introduced

**Formatters Package**:
```go
// Common interface for all formatters
type Formatter interface {
    Format(rules []AutoAddRule) (string, error)
}

// Factory function
func GetFormatter(format string) (Formatter, error)
```

**Filter Pattern**:
```go
// Functional style filtering
filtered := []AutoAddRule{}
for _, rule := range rules {
    if filterOpts.Matches(&rule) {
        filtered = append(filtered, rule)
    }
}
```

---

## Performance Targets

- **API call**: <200ms p95 (network dependent)
- **Parsing 100 rules**: <100ms
- **Parsing 1000 rules**: <1 second
- **Table formatting**: <50ms for any size
- **JSON export**: <200ms for 1000 rules
- **CSV export**: <300ms for 1000 rules

---

## File Checklist

**New Files** (to be created):
- ✅ internal/models/reservation.go (Phase 2)
- ✅ internal/models/filter.go (Phase 4)
- ✅ internal/client/list.go (Phase 3)
- ✅ internal/commands/list.go (Phase 3)
- ✅ internal/formatters/table.go (Phase 3)
- ✅ internal/formatters/json.go (Phase 5)
- ✅ internal/formatters/csv.go (Phase 5)
- ✅ internal/formatters/tsv.go (Phase 5)
- ✅ tests/integration/list_test.go (Phases 3-5)
- ✅ tests/testdata/responses/enumautoadd_success.xml (Phase 1)
- ✅ tests/testdata/responses/enumautoadd_empty.xml (Phase 1)

**Modified Files**:
- ✅ internal/commands/root.go (register list command)
- ✅ tests/testdata/mock_server.go (add EnumAutoAdd handler)
- ✅ README.md (usage examples)

**Total**: 11 new files, 3 modified files

---

## Questions for Implementation

1. **Table truncation**: Max width for Keywords/Exclusions columns? (Suggested: 30 chars)
2. **CSV channel format**: Comma-separated within quoted cell? (Suggested: yes, "ONID-TSID-SID,ONID-TSID-SID")
3. **File overwrite**: Prompt user or auto-overwrite? (Suggested: auto-overwrite, no prompt)
4. **Error on invalid format**: Should --format xyz error or default to table? (Suggested: error with helpful message)

---

## Success Criteria Summary

**Phase 3 (MVP - US1)**:
- ✅ Can retrieve and display all rules
- ✅ Japanese characters display correctly
- ✅ Empty list handled gracefully
- ✅ Connection errors show helpful messages
- ✅ Integration tests pass

**Phase 4 (US2)**:
- ✅ All filter flags work correctly
- ✅ Combined filters use AND logic
- ✅ Empty filter results show appropriate message
- ✅ Integration tests pass

**Phase 5 (US3)**:
- ✅ JSON, CSV, TSV formats export correctly
- ✅ Files written successfully
- ✅ UTF-8 encoding preserved
- ✅ Integration tests pass

**Phase 6 (Polish)**:
- ✅ All error messages consistent
- ✅ Documentation complete
- ✅ Performance targets met
- ✅ Manual testing against live EMWUI passes

---

**End of Tasks Document**

**Next Command**: `/speckit.implement` to begin implementation, or implement user stories incrementally starting with Phase 1-3 (MVP).
