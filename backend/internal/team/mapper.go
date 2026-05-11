package team

func toTeamLeaguesDTO(t *TeamLeagues) teamLeaguesDTO {
	leagues := make([]leagueDTO, len(t.Leagues))
	for i, l := range t.Leagues {
		leagues[i] = leagueDTO{ID: l.ID, Name: l.Name, EntryRank: l.EntryRank}
	}
	return teamLeaguesDTO{ID: t.ID, Name: t.Name, Leagues: leagues}
}
