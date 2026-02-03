-- 000006_create_embeddings.down.sql
-- 删除向量嵌入表

DROP TRIGGER IF EXISTS update_embeddings_updated_at ON embeddings;
DROP TABLE IF EXISTS embeddings;
-- 注意: 不删除 vector 扩展，因为可能被其他表使用
