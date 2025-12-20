# Quick Start Guide: EpgTimer CLI

**Feature**: `001-epgtimer-cli`
**Created**: 2025-12-20

---

## Installation

### Option 1: Go Install

```bash
go install github.com/epy0n0ff/epgtimer-cli/cmd/epgtimer@latest
```

### Option 2: Binary Download

Download the pre-built binary for your platform from the releases page.

```bash
# Linux/macOS
chmod +x epgtimer
sudo mv epgtimer /usr/local/bin/

# Windows
# Add epgtimer.exe to your PATH
```

---

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

---

## Usage

### Add Automatic Recording Rule

Create a rule to automatically record programs matching keywords:

```bash
epgtimer add \
  --andKey "ニュース" \
  --serviceList "32736-32736-1024,32736-32736-1025"
```

### With Exclusion Keywords

Exclude certain programs from the rule:

```bash
epgtimer add \
  --andKey "ドラマ" \
  --notKey "再放送" \
  --serviceList "32736-32736-1024,32737-32737-1032"
```

### Multiple Channels

Specify multiple channels to search across:

```bash
epgtimer add \
  --andKey "映画" \
  --serviceList "32736-32736-1024,32736-32736-1025,32737-32737-1032,32738-32738-1040"
```

---

## Common Channel IDs (Tokyo Area)

```
NHK総合:        32736-32736-1024
NHK教育:        32736-32736-1025
日本テレビ:     32737-32737-1032
TBS:           32738-32738-1040
フジテレビ:     32739-32739-1048
テレビ朝日:     32740-32740-1056
テレビ東京:     32741-32741-1064
```

**Note**: Channel IDs vary by region. Use the values from your EMWUI interface.

---

## Examples

### Example 1: Record all news programs on NHK

```bash
epgtimer add --andKey "ニュース" --serviceList "32736-32736-1024"
```

### Example 2: Record dramas, excluding reruns

```bash
epgtimer add \
  --andKey "ドラマ" \
  --notKey "再放送" \
  --serviceList "32736-32736-1024,32737-32737-1032,32738-32738-1040"
```

### Example 3: Record specific show on multiple channels

```bash
epgtimer add \
  --andKey "わたしが恋人になれるわけ" \
  --notKey "推しエンタ" \
  --serviceList "32736-32736-1024,32736-32736-1025"
```

---

## Troubleshooting

### Connection Failed

**Error**: "Failed to connect to EMWUI service"

**Solutions**:
1. Check EMWUI_ENDPOINT is set correctly: `echo $EMWUI_ENDPOINT`
2. Verify EpgTimer service is running
3. Test connection: `curl http://localhost:5510/api/`

### Invalid Parameters

**Error**: "andKey is required"

**Solution**: Ensure you provide the `--andKey` parameter with a search keyword

**Error**: "serviceList is required"

**Solution**: Provide at least one channel in `--serviceList`

### Character Encoding Issues

**Problem**: Japanese characters not working

**Solution**: Ensure your terminal supports UTF-8 encoding. The CLI automatically handles URL encoding.

---

## Help

View all available options:

```bash
epgtimer --help
epgtimer add --help
```

---

## Next Steps

- Check EpgTimer GUI to verify rules were created
- Monitor recordings to ensure keywords match desired programs
- Adjust keywords and channel lists as needed

---

Last Updated: 2025-12-20
