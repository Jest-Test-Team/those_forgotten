# Progress Estimation

Estimated on 2026-03-12.

## Overall

- Overall delivery progress: `72%`
- Product discovery / information architecture: `85%`
- Production readiness: `54%`

The repo is beyond scaffold stage and already has coherent web, API, crawler, contracts, CI, and deployment planning. The biggest remaining gap is not breadth of features; it is replacing demo/fallback paths with real production integrations and closing the auth / background-job / persistence gaps.

## Detailed Breakdown

| Area | Progress | Status | Notes |
| --- | ---: | --- | --- |
| Monorepo foundation | 95% | strong | Workspace layout, compose, make targets, CI, Dockerfiles, contracts sync all exist. |
| Next.js public surfaces | 88% | strong | Home, auctions, detail, knowledge, advisors, community, member, admin, SEO/PWA surfaces are in place. |
| API route coverage | 84% | strong | Main v1 endpoints exist and are documented in OpenAPI/Swagger. |
| Auth flow | 68% | partial | Supabase SSR login path exists; API now prefers Bearer JWT for identity, but some browser-side mutations still are not token-backed. |
| RBAC | 70% | partial | Web guard + API admin route guard + DB-backed `user_roles` fallback chain exist; final JWT-only enforcement path is not complete across all writes. |
| Postgres persistence | 58% | partial | Read path and selected member/admin writes support Postgres; not every mutation and queue path is fully persisted. |
| Crawler ingestion | 62% | partial | Per-office adapters, fixtures, checksum payloads, ingest client, and admin status views exist; production scheduler and durable ingestion pipeline still need tightening. |
| Historical pricing | 60% | partial | Schema/API/UI path exists, but depends on forward-filled data and limited persistence depth. |
| Calendar / ICS | 72% | partial | Signed feed endpoint and web UX exist; token management and user-specific issuance need more hardening. |
| Community moderation | 76% | partial | Posting, reporting, admin queue, and resolve action exist; hide/suspend workflow and audit depth remain incomplete. |
| Advisor marketplace | 74% | partial | Directory, lead form, admin lead queue exist; assignment lifecycle and closure flow are still thin. |
| Monetization | 32% | early | Stripe endpoints and product surfaces are scaffolded, but real checkout, webhook verification, entitlement sync, and course access are not finished. |
| Push notifications | 36% | early | Web push subscription capture exists; matching worker, VAPID delivery, retry logic, and delivery telemetry are still missing. |
| Security / compliance | 66% | partial | Disclaimers, CORS, rate limit, auth guard, role bootstrap, and warning labels exist; final token verification, upload constraints, and audit/abuse controls need more work. |
| Tests | 64% | partial | Go unit/smoke tests and crawler tests exist; end-to-end auth/billing/push flows and DB integration depth remain incomplete. |
| Deployability | 73% | partial | Vercel/Supabase/Koyeb target plan, Terraform scaffold, CI, and compose exist; real secrets, apply flow, and runtime observability are not yet closed. |

## What Is Effectively Done

- Repo structure and developer workflow are usable.
- Public browsing experience is present and coherent.
- Admin UI now reads protected admin APIs.
- API contracts and implementation shape are aligned enough for parallel work.
- DB-backed RBAC bootstrap exists through `grant-role`.

## What Is Only Partially Done

- API authentication is improved, but not fully standardized across all server and browser requests.
- Postgres support exists, but not every feature path is durable.
- Community, advisor, and crawler admin flows are functional but not yet full operations tooling.
- Deployment assets are scaffolded, not yet fully proven against a live environment.

## Biggest Remaining Gaps Before “Real MVP Launch”

1. Replace remaining caller-supplied identity flows with verified Supabase session/JWT handling everywhere.
2. Finish Postgres persistence for the remaining mutations and admin operations.
3. Implement real Stripe checkout, webhook verification, and entitlement updates.
4. Implement actual push delivery worker and notification matching pipeline.
5. Tighten crawler scheduling, retries, and ingestion observability in a live environment.
6. Add end-to-end tests for login, admin access, purchases, and notification-critical paths.

## Suggested Next 5 Execution Steps

1. Convert browser-side admin/member mutations to send Supabase Bearer tokens instead of anonymous requests.
2. Complete admin-only API authorization middleware using verified token identity only.
3. Finish Postgres-backed persistence for advisor/admin/crawler queue mutations.
4. Wire Stripe checkout + webhook idempotency + membership/course entitlement state.
5. Add deployment secrets, perform first live deploy, and run smoke tests against the real stack.

## Confidence Notes

- Confidence in the `72%` feature-progress estimate: medium
- Confidence in the `54%` production-readiness estimate: high

Reason:

- The repo already covers most requested surfaces.
- The missing work is concentrated in fewer but more sensitive areas: security, persistence completeness, billing, background delivery, and deployment hardening.
