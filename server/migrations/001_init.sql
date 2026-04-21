CREATE TABLE IF NOT EXISTS schema_migrations (
  version BIGINT PRIMARY KEY,
  applied_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS admins (
  id TEXT PRIMARY KEY,
  email TEXT NOT NULL UNIQUE,
  password_hash TEXT NOT NULL,
  created_at TIMESTAMPTZ NOT NULL,
  updated_at TIMESTAMPTZ NOT NULL
);

CREATE TABLE IF NOT EXISTS accesses (
  id TEXT PRIMARY KEY,
  name TEXT NOT NULL,
  provider TEXT NOT NULL,
  config JSONB NOT NULL DEFAULT '{}'::jsonb,
  reserve TEXT NOT NULL DEFAULT '',
  deleted_at TIMESTAMPTZ NULL,
  created_at TIMESTAMPTZ NOT NULL,
  updated_at TIMESTAMPTZ NOT NULL
);

CREATE TABLE IF NOT EXISTS workflows (
  id TEXT PRIMARY KEY,
  name TEXT NOT NULL,
  description TEXT NOT NULL DEFAULT '',
  trigger TEXT NOT NULL DEFAULT 'manual',
  trigger_cron TEXT NOT NULL DEFAULT '',
  enabled BOOLEAN NOT NULL DEFAULT true,
  graph_draft JSONB NOT NULL DEFAULT '{}'::jsonb,
  graph_content JSONB NOT NULL DEFAULT '{}'::jsonb,
  has_draft BOOLEAN NOT NULL DEFAULT false,
  has_content BOOLEAN NOT NULL DEFAULT false,
  last_run_id TEXT NULL,
  last_run_status TEXT NULL,
  last_run_time TIMESTAMPTZ NULL,
  created_at TIMESTAMPTZ NOT NULL,
  updated_at TIMESTAMPTZ NOT NULL
);

CREATE TABLE IF NOT EXISTS workflow_runs (
  id TEXT PRIMARY KEY,
  workflow_id TEXT NOT NULL REFERENCES workflows(id) ON DELETE CASCADE,
  status TEXT NOT NULL,
  trigger TEXT NOT NULL,
  started_at TIMESTAMPTZ NOT NULL,
  ended_at TIMESTAMPTZ NULL,
  graph JSONB NOT NULL DEFAULT '{}'::jsonb,
  error TEXT NOT NULL DEFAULT '',
  created_at TIMESTAMPTZ NOT NULL,
  updated_at TIMESTAMPTZ NOT NULL
);

DO $$
BEGIN
  IF NOT EXISTS (SELECT 1 FROM pg_constraint WHERE conname = 'fk_workflows_last_run') THEN
    ALTER TABLE workflows
      ADD CONSTRAINT fk_workflows_last_run
      FOREIGN KEY (last_run_id) REFERENCES workflow_runs(id)
      DEFERRABLE INITIALLY DEFERRED;
  END IF;
END$$;

CREATE TABLE IF NOT EXISTS certificates (
  id TEXT PRIMARY KEY,
  source TEXT NOT NULL,
  subject_alt_names TEXT NOT NULL DEFAULT '',
  serial_number TEXT NOT NULL DEFAULT '',
  certificate TEXT NOT NULL DEFAULT '',
  private_key TEXT NOT NULL DEFAULT '',
  issuer_org TEXT NOT NULL DEFAULT '',
  key_algorithm TEXT NOT NULL DEFAULT '',
  validity_not_after TIMESTAMPTZ NULL,
  is_revoked BOOLEAN NOT NULL DEFAULT false,
  workflow_id TEXT NOT NULL DEFAULT '',
  workflow_run_id TEXT NOT NULL DEFAULT '',
  created_at TIMESTAMPTZ NOT NULL,
  updated_at TIMESTAMPTZ NOT NULL
);

CREATE TABLE IF NOT EXISTS settings (
  name TEXT PRIMARY KEY,
  content JSONB NOT NULL DEFAULT '{}'::jsonb,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
