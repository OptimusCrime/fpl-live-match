package fixture

import "github.com/OptimusCrime/fpl-live-match/backend/internal/bootstrap"

// FixtureStatus returns the display status of a fixture: "finished", "live", or "upcoming".
func FixtureStatus(f Fixture) string {
	if f.Finished {
		return "finished"
	}
	if f.Started {
		return "live"
	}
	return "upcoming"
}

func toFixtureDTO(f Fixture, teams map[int]bootstrap.Team) fixtureDTO {
	th := teams[f.TeamH]
	ta := teams[f.TeamA]
	return fixtureDTO{
		ID:          f.ID,
		TeamH:       f.TeamH,
		TeamA:       f.TeamA,
		TeamHName:   th.Name,
		TeamAName:   ta.Name,
		TeamHShort:  th.ShortName,
		TeamAShort:  ta.ShortName,
		TeamHCode:   th.Code,
		TeamACode:   ta.Code,
		ScoreH:      f.TeamHScore,
		ScoreA:      f.TeamAScore,
		Status:      FixtureStatus(f),
		KickoffTime: f.KickoffTime,
		Minutes:     f.Minutes,
	}
}

func toFixturesDTO(gw *FixturesResult) fixturesDTO {
	fixtures := make([]fixtureDTO, len(gw.Fixtures))
	for i, f := range gw.Fixtures {
		fixtures[i] = toFixtureDTO(f, gw.Teams)
	}
	return fixturesDTO{GW: gw.ID, DeadlineTime: gw.DeadlineTime, Fixtures: fixtures}
}

