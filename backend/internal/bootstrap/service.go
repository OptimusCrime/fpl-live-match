// Package bootstrap provides access to the FPL bootstrap data (teams, players, gameweek info).
package bootstrap

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/OptimusCrime/fpl-live-match/backend/internal/fpl"
)

// Service fetches FPL bootstrap data via a (potentially cached) FPL client.
type Service struct {
	fpl fpl.Client
}

// New creates a new bootstrap Service.
func New(fpl fpl.Client) *Service {
	return &Service{fpl: fpl}
}

func (s *Service) Teams(ctx context.Context) (map[int]Team, error) {
	b, err := s.fpl.Bootstrap(ctx)
	if err != nil {
		return nil, fmt.Errorf("load bootstrap for teams: %w", err)
	}
	return mapTeams(b.Teams), nil
}

func (s *Service) Players(ctx context.Context) ([]Player, error) {
	b, err := s.fpl.Bootstrap(ctx)
	if err != nil {
		return nil, fmt.Errorf("load bootstrap for players: %w", err)
	}
	return mapPlayers(b.Elements), nil
}

func (s *Service) Positions(ctx context.Context) (map[int]string, error) {
	b, err := s.fpl.Bootstrap(ctx)
	if err != nil {
		return nil, fmt.Errorf("load bootstrap for positions: %w", err)
	}
	return mapPositions(b.ElementTypes), nil
}

func (s *Service) CurrentGameweek(ctx context.Context) (GameweekInfo, error) {
	b, err := s.fpl.Bootstrap(ctx)
	if err != nil {
		return GameweekInfo{}, fmt.Errorf("load bootstrap for current gameweek: %w", err)
	}
	gw := currentGameweek(b.Events)
	if gw.ID == 0 {
		slog.WarnContext(ctx, "no current or finished gameweek found in bootstrap", "events", len(b.Events))
	}
	return gw, nil
}
