CREATE TABLE IF NOT EXISTS accounts (
    id serial PRIMARY KEY,
    email VARCHAR(255),
    encrypted_password VARCHAR(255),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL
);

CREATE OR REPLACE FUNCTION trigger_set_updated_at()
RETURNS TRIGGER AS $$
BEGIN
  NEW.updated_at = NOW();
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE TRIGGER auto_change_updated_at
BEFORE UPDATE ON accounts
FOR EACH ROW
EXECUTE PROCEDURE trigger_set_updated_at();

CREATE UNIQUE INDEX IF NOT EXISTS accounts_email_index ON accounts USING btree (email);
