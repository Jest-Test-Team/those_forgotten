# Backend Evolution Guide

This guide answers the next-stage architecture questions for the platform:

- where Rust adds real value
- when to introduce Kafka and Redis
- how to harden against DDoS, CSRF, and SQL injection
- when to keep JSON and when to add gRPC, WebSocket, or MQTT

## 1. Rust Service Candidates

Do not rewrite the main Go CRUD API first. The best Rust candidates are the narrow, high-throughput, or security-sensitive paths.

### Best Rust Targets

1. `services/matcher-rs`
- Purpose: keyword matching, regulated-item tagging, change diffing, notification candidate generation
- Why Rust: CPU-bound string processing, predictable memory use, easier fuzzing
- Recommended interfaces:
  - consume Kafka `auction.normalized`
  - emit Kafka `auction.matched`

2. `services/policy-rs`
- Purpose: regulated-item rules, compliance flags, import-license warnings, checksum policy validation
- Why Rust: deterministic rules engine with strong safety guarantees
- Recommended interfaces:
  - gRPC internal service or Kafka consumer/producer

3. `services/push-rs`
- Purpose: Web Push / VAPID delivery worker
- Why Rust: good async I/O, easy concurrency control, lower overhead for large fan-out
- Recommended interfaces:
  - consume Kafka `notification.jobs`
  - update delivery status via gRPC or direct DB write

4. `services/feed-rs`
- Purpose: ICS generation, signed feed validation, public feed serving
- Why Rust: small attack surface, easy to isolate, straightforward binary

### Not Worth Rewriting Yet

1. `services/api-go`
- Go Echo is already sufficient for the current request profile.
- Rewriting the control plane now adds delivery risk without material performance gain.

2. `services/crawler`
- Python + Playwright is the right tradeoff for DOM drift and scraping iteration speed.
- Replacing this with Rust would slow you down.

## 2. Kafka And Redis

Use them for different jobs.

### Redis Should Handle

- rate limiting counters
- short-lived cache
- distributed locks
- idempotency keys
- websocket session fan-out
- small delayed jobs or retry metadata

Redis is not the best source of truth for the platform's durable event history.

### Kafka Should Handle

- crawler normalized events
- auction change events
- keyword match events
- notification jobs
- billing/entitlement events
- audit/event replay

For this repo, Kafka is justified when:

- crawler volume increases
- notification fan-out becomes bursty
- multiple downstream services need the same events
- replay/debugging matters

### Recommended Event Topics

- `auction.normalized`
- `auction.changed`
- `auction.matched`
- `notification.jobs`
- `notification.delivery`
- `billing.checkout.completed`
- `billing.entitlement.updated`

### Redis vs Kafka By Concrete Function

- Redis:
  - IP / token-bucket rate limit counters
  - short-lived cache for admin dashboards
  - distributed locks around scheduled crawler triggers
  - websocket fan-out presence and connection routing
  - idempotency keys for Stripe webhook replay handling
- Kafka / Redpanda:
  - crawler normalized payload publication
  - change-log fan-out to matcher / policy / push workers
  - notification job creation and retry streams
  - billing / entitlement event propagation
  - replayable audit/event debugging

### Practical Recommendation

- Keep Redis now.
- Add Redpanda/Kafka before splitting heavy async work into more services.
- Do not replace Postgres with Kafka or Redis. Postgres stays the source of truth.

## 3. DDoS, CSRF, SQL Injection, And General Hardening

## DDoS

Recommended layers:

1. edge protection
- Cloudflare, Fastly, or Koyeb edge protection in front of public endpoints

2. app rate limiting
- Redis-backed IP and token bucket limits
- separate limits for:
  - anonymous browse
  - auth endpoints
  - admin APIs
  - crawler ingest

3. origin hardening
- trusted proxy configuration
- body size limits
- request timeout and header timeout
- connection concurrency caps for webhook and push paths

