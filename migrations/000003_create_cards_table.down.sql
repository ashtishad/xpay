ALTER TABLE IF EXISTS cards DROP CONSTRAINT IF EXISTS check_expiry_date_not_past;

DROP TRIGGER IF EXISTS update_card_updated_at_trigger ON cards;

DROP INDEX IF EXISTS idx_cards_user_provider_type;
DROP INDEX IF EXISTS idx_cards_wallet_id;
DROP INDEX IF EXISTS idx_cards_user_id;

DROP TABLE IF EXISTS cards;

DROP TYPE IF EXISTS card_status;
DROP TYPE IF EXISTS card_type;
DROP TYPE IF EXISTS card_provider;
