# Deployment Plan

## Verdict

Recommended MVP split:

- `apps/web` -> Vercel Hobby
- `Auth + Postgres + Storage` -> Supabase Free
- `services/api-go` -> Koyeb Starter
- `services/crawler` -> Koyeb Starter only for manual/on-demand runs, or move to GitHub Actions cron if free-instance hours become tight

Rejected for the free-tier plan:

- Heroku: not suitable for this repo's free deployment path because current pricing is paid-first and does not provide a practical always-free app runtime
- Fly.io: technically workable, but not the best default for a zero-budget MVP because the current model is credits/allowance oriented rather than a clear always-free service shape

## Why This Split Fits

### Web on Vercel Hobby

- Good fit for Next.js App Router, preview deployments, and image/static delivery
- Clean monorepo support through `root_directory = "apps/web"`
- Best place for the public web surface and admin UI

### Supabase Free for Auth + DB + Storage

- Good fit for Google OAuth, Postgres, object storage, and row-level security
- Reduces the amount of stateful infrastructure we have to self-host
- Good enough for MVP and internal/admin testing

Caveats:

- Free-tier capacity is enough for MVP, not for sustained crawler-heavy production load
- Inactivity / project pausing behavior must be expected in dev and low-traffic environments

### Koyeb Starter for `api-go`

- Fits the Echo API shape well
- Works with Docker image deploys and env-driven config
- Simpler than running a whole VPS for one API service

### Koyeb for `crawler`

- Acceptable only if the crawler is not a 24/7 always-on worker on free capacity
- Better as:
  - manual admin trigger target
  - scheduled lightweight pull
  - or replaced by GitHub Actions cron for the polling loop

For this repo, I would not treat free Koyeb as a durable "always-on API plus always-on crawler" platform forever. It is suitable for MVP iteration, not a stable long-term zero-cost production backend.

## Updated Rollout Plan

### Phase 1

- Deploy `apps/web` to Vercel Hobby
- Set `ADMIN_EMAILS` on the web deployment only as a temporary/fallback admin gate
- Create Supabase project and wire:
  - Auth
  - Postgres
  - Storage
- Deploy `services/api-go` to Koyeb Starter
- Run `make grant-role` against the production database to create the first DB-backed admin role
- Keep crawler fixture/manual in dev until API + DB env is stable

### Phase 2

- Move crawler to scheduled remote execution
- Prefer GitHub Actions cron first for free predictable scheduling
- Keep Koyeb crawler service only if the schedule/runtime profile still fits free limits
- Scaffold added: `.github/workflows/crawler-schedule.yml`

### Phase 3

- Move push matching and ingestion heavy paths off the request path
- Add proper worker separation
- Revisit OCI Always Free or a paid container platform if sustained traffic starts to exceed free limits

## Rust Split Evaluation

There are parts of the current Go/Python stack that are reasonable Rust candidates, but not all of them are worth splitting.

### Strong Rust Candidates

1. `matcher-rs`
- Purpose: keyword matching, regulated-item tagging, change diffing, notification fan-out preparation
- Why Rust fits: high-throughput string matching, lower memory overhead, strong safety for untrusted ingested text
- Benefit: keeps the API request path smaller and moves CPU-bound matching into a safer service

2. `policy-rs`
- Purpose: auction warning/tag derivation, compliance rules, file/content checksum validation
- Why Rust fits: deterministic logic, easy to harden, low runtime footprint
- Benefit: safer and easier to fuzz than embedding all rule logic in the API server

3. `ics-rs` or `feed-rs`
- Purpose: signed ICS/calendar feed generation and token validation
- Why Rust fits: tiny attack surface, straightforward binary, easy to isolate
- Benefit: stronger isolation for a public unauthenticated-ish endpoint

### Not Worth Splitting Yet

1. Main HTTP CRUD API
- Current Go Echo API is already efficient enough
- Rewriting it in Rust now would slow product delivery more than it improves runtime

2. HTML crawler
- The current crawler depends on Python/Playwright
- Replacing it with Rust now would increase complexity and reduce iteration speed

## Rust Plan Addition

Add these optional phase items after the deployment baseline is stable:

1. `services/matcher-rs`
- Input: normalized auction rows or DB change events
- Output: matched subscription ids, derived warning tags, notification jobs

2. `services/policy-rs`
- Input: crawler payloads
- Output: validated/normalized policy metadata and checksum decisions

3. Keep `api-go` as orchestration and CRUD surface
- Go remains the control plane
- Rust handles CPU-bound or security-sensitive paths

## Source Notes

Current suitability is based on the providers' official pricing / docs pages as checked on 2026-03-11:

- Vercel pricing: https://vercel.com/pricing
- Supabase pricing: https://supabase.com/pricing
- Koyeb pricing: https://www.koyeb.com/pricing
- Koyeb Terraform announcement/examples: https://www.koyeb.com/blog/terraform-provider-for-managing-applications-on-koyeb
- Fly.io pricing docs: https://fly.io/docs/about/pricing/
- Heroku pricing: https://www.heroku.com/pricing
