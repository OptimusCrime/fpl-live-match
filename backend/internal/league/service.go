// Package league provides FPL league ownership and player owner data.
package league

import (
	"context"
	"strconv"
	"sync"

	"github.com/OptimusCrime/fpl-live-match/backend/internal/bootstrap"
	"github.com/OptimusCrime/fpl-live-match/backend/internal/fpl"
)

// Ownership holds pick counts and owner details across all league members for the current gameweek.
type Ownership struct {
	Counts       map[int]int          `json:"counts"`
	Owners       map[int][]OwnerEntry `json:"owners"`
	TotalEntries int                  `json:"totalEntries"`
}

type OwnerEntry struct {
	EntryName  string
	PlayerName string
	Rank       int
}

// PlayerOwners holds a list of league members who own a specific player.
type PlayerOwners struct {
	PlayerID     int
	Owners       []OwnerEntry
	TotalEntries int
}

// Service fetches league ownership and player owner data.
type Service struct {
	fpl       fpl.Client
	bootstrap *bootstrap.Service
}

// NewService creates a new league Service.
func NewService(fpl fpl.Client, bootstrap *bootstrap.Service) *Service {
	return &Service{fpl: fpl, bootstrap: bootstrap}
}

func (s *Service) LoadOwnership(ctx context.Context, leagueID string) (*Ownership, error) {
	gwInfo, err := s.bootstrap.CurrentGameweek(ctx)
	if err != nil {
		return nil, err
	}

	allEntries, err := s.fetchAllStandings(ctx, leagueID)
	if err != nil {
		return nil, err
	}

	type result struct {
		picks *fpl.Picks
		err   error
	}
	results := make([]result, len(allEntries))
	sem := make(chan struct{}, 5)
	var wg sync.WaitGroup

	for i, entry := range allEntries {
		wg.Add(1)
		go func(idx int, entryID int) {
			defer wg.Done()
			sem <- struct{}{}
			defer func() { <-sem }()
			p, err := s.fpl.Picks(ctx, strconv.Itoa(entryID), gwInfo.ID)
			results[idx] = result{p, err}
		}(i, entry.Entry)
	}
	wg.Wait()

	counts := make(map[int]int)
	owners := make(map[int][]OwnerEntry)
	for i, r := range results {
		if r.err != nil || r.picks == nil {
			continue
		}
		entry := allEntries[i]
		for _, pick := range r.picks.Picks {
			counts[pick.Element]++
			owners[pick.Element] = append(owners[pick.Element], OwnerEntry{
				EntryName:  entry.EntryName,
				PlayerName: entry.PlayerName,
				Rank:       entry.Rank,
			})
		}
	}

	return &Ownership{Counts: counts, Owners: owners, TotalEntries: len(allEntries)}, nil
}

func (s *Service) PlayerOwners(ctx context.Context, leagueID string, elementID int) (*PlayerOwners, error) {
	ownership, err := s.LoadOwnership(ctx, leagueID)
	if err != nil {
		return nil, err
	}

	owners := ownership.Owners[elementID]
	if owners == nil {
		owners = []OwnerEntry{}
	}

	return &PlayerOwners{PlayerID: elementID, Owners: owners, TotalEntries: ownership.TotalEntries}, nil
}

// maxStandingsEntries caps how many league members we fetch picks for.
// FPL returns 50 entries per page, so this is 4 pages at most.
const maxStandingsEntries = 200

func (s *Service) fetchAllStandings(ctx context.Context, leagueID string) ([]fpl.StandingEntry, error) {
	var allEntries []fpl.StandingEntry
	for page := 1; ; page++ {
		standings, err := s.fpl.LeagueStandings(ctx, leagueID, page)
		if err != nil {
			return nil, err
		}
		allEntries = append(allEntries, standings.Standings.Results...)
		if !standings.Standings.HasNext || len(allEntries) >= maxStandingsEntries {
			break
		}
	}
	if len(allEntries) > maxStandingsEntries {
		allEntries = allEntries[:maxStandingsEntries]
	}
	return allEntries, nil
}
