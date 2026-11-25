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
- **Frontend:** React 19 + TypeScript + Vite (`frontend/`)

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
└── frontend/                     # React 19 + TypeScript + Vite app
```

---

## Quickstart

### What you need
- **Frontend only:** Bun (or npm). No Docker, no database.
- **Full stack:** Go 1.24+, Docker & Docker Compose for Postgres. `make` is optional.

### 1. Just the app — fastest way to see it

```bash
cd frontend
bun install
bun run dev
```

Open <http://127.0.0.1:5173>. Sign-in is switched off for local development and the sheet starts
pre-filled with placeholder data, so the whole app is clickable without any backend.

The demo sheet is memory-only: edits and uploaded frames reset when you refresh. To bring the
real OTP ticket back, set `VITE_AUTH_BYPASS=false` in `frontend/.env.development` and restart.

### 2. Full stack — API + Postgres

You do **not** need `docker compose build`: the compose file only runs Postgres, and the API runs
straight from Go.

1. **Environment:**
   ```bash
   cp .env.example .env
   ```

2. **Database:**
   ```bash
   docker compose up -d postgres
   ```

3. **Migrations:**
   ```bash
   go run github.com/pressly/goose/v3/cmd/goose -dir backend/migrations postgres \
     "postgres://crewly:crewly_secret@localhost:5432/crewly?sslmode=disable" up
   ```

4. **API:**
   ```bash
   go run ./backend/cmd/api
   ```

5. **Check it's up:**
   ```bash
   curl -i http://localhost:8080/health
   ```
   Expect `{"status":"ok","database":"up"}`.

With Postgres and the API running, set `VITE_AUTH_BYPASS=false` in `frontend/.env.development`
and restart `bun run dev` to test the real phone-code login. The 6-digit code is printed in the
API log (dev SMS stub).

`make db-up`, `make migrate` and `make run` do the same three steps if `make` is installed.

### Tests and code generation

```bash
go test ./backend/...                                      # backend tests
go run github.com/sqlc-dev/sqlc/cmd/sqlc generate          # after editing backend/queries/
cd frontend && bun run build                              # typecheck + production build
```

---

## Environment Variables

| Variable | Default | Description |
|---|---|---|
| `PORT` | `8080` | Port for the HTTP API server |
| `ENV` | `development` | Environment (`development`, `production`) |
| `DATABASE_URL` | `postgres://crewly:crewly_secret@localhost:5432/crewly?sslmode=disable` | PostgreSQL connection string |
| `JWT_SECRET` | `dev_jwt_secret_crewly_market_2026_super_secure` | Secret key for signing JWT tokens |
| `PUBLIC_BASE_URL` | `http://localhost:8080` | Public origin used for dev photo upload URLs |

---

## Design & Documentation

- [checklist.md](checklist.md): Step-by-step phased build plan and done criteria.
- [design.md](design.md): Architectural rules, layer conventions, and core mechanics.
- [data-diagram.md](data-diagram.md): Complete schema ERD, constraints, and lifecycle diagrams.
- [backend/api/openapi.yaml](backend/api/openapi.yaml): OpenAPI contract for client consumption.
