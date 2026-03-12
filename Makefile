SHELL := /bin/zsh

.PHONY: bootstrap dev lint test typecheck api-test crawler-test web-test rust-check sync-contracts stack-build stack-up stack-web stack-logs stack-streaming-up stack-streaming-logs stack-prod-config stack-prod-up stack-prod-logs stack-down api-dev crawler-dev crawler-post-fixtures grant-role seed-db notify-worker smoke-check

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

rust-check:
	cargo check

sync-contracts:
	corepack pnpm --filter @customs/contracts build

stack-up:
	docker compose up -d postgres redis api

stack-build:
	docker compose build

stack-web:
	docker compose up -d web crawler

stack-logs:
	docker compose logs -f --tail=100 api web crawler

stack-streaming-up:
	docker compose --profile streaming up -d redpanda

stack-streaming-logs:
	docker compose --profile streaming logs -f --tail=100 redpanda

stack-prod-config:
	docker compose -f docker-compose.yml -f docker-compose.prod.yml config

stack-prod-up:
	docker compose -f docker-compose.yml -f docker-compose.prod.yml up -d

stack-prod-logs:
	docker compose -f docker-compose.yml -f docker-compose.prod.yml logs -f --tail=100 api web crawler

stack-down:
	docker compose down

api-dev:
	cd services/api-go && go run ./cmd/api

crawler-dev:
	cd services/crawler && . .venv/bin/activate && python -m crawler.cli --fixtures fixtures

crawler-post-fixtures:
	cd services/crawler && . .venv/bin/activate && python -m crawler.cli --fixtures fixtures --post --endpoint $${INGEST_ENDPOINT:-http://localhost:8080/internal/ingest/auctions} --token $${INTERNAL_INGEST_TOKEN:-demo-token}

grant-role:
	cd services/api-go && DATABASE_URL=$${DATABASE_URL:?set DATABASE_URL} go run ./cmd/grant-role --email $${EMAIL:?set EMAIL} --role $${ROLE:-admin} --name "$${NAME:-}"

seed-db:
	cd services/api-go && DATABASE_URL=$${DATABASE_URL:?set DATABASE_URL} go run ./cmd/seed

notify-worker:
	cd services/api-go && DATABASE_URL=$${DATABASE_URL:?set DATABASE_URL} NOTIFICATION_BATCH_SIZE=$${NOTIFICATION_BATCH_SIZE:-20} go run ./cmd/notify-worker

smoke-check:
	bash scripts/smoke-check.sh
