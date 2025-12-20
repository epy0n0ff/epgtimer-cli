# Tasks: EpgTimer EMWUI CLI Interface

**Input**: Design documents from `/specs/001-epgtimer-cli/`
**Prerequisites**: plan.md (required), spec.md (required for user stories), research.md, data-model.md, contracts/

**Tests**: Tests are NOT explicitly requested in the specification. Test tasks are included as optional for validation.

**Organization**: Tasks are grouped by user story to enable independent implementation and testing.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: Which user story this task belongs to (e.g., US1)
- Include exact file paths in descriptions

## Path Conventions

This is a single CLI application with structure:
- `cmd/epgtimer/` - Entry point
- `internal/` - Internal packages
- `tests/` - Test files

---

## Phase 1: Setup (Shared Infrastructure)

**Purpose**: Project initialization and basic structure

- [x] T001 Initialize Go module with `go mod init github.com/epy0n0ff/epgtimer-cli`
- [x] T002 Add Cobra dependency with `go get github.com/spf13/cobra@latest`
- [x] T003 [P] Create directory structure: cmd/epgtimer/, internal/models/, internal/client/, internal/commands/, internal/validator/, tests/integration/, tests/testdata/
- [x] T004 [P] Create go.sum and verify dependencies
- [x] T005 [P] Create .gitignore for Go project (bin/, *.exe, go.work, etc.)
- [x] T006 [P] Create Makefile with build, test, install, and clean targets

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Core infrastructure that MUST be complete before user story implementation

**⚠️ CRITICAL**: No user story work can begin until this phase is complete

- [ ] T007 Create AutoAddRuleRequest model struct in internal/models/autoadd.go with andKey, notKey, serviceList, and all default parameters from data-model.md
- [ ] T008 [P] Create ServiceListEntry model struct in internal/models/channel.go with ONID-TSID-SID parsing logic
- [ ] T009 [P] Create AutoAddRuleResponse model struct in internal/models/response.go
- [ ] T010 Implement NewAutoAddRuleRequest() constructor with defaults from curl sample in internal/models/autoadd.go
- [ ] T011 Implement Validate() method for AutoAddRuleRequest in internal/models/autoadd.go
- [ ] T012 Implement ToFormData() method for AutoAddRuleRequest to generate application/x-www-form-urlencoded format in internal/models/autoadd.go
- [ ] T013 [P] Create HTTP client wrapper in internal/client/client.go with timeout configuration (10 seconds)
- [ ] T014 [P] Implement ParseServiceListEntry() function in internal/models/channel.go to validate ONID-TSID-SID format
- [ ] T015 Create root command structure with Cobra in internal/commands/root.go including version flag and EMWUI_ENDPOINT config
- [ ] T016 Create main.go entry point in cmd/epgtimer/main.go that initializes Cobra and calls root command

**Checkpoint**: Foundation ready - user story implementation can now begin

---

## Phase 3: User Story 1 - Add Automatic Recording Rule (Priority: P1) 🎯 MVP

**Goal**: Enable users to create keyword-based automatic recording rules via CLI

**Independent Test**: Run `epgtimer add --andKey "ニュース" --serviceList "32736-32736-1024"` and verify rule is created in EMWUI

### Implementation for User Story 1

