# cubox-cli

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Go Version](https://img.shields.io/badge/go-%3E%3D1.23-blue.svg)](https://go.dev/)
[![npm version](https://img.shields.io/npm/v/cubox-cli.svg)](https://www.npmjs.com/package/cubox-cli)

[English](./README.md) | [中文版](./README.zh.md)

The official [Cubox](https://cubox.pro) CLI tool — built for humans and AI Agents. Manage your bookmarks, browse collections, and read saved content from the terminal.

[Install](#installation) · [Auth](#authentication) · [Commands](#commands) · [AI Agent](#quick-start-ai-agent) · [Examples](#examples) · [Development](#development)

## Features

| Category | Capabilities |
|----------|-------------|
| Groups | List and browse bookmark folders |
| Tags | List and browse tag hierarchy |
| Cards | Filter cards by group, type, tag, starred/read/annotated status; cursor-based pagination |
| Content | Read full article content in markdown format |

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
   - China: https://cubox.pro/web/settings/extensions
   - International: https://cubox.cc/web/settings/extensions
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
cubox-cli card content --id CARD_ID
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

| Command | Description |
|---------|-------------|
| `cubox-cli auth login` | Interactive login (server selection + token input) |
| `cubox-cli auth login --server cubox.pro --token TOKEN` | Non-interactive login (for agents) |
| `cubox-cli auth status` | Show current server, masked token, connection test |
| `cubox-cli auth logout` | Remove saved credentials |

Credentials are stored at `~/.config/cubox-cli/config.json`.

## Commands

### Output Formats

All commands support the `-o` / `--output` flag:

| Flag | Description |
|------|-------------|
| `-o json` | Compact JSON (default, agent-friendly) |
| `-o pretty` | Indented JSON |
| `-o text` | Human-readable text/tree output |

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

Filter and list bookmark cards.

```bash
cubox-cli card list [flags]
```

| Flag | Description |
|------|-------------|
| `--group ID,...` | Filter by group/folder IDs |
| `--type TYPE,...` | Filter by type: `Article`, `Snippet`, `Memo`, `Image`, `Audio`, `Video`, `File` |
| `--tag ID,...` | Filter by tag IDs |
| `--starred` | Only starred cards |
| `--read` | Only read cards |
| `--unread` | Only unread cards |
| `--annotated` | Only cards with highlights |
| `--limit N` | Page size (default 50) |
| `--cursor ID,TIME` | Resume from last card for pagination |
| `--all` | Auto-paginate to fetch all results |

**JSON output fields:** `id`, `title`, `description`, `article_title`, `domain`, `type`, `tags`, `url`, `cubox_url`, `words_count`, `create_time`, `update_time`, `highlights`

### `cubox-cli card content --id ID`

Get the full article content in markdown format.

```bash
cubox-cli card content --id 7247925101516031380
```

By default outputs raw markdown. Use `-o pretty` for a JSON envelope.

## Examples

### List all starred articles

```bash
cubox-cli card list --starred --type Article -o pretty
```

### Browse cards in a specific folder

```bash
# Find the folder ID
cubox-cli group list -o text

# List cards in that folder
cubox-cli card list --group 7230156249357091393 --limit 10
```

### Read a saved article

```bash
cubox-cli card content --id 7247925101516031380
```

### Fetch all annotated cards with highlights

```bash
cubox-cli card list --annotated --all -o pretty
```

### Pagination

```bash
# First page
cubox-cli card list --limit 5

# Use the last card's ID and update_time for the next page
cubox-cli card list --limit 5 --cursor "7247925102807877551,2024-12-04T16:23:01:347+08:00"
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
    card.go               # card list, card content
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
