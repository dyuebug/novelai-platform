-- 000007_create_chapter_versions.up.sql
-- 创建章节版本表

CREATE TABLE chapter_versions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    chapter_id UUID NOT NULL REFERENCES chapters(id) ON DELETE CASCADE,
    version_number INTEGER NOT NULL,
    content TEXT NOT NULL,
    word_count INTEGER DEFAULT 0,
    source VARCHAR(50) DEFAULT 'manual',
    ai_provider VARCHAR(50),
    ai_model VARCHAR(100),
    prompt_used TEXT,
    created_at TIMESTAMP DEFAULT NOW(),
    created_by UUID REFERENCES users(id) ON DELETE SET NULL,
    UNIQUE(chapter_id, version_number)
);

CREATE INDEX idx_chapter_versions_chapter ON chapter_versions(chapter_id);
CREATE INDEX idx_chapter_versions_number ON chapter_versions(chapter_id, version_number);