- [ ] T017 [US1] Implement SetAutoAdd API method in internal/client/add.go that POSTs form data to /api/SetAutoAdd?id=0
- [ ] T018 [US1] Implement response parsing in internal/client/add.go to detect success/error from EMWUI response
- [ ] T019 [US1] Create add command structure in internal/commands/add.go with Cobra
- [ ] T020 [US1] Add --andKey flag (required, string) to add command in internal/commands/add.go
- [ ] T021 [P] [US1] Add --notKey flag (optional, string) to add command in internal/commands/add.go
- [ ] T022 [P] [US1] Add --serviceList flag (required, comma-separated string list) to add command in internal/commands/add.go
- [ ] T023 [US1] Implement command execution logic in internal/commands/add.go: parse flags → validate → create request → call API → display result
- [ ] T024 [US1] Implement EMWUI_ENDPOINT environment variable reading with validation in internal/commands/add.go
- [ ] T025 [US1] Add success message output "Automatic recording rule created successfully" in internal/commands/add.go
- [ ] T026 [US1] Add error handling for missing EMWUI_ENDPOINT with actionable message in internal/commands/add.go
- [ ] T027 [US1] Add error handling for connection failures with actionable message in internal/commands/add.go
- [ ] T028 [US1] Add error handling for validation errors with field-specific messages in internal/commands/add.go
- [ ] T029 [US1] Add error handling for API errors (duplicate, invalid params) in internal/commands/add.go
- [ ] T030 [US1] Register add command with root command in internal/commands/root.go

**Checkpoint**: User Story 1 complete - CLI can add automatic recording rules

---

## Phase 4: Integration Testing & Validation (Optional)

**Purpose**: Validate complete CLI functionality with mock EMWUI server

- [x] T031 [P] Create mock EMWUI server in tests/testdata/mock_server.go that responds to SetAutoAdd requests
- [x] T032 [P] Create integration test in tests/integration/add_test.go for successful rule creation
- [x] T033 [P] Create integration test in tests/integration/add_test.go for missing andKey error
- [x] T034 [P] Create integration test in tests/integration/add_test.go for missing serviceList error
- [x] T035 [P] Create integration test in tests/integration/add_test.go for invalid serviceList format
- [x] T036 [P] Create integration test in tests/integration/add_test.go for connection error
- [x] T037 [P] Create integration test in tests/integration/add_test.go for Japanese keyword encoding
- [x] T038 [P] Create test fixtures in tests/testdata/responses/ with success and error responses

---

## Phase 5: Polish & Cross-Cutting Concerns

**Purpose**: Documentation, build tooling, and final polish

- [x] T039 [P] Create README.md with installation instructions, usage examples, and EMWUI_ENDPOINT configuration
- [x] T040 [P] Add usage examples with Japanese keywords to README.md
- [x] T041 [P] Add common channel IDs table to README.md (NHK総合, 日本テレビ, etc.)
- [x] T042 [P] Update Makefile build target to create cross-platform binaries (Linux, macOS, Windows)
- [x] T043 [P] Add Makefile install target that copies binary to /usr/local/bin or equivalent
- [x] T044 [P] Test CLI with actual EMWUI service if available
- [x] T045 Verify quickstart.md scenarios work correctly
- [x] T046 [P] Add --help text with detailed parameter descriptions and examples
- [x] T047 [P] Add --version flag output with version information
- [x] T048 Final validation: run all integration tests, test binary installation, verify Japanese text handling

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: No dependencies - can start immediately
- **Foundational (Phase 2)**: Depends on Setup completion - BLOCKS all user stories
- **User Story 1 (Phase 3)**: Depends on Foundational phase completion
- **Integration Testing (Phase 4)**: Depends on User Story 1 completion (optional)
- **Polish (Phase 5)**: Can overlap with Phase 4, depends on User Story 1 core functionality

### Task Dependencies Within Phases

**Phase 2 (Foundational)**:
- T010 depends on T007 (need struct before constructor)
- T011 depends on T007 (need struct before Validate)
- T012 depends on T007 (need struct before ToFormData)
- T014 depends on T008 (need ServiceListEntry struct)
- T016 depends on T015 (need root command)

**Phase 3 (User Story 1)**:
- T017-T018 can be done in parallel (client implementation)
- T019-T022 can be done in parallel (command setup and flags)
- T023 depends on T017-T022 (need API and command structure)
- T024-T029 can be done in parallel after T023 (error handling)
- T030 depends on T019 (need add command before registering)

