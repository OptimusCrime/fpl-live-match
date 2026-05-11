package team

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/OptimusCrime/fpl-live-match/backend/internal/render"
	"github.com/OptimusCrime/fpl-live-match/backend/internal/resterr"
)

type leagueDTO struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	EntryRank int    `json:"entryRank"`
}

type teamLeaguesDTO struct {
	ID      int         `json:"id"`
	Name    string      `json:"name"`
	Leagues []leagueDTO `json:"leagues"`
}

// Handler serves HTTP requests for team data.
type Handler struct {
	s *Service
}

// NewHandler creates a new team Handler.
func NewHandler(s *Service) *Handler {
	return &Handler{s: s}
}

func (h *Handler) Leagues(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := chi.URLParam(r, "id")

	result, err := h.s.Leagues(ctx, id)
	if err != nil {
		render.JSON(w, r, resterr.New("could not fetch team: "+err.Error(), http.StatusNotFound))
		return
	}

	render.JSON(w, r, toTeamLeaguesDTO(result))
}
