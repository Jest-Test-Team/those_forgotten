-- +goose Up
CREATE TABLE crawler_runs (
  id UUID PRIMARY KEY,
  source TEXT NOT NULL,
  office TEXT NOT NULL,
  checksum TEXT NOT NULL,
  row_count INTEGER NOT NULL DEFAULT 0,
  status TEXT NOT NULL DEFAULT 'healthy',
  trigger_source TEXT NOT NULL,
  ran_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  next_run_at TIMESTAMPTZ
);

CREATE INDEX crawler_runs_office_ran_at_idx ON crawler_runs (office, ran_at DESC);

-- +goose Down
DROP INDEX IF EXISTS crawler_runs_office_ran_at_idx;
DROP TABLE IF EXISTS crawler_runs;
