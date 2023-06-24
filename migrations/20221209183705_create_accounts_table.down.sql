DROP TRIGGER IF EXISTS auto_change_updated_at ON accounts;
DROP FUNCTION IF EXISTS trigger_set_updated_at;
DROP INDEX IF EXISTS accounts_email_index;

DROP TABLE IF EXISTS accounts;
