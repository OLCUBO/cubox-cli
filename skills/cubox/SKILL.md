---

## name: cubox
version: 1.0.0
description: "Cubox CLI: manage Cubox bookmarks — list groups/tags, filter cards, read card content. Use when the user wants to search, browse, or read their Cubox bookmarks, or needs to query Cubox data from the CLI."

# cubox-cli

Manage Cubox bookmarks via the `cubox-cli` command-line tool.

## First-time Setup

If `cubox-cli` is not installed, install it:

```bash
npm install -g cubox-cli
```

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

### Filter Cards

```bash
cubox-cli card list [flags]
```

Flags:

- `--group ID,...` — filter by group IDs
- `--type Article,Snippet,...` — filter by type (Article, Snippet, Memo, Image, Audio, Video, File)
- `--tag ID,...` — filter by tag IDs
- `--starred` — starred cards only
- `--read` / `--unread` — filter by read status
- `--annotated` — cards with highlights only
- `--limit N` — page size (default 50)
- `--cursor CARD_ID,UPDATE_TIME` — resume pagination from last card
- `--all` — auto-paginate all results

Returns: `[{ "id", "title", "description", "type", "tags", "url", "highlights", ... }]`

### Get Card Content

```bash
cubox-cli card content --id CARD_ID
```

Returns the full article content as markdown text.

## Common Workflows

### Browse and read a bookmark

```bash
# 1. List groups to find the target folder
cubox-cli group list

# 2. Filter cards in that group
cubox-cli card list --group GROUP_ID --limit 10

# 3. Read the content
cubox-cli card content --id CARD_ID
```

### Find starred articles

```bash
cubox-cli card list --starred --type Article
```

### Export all annotated cards

```bash
cubox-cli card list --annotated --all
```

## Notes

- Pagination uses cursor-based approach. When `--all` is not used, the response may be partial. Use `--cursor` with the last card's ID and update_time to continue.
- The `nested_name` field in groups and tags shows the full path (e.g. `"Parent/Child"`).
- Config is stored at `~/.config/cubox-cli/config.json`.

