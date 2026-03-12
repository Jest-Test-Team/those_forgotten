# Terraform Scaffold

This folder contains the first-pass IaC scaffold for the recommended MVP deployment split:

- Vercel Hobby for `apps/web`
- Supabase Free for auth, Postgres, and storage
- Koyeb Starter for `services/api-go`
- Optional Koyeb Starter service for `services/crawler`

## Notes

- Heroku is intentionally excluded from the free-tier plan because it is not a practical free runtime for this repo anymore.
- Fly.io is intentionally excluded from this first scaffold because the repo's MVP target is better served by Vercel + Supabase + Koyeb and Fly's cost model is less predictable for a zero-budget default.
- The Koyeb Terraform provider exists, but provider/resource details should be revalidated before first apply because provider capabilities can change faster than Vercel/Supabase.

## Usage

```bash
cd infra/terraform
cp terraform.tfvars.example terraform.tfvars
terraform init
terraform plan
```

## Expected Inputs

- `vercel_api_token`
- `github_repo`
- `supabase_access_token`
- `supabase_organization_id`
- `supabase_region`
- `supabase_jwt_secret`
- `koyeb_token`
- `koyeb_org_id`

For the current Supabase project in ap-south-1, `database_url` should use the transaction pooler format on port `6543`:

```bash
postgresql://postgres.mluxmwdbjunrqgyuqizn:YOUR_PASSWORD@aws-1-ap-south-1.pooler.supabase.com:6543/postgres?sslmode=require
```

## Follow-up Before Apply

1. Replace demo env values with real production URLs and secrets.
2. Publish API and crawler images or switch Koyeb resources to repo-based deploy settings.
3. Confirm the exact Koyeb provider version and schema before first apply.
4. After the stack is up, grant the first admin role through the API CLI:

```bash
DATABASE_URL=postgresql://postgres.mluxmwdbjunrqgyuqizn:YOUR_PASSWORD@aws-1-ap-south-1.pooler.supabase.com:6543/postgres?sslmode=require EMAIL=admin@example.com ROLE=admin make grant-role
```

`ADMIN_EMAILS` remains in the web deployment as an emergency fallback, but the preferred path is now DB-backed `user_roles`.

`SUPABASE_JWT_SECRET` must also be set on `services/api-go` so admin API routes can verify Supabase Bearer tokens.
