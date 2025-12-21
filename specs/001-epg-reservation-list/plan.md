# Implementation Plan: EPG Reservation List Retrieval

**Branch**: `001-epg-reservation-list` | **Date**: 2025-12-21 | **Spec**: [spec.md](./spec.md)
**Input**: Feature specification from `/specs/001-epg-reservation-list/spec.md`

**Note**: This template is filled in by the `/speckit.plan` command. See `.specify/templates/commands/plan.md` for the execution workflow.

## Summary

Implement a CLI command to retrieve and display the list of automatic recording rules from EpgTimer's EMWUI service using the GET /api/EnumAutoAdd endpoint. The feature will support filtering by criteria (channel, date, status) and exporting to multiple formats (CSV, JSON, TSV) as requested by the user. This extends the existing epgtimer-cli tool which already supports adding automatic recording rules via SetAutoAdd.

## Technical Context

**Language/Version**: Go 1.24.3
**Primary Dependencies**: github.com/spf13/cobra (v1.10.2), encoding/xml, encoding/json, encoding/csv
**Storage**: N/A (reads from remote EMWUI API)
**Testing**: Go testing package, integration tests with mock HTTP server
**Target Platform**: Cross-platform CLI (Linux, Windows via WSL, macOS)
**Project Type**: Single CLI project
**Performance Goals**: <3 seconds response time for 1-100 reservations, <10 seconds for 1000+ reservations
**Constraints**: <200ms p95 for API calls (network dependent), handle Japanese UTF-8 encoding correctly
**Scale/Scope**: Support 1000+ reservation rules, filter operations must be efficient on client side

**Known Technical Details:**
- EMWUI API returns XML responses (format needs research via EnumAutoAdd endpoint)
- Existing client supports dynamic CSRF token fetching (not needed for GET requests)
- Existing models use encoding/xml for response parsing
- CLI follows Cobra command pattern with flags for user input
- Need to research: EnumAutoAdd response structure, available fields, filtering approach (server-side vs client-side)

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

**Initial Check (Pre-Research)**: ✅ PASS
**Post-Design Check**: ✅ PASS

**Status**: ✅ PASS (No constitution file exists yet - project is still establishing patterns)

**Notes:**
- This is the second feature in the project (first: SetAutoAdd)
- Following established patterns from existing code
- No architectural violations detected
- Extending existing client library with new list operation
- Added new `formatters` package for output formatting (table, JSON, CSV, TSV)
- No external dependencies added (uses Go standard library only)
- Maintains single CLI project structure
- Client-side filtering approach consistent with read-only API design

## Project Structure

### Documentation (this feature)

```text
specs/[###-feature]/
├── plan.md              # This file (/speckit.plan command output)
├── research.md          # Phase 0 output (/speckit.plan command)
├── data-model.md        # Phase 1 output (/speckit.plan command)
├── quickstart.md        # Phase 1 output (/speckit.plan command)
├── contracts/           # Phase 1 output (/speckit.plan command)
└── tasks.md             # Phase 2 output (/speckit.tasks command - NOT created by /speckit.plan)
```

### Source Code (repository root)

```text
cmd/
└── epgtimer/
    └── main.go           # CLI entry point

internal/
├── client/
│   ├── client.go         # HTTP client wrapper with ctok fetching
│   ├── add.go            # SetAutoAdd API (existing)
│   └── list.go           # EnumAutoAdd API (NEW for this feature)
├── commands/
│   ├── root.go           # Root command setup
│   ├── add.go            # Add command (existing)
│   └── list.go           # List command (NEW for this feature)
├── models/
│   ├── autoadd.go        # SetAutoAdd request/response (existing)
│   ├── reservation.go    # Reservation model (NEW for this feature)
│   └── channel.go        # Channel parsing utilities (existing)
└── formatters/           # (NEW package for this feature)
    ├── table.go          # Human-readable table output
    ├── csv.go            # CSV export
    ├── json.go           # JSON export
    └── tsv.go            # TSV export

tests/
├── integration/
│   ├── add_test.go       # SetAutoAdd integration tests (existing)
│   └── list_test.go      # EnumAutoAdd integration tests (NEW)
└── testdata/
    ├── mock_server.go    # Mock EMWUI server (existing, extend)
    └── responses/
        ├── enumautoadd_success.xml  # (NEW sample response)
        └── enumautoadd_empty.xml    # (NEW empty list response)
```

**Structure Decision**: Following the existing single CLI project structure. This feature extends the internal packages with new list command, API client method, response models, and a new formatters package for output formats (CSV, JSON, TSV). The existing test infrastructure (mock server, integration tests) will be extended to cover the new endpoint.

## Complexity Tracking

> **Fill ONLY if Constitution Check has violations that must be justified**

| Violation | Why Needed | Simpler Alternative Rejected Because |
|-----------|------------|-------------------------------------|
| [e.g., 4th project] | [current need] | [why 3 projects insufficient] |
| [e.g., Repository pattern] | [specific problem] | [why direct DB access insufficient] |
