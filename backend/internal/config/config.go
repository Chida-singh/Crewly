package config

import (
	"bufio"
	"os"
	"strings"
)

type Config struct {
	Port        string
	Env         string
	DatabaseURL   string
	JWTSecret     string
	PublicBaseURL string
}

// Load reads config from environment variables, optionally reading a .env file if present.
func Load() (*Config, error) {
	// Simple .env parser to avoid external dependencies for basic dev loading
	loadDotEnv(".env")

	cfg := &Config{
		Port:        getEnv("PORT", "8080"),
		Env:         getEnv("ENV", "development"),
		DatabaseURL:   getEnv("DATABASE_URL", "postgres://crewly:crewly_secret@localhost:5432/crewly?sslmode=disable"),
		JWTSecret:     getEnv("JWT_SECRET", "dev_jwt_secret_crewly_market_2026_super_secure"),
		PublicBaseURL: getEnv("PUBLIC_BASE_URL", "http://localhost:8080"),
	}

	return cfg, nil
}

func getEnv(key, defaultVal string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultVal
}

func loadDotEnv(filepath string) {
	file, err := os.Open(filepath)
	if err != nil {
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			k := strings.TrimSpace(parts[0])
			v := strings.TrimSpace(parts[1])
			// Only set if not already set in environment
			if os.Getenv(k) == "" {
				os.Setenv(k, v)
			}
		}
	}
}

