package fpl

import (
	"errors"
	"testing"
	"time"
)

// --- gameweekTTL ---

func TestGameweekTTLNilSlice(t *testing.T) {
	if got := gameweekTTL(nil); got != ttlActive {
		t.Fatalf("nil fixtures: want ttlActive (%v), got %v", ttlActive, got)
	}
}

func TestGameweekTTLEmptySlice(t *testing.T) {
	if got := gameweekTTL([]Fixture{}); got != ttlActive {
		t.Fatalf("empty fixtures: want ttlActive (%v), got %v", ttlActive, got)
	}
}

func TestGameweekTTLAllFinished(t *testing.T) {
	fixtures := []Fixture{{Finished: true}, {Finished: true}, {Finished: true}}
	if got := gameweekTTL(fixtures); got != ttlCompleted {
		t.Fatalf("all finished: want ttlCompleted (%v), got %v", ttlCompleted, got)
	}
}

func TestGameweekTTLOneUnfinished(t *testing.T) {
	fixtures := []Fixture{{Finished: true}, {Finished: false}}
	if got := gameweekTTL(fixtures); got != ttlActive {
		t.Fatalf("one unfinished: want ttlActive (%v), got %v", ttlActive, got)
	}
}

func TestGameweekTTLAllUnfinished(t *testing.T) {
	fixtures := []Fixture{{Finished: false}, {Finished: false}}
	if got := gameweekTTL(fixtures); got != ttlActive {
		t.Fatalf("all unfinished: want ttlActive (%v), got %v", ttlActive, got)
	}
}

// --- picksTTL ---

func TestPicksTTLOnError(t *testing.T) {
	got := picksTTL(nil, errors.New("upstream down"), 38)
	if got != ttlActive {
		t.Fatalf("error case: want ttlActive, got %v", got)
	}
}

func TestPicksTTLBeforeDeadline(t *testing.T) {
	future := time.Now().Add(time.Hour).Format(time.RFC3339)
	b := &Bootstrap{Events: []Event{{ID: 38, DeadlineTime: future}}}
	if got := picksTTL(b, nil, 38); got != ttlActive {
		t.Fatalf("before deadline: want ttlActive, got %v", got)
	}
}

func TestPicksTTLAfterDeadline(t *testing.T) {
	past := time.Now().Add(-time.Hour).Format(time.RFC3339)
	b := &Bootstrap{Events: []Event{{ID: 38, DeadlineTime: past}}}
	if got := picksTTL(b, nil, 38); got != ttlCompleted {
		t.Fatalf("after deadline: want ttlCompleted, got %v", got)
	}
}

func TestPicksTTLNoMatchingEvent(t *testing.T) {
	b := &Bootstrap{Events: []Event{{ID: 1}}}
	if got := picksTTL(b, nil, 99); got != ttlActive {
		t.Fatalf("no matching GW: want ttlActive, got %v", got)
	}
}

func TestPicksTTLBadDeadlineFormat(t *testing.T) {
	b := &Bootstrap{Events: []Event{{ID: 5, DeadlineTime: "not-a-date"}}}
	if got := picksTTL(b, nil, 5); got != ttlActive {
		t.Fatalf("invalid deadline: want ttlActive, got %v", got)
	}
}

func TestPicksTTLEmptyEvents(t *testing.T) {
	b := &Bootstrap{}
	if got := picksTTL(b, nil, 1); got != ttlActive {
		t.Fatalf("empty events: want ttlActive, got %v", got)
	}
}
