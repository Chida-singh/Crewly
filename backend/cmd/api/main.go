package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"crewly/backend/internal/config"
	"crewly/backend/internal/platform/db"
	"crewly/backend/internal/platform/httpx"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	slog.SetDefault(logger)

	cfg, err := config.Load()
	if err != nil {
		logger.Error("failed to load configuration", "error", err)
		os.Exit(1)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Initialize database pool
	pool, err := db.NewPool(ctx, cfg.DatabaseURL)
	if err != nil {
		logger.Warn("could not connect to database on startup", "error", err)
	} else {
		logger.Info("connected to postgres database")
		defer pool.Close()
	}

	r := chi.NewRouter()

	// Standard middleware
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))

	// Health check endpoint
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		status := "ok"
		dbStatus := "up"
		statusCode := http.StatusOK

		if pool == nil {
			status = "degraded"
			dbStatus = "down"
			statusCode = http.StatusServiceUnavailable
		} else {
			pingCtx, pingCancel := context.WithTimeout(r.Context(), 2*time.Second)
			defer pingCancel()

			if err := pool.Ping(pingCtx); err != nil {
				status = "degraded"
				dbStatus = "down"
				statusCode = http.StatusServiceUnavailable
			}
		}

		httpx.JSON(w, statusCode, map[string]string{
			"status":   status,
			"database": dbStatus,
		})
	})

	server := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Server shutdown listener
	shutdownErr := make(chan error, 1)
	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
		<-sigChan

		logger.Info("shutting down server...")
		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer shutdownCancel()

		shutdownErr <- server.Shutdown(shutdownCtx)
	}()

	logger.Info(fmt.Sprintf("starting server on port %s (env: %s)", cfg.Port, cfg.Env))
	if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		logger.Error("server error", "error", err)
		os.Exit(1)
	}

	if err := <-shutdownErr; err != nil {
		logger.Error("error during graceful shutdown", "error", err)
	}
	logger.Info("server stopped gracefully")
}

