CREWLY

Talent Marketplace: Problem Statement, Market Scan and Build Plan
Prepared for Chida | November 2026
1. Problem statement
[App name] is a direct hiring marketplace for India's creative workforce. Anyone with a job (a brand, photographer, event planner or individual) posts a request. Talent applies, or is invited. The client shortlists through a swipe-style interface, and matches unlock chat, contracts and protected payments. We replace the manager's commission, not the manager's trust.
• Launch: models, India only
• Scale: photographers, videographers, makeup artists, crew and support staff
• Positioning: a project that is both a resume piece and a startup attempt
2. Does something like this already exist?
Yes, the market is not empty, but the angle is still open.
Platform
What it does
Verdict
Dazzlerr
Indian casting and networking platform connecting artists directly with casting directors, agencies, production houses and brands
Closest competitor. Listing-board model, no swipe UX found
CastYou
Free sign-up, 0% commission, no registration fee
Sets the price anchor: talent expects free
Promodeling
Free app where models, photographers and makeup artists find shoots and network
Global, not India-focused
MODEL IN
Connects e-commerce stores, brands, photographers and creators with models
Same niche as the launch segment
ModelManagement.com
Founded 2017, links models, influencers, photographers and brands, with $5.0M funding
Proof investors fund this space
None of these surfaced as swipe-first or escrow-backed, but the search was shallow. Check the app stores directly before pitching "no competitors."
The real problem: photographers still describe sourcing as DMing hundreds of models on Instagram, getting few replies, and chasing follow-ups. That is the wedge.
Differentiators:
1. Swipe-based shortlisting
2. Payment protection and contracts (what managers actually provide)
3. One app for all talent types, later
3. Design decisions, with verdicts
Swipe mechanics: two modes, ship one first
1. Apply mode (launch): talent applies, so intent is confirmed. The client swipes on applicants.
2. Invite mode (phase 2): the client swipes talent who have not applied, filtered by criteria. A match happens only if the talent accepts.
3. Verdict: applying means willing. Build Apply mode only first, since Invite mode doubles the complexity.
Cards and filters
• Card: 1 photo, 1-line bio, city, rate range, verified badge
• Full portfolio and reviews unlock after the match
• Filters: city, budget, height/age range, experience, availability date
Payments: deposit-escrow hybrid
1. Client pays a 20-30% deposit into the platform at match time.
2. The balance is released after the shoot is confirmed.
3. Verdict: the deposit stops no-shows and gives the platform leverage. Use an Indian payment gateway with marketplace/split-payout support (check Razorpay's current offering).
Off-platform dealing: outcompete it, don't fight it
1. Do not block phone numbers (it is unenforceable).
2. Make the platform worth staying on: payment protection, contracts, and reviews that count only for in-platform bookings.
3. Unlock chat at match, but share phone numbers only after the deposit.
Verification and disputes: the moat
1. Talent: phone OTP, selfie match against portfolio photos, linked Instagram.
2. Clients: verify too. Fake casting calls are a known risk for models.
3. No-show: deposit forfeited plus a strike. 3 strikes means removal.
4. Disputes: 48-hour window, evidence uploaded in-app, manual review at launch.
Monetization: free for talent, paid by clients
• ~5% service fee on escrowed bookings, paid by the client. CastYou's 0% anchor makes charging talent a hard sell.
• Later: paid job boosts and featured profiles.
4. Suggested plan
Phase
Timeline
Deliverable
1. MVP
Weeks 1-6
Auth, profiles, job posts, apply, swipe, match, chat
2. Trust
Weeks 7-10
Verification, deposit escrow, reviews
3. Pilot
Weeks 11-14
30-50 models and 10 clients in one city
4. Expand
Month 4+
Invite mode, photographers, more cities
Stack verdict: React Native (Expo) + Go + PostgreSQL. Swipe UX needs mobile-first, and Go handles matching and payment logic cleanly. Skip microservices: a single monolith is enough for the MVP.
Resume angle: the escrow flow, verification pipeline and swipe matching make strong, specific bullets.
5. Open questions
1. First client segment? Suggested: e-commerce and D2C sellers and small photographers, since they hire models constantly.
2. Pilot city? Bengaluru is the obvious choice.
3. How to get the first 30 models? Instagram outreach, college fashion clubs, or a local photographer partner.
4. Minors? It changes legal and consent requirements significantly. Recommendation: 18+ only at launch.
5. Solo or team? It decides whether the 14-week plan is realistic.
6. Next step
Answer the open questions, then share the README so it can be updated with everything above.

---

## 7. Developer Quickstart

### Prerequisites
- Go 1.24+
- Docker & Docker Compose
- Make (optional, direct commands also supported)

### Running Locally
1. **Start PostgreSQL database:**
   ```bash
   make db-up
   # Or: docker compose up -d postgres
   ```

2. **Run migrations:**
   ```bash
   make migrate
   ```

3. **Start the API server:**
   ```bash
   make run
   # Or: go run ./backend/cmd/api
   ```

4. **Verify health endpoint:**
   ```bash
   curl -i http://localhost:8080/health
   ```

5. **Run test suite:**
   ```bash
   make test
   # Or: go test -v ./backend/...
   ```
