-- +goose Up
CREATE TABLE IF NOT EXISTS users(
  user_id uuid PRIMARY KEY,
  email text NOT NULL,
  email_normalized text NOT NULL UNIQUE,
  password_hash text NOT NULL,
  status text NOT NULL DEFAULT 'pending_verification',
  email_verified_at timestamptz NULL,
  display_name text NOT NULL DEFAULT '',
  locale text NOT NULL DEFAULT 'ru-RU',
  timezone text NOT NULL DEFAULT 'UTC',
  roles text[] NOT NULL DEFAULT ARRAY['user']::text[],
  is_two_factor_enabled boolean NOT NULL DEFAULT FALSE,
  totp_secret_ciphertext bytea NULL,
  profile jsonb NOT NULL DEFAULT '{}'::jsonb,
  settings jsonb NOT NULL DEFAULT '{}'::jsonb,
  session_version bigint NOT NULL DEFAULT 1,
  last_login_at timestamptz NULL,
  blocked_reason text NOT NULL DEFAULT '',
  blocked_at timestamptz NULL,
  created_at timestamptz NOT NULL DEFAULT now(),
  updated_at timestamptz NOT NULL DEFAULT now(),
  CONSTRAINT chk_users_status CHECK (status IN ('pending_verification', 'active', 'blocked'))
);
CREATE INDEX IF NOT EXISTS idx_users_status ON users(status);
CREATE INDEX IF NOT EXISTS idx_users_created_at ON users(created_at DESC);
CREATE TABLE IF NOT EXISTS auth_action_tokens(
  token_id uuid PRIMARY KEY,
  user_id uuid NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,
  purpose text NOT NULL,
  token_hash text NOT NULL UNIQUE,
  expires_at timestamptz NOT NULL,
  consumed_at timestamptz NULL,
  delivery_value text NOT NULL DEFAULT '',
  meta jsonb NOT NULL DEFAULT '{}'::jsonb,
  created_at timestamptz NOT NULL DEFAULT now(),
  CONSTRAINT chk_auth_action_tokens_purpose CHECK (purpose IN ('verify_email', 'password_reset'))
);
CREATE INDEX IF NOT EXISTS idx_auth_action_tokens_user_id ON auth_action_tokens(user_id, purpose, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_auth_action_tokens_expires_at ON auth_action_tokens(expires_at);
CREATE TABLE IF NOT EXISTS auth_sessions(
  session_id uuid PRIMARY KEY,
  user_id uuid NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,
  family_id uuid NOT NULL,
  refresh_jti text NOT NULL UNIQUE,
  refresh_token_hash text NOT NULL UNIQUE,
  csrf_token_hash text NOT NULL,
  fingerprint_hash text NOT NULL,
  device_label text NOT NULL DEFAULT '',
  ip inet NULL,
  user_agent text NOT NULL DEFAULT '',
  mfa_verified boolean NOT NULL DEFAULT FALSE,
  replaced_by_session_id uuid NULL REFERENCES auth_sessions(session_id) ON DELETE SET NULL,
  last_seen_at timestamptz NOT NULL DEFAULT now(),
  expires_at timestamptz NOT NULL,
  revoked_at timestamptz NULL,
  revoke_reason text NOT NULL DEFAULT '',
  compromised_at timestamptz NULL,
  created_at timestamptz NOT NULL DEFAULT now(),
  updated_at timestamptz NOT NULL DEFAULT now()
);
CREATE INDEX IF NOT EXISTS idx_auth_sessions_user_id ON auth_sessions(user_id, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_auth_sessions_family_id ON auth_sessions(family_id);
CREATE INDEX IF NOT EXISTS idx_auth_sessions_active ON auth_sessions(user_id, revoked_at, expires_at);
CREATE TABLE IF NOT EXISTS auth_backup_codes(
  backup_code_id uuid PRIMARY KEY,
  user_id uuid NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,
  code_hash text NOT NULL,
  consumed_at timestamptz NULL,
  created_at timestamptz NOT NULL DEFAULT now()
);
CREATE INDEX IF NOT EXISTS idx_auth_backup_codes_user_id ON auth_backup_codes(user_id, consumed_at);
CREATE TABLE IF NOT EXISTS auth_login_attempts(
  attempt_id uuid PRIMARY KEY,
  user_id uuid NULL REFERENCES users(user_id) ON DELETE SET NULL,
  email_normalized text NOT NULL,
  ip inet NULL,
  user_agent text NOT NULL DEFAULT '',
  outcome text NOT NULL,
  failure_reason text NOT NULL DEFAULT '',
  risk_score integer NOT NULL DEFAULT 0,
  fingerprint_hash text NOT NULL DEFAULT '',
  suspicious_reason text NOT NULL DEFAULT '',
  created_at timestamptz NOT NULL DEFAULT now()
);
CREATE INDEX IF NOT EXISTS idx_auth_login_attempts_email_created ON auth_login_attempts(email_normalized, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_auth_login_attempts_user_created ON auth_login_attempts(user_id, created_at DESC);
CREATE TABLE IF NOT EXISTS auth_security_events(
  event_id uuid PRIMARY KEY,
  user_id uuid NULL REFERENCES users(user_id) ON DELETE SET NULL,
  actor_user_id uuid NULL REFERENCES users(user_id) ON DELETE SET NULL,
  session_id uuid NULL REFERENCES auth_sessions(session_id) ON DELETE SET NULL,
  category text NOT NULL,
  event_type text NOT NULL,
  severity text NOT NULL,
  ip inet NULL,
  user_agent text NOT NULL DEFAULT '',
  payload jsonb NOT NULL DEFAULT '{}'::jsonb,
  created_at timestamptz NOT NULL DEFAULT now()
);
CREATE INDEX IF NOT EXISTS idx_auth_security_events_user_id ON auth_security_events(user_id, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_auth_security_events_category ON auth_security_events(category, created_at DESC);
-- +goose StatementBegin
CREATE OR REPLACE FUNCTION set_auth_updated_at()
  RETURNS TRIGGER
  AS $fn$
BEGIN
  NEW.updated_at = now();
  RETURN NEW;
END;
$fn$
LANGUAGE plpgsql;
-- +goose StatementEnd
-- +goose StatementBegin
DROP TRIGGER IF EXISTS trg_users_updated_at ON users;
CREATE TRIGGER trg_users_updated_at
  BEFORE UPDATE ON users
  FOR EACH ROW
  EXECUTE FUNCTION set_auth_updated_at();
DROP TRIGGER IF EXISTS trg_auth_sessions_updated_at ON auth_sessions;
CREATE TRIGGER trg_auth_sessions_updated_at
  BEFORE UPDATE ON auth_sessions
  FOR EACH ROW
  EXECUTE FUNCTION set_auth_updated_at();
-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TRIGGER IF EXISTS trg_auth_sessions_updated_at ON auth_sessions;
DROP TRIGGER IF EXISTS trg_users_updated_at ON users;
DROP FUNCTION IF EXISTS set_auth_updated_at();
-- +goose StatementEnd
DROP TABLE IF EXISTS auth_security_events;
DROP TABLE IF EXISTS auth_login_attempts;
DROP TABLE IF EXISTS auth_backup_codes;
DROP TABLE IF EXISTS auth_sessions;
DROP TABLE IF EXISTS auth_action_tokens;
DROP TABLE IF EXISTS users;
