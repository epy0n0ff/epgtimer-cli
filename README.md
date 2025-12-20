# EpgTimer CLI

Command-line interface for EpgTimer's EMWUI to manage automatic recording rules via keyword-based searches.

## Features

- Add automatic recording rules based on search keywords
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

### Basic Command

```bash
epgtimer add --andKey "search keywords" --serviceList "ONID-TSID-SID"
```

### Examples

#### Example 1: Record all news programs on NHK

```bash
epgtimer add --andKey "ニュース" --serviceList "32736-32736-1024"
```

#### Example 2: Record dramas, excluding reruns

```bash
epgtimer add \
  --andKey "ドラマ" \
  --notKey "再放送" \
  --serviceList "32736-32736-1024,32737-32737-1032,32738-32738-1040"
```

#### Example 3: Record movies on multiple channels

```bash
epgtimer add \
  --andKey "映画" \
  --serviceList "32736-32736-1024,32736-32736-1025,32737-32737-1032,32738-32738-1040"
```

#### Example 4: Record specific show with exclusions

```bash
epgtimer add \
  --andKey "わたしが恋人になれるわけ" \
  --notKey "推しエンタ" \
  --serviceList "32736-32736-1024,32736-32736-1025"
```

### Command Options

- `--andKey` (required): Search keywords - programs must contain these keywords in the title
- `--notKey` (optional): Exclusion keywords - programs must NOT contain these keywords
- `--serviceList` (required): Comma-separated list of channels in "ONID-TSID-SID" format
- `--endpoint` (optional): Override EMWUI_ENDPOINT environment variable

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
│   ├── models/            # Data models
│   ├── client/            # HTTP client for EMWUI API
│   └── commands/          # CLI commands
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
- **API**: EMWUI SetAutoAdd endpoint
- **Character Encoding**: UTF-8 (automatic URL encoding)
- **HTTP Timeout**: 10 seconds
- **CSRF Protection**: Automatically fetches ctok token from `/EMWUI/autoaddepg.html` before each request

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

**Version**: 0.1.0
**Last Updated**: 2025-12-20