**Phase 4 (Integration Testing)**:
- T031 is independent (mock server)
- T032-T037 all depend on T031 (need mock server)
- T038 is independent (test fixtures)

**Phase 5 (Polish)**:
- All tasks marked [P] can run in parallel
- T045 depends on Phase 3 completion
- T048 depends on all other Phase 5 tasks

### Parallel Opportunities

**Phase 1**: T003, T004, T005, T006 can all run in parallel

**Phase 2**: T008, T009, T013, T014 can run in parallel

**Phase 3**:
- T017-T018 together (client)
- T019-T022 together (command flags)
- T024-T029 together (error handling)

**Phase 4**: T031-T038 can all run in parallel (independent test files)

**Phase 5**: T039-T044, T046-T047 can all run in parallel

---

## Parallel Example: User Story 1

```bash
# Parallel group 1 (client implementation):
Task: "Implement SetAutoAdd API method in internal/client/add.go"
Task: "Implement response parsing in internal/client/add.go"

# Parallel group 2 (command flags):
Task: "Create add command structure in internal/commands/add.go"
Task: "Add --andKey flag to add command"
Task: "Add --notKey flag to add command"
Task: "Add --serviceList flag to add command"

# Parallel group 3 (error handling):
Task: "Add error handling for missing EMWUI_ENDPOINT"
Task: "Add error handling for connection failures"
Task: "Add error handling for validation errors"
Task: "Add error handling for API errors"
```

---

## Implementation Strategy

### MVP First (User Story 1 Only)

1. Complete Phase 1: Setup (T001-T006)
2. Complete Phase 2: Foundational (T007-T016) - CRITICAL
3. Complete Phase 3: User Story 1 (T017-T030)
4. **STOP and VALIDATE**:
   - Set EMWUI_ENDPOINT environment variable
   - Run `epgtimer add --andKey "test" --serviceList "32736-32736-1024"`
   - Verify rule appears in EMWUI interface
5. Optionally add integration tests (Phase 4)
6. Add documentation and polish (Phase 5)

### Incremental Delivery

1. **Foundation** (Phase 1-2): ~16 tasks → Project compiles, models validate
2. **MVP** (Phase 3): ~14 tasks → Full CLI functionality working
3. **Quality** (Phase 4): ~8 tasks → Comprehensive test coverage (optional)
4. **Polish** (Phase 5): ~10 tasks → Production-ready with docs

### Single Developer Timeline

- Phase 1: 1-2 hours (setup)
- Phase 2: 3-4 hours (models and validation)
- Phase 3: 4-5 hours (CLI implementation)
- Phase 4: 2-3 hours (testing - optional)
- Phase 5: 2-3 hours (documentation and polish)

**Total estimate**: 12-17 hours for complete implementation

---

## Notes

- [P] tasks = different files, no dependencies, can run in parallel
- [US1] label = belongs to User Story 1 (only one user story in this feature)
- User Story 1 should be independently testable after Phase 3
- Integration tests (Phase 4) are optional but recommended
- Commit after each task or logical group
- Test Japanese keyword encoding early (T037)
- Verify EMWUI_ENDPOINT handling before API calls
- Use curl sample defaults exactly as specified in data-model.md
- Stop at any checkpoint to validate independently

---

## Task Summary

**Total Tasks**: 48
- Phase 1 (Setup): 6 tasks
- Phase 2 (Foundational): 10 tasks
- Phase 3 (User Story 1): 14 tasks
- Phase 4 (Integration Testing): 8 tasks
- Phase 5 (Polish): 10 tasks

**Parallel Opportunities**: 24 tasks marked [P]

**User Story Breakdown**:
- User Story 1 (P1): 14 implementation tasks + supporting infrastructure

**Independent Test Criteria**:
- User Story 1: Execute `epgtimer add` command with valid parameters and verify rule creation in EMWUI interface

**Suggested MVP Scope**: Phases 1-3 (30 tasks) delivers fully functional CLI for adding automatic recording rules
