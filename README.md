# cubox-cli

[License: MIT](https://opensource.org/licenses/MIT)
[Go Version](https://go.dev/)
[npm version](https://www.npmjs.com/package/cubox-cli)

[English](./README.md) | [中文版](./README.zh.md)

The official Cubox CLI. Save, search, read, and use what you read with AI. Your personal reading memory, now usable.

[Install](#installation) · [Auth](#authentication) · [Commands](#commands) · [AI Agent](#quick-start-ai-agent) · [Examples](#examples) · [Development](#development)

## Features


| Category | Capabilities                                                                          |
| -------- | ------------------------------------------------------------------------------------- |
| Folders  | List and browse card folders                                                      |
| Tags     | List and browse tag hierarchy; rename, batch delete, and merge tags                   |
| Cards    | Filter/search cards by folder, tag, starred/read/annotated/archived status, keyword, time range |
| RAG      | Semantic search via natural language query (intent-based retrieval)                    |
| Content  | Read full card detail with article content (markdown), annotations, and AI insight    |
| Save     | Save web pages with optional title/description, batch via JSON            |
| Update   | Star/unstar, mark read/unread, move to folder, manage tags                            |
| Archive  | Batch archive cards or restore (unarchive) them into a folder                         |
| Delete   | Delete cards by ID, with dry-run preview support                             |
| Annotations | List and search annotations across all cards                                      |


## Installation

### Requirements

- Node.js 16+ (for `npm`/`npx`)
- Go 1.23+ (only for building from source)

### Quick Start (Human Users)

> **For AI Agents:** jump to [Quick Start (AI Agent)](#quick-start-ai-agent) below.

#### Install

**Option 1 — From npm (recommended):**

```bash
# Install CLI
npm install -g cubox-cli

# Install CLI Skill (required for AI Agents)
npx skills add OLCUBO/cubox-cli -g -y
```

**Option 2 — From source:**

```bash
git clone https://github.com/OLCUBO/cubox-cli.git
cd cubox-cli
make install

# Install CLI Skill (required for AI Agents)
npx skills add OLCUBO/cubox-cli -g -y
```

#### Log in

```bash
cubox-cli auth login
```

The CLI will guide you through:

1. **Sign in to Cubox web** associated with your account:
   - For .pro users: [https://cubox.pro/web/settings/extensions](https://cubox.pro/web/settings/extensions)
   - For international .cc users: [https://cubox.cc/web/settings/extensions](https://cubox.cc/web/settings/extensions)
2. **Enable the API Extension** — go to Extensions, locate the API Extension, enable it, and copy your unique link (e.g. `https://cubox.pro/c/api/save/abcd12345`).
3. **Paste the link** — the CLI extracts the token and server automatically.

#### Start using

```bash
# List your folders
cubox-cli folder list

# List your tags
cubox-cli tag list

# Browse recent cards
cubox-cli card list --limit 5 -o pretty

# Read an article
cubox-cli card detail --id CARD_ID
```

### Quick Start (AI Agent)

> These steps are for AI Agents helping users set up cubox-cli.

**Step 1 — Install**

```bash
# Install CLI
npm install -g cubox-cli

# Install CLI Skill (required)
npx skills add OLCUBO/cubox-cli -g -y
```

**Step 2 — Authenticate**

Ask the user to sign in to their Cubox web and get the API link:

1. Please sign in to the Cubox web associated with your account:
   - For .pro users: `https://cubox.pro/web/settings/extensions`
   - For international .cc users: `https://cubox.cc/web/settings/extensions`
2. Go to Extensions, locate the API Extension, enable it, and copy your unique link (e.g., `https://cubox.pro/c/api/save/abcd12345`).

The agent **must not** ask the user to paste the token into chat or embed it into argv. Instead, direct the user to run one of these themselves in their terminal:

```bash
# Interactive (recommended for humans)
cubox-cli auth login

# Agent / CI without persistence — token lives only in the shell environment
export CUBOX_SERVER=cubox.pro
export CUBOX_TOKEN=... # or the full API link URL
cubox-cli folder list

# Non-interactive persisted login — token piped via stdin, never in argv/ps/history
printf '%s' "$TOKEN" | cubox-cli auth login --server cubox.pro --token-stdin
```

The legacy `--token TOKEN` flag still works but exposes the token to shell history and `ps`; only use it in controlled environments.

**Step 3 — Verify**

```bash
cubox-cli auth status
```

**Step 4 — Use**

```bash
cubox-cli folder list
cubox-cli card list --limit 10
```

## Authentication


| Command                                                              | Description                                                              |
| -------------------------------------------------------------------- | ------------------------------------------------------------------------ |
| `cubox-cli auth login`                                               | Interactive login (server selection + token input)                       |
| `printf '%s' "$TOKEN" \| cubox-cli auth login --server cubox.pro --token-stdin` | Non-interactive login via stdin (recommended for agents)       |
| `CUBOX_SERVER=cubox.pro CUBOX_TOKEN=... cubox-cli ...`               | Transient env-var auth (no file written; recommended for CI / sandboxes) |
| `cubox-cli auth login --server cubox.pro --token TOKEN`              | Legacy argv form — leaks to shell history/ps, avoid when possible        |
| `cubox-cli auth status`                                              | Show current server, masked token, connection test                       |
| `cubox-cli auth logout`                                              | Remove saved credentials                                                 |


Credentials are stored at `~/.config/cubox-cli/config.json`. `CUBOX_TOKEN` and `CUBOX_SERVER` override the on-disk config when set, and are enough on their own when no config file exists.

## Commands

### Output Formats

All commands support the `-o` / `--output` flag:


| Flag        | Description                            |
| ----------- | -------------------------------------- |
| `-o json`   | Compact JSON (default, agent-friendly) |
| `-o pretty` | Indented JSON                          |
| `-o text`   | Human-readable text/tree output        |


### `cubox-cli folder list`

List all folders.

```bash
cubox-cli folder list
cubox-cli folder list -o text
```

**JSON output fields:** `id`, `nested_name`, `name`, `parent_id`, `uncategorized`

### `cubox-cli tag list`

List all tags.

```bash
cubox-cli tag list
cubox-cli tag list -o text
```

**JSON output fields:** `id`, `nested_name`, `name`, `parent_id`

### `cubox-cli tag update`

Rename a tag by ID. The new name applies to the leaf segment only — nested children stay attached and reachable through the new path automatically.

```bash
cubox-cli tag update --id TAG_ID --new-name NEW_NAME
```

| Flag                | Description                                          |
| ------------------- | ---------------------------------------------------- |
| `--id ID`           | Tag ID to rename (required)                          |
| `--new-name NAME`   | New leaf name (required; must not contain `/`)       |

### `cubox-cli tag delete`

Batch delete one or more tags by ID. Cards previously tagged with the deleted tag(s) remain; only the tag-card association is removed.

```bash
cubox-cli tag delete --id TAG_ID[,ID2,...]
```

| Flag          | Description                                          |
| ------------- | ---------------------------------------------------- |
| `--id ID,...` | Tag IDs to delete (comma-separated, required)        |

### `cubox-cli tag merge`

Merge one or more source tags into a target tag. All cards associated with the source tags are re-tagged onto the target, and the source tags are then deleted.

```bash
cubox-cli tag merge --source SRC_ID[,ID2,...] --target TARGET_ID
```

| Flag              | Description                                          |
| ----------------- | ---------------------------------------------------- |
| `--source ID,...` | Source tag IDs to merge (comma-separated, required)  |
| `--target ID`     | Target tag ID to merge into (required)               |

### `cubox-cli card list`

Filter and search cards. Supports keyword search with page-based pagination, and cursor-based pagination for browsing.

```bash
cubox-cli card list [flags]
```


| Flag                | Description                                                        |
| ------------------- | ------------------------------------------------------------------ |
| `--folder ID,...`   | Filter by folder IDs                                               |
| `--tag ID,...`      | Filter by tag IDs                                                  |
| `--starred`         | Only starred cards                                                 |
| `--read`            | Only read cards                                                    |
| `--unread`          | Only unread cards                                                  |
| `--annotated`       | Only cards with annotations                                        |
| `--archived`        | Only archived cards (default: non-archived only)                   |
| `--keyword TEXT`    | Search by keyword                                                  |
| `--start-time TIME` | Filter by create time start (e.g. `2026-01-01T00:00:00:000+08:00`) |
| `--end-time TIME`   | Filter by create time end                                          |
| `--limit N`         | Page size (default 50)                                             |
| `--last-id CARD_ID` | Cursor pagination for browsing (non-search)                        |
| `--page N`          | Page number for search (1-based, used with `--keyword`)            |
| `--all`             | Auto-paginate to fetch all results                                 |


**Pagination:** When `--keyword` is set, use `--page` for pagination. Otherwise, use `--last-id` with the last card's ID.

### `cubox-cli card detail --id ID`

Get full card detail including article content (markdown), author, annotations, and AI-generated insight (summary + Q&A).

```bash
cubox-cli card detail --id 7247925101516031380
cubox-cli card detail --id 7247925101516031380 -o pretty
```

Use `-o text` to output only the markdown content.

### `cubox-cli card rag --query TEXT`

Semantic search using natural language via RAG (Retrieval-Augmented Generation). Unlike keyword search, RAG understands intent and returns conceptually relevant cards even when exact words don't match.

```bash
cubox-cli card rag --query "Java实现数据库图片上传功能"
cubox-cli card rag --query "how to build a REST API with authentication" -o pretty
```

| Flag            | Description                            |
| --------------- | -------------------------------------- |
| `--query TEXT`  | Natural language query text (required) |

**When to use RAG vs keyword search:**

- **`card list --keyword`** — exact terms, known titles, domain names, short phrases
- **`card rag --query`** — questions, topic exploration, conceptual or fuzzy queries

### `cubox-cli save`

Save one or more web pages as bookmarks. Supports three input modes.

```bash
# Simple — URL arguments
cubox-cli save https://example.com
cubox-cli save https://a.com https://b.com --folder "Reading List"

# Single with metadata
cubox-cli save https://example.com --title "My Page" --desc "Interesting read"

# Batch via JSON
cubox-cli save --json '[{"url":"https://a.com","title":"Title A"},{"url":"https://b.com"}]' --tag tech,AI/LLM
```

| Flag             | Description                                                         |
| ---------------- | ------------------------------------------------------------------- |
| `--title TEXT`   | Title for the saved page (single URL mode only)                     |
| `--desc TEXT`    | Description for the saved page (single URL mode only)               |
| `--json JSON`    | Batch card entries as JSON array `[{"url","title","description"}]`   |
| `--folder NAME`  | Target folder by name (e.g. `"parent/child"`)                       |
| `--tag NAME,...` | Tag names (comma-separated, supports nested like `"parent/child"`)  |

### `cubox-cli update`

Update a card's properties.

```bash
cubox-cli update --id CARD_ID [flags]
```


| Flag                   | Description                                                            |
| ---------------------- | ---------------------------------------------------------------------- |
| `--star` / `--unstar`  | Toggle star                                                            |
| `--read` / `--unread`  | Toggle read status                                                     |
| `--folder NAME`        | Move to folder by name (e.g. `"parent/child"`; `""` = Uncategorized)   |
| `--tag NAME,...`       | Replace all tags by name (comma-separated, supports nested like `"parent/child"`) |
| `--add-tag NAME,...`   | Add tags without removing existing ones                                |
| `--remove-tag NAME,...`| Remove specific tags only                                              |
| `--title TEXT`         | Update title                                                           |
| `--description TEXT`   | Update description                                                     |

> Archive / unarchive are batch operations and live in their own commands — see [`cubox-cli archive`](#cubox-cli-archive) below.

### `cubox-cli archive`

Archive one or more cards by ID. Archived cards are hidden from the default `card list` (use `card list --archived` to view them).

```bash
cubox-cli archive --id CARD_ID[,ID2,...]
```

| Flag          | Description                                              |
| ------------- | -------------------------------------------------------- |
| `--id ID,...` | Card IDs to archive (comma-separated, required)          |

### `cubox-cli unarchive`

Restore one or more archived cards by moving them into a non-archived folder. The destination folder is required.

```bash
cubox-cli unarchive --id CARD_ID[,ID2,...] --folder NAME
```

| Flag             | Description                                                                |
| ---------------- | -------------------------------------------------------------------------- |
| `--id ID,...`    | Card IDs to unarchive (comma-separated, required)                          |
| `--folder NAME`  | Destination folder by name (required; `""` = Uncategorized; `"parent/child"` for nested) |

The folder name is resolved client-side via `folder list`; if the name does not match an existing folder the command fails with a helpful error.

### `cubox-cli delete`

Delete one or more cards by ID. Supports `--dry-run` to preview what would be deleted before committing.

```bash
cubox-cli delete --id CARD_ID [flags]
```


| Flag            | Description                                            |
| --------------- | ------------------------------------------------------ |
| `--id ID,...`   | Card IDs to delete (comma-separated, required)         |
| `--dry-run`     | Preview cards to be deleted without actually deleting   |


**Dry Run:** Always use `--dry-run` first to preview which cards will be deleted. For ≤ 3 cards the preview includes title and URL; for larger batches only the count is shown to avoid expensive per-card lookups.

```bash
# Preview
cubox-cli delete --id 7435692934957108160,7435691601617225646 --dry-run

# Execute after confirming
cubox-cli delete --id 7435692934957108160,7435691601617225646
```

### `cubox-cli annotation list`

List and search annotations across all cards.

```bash
cubox-cli annotation list [flags]
```


| Flag                | Description                                                  |
| ------------------- | ------------------------------------------------------------ |
| `--color COLOR,...` | Filter by color: `Yellow`, `Green`, `Blue`, `Pink`, `Purple` |
| `--keyword TEXT`    | Search annotations                                          |
| `--start-time TIME` | Filter by start time                                         |
| `--end-time TIME`   | Filter by end time                                           |
| `--limit N`         | Page size (default 50)                                       |
| `--last-id ID`      | Cursor pagination                                            |
| `--all`             | Auto-paginate to fetch all results                           |


## Examples

### Search for articles (keyword)

```bash
cubox-cli card list --keyword "machine learning" --page 1 -o pretty
```

### Semantic search (RAG)

```bash
cubox-cli card rag --query "articles about building REST APIs with authentication" -o pretty
```

### Browse cards in a specific folder

```bash
cubox-cli folder list -o text
cubox-cli card list --folder 7230156249357091393 --limit 10
```

### Read a saved article with AI insight

```bash
cubox-cli card detail --id 7247925101516031380 -o pretty
```

### Save a page and star it

```bash
cubox-cli save https://example.com --title "Example Site"
cubox-cli update --id CARD_ID --star --read
```

### Archive cards and list / restore them later

```bash
# Archive a batch of cards
cubox-cli archive --id 7444025677600260245,7443973659296793971

# Browse archived cards
cubox-cli card list --archived --limit 10

# Restore (move back) into a non-archived folder
cubox-cli unarchive --id 7444025677600260245,7443973659296793971 --folder "Reading List"
```

### Tidy up tags (rename, delete, merge)

```bash
# Rename a tag (leaf only — children follow automatically)
cubox-cli tag update --id 7295070793040398540 --new-name link

# Batch delete tags (cards keep their other tags)
cubox-cli tag delete --id 7444025677600260245,7443973659296793971

# Merge two tags into a target tag (cards re-tagged, sources deleted)
cubox-cli tag merge --source 7342187912403881105,7342187917722258501 --target 7247925099053977508
```

### Delete cards (with dry-run)

```bash
cubox-cli delete --id 7435692934957108160 --dry-run -o pretty
cubox-cli delete --id 7435692934957108160
```

### Export all annotations

```bash
cubox-cli annotation list --all -o pretty
```

### Cursor-based pagination (browsing)

```bash
cubox-cli card list --limit 5
# Use last card's ID for the next page
cubox-cli card list --limit 5 --last-id 7433152100604841820
```

### Search pagination

```bash
cubox-cli card list --keyword "AI" --page 1
cubox-cli card list --keyword "AI" --page 2
```

## Development

### Build from source

```bash
git clone https://github.com/OLCUBO/cubox-cli.git
cd cubox-cli
make build        # build for current platform
make build-all    # cross-compile for all platforms
make release      # create release archives
```

### Project structure

```
cubox-cli/
  main.go                 # entry point
  cmd/                    # cobra commands
    root.go               # root command, --output flag
    auth.go               # auth login/status/logout
    folder.go             # folder list
    tag.go                # tag list
    card.go               # card list, card detail
    save.go               # save web pages
    update.go             # update card
    archive.go            # batch archive / unarchive cards
    delete.go             # delete cards (with dry-run)
    annotation.go         # annotation list
    version.go            # version
  internal/
    client/               # HTTP client + API types
    config/               # config file management
  scripts/                # npm distribution wrapper
  skills/cubox/           # AI Agent skill
  .github/workflows/      # CI/CD
```

## License

[MIT](LICENSE)