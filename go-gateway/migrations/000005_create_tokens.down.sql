-- 000005_create_tokens.down.sql
-- 回滚令牌表

DROP TABLE IF EXISTS password_reset_tokens;
DROP TABLE IF EXISTS refresh_tokens;
