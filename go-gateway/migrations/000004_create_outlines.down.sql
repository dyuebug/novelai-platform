-- 000004_create_outlines.down.sql
-- 回滚大纲表

DROP TRIGGER IF EXISTS update_outlines_updated_at ON outlines;
DROP TABLE IF EXISTS outlines;
