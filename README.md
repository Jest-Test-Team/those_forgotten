# Customs Auction Platform

Monorepo for a zh-TW customs auction discovery platform.

## Workspaces

- `apps/web`: Next.js App Router PWA
- `services/api-go`: Go Echo API
- `services/crawler`: Python Playwright crawler
- `packages/contracts`: OpenAPI contract package
- `packages/shared-proto`: internal Protobuf contracts
- `infra`: container images for web, api, and crawler
- `services/*-rs`: tonic-based Rust internal services

## Quick Start

```bash
make bootstrap
make dev
```

## Service Commands

```bash
cd apps/web && npm run dev
cd services/api-go && go run ./cmd/api
cd services/crawler && source .venv/bin/activate && python -m crawler.cli --fixtures fixtures
cargo check
```

## Docker Compose

```bash
docker compose up --build
```

This starts `postgres`, `redis`, `api`, `web`, and `crawler` with the demo ingestion token and internal service wiring.

Production-like override:

```bash
docker compose -f docker-compose.yml -f docker-compose.prod.yml config
docker compose -f docker-compose.yml -f docker-compose.prod.yml up -d
```

## Local Stack

```bash
make stack-up
make api-dev
make crawler-dev
make crawler-post-fixtures
make stack-down
```

Production-like stack helpers:

```bash
make stack-prod-config
make stack-prod-up
make stack-prod-logs
```

## Robot Live Tests

```bash
make robot-bootstrap
make robot-live-test
```

Default targets hit:

- `https://comfortable-adelind-dennis-team-1d8552ab.koyeb.app`
- `https://those-forgotten.vercel.app`

Override with:

```bash
ROBOT_API_BASE_URL=https://api.example.com ROBOT_WEB_BASE_URL=https://app.example.com make robot-live-test
```

## Role Management

```bash
DATABASE_URL=postgresql://postgres.mluxmwdbjunrqgyuqizn:YOUR_PASSWORD@aws-1-ap-south-1.pooler.supabase.com:6543/postgres?sslmode=require EMAIL=admin@example.com ROLE=admin make grant-role
```

## Seed Demo Data

```bash
DATABASE_URL=postgresql://postgres.mluxmwdbjunrqgyuqizn:YOUR_PASSWORD@aws-1-ap-south-1.pooler.supabase.com:6543/postgres?sslmode=require make seed-db
```

Manual migration helpers:

```bash
DATABASE_URL=postgresql://postgres.mluxmwdbjunrqgyuqizn:YOUR_PASSWORD@aws-1-ap-south-1.pooler.supabase.com:6543/postgres?sslmode=require make migrate-up
DATABASE_URL=postgresql://postgres.mluxmwdbjunrqgyuqizn:YOUR_PASSWORD@aws-1-ap-south-1.pooler.supabase.com:6543/postgres?sslmode=require make migrate-status
```

If `go run ./cmd/seed` reports `relation "profiles" does not exist`, the schema has not been migrated yet.

## Notifications Worker

```bash
DATABASE_URL=postgresql://postgres.mluxmwdbjunrqgyuqizn:YOUR_PASSWORD@aws-1-ap-south-1.pooler.supabase.com:6543/postgres?sslmode=require make notify-worker
```

## Smoke Check

```bash
API_BASE_URL=https://api.example.com WEB_BASE_URL=https://app.example.com make smoke-check
```

## Contracts

```bash
make sync-contracts
```

This copies `packages/contracts/openapi.yaml` to `services/api-go/docs/swagger.yaml`.

Swagger UI for API testing is exposed by `api-go` at:

- `/swagger`
- `/swagger.yaml`

## Deployment

- Deployment assessment and rollout plan: `packages/docs/deployment-plan.md`
- Koyeb deployment fix guide: `packages/docs/koyeb-deploy.md`
- Detailed progress estimate: `packages/docs/progress-estimation.md`
- Environment matrix: `packages/docs/env-matrix.md`
- Backend evolution guide: `packages/docs/backend-evolution.md`
- Hosted secrets checklist: `packages/docs/secrets-checklist.md`
- Vercel env verification checklist: `packages/docs/vercel-env-verification.md`
- Swagger testing guide: `packages/docs/swagger-testing.md`
- Database bootstrap guide: `packages/docs/database-bootstrap.md`
- Terraform scaffold: `infra/terraform`


## manager
api - koyeb : pcleegood@gmail.com
