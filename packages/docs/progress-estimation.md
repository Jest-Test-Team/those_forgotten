# Progress Estimation

Estimated on 2026-03-12.

## Overall

- Overall delivery progress: `90%`
- Product discovery / information architecture: `91%`
- Production readiness: `75%`

The repo now covers nearly all planned v1 product surfaces and the main backend control paths. The remaining gap is no longer feature breadth; it is the last mile of production hardening: removing the final fallback auth paths, replacing simulated billing/push delivery with live provider calls, and proving the stack in a real deployed environment.

## Detailed Breakdown

| Area | Progress | Status | Notes |
| --- | ---: | --- | --- |
| Monorepo foundation | 96% | strong | Workspace layout, compose, make targets, CI, Dockerfiles, contracts sync, and deploy-oriented env matrices now exist. |
| Next.js public surfaces | 90% | strong | Public, member, and admin surfaces are coherent, mobile-first, and now use token-backed loaders/mutations for the protected paths. |
| API route coverage | 90% | strong | Main v1 endpoints exist, protected writes now enforce member/admin auth, and readiness/health endpoints support deployment verification. |
| Auth flow | 84% | strong | Supabase SSR login exists, browser/server protected requests now send Bearer tokens, and API prefers verified JWT identity before fallback paths. |
| RBAC | 85% | strong | Web guard, API admin guard, and DB-backed `user_roles` resolution are in place; remaining work is removing the last compatibility fallbacks. |
| Postgres persistence | 80% | strong | Member/admin writes, crawler runs, ingest change logs, Stripe entitlements, and notification jobs now have durable Postgres paths. |
| Crawler ingestion | 85% | strong | Per-office adapters, checksum validation, DB-backed ingest, change log persistence, crawler run tracking, notification job creation, and deploy-time ingest URL aliasing are wired. |
| Historical pricing | 60% | partial | Schema/API/UI path exists, but depends on forward-filled data and limited persistence depth. |
| Calendar / ICS | 72% | partial | Signed feed endpoint and web UX exist; token management and user-specific issuance need more hardening. |
| Community moderation | 76% | partial | Posting, reporting, admin queue, and resolve action exist; hide/suspend workflow and audit depth remain incomplete. |
| Advisor marketplace | 74% | partial | Directory, lead form, admin lead queue exist; assignment lifecycle and closure flow are still thin. |
| Monetization | 72% | partial | Typed Stripe checkout payloads, webhook processing, and membership/course entitlement persistence exist; live Stripe API session creation is still the main missing piece. |
| Push notifications | 76% | partial | Web push subscriptions are persisted, ingest creates notification jobs, and a worker can claim/process the queue; real VAPID delivery and retries still need to replace simulate mode. |
| Security / compliance | 80% | strong | Disclaimers, CORS, rate limit, auth guard, JWT preference, RBAC bootstrap, and warning labels are in place; remaining work is removing compatibility fallbacks and hardening uploads/audit depth. |
| Tests | 76% | partial | Go unit/smoke coverage now includes auth, billing, crawler persistence, and notification queue paths; browser E2E and live deploy smoke execution are still pending. |
| Deployability | 86% | strong | Vercel/Supabase/Koyeb plan, Terraform scaffold, CI, deploy-smoke workflow, `readyz`, smoke scripts, Koyeb fix guide, and environment matrix are in place; real secret provisioning and first live apply are still pending. |

## What Is Effectively Done

- Repo structure and developer workflow are usable.
- Public browsing experience is present and coherent.
- Admin UI now reads protected admin APIs.
- API contracts and implementation shape are aligned enough for parallel work.
- DB-backed RBAC bootstrap exists through `grant-role`.
- Stripe entitlement persistence and notification worker queueing are implemented.
- Deploy smoke checks and runtime readiness reporting now exist.
- Service-by-service env examples now align with the actual deploy targets and crawler ingest aliases.

## What Is Only Partially Done

- Some compatibility fallbacks remain while the stack transitions fully to verified token identity.
- Stripe checkout still returns a structured session URL rather than creating a live remote Stripe session.
- Notification delivery currently processes a real queue in simulate mode instead of full browser push delivery.
- Deployment assets and env documentation are present, but a full live environment run has not been executed from this repo yet.

## Biggest Remaining Gaps Before “Real MVP Launch”

1. Remove the last header/query-based identity fallbacks and rely on verified Supabase JWT/session identity only.
2. Replace structured Stripe checkout URLs with real Stripe API session creation.
3. Replace notification simulate mode with real Web Push / VAPID delivery and retry telemetry.
4. Tighten remaining admin lifecycle actions such as hide/suspend/assign/resolve end-to-end.
5. Perform the first real deployment with secrets applied and run post-deploy smoke checks.
6. Add end-to-end tests for login, admin access, purchases, and notification-critical paths.

## Suggested Next 5 Execution Steps

1. Replace the remaining compatibility auth fallbacks with strict Bearer-token enforcement.
2. Create real Stripe checkout sessions from `POST /v1/stripe/checkout`.
3. Swap the notification worker from simulate mode to VAPID-backed delivery.
4. Add admin lifecycle actions for advisor assignment and community hide/suspend.
5. Apply the stack to a real environment using the env matrix and run `.github/workflows/deploy-smoke.yml`.

## Confidence Notes

- Confidence in the `90%` feature-progress estimate: medium
- Confidence in the `75%` production-readiness estimate: high

Reason:

- The repo already covers most requested surfaces.
- The missing work is concentrated in fewer but more sensitive areas: security, persistence completeness, billing, background delivery, and deployment hardening.
