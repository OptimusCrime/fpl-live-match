// Package fpl provides an HTTP client for the Fantasy Premier League API.
package fpl

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

const baseURL = "https://fantasy.premierleague.com/api"

// Mimic a real browser to avoid 403s from FPL's CDN.
const userAgent = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/124.0.0.0 Safari/537.36"

// Client is the interface for fetching FPL data. Callers should not care whether
// results come from the FPL API directly or from a cache.
type Client interface {
	Bootstrap(ctx context.Context) (*Bootstrap, error)
	Fixtures(ctx context.Context, gw int) ([]Fixture, error)
	Live(ctx context.Context, gw int) (*LiveResponse, error)
	Entry(ctx context.Context, id string) (*Entry, error)
	Picks(ctx context.Context, entryID string, gw int) (*Picks, error)
	LeagueStandings(ctx context.Context, leagueID string, page int) (*LeagueStandings, error)
}

// HTTPClient is the raw FPL API client with no caching.
type HTTPClient struct {
	http *http.Client
	cb   circuitBreaker
}

// NewHTTPClient creates a new raw FPL API client with a 30-second timeout.
func NewHTTPClient() *HTTPClient {
	return &HTTPClient{
		http: &http.Client{Timeout: 30 * time.Second},
	}
}

func (c *HTTPClient) get(ctx context.Context, path string, v any) error {
	if !c.cb.allow() {
		return ErrCircuitOpen
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, baseURL+path, nil)
	if err != nil {
		return err
	}
	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("Accept", "application/json, text/javascript, */*; q=0.01")
	req.Header.Set("Accept-Language", "en-GB,en;q=0.9")
	req.Header.Set("Referer", "https://fantasy.premierleague.com/")

	resp, err := c.http.Do(req)
	if err != nil {
		c.cb.failure()
		return fmt.Errorf("fpl: GET %s: %w", path, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		// Count rate-limit and server errors against the circuit; not 4xx client errors.
		if resp.StatusCode == http.StatusTooManyRequests || resp.StatusCode >= 500 {
			c.cb.failure()
		}
		return fmt.Errorf("fpl: GET %s: status %d", path, resp.StatusCode)
	}

	if err := json.NewDecoder(resp.Body).Decode(v); err != nil {
		return err
	}
	c.cb.success()
	return nil
}

func (c *HTTPClient) Bootstrap(ctx context.Context) (*Bootstrap, error) {
	var b Bootstrap
	return &b, c.get(ctx, "/bootstrap-static/", &b)
}

func (c *HTTPClient) Fixtures(ctx context.Context, gw int) ([]Fixture, error) {
	var f []Fixture
	return f, c.get(ctx, fmt.Sprintf("/fixtures/?event=%d", gw), &f)
}

func (c *HTTPClient) Live(ctx context.Context, gw int) (*LiveResponse, error) {
	var l LiveResponse
	return &l, c.get(ctx, fmt.Sprintf("/event/%d/live/", gw), &l)
}

func (c *HTTPClient) Entry(ctx context.Context, id string) (*Entry, error) {
	var e Entry
	return &e, c.get(ctx, fmt.Sprintf("/entry/%s/", id), &e)
}

func (c *HTTPClient) Picks(ctx context.Context, entryID string, gw int) (*Picks, error) {
	var p Picks
	return &p, c.get(ctx, fmt.Sprintf("/entry/%s/event/%d/picks/", entryID, gw), &p)
}

func (c *HTTPClient) LeagueStandings(ctx context.Context, leagueID string, page int) (*LeagueStandings, error) {
	var s LeagueStandings
	path := fmt.Sprintf("/leagues-classic/%s/standings/?page_new_entries=1&page_standings=%d&phase=1", leagueID, page)
	return &s, c.get(ctx, path, &s)
}
