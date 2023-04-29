CREATE TABLE IF NOT EXISTS accounts (
    id serial PRIMARY KEY,
    username VARCHAR(255) UNIQUE,
    password VARCHAR(255),
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

CREATE OR REPLACE TRIGGER set_created_at
BEFORE UPDATE ON accounts
FOR EACH ROW
EXECUTE PROCEDURE trigger_set_updated_at();

CREATE INDEX IF NOT EXISTS accounts_username_index ON accounts USING hash (username);
