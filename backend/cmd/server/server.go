package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"

	"github.com/OptimusCrime/fpl-live-match/backend/internal/bootstrap"
	"github.com/OptimusCrime/fpl-live-match/backend/internal/fixture"
	"github.com/OptimusCrime/fpl-live-match/backend/internal/fpl"
	"github.com/OptimusCrime/fpl-live-match/backend/internal/league"
	"github.com/OptimusCrime/fpl-live-match/backend/internal/leaguematch"
	redisclient "github.com/OptimusCrime/fpl-live-match/backend/internal/redis"
	"github.com/OptimusCrime/fpl-live-match/backend/internal/team"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8199"
	}

	logLevel := parseLogLevel(os.Getenv("LOG_LEVEL"))
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: logLevel}))
	slog.SetDefault(logger)
	slog.Info("logger configured", "level", logLevel)

	redisCache, err := redisclient.New()
	if err != nil {
		slog.Error("failed to connect to redis", "err", err)
		os.Exit(1)
	}
	defer redisCache.Close()

	fplClient := fpl.NewCachedClient(redisCache)

	bootstrapService := bootstrap.New(fplClient)
	fixtureService := fixture.NewService(fplClient, bootstrapService)
	leagueService := league.NewService(fplClient, bootstrapService)
	leagueMatchService := leaguematch.NewService(fplClient, bootstrapService, fixtureService)

	fixtureHandler := fixture.NewHandler(fixtureService)
	leagueHandler := league.NewHandler(leagueService)
	teamHandler := team.NewHandler(team.NewService(fplClient))
	leagueMatchHandler := leaguematch.NewHandler(leagueMatchService)

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{"http://localhost:5173", "https://fpl.optimuscrime.net"},
		AllowedMethods: []string{"GET", "OPTIONS"},
		AllowedHeaders: []string{"Accept", "Content-Type", "X-CURRENT-PROFILE-ID"},
	}))

	r.Get("/api/fixtures", fixtureHandler.Fixtures)
	r.Get("/api/team/{id}/leagues", teamHandler.Leagues)
	r.Get("/api/player/{id}/owners", leagueHandler.PlayerOwners)
	r.Get("/api/league/{leagueId}/match/{matchId}", leagueMatchHandler.LeagueMatch)

	srv := &http.Server{
		Addr:         ":" + port,
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 60 * time.Second,
		IdleTimeout:  120 * time.Second,
		ErrorLog:     slog.NewLogLogger(logger.Handler(), slog.LevelError),
	}

	go func() {
		slog.Info("starting server", "port", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("server error", "err", err)
			os.Exit(1)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	slog.Info("shutting down server")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		slog.Error("shutdown error", "err", err)
	}
}

// parseLogLevel maps LOG_LEVEL (debug, info, warn, error) to a slog level, defaulting to info.
func parseLogLevel(v string) slog.Level {
	switch strings.ToLower(v) {
	case "debug":
		return slog.LevelDebug
	case "warn", "warning":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}
