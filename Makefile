.PHONY: help db-up db-down migrate migrate-down sqlc test run

DB_URL ?= postgres://crewly:crewly_secret@localhost:5432/crewly?sslmode=disable

help:
	@echo "Available commands:"
	@echo "  make db-up        - Start PostgreSQL container"
	@echo "  make db-down      - Stop PostgreSQL container"
	@echo "  make migrate      - Run database migrations"
	@echo "  make migrate-down - Roll back latest database migration"
	@echo "  make sqlc         - Generate typed Go code from SQL queries"
	@echo "  make test         - Run all backend tests"
	@echo "  make run          - Run API server"

db-up:
	docker compose up -d postgres

db-down:
	docker compose down

migrate:
	go run github.com/pressly/goose/v3/cmd/goose -dir backend/migrations postgres "$(DB_URL)" up

migrate-down:
	go run github.com/pressly/goose/v3/cmd/goose -dir backend/migrations postgres "$(DB_URL)" down

sqlc:
	go run github.com/sqlc-dev/sqlc/cmd/sqlc generate

test:
	go test -v ./backend/...

run:
	go run ./backend/cmd/api

