CREATE TABLE IF NOT EXISTS play_users (
  id SERIAL PRIMARY KEY,
  account_id INTEGER UNIQUE NOT NULL,
  recently_joined_rooms jsonb NOT NULL DEFAULT '[]'::jsonb,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL,
  updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL
);

CREATE OR REPLACE TRIGGER auto_change_updated_at
BEFORE UPDATE ON play_users
FOR EACH ROW
EXECUTE PROCEDURE trigger_set_updated_at();

---
CREATE TYPE visibility_type AS ENUM (
  'public',
  'member'
);

CREATE TABLE IF NOT EXISTS media_rooms (
  id SERIAL PRIMARY KEY,
  owner_id INTEGER NOT NULL,-- REFERENCES play_users (account_id) ON DELETE CASCADE,
  visibility visibility_type NOT NULL,
  invite_code VARCHAR(32) NULL,
  name VARCHAR(256) NOT NULL,
  slug VARCHAR(256) NOT NULL,
  description VARCHAR(512) NOT NULL,
  cover VARCHAR(256) NOT NULL,
  audience_count INT NOT NULL DEFAULT 0,
  audiences jsonb NOT NULL DEFAULT '[]'::jsonb,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL,
  updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL
);

CREATE OR REPLACE TRIGGER auto_change_updated_at
BEFORE UPDATE ON media_rooms
FOR EACH ROW
EXECUTE PROCEDURE trigger_set_updated_at();

CREATE INDEX IF NOT EXISTS media_rooms_visibility_index ON media_rooms USING hash (visibility);
CREATE INDEX IF NOT EXISTS media_rooms_owner_index ON media_rooms USING hash (owner_id);
CREATE UNIQUE INDEX IF NOT EXISTS media_rooms_slug_index ON media_rooms USING btree (slug);
