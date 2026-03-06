# Command: join

Join a wasteland — register as a rig in the HOP federation.

**Args**: `[upstream]` (default: `hop/wl-commons`)

You can join any wasteland by specifying its DoltHub path:
- `/wasteland join` — join the root commons (hop/wl-commons)
- `/wasteland join grab/wl-commons` — join Grab's wasteland
- `/wasteland join alice-dev/wl-commons` — join Alice's wasteland

Your rig can participate in multiple wastelands simultaneously.

## Step 1: Check Prerequisites

```bash
dolt version
```

If dolt is not installed, tell the user:
- macOS: `brew install dolt`
- Linux: `curl -L https://github.com/dolthub/dolt/releases/latest/download/install.sh | bash`
- Or see https://docs.dolthub.com/introduction/installation

```bash
dolt creds ls
```

If no credentials, tell user to run `dolt login` first.

## Step 2: Gather Identity

Check if `~/.hop/config.json` already exists:

```bash
cat ~/.hop/config.json 2>/dev/null
```

If it exists and has a handle, the user is already registered. Show their
config and check if they're already in the target wasteland:
- If already joined this wasteland: tell user and offer to re-sync
- If not yet joined: proceed to add this wasteland (keep existing identity)

If it doesn't exist, ask the user for:
- **Handle**: Their rig name (suggest their DoltHub username or GitHub username)
- **Display name**: Human-readable name (suggest: "Alice's Workshop" style)
- **Type**: human, agent, or org (default: human)
- **Email**: Contact email (for the rigs table)

Also determine their DoltHub org:

```bash
dolt creds ls
```

## Step 3: Create MVR Home

```bash
mkdir -p ~/.hop/commons
```

## Step 4: Fork the Commons

Parse upstream into org and db name (split on `/`).

Fork the upstream commons to the user's DoltHub org via the DoltHub API:

```bash
curl -s -X POST "https://www.dolthub.com/api/v1alpha1/database/fork" \
  -H "Content-Type: application/json" \
  -H "authorization: token $DOLTHUB_TOKEN" \
  -d '{
    "owner_name": "USER_DOLTHUB_ORG",
    "new_repo_name": "UPSTREAM_DB",
    "from_owner": "UPSTREAM_ORG",
    "from_repo_name": "UPSTREAM_DB"
  }'
```

If the fork already exists (error contains "already exists"), that's fine.

The DOLTHUB_TOKEN can come from environment variable DOLTHUB_TOKEN, or
extract it from the dolt credentials:

```bash
dolt creds ls
```

If you can't find a token, ask the user to set DOLTHUB_TOKEN or get one
from https://www.dolthub.com/settings/tokens

## Step 5: Clone the Fork

```bash
dolt clone "USER_DOLTHUB_ORG/UPSTREAM_DB" ~/.hop/commons/UPSTREAM_ORG/UPSTREAM_DB
```

If already cloned (`.dolt` directory exists), skip.

## Step 6: Add Upstream Remote

```bash
cd ~/.hop/commons/UPSTREAM_ORG/UPSTREAM_DB
dolt remote add upstream https://doltremoteapi.dolthub.com/UPSTREAM_ORG/UPSTREAM_DB
```

If upstream already exists, that's fine.

## Step 7: Register as a Rig

```bash
cd ~/.hop/commons/UPSTREAM_ORG/UPSTREAM_DB
dolt sql -q "INSERT INTO rigs (handle, display_name, dolthub_org, owner_email, gt_version, trust_level, registered_at, last_seen) VALUES ('HANDLE', 'DISPLAY_NAME', 'DOLTHUB_ORG', 'EMAIL', 'mvr-0.1', 1, NOW(), NOW()) ON DUPLICATE KEY UPDATE last_seen = NOW(), gt_version = 'mvr-0.1'"
dolt add .
dolt commit -m "Register rig: HANDLE"
```

## Step 8: Push Registration

```bash
cd ~/.hop/commons/UPSTREAM_ORG/UPSTREAM_DB
dolt push origin main
```

## Step 9: Save Config

If `~/.hop/config.json` already exists (joining additional wasteland),
read the existing config, append the new wasteland to the `wastelands`
array, and write back. Do NOT overwrite identity fields (handle, type, etc.).

If creating a new config, write `~/.hop/config.json`:

```json
{
  "handle": "USER_HANDLE",
  "display_name": "USER_DISPLAY_NAME",
  "type": "human",
  "dolthub_org": "DOLTHUB_ORG",
  "email": "USER_EMAIL",
  "wastelands": [
    {
      "upstream": "UPSTREAM_ORG/UPSTREAM_DB",
      "fork": "DOLTHUB_ORG/UPSTREAM_DB",
      "local_dir": "~/.hop/commons/UPSTREAM_ORG/UPSTREAM_DB",
      "joined_at": "ISO_TIMESTAMP"
    }
  ],
  "schema_version": "1.0",
  "mvr_version": "0.1"
}
```

When appending, add a new entry to the `wastelands` array:

```json
{
  "upstream": "UPSTREAM_ORG/UPSTREAM_DB",
  "fork": "DOLTHUB_ORG/UPSTREAM_DB",
  "local_dir": "~/.hop/commons/UPSTREAM_ORG/UPSTREAM_DB",
  "joined_at": "ISO_TIMESTAMP"
}
```

## Step 10: Confirm

Print a summary:

```
MVR Node Registered

  Handle:     USER_HANDLE
  Type:       human
  DoltHub:    DOLTHUB_ORG/UPSTREAM_DB
  Upstream:   UPSTREAM_ORG/UPSTREAM_DB
  Local:      ~/.hop/commons/UPSTREAM_ORG/UPSTREAM_DB

  You are now a rig in the Wasteland.

  Next steps:
    /wasteland browse   — see the wanted board
    /wasteland claim    — claim a task
    /wasteland done     — submit completed work
```
