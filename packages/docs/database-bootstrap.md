# Database Bootstrap

Use this flow when the target Postgres database is empty.

## 1. Apply Schema Migrations

```bash
DATABASE_URL=postgresql://postgres.mluxmwdbjunrqgyuqizn:YOUR_PASSWORD@aws-1-ap-south-1.pooler.supabase.com:6543/postgres?sslmode=require make migrate-up
```

Check migration state:

```bash
DATABASE_URL=postgresql://postgres.mluxmwdbjunrqgyuqizn:YOUR_PASSWORD@aws-1-ap-south-1.pooler.supabase.com:6543/postgres?sslmode=require make migrate-status
```

## 2. Seed Demo Data

```bash
DATABASE_URL=postgresql://postgres.mluxmwdbjunrqgyuqizn:YOUR_PASSWORD@aws-1-ap-south-1.pooler.supabase.com:6543/postgres?sslmode=require make seed-db
```

`make seed-db` now runs `migrate up` before `seed`.

## Common Failure

If you run:

```bash
go run ./cmd/seed
```

and get:

```text
ERROR: relation "profiles" does not exist
```

that means the schema has not been created yet. Run `go run ./cmd/migrate up` or `make migrate-up` first.
