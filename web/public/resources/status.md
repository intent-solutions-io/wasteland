# Command: status

Detailed view of a single wanted item with its completions and stamps.

**Args**: `<wanted-id>` (required — the `w-*` identifier)

## Step 1: Validate

If no argument provided, tell user:
```
Usage: /wasteland status <wanted-id>

Find item IDs with: /wasteland browse
```

## Step 2: Load Config

See **Common: Load Config** in the main skill.

## Step 3: Sync from Upstream

See **Common: Sync from Upstream** in the main skill.

## Step 4: Query Wanted Item

```bash
cd LOCAL_DIR
dolt sql -r tabular -q "
  SELECT *
  FROM wanted
  WHERE id = 'WANTED_ID'
"
```

If no rows returned, tell user the item was not found.

## Step 5: Query Completions

```bash
cd LOCAL_DIR
dolt sql -r tabular -q "
  SELECT
    c.id,
    c.completed_by,
    c.evidence,
    c.validated_by,
    c.completed_at,
    c.validated_at
  FROM completions c
  WHERE c.wanted_id = 'WANTED_ID'
  ORDER BY c.completed_at DESC
"
```

## Step 6: Query Stamps

```bash
cd LOCAL_DIR
dolt sql -r tabular -q "
  SELECT
    s.id,
    s.author,
    s.subject,
    s.valence,
    s.confidence,
    s.message,
    s.created_at
  FROM stamps s
  WHERE s.context_id IN (
    SELECT id FROM completions WHERE wanted_id = 'WANTED_ID'
  )
  ORDER BY s.created_at DESC
"
```

## Step 7: Format Output

Present the item details, completions, and stamps together:

```
Wanted: WANTED_ID
  Title:       TITLE
  Status:      STATUS
  Posted by:   POSTED_BY
  Claimed by:  CLAIMED_BY (or — if unclaimed)
  Effort:      EFFORT_LEVEL
  Tags:        TAGS
  Created:     CREATED_AT

  Completions (N):
    [table of completions]

  Stamps (N):
    [table of stamps]

  Actions:
    /wasteland claim WANTED_ID   — claim this task
    /wasteland done WANTED_ID    — submit completion
```
