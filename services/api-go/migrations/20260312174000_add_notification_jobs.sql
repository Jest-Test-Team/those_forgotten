-- +goose Up
CREATE TABLE notification_jobs (
  id UUID PRIMARY KEY,
  auction_lot_id UUID NOT NULL REFERENCES auction_lots(id) ON DELETE CASCADE,
  keyword TEXT NOT NULL,
  endpoint TEXT NOT NULL,
  p256dh TEXT NOT NULL,
  auth_secret TEXT NOT NULL,
  payload JSONB NOT NULL,
  status TEXT NOT NULL DEFAULT 'pending',
  last_error TEXT,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  delivered_at TIMESTAMPTZ
);

CREATE INDEX notification_jobs_status_created_at_idx ON notification_jobs (status, created_at);

-- +goose Down
DROP INDEX IF EXISTS notification_jobs_status_created_at_idx;
DROP TABLE IF EXISTS notification_jobs;
