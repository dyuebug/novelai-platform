-- 000006_create_embeddings.up.sql
-- 创建向量嵌入表 (使用 pgvector 扩展)

-- 启用 pgvector 扩展
CREATE EXTENSION IF NOT EXISTS vector;

-- 向量嵌入表
CREATE TABLE embeddings (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    project_id UUID NOT NULL REFERENCES projects(id) ON DELETE CASCADE,

    -- 内容类型: settings, summaries, chunks, characters
    content_type VARCHAR(50) NOT NULL,
    -- 源内容 ID (章节ID、角色ID等)
    content_id UUID NOT NULL,

    -- 原始内容 (用于返回检索结果)
    content TEXT NOT NULL,
    -- 章节编号 (仅 chunks/summaries 类型使用)
    chapter_number INTEGER,

    -- 向量嵌入 (1536 维度，适配 text-embedding-3-small)
    embedding vector(1536) NOT NULL,

    -- 元数据
    metadata JSONB DEFAULT '{}',

    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- 索引
CREATE INDEX idx_embeddings_project ON embeddings(project_id);
CREATE INDEX idx_embeddings_type ON embeddings(content_type);
CREATE INDEX idx_embeddings_content_id ON embeddings(content_id);
CREATE INDEX idx_embeddings_chapter ON embeddings(chapter_number) WHERE chapter_number IS NOT NULL;

-- 向量相似度索引 (IVFFlat，适合中等规模数据)
CREATE INDEX idx_embeddings_vector ON embeddings USING ivfflat (embedding vector_cosine_ops) WITH (lists = 100);

-- 复合索引：按项目和类型检索
CREATE INDEX idx_embeddings_project_type ON embeddings(project_id, content_type);

-- 更新时间触发器
CREATE TRIGGER update_embeddings_updated_at
    BEFORE UPDATE ON embeddings
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();
