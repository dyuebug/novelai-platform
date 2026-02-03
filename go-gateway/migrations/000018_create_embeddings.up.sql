-- 000018_create_embeddings.up.sql
-- 创建向量嵌入表 (需要 pgvector 扩展)

-- 启用 pgvector 扩展
CREATE EXTENSION IF NOT EXISTS vector;

CREATE TABLE embeddings (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    project_id UUID NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
    content_type VARCHAR(50) NOT NULL,
    content_id UUID NOT NULL,
    content TEXT NOT NULL,
    embedding vector(1536),
    chapter_number INTEGER,
    metadata JSONB DEFAULT '{}',
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    UNIQUE(project_id, content_type, content_id)
);

CREATE INDEX idx_embeddings_project ON embeddings(project_id);
CREATE INDEX idx_embeddings_content_type ON embeddings(content_type);
CREATE INDEX idx_embeddings_chapter ON embeddings(chapter_number);

-- 创建向量索引 (IVFFlat 用于近似最近邻搜索)
CREATE INDEX idx_embeddings_vector ON embeddings USING ivfflat (embedding vector_cosine_ops) WITH (lists = 100);
