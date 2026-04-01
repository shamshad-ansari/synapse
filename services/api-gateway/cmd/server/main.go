package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"

	"github.com/shamshad-ansari/synapse/services/api-gateway/internal/ai"
	"github.com/shamshad-ansari/synapse/services/api-gateway/internal/config"
	"github.com/shamshad-ansari/synapse/services/api-gateway/internal/repository"
	"github.com/shamshad-ansari/synapse/services/api-gateway/internal/service"
	"github.com/shamshad-ansari/synapse/services/api-gateway/internal/transport/http/handlers"
	router "github.com/shamshad-ansari/synapse/services/api-gateway/internal/transport/http"
)

func main() {
	cfg := config.Load()
	cfg.Validate()

	var (
		logger *zap.Logger
		err    error
	)
	if cfg.AppEnv == "prod" {
		logger, err = zap.NewProduction()
	} else {
		logger, err = zap.NewDevelopment()
	}
	if err != nil {
		panic("failed to init logger: " + err.Error())
	}
	defer logger.Sync() //nolint:errcheck

	// Connect pgxpool with a 5s connect timeout.
	poolCfg, err := pgxpool.ParseConfig(cfg.DatabaseURL)
	if err != nil {
		logger.Fatal("failed to parse DATABASE_URL", zap.Error(err))
	}
	poolCfg.MaxConns = 20
	poolCfg.MinConns = 2

	connectCtx, connectCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer connectCancel()

	pool, err := pgxpool.NewWithConfig(connectCtx, poolCfg)
	if err != nil {
		logger.Fatal("failed to connect to database", zap.Error(err))
	}
	defer pool.Close()

	// Connect Redis.
	redisOpts, err := redis.ParseURL(cfg.RedisURL)
	if err != nil {
		logger.Fatal("failed to parse REDIS_URL", zap.Error(err))
	}
	redisClient := redis.NewClient(redisOpts)
	defer redisClient.Close()

	userRepo := repository.NewPostgresUserRepo(pool)
	learningRepo := repository.NewPostgresLearningRepo(pool)
	authSvc := service.NewAuthService(userRepo, &cfg)
	learningSvc := service.NewLearningService(learningRepo)
	autopilot := &handlers.AutopilotHandler{DB: pool}

	gen := newFlashcardGenerator(&cfg)
	emb := newTextEmbedder(&cfg)
	if gen == nil || emb == nil {
		logger.Warn("AI flashcard generation disabled: set AI_GENERATION_PROVIDER and AI_EMBEDDING_PROVIDER with matching API keys (see .env.example). " +
			"Do not switch embedding providers in production without re-embedding stored vectors.")
	}

	noteAI := newNoteAI(&cfg)

	learningH := &handlers.LearningHandler{
		Service:  learningSvc,
		Repo:     learningRepo,
		Embed:    emb,
		AIClient: noteAI,
	}

	var aiSvc *service.AIService
	if gen != nil && emb != nil {
		aiSvc = service.NewAIService(learningRepo, gen, emb, logger)
	}
	aiHandler := &handlers.AIHandler{AIService: aiSvc}

	r := router.NewRouter(&cfg, pool, logger, authSvc, learningH, autopilot, aiHandler)

	srv := &http.Server{
		Addr:         ":" + cfg.HTTPPort,
		Handler:      r,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	logger.Info("synapse api-gateway starting on :" + cfg.HTTPPort)

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("server error", zap.Error(err))
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	logger.Info("shutdown signal received")

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer shutdownCancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		logger.Error("graceful shutdown failed", zap.Error(err))
	}
	logger.Info("server stopped")
}

func newFlashcardGenerator(cfg *config.Config) ai.FlashcardGenerator {
	switch cfg.AIGenerationProvider {
	case "anthropic":
		if cfg.AnthropicAPIKey == "" {
			return nil
		}
		return ai.NewAnthropicClient(cfg.AnthropicAPIKey)
	case "gemini":
		if cfg.GeminiAPIKey == "" {
			return nil
		}
		return ai.NewGeminiClient(cfg.GeminiAPIKey, cfg.GeminiModel)
	default:
		return nil
	}
}

func newTextEmbedder(cfg *config.Config) ai.TextEmbedder {
	switch cfg.AIEmbeddingProvider {
	case "openai":
		if cfg.OpenAIAPIKey == "" {
			return nil
		}
		return ai.NewOpenAIEmbedClient(cfg.OpenAIAPIKey)
	case "gemini":
		if cfg.GeminiAPIKey == "" {
			return nil
		}
		return ai.NewGeminiEmbedClient(cfg.GeminiAPIKey, cfg.GeminiEmbeddingModel)
	default:
		return nil
	}
}

func newNoteAI(cfg *config.Config) ai.Completer {
	switch cfg.AIGenerationProvider {
	case "anthropic":
		if cfg.AnthropicAPIKey == "" {
			return nil
		}
		return ai.NewAnthropicClient(cfg.AnthropicAPIKey)
	case "gemini":
		if cfg.GeminiAPIKey == "" {
			return nil
		}
		return ai.NewGeminiClient(cfg.GeminiAPIKey, cfg.GeminiModel)
	default:
		return nil
	}
}
