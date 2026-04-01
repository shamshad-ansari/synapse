package config

import (
	"os"
	"strconv"
	"strings"
)

type Config struct {
	AppEnv               string
	HTTPPort             string
	DatabaseURL          string
	JWTSecret            string
	JWTExpiryHours       int
	JWTRefreshExpiryDays int
	CORSAllowedOrigins   string
	RedisURL             string
	AnthropicAPIKey      string
	OpenAIAPIKey         string
	// AIGenerationProvider: anthropic | gemini
	AIGenerationProvider string
	// AIEmbeddingProvider: openai | gemini (vectors must remain 1536-dim for pgvector)
	AIEmbeddingProvider string
	GeminiAPIKey        string
	GeminiModel         string
	GeminiEmbeddingModel string
}

func Load() Config {
	return Config{
		AppEnv:               getEnv("APP_ENV", "dev"),
		HTTPPort:             getEnv("HTTP_PORT", "8080"),
		DatabaseURL:          getEnv("DATABASE_URL", ""),
		JWTSecret:            getEnv("JWT_SECRET", ""),
		JWTExpiryHours:       getEnvInt("JWT_EXPIRY_HOURS", 24),
		JWTRefreshExpiryDays: getEnvInt("JWT_REFRESH_EXPIRY_DAYS", 7),
		CORSAllowedOrigins:   getEnv("CORS_ALLOWED_ORIGINS", "http://localhost:4200,http://127.0.0.1:4200"),
		RedisURL:             getEnv("REDIS_URL", "redis://localhost:6379"),
		AnthropicAPIKey:      getEnv("ANTHROPIC_API_KEY", ""),
		OpenAIAPIKey:         getEnv("OPENAI_API_KEY", ""),
		AIGenerationProvider: strings.ToLower(strings.TrimSpace(getEnv("AI_GENERATION_PROVIDER", "anthropic"))),
		AIEmbeddingProvider:  strings.ToLower(strings.TrimSpace(getEnv("AI_EMBEDDING_PROVIDER", "openai"))),
		GeminiAPIKey:         getEnv("GEMINI_API_KEY", ""),
		GeminiModel:          getEnv("GEMINI_MODEL", ""),
		GeminiEmbeddingModel: getEnv("GEMINI_EMBEDDING_MODEL", ""),
	}
}

// Validate panics if any required field is missing.
func (c Config) Validate() {
	if c.DatabaseURL == "" {
		panic("missing required env var: DATABASE_URL")
	}
	if c.JWTSecret == "" {
		panic("missing required env var: JWT_SECRET")
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func getEnvInt(key string, fallback int) int {
	v := os.Getenv(key)
	if v == "" {
		return fallback
	}
	n, err := strconv.Atoi(v)
	if err != nil {
		return fallback
	}
	return n
}
