# Environment Matrix

This matrix lists the runtime environment variables used by each deploy target.

## `apps/web`

Local / Vercel:

- `NEXT_PUBLIC_SITE_URL`
- `NEXT_PUBLIC_API_BASE_URL`
- `NEXT_PUBLIC_SUPABASE_URL`
- `NEXT_PUBLIC_SUPABASE_ANON_KEY`
- `NEXT_PUBLIC_SUPABASE_AUTH_REDIRECT`
- `ADMIN_EMAILS`

Server-only optional override:

- `API_BASE_URL`

Notes:

- `API_BASE_URL` is used by server-side loaders and SSR auth helpers. If unset, the web app falls back to `NEXT_PUBLIC_API_BASE_URL`.
- `ADMIN_EMAILS` is now a fallback only. Preferred admin authority comes from API role resolution and `user_roles`.

## `services/api-go`

- `APP_ENV`
- `PORT`
- `WEB_ORIGIN`
- `DATABASE_URL`
- `REDIS_URL`
- `KAFKA_BROKERS`
- `INTERNAL_INGEST_TOKEN`
- `SUPABASE_JWT_SECRET`
- `STRIPE_CHECKOUT_BASE_URL`
- `STRIPE_SECRET_KEY`
- `STRIPE_WEBHOOK_SECRET`
- `STRIPE_SUCCESS_URL`
- `STRIPE_CANCEL_URL`
- `STRIPE_MEMBERSHIP_PRICE_ID`
- `STRIPE_COURSE_PRICE_ID`
- `VAPID_PUBLIC_KEY`
- `VAPID_PRIVATE_KEY`
- `VAPID_SUBJECT`
- `GRPC_PORT`
- `WS_ALLOWED_ORIGINS`
- `CSRF_TRUSTED_ORIGINS`
- `TRUSTED_PROXY_CIDRS`
- `RATE_LIMIT_RPS`

Notes:

- `DATABASE_URL` is required for production persistence, admin RBAC bootstrap, notification jobs, and worker processing.
- `REDIS_URL` remains the right place for rate limit counters, locks, short-lived cache, and websocket session fan-out.
- `KAFKA_BROKERS` is the recommended event backbone once crawler ingest and notification fan-out move to durable streaming.
- `SUPABASE_JWT_SECRET` is required if admin/member API routes should trust Bearer tokens from Supabase.
- If `STRIPE_SECRET_KEY`, price ids, and redirect URLs are present, `POST /v1/stripe/checkout` now creates a live Stripe Checkout Session.
- `STRIPE_CHECKOUT_BASE_URL` remains as the fallback path when live Stripe credentials are not present.
- If the three `VAPID_*` variables are present, `make notify-worker` sends real Web Push notifications; otherwise it falls back to simulate mode.

## `services/crawler`

- `INGEST_URL`
- `INTERNAL_INGEST_TOKEN`
- `POST_TO_API`
- `FIXTURES`

Backward-compatible alias:

- `INGEST_ENDPOINT`

Notes:

- Koyeb/Terraform currently uses `INGEST_URL`; the crawler now supports both names.
- `POST_TO_API=true` is required if the crawler container should send normalized payloads to `api-go`.

## Current Hosted Split

- Vercel: `apps/web`
- Supabase: Auth, Postgres, Storage
- Koyeb: `services/api-go`
- GitHub Actions cron or optional Koyeb service: `services/crawler`
