package commons

import (
	"fmt"
	"strconv"
	"strings"
)

// StampDetail holds individual stamp data nested under a rig.
type StampDetail struct {
	Author      string   `json:"author"`
	Severity    string   `json:"severity"`
	Quality     float64  `json:"quality"`
	Reliability float64  `json:"reliability"`
	SkillTags   []string `json:"skill_tags,omitempty"`
	Message     string   `json:"message,omitempty"`
	CreatedAt   string   `json:"created_at"`
}

// CompletionDetail holds individual completion data nested under a rig.
type CompletionDetail struct {
	WantedID    string `json:"wanted_id"`
	WantedTitle string `json:"wanted_title,omitempty"`
	CompletedAt string `json:"completed_at"`
	ValidatedAt string `json:"validated_at,omitempty"`
}

// BadgeDetail holds individual badge data nested under a rig.
type BadgeDetail struct {
	BadgeType string `json:"badge_type"`
	AwardedAt string `json:"awarded_at"`
}

// ScoreboardDetailEntry extends ScoreboardEntry with nested detail data.
type ScoreboardDetailEntry struct {
	ScoreboardEntry
	RegisteredAt      string             `json:"registered_at,omitempty"`
	RigType           string             `json:"rig_type,omitempty"`
	RootStamps        int                `json:"root_stamps"`
	BranchStamps      int                `json:"branch_stamps"`
	LeafStamps        int                `json:"leaf_stamps"`
	Stamps            []StampDetail      `json:"stamps"`
	CompletionHistory []CompletionDetail `json:"completion_history"`
	Badges            []BadgeDetail      `json:"badges"`
}

// QueryScoreboardDetail returns per-rig detail with nested stamps, completions, and badges.
func QueryScoreboardDetail(db DB, limit int) ([]ScoreboardDetailEntry, error) {
	// Start with the base scoreboard data.
	base, err := QueryScoreboard(db, limit)
	if err != nil {
		return nil, fmt.Errorf("querying scoreboard base: %w", err)
	}
	if len(base) == 0 {
		return nil, nil
	}

	entries := make([]ScoreboardDetailEntry, len(base))
	for i, b := range base {
		entries[i].ScoreboardEntry = b
	}

	handles := make([]string, len(entries))
	for i, e := range entries {
		handles[i] = fmt.Sprintf("'%s'", EscapeSQL(e.RigHandle))
	}
	inClause := strings.Join(handles, ",")

	// Rig metadata (registered_at, rig_type).
	if err := populateDetailRigMeta(db, entries, inClause); err != nil {
		return nil, fmt.Errorf("querying rig metadata: %w", err)
	}

	// Per-severity stamp counts.
	if err := populateDetailSeverityCounts(db, entries, inClause); err != nil {
		return nil, fmt.Errorf("querying severity counts: %w", err)
	}

	// Individual stamps (capped at 50 per rig).
	if err := populateDetailStamps(db, entries, inClause); err != nil {
		return nil, fmt.Errorf("querying detail stamps: %w", err)
	}

	// Completion history (capped at 50 per rig).
	if err := populateDetailCompletions(db, entries, inClause); err != nil {
		return nil, fmt.Errorf("querying detail completions: %w", err)
	}

	// Badges.
	if err := populateDetailBadges(db, entries, inClause); err != nil {
		return nil, fmt.Errorf("querying detail badges: %w", err)
	}

	return entries, nil
}

func populateDetailRigMeta(db DB, entries []ScoreboardDetailEntry, inClause string) error {
	query := fmt.Sprintf(`SELECT handle, COALESCE(registered_at, '') AS registered_at, COALESCE(rig_type, '') AS rig_type
FROM rigs
WHERE handle IN (%s)`, inClause)

	output, err := db.Query(query, "")
	if err != nil {
		return err
	}

	rows := parseSimpleCSV(output)
	meta := make(map[string][2]string) // handle -> [registered_at, rig_type]
	for _, row := range rows {
		meta[row["handle"]] = [2]string{row["registered_at"], row["rig_type"]}
	}

	for i := range entries {
		if m, ok := meta[entries[i].RigHandle]; ok {
			entries[i].RegisteredAt = m[0]
			entries[i].RigType = m[1]
		}
	}
	return nil
}

