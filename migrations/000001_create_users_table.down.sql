DROP TRIGGER IF EXISTS trigger_users_update_updated_at ON users;
DROP TABLE IF EXISTS users;
DROP FUNCTION IF EXISTS update_updated_at();
DROP TYPE IF EXISTS user_status;
DROP TYPE IF EXISTS user_role;
