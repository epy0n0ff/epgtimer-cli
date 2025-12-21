# Quickstart Guide: EPG Reservation List Retrieval

**Feature**: EPG Reservation List Retrieval (001-epg-reservation-list)
**Estimated Reading Time**: 5 minutes
**Prerequisites**: EpgTimer with EMWUI running, epgtimer CLI installed

## What This Feature Does

This feature allows you to retrieve and view your automatic recording rules from EpgTimer via the command line. You can:
- List all automatic recording rules
- Filter rules by keywords, channels, or status
- Export rules to CSV, JSON, or TSV formats
- View rules in a human-readable table format

## Quick Examples

### 1. List All Rules (Basic Usage)

```bash
# List all automatic recording rules in table format
epgtimer list
```

**Output**:
```
ID  Enabled  Keywords             Exclusions  Channels
1   Yes      サイエンスZERO       [再]        62 channels
2   Yes      ブラタモリ           [再]        62 channels
3   Yes      うまいッ！           [再]        62 channels
4   Yes      NHKスペシャル        プレマップ  62 channels
```

---

### 2. Filter by Keyword

```bash
# Show only rules that contain "ニュース" in keywords
epgtimer list --andKey ニュース
```

**Output**:
```
ID  Enabled  Keywords             Exclusions  Channels
12  Yes      NHKニュース          [再]        15 channels
34  Yes      ニュース7           [再]        10 channels
```

---

### 3. Export to JSON

```bash
# Export all rules to JSON file
epgtimer list --format json --output rules.json
```

**Output** (rules.json):
```json
[
  {
    "id": 1,
    "search": {
      "disabled": 0,
      "and_key": "サイエンスZERO",
      "not_key": "[再]",
      "regex": 0,
      "channels": [
        {"onid": 32736, "tsid": 32736, "sid": 1024}
      ]
    },
    "recording": {
      "rec_mode": 1,
      "priority": 2
    }
  }
]
```

---

### 4. Export to CSV for Spreadsheet Analysis

```bash
# Export to CSV file
epgtimer list --format csv --output rules.csv
```

**Output** (rules.csv):
```csv
ID,Enabled,AndKey,NotKey,RegExp,Channels,ChannelCount,Priority,RecMode
1,true,サイエンスZERO,[再],false,32736-32736-1024;32736-32736-1025;...,61,2,1
```

You can now open `rules.csv` in Excel, Google Sheets, or any spreadsheet application.

---

### 5. Filter and Export

```bash
# Show only enabled rules and export to TSV
epgtimer list --enabled --format tsv --output enabled-rules.tsv
```

---

## Common Use Cases

### Use Case 1: Check What's Being Recorded

**Goal**: See all active recording rules to verify your setup

```bash
epgtimer list
```

Review the table output to confirm all desired programs are configured.

---

### Use Case 2: Find Specific Rule

**Goal**: Find rules for a specific program

```bash
# Search by keyword
epgtimer list --andKey ブラタモリ
```

---

### Use Case 3: Audit Disabled Rules

**Goal**: See which rules are disabled

```bash
# Show only disabled rules
epgtimer list --disabled
```

**Output**:
```
ID  Enabled  Keywords             Exclusions  Channels
56  No       古い番組             [再]        10 channels
78  No       テスト               なし        5 channels
```

---

### Use Case 4: Backup Rules to File

**Goal**: Export all rules for backup or analysis

```bash
# Export to JSON (preserves full structure)
epgtimer list --format json --output backup-$(date +%Y%m%d).json
```

Creates file like `backup-20251221.json` with all rule data.

---

### Use Case 5: Analyze Rules in Spreadsheet

**Goal**: Open rules in Excel for filtering/sorting

```bash
# Export to CSV
epgtimer list --format csv --output rules.csv
```

Then open `rules.csv` in your spreadsheet application:
1. Open Excel/Google Sheets
2. File → Open → Select `rules.csv`
3. Use built-in filters to sort by keyword, channel count, etc.

---

### Use Case 6: Check Channel-Specific Rules

**Goal**: See which rules apply to a specific channel

