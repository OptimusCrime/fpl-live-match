package fixture

import (
	"testing"

	"github.com/OptimusCrime/fpl-live-match/backend/internal/fpl"
)

// --- mostCommonStarterMinutes ---

func TestMostCommonStarterMinutesEmpty(t *testing.T) {
	if got := mostCommonStarterMinutes(nil, nil); got != 0 {
		t.Fatalf("empty slices: want 0, got %d", got)
	}
}

func TestMostCommonStarterMinutesNoStarters(t *testing.T) {
	home := []MatchPlayer{{IsStarter: false, Minutes: 60}}
	away := []MatchPlayer{{IsStarter: false, Minutes: 45}}
	if got := mostCommonStarterMinutes(home, away); got != 0 {
		t.Fatalf("no starters: want 0, got %d", got)
	}
}

func TestMostCommonStarterMinutesStartersWithZeroMinutes(t *testing.T) {
	// Starters with 0 minutes are excluded (match hasn't started for them).
	home := []MatchPlayer{{IsStarter: true, Minutes: 0}}
	if got := mostCommonStarterMinutes(home, nil); got != 0 {
		t.Fatalf("starter with 0 minutes: want 0, got %d", got)
	}
}

func TestMostCommonStarterMinutesSingleStarter(t *testing.T) {
	home := []MatchPlayer{{IsStarter: true, Minutes: 72}}
	if got := mostCommonStarterMinutes(home, nil); got != 72 {
		t.Fatalf("single starter: want 72, got %d", got)
	}
}

func TestMostCommonStarterMinutesCommonValue(t *testing.T) {
	home := []MatchPlayer{
		{IsStarter: true, Minutes: 90},
		{IsStarter: true, Minutes: 90},
		{IsStarter: true, Minutes: 72}, // substituted off
	}
	away := []MatchPlayer{
		{IsStarter: true, Minutes: 90},
	}
	if got := mostCommonStarterMinutes(home, away); got != 90 {
		t.Fatalf("want most-common 90, got %d", got)
	}
}

// When two minute values appear equally often, the higher one wins.
func TestMostCommonStarterMinutesTieBreakHigher(t *testing.T) {
	home := []MatchPlayer{
		{IsStarter: true, Minutes: 60},
		{IsStarter: true, Minutes: 90},
	}
	if got := mostCommonStarterMinutes(home, nil); got != 90 {
		t.Fatalf("tie-break: want higher value 90, got %d", got)
	}
}

// --- applySubStatus ---

func TestApplySubStatusNonStarterWithMinutes(t *testing.T) {
	players := []MatchPlayer{{IsStarter: false, Minutes: 30}}
	applySubStatus(players, 90)
	if !players[0].SubbedOn {
		t.Fatal("non-starter with minutes should be SubbedOn")
	}
	if players[0].SubbedOff {
		t.Fatal("non-starter should not be SubbedOff")
	}
}

func TestApplySubStatusStarterWithFewerMinutes(t *testing.T) {
	players := []MatchPlayer{{IsStarter: true, Minutes: 60}}
	applySubStatus(players, 90)
	if !players[0].SubbedOff {
		t.Fatal("starter with fewer minutes than current should be SubbedOff")
	}
	if players[0].SubbedOn {
		t.Fatal("starter should not be SubbedOn")
	}
}

func TestApplySubStatusStarterWithFullMinutes(t *testing.T) {
	players := []MatchPlayer{{IsStarter: true, Minutes: 90}}
	applySubStatus(players, 90)
	if players[0].SubbedOff {
		t.Fatal("starter with full minutes should not be SubbedOff")
	}
}

func TestApplySubStatusZeroCurrentMins(t *testing.T) {
	// If the match hasn't started yet (currentMins == 0), nobody should be SubbedOff.
	players := []MatchPlayer{{IsStarter: true, Minutes: 0}}
	applySubStatus(players, 0)
	if players[0].SubbedOff {
		t.Fatal("should not be SubbedOff when currentMins is 0")
	}
}

func TestApplySubStatusNonStarterNoMinutes(t *testing.T) {
	players := []MatchPlayer{{IsStarter: false, Minutes: 0}}
	applySubStatus(players, 90)
	if players[0].SubbedOn {
		t.Fatal("non-starter with 0 minutes should not be SubbedOn")
	}
}

// --- fromFPLFixture ---

func TestFromFPLFixtureCopiesAllFields(t *testing.T) {
	score := 2
	raw := fpl.Fixture{
		ID:          42,
		TeamH:       7,
		TeamA:       3,
		Event:       30,
		TeamHScore:  &score,
		TeamAScore:  nil,
		KickoffTime: "2025-04-12T15:00:00Z",
		Minutes:     67,
		Finished:    false,
		Started:     true,
	}
	got := fromFPLFixture(raw)

	if got.ID != 42 || got.TeamH != 7 || got.TeamA != 3 || got.Event != 30 {
		t.Fatalf("basic fields mismatch: %+v", got)
	}
	if got.TeamHScore == nil || *got.TeamHScore != 2 {
		t.Fatal("TeamHScore not copied correctly")
	}
	if got.TeamAScore != nil {
		t.Fatal("TeamAScore should be nil")
	}
	if got.KickoffTime != raw.KickoffTime {
		t.Fatalf("KickoffTime mismatch: %q vs %q", got.KickoffTime, raw.KickoffTime)
	}
	if got.Minutes != 67 || got.Finished || !got.Started {
		t.Fatalf("minutes/status fields mismatch: %+v", got)
	}
}

// --- photoURL ---

func TestPhotoURL(t *testing.T) {
	got := photoURL(12345)
	want := "https://resources.premierleague.com/premierleague/photos/players/110x140/p12345.png"
	if got != want {
		t.Fatalf("photoURL: want %q, got %q", want, got)
	}
}
