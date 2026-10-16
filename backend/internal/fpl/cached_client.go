package fpl

import (
	"context"
	"encoding/json"
	"errors"
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
	var cached Bootstrap
	if c.readCache(ctx, keyBootstrap, &cached) {
		return &cached, nil
	}
	b, err := c.http.Bootstrap(ctx)
	if err != nil {
		return nil, err
	}
	c.writeCache(ctx, keyBootstrap, b, ttlBootstrap)
	return b, nil
}

func (c *CachedClient) Fixtures(ctx context.Context, gw int) ([]Fixture, error) {
	key := fmt.Sprintf(keyFixtures, gw)
	var cached []Fixture
	if c.readCache(ctx, key, &cached) {
		return cached, nil
	}
	fixtures, err := c.http.Fixtures(ctx, gw)
	if err != nil {
		return nil, err
	}
	c.writeCache(ctx, key, fixtures, gameweekTTL(fixtures))
	return fixtures, nil
}

func (c *CachedClient) Live(ctx context.Context, gw int) (*LiveResponse, error) {
	key := fmt.Sprintf(keyLive, gw)
	var cached LiveResponse
	if c.readCache(ctx, key, &cached) {
		return &cached, nil
	}
	live, err := c.http.Live(ctx, gw)
	if err != nil {
		return nil, err
	}
	// Use fixture status to determine TTL; Fixtures(gw) will be a cache hit here.
	ttl := ttlActive
	if fixtures, err := c.Fixtures(ctx, gw); err == nil {
		ttl = gameweekTTL(fixtures)
	} else {
		slog.WarnContext(ctx, "could not load fixtures to determine live cache ttl, using short ttl", "gw", gw, "err", err)
	}
	c.writeCache(ctx, key, live, ttl)
	return live, nil
}

func (c *CachedClient) Entry(ctx context.Context, id string) (*Entry, error) {
	key := fmt.Sprintf(keyEntry, id)
	var cached Entry
	if c.readCache(ctx, key, &cached) {
		return &cached, nil
	}
	e, err := c.http.Entry(ctx, id)
	if err != nil {
		return nil, err
	}
	c.writeCache(ctx, key, e, ttlEntry)
	return e, nil
}

func (c *CachedClient) Picks(ctx context.Context, entryID string, gw int) (*Picks, error) {
	key := fmt.Sprintf(keyPicks, entryID, gw)
	var cached Picks
	if c.readCache(ctx, key, &cached) {
		return &cached, nil
	}
	p, err := c.http.Picks(ctx, entryID, gw)
	if err != nil {
		return nil, err
	}
	bootstrap, bootstrapErr := c.Bootstrap(ctx)
	if bootstrapErr != nil {
		slog.WarnContext(ctx, "could not load bootstrap to determine picks cache ttl, using short ttl", "gw", gw, "err", bootstrapErr)
	}
	c.writeCache(ctx, key, p, picksTTL(bootstrap, bootstrapErr, gw))
	return p, nil
}

// readCache loads key from Redis into v. It returns false on a cache miss or
// any error, in which case the caller should fetch from upstream.
func (c *CachedClient) readCache(ctx context.Context, key string, v any) bool {
	val, err := c.cache.Get(ctx, key).Bytes()
	if err != nil {
		if !errors.Is(err, redis.Nil) {
			slog.ErrorContext(ctx, "cache read failed, re-fetching", "key", key, "err", err)
		}
		return false
	}
	if err := json.Unmarshal(val, v); err != nil {
		slog.WarnContext(ctx, "cache unmarshal failed, re-fetching", "key", key, "err", err)
		return false
	}
	return true
}

// writeCache stores v in Redis under key. Failures are logged but not returned,
// since the caller already has the data it needs.
func (c *CachedClient) writeCache(ctx context.Context, key string, v any, ttl time.Duration) {
	data, err := json.Marshal(v)
	if err != nil {
		slog.ErrorContext(ctx, "cache marshal failed, not caching", "key", key, "err", err)
		return
	}
	if err := c.cache.Set(ctx, key, data, ttl).Err(); err != nil {
		slog.ErrorContext(ctx, "cache write failed", "key", key, "ttl", ttl, "err", err)
	}
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
	var cached LeagueStandings
	if c.readCache(ctx, key, &cached) {
		return &cached, nil
	}
	s, err := c.http.LeagueStandings(ctx, leagueID, page)
	if err != nil {
		return nil, err
	}
	c.writeCache(ctx, key, s, ttlLeagueStandings)
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
