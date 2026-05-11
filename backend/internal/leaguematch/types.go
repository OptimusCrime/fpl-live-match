package leaguematch

import (
	"github.com/OptimusCrime/fpl-live-match/backend/internal/bootstrap"
	"github.com/OptimusCrime/fpl-live-match/backend/internal/fixture"
)

const MaxLeagueMembers = 100

type Member struct {
	EntryID    int    `json:"entryId"`
	EntryName  string `json:"entryName"`
	PlayerName string `json:"playerName"`
	Rank       int    `json:"rank"`
}

type MemberPick struct {
	Element       int  `json:"element"`
	IsCaptain     bool `json:"isCaptain"`
	IsViceCaptain bool `json:"isViceCaptain"`
	Multiplier    int  `json:"multiplier"`
}

// AllPicks maps entryID to that entry's picks.
type AllPicks map[int][]MemberPick

type LeagueMatchPlayer struct {
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
	OwnedCount int
	IsOwned    bool
}

type LeagueMatchResult struct {
	Fixture      fixture.Fixture
	Teams        map[int]bootstrap.Team
	Home         []LeagueMatchPlayer
	Away         []LeagueMatchPlayer
	TotalEntries int
}
