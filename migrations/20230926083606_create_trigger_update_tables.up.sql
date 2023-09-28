CREATE OR REPLACE FUNCTION table_update_notify() RETURNS trigger AS $$
DECLARE
  id bigint;
BEGIN
  IF TG_OP = 'INSERT' OR TG_OP = 'UPDATE' THEN
    id = NEW.id;
  ELSE
    id = OLD.id;
  END IF;
  PERFORM pg_notify('table_update', json_build_object('table', TG_TABLE_NAME, 'id', id, 'type', TG_OP, 'data', to_json(COALESCE(NEW, OLD))::text)::text);
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;


DROP TRIGGER IF EXISTS accounts_notify_update ON accounts;
CREATE TRIGGER accounts_notify_update AFTER UPDATE ON accounts FOR EACH ROW EXECUTE PROCEDURE table_update_notify();

DROP TRIGGER IF EXISTS accounts_notify_insert ON accounts;
CREATE TRIGGER accounts_notify_insert AFTER INSERT ON accounts FOR EACH ROW EXECUTE PROCEDURE table_update_notify();

DROP TRIGGER IF EXISTS accounts_notify_delete ON accounts;
CREATE TRIGGER accounts_notify_delete AFTER DELETE ON accounts FOR EACH ROW EXECUTE PROCEDURE table_update_notify();

DROP TRIGGER IF EXISTS profiles_notify_update ON profiles;
CREATE TRIGGER profiles_notify_update AFTER UPDATE ON profiles FOR EACH ROW EXECUTE PROCEDURE table_update_notify();

DROP TRIGGER IF EXISTS profiles_notify_insert ON profiles;
CREATE TRIGGER profiles_notify_insert AFTER INSERT ON profiles FOR EACH ROW EXECUTE PROCEDURE table_update_notify();

DROP TRIGGER IF EXISTS profiles_notify_delete ON profiles;
CREATE TRIGGER profiles_notify_delete AFTER DELETE ON profiles FOR EACH ROW EXECUTE PROCEDURE table_update_notify();

DROP TRIGGER IF EXISTS media_rooms_notify_update ON media_rooms;
CREATE TRIGGER media_rooms_notify_update AFTER UPDATE ON media_rooms FOR EACH ROW EXECUTE PROCEDURE table_update_notify();

DROP TRIGGER IF EXISTS media_rooms_notify_insert ON media_rooms;
CREATE TRIGGER media_rooms_notify_insert AFTER INSERT ON media_rooms FOR EACH ROW EXECUTE PROCEDURE table_update_notify();

DROP TRIGGER IF EXISTS media_rooms_notify_delete ON media_rooms;
CREATE TRIGGER media_rooms_notify_delete AFTER DELETE ON media_rooms FOR EACH ROW EXECUTE PROCEDURE table_update_notify();
