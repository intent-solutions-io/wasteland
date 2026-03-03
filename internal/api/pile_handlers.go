package api

import (
	"errors"
	"net/http"

	"github.com/julianknutsen/wasteland/internal/pile"
)

const maxProfileSearchLimit = 100

// handleProfile serves GET /api/profile/{handle}
// Returns a full developer profile from hop/the-pile.
// No auth required — profile lookups are public read-only data.
func (s *Server) handleProfile(w http.ResponseWriter, r *http.Request) {
	handle := r.PathValue("handle")
	if handle == "" {
		writeError(w, http.StatusBadRequest, "handle is required")
		return
	}

	if s.pile == nil {
		writeError(w, http.StatusServiceUnavailable, "profile service not configured")
		return
	}

	profile, err := pile.QueryProfile(s.pile, handle)
	if err != nil {
		if errors.Is(err, pile.ErrProfileNotFound) {
			writeError(w, http.StatusNotFound, err.Error())
		} else {
			writeError(w, http.StatusBadGateway, "upstream profile service error")
		}
		return
	}
	writeJSON(w, http.StatusOK, profile)
}

// handleProfileSearch serves GET /api/profile?q=search
// Searches for profiles matching the query string.
func (s *Server) handleProfileSearch(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query().Get("q")
	if q == "" {
		writeError(w, http.StatusBadRequest, "q parameter is required")
		return
	}

	if s.pile == nil {
		writeError(w, http.StatusServiceUnavailable, "profile service not configured")
		return
	}

	limit := parseIntParam(r, "limit", 20)
	if limit > maxProfileSearchLimit {
		limit = maxProfileSearchLimit
	}
	results, err := pile.SearchProfiles(s.pile, q, limit)
	if err != nil {
		writeError(w, http.StatusBadGateway, "upstream profile service error")
		return
	}
	writeJSON(w, http.StatusOK, results)
}
