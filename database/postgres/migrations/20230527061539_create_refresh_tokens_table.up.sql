CREATE TABLE IF NOT EXISTS refresh_tokens (
  id serial PRIMARY KEY,
  account_id INTEGER NOT NULL REFERENCES accounts (id) ON DELETE CASCADE,
  encrypted_token VARCHAR(256) NOT NULL,
  revoked BOOLEAN DEFAULT FALSE NOT NULL,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL
);

CREATE INDEX IF NOT EXISTS refresh_tokens_token ON refresh_tokens USING hash (encrypted_token);
CREATE INDEX IF NOT EXISTS refresh_tokens_account_id ON refresh_tokens USING hash (account_id);
