package api

import (
	"net/http"
	"time"

	"github.com/julianknutsen/wasteland/internal/commons"
)

// ScoreboardDumpResponse is the JSON response for GET /api/scoreboard/dump.
type ScoreboardDumpResponse struct {
	Rigs        []commons.RigRow        `json:"rigs"`
	Stamps      []commons.StampRow      `json:"stamps"`
	Completions []commons.CompletionRow `json:"completions"`
	Wanted      []commons.WantedRow     `json:"wanted"`
	Badges      []commons.BadgeRow      `json:"badges"`
	UpdatedAt   string                  `json:"updated_at"`
}

func toScoreboardDumpResponse(dump *commons.ScoreboardDump) *ScoreboardDumpResponse {
	return &ScoreboardDumpResponse{
		Rigs:        dump.Rigs,
		Stamps:      dump.Stamps,
		Completions: dump.Completions,
		Wanted:      dump.Wanted,
		Badges:      dump.Badges,
		UpdatedAt:   time.Now().UTC().Format(time.RFC3339),
	}
}

// handleScoreboardDump serves the cached scoreboard dump JSON with CORS headers.
func (s *Server) handleScoreboardDump(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	if s.scoreboardDump == nil {
		writeError(w, http.StatusServiceUnavailable, "scoreboard dump not configured")
		return
	}

	data := s.scoreboardDump.Get()
	if data == nil {
		writeError(w, http.StatusServiceUnavailable, "scoreboard dump data unavailable")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "public, max-age=300")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(data)
}
