# Command: sync

Pull upstream changes into the local fork and show a summary.

## Step 1: Load Config

See **Common: Load Config** in the main skill. If no config, tell user to run
`/wasteland join` first.

## Step 2: Pull Upstream

```bash
cd LOCAL_DIR
dolt pull upstream main
```

If this fails (merge conflict), note the error but continue with local
data — it may be slightly stale.

## Step 3: Show Summary

```bash
cd LOCAL_DIR
dolt sql -r tabular -q "
  SELECT
    status,
    COUNT(*) as count
  FROM wanted
  GROUP BY status
  ORDER BY
    CASE status WHEN 'open' THEN 0 WHEN 'claimed' THEN 1 WHEN 'in_review' THEN 2 ELSE 3 END
"
```

Print the counts and confirm sync completed:

```
Synced from upstream.

  open:       N
  claimed:    N
  in_review:  N

  Browse the board: /wasteland browse
```
