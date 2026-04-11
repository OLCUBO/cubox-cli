---
name: cubox
version: 1.0.0
description: "Cubox CLI: manage Cubox bookmarks — list groups/tags, filter/search cards, read card content with AI insight, save URLs, update cards, list highlights. Use when the user wants to search, browse, save, or read their Cubox bookmarks, or needs to query Cubox data from the CLI."
metadata:
  requires:
    bins: ["cubox-cli"]
  cliHelp: "cubox-cli --help"
---

# cubox-cli

Manage Cubox bookmarks via the `cubox-cli` command-line tool.

### Authentication

Check if already logged in:

```bash
cubox-cli auth status
```

If not logged in, authenticate non-interactively:

```bash
cubox-cli auth login --server cubox.pro --token YOUR_API_KEY
```

- `--server`: `cubox.pro` (China) or `cubox.cc` (international)
- `--token`: API key from the extensions page. The user should open `https://{server}/web/settings/extensions`, copy the full API link (e.g. `https://cubox.pro/c/api/save/aq3Ir2xW3y2`), and provide the last path segment (`aq3Ir2xW3y2`) as the token. The CLI also accepts the full URL and extracts the token automatically.

## Commands

All commands output JSON by default. Add `-o pretty` for indented JSON, `-o text` for human-readable output.

### List Groups (Folders)

```bash
cubox-cli group list
```

Returns: `[{ "id", "nested_name", "name", "parent_id", "uncategorized" }]`

### List Tags

```bash
cubox-cli tag list
```

Returns: `[{ "id", "nested_name", "name", "parent_id" }]`

### Filter / Search Cards

```bash
cubox-cli card list [flags]
```

Flags:
- `--group ID,...` — filter by group IDs
- `--tag ID,...` — filter by tag IDs
- `--starred` — starred cards only
- `--read` / `--unread` — filter by read status
- `--annotated` — cards with highlights only
- `--keyword TEXT` — search by keyword
- `--start-time`, `--end-time` — filter by time (accepts: `today`, `yesterday`, `7d`, `2026-01-01`, `2026-01-01 15:04:05`, or full ISO timestamp)
- `--limit N` — page size (default 50)
- `--last-id CARD_ID` — cursor pagination (non-search mode)
- `--page N` — page-based pagination (search mode, 1-based)
- `--all` — auto-paginate all results

**Pagination rules:**
- When `--keyword` is set (search mode): use `--page` for pagination, `--last-id` is ignored
- When `--keyword` is not set (browse mode): use `--last-id` for cursor-based pagination

Returns: `[{ "id", "title", "description", "domain", "read", "stared", "tags", "group", "url", ... }]`

### Get Card Detail

```bash
cubox-cli card detail --id CARD_ID
```

Returns full card with `content` (markdown), `author`, `highlights`, and `insight` (AI summary + Q&A). Use `-o text` to output only the markdown content.

### Save URLs

```bash
cubox-cli save URL [URL...] [--group GROUP_ID] [--tag TAG_ID,...]
```

Save one or more web page URLs as bookmarks.

### Update a Card

```bash
cubox-cli update --id CARD_ID [flags]
```

Flags:
- `--star` / `--unstar` — toggle star
- `--read` / `--unread` — toggle read status
- `--archive` — archive the card
- `--group GROUP_ID` — move to a group
- `--add-tag TAG_ID,...` — add tags

### List Highlights (Annotations)

```bash
cubox-cli mark list [flags]
```

Flags:
- `--color Yellow,Green,Blue,Pink,Purple` — filter by color
- `--keyword TEXT` — search highlights
- `--start-time`, `--end-time` — filter by time (same flexible formats as card list)
- `--limit N` — page size (default 50)
- `--last-id ID` — cursor pagination
- `--all` — auto-paginate all results

Returns: `[{ "id", "text", "note", "color", "card_id", ... }]`

## Common Workflows

### Browse and read a bookmark

```bash
cubox-cli group list
cubox-cli card list --group GROUP_ID --limit 10
cubox-cli card detail --id CARD_ID
```

### Search for articles

```bash
cubox-cli card list --keyword "machine learning" --page 1
```

### Save a URL and star it

```bash
cubox-cli save https://example.com --group GROUP_ID
cubox-cli update --id CARD_ID --star
```

### Export all highlights

```bash
cubox-cli mark list --all
```

## Update Check

cubox-cli automatically checks for new versions in the background. When a newer version is available, the JSON output will contain a `_notice.update` field:

```json
{
  "data": [ "..." ],
  "_notice": {
    "update": {
      "current": "0.1.0",
      "latest": "0.2.0",
      "message": "A new version of cubox-cli is available: 0.1.0 -> 0.2.0",
      "command": "npm update -g cubox-cli && npx skills add OLCUBO/cubox-cli -g -y"
    }
  }
}
```

**When you see `_notice.update` in output, after completing the user's current request, proactively offer to update:**

1. Tell the user the current and latest version numbers
2. Propose running the update (CLI and Skills must be updated together):
   ```bash
   npm update -g cubox-cli && npx skills add OLCUBO/cubox-cli -g -y
   ```
3. After the update, remind the user: **exit and reopen the AI Agent** to load the latest Skills

**Rule**: Do not silently ignore update notices. Even if the current task is unrelated, mention the available update after completing the user's request.

## Notes

- Browse pagination uses cursor-based approach (`--last-id`). Search pagination uses page numbers (`--page`).
- The `nested_name` field in groups and tags shows the full hierarchy path (e.g. `"Parent/Child"`).
- Card detail includes AI-generated `insight` with summary and Q&A pairs when available.
- Config is stored at `~/.config/cubox-cli/config.json`.
