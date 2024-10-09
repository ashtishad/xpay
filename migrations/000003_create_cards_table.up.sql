CREATE TYPE card_provider AS ENUM ('visa', 'mastercard', 'amex');
CREATE TYPE card_type AS ENUM ('credit', 'debit');
CREATE TYPE card_status AS ENUM ('active', 'inactive', 'deleted');

CREATE TABLE IF NOT EXISTS cards (
    id BIGSERIAL PRIMARY KEY,
    uuid UUID UNIQUE NOT NULL DEFAULT gen_random_uuid(),
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    wallet_id BIGINT NOT NULL REFERENCES wallets(id) ON DELETE CASCADE,
    encrypted_card_number BYTEA NOT NULL,
    provider card_provider NOT NULL,
    type card_type NOT NULL,
    last_four VARCHAR(4) NOT NULL,
    expiry_date DATE NOT NULL,
    status card_status NOT NULL DEFAULT 'active',
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_cards_user_id ON cards(user_id);
CREATE INDEX idx_cards_wallet_id ON cards(wallet_id);

-- Ensure a user can only have one active card per provider and type combination (Partial index and Unique constraint)
CREATE UNIQUE INDEX idx_cards_user_provider_type ON cards(user_id, provider, type) WHERE status = 'active';

-- Add check constraint to ensure expiry_date is not in the past
ALTER TABLE cards ADD CONSTRAINT check_expiry_date_not_past
    CHECK (expiry_date >= CURRENT_DATE);

-- Create a trigger to automatically update the updated_at column
CREATE TRIGGER update_card_updated_at_trigger
BEFORE UPDATE ON cards
FOR EACH ROW
EXECUTE FUNCTION update_updated_at();
