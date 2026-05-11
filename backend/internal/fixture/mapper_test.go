package fixture

import "testing"

func TestFixtureStatusFinished(t *testing.T) {
	f := Fixture{Finished: true, Started: true}
	if got := FixtureStatus(f); got != "finished" {
		t.Fatalf("want finished, got %q", got)
	}
}

func TestFixtureStatusLive(t *testing.T) {
	f := Fixture{Finished: false, Started: true}
	if got := FixtureStatus(f); got != "live" {
		t.Fatalf("want live, got %q", got)
	}
}

func TestFixtureStatusUpcoming(t *testing.T) {
	f := Fixture{Finished: false, Started: false}
	if got := FixtureStatus(f); got != "upcoming" {
		t.Fatalf("want upcoming, got %q", got)
	}
}

// Finished takes precedence even if Started is somehow false (defensive check).
func TestFixtureStatusFinishedTakesPrecedence(t *testing.T) {
	f := Fixture{Finished: true, Started: false}
	if got := FixtureStatus(f); got != "finished" {
		t.Fatalf("Finished should take precedence over Started=false, got %q", got)
	}
}
