# Implementation Plan: EpgTimer EMWUI CLI Interface

**Branch**: `001-epgtimer-cli` | **Date**: 2025-12-20 | **Spec**: [spec.md](spec.md)
**Input**: Feature specification from `/specs/001-epgtimer-cli/spec.md`

**Note**: This template is filled in by the `/speckit.plan` command. See `.specify/templates/commands/plan.md` for the execution workflow.

## Summary

Implement a command-line interface for EpgTimer's EMWUI service that enables users to create automatic recording rules based on keyword search criteria. The CLI will be written in Go 1.24, connect to the EMWUI SetAutoAdd API endpoint without authentication, and provide clear success/error messages with comprehensive error handling. This enables automation of recurring program recording based on title keywords and channel selection.

**API Endpoint**: `POST /api/SetAutoAdd` with form data
**Key Parameters**: andKey (search keywords), notKey (exclusion keywords), serviceList (channel list)
**Default Values**: All other parameters use defaults from the provided curl sample

## Technical Context

**Language/Version**: Go 1.24
**Primary Dependencies**:
- Standard library (net/http, flag for CLI, time for date/time handling)
- github.com/spf13/cobra (CLI framework with subcommands and help generation)
- No table formatting needed (only add command with success/error messages)

**Storage**: N/A (stateless CLI client)
**Testing**: Go standard testing framework (testing package), table-driven tests
**Target Platform**: Cross-platform (Linux, macOS, Windows) - compiled Go binary
**Project Type**: Single CLI application
**Performance Goals**:
- API response handling < 3 seconds (list up to 100 items)
- Command execution latency < 100ms (excluding network)
- Binary size < 10MB (single statically compiled executable)

**Constraints**:
- No authentication required (trusted network)
- Must handle Japanese character encoding properly (UTF-8)
- Terminal output must adapt to various terminal widths
- Network timeout handling (5-10 second default)
- Configuration via environment variable (EMWUI_ENDPOINT)

**Scale/Scope**:
- Single binary CLI tool
- 1 primary command (add)
- Simple command-line interface with parameter validation
- ~5-8 source files, ~500-800 LOC estimated

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

**Status**: ⚠️ No project constitution found

The constitution file (`.specify/memory/constitution.md`) contains only a template. This should be filled in with project-specific principles before proceeding to ensure consistent architecture and development practices.

**Recommended Constitution Principles for CLI Tools**:
1. **Single Binary Principle**: Ship as single statically-compiled binary with no runtime dependencies
2. **12-Factor Config**: Configuration via environment variables, never hardcode endpoints
3. **UNIX Philosophy**: Do one thing well, compose with other tools, clear exit codes
4. **Error Handling**: All errors must be actionable with clear messages indicating what failed and how to fix
5. **Test Coverage**: Unit tests for business logic, integration tests against mock EMWUI API

**Gates Assumed for This Feature**:
- ✅ Stateless design (no persistent storage)
- ✅ Standard library first (minimize dependencies)
- ✅ Cross-platform compatibility
- ✅ Clear separation: API client, business logic, CLI presentation

## Project Structure

### Documentation (this feature)

```text
specs/001-epgtimer-cli/
├── plan.md              # This file (/speckit.plan command output)
├── research.md          # Phase 0 output (/speckit.plan command)
├── data-model.md        # Phase 1 output (/speckit.plan command)
├── quickstart.md        # Phase 1 output (/speckit.plan command)
├── contracts/           # Phase 1 output (/speckit.plan command)
│   └── emwui-api.yaml   # OpenAPI spec or equivalent for EMWUI endpoints
└── tasks.md             # Phase 2 output (/speckit.tasks command - NOT created by /speckit.plan)
```

### Source Code (repository root)

```text
# Single CLI application structure
cmd/
└── epgtimer/
    └── main.go          # Entry point, CLI framework setup

internal/
├── client/              # EMWUI HTTP client
│   ├── client.go        # HTTP client wrapper
│   ├── add.go           # Add reservation API method
│   └── types.go         # Request/response types
├── models/              # Domain models
│   ├── reservation.go   # Reservation entity (for add requests)
│   └── channel.go       # Channel entity
├── commands/            # CLI command handlers
│   ├── add.go           # Add reservation command
│   └── root.go          # Root command and global flags
└── validator/           # Input validation
    └── validate.go      # Parameter validation logic

tests/
├── integration/         # Tests against mock EMWUI server
│   └── add_test.go      # Integration tests for add command
└── testdata/            # Mock responses and fixtures
    └── responses/

go.mod
go.sum
Makefile                 # Build, test, install targets
README.md
```

