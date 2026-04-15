---

name: cubox
version: 1.0.0
description: "Cubox CLI is a callable personal reading memory system that enables you to search, read, and use saved content, perform semantic (RAG-based) queries, access articles, highlights, and metadata, save URLs, update content states, and retrieve annotations and structure such as folders and tags. Use this tool when a task depends on the user’s reading history or requires context from their Cubox library."
metadata:
  requires:
    bins: ["cubox-cli"]
  cliHelp: "cubox-cli --help"
---

# cubox-cli

Manage Cubox bookmarks via the `cubox-cli` command-line tool.

## Authentication

If any command fails with "not logged in", run `cubox-cli auth login` and follow the interactive prompts.

## Commands

All commands output JSON by default. Add `-o pretty` for indented JSON, `-o text` for human-readable output.

### List Folders

```bash
cubox-cli folder list
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

- `--folder ID,...` — filter by folder IDs
- `--tag ID,...` — filter by tag IDs
- `--starred` — starred cards only
- `--read` / `--unread` — filter by read status
- `--annotated` — cards with annotations only
- `--keyword TEXT` — search by keyword
- `--start-time`, `--end-time` — filter by time range (see **Time filtering** below)
- `--limit N` — page size (default 50)
- `--last-id CARD_ID` — cursor pagination (non-search mode)
- `--page N` — page-based pagination (search mode, 1-based)
- `--all` — auto-paginate all results

**Pagination rules:**

- When `--keyword` is set (search mode): use `--page` for pagination, `--last-id` is ignored
- When `--keyword` is not set (browse mode): use `--last-id` for cursor-based pagination

Returns: `[{ "id", "title", "description", "domain", "read", "starred", "tags", "folder", "url", ... }]`

### Get Card Detail

```bash
cubox-cli card detail --id CARD_ID
```

Returns full card with `content` (markdown), `author`, `annotations`, and `insight` (AI summary + Q&A). Use `-o text` to output only the markdown content.

### RAG Semantic Search

```bash
cubox-cli card rag --query "QUERY_TEXT"
```

Semantic search via natural language. Unlike `--keyword`, RAG understands intent and returns conceptually relevant cards. **[Must-read: RAG workflow](references/card-rag-workflow.md)** — covers when to use RAG vs keyword, query refinement, progressive detail fetching, and re-ranking.

Returns: `[{ "id", "title", "description", "domain", "tags", "folder", "url", ... }]` (same Card shape as `card list`)

### Save Web Pages

```bash
cubox-cli save URL [URL...] [--title TEXT] [--desc TEXT] [--folder NAME] [--tag NAME,...]
cubox-cli save --json '[{"url":"...","title":"...","description":"..."}]' [--folder NAME] [--tag NAME,...]
```

Save one or more web pages as bookmarks. Three input modes:

- **URL arguments** — simple: `cubox-cli save https://example.com https://b.com`
- **Single with metadata** — `cubox-cli save https://example.com --title "My Page" --desc "A description"`
- **Batch via JSON** — `cubox-cli save --json '[{"url":"https://a.com","title":"Title A"}]'`

Folders and tags are specified **by name** (not ID), including nested paths like `"parent/child"`.

### Update a Card

```bash
cubox-cli update --id CARD_ID [flags]
```

Flags:

- `--star` / `--unstar` — toggle star
- `--read` / `--unread` — toggle read status
- `--archive` — archive the card
- `--folder NAME` — move to folder by name (e.g. `"parent/child"`; `""` = Uncategorized)
- `--tag NAME,...` — set tags by name (comma-separated, supports nested like `"parent/child"`)
- `--title TEXT` — update title
- `--description TEXT` — update description

Folders and tags are specified **by name** (not ID). No need to query IDs first.

### Delete Cards

```bash
cubox-cli delete --id CARD_ID [--id ID2,...] [--dry-run]
```

Delete cards by ID. **Always `--dry-run` first.** **[Must-read: Dry Run Policy](references/card-delete.md)** — agents must preview before deleting.

### List Annotations

```bash
cubox-cli annotation list [flags]
```

Flags:

- `--color Yellow,Green,Blue,Pink,Purple` — filter by color
- `--keyword TEXT` — search annotations
- `--start-time`, `--end-time` — filter by time range (same formats and rules as card list)
- `--limit N` — page size (default 50)
- `--last-id ID` — cursor pagination
- `--all` — auto-paginate all results

Returns: `[{ "id", "text", "note", "color", "card_id", ... }]`

### Cubox Deep Links

Construct clickable Cubox links from any resource ID (card, folder, tag). No API call needed — just the ID + server. **[Must-read: Deep Links](references/deep-links.md)** — URL patterns, scheme rules, and examples.

Default: `https://{server}/web/card/{ID}` — use `cubox://` scheme only when explicitly requested.

## Time filtering

`--start-time` and `--end-time` accept flexible shorthand values. The CLI automatically resolves day-level inputs to the correct boundary:

- `--start-time` resolves to **start of day** (00:00:00.000)
- `--end-time` resolves to **end of day** (23:59:59.999)

Accepted formats: `today`, `yesterday`, `now`, `7d` (7 days ago), `2026-01-01`, `2026-01-01 15:04:05`, or full ISO timestamp.

Common time query patterns:

| Intent | Command |
|--------|---------|
| Today's cards | `--start-time today --end-time today` |
| Yesterday's cards | `--start-time yesterday --end-time yesterday` |
| Last 7 days | `--start-time 7d --end-time today` |
| Since a date | `--start-time 2026-01-01` |
| Up to now | `--end-time now` |

## Common Workflows

### Browse and read a card detail

```bash
cubox-cli folder list
cubox-cli card list --folder FOLDER_ID --limit 10
cubox-cli card detail --id CARD_ID
```

### Search for articles

```bash
cubox-cli card list --keyword "machine learning" --page 1
```

### Save a page and star it

```bash
cubox-cli save https://example.com --title "Example" --folder "Reading List"
cubox-cli update --id CARD_ID --star
```

### List cards with Cubox links

```bash
cubox-cli auth status          # determine server (cubox.pro or cubox.cc)
cubox-cli card list --limit 5  # get cards, then append link from ID
# For card ID 7247925101516031380 on cubox.pro:
# → https://cubox.pro/web/card/7247925101516031380
```

### Export all annotations

```bash
cubox-cli annotation list --all
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

## Security Rules

- Never expose sensitive credentials in plain text (API key/token, session data, auth headers).
- Treat Cubox API tokens as local secrets. Do not commit or copy them into repository files, screenshots, or shared notes.
- Before any write/destructive action (`save`, `update`, `delete`), confirm user intent first. For deletion, always run `--dry-run` and present the preview before execution.
- When demonstrating commands, use placeholders (for example `YOUR_API_KEY`) instead of real values.
- Avoid leaving secrets in shell history where possible (for example, prefer temporary environment variables and clear them after use).
- If credentials are suspected to be leaked, instruct the user to rotate the Cubox API token from the extensions page immediately.

## Notes

- Browse pagination uses cursor-based approach (`--last-id`). Search pagination uses page numbers (`--page`).
- The `nested_name` field in folders and tags shows the full hierarchy path (e.g. `"Parent/Child"`).
- Card detail includes AI-generated `insight` with summary and Q&A pairs when available.
- Config is stored at `~/.config/cubox-cli/config.json`.

