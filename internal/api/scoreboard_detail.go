package api

import (
	"net/http"
	"time"

	"github.com/julianknutsen/wasteland/internal/commons"
)

// ScoreboardDetailResponse is the JSON response for GET /api/scoreboard/detail.
type ScoreboardDetailResponse struct {
	Entries   []ScoreboardDetailEntryJSON `json:"entries"`
	UpdatedAt string                      `json:"updated_at"`
}

// ScoreboardDetailEntryJSON is the JSON representation of a per-rig detail entry.
type ScoreboardDetailEntryJSON struct {
	RigHandle      string                     `json:"rig_handle"`
	DisplayName    string                     `json:"display_name,omitempty"`
	TrustTier      string                     `json:"trust_tier"`
	StampCount     int                        `json:"stamp_count"`
	WeightedScore  int                        `json:"weighted_score"`
	UniqueTowns    int                        `json:"unique_towns"`
	Completions    int                        `json:"completions"`
	AvgQuality     float64                    `json:"avg_quality"`
	AvgReliability float64                    `json:"avg_reliability"`
	TopSkills      []string                   `json:"top_skills,omitempty"`
	RegisteredAt   string                     `json:"registered_at,omitempty"`
	RigType        string                     `json:"rig_type,omitempty"`
	RootStamps     int                        `json:"root_stamps"`
	BranchStamps   int                        `json:"branch_stamps"`
	LeafStamps     int                        `json:"leaf_stamps"`
	Stamps         []commons.StampDetail      `json:"stamps"`
	CompHistory    []commons.CompletionDetail `json:"completion_history"`
	Badges         []commons.BadgeDetail      `json:"badges"`
}

func toScoreboardDetailResponse(entries []commons.ScoreboardDetailEntry) *ScoreboardDetailResponse {
	items := make([]ScoreboardDetailEntryJSON, len(entries))
	for i, e := range entries {
		items[i] = ScoreboardDetailEntryJSON{
			RigHandle:      e.RigHandle,
			DisplayName:    e.DisplayName,
			TrustTier:      e.TrustTier,
			StampCount:     e.StampCount,
			WeightedScore:  e.WeightedScore,
			UniqueTowns:    e.UniqueTowns,
			Completions:    e.Completions,
			AvgQuality:     e.AvgQuality,
			AvgReliability: e.AvgReliab,
			TopSkills:      e.TopSkills,
			RegisteredAt:   e.RegisteredAt,
			RigType:        e.RigType,
			RootStamps:     e.RootStamps,
			BranchStamps:   e.BranchStamps,
			LeafStamps:     e.LeafStamps,
			Stamps:         e.Stamps,
			CompHistory:    e.CompletionHistory,
			Badges:         e.Badges,
		}
	}
	return &ScoreboardDetailResponse{
		Entries:   items,
		UpdatedAt: time.Now().UTC().Format(time.RFC3339),
	}
}

// handleScoreboardDetail serves the cached scoreboard detail JSON with CORS headers.
func (s *Server) handleScoreboardDetail(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	if s.scoreboardDetail == nil {
		writeError(w, http.StatusServiceUnavailable, "scoreboard detail not configured")
		return
	}

	data := s.scoreboardDetail.Get()
	if data == nil {
		writeError(w, http.StatusServiceUnavailable, "scoreboard detail data unavailable")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "public, max-age=300")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(data)
}