**Structure Decision**: Single project structure is appropriate for a CLI tool. The `internal/` package prevents external imports while keeping code organized. The `cmd/` directory follows Go project layout conventions. Tests are separated by type (unit tests alongside source files, integration tests in dedicated directory).

## Complexity Tracking

> **Fill ONLY if Constitution Check has violations that must be justified**

N/A - No constitution violations identified (constitution template not yet filled).

## Phase 0: Research & Unknowns

### Research Tasks

The following unknowns from Technical Context have been investigated:

1. **EMWUI API Protocol & Endpoints** ✅
   - Finding: Custom HTTP-based API with Lua CGI endpoints
   - Format: Plain text responses (NOT REST/JSON)
   - Add endpoint: `/api/AddReserveManual.lua` with query parameters
   - No list endpoint available (confirmed by user)
   - Errors communicated via HTTP status + "NG" in response body

2. **CLI Framework Selection** ✅
   - Decision: Cobra (github.com/spf13/cobra)
   - Rationale: Industry standard, excellent help generation, good for future expansion
   - Binary size acceptable (~2-3 MB) for features provided

3. **Japanese Character Encoding** ✅
   - Go handles UTF-8 natively
   - URL encoding required for Japanese titles in query parameters
   - Use `net/url.QueryEscape()` for parameter encoding

4. **EpgTimer EMWUI Documentation** ✅
   - Primary source: xtne6f/EDCB GitHub repository
   - Lua API scripts provide implementation details
   - No formal API versioning - compatibility detection needed

**Output**: All findings consolidated in `research.md`

## Phase 1: Design Artifacts

*Prerequisites: research.md complete*

### Data Model (data-model.md)

Define Go structs for:
- **AutoAddRuleRequest**: andKey, notKey, serviceList, and all default parameters from curl sample
- **AutoAddRuleResponse**: Success status, error message if failed
- **ServiceListEntry**: ONID-TSID-SID format string (e.g., "32736-32736-1024")
- Validation rules for each field (required, format, constraints)
- Form data encoding logic for POST request

### API Contracts (contracts/)

Document EMWUI SetAutoAdd endpoint:
- Endpoint: `POST /api/SetAutoAdd?id=0`
- Content-Type: `application/x-www-form-urlencoded`
- Form parameters: andKey, notKey, serviceList (multiple), and all defaults from curl sample
- Response format: Plain text or JSON (verify with curl sample)
- HTTP status codes and error detection patterns
- Example request/response pairs based on provided curl sample

### Quickstart (quickstart.md)

User-facing guide:
1. Installation (go install or binary download)
2. Configuration (setting EMWUI_ENDPOINT environment variable)
3. Add command usage: andKey, notKey (optional), serviceList
4. Examples with Japanese keywords ("ニュース", "ドラマ", etc.)
5. Channel list format and examples
6. Troubleshooting common errors (connection, invalid parameters, encoding issues)

## Phase 2: Implementation Tasks

*Generated by `/speckit.tasks` command - not part of this plan output*

The tasks.md file will break down implementation into:
- Project initialization (go.mod, directory structure)
- EMWUI client implementation (HTTP POST requests to SetAutoAdd endpoint)
- Domain models (AutoAddRuleRequest, validation, form data encoding)
- CLI add command (parameter parsing for andKey/notKey/serviceList, validation, API calls)
- Default parameter management (loading defaults from curl sample)
- Error handling (connection errors, validation errors, API errors)
- Japanese text encoding (form data encoding for UTF-8 keywords)
- ServiceList format handling (ONID-TSID-SID triplets)
- Integration tests (mock EMWUI server)
- Documentation (README, usage examples with keywords and channels)
- Build and release setup (Makefile, cross-platform builds)
