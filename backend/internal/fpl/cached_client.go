package fpl

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	"github.com/redis/go-redis/v9"
)

const (
	keyBootstrap      = "bootstrap:raw"
	keyFixtures       = "fixtures:%d"
	keyLive           = "live:%d"
	keyEntry          = "entry:%s"
	keyPicks          = "picks:%s:%d"
	keyLeagueStandings = "league_standings:%s:%d"

	ttlBootstrap       = 5 * time.Minute
	ttlActive          = 30 * time.Second
	ttlCompleted       = 4 * time.Hour
	ttlEntry           = 60 * time.Minute
	ttlLeagueStandings = 24 * time.Hour
)

// CachedClient wraps HTTPClient and transparently caches all responses in Redis.
type CachedClient struct {
	http  *HTTPClient
	cache *redis.Client
}

// NewCachedClient creates a CachedClient backed by the given Redis instance.
func NewCachedClient(cache *redis.Client) *CachedClient {
	return &CachedClient{http: NewHTTPClient(), cache: cache}
}

func (c *CachedClient) Bootstrap(ctx context.Context) (*Bootstrap, error) {
	if val, err := c.cache.Get(ctx, keyBootstrap).Result(); err == nil {
		var b Bootstrap
		if err := json.Unmarshal([]byte(val), &b); err == nil {
			return &b, nil
		}
		slog.Warn("cache unmarshal failed, re-fetching", "key", keyBootstrap)
	}
	b, err := c.http.Bootstrap(ctx)
	if err != nil {
		return nil, err
	}
	if data, err := json.Marshal(b); err == nil {
		c.cache.Set(ctx, keyBootstrap, data, ttlBootstrap)
	}
	return b, nil
}

func (c *CachedClient) Fixtures(ctx context.Context, gw int) ([]Fixture, error) {
	key := fmt.Sprintf(keyFixtures, gw)
	if val, err := c.cache.Get(ctx, key).Result(); err == nil {
		var fixtures []Fixture
		if err := json.Unmarshal([]byte(val), &fixtures); err == nil {
			return fixtures, nil
		}
		slog.Warn("cache unmarshal failed, re-fetching", "key", key)
	}
	fixtures, err := c.http.Fixtures(ctx, gw)
	if err != nil {
		return nil, err
	}
	if data, err := json.Marshal(fixtures); err == nil {
		c.cache.Set(ctx, key, data, gameweekTTL(fixtures))
	}
	return fixtures, nil
}

func (c *CachedClient) Live(ctx context.Context, gw int) (*LiveResponse, error) {
	key := fmt.Sprintf(keyLive, gw)
	if val, err := c.cache.Get(ctx, key).Result(); err == nil {
		var live LiveResponse
		if err := json.Unmarshal([]byte(val), &live); err == nil {
			return &live, nil
		}
		slog.Warn("cache unmarshal failed, re-fetching", "key", key)
	}
	live, err := c.http.Live(ctx, gw)
	if err != nil {
		return nil, err
	}
	// Use fixture status to determine TTL; Fixtures(gw) will be a cache hit here.
	ttl := ttlActive
	if fixtures, err := c.Fixtures(ctx, gw); err == nil {
		ttl = gameweekTTL(fixtures)
	}
	if data, err := json.Marshal(live); err == nil {
		c.cache.Set(ctx, key, data, ttl)
	}
	return live, nil
}

func (c *CachedClient) Entry(ctx context.Context, id string) (*Entry, error) {
	key := fmt.Sprintf(keyEntry, id)
	if val, err := c.cache.Get(ctx, key).Result(); err == nil {
		var e Entry
		if err := json.Unmarshal([]byte(val), &e); err == nil {
			return &e, nil
		}
		slog.Warn("cache unmarshal failed, re-fetching", "key", key)
	}
	e, err := c.http.Entry(ctx, id)
	if err != nil {
		return nil, err
	}
	if data, err := json.Marshal(e); err == nil {
		c.cache.Set(ctx, key, data, ttlEntry)
	}
	return e, nil
}

func (c *CachedClient) Picks(ctx context.Context, entryID string, gw int) (*Picks, error) {
	key := fmt.Sprintf(keyPicks, entryID, gw)
	if val, err := c.cache.Get(ctx, key).Result(); err == nil {
		var p Picks
		if err := json.Unmarshal([]byte(val), &p); err == nil {
			return &p, nil
		}
		slog.Warn("cache unmarshal failed, re-fetching", "key", key)
	}
	p, err := c.http.Picks(ctx, entryID, gw)
	if err != nil {
		return nil, err
	}
	if data, err := json.Marshal(p); err == nil {
		bootstrap, bootstrapErr := c.Bootstrap(ctx)
		c.cache.Set(ctx, key, data, picksTTL(bootstrap, bootstrapErr, gw))
	}
	return p, nil
}

// picksTTL returns a short TTL before the GW deadline (captain swaps still possible)
// and a long TTL once the deadline has passed (picks are locked for the rest of the GW).
func picksTTL(bootstrap *Bootstrap, err error, gw int) time.Duration {
	if err != nil {
		return ttlActive
	}
	for _, e := range bootstrap.Events {
		if e.ID != gw {
			continue
		}
		deadline, parseErr := time.Parse(time.RFC3339, e.DeadlineTime)
		if parseErr != nil {
			return ttlActive
		}
		if time.Now().Before(deadline) {
			return ttlActive
		}
		return ttlCompleted
	}
	return ttlActive
}

func (c *CachedClient) LeagueStandings(ctx context.Context, leagueID string, page int) (*LeagueStandings, error) {
	key := fmt.Sprintf(keyLeagueStandings, leagueID, page)
	if val, err := c.cache.Get(ctx, key).Result(); err == nil {
		var s LeagueStandings
		if err := json.Unmarshal([]byte(val), &s); err == nil {
			return &s, nil
		}
		slog.Warn("cache unmarshal failed, re-fetching", "key", key)
	}
	s, err := c.http.LeagueStandings(ctx, leagueID, page)
	if err != nil {
		return nil, err
	}
	if data, err := json.Marshal(s); err == nil {
		c.cache.Set(ctx, key, data, ttlLeagueStandings)
	}
	return s, nil
}

// gameweekTTL returns a short TTL while any fixture is still ongoing, and a
// long TTL once all fixtures in the gameweek are finished.
func gameweekTTL(fixtures []Fixture) time.Duration {
	if len(fixtures) == 0 {
		return ttlActive
	}
	for _, f := range fixtures {
		if !f.Finished {
			return ttlActive
		}
	}
	return ttlCompleted
}
