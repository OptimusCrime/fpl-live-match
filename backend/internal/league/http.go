package league

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"github.com/OptimusCrime/fpl-live-match/backend/internal/fpl"
	"github.com/OptimusCrime/fpl-live-match/backend/internal/render"
	"github.com/OptimusCrime/fpl-live-match/backend/internal/resterr"
)

type playerOwnerDTO struct {
	EntryName  string `json:"entryName"`
	PlayerName string `json:"playerName"`
	Rank       int    `json:"rank"`
}

type playerOwnersDTO struct {
	PlayerID     int              `json:"playerId"`
	Owners       []playerOwnerDTO `json:"owners"`
	TotalEntries int              `json:"totalEntries"`
}

// Handler serves HTTP requests for league ownership data.
type Handler struct {
	s *Service
}

// NewHandler creates a new league Handler.
func NewHandler(s *Service) *Handler {
	return &Handler{s: s}
}

func (h *Handler) PlayerOwners(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	elementID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		render.JSON(w, r, resterr.New("invalid id", http.StatusBadRequest))
		return
	}
	leagueID := r.URL.Query().Get("league")
	if leagueID == "" {
		render.JSON(w, r, resterr.New("league query param required", http.StatusBadRequest))
		return
	}

	result, err := h.s.PlayerOwners(ctx, leagueID, elementID)
	if err != nil {
		if errors.Is(err, fpl.ErrCircuitOpen) {
			render.JSON(w, r, resterr.New("FPL API temporarily unavailable", http.StatusServiceUnavailable))
			return
		}
		render.JSON(w, r, resterr.New("could not fetch ownership: "+err.Error(), http.StatusBadGateway))
		return
	}

	render.JSON(w, r, toPlayerOwnersDTO(result))
}
