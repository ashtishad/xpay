CREATE TYPE wallet_status AS ENUM ('active', 'inactive', 'blocked');
CREATE TYPE wallet_currency AS ENUM ('USD');

CREATE TABLE IF NOT EXISTS wallets (
    id BIGSERIAL PRIMARY KEY,
    uuid UUID UNIQUE NOT NULL DEFAULT gen_random_uuid(),
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    balance BIGINT NOT NULL DEFAULT 0 CHECK (balance >= 0),
    currency wallet_currency NOT NULL DEFAULT 'USD',
    status wallet_status NOT NULL DEFAULT 'active',
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

ALTER TABLE wallets ADD CONSTRAINT unique_user_currency UNIQUE (user_id, currency);

CREATE TRIGGER update_wallet_updated_at_trigger
BEFORE UPDATE ON wallets
FOR EACH ROW
EXECUTE FUNCTION update_updated_at();