func populateDetailSeverityCounts(db DB, entries []ScoreboardDetailEntry, inClause string) error {
	query := fmt.Sprintf(`SELECT subject, severity, COUNT(*) AS cnt
FROM stamps
WHERE subject IN (%s)
GROUP BY subject, severity`, inClause)

	output, err := db.Query(query, "")
	if err != nil {
		return err
	}

	rows := parseSimpleCSV(output)
	type sevCounts struct {
		root, branch, leaf int
	}
	perRig := make(map[string]*sevCounts)
	for _, row := range rows {
		rig := row["subject"]
		if perRig[rig] == nil {
			perRig[rig] = &sevCounts{}
		}
		cnt, err := strconv.Atoi(row["cnt"])
		if err != nil {
			return fmt.Errorf("parsing severity count for %q: %w", rig, err)
		}
		switch row["severity"] {
		case "root":
			perRig[rig].root = cnt
		case "branch":
			perRig[rig].branch = cnt
		case "leaf":
			perRig[rig].leaf = cnt
		}
	}

	for i := range entries {
		if sc, ok := perRig[entries[i].RigHandle]; ok {
			entries[i].RootStamps = sc.root
			entries[i].BranchStamps = sc.branch
			entries[i].LeafStamps = sc.leaf
		}
	}
	return nil
}

func populateDetailStamps(db DB, entries []ScoreboardDetailEntry, inClause string) error {
	// Use ROW_NUMBER to cap at 50 per rig.
	query := fmt.Sprintf(`SELECT subject, author, severity,
  COALESCE(JSON_EXTRACT(valence, '$.quality'), 0) AS quality,
  COALESCE(JSON_EXTRACT(valence, '$.reliability'), 0) AS reliability,
  COALESCE(skill_tags, '') AS skill_tags,
  COALESCE(message, '') AS message,
  COALESCE(created_at, '') AS created_at
FROM stamps
WHERE subject IN (%s)
ORDER BY subject, created_at DESC`, inClause)

	output, err := db.Query(query, "")
	if err != nil {
		return err
	}

	rows := parseSimpleCSV(output)

	perRig := make(map[string][]StampDetail)
	for _, row := range rows {
		rig := row["subject"]
		if len(perRig[rig]) >= 50 {
			continue
		}
		q, err := strconv.ParseFloat(row["quality"], 64)
		if err != nil {
			return fmt.Errorf("parsing quality for %q: %w", rig, err)
		}
		r, err := strconv.ParseFloat(row["reliability"], 64)
		if err != nil {
			return fmt.Errorf("parsing reliability for %q: %w", rig, err)
		}
		sd := StampDetail{
			Author:      row["author"],
			Severity:    row["severity"],
			Quality:     q,
			Reliability: r,
			Message:     row["message"],
			CreatedAt:   row["created_at"],
		}
		tags := parseTagsJSON(row["skill_tags"])
		if len(tags) > 0 {
			sd.SkillTags = tags
		}
		perRig[rig] = append(perRig[rig], sd)
	}

	for i := range entries {
		stamps := perRig[entries[i].RigHandle]
		if stamps == nil {
			stamps = []StampDetail{}
		}
		entries[i].Stamps = stamps
	}
	return nil
}

func populateDetailCompletions(db DB, entries []ScoreboardDetailEntry, inClause string) error {
	query := fmt.Sprintf(`SELECT c.completed_by, c.wanted_id,
  COALESCE(w.title, '') AS wanted_title,
  COALESCE(c.completed_at, '') AS completed_at,
  COALESCE(c.validated_at, '') AS validated_at
FROM completions c
LEFT JOIN wanted w ON c.wanted_id = w.id
WHERE c.completed_by IN (%s)
ORDER BY c.completed_by, c.completed_at DESC`, inClause)

	output, err := db.Query(query, "")
	if err != nil {
		return err
	}

	rows := parseSimpleCSV(output)

	perRig := make(map[string][]CompletionDetail)
	for _, row := range rows {
		rig := row["completed_by"]
		if len(perRig[rig]) >= 50 {
			continue
		}
		perRig[rig] = append(perRig[rig], CompletionDetail{
			WantedID:    row["wanted_id"],
			WantedTitle: row["wanted_title"],
			CompletedAt: row["completed_at"],
			ValidatedAt: row["validated_at"],
		})
	}

	for i := range entries {
		hist := perRig[entries[i].RigHandle]
		if hist == nil {
			hist = []CompletionDetail{}
		}
		entries[i].CompletionHistory = hist
	}
	return nil
}

func populateDetailBadges(db DB, entries []ScoreboardDetailEntry, inClause string) error {
	query := fmt.Sprintf(`SELECT rig_handle, badge_type, COALESCE(awarded_at, '') AS awarded_at
FROM badges
WHERE rig_handle IN (%s)
ORDER BY rig_handle, awarded_at DESC`, inClause)

	output, err := db.Query(query, "")
	if err != nil {
		return err
	}

	rows := parseSimpleCSV(output)

	perRig := make(map[string][]BadgeDetail)
	for _, row := range rows {
		perRig[row["rig_handle"]] = append(perRig[row["rig_handle"]], BadgeDetail{
			BadgeType: row["badge_type"],
			AwardedAt: row["awarded_at"],
		})
	}

	for i := range entries {
		badges := perRig[entries[i].RigHandle]
		if badges == nil {
			badges = []BadgeDetail{}
		}
		entries[i].Badges = badges
	}
	return nil
}
