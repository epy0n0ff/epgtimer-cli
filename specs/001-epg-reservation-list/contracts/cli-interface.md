# CLI Interface Contract: List Command

**Feature**: EPG Reservation List Retrieval (001-epg-reservation-list)
**Command**: `epgtimer list`
**Version**: 1.0.0

## Command Signature

```bash
epgtimer list [flags]
```

## Description

Retrieve and display automatic recording rules from EpgTimer's EMWUI service. Supports filtering by keywords, channels, and status, with output in multiple formats (table, JSON, CSV, TSV).

## Flags

### Filter Flags

| Flag | Type | Default | Description | Example |
|------|------|---------|-------------|---------|
| `--andKey <keyword>` | string | "" | Filter by search keyword (substring, case-insensitive) | `--andKey ニュース` |
| `--channel <id>` | string | "" | Filter by channel (ONID-TSID-SID format) | `--channel 32736-32736-1024` |
| `--enabled` | bool | false | Show only enabled rules | `--enabled` |
| `--disabled` | bool | false | Show only disabled rules | `--disabled` |
| `--regex` | bool | false | Show only regex-enabled rules | `--regex` |

### Output Flags

| Flag | Type | Default | Description | Example |
|------|------|---------|-------------|---------|
| `--format <fmt>` | string | "table" | Output format: table, json, csv, tsv | `--format json` |
| `--output <file>` | string | "-" (stdout) | Output file path (use "-" for stdout) | `--output rules.csv` |

### Connection Flags (Inherited from root)

| Flag | Type | Default | Description | Example |
|------|------|---------|-------------|---------|
| `--endpoint <url>` | string | env:EMWUI_ENDPOINT | EMWUI service endpoint | `--endpoint http://localhost:5510` |

## Output Formats

### Table Format (Default)

Human-readable table for terminal display.

**Columns**: ID | Enabled | Keywords | Exclusions | Channels

**Width Limits**:
- Keywords: Truncated at 30 chars with "..."
- Exclusions: Truncated at 30 chars with "..."
- Channels: Shows count (e.g., "62 channels")

**Example**:
```
ID  Enabled  Keywords             Exclusions  Channels
1   Yes      サイエンスZERO       [再]        62 channels
2   Yes      ブラタモリ           [再]        62 channels
3   Yes      うまいッ！           [再]        62 channels
```

### JSON Format

Machine-readable JSON with full structure preservation.

**Structure**: Array of objects with nested search and recording settings

**Example**:
```json
[
  {
    "id": 1,
    "search": {
      "disabled": 0,
      "case_sensitive": 0,
      "and_key": "サイエンスZERO",
      "not_key": "[再]",
      "regex": 0,
      "title_only": 0,
      "duration_min": 0,
      "duration_max": 0,
      "channels": [
        {"onid": 32736, "tsid": 32736, "sid": 1024},
        {"onid": 32736, "tsid": 32736, "sid": 1025}
      ]
    },
    "recording": {
      "rec_mode": 1,
      "priority": 2,
      "auto_follow": 1
    }
  }
]
```

### CSV Format

Spreadsheet-compatible CSV with headers.

**Headers**: id, enabled, and_key, not_key, regex, title_only, duration_min, duration_max, channel_count, channels, rec_mode, priority, auto_follow

**Escaping**: Follows RFC 4180 (quotes around fields with commas, double-quotes escaped)

**Boolean Representation**: "true"/"false" strings

**Channel Representation**: Comma-separated string (quoted)

**Example**:
```csv
id,enabled,and_key,not_key,regex,title_only,duration_min,duration_max,channel_count,channels,rec_mode,priority,auto_follow
1,true,"サイエンスZERO","[再]",false,false,0,0,62,"32736-32736-1024,32736-32736-1025,32737-32737-1032",1,2,true
2,true,"ブラタモリ","[再]",false,false,0,0,62,"32736-32736-1024,32736-32736-1025,32737-32737-1032",1,2,true
```

### TSV Format

Tab-separated values (similar to CSV but tab-delimited).

**Headers**: Same as CSV

**Delimiter**: Tab character (\t)

**No Escaping**: Fields not quoted (assumes no tabs in data)

**Example**:
```tsv
id	enabled	and_key	not_key	regex	title_only	duration_min	duration_max	channel_count	channels	rec_mode	priority	auto_follow
1	true	サイエンスZERO	[再]	false	false	0	0	62	32736-32736-1024,32736-32736-1025,32737-32737-1032	1	2	true
2	true	ブラタモリ	[再]	false	false	0	0	62	32736-32736-1024,32736-32736-1025,32737-32737-1032	1	2	true
```

## Exit Codes

| Code | Meaning | Description |
|------|---------|-------------|
| 0 | Success | Rules retrieved and displayed successfully |
| 1 | API Error | Failed to connect to EMWUI or parse response |
| 2 | Invalid Arguments | Invalid flag combination or format specified |
| 3 | File Error | Failed to write output file |

