package fixture

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/OptimusCrime/fpl-live-match/backend/internal/bootstrap"
	"github.com/OptimusCrime/fpl-live-match/backend/internal/fpl"
)

// ErrFixtureNotFound is returned when a requested fixture does not exist.
var ErrFixtureNotFound = errors.New("fixture not found")

// Fixture represents a single FPL fixture (match).
type Fixture struct {
	ID          int
	TeamH       int
	TeamA       int
	Event       int
	TeamHScore  *int
	TeamAScore  *int
	KickoffTime string
	Minutes     int
	Finished    bool
	Started     bool
}

// FixturesResult holds the current gameweek's fixtures and associated team data.
type FixturesResult struct {
	ID           int
	DeadlineTime string
	Fixtures     []Fixture
	Teams        map[int]bootstrap.Team
}

type MatchPlayer struct {
	ID         int
	WebName    string
	PhotoURL   string
	Position   string
	LivePoints int
	Minutes    int
	Started    bool
	IsStarter  bool
	SubbedOn   bool
	SubbedOff  bool
	Team       int
}

// MatchResult holds the live data for a single fixture, including players on each side.
type MatchResult struct {
	Fixture Fixture
	Teams   map[int]bootstrap.Team
	Home    []MatchPlayer
	Away    []MatchPlayer
}

// Service fetches fixture and live match data.
type Service struct {
	fpl       fpl.Client
	bootstrap *bootstrap.Service
}

// NewService creates a new fixture Service.
func NewService(fpl fpl.Client, bootstrap *bootstrap.Service) *Service {
	return &Service{fpl: fpl, bootstrap: bootstrap}
}

func (s *Service) Fixtures(ctx context.Context) (*FixturesResult, error) {
	gwInfo, err := s.bootstrap.CurrentGameweek(ctx)
	if err != nil {
		return nil, err
	}

	rawFixtures, err := s.fpl.Fixtures(ctx, gwInfo.ID)
	if err != nil {
		return nil, err
	}

	teams, err := s.bootstrap.Teams(ctx)
	if err != nil {
		return nil, err
	}

	fixtures := make([]Fixture, 0)
	for _, f := range rawFixtures {
		if f.Event != gwInfo.ID {
			continue
		}
		fixtures = append(fixtures, fromFPLFixture(f))
	}

	return &FixturesResult{
		ID:           gwInfo.ID,
		DeadlineTime: gwInfo.DeadlineTime,
		Fixtures:     fixtures,
		Teams:        teams,
	}, nil
}

func (s *Service) Match(ctx context.Context, fixtureID int, leagueOwned map[int]bool) (*MatchResult, error) {
	gwInfo, err := s.bootstrap.CurrentGameweek(ctx)
	if err != nil {
		return nil, err
	}

	rawFixtures, err := s.fpl.Fixtures(ctx, gwInfo.ID)
	if err != nil {
		return nil, err
	}

	var raw *fpl.Fixture
	for i := range rawFixtures {
		if rawFixtures[i].ID == fixtureID {
			raw = &rawFixtures[i]
			break
		}
	}
	if raw == nil {
		return nil, ErrFixtureNotFound
	}

	live, err := s.fpl.Live(ctx, raw.Event)
	if err != nil {
		return nil, err
	}

	teams, err := s.bootstrap.Teams(ctx)
	if err != nil {
		return nil, err
	}

	players, err := s.bootstrap.Players(ctx)
	if err != nil {
		return nil, err
	}

	positions, err := s.bootstrap.Positions(ctx)
	if err != nil {
		return nil, err
	}

	liveMap := make(map[int]fpl.LiveStats, len(live.Elements))
	for _, el := range live.Elements {
		liveMap[el.ID] = el.Stats
	}

	fixture := fromFPLFixture(*raw)
	home := make([]MatchPlayer, 0)
	away := make([]MatchPlayer, 0)

	for _, el := range players {
		if el.Team != fixture.TeamH && el.Team != fixture.TeamA {
			continue
		}
		stats := liveMap[el.ID]

		points := stats.TotalPoints
		if points == 0 && fixture.Finished {
			points = el.EventPoints
		}

		isPlayed := stats.Minutes > 0 || (fixture.Finished && points > 0)
		if !isPlayed && !leagueOwned[el.ID] {
			continue
		}

		code := el.Code
		if el.Photo != "" {
			if n, err := strconv.Atoi(strings.TrimSuffix(el.Photo, ".jpg")); err == nil {
				code = n
			}
		}

		started := stats.Minutes > 0 || (fixture.Finished && points > 0)
		mp := MatchPlayer{
			ID:         el.ID,
			WebName:    el.WebName,
			PhotoURL:   photoURL(code),
			Position:   positions[el.ElementType],
			LivePoints: points,
			Minutes:    stats.Minutes,
			Started:    started,
			IsStarter:  stats.Starts == 1,
			Team:       el.Team,
		}

		if el.Team == fixture.TeamH {
			home = append(home, mp)
		} else {
			away = append(away, mp)
		}
	}

	currentMins := mostCommonStarterMinutes(home, away)
	applySubStatus(home, currentMins)
	applySubStatus(away, currentMins)

	return &MatchResult{Fixture: fixture, Teams: teams, Home: home, Away: away}, nil
}

// mostCommonStarterMinutes finds the most frequent minutes value among starters
// across both teams, which approximates the current match minute.
func mostCommonStarterMinutes(home, away []MatchPlayer) int {
	freq := make(map[int]int)
	for _, p := range home {
		if p.IsStarter && p.Minutes > 0 {
			freq[p.Minutes]++
		}
	}
	for _, p := range away {
		if p.IsStarter && p.Minutes > 0 {
			freq[p.Minutes]++
		}
	}
	best, bestCount := 0, 0
	for mins, count := range freq {
		if count > bestCount || (count == bestCount && mins > best) {
			best = mins
			bestCount = count
		}
	}
	return best
}

func applySubStatus(players []MatchPlayer, currentMins int) {
	for i := range players {
		players[i].SubbedOn = !players[i].IsStarter && players[i].Minutes > 0
		players[i].SubbedOff = players[i].IsStarter && currentMins > 0 && players[i].Minutes < currentMins
	}
}

func fromFPLFixture(f fpl.Fixture) Fixture {
	return Fixture{
		ID:          f.ID,
		TeamH:       f.TeamH,
		TeamA:       f.TeamA,
		Event:       f.Event,
		TeamHScore:  f.TeamHScore,
		TeamAScore:  f.TeamAScore,
		KickoffTime: f.KickoffTime,
		Minutes:     f.Minutes,
		Finished:    f.Finished,
		Started:     f.Started,
	}
}

func photoURL(code int) string {
	return fmt.Sprintf("https://resources.premierleague.com/premierleague/photos/players/110x140/p%d.png", code)
}
