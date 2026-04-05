-- +goose Up
CREATE TABLE IF NOT EXISTS users (
  user_id UUID PRIMARY KEY,
  email TEXT NOT NULL,
  email_normalized TEXT NOT NULL UNIQUE,
  password_hash TEXT NOT NULL,
  status TEXT NOT NULL DEFAULT 'pending_verification',
  email_verified_at TIMESTAMPTZ NULL,
  display_name TEXT NOT NULL DEFAULT '',
  locale TEXT NOT NULL DEFAULT 'ru-RU',
  timezone TEXT NOT NULL DEFAULT 'UTC',
  roles TEXT[] NOT NULL DEFAULT ARRAY['user']::TEXT[],
  is_two_factor_enabled BOOLEAN NOT NULL DEFAULT FALSE,
  totp_secret_ciphertext BYTEA NULL,
  profile JSONB NOT NULL DEFAULT '{}'::jsonb,
  settings JSONB NOT NULL DEFAULT '{}'::jsonb,
  session_version BIGINT NOT NULL DEFAULT 1,
  last_login_at TIMESTAMPTZ NULL,
  blocked_reason TEXT NOT NULL DEFAULT '',
  blocked_at TIMESTAMPTZ NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  CONSTRAINT chk_users_status CHECK (status IN ('pending_verification', 'active', 'blocked'))
);

CREATE INDEX IF NOT EXISTS idx_users_status ON users (status);
CREATE INDEX IF NOT EXISTS idx_users_created_at ON users (created_at DESC);

CREATE TABLE IF NOT EXISTS auth_action_tokens (
  token_id UUID PRIMARY KEY,
  user_id UUID NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,
  purpose TEXT NOT NULL,
  token_hash TEXT NOT NULL UNIQUE,
  expires_at TIMESTAMPTZ NOT NULL,
  consumed_at TIMESTAMPTZ NULL,
  delivery_value TEXT NOT NULL DEFAULT '',
  meta JSONB NOT NULL DEFAULT '{}'::jsonb,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  CONSTRAINT chk_auth_action_tokens_purpose CHECK (purpose IN ('verify_email', 'password_reset'))
);

CREATE INDEX IF NOT EXISTS idx_auth_action_tokens_user_id ON auth_action_tokens (user_id, purpose, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_auth_action_tokens_expires_at ON auth_action_tokens (expires_at);

CREATE TABLE IF NOT EXISTS auth_sessions (
  session_id UUID PRIMARY KEY,
  user_id UUID NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,
  family_id UUID NOT NULL,
  refresh_jti TEXT NOT NULL UNIQUE,
  refresh_token_hash TEXT NOT NULL UNIQUE,
  csrf_token_hash TEXT NOT NULL,
  fingerprint_hash TEXT NOT NULL,
  device_label TEXT NOT NULL DEFAULT '',
  ip INET NULL,
  user_agent TEXT NOT NULL DEFAULT '',
  mfa_verified BOOLEAN NOT NULL DEFAULT FALSE,
  replaced_by_session_id UUID NULL REFERENCES auth_sessions(session_id) ON DELETE SET NULL,
  last_seen_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  expires_at TIMESTAMPTZ NOT NULL,
  revoked_at TIMESTAMPTZ NULL,
  revoke_reason TEXT NOT NULL DEFAULT '',
  compromised_at TIMESTAMPTZ NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_auth_sessions_user_id ON auth_sessions (user_id, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_auth_sessions_family_id ON auth_sessions (family_id);
CREATE INDEX IF NOT EXISTS idx_auth_sessions_active ON auth_sessions (user_id, revoked_at, expires_at);

CREATE TABLE IF NOT EXISTS auth_backup_codes (
  backup_code_id UUID PRIMARY KEY,
  user_id UUID NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,
  code_hash TEXT NOT NULL,
  consumed_at TIMESTAMPTZ NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_auth_backup_codes_user_id ON auth_backup_codes (user_id, consumed_at);

CREATE TABLE IF NOT EXISTS auth_login_attempts (
  attempt_id UUID PRIMARY KEY,
  user_id UUID NULL REFERENCES users(user_id) ON DELETE SET NULL,
  email_normalized TEXT NOT NULL,
  ip INET NULL,
  user_agent TEXT NOT NULL DEFAULT '',
  outcome TEXT NOT NULL,
  failure_reason TEXT NOT NULL DEFAULT '',
  risk_score INTEGER NOT NULL DEFAULT 0,
  fingerprint_hash TEXT NOT NULL DEFAULT '',
  suspicious_reason TEXT NOT NULL DEFAULT '',
  created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_auth_login_attempts_email_created ON auth_login_attempts (email_normalized, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_auth_login_attempts_user_created ON auth_login_attempts (user_id, created_at DESC);

CREATE TABLE IF NOT EXISTS auth_security_events (
  event_id UUID PRIMARY KEY,
  user_id UUID NULL REFERENCES users(user_id) ON DELETE SET NULL,
  actor_user_id UUID NULL REFERENCES users(user_id) ON DELETE SET NULL,
  session_id UUID NULL REFERENCES auth_sessions(session_id) ON DELETE SET NULL,
  category TEXT NOT NULL,
  event_type TEXT NOT NULL,
  severity TEXT NOT NULL,
  ip INET NULL,
  user_agent TEXT NOT NULL DEFAULT '',
  payload JSONB NOT NULL DEFAULT '{}'::jsonb,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_auth_security_events_user_id ON auth_security_events (user_id, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_auth_security_events_category ON auth_security_events (category, created_at DESC);

-- +goose StatementBegin
CREATE OR REPLACE FUNCTION set_auth_updated_at() RETURNS TRIGGER AS $fn$
BEGIN
  NEW.updated_at = now();
  RETURN NEW;
END;
$fn$ LANGUAGE plpgsql;
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