## Error Handling

### Connection Errors

**Scenario**: EMWUI service unreachable

**Output** (stderr):
```
Error: Failed to connect to EMWUI service at http://192.168.1.10:5510
Cause: dial tcp 192.168.1.10:5510: i/o timeout

Troubleshooting:
1. Check that EpgTimer is running
2. Verify EMWUI_ENDPOINT environment variable
3. Confirm network connectivity
```

**Exit Code**: 1

### Parse Errors

**Scenario**: Malformed XML response

**Output** (stderr):
```
Error: Failed to parse XML response from EnumAutoAdd API
Cause: XML syntax error at line 15: unexpected EOF

This may indicate an API incompatibility. Response body:
<?xml version="1.0" encoding="UTF-8" ?><entry><total>230</total><index>0</index><count>...
```

**Exit Code**: 1

### Empty Results

**Scenario**: No rules found (or all filtered out)

**Output** (stdout, table format):
```
No automatic recording rules found.
```

**Output** (JSON format):
```json
[]
```

**Output** (CSV/TSV format):
```csv
id,enabled,and_key,not_key,regex,title_only,duration_min,duration_max,channel_count,channels,rec_mode,priority,auto_follow
```

**Exit Code**: 0 (not an error)

### Invalid Flag Combinations

**Scenario**: Both --enabled and --disabled specified

**Output** (stderr):
```
Error: Cannot specify both --enabled and --disabled flags
```

**Exit Code**: 2

**Scenario**: Invalid format specified

**Output** (stderr):
```
Error: Invalid format "yaml". Supported formats: table, json, csv, tsv
```

**Exit Code**: 2

### File Write Errors

**Scenario**: Cannot write to output file

**Output** (stderr):
```
Error: Failed to write output to file: rules.csv
Cause: open rules.csv: permission denied
```

**Exit Code**: 3

## Usage Examples

### Basic Usage

```bash
# List all rules (table format)
epgtimer list

# List all rules as JSON
epgtimer list --format json

# List only enabled rules
epgtimer list --enabled

# Filter by keyword
epgtimer list --andKey ニュース

# Filter by channel
epgtimer list --channel 32736-32736-1024

# Combine filters
epgtimer list --enabled --andKey ドラマ
```

### Export to File

```bash
# Export to CSV file
epgtimer list --format csv --output rules.csv

# Export to JSON file
epgtimer list --format json --output rules.json

# Export filtered results to TSV
epgtimer list --enabled --format tsv --output enabled-rules.tsv
```

### Advanced Usage

```bash
# List only regex rules as JSON
epgtimer list --regex --format json

# List disabled rules and export to CSV
epgtimer list --disabled --format csv --output disabled-rules.csv

# Filter by keyword and export to file
epgtimer list --andKey NHKスペシャル --format json --output nhk-special.json
```

## Performance Characteristics

### Response Time

- **Small datasets** (1-100 rules): < 1 second
- **Medium datasets** (100-500 rules): 1-2 seconds
- **Large datasets** (500-1000 rules): 2-3 seconds
- **Very large datasets** (1000+ rules): 3-10 seconds

**Note**: Time includes network request, XML parsing, filtering, and formatting.

### Memory Usage

- **Typical**: ~5-10 MB (for 100-500 rules)
- **Large**: ~20-30 MB (for 1000+ rules)

### Network Traffic

- **Typical request**: ~50-200 KB (for 100-500 rules with 20-60 channels each)
- **Large request**: ~500 KB - 1 MB (for 1000+ rules)

## Compatibility

### EMWUI API Version

- Tested with: EpgTimer xtne6f/EDCB fork (2024-2025)
- Endpoint: GET /api/EnumAutoAdd
- Response format: XML (UTF-8)

### Go Version

- Minimum: Go 1.24
- Uses standard library only (no external dependencies for parsing/formatting)

### Terminal Compatibility

- UTF-8 terminal required for Japanese character display
- Table format tested on: Linux terminal, Windows Terminal, macOS Terminal
- Width: Assumes minimum 80 character width

## Security Considerations

### Authentication

- No authentication required (EMWUI runs on trusted local network)
- Endpoint should not be exposed to public internet

### Input Validation

- Filter values sanitized before use (no command injection risk)
- XML parsing uses safe Go standard library (no XXE vulnerabilities)
- File paths validated before writing (no path traversal)

### Data Exposure

- Output may contain sensitive recording preferences
- Exported files should be protected with appropriate permissions
- No passwords or credentials stored or transmitted

## Future Enhancements (Out of Scope)

- Server-side filtering (if EMWUI adds query parameters)
- Sorting options (--sort-by id|keyword|priority)
- Column selection for table format (--columns id,keyword,channels)
- Verbose mode (--verbose) to show all fields in table
- Pagination (--page 1 --per-page 50) for very large datasets
- Interactive mode (TUI with arrow key navigation)
