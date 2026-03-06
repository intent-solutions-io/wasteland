# Command: create

Create your own wasteland — a new DoltHub database from the MVR schema.

**Args**: `[owner/name]` (optional — will prompt if not provided)

Anyone can create a wasteland. You become its first rig and maintainer
(trust_level=3). Your wasteland is registered in the root commons
(`hop/wl-commons`) via PR, making it discoverable by the federation.

## Step 1: Check Prerequisites

```bash
dolt version
```

If dolt is not installed, tell the user:
- macOS: `brew install dolt`
- Linux: `curl -L https://github.com/dolthub/dolt/releases/latest/download/install.sh | bash`

```bash
dolt creds ls
```

If no credentials, tell user to run `dolt login` first.

## Step 2: Gather Details

If database path not provided in arguments, ask for:
- **Owner**: DoltHub org name (suggest their DoltHub username)
- **Database name**: Usually `wl-commons` (conventional name)

Then ask for:
- **Wasteland name**: Human-readable name (e.g., "Acme Engineering", "Indie Builders")
- **Description**: Optional description for DoltHub
- **Display name**: Your display name for the rigs table
- **Email**: Contact email

Also determine their DoltHub org from credentials:

```bash
dolt creds ls
```

## Step 3: Verify Database Doesn't Exist

```bash
curl -s "https://www.dolthub.com/api/v1alpha1/OWNER/DB_NAME" \
  -H "authorization: token $DOLTHUB_TOKEN" | head -5
```

If it exists, tell the user and suggest `/wasteland join OWNER/DB_NAME` instead.

## Step 4: Create Database on DoltHub

```bash
curl -s -X POST "https://www.dolthub.com/api/v1alpha1/database" \
  -H "Content-Type: application/json" \
  -H "authorization: token $DOLTHUB_TOKEN" \
  -d '{
    "ownerName": "OWNER",
    "repoName": "DB_NAME",
    "visibility": "public",
    "description": "Wasteland: WASTELAND_NAME — a HOP federation commons"
  }'
```

## Step 5: Initialize Schema from Template

Create a temp dolt database and apply the MVR schema (see `{baseDir}/resources/mvr-schema.md`):

```bash
TMPDIR=$(mktemp -d)
cd $TMPDIR
dolt init --name OWNER --email EMAIL

# Apply MVR schema via heredoc
dolt sql <<'SCHEMA'
-- (paste the full schema from {baseDir}/resources/mvr-schema.md)
SCHEMA
```

## Step 6: Configure Wasteland Metadata

```bash
cd $TMPDIR
dolt sql -q "REPLACE INTO _meta (\`key\`, value) VALUES ('wasteland_name', 'WASTELAND_NAME')"
dolt sql -q "REPLACE INTO _meta (\`key\`, value) VALUES ('created_by', 'HANDLE')"
dolt sql -q "REPLACE INTO _meta (\`key\`, value) VALUES ('upstream', 'hop/wl-commons')"
dolt sql -q "REPLACE INTO _meta (\`key\`, value) VALUES ('phase1_mode', 'wild_west')"
dolt sql -q "REPLACE INTO _meta (\`key\`, value) VALUES ('genesis_validators', '[\"HANDLE\"]')"

dolt add .
dolt commit -m "Initialize WASTELAND_NAME wasteland from MVR schema v1.1"
```

## Step 7: Register Creator as First Rig

```bash
cd $TMPDIR
dolt sql -q "INSERT INTO rigs (handle, display_name, dolthub_org, owner_email, gt_version, rig_type, trust_level, registered_at, last_seen) VALUES ('HANDLE', 'DISPLAY_NAME', 'OWNER', 'EMAIL', 'mvr-0.1', 'human', 3, NOW(), NOW())"
dolt add rigs
dolt commit -m "Register creator: HANDLE (maintainer)"
```

The creator gets trust_level=3 (maintainer) — they can validate completions,
merge PRs, and manage the wasteland.

## Step 8: Push to DoltHub

```bash
cd $TMPDIR
dolt remote add origin https://doltremoteapi.dolthub.com/OWNER/DB_NAME
dolt push origin main
```

## Step 9: Register in Root Commons

Register the new wasteland in the root commons (`hop/wl-commons`)
via the `chain_meta` table.

```bash
CHAIN_ID="wl-$(openssl rand -hex 8)"

ROOT_TMP=$(mktemp -d)
dolt clone hop/wl-commons $ROOT_TMP
cd $ROOT_TMP

dolt checkout -b "register-wasteland/OWNER/DB_NAME"

dolt sql -q "INSERT INTO chain_meta (chain_id, chain_type, parent_chain_id, hop_uri, dolt_database, created_at) VALUES ('$CHAIN_ID', 'community', NULL, 'hop://OWNER/DB_NAME', 'OWNER/DB_NAME', NOW())"
dolt add chain_meta
dolt commit -m "Register wasteland: WASTELAND_NAME (OWNER/DB_NAME)"

dolt push origin "register-wasteland/OWNER/DB_NAME"
```

Then open a DoltHub PR from the registration branch to main on
`hop/wl-commons`. If the user has a fork, push the branch there
and open the PR from the fork.

If root registration fails, it's non-fatal. The wasteland works without it —
it just won't be discoverable in the root directory yet.

## Step 10: Clean Up and Save Config

Update `~/.hop/config.json` to track the new wasteland.

If the config file exists, add the new wasteland to the `wastelands` array.
If it doesn't exist, create a new config:

```json
{
  "handle": "HANDLE",
  "wastelands": [
    {
      "upstream": "OWNER/DB_NAME",
      "fork": "OWNER/DB_NAME",
      "local_dir": "~/.hop/commons/OWNER/DB_NAME",
      "joined_at": "ISO_TIMESTAMP",
      "is_owner": true
    }
  ]
}
```

Clean up temp directories.

## Step 11: Confirm

```
Wasteland Created: WASTELAND_NAME

  Database:     OWNER/DB_NAME (DoltHub)
  Chain ID:     CHAIN_ID
  Creator:      HANDLE (maintainer, trust_level=3)
  Root:         registered (PR: URL) | not registered (standalone)

  Others can join with:
    /wasteland join OWNER/DB_NAME

  Your wasteland commands:
    /wasteland browse          — see the wanted board
    /wasteland post            — post work to your board
    /wasteland claim <id>      — claim a wanted item
    /wasteland done <id>       — submit completed work
```
