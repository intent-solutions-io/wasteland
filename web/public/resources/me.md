# Command: me

Personal dashboard — shows your claimed tasks, completions, stamps, and
badges in one view.

## Step 1: Load Config

See **Common: Load Config** in the main skill. Extract USER_HANDLE from `handle`.

## Step 2: Sync from Upstream

See **Common: Sync from Upstream** in the main skill.

## Step 3: Show Claimed Items

```bash
cd LOCAL_DIR
dolt sql -r tabular -q "
  SELECT id, title, status, effort_level, updated_at
  FROM wanted
  WHERE claimed_by = 'USER_HANDLE' AND status IN ('claimed', 'in_review')
  ORDER BY updated_at DESC
"
```

## Step 4: Show Completions

```bash
cd LOCAL_DIR
dolt sql -r tabular -q "
  SELECT
    c.id,
    c.wanted_id,
    w.title as task,
    c.evidence,
    c.completed_at
  FROM completions c
  LEFT JOIN wanted w ON c.wanted_id = w.id
  WHERE c.completed_by = 'USER_HANDLE'
  ORDER BY c.completed_at DESC
"
```

## Step 5: Show Stamps Received

```bash
cd LOCAL_DIR
dolt sql -r tabular -q "
  SELECT
    s.id,
    s.author,
    s.valence,
    s.confidence,
    s.severity,
    s.message,
    s.created_at
  FROM stamps s
  WHERE s.subject = 'USER_HANDLE'
  ORDER BY s.created_at DESC
"
```

## Step 6: Show Badges

```bash
cd LOCAL_DIR
dolt sql -r tabular -q "
  SELECT badge_type, evidence, awarded_at
  FROM badges
  WHERE rig_handle = 'USER_HANDLE'
  ORDER BY awarded_at DESC
"
```

## Step 7: Format Dashboard

Present the results as a personal dashboard summary:

```
Personal Dashboard: USER_HANDLE

  Active Claims (N):
    [table of claimed items]

  Completions (N):
    [table of completions]

  Stamps Received (N):
    [table of stamps]

  Badges (N):
    [table of badges]

  Next steps:
    /wasteland browse   — find more work
    /wasteland done <id> — submit a completion
```
