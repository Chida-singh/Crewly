package config_test

import (
	"os"
	"testing"

	"crewly/backend/internal/config"
)

func TestConfigLoadDefaults(t *testing.T) {
	os.Unsetenv("PORT")
	os.Unsetenv("ENV")
	os.Unsetenv("DATABASE_URL")
	os.Unsetenv("JWT_SECRET")

	cfg, err := config.Load()
	if err != nil {
		t.Fatalf("unexpected error loading config: %v", err)
	}

	if cfg.Port != "8080" {
		t.Errorf("expected default port 8080, got %s", cfg.Port)
	}
	if cfg.Env != "development" {
		t.Errorf("expected default env development, got %s", cfg.Env)
	}
	if cfg.DatabaseURL == "" {
		t.Error("expected default DatabaseURL to not be empty")
	}
	if cfg.JWTSecret == "" {
		t.Error("expected default JWTSecret to not be empty")
	}
}

func TestConfigLoadCustomEnv(t *testing.T) {
	t.Setenv("PORT", "9090")
	t.Setenv("ENV", "production")
	t.Setenv("DATABASE_URL", "postgres://custom:pass@db:5432/customdb")
	t.Setenv("JWT_SECRET", "super-secret-key")

	cfg, err := config.Load()
	if err != nil {
		t.Fatalf("unexpected error loading config: %v", err)
	}

	if cfg.Port != "9090" {
		t.Errorf("expected port 9090, got %s", cfg.Port)
	}
	if cfg.Env != "production" {
		t.Errorf("expected env production, got %s", cfg.Env)
	}
	if cfg.DatabaseURL != "postgres://custom:pass@db:5432/customdb" {
		t.Errorf("expected custom database url, got %s", cfg.DatabaseURL)
	}
	if cfg.JWTSecret != "super-secret-key" {
		t.Errorf("expected custom jwt secret, got %s", cfg.JWTSecret)
	}
}

