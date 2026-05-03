ALTER TABLE admins
  ADD COLUMN IF NOT EXISTS role TEXT NOT NULL DEFAULT 'admin',
  ADD COLUMN IF NOT EXISTS status TEXT NOT NULL DEFAULT 'active';

UPDATE admins SET role='admin' WHERE role IS NULL OR role='';
UPDATE admins SET status='active' WHERE status IS NULL OR status='';

ALTER TABLE accesses ADD COLUMN IF NOT EXISTS owner_user_id TEXT;
ALTER TABLE workflows ADD COLUMN IF NOT EXISTS owner_user_id TEXT;
ALTER TABLE certificates ADD COLUMN IF NOT EXISTS owner_user_id TEXT;
ALTER TABLE workflow_runs ADD COLUMN IF NOT EXISTS owner_user_id TEXT;

DO $$
DECLARE
  bootstrap_id TEXT;
BEGIN
  SELECT id INTO bootstrap_id FROM admins ORDER BY created_at ASC LIMIT 1;

  IF bootstrap_id IS NOT NULL THEN
    UPDATE accesses SET owner_user_id=bootstrap_id WHERE owner_user_id IS NULL;
    UPDATE workflows SET owner_user_id=bootstrap_id WHERE owner_user_id IS NULL;
    UPDATE certificates SET owner_user_id=bootstrap_id WHERE owner_user_id IS NULL;
    UPDATE workflow_runs SET owner_user_id=bootstrap_id WHERE owner_user_id IS NULL;
  END IF;
END $$;

DO $$
BEGIN
  IF NOT EXISTS (SELECT 1 FROM pg_constraint WHERE conname = 'fk_accesses_owner_user') THEN
    ALTER TABLE accesses ADD CONSTRAINT fk_accesses_owner_user FOREIGN KEY (owner_user_id) REFERENCES admins(id) ON DELETE SET NULL;
  END IF;
  IF NOT EXISTS (SELECT 1 FROM pg_constraint WHERE conname = 'fk_workflows_owner_user') THEN
    ALTER TABLE workflows ADD CONSTRAINT fk_workflows_owner_user FOREIGN KEY (owner_user_id) REFERENCES admins(id) ON DELETE SET NULL;
  END IF;
  IF NOT EXISTS (SELECT 1 FROM pg_constraint WHERE conname = 'fk_certificates_owner_user') THEN
    ALTER TABLE certificates ADD CONSTRAINT fk_certificates_owner_user FOREIGN KEY (owner_user_id) REFERENCES admins(id) ON DELETE SET NULL;
  END IF;
  IF NOT EXISTS (SELECT 1 FROM pg_constraint WHERE conname = 'fk_workflow_runs_owner_user') THEN
    ALTER TABLE workflow_runs ADD CONSTRAINT fk_workflow_runs_owner_user FOREIGN KEY (owner_user_id) REFERENCES admins(id) ON DELETE SET NULL;
  END IF;
END $$;

CREATE INDEX IF NOT EXISTS idx_accesses_owner_user_id ON accesses(owner_user_id);
CREATE INDEX IF NOT EXISTS idx_workflows_owner_user_id ON workflows(owner_user_id);
CREATE INDEX IF NOT EXISTS idx_certificates_owner_user_id ON certificates(owner_user_id);
CREATE INDEX IF NOT EXISTS idx_workflow_runs_owner_user_id ON workflow_runs(owner_user_id);

CREATE TABLE IF NOT EXISTS api_keys (
  id TEXT PRIMARY KEY,
  user_id TEXT NOT NULL REFERENCES admins(id) ON DELETE CASCADE,
  name TEXT NOT NULL,
  prefix TEXT NOT NULL,
  key_hash TEXT NOT NULL,
  status TEXT NOT NULL DEFAULT 'active',
  expires_at TIMESTAMPTZ NULL,
  last_used_at TIMESTAMPTZ NULL,
  revoked_at TIMESTAMPTZ NULL,
  created_at TIMESTAMPTZ NOT NULL,
  updated_at TIMESTAMPTZ NOT NULL
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_api_keys_prefix ON api_keys(prefix);
CREATE INDEX IF NOT EXISTS idx_api_keys_user_id ON api_keys(user_id);
CREATE INDEX IF NOT EXISTS idx_api_keys_status ON api_keys(status);
