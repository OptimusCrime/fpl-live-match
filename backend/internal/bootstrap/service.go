// Package bootstrap provides access to the FPL bootstrap data (teams, players, gameweek info).
package bootstrap

import (
	"context"

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
		return nil, err
	}
	return mapTeams(b.Teams), nil
}

func (s *Service) Players(ctx context.Context) ([]Player, error) {
	b, err := s.fpl.Bootstrap(ctx)
	if err != nil {
		return nil, err
	}
	return mapPlayers(b.Elements), nil
}

func (s *Service) Positions(ctx context.Context) (map[int]string, error) {
	b, err := s.fpl.Bootstrap(ctx)
	if err != nil {
		return nil, err
	}
	return mapPositions(b.ElementTypes), nil
}

func (s *Service) CurrentGameweek(ctx context.Context) (GameweekInfo, error) {
	b, err := s.fpl.Bootstrap(ctx)
	if err != nil {
		return GameweekInfo{}, err
	}
	return currentGameweek(b.Events), nil
}
