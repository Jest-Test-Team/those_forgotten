SHELL := /bin/zsh

.PHONY: bootstrap dev lint test typecheck api-test crawler-test web-test sync-contracts

bootstrap:
	corepack pnpm install
	cd apps/web && npm install
	cd services/crawler && python3 -m venv .venv && . .venv/bin/activate && pip install -e ".[dev]"
	cd services/api-go && go mod tidy
	corepack pnpm --filter @customs/contracts build

dev:
	corepack pnpm dev

lint:
	corepack pnpm lint
	cd services/api-go && go test ./...
	cd services/crawler && . .venv/bin/activate && pytest

test:
	corepack pnpm test
	cd services/api-go && go test ./...
	cd services/crawler && . .venv/bin/activate && pytest

typecheck:
	corepack pnpm typecheck

api-test:
	cd services/api-go && go test ./...

crawler-test:
	cd services/crawler && . .venv/bin/activate && pytest

web-test:
	cd apps/web && npm test

sync-contracts:
	corepack pnpm --filter @customs/contracts build
