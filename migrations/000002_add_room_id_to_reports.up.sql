-- 000002_add_room_id_to_reports.up.sql
ALTER TABLE reports ADD COLUMN IF NOT EXISTS room_id INTEGER;
CREATE INDEX IF NOT EXISTS idx_reports_room_id ON reports(room_id);
