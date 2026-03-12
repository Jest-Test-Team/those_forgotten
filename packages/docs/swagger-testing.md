# Swagger Testing Guide

The API now serves Swagger UI directly, so you can test routes without wiring a separate tool first.

## Hosted Or Local URLs

- UI: `http://localhost:8080/swagger`
- Raw spec: `http://localhost:8080/swagger.yaml`

If the API is hosted, replace the host with the Koyeb domain.

## Local Flow

1. Start the API:

```bash
make stack-up
```

2. Open:

```text
http://localhost:8080/swagger
```

3. If you want seeded Postgres-backed responses first:

```bash
DATABASE_URL=postgresql://postgres.mluxmwdbjunrqgyuqizn:YOUR_PASSWORD@aws-1-ap-south-1.pooler.supabase.com:6543/postgres?sslmode=require make seed-db
```

## Useful Test Routes

No auth required:

- `GET /healthz`
- `GET /readyz`
- `GET /v1/auctions`
- `GET /v1/auctions/{id}`
- `GET /v1/auctions/{id}/history`
- `GET /v1/community/posts`
- `GET /v1/advisors`
- `GET /v1/courses`

Requires Bearer token once `SUPABASE_JWT_SECRET` is configured:

- `GET /v1/keyword-subscriptions`
- `POST /v1/keyword-subscriptions`
- `POST /v1/community/posts`
- `POST /v1/stripe/checkout`
- `GET /v1/admin/community-reports`

## Koyeb Note

The API container now copies `docs/swagger.yaml` into the runtime image, so `/swagger` and `/swagger.yaml` work in Koyeb as long as the service itself is healthy.
