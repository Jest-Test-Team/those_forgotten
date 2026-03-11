SHELL := /bin/zsh

.PHONY: bootstrap dev lint test typecheck api-test crawler-test web-test

bootstrap:
	pnpm install
	cd apps/web && pnpm install
	cd services/crawler && python3 -m venv .venv && . .venv/bin/activate && pip install -e ".[dev]"
	cd services/api-go && go mod tidy

dev:
	pnpm dev

lint:
	pnpm lint
	cd services/api-go && go test ./...
	cd services/crawler && . .venv/bin/activate && pytest

test:
	pnpm test
	cd services/api-go && go test ./...
	cd services/crawler && . .venv/bin/activate && pytest

typecheck:
	pnpm typecheck

api-test:
	cd services/api-go && go test ./...

crawler-test:
	cd services/crawler && . .venv/bin/activate && pytest

web-test:
	cd apps/web && pnpm test
