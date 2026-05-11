package leaguematch

import "github.com/OptimusCrime/fpl-live-match/backend/internal/fixture"

func toLeagueMatchFixtureDTO(r *LeagueMatchResult) leagueMatchFixtureDTO {
	th := r.Teams[r.Fixture.TeamH]
	ta := r.Teams[r.Fixture.TeamA]
	return leagueMatchFixtureDTO{
		ID:          r.Fixture.ID,
		TeamH:       r.Fixture.TeamH,
		TeamA:       r.Fixture.TeamA,
		TeamHName:   th.Name,
		TeamAName:   ta.Name,
		TeamHShort:  th.ShortName,
		TeamAShort:  ta.ShortName,
		TeamHCode:   th.Code,
		TeamACode:   ta.Code,
		ScoreH:      r.Fixture.TeamHScore,
		ScoreA:      r.Fixture.TeamAScore,
		Status:      fixture.FixtureStatus(r.Fixture),
		KickoffTime: r.Fixture.KickoffTime,
		Minutes:     r.Fixture.Minutes,
	}
}

func toLeagueMatchPlayerDTO(p LeagueMatchPlayer) leagueMatchPlayerDTO {
	return leagueMatchPlayerDTO{
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
		OwnedCount: p.OwnedCount,
		IsOwned:    p.IsOwned,
	}
}

func toLeagueMatchTeamDTO(r *LeagueMatchResult, teamID int, players []LeagueMatchPlayer) leagueMatchTeamDTO {
	team := r.Teams[teamID]
	dtoPlayers := make([]leagueMatchPlayerDTO, len(players))
	for i, p := range players {
		dtoPlayers[i] = toLeagueMatchPlayerDTO(p)
	}
	return leagueMatchTeamDTO{
		ID:        team.ID,
		Name:      team.Name,
		ShortName: team.ShortName,
		Code:      team.Code,
		Players:   dtoPlayers,
	}
}

func toLeagueMatchDTO(r *LeagueMatchResult) leagueMatchDTO {
	return leagueMatchDTO{
		Fixture:      toLeagueMatchFixtureDTO(r),
		Home:         toLeagueMatchTeamDTO(r, r.Fixture.TeamH, r.Home),
		Away:         toLeagueMatchTeamDTO(r, r.Fixture.TeamA, r.Away),
		TotalEntries: r.TotalEntries,
	}
}