```bash
# Filter by channel (ONID-TSID-SID format)
epgtimer list --channel 32736-32736-1024
```

**Output**: Shows only rules that include channel 32736-32736-1024 in their serviceList.

---

## Configuration

### EMWUI Endpoint

The CLI needs to know where your EpgTimer EMWUI service is running.

**Option 1: Environment Variable** (Recommended)

```bash
export EMWUI_ENDPOINT="http://192.168.1.10:5510"
epgtimer list
```

**Option 2: Command-Line Flag**

```bash
epgtimer list --endpoint http://192.168.1.10:5510
```

**Default**: If not specified, CLI will look for `EMWUI_ENDPOINT` environment variable.

---

## All Available Flags

### Filter Flags

| Flag | Description | Example |
|------|-------------|---------|
| `--andKey <keyword>` | Filter by search keyword (substring match) | `--andKey ニュース` |
| `--channel <id>` | Filter by channel (ONID-TSID-SID format) | `--channel 32736-32736-1024` |
| `--enabled` | Show only enabled rules | `--enabled` |
| `--disabled` | Show only disabled rules | `--disabled` |
| `--regex` | Show only regex-enabled rules | `--regex` |

### Output Flags

| Flag | Description | Example |
|------|-------------|---------|
| `--format <fmt>` | Output format: table, json, csv, tsv | `--format json` |
| `--output <file>` | Output file path ("-" for stdout) | `--output rules.csv` |

### Connection Flags

| Flag | Description | Example |
|------|-------------|---------|
| `--endpoint <url>` | EMWUI service endpoint | `--endpoint http://localhost:5510` |

---

## Combining Filters

You can combine multiple filters to narrow down results:

```bash
# Enabled rules with keyword "ドラマ"
epgtimer list --enabled --andKey ドラマ

# Regex rules for specific channel
epgtimer list --regex --channel 32736-32736-1024

# Enabled rules exported to CSV
epgtimer list --enabled --format csv --output enabled.csv
```

**Filter Logic**: All filters are AND-ed together (rule must match ALL specified filters).

---

## Output Formats Comparison

| Format | Use Case | Pros | Cons |
|--------|----------|------|------|
| **table** | Quick terminal viewing | Human-readable, no extra tools | Truncated long values |
| **json** | Programmatic processing | Full structure, nested data | Verbose, needs parser |
| **csv** | Spreadsheet analysis | Excel/Sheets compatible | Flattened structure |
| **tsv** | Data import/export | Simple delimiter | Less common than CSV |

**Recommendation**: Use `table` for quick checks, `json` for backups, `csv` for spreadsheet analysis.

---

## Troubleshooting

### Problem: "Failed to connect to EMWUI service"

**Cause**: EMWUI endpoint not configured or unreachable

**Solution**:
1. Check EpgTimer is running
2. Verify endpoint URL: `echo $EMWUI_ENDPOINT`
3. Test connection: `curl http://192.168.1.10:5510/api/EnumAutoAdd`

---

### Problem: "No automatic recording rules found"

**Cause**: No rules configured OR all rules filtered out

**Solution**:
1. Check EpgTimer UI to confirm rules exist
2. Remove filters and try: `epgtimer list` (no flags)
3. If still empty, add rules via EpgTimer UI or `epgtimer add` command

---

### Problem: Japanese characters display as "???"

**Cause**: Terminal doesn't support UTF-8

**Solution**:
1. Ensure terminal locale is UTF-8: `locale` (should show UTF-8)
2. Set locale: `export LANG=ja_JP.UTF-8` or `export LANG=en_US.UTF-8`
3. Use modern terminal (Windows Terminal, iTerm2, GNOME Terminal)

---

### Problem: "Failed to write output to file"

**Cause**: Permission denied or disk full

**Solution**:
1. Check file permissions: `ls -la rules.csv`
2. Try different directory: `epgtimer list --format csv --output ~/rules.csv`
3. Check disk space: `df -h`

---

## Next Steps

### Related Commands

