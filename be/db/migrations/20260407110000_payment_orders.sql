-- +goose Up
-- +goose StatementBegin
CREATE TABLE payment_orders(
  order_id uuid PRIMARY KEY,
  public_token text NOT NULL UNIQUE,
  user_id uuid NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,
  user_email text NOT NULL,
  user_display_name text NOT NULL DEFAULT '',
  provider text NOT NULL,
  plan_code text NOT NULL,
  base_amount bigint NOT NULL DEFAULT 0,
  amount bigint NOT NULL DEFAULT 0,
  tip_amount bigint NOT NULL DEFAULT 0,
  currency text NOT NULL DEFAULT 'RUB',
  subscription_type text NOT NULL DEFAULT 'one_time',
  status text NOT NULL DEFAULT 'pending',
  provider_status text NOT NULL DEFAULT '',
  provider_payment_id text NOT NULL DEFAULT '',
  provider_error_code text NOT NULL DEFAULT '',
  provider_message text NOT NULL DEFAULT '',
  provider_response jsonb NOT NULL DEFAULT '{}'::jsonb,
  provider_fee_percent double precision NOT NULL DEFAULT 0,
  fee_amount bigint NOT NULL DEFAULT 0,
  net_amount bigint NOT NULL DEFAULT 0,
  access_from timestamptz NULL,
  access_until timestamptz NULL,
  return_to text NOT NULL DEFAULT '',
  success_url text NOT NULL DEFAULT '',
  fail_url text NOT NULL DEFAULT '',
  payment_url text NOT NULL DEFAULT '',
  last_checked_at timestamptz NULL,
  notified_at timestamptz NULL,
  paid_at timestamptz NULL,
  failed_at timestamptz NULL,
  refunded_at timestamptz NULL,
  created_at timestamptz NOT NULL DEFAULT NOW(),
  updated_at timestamptz NOT NULL DEFAULT NOW(),
  CONSTRAINT payment_orders_provider_check CHECK (provider IN ('tbank')),
  CONSTRAINT payment_orders_plan_code_check CHECK (plan_code IN ('pro', 'donation')),
  CONSTRAINT payment_orders_subscription_type_check CHECK (subscription_type IN ('one_time')),
  CONSTRAINT payment_orders_status_check CHECK (status IN ('pending', 'success', 'failed', 'refunded', 'canceled')),
  CONSTRAINT payment_orders_amounts_check CHECK (base_amount >= 0 AND amount >= 0 AND tip_amount >= 0 AND fee_amount >= 0 AND net_amount >= 0 AND amount = base_amount + tip_amount),
  CONSTRAINT payment_orders_access_range_check CHECK (access_from IS NULL OR access_until IS NULL OR access_until > access_from)
);
CREATE INDEX idx_payment_orders_user_id ON payment_orders(user_id);
CREATE INDEX idx_payment_orders_status ON payment_orders(status);
CREATE INDEX idx_payment_orders_plan_code ON payment_orders(plan_code);
CREATE INDEX idx_payment_orders_provider ON payment_orders(provider);
CREATE INDEX idx_payment_orders_provider_payment_id ON payment_orders(provider_payment_id);
CREATE INDEX idx_payment_orders_user_id_created_at_desc ON payment_orders(user_id, created_at DESC);
CREATE INDEX idx_payment_orders_user_id_access_until_desc ON payment_orders(user_id, access_until DESC);
CREATE INDEX idx_payment_orders_created_at_desc ON payment_orders(created_at DESC);
CREATE INDEX idx_payment_orders_paid_at_desc ON payment_orders(paid_at DESC);
CREATE INDEX idx_payment_orders_public_token ON payment_orders(public_token);
CREATE INDEX idx_payment_orders_user_email_trgm ON payment_orders USING GIN(lower(user_email) gin_trgm_ops);
CREATE INDEX idx_payment_orders_user_display_name_trgm ON payment_orders USING GIN(lower(user_display_name) gin_trgm_ops);
CREATE FUNCTION set_payment_orders_updated_at()
  RETURNS TRIGGER
  AS $$
BEGIN
  NEW.updated_at := NOW();
  RETURN NEW;
END;
$$
LANGUAGE plpgsql;
CREATE TRIGGER trg_payment_orders_updated_at
  BEFORE UPDATE ON payment_orders
  FOR EACH ROW
  EXECUTE FUNCTION set_payment_orders_updated_at();
-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TRIGGER trg_payment_orders_updated_at ON payment_orders;
DROP FUNCTION set_payment_orders_updated_at();
DROP INDEX idx_payment_orders_user_display_name_trgm;
DROP INDEX idx_payment_orders_user_email_trgm;
DROP INDEX idx_payment_orders_public_token;
DROP INDEX idx_payment_orders_paid_at_desc;
DROP INDEX idx_payment_orders_created_at_desc;
DROP INDEX idx_payment_orders_user_id_access_until_desc;
DROP INDEX idx_payment_orders_user_id_created_at_desc;
DROP INDEX idx_payment_orders_provider_payment_id;
DROP INDEX idx_payment_orders_provider;
DROP INDEX idx_payment_orders_plan_code;
DROP INDEX idx_payment_orders_status;
DROP INDEX idx_payment_orders_user_id;
DROP TABLE payment_orders;
-- +goose StatementEnd
