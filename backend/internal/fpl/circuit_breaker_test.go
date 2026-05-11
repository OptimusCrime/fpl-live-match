package fpl

import (
	"testing"
	"time"
)

func TestCircuitBreakerInitiallyAllows(t *testing.T) {
	cb := &circuitBreaker{}
	if !cb.allow() {
		t.Fatal("fresh circuit breaker should allow requests")
	}
}

func TestCircuitBreakerAllowsBelowThreshold(t *testing.T) {
	cb := &circuitBreaker{}
	for i := 0; i < cbThreshold-1; i++ {
		cb.failure()
	}
	if !cb.allow() {
		t.Fatalf("circuit should stay closed with %d failures (threshold %d)", cbThreshold-1, cbThreshold)
	}
}

func TestCircuitBreakerOpensAtThreshold(t *testing.T) {
	cb := &circuitBreaker{}
	for i := 0; i < cbThreshold; i++ {
		cb.failure()
	}
	if cb.allow() {
		t.Fatal("circuit should be open after reaching threshold")
	}
}

func TestCircuitBreakerSuccessResets(t *testing.T) {
	cb := &circuitBreaker{}
	for i := 0; i < cbThreshold; i++ {
		cb.failure()
	}
	cb.success()
	if !cb.allow() {
		t.Fatal("circuit should close after success")
	}
	if cb.failures != 0 {
		t.Fatalf("failures should be 0 after success, got %d", cb.failures)
	}
	if cb.backoff != 0 {
		t.Fatalf("backoff should be 0 after success, got %v", cb.backoff)
	}
	if !cb.openUntil.IsZero() {
		t.Fatalf("openUntil should be zero after success, got %v", cb.openUntil)
	}
}

// TestCircuitBreakerHalfOpen verifies that a probe is allowed once the backoff window expires.
func TestCircuitBreakerHalfOpen(t *testing.T) {
	cb := &circuitBreaker{}
	for i := 0; i < cbThreshold; i++ {
		cb.failure()
	}
	cb.mu.Lock()
	cb.openUntil = time.Now().Add(-time.Millisecond)
	cb.mu.Unlock()

	if !cb.allow() {
		t.Fatal("circuit should allow one probe after backoff window expires")
	}
}

func TestCircuitBreakerInitialBackoff(t *testing.T) {
	cb := &circuitBreaker{}
	for i := 0; i < cbThreshold; i++ {
		cb.failure()
	}
	if cb.backoff != cbInitBackoff {
		t.Fatalf("expected initial backoff %v, got %v", cbInitBackoff, cb.backoff)
	}
}

func TestCircuitBreakerBackoffDoubles(t *testing.T) {
	cb := &circuitBreaker{}
	for i := 0; i < cbThreshold; i++ {
		cb.failure()
	}
	// Expire window so the probe attempt can proceed, then fail again.
	cb.mu.Lock()
	cb.openUntil = time.Now().Add(-time.Millisecond)
	cb.mu.Unlock()
	cb.failure()

	if cb.backoff != cbInitBackoff*2 {
		t.Fatalf("expected doubled backoff %v, got %v", cbInitBackoff*2, cb.backoff)
	}
}

func TestCircuitBreakerBackoffCapped(t *testing.T) {
	cb := &circuitBreaker{}
	for i := 0; i < cbThreshold; i++ {
		cb.failure()
	}
	// Keep failing beyond the doubling limit.
	for i := 0; i < 10; i++ {
		cb.mu.Lock()
		cb.openUntil = time.Now().Add(-time.Millisecond)
		cb.mu.Unlock()
		cb.failure()
	}
	if cb.backoff != cbMaxBackoff {
		t.Fatalf("expected backoff to saturate at %v, got %v", cbMaxBackoff, cb.backoff)
	}
}
