// Package team provides FPL team and league membership data.
package team

import (
	"context"
	"fmt"

	"github.com/OptimusCrime/fpl-live-match/backend/internal/fpl"
)

// League represents an FPL classic league that a team belongs to.
type League struct {
	ID        int
	Name      string
	EntryRank int
}

// TeamLeagues holds a team's identity and the classic leagues it belongs to.
type TeamLeagues struct {
	ID      int
	Name    string
	Leagues []League
}

// Service fetches FPL team and league data.
type Service struct {
	fpl fpl.Client
}

// NewService creates a new team Service.
func NewService(fpl fpl.Client) *Service {
	return &Service{fpl: fpl}
}

func (s *Service) Leagues(ctx context.Context, id string) (*TeamLeagues, error) {
	e, err := s.fpl.Entry(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("fetch entry %s: %w", id, err)
	}

	leagues := make([]League, len(e.Leagues.Classic))
	for i, l := range e.Leagues.Classic {
		leagues[i] = League{ID: l.ID, Name: l.Name, EntryRank: l.EntryRank}
	}
	return &TeamLeagues{ID: e.ID, Name: e.Name, Leagues: leagues}, nil
}
