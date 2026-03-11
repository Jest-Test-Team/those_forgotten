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

## Contracts

```bash
make sync-contracts
```

This copies `packages/contracts/openapi.yaml` to `services/api-go/docs/swagger.yaml`.
