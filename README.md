# cubox-cli

[License: MIT](https://opensource.org/licenses/MIT)
[Go Version](https://go.dev/)
[npm version](https://www.npmjs.com/package/cubox-cli)

[English](./README.md) | [中文版](./README.zh.md)

The official [Cubox](https://cubox.pro) CLI tool — built for humans and AI Agents. Manage your bookmarks, browse collections, and read saved content from the terminal.

[Install](#installation) · [Auth](#authentication) · [Commands](#commands) · [AI Agent](#quick-start-ai-agent) · [Examples](#examples) · [Development](#development)

## Features


| Category | Capabilities                                                                          |
| -------- | ------------------------------------------------------------------------------------- |
| Groups   | List and browse bookmark folders                                                      |
| Tags     | List and browse tag hierarchy                                                         |
| Cards    | Filter/search cards by group, tag, starred/read/annotated status, keyword, time range |
| Content  | Read full card detail with article content (markdown), highlights, and AI insight     |
| Save     | Save web page URLs as bookmarks                                                       |
| Update   | Star/unstar, mark read/unread, archive, move to group, add tags                       |
| Delete   | Delete bookmark cards by ID, with dry-run preview support                             |
| Marks    | List and search highlights/annotations across all cards                               |


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

1. **Select your server** — `cubox.pro` (China) or `cubox.cc` (international)
2. **Get your API key** — the CLI shows the URL to open:
  - China: [https://cubox.pro/web/settings/extensions](https://cubox.pro/web/settings/extensions)
  - International: [https://cubox.cc/web/settings/extensions](https://cubox.cc/web/settings/extensions)
3. **Paste the API link** — copy the full link (e.g. `https://cubox.pro/c/api/save/abcdefg`) and paste it. The CLI extracts the token automatically.

#### Start using

```bash
# List your folders
cubox-cli group list

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

Ask the user which server they use (`cubox.pro` or `cubox.cc`), then instruct them to open the extensions settings page and copy their API link:

- China: `https://cubox.pro/web/settings/extensions`
- International: `https://cubox.cc/web/settings/extensions`

Once the user provides the API link or token, run:

```bash
cubox-cli auth login --server cubox.pro --token TOKEN
```

The `--token` flag accepts either the full API link URL or just the token string.

**Step 3 — Verify**

```bash
cubox-cli auth status
```

**Step 4 — Use**

```bash
cubox-cli group list
cubox-cli card list --limit 10
```

## Authentication


| Command                                                 | Description                                        |
| ------------------------------------------------------- | -------------------------------------------------- |
| `cubox-cli auth login`                                  | Interactive login (server selection + token input) |
| `cubox-cli auth login --server cubox.pro --token TOKEN` | Non-interactive login (for agents)                 |
| `cubox-cli auth status`                                 | Show current server, masked token, connection test |
| `cubox-cli auth logout`                                 | Remove saved credentials                           |


Credentials are stored at `~/.config/cubox-cli/config.json`.

## Commands

### Output Formats

All commands support the `-o` / `--output` flag:


| Flag        | Description                            |
| ----------- | -------------------------------------- |
| `-o json`   | Compact JSON (default, agent-friendly) |
| `-o pretty` | Indented JSON                          |
| `-o text`   | Human-readable text/tree output        |


### `cubox-cli group list`

List all bookmark groups (folders).

```bash
cubox-cli group list
cubox-cli group list -o text
```

**JSON output fields:** `id`, `nested_name`, `name`, `parent_id`, `uncategorized`

### `cubox-cli tag list`

List all tags.

```bash
cubox-cli tag list
cubox-cli tag list -o text
```

**JSON output fields:** `id`, `nested_name`, `name`, `parent_id`

### `cubox-cli card list`

Filter and search bookmark cards. Supports keyword search with page-based pagination, and cursor-based pagination for browsing.

```bash
cubox-cli card list [flags]
```


| Flag                | Description                                                        |
| ------------------- | ------------------------------------------------------------------ |
| `--group ID,...`    | Filter by group/folder IDs                                         |
| `--tag ID,...`      | Filter by tag IDs                                                  |
| `--starred`         | Only starred cards                                                 |
| `--read`            | Only read cards                                                    |
| `--unread`          | Only unread cards                                                  |
| `--annotated`       | Only cards with highlights                                         |
| `--keyword TEXT`    | Search by keyword                                                  |
| `--start-time TIME` | Filter by create time start (e.g. `2026-01-01T00:00:00:000+08:00`) |
| `--end-time TIME`   | Filter by create time end                                          |
| `--limit N`         | Page size (default 50)                                             |
| `--last-id CARD_ID` | Cursor pagination for browsing (non-search)                        |
| `--page N`          | Page number for search (1-based, used with `--keyword`)            |
| `--all`             | Auto-paginate to fetch all results                                 |


**Pagination:** When `--keyword` is set, use `--page` for pagination. Otherwise, use `--last-id` with the last card's ID.

### `cubox-cli card detail --id ID`

Get full card detail including article content (markdown), author, highlights, and AI-generated insight (summary + Q&A).

```bash
cubox-cli card detail --id 7247925101516031380
cubox-cli card detail --id 7247925101516031380 -o pretty
```

Use `-o text` to output only the markdown content.

### `cubox-cli save`

Save one or more web page URLs as bookmarks.

```bash
cubox-cli save https://example.com
cubox-cli save https://a.com https://b.com --group GROUP_ID
cubox-cli save https://example.com --tag TAG_ID1,TAG_ID2
```

### `cubox-cli update`

Update a card's properties.

```bash
cubox-cli update --id CARD_ID [flags]
```


| Flag                   | Description            |
| ---------------------- | ---------------------- |
| `--star` / `--unstar`  | Toggle star            |
| `--read` / `--unread`  | Toggle read status     |
| `--archive`            | Archive the card       |
| `--group GROUP_ID`     | Move to a group/folder |
| `--add-tag TAG_ID,...` | Add tags               |


### `cubox-cli delete`

Delete one or more bookmark cards by ID. Supports `--dry-run` to preview what would be deleted before committing.

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

### `cubox-cli mark list`

List and search highlights (annotations) across all cards.

```bash
cubox-cli mark list [flags]
```


| Flag                | Description                                                  |
| ------------------- | ------------------------------------------------------------ |
| `--color COLOR,...` | Filter by color: `Yellow`, `Green`, `Blue`, `Pink`, `Purple` |
| `--keyword TEXT`    | Search highlights                                            |
| `--start-time TIME` | Filter by start time                                         |
| `--end-time TIME`   | Filter by end time                                           |
| `--limit N`         | Page size (default 50)                                       |
| `--last-id ID`      | Cursor pagination                                            |
| `--all`             | Auto-paginate to fetch all results                           |


## Examples

### Search for articles

```bash
cubox-cli card list --keyword "machine learning" --page 1 -o pretty
```

### Browse cards in a specific folder

```bash
cubox-cli group list -o text
cubox-cli card list --group 7230156249357091393 --limit 10
```

### Read a saved article with AI insight

```bash
cubox-cli card detail --id 7247925101516031380 -o pretty
```

### Save a URL and star it

```bash
cubox-cli save https://example.com
cubox-cli update --id CARD_ID --star --read
```

### Delete cards (with dry-run)

```bash
cubox-cli delete --id 7435692934957108160 --dry-run -o pretty
cubox-cli delete --id 7435692934957108160
```

### Export all highlights

```bash
cubox-cli mark list --all -o pretty
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
    group.go              # group list
    tag.go                # tag list
    card.go               # card list, card detail
    save.go               # save URLs
    update.go             # update card
    delete.go             # delete cards (with dry-run)
    mark.go               # mark (highlight) list
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