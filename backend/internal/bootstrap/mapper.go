package bootstrap

import "github.com/OptimusCrime/fpl-live-match/backend/internal/fpl"

func mapTeams(src []fpl.Team) map[int]Team {
	out := make(map[int]Team, len(src))
	for _, t := range src {
		out[t.ID] = Team{
			ID:        t.ID,
			Name:      t.Name,
			ShortName: t.ShortName,
			Code:      t.Code,
		}
	}
	return out
}

func mapPlayers(src []fpl.Element) []Player {
	out := make([]Player, len(src))
	for i, e := range src {
		out[i] = Player{
			ID:          e.ID,
			Code:        e.Code,
			WebName:     e.WebName,
			Team:        e.Team,
			ElementType: e.ElementType,
			Photo:       e.Photo,
			EventPoints: e.EventPoints,
		}
	}
	return out
}

func mapPositions(src []fpl.ElementType) map[int]string {
	out := make(map[int]string, len(src))
	for _, et := range src {
		out[et.ID] = et.SingularShort
	}
	return out
}

func currentGameweek(events []fpl.Event) GameweekInfo {
	for _, e := range events {
		if e.IsCurrent {
			return GameweekInfo{ID: e.ID, DeadlineTime: e.DeadlineTime}
		}
	}
	// Fallback: last finished event.
	var last fpl.Event
	for _, e := range events {
		if e.Finished && e.ID > last.ID {
			last = e
		}
	}
	return GameweekInfo{ID: last.ID, DeadlineTime: last.DeadlineTime}
}
