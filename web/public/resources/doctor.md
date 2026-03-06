# Command: doctor

Verify the wasteland setup is functional — checks prerequisites,
configuration, and connectivity.

## Step 1: Check Dolt

```bash
dolt version
```

If dolt is not installed, report FAIL and tell user:
- macOS: `brew install dolt`
- Linux: `curl -L https://github.com/dolthub/dolt/releases/latest/download/install.sh | bash`

## Step 2: Check Config

```bash
cat ~/.hop/config.json
```

If missing, report FAIL and tell user to run `/wasteland join` first.
If present, verify it contains `handle` and at least one entry in
`wastelands[]`.

## Step 3: Check Local Clone

Verify LOCAL_DIR exists and contains a `.dolt` directory:

```bash
ls -d LOCAL_DIR/.dolt
```

If missing, report FAIL — the local clone may need to be re-created
via `/wasteland join`.

## Step 4: Check Remotes

```bash
cd LOCAL_DIR
dolt remote -v
```

Verify both `origin` (user's fork) and `upstream` (the commons source)
are configured. Report FAIL for any missing remote.

## Step 5: Check Connectivity

```bash
cd LOCAL_DIR
dolt fetch upstream 2>&1
```

If fetch succeeds, connectivity is good. If it fails, report FAIL with
the error message.

## Step 6: Print Summary

```
Wasteland Doctor

  [PASS] dolt installed (vX.Y.Z)
  [PASS] config exists (~/.hop/config.json)
  [PASS] local clone (LOCAL_DIR)
  [PASS] remotes configured (origin + upstream)
  [PASS] connectivity (upstream reachable)

  All checks passed. Your wasteland setup is healthy.
```

Or for failures:

```
  [FAIL] config exists — run /wasteland join first
```
