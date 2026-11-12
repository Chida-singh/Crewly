package db_test

import (
	"context"
	"os"
	"testing"
	"time"

	"crewly/backend/internal/platform/db"
)

func TestDatabasePing(t *testing.T) {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = "postgres://crewly:crewly_secret@localhost:5432/crewly?sslmode=disable"
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	pool, err := db.NewPool(ctx, dbURL)
	if err != nil {
		t.Skipf("skipping live database test: postgres not reachable at %s (%v)", dbURL, err)
		return
	}
	defer pool.Close()

	var result int
	if err := pool.QueryRow(ctx, "SELECT 1").Scan(&result); err != nil {
		t.Fatalf("failed to query database: %v", err)
	}

	if result != 1 {
		t.Errorf("expected 1 from SELECT 1, got %d", result)
	}
}

