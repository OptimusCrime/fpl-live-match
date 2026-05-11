// Package fixture provides FPL fixture and live match data.
package fixture

import (
	"net/http"

	"github.com/OptimusCrime/fpl-live-match/backend/internal/render"
	"github.com/OptimusCrime/fpl-live-match/backend/internal/resterr"
)

type fixtureDTO struct {
	ID          int    `json:"id"`
	TeamH       int    `json:"teamH"`
	TeamA       int    `json:"teamA"`
	TeamHName   string `json:"teamHName"`
	TeamAName   string `json:"teamAName"`
	TeamHShort  string `json:"teamHShort"`
	TeamAShort  string `json:"teamAShort"`
	TeamHCode   int    `json:"teamHCode"`
	TeamACode   int    `json:"teamACode"`
	ScoreH      *int   `json:"scoreH"`
	ScoreA      *int   `json:"scoreA"`
	Status      string `json:"status"`
	KickoffTime string `json:"kickoffTime"`
	Minutes     int    `json:"minutes"`
}

type fixturesDTO struct {
	GW           int          `json:"gw"`
	DeadlineTime string       `json:"deadlineTime"`
	Fixtures     []fixtureDTO `json:"fixtures"`
}

// Handler serves HTTP requests for fixture data.
type Handler struct {
	s *Service
}

// NewHandler creates a new fixture Handler.
func NewHandler(s *Service) *Handler {
	return &Handler{s: s}
}

func (h *Handler) Fixtures(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	result, err := h.s.Fixtures(ctx)
	if err != nil {
		render.JSON(w, r, resterr.FromErr(err, http.StatusInternalServerError))
		return
	}

	render.JSON(w, r, toFixturesDTO(result))
}
