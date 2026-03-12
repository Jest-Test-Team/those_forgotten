# Vercel Env Verification

Use this checklist right after updating Production env vars for `apps/web` on Vercel.

## Required Web Env

```bash
NEXT_PUBLIC_SITE_URL=https://those-forgotten.vercel.app
NEXT_PUBLIC_API_BASE_URL=https://comfortable-adelind-dennis-team-1d8552ab.koyeb.app
API_BASE_URL=https://comfortable-adelind-dennis-team-1d8552ab.koyeb.app
NEXT_PUBLIC_SUPABASE_URL=https://<your-project-id>.supabase.co
NEXT_PUBLIC_SUPABASE_ANON_KEY=<your-supabase-anon-key>
NEXT_PUBLIC_SUPABASE_AUTH_REDIRECT=/auth/callback
ADMIN_EMAILS=<admin-email>
```

After saving env vars in Vercel, trigger a new Production deployment.

## Curl Checks

Run these after the new deployment reaches `Ready`.

### 1. Login page should no longer show auth-disabled fallback

```bash
curl -sS https://those-forgotten.vercel.app/auth/login | rg "尚未設定登入環境|NEXT_PUBLIC_SUPABASE_URL|NEXT_PUBLIC_SUPABASE_ANON_KEY"
```

Expected:

- no output

### 2. Login page should render the real sign-in CTA

```bash
curl -sS https://those-forgotten.vercel.app/auth/login | rg "Google 登入"
```

Expected:

- output contains `Google 登入`

### 3. Unauthenticated admin access should redirect to login

```bash
curl -I -sS https://those-forgotten.vercel.app/admin
```

Expected:

- `302` or `307`
- `location: /auth/login?next=/admin`

### 4. API base URL should be reachable from the browser-facing deployment

```bash
curl -sS https://comfortable-adelind-dennis-team-1d8552ab.koyeb.app/healthz
```

Expected:

- JSON response with `status`

## Robot Checks

Bootstrap once:

```bash
make robot-bootstrap
```

Run live smoke:

```bash
make robot-live-test
```

If you want to force explicit targets:

```bash
ROBOT_API_BASE_URL=https://comfortable-adelind-dennis-team-1d8552ab.koyeb.app \
ROBOT_WEB_BASE_URL=https://those-forgotten.vercel.app \
make robot-live-test
```

Expected:

- `14 tests, 14 passed, 0 failed`

## Functional Spot Checks

1. Open `/auth/login` in the browser.
2. Confirm the button says `Google 登入`, not `設定登入`.
3. Open `/admin` in an incognito window.
4. Confirm you are redirected to `/auth/login?next=/admin`.
5. Sign in with an email included in `ADMIN_EMAILS`.
6. Confirm `/admin` loads the dashboard instead of the auth-disabled fallback.

## Failure Mapping

- If `/auth/login` still shows `尚未設定登入環境`, Vercel env is still missing or malformed.
- If `/admin` loads but shows `API Repository unreachable`, `API_BASE_URL` is missing or the API is unreachable from Vercel.
- If `/admin` redirects to `/member?denied=admin`, auth is working and the account is not in admin role resolution.
