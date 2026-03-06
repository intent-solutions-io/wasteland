# Command: browse

Browse the wanted board — see available work.

**Args**: `[filter]` (optional — filter by status, tag, or keyword)

## Step 1: Load Config

See **Common: Load Config** in the main skill. If no config, tell user to run
`/wasteland join` first.

## Step 2: Sync from Upstream

See **Common: Sync from Upstream** in the main skill.

## Step 3: Query the Wanted Board

```bash
cd LOCAL_DIR
dolt sql -r tabular -q "
  SELECT
    id,
    title,
    COALESCE(status, 'open') as status,
    COALESCE(effort_level, 'medium') as effort,
    COALESCE(posted_by, '—') as posted_by,
    COALESCE(claimed_by, '—') as claimed_by,
    COALESCE(JSON_EXTRACT(tags, '$'), '[]') as tags
  FROM wanted
  ORDER BY
    CASE status WHEN 'open' THEN 0 WHEN 'claimed' THEN 1 ELSE 2 END,
    priority ASC,
    created_at DESC
"
```

## Step 4: Format Output

Present results as a clean table. Group by status:

**Open** — available to claim
**Claimed** — someone is working on it
**In Review** — completed, awaiting validation

If a filter argument was provided:
- If it matches a status (open/claimed/in_review), filter by status
- Otherwise, search title, tags, and project fields for the keyword

## Step 5: Show Rig Registry (optional)

If the user asks or if the board is empty, also show registered rigs:

```bash
cd LOCAL_DIR
dolt sql -r tabular -q "
  SELECT handle, display_name, trust_level, registered_at
  FROM rigs
  ORDER BY registered_at DESC
  LIMIT 20
"
```

## Step 6: Show Character Sheet (optional)

If the user asks about their own profile:

```bash
cd LOCAL_DIR
dolt sql -r tabular -q "
  SELECT
    c.id,
    c.wanted_id,
    w.title as task,
    c.completed_at
  FROM completions c
  LEFT JOIN wanted w ON c.wanted_id = w.id
  WHERE c.completed_by = 'USER_HANDLE'
  ORDER BY c.completed_at DESC
"
```

And their stamps:

```bash
cd LOCAL_DIR
dolt sql -r tabular -q "
  SELECT
    s.id,
    s.author,
    s.valence,
    s.confidence,
    s.severity,
    s.created_at
  FROM stamps s
  WHERE s.context_id IN (
    SELECT id FROM completions WHERE completed_by = 'USER_HANDLE'
  )
  ORDER BY s.created_at DESC
"
```
