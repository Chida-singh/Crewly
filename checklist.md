# Crewly: Build Checklist

Backend-first MVP. Stack: Go (chi, pgx, sqlc, goose) + PostgreSQL. Frontend (Expo) comes later.
Read `design.md` for architecture and `data-diagram.md` for the schema before starting any task.

**Rule for the agent:** work one week at a time. Do not start a later week until every box in the current week is checked and its "Done when" test passes.

---

## Phase 1: MVP (Weeks 1-6)

### Week 1: Scaffold
- [x] `git init`, add `README.md`, `.gitignore`, `Makefile`
- [x] Create folder layout from `design.md` (`backend/cmd/api`, `backend/internal/...`)
- [x] `docker-compose.yml` with PostgreSQL 16
- [x] `go mod init`, add chi, pgx, goose, sqlc
- [x] Copy `0001_init.sql` into `backend/migrations/` and run `make migrate`
- [x] Config loader (`internal/config`) reading env vars
- [x] chi server in `cmd/api/main.go` with `GET /health` (also pings DB)
- [x] `api/openapi.yaml` skeleton with `/health`
- [x] First test that hits the database
- **Done when:** `make run` serves `/health` and the migration applies cleanly on a fresh database.

### Week 2: Auth and profiles
- [ ] `POST /auth/otp`: generate OTP, store hashed with expiry, rate limit per phone
- [ ] SMS client interface with a **dev stub that logs the OTP**
- [ ] `POST /auth/verify`: check OTP, create user if new, return JWT
- [ ] JWT middleware that puts `user_id` and `role` in the request context
- [ ] `GET/PUT /me/profile` for talent and client roles
- [ ] `POST /me/photos`: presigned upload URL, save URL in `profile_photos`
- [ ] Enforce 18+ at profile creation (DB check already exists)
- [ ] Unit tests for the auth service
- **Done when:** you can sign up as talent and as client using curl only.

### Week 3: Jobs and applications
- [ ] `POST /jobs` (client only)
- [ ] `GET /jobs` with filters: city, budget, date, talent type; cursor pagination
- [ ] `GET /jobs/{id}`
- [ ] `POST /jobs/{id}/apply` (talent only, job must be `open`)
- [ ] Duplicate apply returns a clear 409 (`UNIQUE (job_id, talent_id)`)
- [ ] Integration tests for the apply flow
- **Done when:** talent applies once and a second attempt is rejected.

### Week 4: Swipe and match
- [ ] `GET /jobs/{id}/applicants`: swipe deck, only status `applied`, cursor paginated, owner-only
- [ ] Card payload: 1 photo, 1-line bio, city, rate range, verified flag
- [ ] `POST /applications/{id}/decision` with `shortlist` or `reject`
- [ ] Decision runs in **one transaction**: update with `WHERE status='applied'`, then insert into `matches` on shortlist
- [ ] Double-tap safety: repeat decision returns the existing result, not a second match
- [ ] `GET /matches` for both roles
- [ ] Full profile and portfolio unlock only after a match
- **Done when:** one right swipe creates exactly **1** match, even under concurrent requests.

### Week 5: Chat
- [ ] WebSocket hub in `internal/chat` (`WS /ws`, JWT auth on connect)
- [ ] Persist every message to `messages`
- [ ] `GET /matches/{id}/messages` with pagination
- [ ] Only the two match participants can read or send
- [ ] Phone numbers are never exposed in chat payloads
- **Done when:** two clients exchange messages in real time and history loads after reconnect.

### Week 6: Harden and ship
- [ ] Integration test suite covering signup to match to chat
- [ ] `make seed` script: 10 talent, 3 clients, 5 jobs
- [ ] Rate limiting middleware on auth and apply endpoints
- [ ] Structured logging with request IDs
- [ ] Postman collection exported from `openapi.yaml`
- [ ] First deploy (Fly.io or Railway) with managed Postgres
- [ ] README updated with run instructions
- **Done when:** a fresh clone runs end to end in under **10 minutes**.

---

## Phase 2: Trust (Weeks 7-10)
- [ ] Talent verification: OTP (done), selfie match against portfolio photos, linked Instagram
- [ ] Client verification to block fake casting calls
- [ ] Razorpay integration: deposit of **20-30%** at match time
- [ ] Payment states: `pending`, `held`, `released`, `refunded`, `failed`
- [ ] Release balance after shoot confirmation
- [ ] Reviews that count only for in-platform bookings
- [ ] No-show handling: deposit forfeited plus a strike, **3 strikes** means removal
- [ ] Dispute flow: 48-hour window, evidence upload, manual review
- [ ] Phone numbers shared only after the deposit

## Phase 3: Pilot (Weeks 11-14)
- [ ] One city (Bengaluru)
- [ ] **30-50 models** onboarded
- [ ] **10 clients** posting real jobs
- [ ] Track: applications per job, match rate, time to first match, deposit completion

## Phase 4: Expand (Month 4+)
- [ ] Invite mode (client swipes talent who have not applied; talent must accept)
- [ ] More talent types: photographers, videographers, makeup artists, crew
- [ ] More cities
- [ ] Paid job boosts and featured profiles
- [ ] Mobile app (Expo) against the stable `openapi.yaml`

---

## Open decisions
- [ ] First client segment (suggested: D2C and e-commerce sellers, small photographers)
- [ ] Pilot city (suggested: Bengaluru)
- [ ] How to source the first 30 models
- [ ] Platform fee: **~5%** client-side service fee
- [ ] Solo or team (decides if the 14-week plan holds)
