# Hosted Secrets Checklist

Use this checklist before the first live deployment. It is intentionally grouped by hosting boundary, not by source file.

## Vercel `apps/web`

Required:

- `NEXT_PUBLIC_SITE_URL`
- `NEXT_PUBLIC_API_BASE_URL`
- `NEXT_PUBLIC_SUPABASE_URL`
- `NEXT_PUBLIC_SUPABASE_ANON_KEY`

Recommended:

- `API_BASE_URL`
- `NEXT_PUBLIC_SUPABASE_AUTH_REDIRECT`
- `ADMIN_EMAILS`

Example:

```bash
NEXT_PUBLIC_SITE_URL=https://app.example.com
NEXT_PUBLIC_API_BASE_URL=https://api.example.com
API_BASE_URL=https://api.example.com
NEXT_PUBLIC_SUPABASE_URL=https://your-project.supabase.co
NEXT_PUBLIC_SUPABASE_ANON_KEY=eyJ...
NEXT_PUBLIC_SUPABASE_AUTH_REDIRECT=/auth/callback
ADMIN_EMAILS=admin@example.com
```

## Koyeb `services/api-go`

Required:

- `PORT=8080`
- `WEB_ORIGIN`
- `DATABASE_URL`
- `REDIS_URL`
- `INTERNAL_INGEST_TOKEN`
- `SUPABASE_JWT_SECRET`
- `STRIPE_WEBHOOK_SECRET`

Required for live Stripe:

- `STRIPE_SECRET_KEY`
- `STRIPE_SUCCESS_URL`
- `STRIPE_CANCEL_URL`
- `STRIPE_MEMBERSHIP_PRICE_ID`
- `STRIPE_COURSE_PRICE_ID`

Required for live Web Push:

- `VAPID_PUBLIC_KEY`
- `VAPID_PRIVATE_KEY`
- `VAPID_SUBJECT`

Optional hardening / future transport:

- `ADMIN_EMAILS`
- `KAFKA_BROKERS`
- `GRPC_PORT`
- `WS_ALLOWED_ORIGINS`
- `CSRF_TRUSTED_ORIGINS`
- `TRUSTED_PROXY_CIDRS`
- `RATE_LIMIT_RPS`

Supabase transaction pooler example:

```bash
DATABASE_URL=postgresql://postgres.mluxmwdbjunrqgyuqizn:YOUR_PASSWORD@aws-1-ap-south-1.pooler.supabase.com:6543/postgres?sslmode=require
```

## Koyeb Or GitHub Actions `services/crawler`

Required:

- `INGEST_URL`
- `INTERNAL_INGEST_TOKEN`

Recommended:

- `POST_TO_API=true`
- `FIXTURES=fixtures`

## Supabase Project

Collect and store:

- Project URL
- Anon key
- JWT secret
- Transaction pooler connection string
- Google OAuth client id / secret
- Storage bucket names

If web and API are already live, also verify these Supabase-side settings:

- Auth redirect URLs include `https://app.example.com/auth/callback`
- Site URL matches the web origin
- RLS and service-role usage are limited to crawler / backend contexts only

## First Deploy Gate

Before calling the hosted stack ready:

1. `web`, `api`, and `crawler` all have their secrets provisioned.
2. `DATABASE_URL` uses the Supabase pooler URL with `sslmode=require`.
3. `SUPABASE_JWT_SECRET` on `api-go` matches the current Supabase project.
4. Stripe webhook secret is set in Koyeb and the endpoint is registered in Stripe.
5. VAPID keys are provisioned if live push delivery is expected.
