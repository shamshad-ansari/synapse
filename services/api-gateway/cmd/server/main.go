package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/shamshad-ansari/synapse/services/api-gateway/internal/config"
	"github.com/shamshad-ansari/synapse/services/api-gateway/internal/db"
	"github.com/shamshad-ansari/synapse/services/api-gateway/internal/http/router"
)

func main() {
	cfg := config.MustLoad()

	zerolog.TimeFieldFormat = time.RFC3339Nano
	logger := log.Output(zerolog.ConsoleWriter{Out: os.Stdout}).With().
		Str("service", "api-gateway").
		Str("env", cfg.Env).
		Logger()

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	pool, err := db.NewPool(ctx, cfg.DatabaseURL)
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to connect to database")
	}
	defer pool.Close()

	r := router.New(cfg, pool, logger)

	srv := &http.Server{
		Addr:              ":" + cfg.HTTPPort,
		Handler:           r,
		ReadHeaderTimeout: 5 * time.Second,
	}

	go func() {
		logger.Info().Str("port", cfg.HTTPPort).Msg("api-gateway started")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal().Err(err).Msg("server failed")
		}
	}()

	<-ctx.Done()
	logger.Info().Msg("shutdown signal received")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		logger.Error().Err(err).Msg("graceful shutdown failed")
	}
	logger.Info().Msg("server stopped")
}