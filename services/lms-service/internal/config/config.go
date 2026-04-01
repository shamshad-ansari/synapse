package config

import (
	"encoding/base64"
	"os"
	"strconv"
	"strings"
)

type Config struct {
	AppEnv             string
	HTTPPort           string
	DatabaseURL        string
	JWTSecret          string
	CORSAllowedOrigins string
	RedisURL           string
	CanvasClientID     string
	CanvasClientSecret string
	EncryptionKey      []byte
	FrontendURL        string
	CanvasRedirectURI  string
	CanvasInternalURL  string
	SyncIntervalMinutes int
	SyncStartupDelaySec int
}

// Load reads configuration from environment variables with sensible defaults.
func Load() Config {
	keyB64 := strings.TrimSpace(getEnv("ENCRYPTION_KEY", ""))
	keyB64 = strings.Trim(keyB64, `"'`)
	var key []byte
	if keyB64 != "" {
		var err error
		key, err = base64.StdEncoding.DecodeString(keyB64)
		if err != nil {
			if len(keyB64) == 32 {
				// Dev convenience fallback: allow a raw 32-byte key.
				key = []byte(keyB64)
			} else {
				panic("ENCRYPTION_KEY must be valid base64 (or a raw 32-byte string): " + err.Error())
			}
		}
	}

	return Config{
		AppEnv:             getEnv("APP_ENV", "dev"),
		HTTPPort:           getEnv("HTTP_PORT", "8081"),
		DatabaseURL:        getEnv("DATABASE_URL", ""),
		JWTSecret:          getEnv("JWT_SECRET", ""),
		CORSAllowedOrigins: getEnv("CORS_ALLOWED_ORIGINS", "http://localhost:4200"),
		RedisURL:           getEnv("REDIS_URL", "redis://localhost:6379"),
		CanvasClientID:     getEnv("CANVAS_CLIENT_ID", ""),
		CanvasClientSecret: getEnv("CANVAS_CLIENT_SECRET", ""),
		EncryptionKey:      key,
		FrontendURL:        getEnv("FRONTEND_URL", "http://localhost:4200"),
		CanvasRedirectURI:  getEnv("CANVAS_REDIRECT_URI", "http://localhost:8081/v1/lms/callback/canvas"),
		CanvasInternalURL:  getEnv("CANVAS_INTERNAL_URL", ""),
		SyncIntervalMinutes: getEnvInt("LMS_SYNC_INTERVAL_MINUTES", 15),
		SyncStartupDelaySec: getEnvInt("LMS_SYNC_STARTUP_DELAY_SEC", 20),
	}
}

// Validate panics if any required field is missing or invalid.
func (c Config) Validate() {
	if c.DatabaseURL == "" {
		panic("missing required env var: DATABASE_URL")
	}
	if c.JWTSecret == "" {
		panic("missing required env var: JWT_SECRET")
	}
	if c.CanvasClientID == "" {
		panic("missing required env var: CANVAS_CLIENT_ID")
	}
	if c.CanvasClientSecret == "" {
		panic("missing required env var: CANVAS_CLIENT_SECRET")
	}
	if len(c.EncryptionKey) != 32 {
		panic("ENCRYPTION_KEY must decode to exactly 32 bytes")
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func getEnvInt(key string, fallback int) int {
	raw := strings.TrimSpace(getEnv(key, ""))
	if raw == "" {
		return fallback
	}
	n, err := strconv.Atoi(raw)
	if err != nil || n <= 0 {
		return fallback
	}
	return n
}
