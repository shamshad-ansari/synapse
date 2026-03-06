package config

import (
	"os"
)

type Config struct {
	Env               string
	HTTPPort          string
	DatabaseURL       string
	CORSAllowedOrigin string
}

func MustLoad() Config {
	cfg := Config{
		Env:               getEnv("APP_ENV", "dev"),
		HTTPPort:          getEnv("HTTP_PORT", "8080"),
		DatabaseURL:       mustGetEnv("DATABASE_URL"),
		CORSAllowedOrigin: getEnv("CORS_ALLOWED_ORIGINS", "http://localhost:4200"),
	}
	return cfg
}

func getEnv(key, fallback string) string {
	v := os.Getenv(key)
	if v == "" {
		return fallback
	}
	return v
}

func mustGetEnv(key string) string {
	v := os.Getenv(key)
	if v == "" {
		panic("missing required env var: " + key)
	}
	return v
}