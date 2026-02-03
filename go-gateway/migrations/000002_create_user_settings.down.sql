-- 000002_create_user_settings.down.sql
-- 回滚用户设置表

DROP TRIGGER IF EXISTS update_user_settings_updated_at ON user_settings;
DROP TABLE IF EXISTS user_settings;
