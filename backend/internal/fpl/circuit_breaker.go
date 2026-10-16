package fpl

import (
	"errors"
	"log/slog"
	"sync"
	"time"
)

// ErrCircuitOpen is returned when the circuit breaker is tripped and the FPL
// API is being held off. Callers should map this to 503 Service Unavailable.
var ErrCircuitOpen = errors.New("fpl: circuit breaker open, upstream unavailable")

const (
	cbThreshold  = 5               // consecutive failures before opening
	cbInitBackoff = 30 * time.Second
	cbMaxBackoff  = 5 * time.Minute
)

// circuitBreaker is a simple three-state (closed → open → half-open) guard.
// It is embedded in HTTPClient so there is one breaker per process.
type circuitBreaker struct {
	mu        sync.Mutex
	failures  int
	openUntil time.Time
	backoff   time.Duration
}

// allow returns true if a request may proceed. When the circuit is open and
// the backoff period has elapsed, it returns true once to probe the upstream
// (half-open). On the next failure the backoff doubles; on success it resets.
func (cb *circuitBreaker) allow() bool {
	cb.mu.Lock()
	defer cb.mu.Unlock()
	if cb.failures < cbThreshold {
		return true
	}
	// Open → let one probe through once the backoff expires.
	return time.Now().After(cb.openUntil)
}

func (cb *circuitBreaker) success() {
	cb.mu.Lock()
	defer cb.mu.Unlock()
	if cb.failures >= cbThreshold {
		slog.Info("fpl circuit breaker closed, upstream recovered", "failures", cb.failures)
	}
	cb.failures = 0
	cb.backoff = 0
	cb.openUntil = time.Time{}
}

// failure records one upstream error. Only 429 and 5xx responses (plus network
// errors) should be treated as failures; 4xx client errors are not counted.
func (cb *circuitBreaker) failure() {
	cb.mu.Lock()
	defer cb.mu.Unlock()
	cb.failures++
	if cb.failures >= cbThreshold {
		if cb.backoff == 0 {
			cb.backoff = cbInitBackoff
		} else {
			cb.backoff = min(cb.backoff*2, cbMaxBackoff)
		}
		cb.openUntil = time.Now().Add(cb.backoff)
		slog.Error("fpl circuit breaker open",
			"failures", cb.failures,
			"backoff", cb.backoff,
			"open_until", cb.openUntil,
		)
	}
}
