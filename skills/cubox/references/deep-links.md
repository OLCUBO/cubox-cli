# Cubox Deep Links

Construct clickable Cubox links from IDs returned by CLI commands. No extra API call is needed — just combine the ID with the correct URL pattern.

## URL Patterns

### HTTPS links (default)

Requires knowing the user's server (`cubox.pro` or `cubox.cc`). Obtain it from `cubox-cli auth status` or from the config at `~/.config/cubox-cli/config.json`.


| Resource | Pattern                             |
| -------- | ----------------------------------- |
| Card     | `https://{server}/web/card/{ID}`    |
| Folder   | `https://{server}/web/folders/{ID}` |
| Tag      | `https://{server}/web/tags/{ID}`    |


### cubox:// scheme links

Platform-native deep links. Use only when the user explicitly requests `cubox://` links.


| Resource | Pattern                  |
| -------- | ------------------------ |
| Card     | `cubox://card?id={ID}`   |
| Folder   | `cubox://folder?id={ID}` |
| Tag      | `cubox://tag?id={ID}`    |


## Rules

1. **Default to HTTPS links.** Only use `cubox://` scheme when the user explicitly asks for it.
2. **Server detection:** Run `cubox-cli auth status` once per session to determine the server. Cache the result — do not re-run for every link.
3. **Batch construction:** When listing multiple items, construct all links locally from the IDs. Do not make additional API calls.

## Examples

### List cards with Cubox links

```bash
cubox-cli card list --starred --limit 5
```

From the output, for each card with `"id": "7247925101516031380"` on server `cubox.pro`:

- HTTPS: `https://cubox.pro/web/card/7247925101516031380`
- cubox://: `cubox://card?id=7247925101516031380`

### List folders with links

```bash
cubox-cli folder list
```

For a folder with `"id": "7230156249357091393"` on server `cubox.cc`:

- HTTPS: `https://cubox.cc/web/folders/7230156249357091393`
- cubox://: `cubox://folder?id=7230156249357091393`

### List tags with links

```bash
cubox-cli tag list
```

For a tag with `"id": "7230156249461949232"` on server `cubox.pro`:

- HTTPS: `https://cubox.pro/web/tags/7230156249461949232`
- cubox://: `cubox://tag?id=7230156249461949232`

