# Command: post

Post a wanted item to the board.

**Args**: `[title]` (optional — will prompt if not provided)

## Step 1: Load Config

See **Common: Load Config** in the main skill. If no config, tell user to run
`/wasteland join` first.

## Step 2: Gather Details

If title not provided in arguments, ask for it.

Then ask for:
- **Description**: What needs to be done (can be multi-line)
- **Project**: Project name (optional, e.g., "gastown", "beads", "hop")
- **Type**: bug, feature, docs, design, research, community (default: feature)
- **Effort level**: trivial, small, medium, large, epic (default: medium)
- **Tags**: Comma-separated tags (e.g., "Go,testing,refactor")
- **Sandbox required?**: true/false (default: false)

## Step 3: Generate Wanted ID

```bash
echo "w-$(openssl rand -hex 5)"
```

## Step 4: Insert

```bash
cd LOCAL_DIR
dolt sql -q "INSERT INTO wanted (id, title, description, project, type, priority, tags, posted_by, status, effort_level, sandbox_required, created_at, updated_at) VALUES ('WANTED_ID', 'TITLE', 'DESCRIPTION', PROJECT_OR_NULL, 'TYPE', 2, TAGS_JSON_OR_NULL, 'USER_HANDLE', 'open', 'EFFORT', SANDBOX_BOOL, NOW(), NOW())"
dolt add .
dolt commit -m "Post wanted: TITLE"
dolt push origin main
```

For tags, format as JSON array: `'["Go","testing"]'` or NULL if none.

## Step 5: Confirm

```
Posted: WANTED_ID
  Title:  TITLE
  By:     USER_HANDLE
  Effort: EFFORT_LEVEL
  Tags:   TAG_LIST

  Submit directly:        /wasteland done WANTED_ID
  Or claim it first:      /wasteland claim WANTED_ID
```
