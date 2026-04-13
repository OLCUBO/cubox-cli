# card delete — Dry Run Policy

Delete one or more cards by ID.

```bash
cubox-cli delete --id CARD_ID [--id ID2,...] [--dry-run]
```

## Flags


| Flag          | Description                                                    |
| ------------- | -------------------------------------------------------------- |
| `--id ID,...` | Card IDs to delete (comma-separated, required)                 |
| `--dry-run`   | Preview which cards would be deleted without actually deleting |


## Output

Returns: `{ "dry_run": bool, "count": N, "cards": [...], "message": "..." }`

- When ≤ 3 cards: `cards` array includes `id`, `title`, `url` fetched from server.
- When > 3 cards: `cards` is omitted; only `count` is shown (avoids heavy per-card API calls).

## Dry Run Policy for AI Agents

Deleted items will stay in your "Recently Deleted" folder for 30 days before permanent removal.

1. **Always** run with `--dry-run` first before any real deletion.
2. Present the dry-run result to the user and ask for explicit confirmation.
  - ≤ 3 cards: show the card titles returned by dry-run.
  - 3 cards: tell the user the count (e.g. "Will delete 25 cards").
3. Only after the user confirms, run again without `--dry-run` to perform the actual deletion.
4. Never skip the dry-run step — deletions are irreversible.

## Example workflow

```bash
# Step 1: preview
cubox-cli delete --id 7435692934957108160,7435691601617225646 --dry-run
# Step 2: show user the preview, ask "Delete these 2 cards?"
# Step 3: if confirmed, execute
cubox-cli delete --id 7435692934957108160,7435691601617225646
```

