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
- **Web frontend:** React 19 + TypeScript + Vite (`frontend/`)
- **Mobile frontend (planned):** React Native + Expo (`mobile/`), same Go API

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

## Platforms

Crewly is one product with two surfaces. The API is shared; only the client differs.

| Surface | Folder | Status today |
|---|---|---|
| Desktop / laptop browser | `frontend/` | working |
| Phone browser (responsive web) | `frontend/` | working over your Wi-Fi |
| Native app (React Native + Expo) | `mobile/` | planned, not scaffolded |

Everything below runs on your own machine; nothing is deployed.

## Quickstart

### What you need
- **Web only:** Bun (or npm). No Docker, no database.
- **Full stack:** Go 1.24+, Docker & Docker Compose for Postgres. `make` is optional.
- **Native app (later):** Node 20+, the Expo Go app on your phone.

### 1. On a PC — fastest way to see it

```bash
cd frontend
bun install
bun run dev
```

Open <http://127.0.0.1:5173>. Sign-in is switched off for local development and the sheet starts
pre-filled with placeholder data, so the whole app is clickable without any backend.

The demo sheet is memory-only: edits and uploaded frames reset when you refresh. To bring the
real OTP ticket back, set `VITE_AUTH_BYPASS=false` in `frontend/.env.development` and restart.

### 2. On a phone today — same web app over your Wi-Fi

The layout is responsive down to 375px, so you can use the phone browser with no extra build:

```bash
cd frontend
bun run dev -- --host
```

1. Find your computer's LAN IP (`ipconfig` on Windows, `ifconfig` on macOS, `ip addr` on Linux).
2. On the phone, open `http://<that-ip>:5173`.
3. Both devices must be on the same Wi-Fi. On Windows, allow Node through the firewall for
   private networks the first time it asks.

The API, if you are running it, stays on the computer: the dev server proxies to `127.0.0.1:8080`
for you. This is still a browser tab, not an installed app — no app icon or push notifications yet.

### 3. Full stack — API + Postgres

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

### 4. Native app — later, once `mobile/` exists

The Expo client will reuse the same Go API; nothing here is generated yet. When the folder is
scaffolded, the loop will be:

```bash
cd mobile
bun install
npx expo start
```

- Press `a` for an Android emulator, `i` for the iOS simulator (macOS only), or scan the QR code
  with Expo Go on a physical phone.
- Point the app at your API with `EXPO_PUBLIC_API_URL=http://<LAN-IP>:8080` in `mobile/.env`, and
  keep the phone on the same Wi-Fi as the computer.
- Hot reload works the same way as the web app; changes appear on the phone without a rebuild.

Ask for the mobile scaffold when you are ready and this section becomes real commands.

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
