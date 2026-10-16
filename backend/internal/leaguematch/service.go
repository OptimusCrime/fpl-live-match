// Package leaguematch provides data combining league ownership with live fixture results.
package leaguematch

import (
	"context"
	"fmt"
	"log/slog"
	"strconv"
	"sync"

	"github.com/OptimusCrime/fpl-live-match/backend/internal/bootstrap"
	"github.com/OptimusCrime/fpl-live-match/backend/internal/fixture"
	"github.com/OptimusCrime/fpl-live-match/backend/internal/fpl"
)

// Service fetches league match data enriched with ownership information.
type Service struct {
	fpl       fpl.Client
	bootstrap *bootstrap.Service
	fixture   *fixture.Service
}

// NewService creates a new leaguematch Service.
func NewService(fpl fpl.Client, bootstrap *bootstrap.Service, fixture *fixture.Service) *Service {
	return &Service{fpl: fpl, bootstrap: bootstrap, fixture: fixture}
}

func (s *Service) LeagueMatch(ctx context.Context, leagueID string, matchID int, profileID int) (*LeagueMatchResult, error) {
	gwInfo, err := s.bootstrap.CurrentGameweek(ctx)
	if err != nil {
		return nil, err
	}

	members, err := s.loadLeagueMembers(ctx, leagueID)
	if err != nil {
		return nil, fmt.Errorf("load members for league %s: %w", leagueID, err)
	}

	allPicks, err := s.loadAllPicks(ctx, leagueID, gwInfo.ID, members)
	if err != nil {
		return nil, err
	}

	leagueOwned := make(map[int]bool)
	for _, picks := range allPicks {
		for _, p := range picks {
			leagueOwned[p.Element] = true
		}
	}

	matchResult, err := s.fixture.Match(ctx, matchID, leagueOwned)
	if err != nil {
		return nil, fmt.Errorf("load match %d: %w", matchID, err)
	}

	ownedCount := make(map[int]int)
	for _, picks := range allPicks {
		for _, p := range picks {
			ownedCount[p.Element]++
		}
	}

	profilePickSet := make(map[int]bool)
	if profileID > 0 {
		for _, p := range allPicks[profileID] {
			profilePickSet[p.Element] = true
		}
	}

	home := enrichPlayers(matchResult.Home, ownedCount, profilePickSet)
	away := enrichPlayers(matchResult.Away, ownedCount, profilePickSet)

	return &LeagueMatchResult{
		Fixture:      matchResult.Fixture,
		Teams:        matchResult.Teams,
		Home:         home,
		Away:         away,
		TotalEntries: len(members),
	}, nil
}

func (s *Service) loadLeagueMembers(ctx context.Context, leagueID string) ([]Member, error) {
	var members []Member
	for page := 1; len(members) < MaxLeagueMembers; page++ {
		standings, err := s.fpl.LeagueStandings(ctx, leagueID, page)
		if err != nil {
			return nil, fmt.Errorf("standings page %d: %w", page, err)
		}
		for _, e := range standings.Standings.Results {
			members = append(members, Member{
				EntryID:    e.Entry,
				EntryName:  e.EntryName,
				PlayerName: e.PlayerName,
				Rank:       e.Rank,
			})
			if len(members) >= MaxLeagueMembers {
				break
			}
		}
		if !standings.Standings.HasNext {
			break
		}
	}
	return members, nil
}

func (s *Service) loadAllPicks(ctx context.Context, _ string, gwID int, members []Member) (AllPicks, error) {
	type result struct {
		entryID int
		picks   []MemberPick
	}
	results := make([]result, len(members))
	sem := make(chan struct{}, 5)
	var wg sync.WaitGroup

	for i, m := range members {
		wg.Add(1)
		go func(idx int, entryID int) {
			defer wg.Done()
			sem <- struct{}{}
			defer func() { <-sem }()

			raw, err := s.fpl.Picks(ctx, strconv.Itoa(entryID), gwID)
			if err != nil || raw == nil {
				slog.WarnContext(ctx, "could not fetch picks for league member, excluding from ownership",
					"entry", entryID, "gw", gwID, "err", err)
				results[idx] = result{entryID: entryID}
				return
			}
			picks := make([]MemberPick, 0, len(raw.Picks))
			for _, p := range raw.Picks {
				picks = append(picks, MemberPick{
					Element:       p.Element,
					IsCaptain:     p.IsCaptain,
					IsViceCaptain: p.IsViceCaptain,
					Multiplier:    p.Multiplier,
				})
			}
			results[idx] = result{entryID: entryID, picks: picks}
		}(i, m.EntryID)
	}
	wg.Wait()

	ap := make(AllPicks, len(members))
	for _, r := range results {
		if r.picks != nil {
			ap[r.entryID] = r.picks
		}
	}
	return ap, nil
}

func enrichPlayers(players []fixture.MatchPlayer, ownedCount map[int]int, profilePickSet map[int]bool) []LeagueMatchPlayer {
	result := make([]LeagueMatchPlayer, len(players))
	for i, p := range players {
		result[i] = LeagueMatchPlayer{
			ID:         p.ID,
			WebName:    p.WebName,
			PhotoURL:   p.PhotoURL,
			Position:   p.Position,
			LivePoints: p.LivePoints,
			Minutes:    p.Minutes,
			Started:    p.Started,
			IsStarter:  p.IsStarter,
			SubbedOn:   p.SubbedOn,
			SubbedOff:  p.SubbedOff,
			Team:       p.Team,
			OwnedCount: ownedCount[p.ID],
			IsOwned:    profilePickSet[p.ID],
		}
	}
	return result
}
