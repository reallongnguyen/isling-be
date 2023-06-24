DROP INDEX IF EXISTS media_rooms_visibility_index;
DROP INDEX IF EXISTS media_rooms_owner_index;
DROP INDEX IF EXISTS media_rooms_slug_index;

DROP TRIGGER IF EXISTS auto_change_updated_at ON media_rooms;
DROP TABLE IF EXISTS media_rooms;

DROP TRIGGER IF EXISTS auto_change_updated_at ON play_users;
DROP TABLE IF EXISTS play_users;
DROP TYPE IF EXISTS visibility_type;
