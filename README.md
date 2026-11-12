# Crewly

Direct hiring marketplace for India's creative workforce. Clients post requests, talent applies, clients shortlist through a swipe deck, and matches unlock chat, contracts, and protected escrow payments.

---

## Tech Stack

- **Language:** Go 1.24+
- **HTTP Router:** Chi (`net/http` standard library compatibility)
- **Database:** PostgreSQL 16
- **DB Driver:** `pgx/v5` + `pgxpool`
- **Query Generation:** `sqlc` (type-safe SQL, zero ORM overhead)
- **Migrations:** `goose`
- **API Spec:** OpenAPI 3.0 (`backend/api/openapi.yaml`)
- **Frontend (planned):** React Native / Expo

---

## Project Structure

```
crewly/
├── Makefile                      # Common dev commands
├── docker-compose.yml            # PostgreSQL 16 dev container
├── checklist.md                  # Milestone & weekly build checklist
├── design.md                     # System architecture & core flows
├── data-diagram.md               # Database ER diagram & queries
├── backend/
│   ├── cmd/api/main.go           # Wiring & HTTP server entrypoint
│   ├── internal/
│   │   ├── config/               # Environment & configuration loader
│   │   ├── platform/             # Infrastructure (db pool, httpx, logger)
│   │   ├── auth/                 # OTP generation, verification, JWT
│   │   ├── profile/              # Talent & client profiles
│   │   ├── job/                  # Job postings & filters
│   │   ├── match/                # Swipe deck & shortlist matching
│   │   ├── chat/                 # WebSocket & messaging
│   │   └── payment/              # Escrow deposit & payouts
│   ├── migrations/               # Goose SQL migrations (0001_init.sql)
│   ├── queries/                  # sqlc SQL queries
│   └── api/                      # OpenAPI 3.0 specification
```

---

## Developer Quickstart

### Prerequisites
- [Go](https://go.dev/) (1.24+)
- [Docker](https://www.docker.com/) & Docker Compose
- `make` (optional; commands can be run via `go` / `docker` directly)

### Setup & Run

1. **Clone repository & copy environment variables:**
   ```bash
   cp .env.example .env
   ```

2. **Start PostgreSQL:**
   ```bash
   make db-up
   # Or: docker compose up -d postgres
   ```

3. **Run database migrations:**
   ```bash
   make migrate
   ```

4. **Start API server:**
   ```bash
   make run
   # Or: go run ./backend/cmd/api
   ```

5. **Verify health check:**
   ```bash
   curl -i http://localhost:8080/health
   ```

6. **Run tests:**
   ```bash
   make test
   # Or: go test -v ./backend/...
   ```

7. **Generate typed queries (sqlc):**
   ```bash
   make sqlc
   # Or: sqlc generate
   ```

---

## Environment Variables

| Variable | Default | Description |
|---|---|---|
| `PORT` | `8080` | Port for the HTTP API server |
| `ENV` | `development` | Environment (`development`, `production`) |
| `DATABASE_URL` | `postgres://crewly:crewly_secret@localhost:5432/crewly?sslmode=disable` | PostgreSQL connection string |
| `JWT_SECRET` | `dev_jwt_secret_crewly_market_2026_super_secure` | Secret key for signing JWT tokens |

---

## Design & Documentation

- [checklist.md](checklist.md): Step-by-step phased build plan and done criteria.
- [design.md](design.md): Architectural rules, layer conventions, and core mechanics.
- [data-diagram.md](data-diagram.md): Complete schema ERD, constraints, and lifecycle diagrams.
- [backend/api/openapi.yaml](backend/api/openapi.yaml): OpenAPI contract for client consumption.
