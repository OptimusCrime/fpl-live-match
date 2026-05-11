package leaguematch

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"github.com/OptimusCrime/fpl-live-match/backend/internal/fixture"
	"github.com/OptimusCrime/fpl-live-match/backend/internal/fpl"
	"github.com/OptimusCrime/fpl-live-match/backend/internal/render"
	"github.com/OptimusCrime/fpl-live-match/backend/internal/resterr"
)

type leagueMatchFixtureDTO struct {
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

type leagueMatchPlayerDTO struct {
	ID         int    `json:"id"`
	WebName    string `json:"webName"`
	PhotoURL   string `json:"photoUrl"`
	Position   string `json:"position"`
	LivePoints int    `json:"livePoints"`
	Minutes    int    `json:"minutes"`
	Started    bool   `json:"started"`
	IsStarter  bool   `json:"isStarter"`
	SubbedOn   bool   `json:"subbedOn"`
	SubbedOff  bool   `json:"subbedOff"`
	OwnedCount int    `json:"ownedCount"`
	IsOwned    bool   `json:"isOwned"`
}

type leagueMatchTeamDTO struct {
	ID        int                    `json:"id"`
	Name      string                 `json:"name"`
	ShortName string                 `json:"shortName"`
	Code      int                    `json:"code"`
	Players   []leagueMatchPlayerDTO `json:"players"`
}

type leagueMatchDTO struct {
	Fixture      leagueMatchFixtureDTO `json:"fixture"`
	Home         leagueMatchTeamDTO    `json:"home"`
	Away         leagueMatchTeamDTO    `json:"away"`
	TotalEntries int                   `json:"totalEntries"`
}

// Handler serves HTTP requests for league match data.
type Handler struct {
	s *Service
}

// NewHandler creates a new leaguematch Handler.
func NewHandler(s *Service) *Handler {
	return &Handler{s: s}
}

func (h *Handler) LeagueMatch(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	leagueID := chi.URLParam(r, "leagueId")

	matchID, err := strconv.Atoi(chi.URLParam(r, "matchId"))
	if err != nil {
		render.JSON(w, r, resterr.New("invalid matchId", http.StatusBadRequest))
		return
	}

	profileID, _ := strconv.Atoi(r.Header.Get("X-CURRENT-PROFILE-ID"))

	result, err := h.s.LeagueMatch(ctx, leagueID, matchID, profileID)
	if err != nil {
		if errors.Is(err, fixture.ErrFixtureNotFound) {
			render.JSON(w, r, resterr.New("fixture not found", http.StatusNotFound))
			return
		}
		if errors.Is(err, fpl.ErrCircuitOpen) {
			render.JSON(w, r, resterr.New("FPL API temporarily unavailable", http.StatusServiceUnavailable))
			return
		}
		render.JSON(w, r, resterr.FromErr(err, http.StatusBadGateway))
		return
	}

	render.JSON(w, r, toLeagueMatchDTO(result))
}
