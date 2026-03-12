# Koyeb Deployment Fix

This repo should be deployed to Koyeb using the Dockerfile builder, not the Node.js buildpack.

## Why The Current Deployment Failed

Your failed deployment used:

- builder type: `Buildpack`
- custom build command: `./infra/Dockerfile.api`

That makes Koyeb treat the repository root like a Node.js app. It then tries to resolve the package manager from the root `package.json`, fails to find a root lockfile, and never builds the Go API container. A Dockerfile path is not a shell build command.

## Correct Koyeb Configuration For `services/api-go`

Create a new Koyeb **Web Service** with these settings:

- Deployment method: `GitHub`
- Repository: `Jest-Test-Team/those_forgotten`
- Branch: `main`
- Builder: `Dockerfile`
- Dockerfile path: `infra/Dockerfile.api`
- Exposed port: `8080`
- Health check type: `HTTP`
- Health check path: `/readyz`

Do not set a custom build command for the API service.

## Required Environment Variables

Set these on the Koyeb API service:

- `PORT=8080`
- `DATABASE_URL`
- `WEB_ORIGIN`
- `ADMIN_EMAILS`
- `INTERNAL_INGEST_TOKEN`
- `SUPABASE_JWT_SECRET`
- `STRIPE_CHECKOUT_BASE_URL`
- `STRIPE_WEBHOOK_SECRET`

## Why `/readyz`

The API now exposes:

- `/` for a simple 200 response on the root path
- `/healthz` for basic process health
- `/readyz` for deployment readiness signals

Koyeb recommends HTTP health checks for HTTP services, and the probe path must return a `2xx` or `3xx` status.

## Crawler On Koyeb

If you deploy the crawler on Koyeb too, use a separate service with:

- Builder: `Dockerfile`
- Dockerfile path: `infra/Dockerfile.crawler`

For zero-budget setups, prefer GitHub Actions cron first and keep Koyeb crawler deployment optional.

## Current Recommendation

For the API, either:

1. Recreate the current Koyeb service using the Dockerfile builder.
2. Or edit the service and switch the builder from Buildpack to Dockerfile if the UI allows it.

Do not keep the current `Buildpack + ./infra/Dockerfile.api` configuration.

## Official References

- Dockerfile-based deployment on Koyeb: https://www.koyeb.com/docs/deploy/go
- Koyeb health checks: https://www.koyeb.com/docs/run-and-scale/health-checks