- `epgtimer add`: Add new automatic recording rule
- `epgtimer help`: Show all available commands
- `epgtimer list --help`: Show detailed help for list command

### Advanced Usage

**Scripting Example**: Daily backup of rules

```bash
#!/bin/bash
# backup-rules.sh - Daily backup of recording rules

DATE=$(date +%Y%m%d)
BACKUP_DIR="$HOME/epgtimer-backups"
mkdir -p "$BACKUP_DIR"

epgtimer list --format json --output "$BACKUP_DIR/rules-$DATE.json"

echo "Backup saved to $BACKUP_DIR/rules-$DATE.json"
```

**Cron Job**: Run daily at 3 AM

```cron
0 3 * * * /path/to/backup-rules.sh
```

---

**Filtering in Scripts**:

```bash
#!/bin/bash
# check-disabled-rules.sh - Alert if disabled rules found

COUNT=$(epgtimer list --disabled --format json | jq '. | length')

if [ "$COUNT" -gt 0 ]; then
    echo "WARNING: $COUNT disabled recording rules found"
    epgtimer list --disabled
fi
```

---

## Tips & Best Practices

### Tip 1: Use Aliases for Common Operations

```bash
# Add to ~/.bashrc or ~/.zshrc
alias epl='epgtimer list'
alias epl-json='epgtimer list --format json'
alias epl-enabled='epgtimer list --enabled'

# Usage
epl                    # Quick list
epl-json               # JSON output
epl-enabled            # Only enabled rules
```

---

### Tip 2: Pipe JSON to jq for Analysis

```bash
# Count total rules
epgtimer list --format json | jq '. | length'

# Extract only keywords
epgtimer list --format json | jq '.[].search.and_key'

# Filter in jq (enabled rules with priority > 2)
epgtimer list --format json | jq '.[] | select(.search.disabled == 0 and .recording.priority > 2)'
```

---

### Tip 3: Compare Before/After Changes

```bash
# Before making changes
epgtimer list --format json --output before.json

# ... make changes in EpgTimer UI ...

# After changes
epgtimer list --format json --output after.json

# Compare
diff before.json after.json
```

---

### Tip 4: Regular Backups

Set up automated backups to avoid losing rules:

```bash
# Weekly backup (keeps last 4 weeks)
epgtimer list --format json --output "backup-$(date +%Y-week%V).json"
```

---

## FAQ

**Q: Can I modify rules using the list command?**
A: No, `list` is read-only. Use EpgTimer UI or (future) `epgtimer update` command to modify rules.

**Q: How many rules can I retrieve?**
A: No limit. Tested with 230+ rules successfully. Performance may degrade beyond 1000 rules.

**Q: Does filtering happen on the server or client?**
A: Client-side. The API returns all rules, then CLI filters locally.

**Q: Can I list manual reservations?**
A: Not with this command. This lists automatic recording rules only. Use `/api/EnumReserveInfo` API for manual reservations (feature not yet implemented).

**Q: What if my EpgTimer uses a different port?**
A: Specify the full endpoint: `--endpoint http://192.168.1.10:CUSTOM_PORT`

**Q: Can I filter by multiple keywords?**
A: Not directly. Use JSON output + jq for complex filtering:
```bash
epgtimer list --format json | jq '.[] | select(.search.and_key | contains("ニュース") or contains("ドラマ"))'
```

---

## Getting Help

**Command Help**:
```bash
epgtimer list --help
```

**Full Documentation**:
- Feature Specification: `specs/001-epg-reservation-list/spec.md`
- API Contract: `specs/001-epg-reservation-list/contracts/api-contract.md`
- CLI Contract: `specs/001-epg-reservation-list/contracts/cli-interface.md`

**Issues**:
Report bugs or request features via project issue tracker.

---

## Summary

**Basic Command**:
```bash
epgtimer list
```

**With Filters**:
```bash
epgtimer list --enabled --andKey ニュース
```

**Export**:
```bash
epgtimer list --format csv --output rules.csv
```

**Remember**: `--format` (table|json|csv|tsv) and `--output` for file export, various filter flags to narrow results.
