# EpgTimer CLI

Command-line interface for EpgTimer's EMWUI to manage automatic recording rules via keyword-based searches.

## Features

- Add automatic recording rules based on search keywords
- Delete automatic recording rules by ID
- List and filter existing recording rules
- View available channels with filtering by type and network
- List manual reservations with filtering
- Browse recorded programs with filtering
- View EPG (Electronic Program Guide) for channels
- Export to JSON, CSV, or TSV format
- Support for Japanese keywords and channel names
- Exclusion keywords to filter out unwanted programs
- Multiple channel selection
- Simple configuration via environment variable

## Installation

### Option 1: Build from Source

```bash
# Clone the repository
git clone https://github.com/epy0n0ff/epgtimer-cli.git
cd epgtimer-cli

# Build the binary
make build

# Install to /usr/local/bin
make install
```

### Option 2: Go Install

```bash
go install github.com/epy0n0ff/epgtimer-cli/cmd/epgtimer@latest
```

### Option 3: Download Pre-built Binary

Download the pre-built binary for your platform from the [releases page](https://github.com/epy0n0ff/epgtimer-cli/releases).

```bash
# Linux/macOS
chmod +x epgtimer
sudo mv epgtimer /usr/local/bin/

# Windows
# Add epgtimer.exe to your PATH
```

## Configuration

Set the EMWUI service endpoint:

```bash
# Linux/macOS
export EMWUI_ENDPOINT="http://localhost:5510"

# Windows (PowerShell)
$env:EMWUI_ENDPOINT = "http://localhost:5510"

# Windows (CMD)
set EMWUI_ENDPOINT=http://localhost:5510
```

**Note**: Replace `localhost:5510` with your EpgTimer EMWUI server address.

## Usage

### Commands

#### Add Recording Rule

Add a new automatic recording rule:

```bash
epgtimer add --andKey "search keywords" --serviceList "ONID-TSID-SID"
```

**Options**:
- `--andKey` (required): Search keywords - programs must contain these keywords in the title
- `--notKey` (optional): Exclusion keywords - programs must NOT contain these keywords
- `--serviceList` (required): Comma-separated list of channels in "ONID-TSID-SID" format
- `--endpoint` (optional): Override EMWUI_ENDPOINT environment variable

**Examples**:

```bash
# Record all news programs on NHK
epgtimer add --andKey "ニュース" --serviceList "32736-32736-1024"

# Record dramas, excluding reruns
epgtimer add \
  --andKey "ドラマ" \
  --notKey "再放送" \
  --serviceList "32736-32736-1024,32737-32737-1032,32738-32738-1040"

# Record specific show with exclusions
epgtimer add \
  --andKey "わたしが恋人になれるわけ" \
  --notKey "推しエンタ" \
  --serviceList "32736-32736-1024,32736-32736-1025"
```

#### Delete Recording Rule

Delete an existing automatic recording rule by its ID:

```bash
epgtimer delete [rule-id]
```

**Arguments**:
- `rule-id` (required): The ID of the recording rule to delete

**Options**:
- `--id`: Alternative way to specify the rule ID
- `--endpoint` (optional): Override EMWUI_ENDPOINT environment variable

**Examples**:

```bash
# First, list rules to find the ID
epgtimer list

# Delete rule with ID 334 (using positional argument)
epgtimer delete 334

# Delete rule using --id flag
epgtimer delete --id 334

# Delete rule with custom endpoint
epgtimer delete --endpoint http://192.168.1.10:5510 334
```

**Note**: Use `epgtimer list` to find the rule IDs you want to delete.

#### List Recording Rules

View and filter existing automatic recording rules:

```bash
epgtimer list [flags]
```

**Filter Options**:
- `--andKey`: Filter by search keyword (substring match, case-insensitive)
- `--channel`: Filter by channel (ONID-TSID-SID format)
- `--enabled`: Show only enabled rules
- `--disabled`: Show only disabled rules
- `--regex`: Show only regex-enabled rules

**Export Options**:
- `--format`: Output format - table (default), json, csv, tsv
- `-o, --output`: Output file path (default: stdout)

**Examples**:

```bash
# List all rules (table format)
epgtimer list

# Filter by keyword
epgtimer list --andKey "ニュース"

# Show only enabled rules
epgtimer list --enabled

# Combine filters
epgtimer list --enabled --andKey "ドラマ"

# Export all rules to JSON file
epgtimer list --format json --output rules.json

# Export filtered rules to CSV
epgtimer list --enabled --format csv -o enabled_rules.csv

# Output as TSV to stdout
epgtimer list --format tsv
```

#### List Channels

View and filter available channels/services configured in EpgTimer:

```bash
epgtimer channels [flags]
```

**Filter Options**:
- `--tv`: Show only TV channels (service_type=1)
- `--radio`: Show only radio channels (service_type=2)
- `--data`: Show only data channels (service_type=192)
- `--network`: Filter by network name (substring match, case-insensitive)
- `--name`: Filter by channel name (substring match, case-insensitive)

**Export Options**:
- `--format`: Output format - table (default), json, csv, tsv
- `-o, --output`: Output file path (default: stdout)

**Examples**:

```bash
# List all channels
epgtimer channels

# Show only TV channels
epgtimer channels --tv

# Show only radio channels
epgtimer channels --radio

# Filter by network
epgtimer channels --network "BS Digital"

# Filter by channel name
epgtimer channels --name "NHK"

# Export to JSON file
epgtimer channels --format json --output channels.json

# Export TV channels to CSV
epgtimer channels --tv --format csv -o tv_channels.csv
```

#### List Reservations

View and filter manual recording reservations:

```bash
epgtimer reservations [flags]
```

**Filter Options**:
- `--title`: Filter by program title (substring match, case-insensitive)
- `--station`: Filter by station name (substring match, case-insensitive)
- `--channel`: Filter by channel ID (exact match, format: ONID-TSID-SID)

**Export Options**:
- `--format`: Output format - table (default), json, csv, tsv
- `-o, --output`: Output file path (default: stdout)

**Examples**:

```bash
# List all reservations
epgtimer reservations

# Filter by title
epgtimer reservations --title "ニュース"

# Filter by station
epgtimer reservations --station "NHK"

# Filter by specific channel
epgtimer reservations --channel "32736-32736-1024"

# Export to JSON file
epgtimer reservations --format json --output reservations.json

# Export to CSV
epgtimer reservations --format csv -o reservations.csv
```

#### List Recordings

View and filter recorded programs:

```bash
epgtimer recordings [flags]
```

**Filter Options**:
- `--title`: Filter by program title (substring match, case-insensitive)
- `--station`: Filter by station name (substring match, case-insensitive)
- `--channel`: Filter by channel ID (exact match, format: ONID-TSID-SID)
- `--protected`: Show only protected recordings

**Export Options**:
- `--format`: Output format - table (default), json, csv, tsv
- `-o, --output`: Output file path (default: stdout)

**Note**: The API returns recordings in paginated batches (200 items per request). This command retrieves only the first batch by default.

**Examples**:

```bash
# List recordings
epgtimer recordings

# Filter by title
epgtimer recordings --title "ニュース"

# Filter by station
epgtimer recordings --station "NHK"

# Show only protected recordings
epgtimer recordings --protected

# Export to JSON file
epgtimer recordings --format json --output recordings.json

# Export to CSV
epgtimer recordings --format csv -o recordings.csv
```

#### View EPG (Program Guide)

View EPG (Electronic Program Guide) data for channels:

```bash
epgtimer epg [flags]
```

**Channel Selection**:
- `--channel`: Specific channel in ONID-TSID-SID format
- `--all-channels`: Retrieve EPG for all channels from serviceList_without_local.txt

**Filter Options**:
- `--title`: Filter by program title (substring match, case-insensitive)
- `--genre`: Filter by genre (substring match, case-insensitive)

**Export Options**:
- `--format`: Output format - table (default), json, csv, tsv
- `-o, --output`: Output file path (default: stdout)

**Examples**:

```bash
# View EPG for a specific channel
epgtimer epg --channel "32736-32736-1024"

# View EPG for all channels
epgtimer epg --all-channels

# Filter by title
epgtimer epg --channel "32736-32736-1024" --title "ニュース"

# Filter by genre
epgtimer epg --channel "32736-32736-1024" --genre "ドラマ"

# Export to JSON file
epgtimer epg --channel "32736-32736-1024" --format json --output epg.json

# Export all channels to CSV
epgtimer epg --all-channels --format csv -o epg.csv
```

## Common Channel IDs (Tokyo Area)

| Channel | ONID-TSID-SID |
|---------|---------------|
| NHK総合 | 32736-32736-1024 |
| NHK教育 (Eテレ) | 32736-32736-1025 |
| 日本テレビ | 32737-32737-1032 |
| TBS | 32738-32738-1040 |
| フジテレビ | 32739-32739-1048 |
| テレビ朝日 | 32740-32740-1056 |
| テレビ東京 | 32741-32741-1064 |

**Note**: Channel IDs vary by region. Use the values from your EMWUI interface.

To find your channel IDs:
1. Open EMWUI in your browser
2. Navigate to automatic recording settings
3. Check the channel list - the format will be shown in the HTML

## Help

View all available options:

```bash
epgtimer --help
epgtimer add --help
epgtimer delete --help
epgtimer list --help
epgtimer channels --help
epgtimer reservations --help
epgtimer recordings --help
epgtimer epg --help
epgtimer --version
```

## Troubleshooting

### Connection Failed

**Error**: "Failed to connect to EMWUI service"

**Solutions**:
1. Check EMWUI_ENDPOINT is set correctly:
   ```bash
   echo $EMWUI_ENDPOINT
   ```
2. Verify EpgTimer service is running
3. Test connection:
   ```bash
   curl http://localhost:5510/api/
   ```

### Invalid Parameters

**Error**: "andKey is required"

**Solution**: Ensure you provide the `--andKey` parameter with a search keyword

**Error**: "serviceList is required"

**Solution**: Provide at least one channel in `--serviceList`

**Error**: "invalid channel format"

**Solution**: Ensure channels are in "ONID-TSID-SID" format (e.g., "32736-32736-1024")

### Character Encoding Issues

**Problem**: Japanese characters not working

**Solution**: Ensure your terminal supports UTF-8 encoding. The CLI automatically handles URL encoding for Japanese text.

## Development

### Build

```bash
# Build for current platform
make build

# Build for all platforms (Linux, macOS, Windows)
make build-all

# Run tests
make test

# Run tests with coverage
make test-coverage
```

### Project Structure

```
epgtimer-cli/
├── cmd/epgtimer/          # Main entry point
├── internal/
│   ├── models/            # Data models (rules, filters)
│   ├── client/            # HTTP client for EMWUI API
│   ├── commands/          # CLI commands (add, list)
│   └── formatters/        # Output formatters (table, JSON, CSV, TSV)
├── tests/
│   ├── integration/       # Integration tests
│   └── testdata/          # Test fixtures and mock server
├── specs/                 # Feature specifications
├── Makefile              # Build automation
└── README.md             # This file
```

## Technical Details

- **Language**: Go 1.24
- **CLI Framework**: Cobra
- **API Endpoints**:
  - EMWUI SetAutoAdd - Add automatic recording rules
  - EMWUI EnumAutoAdd - List and retrieve recording rules
  - EMWUI EnumService - List available channels/services
  - EMWUI EnumReserveInfo - List manual recording reservations
  - EMWUI EnumRecInfo - List recorded programs (paginated)
  - EMWUI EnumEventInfo - Retrieve EPG (program guide) data
- **Character Encoding**: UTF-8 (automatic URL encoding)
- **HTTP Timeout**: 10 seconds
- **CSRF Protection**: Automatically fetches ctok token from `/EMWUI/autoaddepg.html` before each request
- **Export Formats**: JSON, CSV, TSV with UTF-8 support
- **Filtering**: Client-side filtering with AND logic for multiple criteria

## Contributing

Contributions are welcome! Please ensure:
1. All tests pass: `make test`
2. Code follows Go conventions
3. Japanese character handling is tested

## License

[Add your license here]

## Acknowledgments

- EpgTimer project for the EMWUI interface
- Cobra CLI framework

---

**Version**: 0.4.0
**Last Updated**: 2025-12-23
