# Customs Auction Platform

Monorepo for a zh-TW customs auction discovery platform.

## Workspaces

- `apps/web`: Next.js App Router PWA
- `services/api-go`: Go Echo API
- `services/crawler`: Python Playwright crawler
- `packages/contracts`: OpenAPI contract package
- `infra`: container images for web, api, and crawler

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
```

## Docker Compose

```bash
docker compose up --build
```

This starts `postgres`, `redis`, `api`, `web`, and `crawler` with the demo ingestion token and internal service wiring.

## Local Stack

```bash
make stack-up
make api-dev
make crawler-dev
make crawler-post-fixtures
make stack-down
```

## Role Management

```bash
DATABASE_URL=postgres://... EMAIL=admin@example.com ROLE=admin make grant-role
```

## Notifications Worker

```bash
DATABASE_URL=postgres://... make notify-worker
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

## Deployment

- Deployment assessment and rollout plan: `packages/docs/deployment-plan.md`
- Koyeb deployment fix guide: `packages/docs/koyeb-deploy.md`
- Detailed progress estimate: `packages/docs/progress-estimation.md`
- Environment matrix: `packages/docs/env-matrix.md`
- Backend evolution guide: `packages/docs/backend-evolution.md`
- Terraform scaffold: `infra/terraform`


## manager
api - koyeb : pcleegood@gmail.com