This repo now includes request body limits, timeout middleware, header timeouts, and a configurable rate-limit RPS baseline in `api-go`.

## CSRF

For this repo:

- Bearer-token JSON APIs are lower-CSRF risk than cookie-only form posts
- SSR auth routes and any future cookie-mutating form actions still need CSRF protection

Recommended approach:

- SameSite `lax` or `strict` cookies where possible
- origin / referer validation on state-changing browser routes
- CSRF token on future cookie-based server actions or admin forms

This repo now includes a baseline origin/referer guard for mutating requests that carry cookies.

## SQL Injection

Current repo already benefits from parameterized queries in `pgx`.

Keep these rules:

- never build SQL with string concatenation for user input
- whitelist sort fields and filter operators
- centralize raw SQL review for admin/search endpoints
- prefer prepared/parameterized queries everywhere

The current `api-go` repository layer continues to use `pgx` parameter binding instead of string-concatenated user SQL.

## Additional Hardening Worth Adding

- CSP for web
- strict upload MIME/type validation
- signed object URLs only
- audit trail for admin mutations
- secret rotation runbook
- webhook replay protection and idempotency keys

## 4. JSON vs Protobuf vs gRPC vs WebSocket vs MQTT

Do not replace every public JSON API just because internal services are getting more advanced.

### Keep JSON For

- public web browse/search/detail endpoints
- SEO-facing pages
- admin dashboard fetches where developer velocity matters
- third-party integrations that expect REST

### Add Protobuf + gRPC For

- internal service-to-service calls
- Rust workers talking to Go control plane
- low-latency internal APIs
- strongly typed async/batch worker coordination

Recommended internal gRPC candidates:

- `matcher-rs <-> api-go`
- `push-rs <-> api-go`
- `policy-rs <-> api-go`

### Add WebSocket For

- admin live moderation queue updates
- crawler health live feed
- notification delivery live dashboard
- member real-time keyword alert stream in the web app

WebSocket is the right choice for browser real-time UX.

### Use MQTT Only If

- you need unreliable networks
- a large mobile/native device fleet needs lightweight pub/sub
- you move beyond browser-first into apps/devices

For this repo as a browser-first auction platform, MQTT is not the default recommendation.

### Practical Recommendation

- Keep external REST/JSON.
- Use Protobuf/gRPC internally.
- Use WebSocket for browser live updates.
- Use Kafka for durable async event streams.
- Do not adopt MQTT unless your client mix changes materially.

## 4.1 Internal Proto Scaffold Added

This repo now includes:

- `packages/shared-proto/customs/platform/v1/platform.proto`
- `services/matcher-rs`
- `services/policy-rs`
- `services/push-rs`
- `services/feed-rs`

The Rust services use `tonic` for gRPC scaffolding. They are intentionally narrow:

- `matcher-rs`: keyword/tag derivation
- `policy-rs`: compliance labeling
- `push-rs`: push payload preparation
- `feed-rs`: ICS rendering

The current scaffold is meant to establish stable contracts and compilation paths first. Durable Kafka consumers, Redis-backed coordination, and production auth between services should be layered on next.

## 5. Suggested Evolution Order

1. Keep `api-go` as the control plane.
2. Add Kafka/Redpanda for durable event streaming.
3. Split `matcher-rs`.
4. Split `push-rs`.
5. Add internal gRPC between Go and Rust services.
6. Add WebSocket gateway for admin/member real-time updates.
7. Add MQTT only if native/mobile/device clients become a core channel.

## 6. 100% Exit Criteria

The repo should only be called `100%` production-ready after all of these are true:

1. live deploy completed successfully
2. deploy smoke checks pass in hosted environment
3. live Stripe checkout + webhook roundtrip verified
4. live VAPID notification delivery verified
5. browser E2E coverage exists for auth, purchase, and notification paths
6. security baseline is verified for rate limiting, CSRF handling, and injection-safe query paths

Before those are true, the repo can be very close to launch, but not honestly `100%`.
