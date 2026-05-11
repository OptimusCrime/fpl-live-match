package bootstrap

import (
	"testing"

	"github.com/OptimusCrime/fpl-live-match/backend/internal/fpl"
)

// --- currentGameweek ---

func TestCurrentGameweekIsCurrentFlag(t *testing.T) {
	events := []fpl.Event{
		{ID: 1, Finished: true},
		{ID: 2, IsCurrent: true, DeadlineTime: "2025-04-12T11:00:00Z"},
		{ID: 3},
	}
	got := currentGameweek(events)
	if got.ID != 2 {
		t.Fatalf("want GW 2 (IsCurrent), got %d", got.ID)
	}
	if got.DeadlineTime != "2025-04-12T11:00:00Z" {
		t.Fatalf("unexpected deadline: %q", got.DeadlineTime)
	}
}

func TestCurrentGameweekFallbackToLastFinished(t *testing.T) {
	events := []fpl.Event{
		{ID: 3, Finished: true},
		{ID: 1, Finished: true},
		{ID: 2, Finished: true},
	}
	got := currentGameweek(events)
	if got.ID != 3 {
		t.Fatalf("fallback should pick last finished GW (3), got %d", got.ID)
	}
}

func TestCurrentGameweekNoEvents(t *testing.T) {
	got := currentGameweek(nil)
	if got.ID != 0 {
		t.Fatalf("empty events: want ID 0, got %d", got.ID)
	}
}

func TestCurrentGameweekNeitherCurrentNorFinished(t *testing.T) {
	events := []fpl.Event{{ID: 5, IsCurrent: false, Finished: false}}
	got := currentGameweek(events)
	// No finished events → fallback returns zero-value (ID 0).
	if got.ID != 0 {
		t.Fatalf("no finished events: want zero-value ID 0, got %d", got.ID)
	}
}

// --- mapTeams ---

func TestMapTeamsBasic(t *testing.T) {
	src := []fpl.Team{
		{ID: 1, Name: "Arsenal", ShortName: "ARS", Code: 3},
		{ID: 2, Name: "Chelsea", ShortName: "CHE", Code: 8},
	}
	got := mapTeams(src)
	if len(got) != 2 {
		t.Fatalf("want 2 teams, got %d", len(got))
	}
	ars := got[1]
	if ars.Name != "Arsenal" || ars.ShortName != "ARS" || ars.Code != 3 {
		t.Fatalf("Arsenal fields wrong: %+v", ars)
	}
}

func TestMapTeamsEmpty(t *testing.T) {
	if got := mapTeams(nil); len(got) != 0 {
		t.Fatalf("nil input: want empty map, got %d entries", len(got))
	}
}

// --- mapPlayers ---

func TestMapPlayersPreservesFields(t *testing.T) {
	src := []fpl.Element{
		{ID: 10, Code: 99, WebName: "Salah", Team: 14, ElementType: 3, Photo: "99.jpg", EventPoints: 12},
	}
	got := mapPlayers(src)
	if len(got) != 1 {
		t.Fatalf("want 1 player, got %d", len(got))
	}
	p := got[0]
	if p.ID != 10 || p.Code != 99 || p.WebName != "Salah" || p.Team != 14 || p.ElementType != 3 || p.Photo != "99.jpg" || p.EventPoints != 12 {
		t.Fatalf("player fields wrong: %+v", p)
	}
}

func TestMapPlayersEmpty(t *testing.T) {
	if got := mapPlayers(nil); len(got) != 0 {
		t.Fatalf("nil input: want empty slice, got %d", len(got))
	}
}

// --- mapPositions ---

func TestMapPositionsBasic(t *testing.T) {
	src := []fpl.ElementType{
		{ID: 1, SingularShort: "GKP"},
		{ID: 2, SingularShort: "DEF"},
		{ID: 3, SingularShort: "MID"},
		{ID: 4, SingularShort: "FWD"},
	}
	got := mapPositions(src)
	if len(got) != 4 {
		t.Fatalf("want 4 positions, got %d", len(got))
	}
	if got[3] != "MID" {
		t.Fatalf("position 3 should be MID, got %q", got[3])
	}
}

func TestMapPositionsEmpty(t *testing.T) {
	if got := mapPositions(nil); len(got) != 0 {
		t.Fatalf("nil input: want empty map, got %d entries", len(got))
	}
}
